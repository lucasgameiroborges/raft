package Acceptor

import (
	"fmt"
	"github.com/paxos/src/pkg/model/message"
	"github.com/paxos/src/pkg/shared/constant"
	"github.com/paxos/src/pkg/shared/util"
)

// handlePrepare processes prepare messages
func (c *Config) handlePrepare(incomingMessage *message.Message) error {
	prepareMessage := &message.Prepare{}
	if err := message.Unmarshal(incomingMessage.Payload, prepareMessage); err != nil {
		return err
	}

	// If the acceptor has already promised a nonce equal to or larger than the nonce supplied in the prepare message,
	// reject it with a NACK message and return the highest nonce that it has already promised
	for _, promise := range c.Acceptor.Promises {
		if prepareMessage.Nonce <= promise.Nonce {
			outgoingMessage := &message.Message{
				Source: c.Acceptor.Port,
				Type:   constant.NACK,
				Payload: message.Nack{
					Nonce: promise.Nonce,
					Round: prepareMessage.Round,
				},
			}
			util.WriteToMultiFile(fmt.Sprintf("acceptor %d--x proposer %d:(%d) Nack", c.Acceptor.Port, incomingMessage.Source, prepareMessage.Nonce))
			if err := util.SendMessage(outgoingMessage, incomingMessage.Source); err != nil {
				return err
			}
			return nil
		}
	}

	// Construct the promise
	promise := message.Promise{
		Nonce: prepareMessage.Nonce,
		Round: prepareMessage.Round,
	}

	// If the acceptor has already accepted a proposal for this round (default of 1 for Basic-Paxos) then include it
	// in its promise to the proposer
	if c.Acceptor.HasAcceptedProposal(prepareMessage.Round) {
		promise.Proposal = c.Acceptor.AcceptedProposals[prepareMessage.Round-1]
	}

	// Add the promise to the acceptors list of promises
	c.Acceptor.AddPromise(promise)

	// Construct the promise message
	outgoingMessage := &message.Message{
		Source:  c.Acceptor.Port,
		Type:    constant.PROMISE,
		Payload: promise,
	}

	// Send the promise message to proposer
	if c.Acceptor.HasAcceptedProposal(prepareMessage.Round) { // Send a promise with a proposal that has already been accepted
		util.WriteToMultiFile(fmt.Sprintf("acceptor %d-->> proposer %d:(%d) Promise: %+v", c.Acceptor.Port, incomingMessage.Source, prepareMessage.Nonce, promise.Proposal))
		if err := util.SendMessage(outgoingMessage, incomingMessage.Source); err != nil {
			return err
		}
	} else { // Send a promise to not accept any nonce equal to or less than the one supplied in the prepare message
		util.WriteToMultiFile(fmt.Sprintf("acceptor %d-->> proposer %d:(%d) Promise", c.Acceptor.Port, incomingMessage.Source, prepareMessage.Nonce))
		if err := util.SendMessage(outgoingMessage, incomingMessage.Source); err != nil {
			return err
		}
	}

	return nil
}

// handleAccept handles accept messages
func (c *Config) handleAccept(incomingMessage *message.Message) error {
	acceptMessage := &message.Accept{}
	if err := message.Unmarshal(incomingMessage.Payload, acceptMessage); err != nil {
		return err
	}

	// If the Acceptor has already promised a nonce equal to or larger than the nonce supplied in the accept message,
	// reject it with a NACK message and return the highest nonce that it has already promised
	for _, promise := range c.Acceptor.Promises {
		if acceptMessage.Nonce < promise.Nonce {
			outgoingMessage := &message.Message{
				Source: c.Acceptor.Port,
				Type:   constant.NACK,
				Payload: message.Nack{
					Nonce: promise.Nonce,
					Round: acceptMessage.Round,
				},
			}
			util.WriteToMultiFile(fmt.Sprintf("acceptor %d--x proposer %d:(%d) Nack", c.Acceptor.Port, incomingMessage.Source, acceptMessage.Nonce))
			if err := util.SendMessage(outgoingMessage, incomingMessage.Source); err != nil {
				return err
			}
		}
	}

	// Add the proposal to the acceptors list of accepted proposals
	c.Acceptor.AddAcceptedProposal(acceptMessage.Value, acceptMessage.Nonce)

	// Construct the accept message
	outgoingMessage := &message.Message{
		Source: c.Acceptor.Port,
		Type:   constant.ACCEPTED,
		Payload: message.Accepted{
			Nonce: acceptMessage.Nonce,
			Value: acceptMessage.Value,
			Round: acceptMessage.Round,
		},
	}

	// Broadcast that a proposal has been accepted by this acceptor for this round to its list of learners
	for _, learner := range c.Acceptor.Learners {
		util.WriteToMultiFile(fmt.Sprintf("acceptor %d-->> learner %d:(%d) Accepted: %s", c.Acceptor.Port, learner, acceptMessage.Nonce, acceptMessage.Value))
		if err := util.SendMessage(outgoingMessage, learner); err != nil {
			return err
		}
	}

	return nil
}
