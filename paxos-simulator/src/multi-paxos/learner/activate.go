package Learner

import (
	"encoding/json"
	"github.com/paxos/src/pkg/model/message"
	"github.com/paxos/src/pkg/model/node"
	"github.com/paxos/src/pkg/shared/constant"
	"log"
	"net"
)

type Config struct {
	Learner node.Learner
}

// Activate activates a learner node
// A learner must be initialized with a port number to be identified with
func Activate(port string) {
	c := &Config{
		Learner: node.Learner{
			Port: port,
		},
	}

	ln, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Printf("Failed to connect to port 9003, error: %v\n ", err)
		return
	}

	log.Printf("Accepting messages on: :9003\n")
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				log.Printf("Error received while listening :9003\n")
			}
		}

		msg := &message.Message{}
		if err := json.NewDecoder(connIn).Decode(msg); err != nil {
			log.Printf("Error decoding %v\n", err)
		}

		switch msg.Type {
		case constant.ACCEPTED:
			if err := c.handleAccepted(msg); err != nil {
				log.Fatalf("Failed to handle an [accepted] message: %v\n", err)
			}
			break
		}
	}
}
