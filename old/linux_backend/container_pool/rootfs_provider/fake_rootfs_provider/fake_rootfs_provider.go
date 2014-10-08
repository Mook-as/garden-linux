// This file was generated by counterfeiter
package fake_rootfs_provider

import (
	"net/url"
	"sync"

	"github.com/cloudfoundry-incubator/garden-linux/old/linux_backend/container_pool/rootfs_provider"
	"github.com/pivotal-golang/lager"
)

type FakeRootFSProvider struct {
	ProvideRootFSStub        func(logger lager.Logger, id string, rootfs *url.URL) (mountpoint string, envvar []string, err error)
	provideRootFSMutex       sync.RWMutex
	provideRootFSArgsForCall []struct {
		logger lager.Logger
		id     string
		rootfs *url.URL
	}
	provideRootFSReturns struct {
		result1 string
		result2 []string
		result3 error
	}
	CleanupRootFSStub        func(logger lager.Logger, id string) error
	cleanupRootFSMutex       sync.RWMutex
	cleanupRootFSArgsForCall []struct {
		logger lager.Logger
		id     string
	}
	cleanupRootFSReturns struct {
		result1 error
	}
}

func (fake *FakeRootFSProvider) ProvideRootFS(logger lager.Logger, id string, rootfs *url.URL) (mountpoint string, envvar []string, err error) {
	fake.provideRootFSMutex.Lock()
	fake.provideRootFSArgsForCall = append(fake.provideRootFSArgsForCall, struct {
		logger lager.Logger
		id     string
		rootfs *url.URL
	}{logger, id, rootfs})
	fake.provideRootFSMutex.Unlock()
	if fake.ProvideRootFSStub != nil {
		return fake.ProvideRootFSStub(logger, id, rootfs)
	} else {
		return fake.provideRootFSReturns.result1, fake.provideRootFSReturns.result2, fake.provideRootFSReturns.result3
	}
}

func (fake *FakeRootFSProvider) ProvideRootFSCallCount() int {
	fake.provideRootFSMutex.RLock()
	defer fake.provideRootFSMutex.RUnlock()
	return len(fake.provideRootFSArgsForCall)
}

func (fake *FakeRootFSProvider) ProvideRootFSArgsForCall(i int) (lager.Logger, string, *url.URL) {
	fake.provideRootFSMutex.RLock()
	defer fake.provideRootFSMutex.RUnlock()
	return fake.provideRootFSArgsForCall[i].logger, fake.provideRootFSArgsForCall[i].id, fake.provideRootFSArgsForCall[i].rootfs
}

func (fake *FakeRootFSProvider) ProvideRootFSReturns(result1 string, result2 []string, result3 error) {
	fake.ProvideRootFSStub = nil
	fake.provideRootFSReturns = struct {
		result1 string
		result2 []string
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeRootFSProvider) CleanupRootFS(logger lager.Logger, id string) error {
	fake.cleanupRootFSMutex.Lock()
	fake.cleanupRootFSArgsForCall = append(fake.cleanupRootFSArgsForCall, struct {
		logger lager.Logger
		id     string
	}{logger, id})
	fake.cleanupRootFSMutex.Unlock()
	if fake.CleanupRootFSStub != nil {
		return fake.CleanupRootFSStub(logger, id)
	} else {
		return fake.cleanupRootFSReturns.result1
	}
}

func (fake *FakeRootFSProvider) CleanupRootFSCallCount() int {
	fake.cleanupRootFSMutex.RLock()
	defer fake.cleanupRootFSMutex.RUnlock()
	return len(fake.cleanupRootFSArgsForCall)
}

func (fake *FakeRootFSProvider) CleanupRootFSArgsForCall(i int) (lager.Logger, string) {
	fake.cleanupRootFSMutex.RLock()
	defer fake.cleanupRootFSMutex.RUnlock()
	return fake.cleanupRootFSArgsForCall[i].logger, fake.cleanupRootFSArgsForCall[i].id
}

func (fake *FakeRootFSProvider) CleanupRootFSReturns(result1 error) {
	fake.CleanupRootFSStub = nil
	fake.cleanupRootFSReturns = struct {
		result1 error
	}{result1}
}

var _ rootfs_provider.RootFSProvider = new(FakeRootFSProvider)