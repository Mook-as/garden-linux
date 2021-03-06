// This file was generated by counterfeiter
package fake_ptyopener

import (
	"os"
	"sync"

	"github.com/cloudfoundry-incubator/garden-linux/container_daemon"
)

type FakePTYOpener struct {
	OpenStub        func() (pty *os.File, tty *os.File, err error)
	openMutex       sync.RWMutex
	openArgsForCall []struct{}
	openReturns     struct {
		result1 *os.File
		result2 *os.File
		result3 error
	}
}

func (fake *FakePTYOpener) Open() (pty *os.File, tty *os.File, err error) {
	fake.openMutex.Lock()
	fake.openArgsForCall = append(fake.openArgsForCall, struct{}{})
	fake.openMutex.Unlock()
	if fake.OpenStub != nil {
		return fake.OpenStub()
	} else {
		return fake.openReturns.result1, fake.openReturns.result2, fake.openReturns.result3
	}
}

func (fake *FakePTYOpener) OpenCallCount() int {
	fake.openMutex.RLock()
	defer fake.openMutex.RUnlock()
	return len(fake.openArgsForCall)
}

func (fake *FakePTYOpener) OpenReturns(result1 *os.File, result2 *os.File, result3 error) {
	fake.OpenStub = nil
	fake.openReturns = struct {
		result1 *os.File
		result2 *os.File
		result3 error
	}{result1, result2, result3}
}

var _ container_daemon.PTYOpener = new(FakePTYOpener)
