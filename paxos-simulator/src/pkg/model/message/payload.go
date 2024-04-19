package message

import "encoding/json"

type Prepare struct {
	Nonce int `json:"nonce"`
	Round int `json:"round"`
}

type Promise struct {
	Nonce    int      `json:"nonce"`
	Round    int      `json:"round"`
	Proposal Proposal `json:"proposal"`
}

type Accept struct {
	Nonce int    `json:"nonce"`
	Round int    `json:"round"`
	Value string `json:"value"`
}

type Accepted struct {
	Nonce int    `json:"nonce"`
	Round int    `json:"round"`
	Value string `json:"value"`
}

type Nack struct {
	Nonce int `json:"nonce"`
	Round int `json:"round"`
}

type Request struct {
	Value string `json:"value"`
}

type Response struct {
	Value string `json:"value"`
}

type Proposal struct {
	Value string `json:"value"`
	Nonce int    `json:"nonce"`
}

// Unmarshal is used to unmarshal payloads
func Unmarshal(in interface{}, out interface{}) error {
	if raw, err := json.Marshal(in); err != nil {
		return err
	} else {
		return json.Unmarshal(raw, &out)
	}
}
