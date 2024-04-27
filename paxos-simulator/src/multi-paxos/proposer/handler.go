package Proposer

import (
	"fmt"
	"github.com/paxos/src/pkg/model/message"
	"github.com/paxos/src/pkg/shared/constant"
	"github.com/paxos/src/pkg/shared/util"
)

// handleRequest processes request messages
func (c *Config) handleRequest(incomingMessage *message.Message) error {
	requestMessage := &message.Request{}
	if err := message.Unmarshal(incomingMessage.Payload, requestMessage); err != nil {
		return err
	}

	// Add the request to the proposers list of proposals
	proposal := c.Proposer.AddProposal(requestMessage.Value)

	// Construct the proposal message for the new round
	outgoingMessage := &message.Message{
		Source: c.Proposer.Port,
		Type:   constant.PREPARE,
		Payload: message.Prepare{
			Nonce: proposal.Nonce,
			Round: len(c.Proposer.Proposals),
		},
	}

	// Broadcast a prepare message to the proposals quorum of acceptors for the new round
	for _, acceptor := range proposal.Quorum {
		fmt.Println("proposer %d->> acceptor %d:(%d) Prepare", c.Proposer.Port, acceptor, proposal.Nonce)
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}

// handlePromise processes promise messages
func (c *Config) handlePromise(incomingMessage *message.Message) error {
	promiseMessage := &message.Promise{}
	if err := message.Unmarshal(incomingMessage.Payload, promiseMessage); err != nil {
		return err
	}

	// Add the promise to the proposers list of promises
	c.Proposer.Proposals[promiseMessage.Round-1].RegisterPromise(*promiseMessage)

	// If the proposer has not received a sufficient number of promises for its current proposal, do nothing
	if c.Proposer.Proposals[promiseMessage.Round-1].HasInsufficientNumberOfPromises() {
		// TODO exit gracefully
		return nil
	}

	// Construct the accept message
	// If the proposer has learned that another proposal has already been accepted for this round, share that with its proposals
	// quorum of acceptors
	payload := message.Accept{
		Nonce: c.Proposer.Proposals[promiseMessage.Round-1].Nonce,
		Round: promiseMessage.Round,
	}
	if c.Proposer.Proposals[promiseMessage.Round-1].HasAcceptedValueToBroadcast() {
		payload.Value = c.Proposer.Proposals[promiseMessage.Round-1].GetAcceptedValueToBroadcast()
	} else {
		payload.Value = c.Proposer.Proposals[promiseMessage.Round-1].Value
	}
	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    constant.ACCEPT,
		Payload: payload,
	}

	// Broadcast an accept message to the proposals quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[promiseMessage.Round-1].Quorum {
		fmt.Println("proposer %d->> acceptor %d:(%d) Accept: %s", c.Proposer.Port, acceptor, payload.Nonce, payload.Value)
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}

// handleNack processes nack messages
func (c *Config) handleNack(incomingMessage *message.Message) error {
	nackMessage := &message.Nack{}
	if err := message.Unmarshal(incomingMessage.Payload, nackMessage); err != nil {
		return err
	}

	// If the nack is less than the proposers current nonce, it is outdated and can be ignored
	if nackMessage.Nonce < c.Proposer.CurrentNonce {
		return nil
	}

	// Construct a new prepare message with a nonce that is greater than the nonce supplied in the nack message
	c.Proposer.CurrentNonce = nackMessage.Nonce + 1
	c.Proposer.Proposals[nackMessage.Round-1].Nonce = c.Proposer.CurrentNonce
	c.Proposer.Proposals[nackMessage.Round-1].Promises = []message.Promise{}
	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    constant.PREPARE,
		Payload: message.Prepare{Nonce: c.Proposer.Proposals[nackMessage.Round-1].Nonce},
	}

	// Broadcast the updated prepare message for this round to the proposals quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[nackMessage.Round-1].Quorum {
		fmt.Println("proposer %d->> acceptor %d:(%d) Prepare", c.Proposer.Port, acceptor, c.Proposer.Proposals[nackMessage.Round-1].Nonce)
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}
