package Proposer

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
	Proposer node.Proposer
}

func Activate(port int, acceptors []int) {
	c := &Config{
		Proposer: node.Proposer{
			Port:      port,
			Acceptors: acceptors,
		},
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", c.Proposer.Port))
	if err != nil {
		log.Printf("Failed to connect to port: %d, error: %v\n ", c.Proposer.Port, err)
		return
	}

	log.Printf("Accepting messages on: 127.0.0.1:%d\n", c.Proposer.Port)
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				log.Printf("Error received while listening 127.0.0.1:%d\n", c.Proposer.Port)
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
