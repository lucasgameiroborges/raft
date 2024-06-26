package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	Acceptor "github.com/paxos/src/multi-paxos/acceptor"
	Learner "github.com/paxos/src/multi-paxos/learner"
	Proposer "github.com/paxos/src/multi-paxos/proposer"
	"github.com/paxos/src/multi-paxos/variable"
	"github.com/paxos/src/pkg/model/message"
	"github.com/paxos/src/pkg/shared/constant"
	"github.com/paxos/src/pkg/shared/util"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func handleAdd(w http.ResponseWriter, r *http.Request) {
	tries := 0
	initialLog := variable.LogSize
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
		Source:  "nowhere",
		Type:    constant.REQUEST,
		Payload: message.Request{Value: str},
	}

	// fmt.Fprintf(w, "sending a message... \n")
	prop := os.Getenv("NODE_ID") + ".raft000.raft-k8s.svc.cluster.local" + ":9001"
	err = util.SendMessage(message1, prop)
	if err != nil {
		// fmt.Fprintf(w, "Failed! %s \n", err.Error())
		return
	}
	// fmt.Fprintf(w, "message sent! \n")

	// Wait some time for Paxos to reach consensus
	//time.Sleep(time.Second / 5)
	startTime := time.Now()
	elapsedTime := time.Since(startTime)
	for {
		if variable.LogSize > initialLog {
			fmt.Fprintf(w, "Value successfuly stored in log\n")
			variable.Round++
			fmt.Fprintf(w, "New round: %s\n", variable.Round)
			break
		}
		elapsedTime = time.Since(startTime)
		if elapsedTime > 1*time.Second {
			http.Error(w, "fail", http.StatusBadRequest)
			return 
		}
		tries = tries + 1
	}
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
		Source:  "nowhere2",
		Type:    constant.REQUEST,
		Payload: message.Request{Value: deleted_value},
	}

	statmsg := fmt.Sprintf("client ->> proposer %s:9001: Request: %v", os.Getenv("NODE_ID"), str)
	util.WriteFile("status", statmsg)
	statmsg = fmt.Sprintf("Note over client,proposer %s:9001: Initialize round 1\n", os.Getenv("NODE_ID"))
	util.WriteFile("status", statmsg)
	prop := os.Getenv("NODE_ID") + ".raft000.raft-k8s.svc.cluster.local" + ":9001"
	err = util.SendMessage(message1, prop)
	if err != nil {
		fmt.Fprintf(w, "Failed! %s \n", err.Error())
		return
	}

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second / 10)

	fmt.Fprintf(w, "Received POST request with body: %s \n", str)
}

func handleLog(w http.ResponseWriter, r *http.Request) {
	// Read the contents of the file
	completeFilePath := "/node/cluster-data/log/" + "log-" + os.Getenv("NODE_ID") + ".txt"
	content, err := ioutil.ReadFile(completeFilePath)
	if err != nil {
		fmt.Fprintf(w, "Error reading file: %v \n", err)
		return
	}

	// Print the contents of the file
	fmt.Fprintf(w, "contents of log file from %s: \n", os.Getenv("NODE_ID"))
	fmt.Fprintf(w, string(content))
}

func handlePrint(w http.ResponseWriter, r *http.Request) {
	// Open the file
	completeFilePath := "/node/cluster-data/log/" + "log-" + os.Getenv("NODE_ID") + ".txt"
	file, err := os.Open(completeFilePath)
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
	fmt.Fprintf(w, "Dictionary stored in %s: \n\n", os.Getenv("NODE_ID"))
	fmt.Fprintf(w, string(prettyJSON))
	fmt.Fprintf(w, "\n")
}

func wipeLogFolder() {
    // Specify the directory path
    directory := "/node/cluster-data/log/"

    // List all files in the directory
    files, err := filepath.Glob(filepath.Join(directory, "*"))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Delete each file
    for _, file := range files {
        err := os.Remove(file)
        if err != nil {
            fmt.Println("Error deleting file:", err)
        }
    }
}

// Initializes an instance of Multi-Paxos with several nodes: one proposer, three acceptors, and one learner
// The instance simulates a scenario where a client submits two requests to the same proposer with different values
// The requests are processed in rounds by the network, and the network arrives to a consensus on both values in their
// respective rounds
func main() {
	wipeLogFolder()
	variable.Round = 1
	variable.LogSize = 0

	var acceptors []string
	var learners []string
	var targetServiceIP string
	var acc string
	var lea string

	for i := 0; i < 3; i++ {
		targetServiceIP = fmt.Sprintf("raft000-%d.raft000.raft-k8s.svc.cluster.local", i)
		acc = targetServiceIP + ":9002"
		lea = targetServiceIP + ":9003"
		acceptors = append(acceptors, acc)
		learners = append(learners, lea)
	}
	proposer := os.Getenv("NODE_ID") + ".raft000.raft-k8s.svc.cluster.local" + ":9001"
	acceptor := os.Getenv("NODE_ID") + ".raft000.raft-k8s.svc.cluster.local" + ":9002"
	learner := os.Getenv("NODE_ID") + ".raft000.raft-k8s.svc.cluster.local" + ":9003"

	go Proposer.Activate(proposer, acceptors)
	go Acceptor.Activate(acceptor, learners)
	go Learner.Activate(learner)

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
