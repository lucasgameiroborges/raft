package util

import (
	"encoding/json"
	"fmt"
	"github.com/paxos/src/pkg/model/message"
	"net"
	"time"
)

func SendMessage(msg *message.Message, dest string) error {
	connOut, err := net.DialTimeout("tcp", dest, time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Println("Could not connect with %s\n", dest)
			return err
		}
	}

	if err := json.NewEncoder(connOut).Encode(msg); err != nil {
		fmt.Println("Could not enncode message: %v\n", msg)
		return err
	}
	return nil
}