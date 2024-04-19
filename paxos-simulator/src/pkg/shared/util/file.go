package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func CreateNewFile(paxosType string) {
	err := ioutil.WriteFile(fmt.Sprintf("./artifacts/%s-paxos-output.txt", paxosType), []byte("sequenceDiagram\n"), 0755)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
}

func WriteToBasicFile(text string) {
	file, err := os.OpenFile("./artifacts/basic-paxos-output.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("Can't write error: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", text))
}

func WriteToMultiFile(text string) {
	file, err := os.OpenFile("./artifacts/multi-paxos-output.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("Can't write error: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", text))
}
