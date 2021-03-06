// Code generated by counterfeiter. DO NOT EDIT.
package dnsresolverfakes

import (
	"bosh-dns/dns/server/records/dnsresolver"
	"sync"
)

type FakeRecordSet struct {
	ResolveStub        func(domain string) ([]string, error)
	resolveMutex       sync.RWMutex
	resolveArgsForCall []struct {
		domain string
	}
	resolveReturns struct {
		result1 []string
		result2 error
	}
	resolveReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRecordSet) Resolve(domain string) ([]string, error) {
	fake.resolveMutex.Lock()
	ret, specificReturn := fake.resolveReturnsOnCall[len(fake.resolveArgsForCall)]
	fake.resolveArgsForCall = append(fake.resolveArgsForCall, struct {
		domain string
	}{domain})
	fake.recordInvocation("Resolve", []interface{}{domain})
	fake.resolveMutex.Unlock()
	if fake.ResolveStub != nil {
		return fake.ResolveStub(domain)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.resolveReturns.result1, fake.resolveReturns.result2
}

func (fake *FakeRecordSet) ResolveCallCount() int {
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	return len(fake.resolveArgsForCall)
}

func (fake *FakeRecordSet) ResolveArgsForCall(i int) string {
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	return fake.resolveArgsForCall[i].domain
}

func (fake *FakeRecordSet) ResolveReturns(result1 []string, result2 error) {
	fake.ResolveStub = nil
	fake.resolveReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeRecordSet) ResolveReturnsOnCall(i int, result1 []string, result2 error) {
	fake.ResolveStub = nil
	if fake.resolveReturnsOnCall == nil {
		fake.resolveReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.resolveReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeRecordSet) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeRecordSet) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ dnsresolver.RecordSet = new(FakeRecordSet)
