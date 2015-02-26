// This file was generated by counterfeiter
package fake_cn_persistor

import (
	"sync"

	"github.com/cloudfoundry-incubator/garden-linux/container_pool"
	"github.com/cloudfoundry-incubator/garden-linux/network/cnet"
)

type FakeCNPersistor struct {
	PersistStub        func(cn cnet.ContainerNetwork, path string) error
	persistMutex       sync.RWMutex
	persistArgsForCall []struct {
		cn   cnet.ContainerNetwork
		path string
	}
	persistReturns struct {
		result1 error
	}
	RecoverStub        func(path string) (cnet.ContainerNetwork, error)
	recoverMutex       sync.RWMutex
	recoverArgsForCall []struct {
		path string
	}
	recoverReturns struct {
		result1 cnet.ContainerNetwork
		result2 error
	}
}

func (fake *FakeCNPersistor) Persist(cn cnet.ContainerNetwork, path string) error {
	fake.persistMutex.Lock()
	fake.persistArgsForCall = append(fake.persistArgsForCall, struct {
		cn   cnet.ContainerNetwork
		path string
	}{cn, path})
	fake.persistMutex.Unlock()
	if fake.PersistStub != nil {
		return fake.PersistStub(cn, path)
	} else {
		return fake.persistReturns.result1
	}
}

func (fake *FakeCNPersistor) PersistCallCount() int {
	fake.persistMutex.RLock()
	defer fake.persistMutex.RUnlock()
	return len(fake.persistArgsForCall)
}

func (fake *FakeCNPersistor) PersistArgsForCall(i int) (cnet.ContainerNetwork, string) {
	fake.persistMutex.RLock()
	defer fake.persistMutex.RUnlock()
	return fake.persistArgsForCall[i].cn, fake.persistArgsForCall[i].path
}

func (fake *FakeCNPersistor) PersistReturns(result1 error) {
	fake.PersistStub = nil
	fake.persistReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCNPersistor) Recover(path string) (cnet.ContainerNetwork, error) {
	fake.recoverMutex.Lock()
	fake.recoverArgsForCall = append(fake.recoverArgsForCall, struct {
		path string
	}{path})
	fake.recoverMutex.Unlock()
	if fake.RecoverStub != nil {
		return fake.RecoverStub(path)
	} else {
		return fake.recoverReturns.result1, fake.recoverReturns.result2
	}
}

func (fake *FakeCNPersistor) RecoverCallCount() int {
	fake.recoverMutex.RLock()
	defer fake.recoverMutex.RUnlock()
	return len(fake.recoverArgsForCall)
}

func (fake *FakeCNPersistor) RecoverArgsForCall(i int) string {
	fake.recoverMutex.RLock()
	defer fake.recoverMutex.RUnlock()
	return fake.recoverArgsForCall[i].path
}

func (fake *FakeCNPersistor) RecoverReturns(result1 cnet.ContainerNetwork, result2 error) {
	fake.RecoverStub = nil
	fake.recoverReturns = struct {
		result1 cnet.ContainerNetwork
		result2 error
	}{result1, result2}
}

var _ container_pool.CNPersistor = new(FakeCNPersistor)
