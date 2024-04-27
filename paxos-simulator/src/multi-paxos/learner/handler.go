package Learner

import (
	"fmt"
	"github.com/paxos/src/pkg/model/message"
	"os"
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
		file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		defer file.Close()
		_, err = file.WriteString(acceptedMessage.Value + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return nil
		}
		c.Learner.Logs = append(c.Learner.Logs, acceptedMessage.Value)
		fmt.Println("learner %d->>+ client: %s was accepted as the value!", c.Learner.Port, c.Learner.Logs[acceptedMessage.Round-1])
	}

	return nil
}
