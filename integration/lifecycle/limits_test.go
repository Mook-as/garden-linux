package lifecycle_test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/garden-linux/integration/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Limits", func() {
	const BTRFS_WAIT_TIME = 90

	var container garden.Container
	var startGardenArgs []string

	var privilegedContainer bool
	var rootfs string

	JustBeforeEach(func() {
		var err error
		client = startGarden(startGardenArgs...)
		container, err = client.Create(garden.ContainerSpec{Privileged: privilegedContainer, RootFSPath: rootfs})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if container != nil {
			Expect(client.Destroy(container.Handle())).To(Succeed())
		}
	})

	BeforeEach(func() {
		privilegedContainer = false
		rootfs = ""
		startGardenArgs = []string{}
	})

	Describe("LimitMemory", func() {
		Context("with a memory limit", func() {
			JustBeforeEach(func() {
				err := container.LimitMemory(garden.MemoryLimits{
					LimitInBytes: 64 * 1024 * 1024,
				})
				Expect(err).ToNot(HaveOccurred())
			})

			Context("when the process writes too much to /dev/shm", func() {
				It("is killed", func() {
					process, err := container.Run(garden.ProcessSpec{
						User: "vcap",
						Path: "dd",
						Args: []string{"if=/dev/urandom", "of=/dev/shm/too-big", "bs=1M", "count=65"},
					}, garden.ProcessIO{})
					Expect(err).ToNot(HaveOccurred())

					Expect(process.Wait()).ToNot(Equal(0))
				})
			})
		})
	})

	Describe("LimitDisk", func() {
		Context("with quotas disabled", func() {
			BeforeEach(func() {
				startGardenArgs = []string{"-disableQuotas=true"}
				rootfs = runner.RootFSPath
				privilegedContainer = true
			})

			Context("and there is a disk limit", func() {
				quotaLimit := garden.DiskLimits{
					ByteSoft: 5 * 1024 * 1024,
					ByteHard: 5 * 1024 * 1024,
				}

				JustBeforeEach(func() {
					Expect(container.LimitDisk(quotaLimit)).To(Succeed())
				})

				It("reports the disk limit size of the container as zero", func() {
					limit, err := container.CurrentDiskLimits()
					Expect(err).ToNot(HaveOccurred())
					Expect(limit).To(Equal(garden.DiskLimits{}))
				})

				Context("and it runs a process that exceeds the limit", func() {
					It("does not kill the process", func() {
						dd, err := container.Run(garden.ProcessSpec{
							User: "vcap",
							Path: "dd",
							Args: []string{"if=/dev/zero", "of=/tmp/some-file", "bs=1M", "count=6"},
						}, garden.ProcessIO{})
						Expect(err).ToNot(HaveOccurred())
						Expect(dd.Wait()).To(Equal(0))
					})
				})
			})
		})

		Context("when quotas are enabled and there is a disk limit", func() {
			var quotaLimit garden.DiskLimits

			if os.Getenv("BTRFS_SUPPORTED") == "" {
				log.Println("btrfs not supported: skipping disk limit tests")
				return
			}

			BeforeEach(func() {
				privilegedContainer = false
				rootfs = runner.RootFSPath
				quotaLimit = garden.DiskLimits{
					ByteSoft: 180 * 1024 * 1024,
					ByteHard: 180 * 1024 * 1024,
				}
			})

			JustBeforeEach(func() {
				err := container.LimitDisk(quotaLimit)
				Expect(err).ToNot(HaveOccurred())
			})

			Context("on a directory rootfs container", func() {
				BeforeEach(func() {
					privilegedContainer = false
					rootfs = runner.RootFSPath
				})

				It("reports correct disk usage", func() {
					var diskUsage uint64
					stdout := gbytes.NewBuffer()

					process, err := container.Run(garden.ProcessSpec{
						User: "vcap",
						Path: "sh",
						Args: []string{"-c", "du -d 0 / | awk ' {print $1 }'"},
					}, garden.ProcessIO{Stdout: stdout})
					Expect(err).ToNot(HaveOccurred())
					Expect(process.Wait()).To(Equal(0))

					_, err = fmt.Sscanf(strings.TrimSpace(string(stdout.Contents())), "%d", &diskUsage)
					Expect(err).ToNot(HaveOccurred())

					process, err = container.Run(garden.ProcessSpec{
						User: "vcap",
						Path: "dd",
						Args: []string{"if=/dev/urandom", "of=/home/vcap/some-file", "bs=1M", "count=10"},
					}, garden.ProcessIO{})
					Expect(err).ToNot(HaveOccurred())
					Expect(process.Wait()).To(Equal(0))

					metrics := func() uint64 {
						metricsAfter, err := container.Metrics()
						Expect(err).ToNot(HaveOccurred())

						return metricsAfter.DiskStat.BytesUsed
					}

					expectedBytes := (diskUsage * 1024) + uint64(10*1024*1024)
					Eventually(metrics, BTRFS_WAIT_TIME, 30).Should(BeNumerically("~", expectedBytes, 1269760))

					process, err = container.Run(garden.ProcessSpec{
						User: "vcap",
						Path: "dd",
						Args: []string{"if=/dev/urandom", "of=/home/vcap/another-file", "bs=1M", "count=10"},
					}, garden.ProcessIO{})
					Expect(err).ToNot(HaveOccurred())
					Expect(process.Wait()).To(Equal(0))

					expectedBytes = (diskUsage * 1024) + uint64(20*1024*1024)
					Eventually(metrics, BTRFS_WAIT_TIME, 30).Should(BeNumerically("~", expectedBytes, 1269760))
				})

				Context("on a rootfs with pre-existing users", func() {
					if os.Getenv("GARDEN_PREEXISTING_USERS_TEST_ROOTFS") == "" {
						log.Println("pre-existing users rootfs not found: skipping some limit disk tests")
						return
					}

					BeforeEach(func() {
						rootfs = os.Getenv("GARDEN_PREEXISTING_USERS_TEST_ROOTFS")
						quotaLimit = garden.DiskLimits{
							ByteSoft: 10 * 1024 * 1024,
							ByteHard: 10 * 1024 * 1024,
						}
					})

					Context("and run a process that exceeds the quota as bob", func() {
						It("kills the process", func() {
							dd, err := container.Run(garden.ProcessSpec{
								User: "bob",
								Path: "dd",
								Args: []string{"if=/dev/zero", "of=/home/bob/test", "bs=1M", "count=11"},
							}, garden.ProcessIO{})
							Expect(err).ToNot(HaveOccurred())
							Expect(dd.Wait()).ToNot(Equal(0))
						})
					})

					Context("and run a process that exceeds the quota as alice", func() {
						It("kills the process", func() {
							dd, err := container.Run(garden.ProcessSpec{
								User: "alice",
								Path: "dd",
								Args: []string{"if=/dev/zero", "of=/home/alice/test", "bs=1M", "count=11"},
							}, garden.ProcessIO{})
							Expect(err).ToNot(HaveOccurred())
							Expect(dd.Wait()).ToNot(Equal(0))
						})
					})
				})
			})

			Context("on a Docker container", func() {
				BeforeEach(func() {
					privilegedContainer = false
					rootfs = "docker:///busybox"
					quotaLimit = garden.DiskLimits{
						ByteSoft: 10 * 1024 * 1024,
						ByteHard: 10 * 1024 * 1024,
					}
				})

				It("reports the correct disk limit size of the container", func() {
					limit, err := container.CurrentDiskLimits()
					Expect(err).ToNot(HaveOccurred())
					Expect(limit).To(Equal(quotaLimit))
				})

				Context("and run a process that exceeds the quota", func() {
					It("kills the process", func() {
						dd, err := container.Run(garden.ProcessSpec{
							User: "vcap",
							Path: "dd",
							Args: []string{"if=/dev/zero", "of=/root/test", "bs=1M", "count=11"},
						}, garden.ProcessIO{})
						Expect(err).ToNot(HaveOccurred())
						Expect(dd.Wait()).ToNot(Equal(0))
					})
				})

				Context("on a rootfs with pre-existing users", func() {
					BeforeEach(func() {
						rootfs = "docker:///cloudfoundry/preexisting_users"
					})

					Context("and run a process that exceeds the quota as bob", func() {
						BeforeEach(func() {
							quotaLimit = garden.DiskLimits{
								ByteSoft: 10 * 1024 * 1024,
								ByteHard: 10 * 1024 * 1024,
							}
						})

						It("kills the process", func() {
							dd, err := container.Run(garden.ProcessSpec{
								User: "bob",
								Path: "dd",
								Args: []string{"if=/dev/zero", "of=/home/bob/test", "bs=1M", "count=11"},
							}, garden.ProcessIO{})
							Expect(err).ToNot(HaveOccurred())
							Expect(dd.Wait()).ToNot(Equal(0))
						})
					})

					Context("and run a process that exceeds the quota as alice", func() {
						BeforeEach(func() {
							quotaLimit = garden.DiskLimits{
								ByteSoft: 10 * 1024 * 1024,
								ByteHard: 10 * 1024 * 1024,
							}
						})

						It("kills the process", func() {
							dd, err := container.Run(garden.ProcessSpec{
								User: "alice",
								Path: "dd",
								Args: []string{"if=/dev/zero", "of=/home/alice/test", "bs=1M", "count=11"},
							}, garden.ProcessIO{})
							Expect(err).ToNot(HaveOccurred())
							Expect(dd.Wait()).ToNot(Equal(0))
						})
					})

					Context("user alice is getting near the set limit", func() {
						JustBeforeEach(func() {
							metrics := func() uint64 {
								metricsAfter, err := container.Metrics()
								Expect(err).ToNot(HaveOccurred())

								return metricsAfter.DiskStat.BytesUsed
							}

							Eventually(metrics, BTRFS_WAIT_TIME, 30).Should(BeNumerically("~", uint64(10*1024*1024), 1024*1024))

							bytesUsed := metrics()

							quotaLimit = garden.DiskLimits{
								ByteSoft: 10*1024*1024 + bytesUsed,
								ByteHard: 10*1024*1024 + bytesUsed,
							}

							err := container.LimitDisk(quotaLimit)
							Expect(err).ToNot(HaveOccurred())

							dd, err := container.Run(garden.ProcessSpec{
								User: "alice",
								Path: "dd",
								Args: []string{"if=/dev/zero", "of=/home/alice/test", "bs=1M", "count=8"},
							}, garden.ProcessIO{
								Stderr: GinkgoWriter,
								Stdout: GinkgoWriter,
							})
							Expect(err).ToNot(HaveOccurred())
							Expect(dd.Wait()).To(Equal(0))
						})

						It("kills the process if user bob tries to exceed the shared limit", func() {
							dd, err := container.Run(garden.ProcessSpec{
								User: "bob",
								Path: "dd",
								Args: []string{"if=/dev/zero", "of=/home/bob/test", "bs=1M", "count=3"},
							}, garden.ProcessIO{
								Stderr: GinkgoWriter,
								Stdout: GinkgoWriter,
							})
							Expect(err).ToNot(HaveOccurred())
							Expect(dd.Wait()).ToNot(Equal(0))
						})
					})
				})

				Context("that is privileged", func() {
					BeforeEach(func() {
						privilegedContainer = true
					})

					Context("and run a process that exceeds the quota as root", func() {
						It("kills the process", func() {
							dd, err := container.Run(garden.ProcessSpec{
								User: "root",
								Path: "dd",
								Args: []string{"if=/dev/zero", "of=/root/test", "bs=1M", "count=11"},
							}, garden.ProcessIO{})
							Expect(err).ToNot(HaveOccurred())
							Expect(dd.Wait()).ToNot(Equal(0))
						})
					})

					Context("and run a process that exceeds the quota as a new user", func() {
						It("kills the process", func() {
							addUser, err := container.Run(garden.ProcessSpec{
								User: "root",
								Path: "adduser",
								Args: []string{"-D", "-g", "", "bob"},
							}, garden.ProcessIO{})
							Expect(err).ToNot(HaveOccurred())
							Expect(addUser.Wait()).To(Equal(0))

							dd, err := container.Run(garden.ProcessSpec{
								User: "bob",
								Path: "dd",
								Args: []string{"if=/dev/zero", "of=/home/bob/test", "bs=1M", "count=11"},
							}, garden.ProcessIO{})
							Expect(err).ToNot(HaveOccurred())
							Expect(dd.Wait()).ToNot(Equal(0))
						})
					})
				})
			})

			Context("when multiple containers are created for the same user", func() {
				var container2 garden.Container
				var err error

				BeforeEach(func() {
					privilegedContainer = false
					rootfs = runner.RootFSPath
					quotaLimit = garden.DiskLimits{
						ByteSoft: 50 * 1024 * 1024,
						ByteHard: 50 * 1024 * 1024,
					}
				})

				JustBeforeEach(func() {
					container2, err = client.Create(garden.ContainerSpec{
						Privileged: privilegedContainer,
						RootFSPath: rootfs,
					})
					Expect(err).ToNot(HaveOccurred())
				})

				AfterEach(func() {
					if container2 != nil {
						Expect(client.Destroy(container2.Handle())).To(Succeed())
					}
				})

				It("gives each container its own quota", func() {
					process, err := container.Run(garden.ProcessSpec{
						User: "vcap",
						Path: "dd",
						Args: []string{"if=/dev/urandom", "of=/home/vcap/some-file", "bs=1M", "count=40"},
					}, garden.ProcessIO{})
					Expect(err).ToNot(HaveOccurred())
					Expect(process.Wait()).To(Equal(0))

					process, err = container2.Run(garden.ProcessSpec{
						User: "vcap",
						Path: "dd",
						Args: []string{"if=/dev/urandom", "of=/home/vcap/some-file", "bs=1M", "count=40"},
					}, garden.ProcessIO{})
					Expect(err).ToNot(HaveOccurred())
					Expect(process.Wait()).To(Equal(0))
				})
			})
		})
	})
})
