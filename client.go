package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println("Received:", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading from server: %v\n", err)
		}
	}()

	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		line := inputScanner.Text()
		if strings.HasPrefix(line, "SUBSCRIBE") || strings.HasPrefix(line, "PUBLISH") || strings.HasPrefix(line, "UNSUBSCRIBE") || strings.HasPrefix(line, "SHOWLIST") {
			_, err := conn.Write([]byte(line + "\n"))
			if err != nil {
				fmt.Printf("Error sending to server: %v\n", err)
				return
			}
		} else {
			fmt.Println("Invalid command. Use SUBSCRIBE <channel>, PUBLISH <channel> <message>, UNSUBSCRIBE <channel>, or SHOWLIST.")
		}
	}
	if err := inputScanner.Err(); err != nil {
		fmt.Printf("Error reading from input: %v\n", err)
	}
}
