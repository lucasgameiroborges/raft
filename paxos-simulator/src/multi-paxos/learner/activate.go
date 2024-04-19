package Learner

import (
	"encoding/json"
	"fmt"
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
func Activate(port int) {
	c := &Config{
		Learner: node.Learner {
			Port: port,
		},
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d",c.Learner.Port))
	if err != nil {
		log.Printf("Failed to connect to port: %d, error: %v\n ", c.Learner.Port, err)
		return
	}

	log.Printf("Accepting messages on: 127.0.0.1:%d\n", c.Learner.Port)
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				log.Printf("Error received while listening 127.0.0.1:%d\n", c.Learner.Port)
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
