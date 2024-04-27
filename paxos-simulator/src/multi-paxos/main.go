package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	Acceptor "github.com/paxos/src/multi-paxos/acceptor"
	Learner "github.com/paxos/src/multi-paxos/learner"
	Proposer "github.com/paxos/src/multi-paxos/proposer"
	"github.com/paxos/src/pkg/model/message"
	"github.com/paxos/src/pkg/shared/constant"
	"github.com/paxos/src/pkg/shared/util"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func handleAdd(w http.ResponseWriter, r *http.Request) {
	// Read the body of the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	str := string(body)

	pattern := `^[a-zA-Z0-9_-]+:[a-zA-Z0-9_-]+$`
	re := regexp.MustCompile(pattern)

	if re.MatchString(str) {
		fmt.Fprintf(w, "String '%s' matches the pattern key:value\n", str)
	} else {
		fmt.Fprintf(w, "String '%s' does not match the pattern key:value, aborting\n", str)
		return 
	}

	message1 := &message.Message{
		Source:  0,
		Type:    constant.REQUEST,
		Payload: message.Request{Value: str},
	}

	fmt.Println("client ->> proposer 9001: Request: %v", str)
	fmt.Println("Note over client,proposer 9001: Initialize round 1\n")
	util.SendMessage(message1, 9001)

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second / 10)

	// Send a response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Received POST request with body: %s \n", str)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	// Read the body of the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	str := string(body)

	deleted_value := "DELETE " + str
	message1 := &message.Message{
		Source:  0,
		Type:    constant.REQUEST,
		Payload: message.Request{Value: deleted_value},
	}

	fmt.Println("client ->> proposer 9001: Request: %v", str)
	fmt.Println("Note over client,proposer 9001: Initialize round 1\n")
	util.SendMessage(message1, 9001)

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second / 10)

	// Send a response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Received POST request with body: %s \n", str)
}

func handleLog(w http.ResponseWriter, r *http.Request) {
	// Read the contents of the file
	content, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Fprintf(w, "Error reading file: %v \n", err)
		return
	}

	// Print the contents of the file
	fmt.Fprintf(w, "contents of log file from this machine: \n")
	fmt.Fprintf(w, string(content))
}

func handlePrint(w http.ResponseWriter, r *http.Request) {
	// Open the file
	file, err := os.Open("log.txt")
	if err != nil {
		fmt.Fprintf(w, "Error reading file: %v \n", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line in the file
	dictionary := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text() // Get the current line
		if strings.HasPrefix(line, "DELETE") {
			parts := strings.Split(line, " ")
			key := parts[1]
			if _, exists := dictionary[key]; exists {
				delete(dictionary, key)
			}
		} else {
			parts := strings.Split(line, ":")
			dictionary[parts[0]] = parts[1]
		}
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(w, "Error scanning file: %v \n", err)
	}

	// Pretty print the map
	prettyJSON, err := json.MarshalIndent(dictionary, "", "    ")
	if err != nil {
		fmt.Fprintf(w, "Error:", err)
		return
	}
	fmt.Fprintf(w, "Dictionary stored in log: \n\n")
	fmt.Fprintf(w, string(prettyJSON))
	fmt.Fprintf(w, "\n")
}

// Initializes an instance of Multi-Paxos with several nodes: one proposer, three acceptors, and one learner
// The instance simulates a scenario where a client submits two requests to the same proposer with different values
// The requests are processed in rounds by the network, and the network arrives to a consensus on both values in their
// respective rounds
func main() {

	fmt.Println("Initializing Multi-Paxos...")

	go Proposer.Activate(9001, []int{9002})
	go Acceptor.Activate(9002, []int{9003})
	go Learner.Activate(9003)

	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/log", handleLog)
	http.HandleFunc("/print", handlePrint)

	fmt.Println("Server Listening.")
	err := http.ListenAndServe(":7777", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed\n")
	} else if err != nil {
		fmt.Println("error starting server: %s\n", err)
		os.Exit(1)
	}
}
