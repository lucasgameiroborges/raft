package util

import (
	"fmt"
	"os"
)

func WriteFile(fileName string, msg string) {
	completeFilePath := "/node/cluster-data/log/" + fileName + "-" + os.Getenv("NODE_ID") + ".txt"
	file, err := os.OpenFile(completeFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(msg + "\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}