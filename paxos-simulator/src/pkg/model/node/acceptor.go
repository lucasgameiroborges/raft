package node

import "github.com/paxos/src/pkg/model/message"

type Acceptor struct {
	Port              int
	Learners          []int
	Promises          []message.Promise
	AcceptedProposals []message.Proposal
}

func (a *Acceptor) HasPromisedGreaterNonceThan(nonce int) bool {
	for _, promise := range a.Promises {
		if nonce <= promise.Nonce {
			return true
		}
	}
	return false
}

func (a *Acceptor) HasAcceptedProposal(round int) bool {
	return round == len(a.AcceptedProposals)
}

func (a *Acceptor) AddPromise(promise message.Promise) {
	a.Promises = append(a.Promises, promise)
}

func (a *Acceptor) AddAcceptedProposal(value string, nonce int) {
	a.AcceptedProposals = append(a.AcceptedProposals, message.Proposal{
		Value: value,
		Nonce: nonce,
	})
}
