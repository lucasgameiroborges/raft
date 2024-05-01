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
	fmt.Println("Aprovaram algo: ", acceptedMessage.Value)
	fmt.Println("Que round pensam que ta: ", acceptedMessage.Round)
	fmt.Println("tamanho do log: ", len(c.Learner.Logs))

	// If the learner has no log of any other value being accepted by the network for this round,
	// log it and inform the client of the accepted value
	if acceptedMessage.Round > len(c.Learner.Logs) {
		util.WriteFile("log", acceptedMessage.Value)
		c.Learner.Logs = append(c.Learner.Logs, acceptedMessage.Value)
	}

	return nil
}
