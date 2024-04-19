package node

import (
	"github.com/paxos/src/pkg/model/message"
	"math/rand"
	"time"
)

type Proposal struct {
	Value    string
	Nonce    int
	Quorum   []int
	Promises []message.Promise
}

type Proposer struct {
	Port         int
	Acceptors    []int
	Proposals    []Proposal
	CurrentNonce int
}

///////////////////////////
//// Proposal Helpers ////
//////////////////////////

func (p *Proposal) NonceDoesNotEqual(nonce int) bool {
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
	nonce := int(0)
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

func (p *Proposer) GetQuorum() []int {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p.Acceptors), func(i, j int) { p.Acceptors[i], p.Acceptors[j] = p.Acceptors[j], p.Acceptors[i] })
	return p.Acceptors[:(len(p.Acceptors)/2)+1]
}

func (p *Proposer) GetNonce() int {
	p.CurrentNonce++
	return p.CurrentNonce
}

func (p *Proposer) AddProposal(value string) Proposal{
	proposal :=  Proposal{
		Value:    value,
		Nonce:    p.GetNonce(),
		Quorum:   p.GetQuorum(),
	}

	p.Proposals = append(p.Proposals, proposal)

	return proposal
}
