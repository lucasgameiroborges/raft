package Acceptor

import (
	"encoding/json"
	"github.com/paxos/src/pkg/model/message"
	"github.com/paxos/src/pkg/model/node"
	"github.com/paxos/src/pkg/shared/constant"
	"log"
	"net"
)

type Config struct {
	Acceptor node.Acceptor
}

// Activate activates an acceptor node
// An acceptor must be initialized with a port number to be identified with and a set of learner nodes
// which it may communicate with
func Activate(port string, learners []string) {
	c := &Config{
		Acceptor: node.Acceptor{
			Port:     port,
			Learners: learners,
		},
	}

	ln, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Printf("Failed to connect to port 9002, error: %v\n ", err)
		return
	}

	log.Printf("Accepting messages on: :9002\n")
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				log.Printf("Error received while listening :9002\n")
			}
		}

		msg := &message.Message{}
		if err := json.NewDecoder(connIn).Decode(msg); err != nil {
			log.Printf("Error decoding %v\n", err)
		}

		switch msg.Type {
		case constant.PREPARE:
			if err := c.handlePrepare(msg); err != nil {
				log.Fatalf("Failed to handle a [prepare]: %v\n", err)
			}
			break
		case constant.ACCEPT:
			if err := c.handleAccept(msg); err != nil {
				log.Fatalf("Failed to handle an [accept]: %v\n", err)
			}
			break
		}
	}
}
