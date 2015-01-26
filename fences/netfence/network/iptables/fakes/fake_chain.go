// This file was generated by counterfeiter
package fakes

import (
	"net"
	"sync"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/garden-linux/fences/netfence/network/iptables"
)

type FakeChain struct {
	SetupStub        func() error
	setupMutex       sync.RWMutex
	setupArgsForCall []struct{}
	setupReturns struct {
		result1 error
	}
	TearDownStub        func() error
	tearDownMutex       sync.RWMutex
	tearDownArgsForCall []struct{}
	tearDownReturns struct {
		result1 error
	}
	AppendRuleStub        func(source string, destination string, jump iptables.Action) error
	appendRuleMutex       sync.RWMutex
	appendRuleArgsForCall []struct {
		source      string
		destination string
		jump        iptables.Action
	}
	appendRuleReturns struct {
		result1 error
	}
	DeleteRuleStub        func(source string, destination string, jump iptables.Action) error
	deleteRuleMutex       sync.RWMutex
	deleteRuleArgsForCall []struct {
		source      string
		destination string
		jump        iptables.Action
	}
	deleteRuleReturns struct {
		result1 error
	}
	AppendNatRuleStub        func(source string, destination string, jump iptables.Action, to net.IP) error
	appendNatRuleMutex       sync.RWMutex
	appendNatRuleArgsForCall []struct {
		source      string
		destination string
		jump        iptables.Action
		to          net.IP
	}
	appendNatRuleReturns struct {
		result1 error
	}
	DeleteNatRuleStub        func(source string, destination string, jump iptables.Action, to net.IP) error
	deleteNatRuleMutex       sync.RWMutex
	deleteNatRuleArgsForCall []struct {
		source      string
		destination string
		jump        iptables.Action
		to          net.IP
	}
	deleteNatRuleReturns struct {
		result1 error
	}
	PrependFilterRuleStub        func(protocol garden.Protocol, dest string, destPort uint32, destPortRange string, destIcmpType, destIcmpCode int32, log bool) error
	prependFilterRuleMutex       sync.RWMutex
	prependFilterRuleArgsForCall []struct {
		protocol      garden.Protocol
		dest          string
		destPort      uint32
		destPortRange string
		destIcmpType  int32
		destIcmpCode  int32
		log           bool
	}
	prependFilterRuleReturns struct {
		result1 error
	}
}

func (fake *FakeChain) Setup() error {
	fake.setupMutex.Lock()
	fake.setupArgsForCall = append(fake.setupArgsForCall, struct{}{})
	fake.setupMutex.Unlock()
	if fake.SetupStub != nil {
		return fake.SetupStub()
	} else {
		return fake.setupReturns.result1
	}
}

func (fake *FakeChain) SetupCallCount() int {
	fake.setupMutex.RLock()
	defer fake.setupMutex.RUnlock()
	return len(fake.setupArgsForCall)
}

func (fake *FakeChain) SetupReturns(result1 error) {
	fake.SetupStub = nil
	fake.setupReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeChain) TearDown() error {
	fake.tearDownMutex.Lock()
	fake.tearDownArgsForCall = append(fake.tearDownArgsForCall, struct{}{})
	fake.tearDownMutex.Unlock()
	if fake.TearDownStub != nil {
		return fake.TearDownStub()
	} else {
		return fake.tearDownReturns.result1
	}
}

func (fake *FakeChain) TearDownCallCount() int {
	fake.tearDownMutex.RLock()
	defer fake.tearDownMutex.RUnlock()
	return len(fake.tearDownArgsForCall)
}

func (fake *FakeChain) TearDownReturns(result1 error) {
	fake.TearDownStub = nil
	fake.tearDownReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeChain) AppendRule(source string, destination string, jump iptables.Action) error {
	fake.appendRuleMutex.Lock()
	fake.appendRuleArgsForCall = append(fake.appendRuleArgsForCall, struct {
		source      string
		destination string
		jump        iptables.Action
	}{source, destination, jump})
	fake.appendRuleMutex.Unlock()
	if fake.AppendRuleStub != nil {
		return fake.AppendRuleStub(source, destination, jump)
	} else {
		return fake.appendRuleReturns.result1
	}
}

func (fake *FakeChain) AppendRuleCallCount() int {
	fake.appendRuleMutex.RLock()
	defer fake.appendRuleMutex.RUnlock()
	return len(fake.appendRuleArgsForCall)
}

func (fake *FakeChain) AppendRuleArgsForCall(i int) (string, string, iptables.Action) {
	fake.appendRuleMutex.RLock()
	defer fake.appendRuleMutex.RUnlock()
	return fake.appendRuleArgsForCall[i].source, fake.appendRuleArgsForCall[i].destination, fake.appendRuleArgsForCall[i].jump
}

func (fake *FakeChain) AppendRuleReturns(result1 error) {
	fake.AppendRuleStub = nil
	fake.appendRuleReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeChain) DeleteRule(source string, destination string, jump iptables.Action) error {
	fake.deleteRuleMutex.Lock()
	fake.deleteRuleArgsForCall = append(fake.deleteRuleArgsForCall, struct {
		source      string
		destination string
		jump        iptables.Action
	}{source, destination, jump})
	fake.deleteRuleMutex.Unlock()
	if fake.DeleteRuleStub != nil {
		return fake.DeleteRuleStub(source, destination, jump)
	} else {
		return fake.deleteRuleReturns.result1
	}
}

func (fake *FakeChain) DeleteRuleCallCount() int {
	fake.deleteRuleMutex.RLock()
	defer fake.deleteRuleMutex.RUnlock()
	return len(fake.deleteRuleArgsForCall)
}

func (fake *FakeChain) DeleteRuleArgsForCall(i int) (string, string, iptables.Action) {
	fake.deleteRuleMutex.RLock()
	defer fake.deleteRuleMutex.RUnlock()
	return fake.deleteRuleArgsForCall[i].source, fake.deleteRuleArgsForCall[i].destination, fake.deleteRuleArgsForCall[i].jump
}

func (fake *FakeChain) DeleteRuleReturns(result1 error) {
	fake.DeleteRuleStub = nil
	fake.deleteRuleReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeChain) AppendNatRule(source string, destination string, jump iptables.Action, to net.IP) error {
	fake.appendNatRuleMutex.Lock()
	fake.appendNatRuleArgsForCall = append(fake.appendNatRuleArgsForCall, struct {
		source      string
		destination string
		jump        iptables.Action
		to          net.IP
	}{source, destination, jump, to})
	fake.appendNatRuleMutex.Unlock()
	if fake.AppendNatRuleStub != nil {
		return fake.AppendNatRuleStub(source, destination, jump, to)
	} else {
		return fake.appendNatRuleReturns.result1
	}
}

func (fake *FakeChain) AppendNatRuleCallCount() int {
	fake.appendNatRuleMutex.RLock()
	defer fake.appendNatRuleMutex.RUnlock()
	return len(fake.appendNatRuleArgsForCall)
}

func (fake *FakeChain) AppendNatRuleArgsForCall(i int) (string, string, iptables.Action, net.IP) {
	fake.appendNatRuleMutex.RLock()
	defer fake.appendNatRuleMutex.RUnlock()
	return fake.appendNatRuleArgsForCall[i].source, fake.appendNatRuleArgsForCall[i].destination, fake.appendNatRuleArgsForCall[i].jump, fake.appendNatRuleArgsForCall[i].to
}

func (fake *FakeChain) AppendNatRuleReturns(result1 error) {
	fake.AppendNatRuleStub = nil
	fake.appendNatRuleReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeChain) DeleteNatRule(source string, destination string, jump iptables.Action, to net.IP) error {
	fake.deleteNatRuleMutex.Lock()
	fake.deleteNatRuleArgsForCall = append(fake.deleteNatRuleArgsForCall, struct {
		source      string
		destination string
		jump        iptables.Action
		to          net.IP
	}{source, destination, jump, to})
	fake.deleteNatRuleMutex.Unlock()
	if fake.DeleteNatRuleStub != nil {
		return fake.DeleteNatRuleStub(source, destination, jump, to)
	} else {
		return fake.deleteNatRuleReturns.result1
	}
}

func (fake *FakeChain) DeleteNatRuleCallCount() int {
	fake.deleteNatRuleMutex.RLock()
	defer fake.deleteNatRuleMutex.RUnlock()
	return len(fake.deleteNatRuleArgsForCall)
}

func (fake *FakeChain) DeleteNatRuleArgsForCall(i int) (string, string, iptables.Action, net.IP) {
	fake.deleteNatRuleMutex.RLock()
	defer fake.deleteNatRuleMutex.RUnlock()
	return fake.deleteNatRuleArgsForCall[i].source, fake.deleteNatRuleArgsForCall[i].destination, fake.deleteNatRuleArgsForCall[i].jump, fake.deleteNatRuleArgsForCall[i].to
}

func (fake *FakeChain) DeleteNatRuleReturns(result1 error) {
	fake.DeleteNatRuleStub = nil
	fake.deleteNatRuleReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeChain) PrependFilterRule(protocol garden.Protocol, dest string, destPort uint32, destPortRange string, destIcmpType int32, destIcmpCode int32, log bool) error {
	fake.prependFilterRuleMutex.Lock()
	fake.prependFilterRuleArgsForCall = append(fake.prependFilterRuleArgsForCall, struct {
		protocol      garden.Protocol
		dest          string
		destPort      uint32
		destPortRange string
		destIcmpType  int32
		destIcmpCode  int32
		log           bool
	}{protocol, dest, destPort, destPortRange, destIcmpType, destIcmpCode, log})
	fake.prependFilterRuleMutex.Unlock()
	if fake.PrependFilterRuleStub != nil {
		return fake.PrependFilterRuleStub(protocol, dest, destPort, destPortRange, destIcmpType, destIcmpCode, log)
	} else {
		return fake.prependFilterRuleReturns.result1
	}
}

func (fake *FakeChain) PrependFilterRuleCallCount() int {
	fake.prependFilterRuleMutex.RLock()
	defer fake.prependFilterRuleMutex.RUnlock()
	return len(fake.prependFilterRuleArgsForCall)
}

func (fake *FakeChain) PrependFilterRuleArgsForCall(i int) (garden.Protocol, string, uint32, string, int32, int32, bool) {
	fake.prependFilterRuleMutex.RLock()
	defer fake.prependFilterRuleMutex.RUnlock()
	return fake.prependFilterRuleArgsForCall[i].protocol, fake.prependFilterRuleArgsForCall[i].dest, fake.prependFilterRuleArgsForCall[i].destPort, fake.prependFilterRuleArgsForCall[i].destPortRange, fake.prependFilterRuleArgsForCall[i].destIcmpType, fake.prependFilterRuleArgsForCall[i].destIcmpCode, fake.prependFilterRuleArgsForCall[i].log
}

func (fake *FakeChain) PrependFilterRuleReturns(result1 error) {
	fake.PrependFilterRuleStub = nil
	fake.prependFilterRuleReturns = struct {
		result1 error
	}{result1}
}

var _ iptables.Chain = new(FakeChain)