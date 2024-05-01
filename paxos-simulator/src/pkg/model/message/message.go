package message

import "github.com/paxos/src/pkg/shared/constant"

type Message struct {
	Source  string           `json:"source"`
	Type    constant.Type `json:"type"`
	Payload interface{}   `json:"payload"`
}
