package node

import (
	"github.com/paxos/src/pkg/model/message"
	"math/rand"
	"time"
)

type Proposal struct {
	Value    string
	Nonce    int64
	Quorum   []string
	Promises []message.Promise
}

type Proposer struct {
	Port         string
	Acceptors    []string
	Proposals    []Proposal
	CurrentNonce int64
}

///////////////////////////
//// Proposal Helpers ////
//////////////////////////

func (p *Proposal) NonceDoesNotEqual(nonce int64) bool {
	return p.Nonce != nonce
}

func (p *Proposal) RegisterPromise(promise message.Promise) {
	p.Promises = append(p.Promises, promise)
}

func (p *Proposal) HasInsufficientNumberOfPromises() bool {
	return len(p.Promises) != len(p.Quorum)
}

func (p *Proposal) HasAcceptedValueToBroadcast() bool {
	for _, promise := range p.Promises {
		if promise.Proposal != (message.Proposal{}) {
			return true
		}
	}
	return false
}

func (p *Proposal) GetAcceptedValueToBroadcast() string {
	nonce := time.Now().UnixNano()
	value := ""

	for _, promise := range p.Promises {
		if nonce < promise.Proposal.Nonce {
			nonce = promise.Proposal.Nonce
			value = promise.Proposal.Value
		}
	}

	return value
}

///////////////////////////
//// Proposer Helpers ////
//////////////////////////

func (p *Proposer) GetQuorum() []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p.Acceptors), func(i, j int) { p.Acceptors[i], p.Acceptors[j] = p.Acceptors[j], p.Acceptors[i] })
	return p.Acceptors[:(len(p.Acceptors)/2)+1]
}

func (p *Proposer) GetNonce() int64 {
	return time.Now().UnixNano()
}

func (p *Proposer) AddProposal(value string) Proposal {
	proposal := Proposal{
		Value:  value,
		Nonce:  p.GetNonce(),
		Quorum: p.GetQuorum(),
	}

	p.Proposals = append(p.Proposals, proposal)

	return proposal
}
