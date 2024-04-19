package Learner

import (
	"fmt"
	"github.com/paxos/src/pkg/model/message"
	"github.com/paxos/src/pkg/shared/util"
)

// handleAccepted processes accepted messages
func (c *Config) handleAccepted(incomingMessage *message.Message) error {

	acceptedMessage := &message.Accepted{}
	if err := message.Unmarshal(incomingMessage.Payload, acceptedMessage); err != nil {
		return err
	}

	// If the learner has no log of any other value being accepted by the network for this round,
	// log it and inform the client of the accepted value
	if acceptedMessage.Round > len(c.Learner.Logs) {
		c.Learner.Logs = append(c.Learner.Logs, acceptedMessage.Value)
		util.WriteToMultiFile(fmt.Sprintf("learner %d->>+ client: %s was accepted as the value!", c.Learner.Port, c.Learner.Logs[acceptedMessage.Round-1]))
	}

	return nil
}
