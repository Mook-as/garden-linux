package lifecycle_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cloudfoundry-incubator/garden"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var (
	dockerRegistryRootFSPath   = os.Getenv("GARDEN_DOCKER_REGISTRY_TEST_ROOTFS")
	dockerRegistryV2RootFSPath = os.Getenv("GARDEN_DOCKER_REGISTRY_V2_TEST_ROOTFS")
)

var _ = Describe("Rootfs container create parameter", func() {
	var container garden.Container
	var args []string

	BeforeEach(func() {
		container = nil
		args = []string{}
	})

	JustBeforeEach(func() {
		client = startGarden(args...)
	})

	AfterEach(func() {
		if container != nil {
			Expect(client.Destroy(container.Handle())).To(Succeed())
		}
	})

	Context("without a default rootfs", func() {
		BeforeEach(func() {
			args = []string{"--rootfs", ""}
		})

		It("without a rootfs in container spec, the container creation fails", func() {
			var err error

			container, err = client.Create(garden.ContainerSpec{RootFSPath: ""})
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(MatchError(ContainSubstring(
				"RootFSPath: is a required parameter, since no default rootfs was provided to the server. To provide a default rootfs, use the --rootfs flag on startup.",
			)))
		})

		It("with a rootfs in container spec, the container is created successfully", func() {
			var err error

			container, err = client.Create(garden.ContainerSpec{RootFSPath: os.Getenv("GARDEN_TEST_ROOTFS")})
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("with a default rootfs", func() {
		It("the container is created successfully", func() {
			var err error

			container, err = client.Create(garden.ContainerSpec{RootFSPath: ""})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("with a docker rootfs URI", func() {
		Context("not containing a host", func() {
			It("succesfully creates the container", func() {
				var err error

				container, err = client.Create(garden.ContainerSpec{RootFSPath: "docker:///busybox"})
				Expect(err).ToNot(HaveOccurred())
			})

			Context("when image does not exist", func() {
				It("should not leak the depot directory", func() {
					_, err := client.Create(
						garden.ContainerSpec{
							RootFSPath: "docker:///cloudfoundry/doesnotexist",
						},
					)
					Expect(err).To(HaveOccurred())

					entries, err := ioutil.ReadDir(client.DepotPath)
					Expect(err).ToNot(HaveOccurred())
					Expect(entries).To(HaveLen(0))
				})
			})

			Context("when the -registry flag targets a non-existing registry", func() {
				BeforeEach(func() {
					args = []string{"--registry", "registry-12.banana-docker.io"}
				})

				It("should fail to create a container", func() {
					var err error

					container, err = client.Create(garden.ContainerSpec{RootFSPath: "docker:///busybox"})
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("containing a host", func() {
			Context("which is valid", func() {
				It("creates the container successfully", func() {
					var err error

					container, err = client.Create(garden.ContainerSpec{RootFSPath: "docker://registry-1.docker.io/busybox"})
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("which is invalid", func() {
				It("the container is not created successfully", func() {
					var err error
					container, err = client.Create(garden.ContainerSpec{RootFSPath: "docker://xindex.docker.io/busybox"})
					Expect(err).To(HaveOccurred())
				})
			})

			Context("which is insecure", func() {
				var (
					dockerRegistry     garden.Container
					dockerRegistryIP   string
					dockerRegistryPort string
				)

				BeforeEach(func() {
					dockerRegistryIP = "10.0.0.2"
					dockerRegistryPort = "5000"
				})

				JustBeforeEach(func() {
					if dockerRegistryV2RootFSPath == "" {
						Skip("GARDEN_DOCKER_REGISTRY_V2_TEST_ROOTFS undefined")
					}

					dockerRegistry = startV2DockerRegistry(dockerRegistryIP, dockerRegistryPort)
				})

				AfterEach(func() {
					if dockerRegistry != nil {
						Expect(client.Destroy(dockerRegistry.Handle())).To(Succeed())
					}
				})

				Context("when the host is listed in -insecureDockerRegistry", func() {
					BeforeEach(func() {
						args = []string{
							"-allowHostAccess=true",
						}
					})

					Context("when the registry is NOT using TLS", func() {
						BeforeEach(func() {
							args = append(
								args,
								"-insecureDockerRegistry",
								fmt.Sprintf("%s:%s", dockerRegistryIP, dockerRegistryPort),
							)
						})

						It("creates the container successfully ", func() {
							_, err := client.Create(garden.ContainerSpec{
								RootFSPath: fmt.Sprintf("docker://%s:%s/busybox", dockerRegistryIP,
									dockerRegistryPort),
							})
							Expect(err).ToNot(HaveOccurred())
						})
					})

					Context("when the registry is in a CIDR", func() {
						BeforeEach(func() {
							args = append(
								args,
								"-insecureDockerRegistry",
								fmt.Sprintf("%s/24", dockerRegistryIP),
							)
						})

						It("creates the container successfully ", func() {
							_, err := client.Create(garden.ContainerSpec{
								RootFSPath: fmt.Sprintf("docker://%s:%s/busybox", dockerRegistryIP, dockerRegistryPort),
							})
							Expect(err).ToNot(HaveOccurred())
						})
					})

					Context("when the registry is using TLS", func() {
						var server *httptest.Server
						var serverURL *url.URL

						BeforeEach(func() {
							proxyTo, err := url.Parse(fmt.Sprintf("http://%s:%s", dockerRegistryIP,
								dockerRegistryPort))
							Expect(err).NotTo(HaveOccurred())

							server = httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(proxyTo))
							serverURL, err = url.Parse(server.URL)
							Expect(err).NotTo(HaveOccurred())

							args = append(
								args,
								"-insecureDockerRegistry",
								serverURL.Host,
							)
						})

						AfterEach(func() {
							server.Close()
						})

						It("creates the container successfully", func() {
							_, err := client.Create(garden.ContainerSpec{
								RootFSPath: fmt.Sprintf("docker://%s/busybox", serverURL.Host),
							})
							Expect(err).ToNot(HaveOccurred())
						})

						Context("and its specified as --registry", func() {
							BeforeEach(func() {
								args = append(args, "--registry", serverURL.Host)
							})

							It("still works when the host is specified", func() {
								_, err := client.Create(garden.ContainerSpec{
									RootFSPath: fmt.Sprintf("docker://%s/busybox", serverURL.Host),
								})
								Expect(err).ToNot(HaveOccurred())
							})

							It("still works using the default host", func() {
								_, err := client.Create(garden.ContainerSpec{
									RootFSPath: fmt.Sprintf("docker:///busybox"),
								})
								Expect(err).ToNot(HaveOccurred())
							})
						})
					})
				})

				Context("when the host is NOT listed in -insecureDockerRegistry", func() {
					It("fails", func() {
						_, err := client.Create(garden.ContainerSpec{
							RootFSPath: fmt.Sprintf("docker://%s:%s/busybox", dockerRegistryIP,
								dockerRegistryPort),
						})

						Expect(err).To(HaveOccurred())
					})
				})
			})
		})
	})

	Context("when the modified timestamp of the rootfs top-level directory changes", func() {
		var (
			rootfspath          string
			privilegedContainer bool
			container2          garden.Container
		)

		JustBeforeEach(func() {
			var err error
			rootfspath = createSmallRootfs()

			container, err = client.Create(garden.ContainerSpec{
				RootFSPath: rootfspath,
				Privileged: privilegedContainer,
			})
			Expect(err).NotTo(HaveOccurred())

			// ls is convenient, but any file modification is sufficient
			ls := filepath.Join(rootfspath, "bin", "ls")
			Expect(exec.Command("cp", ls, rootfspath).Run()).To(Succeed())

			container2, err = client.Create(garden.ContainerSpec{
				RootFSPath: rootfspath,
				Privileged: privilegedContainer,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if container2 != nil {
				Expect(client.Destroy(container2.Handle())).To(Succeed())
			}
		})

		Context("with a non-privileged container", func() {
			BeforeEach(func() {
				privilegedContainer = false
			})

			It("should use the updated rootfs when creating a new container", func() {
				process, err := container2.Run(garden.ProcessSpec{
					Path: "/ls",
					User: "root",
				}, garden.ProcessIO{Stdout: GinkgoWriter, Stderr: GinkgoWriter})
				Expect(err).NotTo(HaveOccurred())

				exitStatus, err := process.Wait()
				Expect(err).NotTo(HaveOccurred())
				Expect(exitStatus).To(Equal(0))
			})
		})
	})

})

func startV2DockerRegistry(dockerRegistryIP string, dockerRegistryPort string) garden.Container {
	dockerRegistry, err := client.Create(
		garden.ContainerSpec{
			RootFSPath: dockerRegistryV2RootFSPath,
			Network:    dockerRegistryIP,
		},
	)
	Expect(err).ToNot(HaveOccurred())

	_, err = dockerRegistry.Run(garden.ProcessSpec{
		User: "root",
		Env: []string{
			"REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY=/opt/docker-registry",
		},
		Path: "/go/bin/registry",
		Args: []string{"/go/src/github.com/docker/distribution/cmd/registry/config.yml"},
	}, garden.ProcessIO{Stdout: GinkgoWriter, Stderr: GinkgoWriter})
	Expect(err).ToNot(HaveOccurred())

	Eventually(
		fmt.Sprintf("http://%s:%s/v2/", dockerRegistryIP, dockerRegistryPort),
		"60s",
	).Should(RespondToGETWith(200))

	return dockerRegistry
}

type statusMatcher struct {
	expectedStatus int

	httpError    error
	actualStatus int
}

func RespondToGETWith(expected int) types.GomegaMatcher {
	return &statusMatcher{expected, nil, 200}
}

func (m *statusMatcher) Match(actual interface{}) (success bool, err error) {
	response, err := http.Get(fmt.Sprintf("%s", actual))
	if err != nil {
		m.httpError = err
		return false, nil
	}

	m.httpError = nil
	m.actualStatus = response.StatusCode
	return response.StatusCode == m.expectedStatus, nil
}

func (m *statusMatcher) FailureMessage(actual interface{}) string {
	if m.httpError != nil {
		return fmt.Sprintf("Expected http request to have status %d but got error: %s", m.expectedStatus, m.httpError.Error())
	}

	return fmt.Sprintf("Expected http status code to be %d but was %d", m.expectedStatus, m.actualStatus)
}

func (m *statusMatcher) NegatedFailureMessage(actual interface{}) string {
	if m.httpError != nil {
		return fmt.Sprintf("Expected http request to have status %d, but got error: %s", m.expectedStatus, m.httpError.Error())
	}

	return fmt.Sprintf("Expected http status code not to be %d", m.expectedStatus)
}

func createSmallRootfs() string {
	rootfs := os.Getenv("GARDEN_PREEXISTING_USERS_TEST_ROOTFS")
	if rootfs == "" {
		Skip("pre-existing users rootfs not found")
	}

	rootfspath, err := ioutil.TempDir("", "rootfs-cache-invalidation")
	Expect(err).NotTo(HaveOccurred())
	cmd := exec.Command("cp", "-rf", rootfs, rootfspath)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	Expect(cmd.Run()).To(Succeed())
	return filepath.Join(rootfspath, filepath.Base(rootfs))
}
