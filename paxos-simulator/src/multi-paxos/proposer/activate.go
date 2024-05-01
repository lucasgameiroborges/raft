package Proposer

import (
	"encoding/json"
	"github.com/paxos/src/pkg/model/message"
	"github.com/paxos/src/pkg/model/node"
	"github.com/paxos/src/pkg/shared/constant"
	"log"
	"net"
)

type Config struct {
	Proposer node.Proposer
}

func Activate(port string, acceptors []string) {
	c := &Config{
		Proposer: node.Proposer{
			Port:      port,
			Acceptors: acceptors,
		},
	}

	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Printf("Failed to connect to port: 9001, error: %v\n ", err)
		return
	}

	log.Printf("Accepting messages on: :9001\n")
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				log.Printf("Error received while listening :9001\n")
			}
		}

		msg := &message.Message{}
		if err := json.NewDecoder(connIn).Decode(msg); err != nil {
			log.Printf("Error decoding %v\n", err)
		}

		switch msg.Type {
		case constant.REQUEST:
			if err := c.handleRequest(msg); err != nil {
				log.Fatalf("Failed to handle a [request]: %v\n", err)
			}
			break
		case constant.PROMISE:
			if err := c.handlePromise(msg); err != nil {
				log.Fatalf("Failed to handle a [promise]: %v\n", err)
			}
			break
		case constant.NACK:
			if err := c.handleNack(msg); err != nil {
				log.Fatalf("Failed to handle a [nack]: %v\n", err)
			}
			break
		}
	}
}
