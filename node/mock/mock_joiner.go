package mock

import (
	"errors"

	"github.com/su225/raft/node/cluster"
)

// ErrDiscovery is a dummy discovery related error
var ErrDiscovery = errors.New("error while discovery")

// MockJoiner represents the mock joiner
// which just returns a set of nodes as
// discovered nodes
type MockJoiner struct {
	ShouldDiscoverySucceed bool
	DiscoveredNodes        []cluster.NodeInfo
}

// DiscoverNodes returns the expected list of discovered nodes or error
// if discovery is set so that it should not be successful
func (joiner *MockJoiner) DiscoverNodes() ([]cluster.NodeInfo, error) {
	if joiner.ShouldDiscoverySucceed {
		return []cluster.NodeInfo{}, ErrDiscovery
	}
	return joiner.DiscoveredNodes, nil
}

// GetDefaultMockJoiner returns mock joiner with sensible default
func GetDefaultMockJoiner(discoveryResult bool) *MockJoiner {
	return &MockJoiner{
		ShouldDiscoverySucceed: discoveryResult,
		DiscoveredNodes: []cluster.NodeInfo{
			SampleNodeInfo0,
			SampleNodeInfo1,
			SampleNodeInfo2,
		},
	}
}
