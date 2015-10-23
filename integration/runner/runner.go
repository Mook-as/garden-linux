package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/cloudfoundry-incubator/garden/client"
	"github.com/cloudfoundry-incubator/garden/client/connection"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var RootFSPath = os.Getenv("GARDEN_TEST_ROOTFS")
var GraphRoot = os.Getenv("GARDEN_TEST_GRAPHPATH")
var BinPath = "../../linux_backend/bin"
var GardenBin = "../../out/garden-linux"

type RunningGarden struct {
	client.Client
	process ifrit.Process
	runner  *ginkgomon.Runner

	Pid int

	tmpdir        string
	GraphRoot     string
	GraphPath     string
	DepotPath     string
	SnapshotsPath string

	logger lager.Logger
}

func Start(argv ...string) *RunningGarden {
	gardenAddr := fmt.Sprintf("/tmp/garden_%d.sock", GinkgoParallelNode())
	return start("unix", gardenAddr, argv...)
}

func start(network, addr string, argv ...string) *RunningGarden {
	tmpDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("test-garden-%d", ginkgo.GinkgoParallelNode()),
	)
	Expect(os.MkdirAll(tmpDir, 0755)).To(Succeed())

	if GraphRoot == "" {
		GraphRoot = filepath.Join(tmpDir, "graph")
	}
	graphPath := filepath.Join(GraphRoot, fmt.Sprintf("node-%d", ginkgo.GinkgoParallelNode()))

	depotPath := filepath.Join(tmpDir, "containers")
	Expect(os.MkdirAll(depotPath, 0755)).To(Succeed())

	snapshotsPath := filepath.Join(tmpDir, "snapshots")
	Expect(os.MkdirAll(snapshotsPath, 0755)).To(Succeed())

	r := &RunningGarden{
		GraphRoot:     GraphRoot,
		GraphPath:     graphPath,
		DepotPath:     depotPath,
		SnapshotsPath: snapshotsPath,
		tmpdir:        tmpDir,
		logger:        lagertest.NewTestLogger("garden-runner"),

		Client: client.New(connection.New(network, addr)),
	}

	c := cmd(depotPath, snapshotsPath, graphPath, network, addr, GardenBin, BinPath, RootFSPath, argv...)
	r.runner = ginkgomon.New(ginkgomon.Config{
		Name:              "garden-linux",
		Command:           c,
		AnsiColorCode:     "31m",
		StartCheck:        "garden-linux.started",
		StartCheckTimeout: 30 * time.Second,
	})

	r.process = ifrit.Invoke(r.runner)
	r.Pid = c.Process.Pid

	return r
}

func (r *RunningGarden) Buffer() *gbytes.Buffer {
	return r.runner.Buffer()
}

func (r *RunningGarden) Kill() error {
	r.process.Signal(syscall.SIGKILL)
	select {
	case err := <-r.process.Wait():
		return err
	case <-time.After(time.Second * 10):
		r.process.Signal(syscall.SIGKILL)
		return errors.New("timed out waiting for garden to shutdown after 10 seconds")
	}
}

func (r *RunningGarden) DestroyAndStop() error {
	if err := r.DestroyContainers(); err != nil {
		return err
	}

	if err := r.Stop(); err != nil {
		return err
	}

	return nil
}

func (r *RunningGarden) Stop() error {
	r.process.Signal(syscall.SIGTERM)
	select {
	case err := <-r.process.Wait():
		return err
	case <-time.After(time.Second * 10):
		r.process.Signal(syscall.SIGKILL)
		return errors.New("timed out waiting for garden to shutdown after 10 seconds")
	}
}

func cmd(depotPath, snapshotsPath, graphPath, network, addr, bin, binPath, RootFSPath string, argv ...string) *exec.Cmd {
	appendDefaultFlag := func(ar []string, key, value string) []string {
		for _, a := range argv {
			if a == key {
				return ar
			}
		}

		if value != "" {
			return append(ar, key, value)
		} else {
			return append(ar, key)
		}
	}

	hasFlag := func(ar []string, key string) bool {
		for _, a := range ar {
			if a == key {
				return true
			}
		}

		return false
	}

	gardenArgs := make([]string, len(argv))
	copy(gardenArgs, argv)

	gardenArgs = appendDefaultFlag(gardenArgs, "--listenNetwork", network)
	gardenArgs = appendDefaultFlag(gardenArgs, "--listenAddr", addr)
	gardenArgs = appendDefaultFlag(gardenArgs, "--bin", binPath)
	if RootFSPath != "" { //rootfs is an optional parameter
		gardenArgs = appendDefaultFlag(gardenArgs, "--rootfs", RootFSPath)
	}
	gardenArgs = appendDefaultFlag(gardenArgs, "--depot", depotPath)
	gardenArgs = appendDefaultFlag(gardenArgs, "--snapshots", snapshotsPath)
	gardenArgs = appendDefaultFlag(gardenArgs, "--graph", graphPath)
	gardenArgs = appendDefaultFlag(gardenArgs, "--logLevel", "debug")
	gardenArgs = appendDefaultFlag(gardenArgs, "--networkPool", fmt.Sprintf("10.250.%d.0/24", ginkgo.GinkgoParallelNode()))
	gardenArgs = appendDefaultFlag(gardenArgs, "--portPoolStart", strconv.Itoa(51000+(1000*ginkgo.GinkgoParallelNode())))
	gardenArgs = appendDefaultFlag(gardenArgs, "--portPoolSize", "1000")
	gardenArgs = appendDefaultFlag(gardenArgs, "--tag", strconv.Itoa(ginkgo.GinkgoParallelNode()))

	if !hasFlag(gardenArgs, "-enableGraphCleanup=false") {
		gardenArgs = appendDefaultFlag(gardenArgs, "--enableGraphCleanup", "")
	}

	gardenArgs = appendDefaultFlag(gardenArgs, "--debugAddr", fmt.Sprintf(":808%d", ginkgo.GinkgoParallelNode()))

	return exec.Command(bin, gardenArgs...)
}

func (r *RunningGarden) Cleanup() {
	err := filepath.Walk(
		filepath.Join(r.GraphPath, "aufs", "mnt"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			return exec.Command("umount", path).Run()
		})
	if err != nil {
		r.logger.Error("unmount graph layers", err)
	}
	if err := exec.Command("umount", filepath.Join(r.GraphPath, "aufs")).Run(); err != nil {
		r.logger.Error("unmount aufs directory", err)
	}
	if err := os.RemoveAll(r.GraphPath); err != nil {
		r.logger.Error("remove graph", err)
	}

	r.logger.Info("cleanup-tempdirs")
	if err := os.RemoveAll(r.tmpdir); err != nil {
		r.logger.Error("cleanup-tempdirs-failed", err, lager.Data{"tmpdir": r.tmpdir})
	} else {
		r.logger.Info("tempdirs-removed")
	}
}

func (r *RunningGarden) DestroyContainers() error {
	containers, err := r.Containers(nil)
	if err != nil {
		return err
	}

	for _, container := range containers {
		err := r.Destroy(container.Handle())
		if err != nil {
			return err
		}
	}

	return nil
}
