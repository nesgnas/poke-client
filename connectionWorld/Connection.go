package connectionWorld

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func ConnecWorld() {
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
		} else if strings.HasPrefix(line, "EXIT") {
			_, err := conn.Write([]byte(line + "\n"))
			if err != nil {
				fmt.Printf("Error sending to server: %v\n", err)
				return
			}

		} else if strings.HasPrefix(line, "GET") {
			_, err := conn.Write([]byte(line + "\n"))
			if err != nil {
				fmt.Printf("Error sending to server: %v\n", err)
				return
			}

			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				text := scanner.Text()
				if isJSON(text) {
					fmt.Println("Received JSON:")
					var prettyJSON bytes.Buffer
					if err := json.Indent(&prettyJSON, []byte(text), "", "  "); err != nil {
						fmt.Println("Error parsing JSON:", err)
						continue
					}

					parts := strings.SplitN(line, " ", 3)

					err := ioutil.WriteFile(parts[1], prettyJSON.Bytes(), 0644)
					if err != nil {
						fmt.Printf("Error writing to file: %v\n", err)
						return
					}
					fmt.Println("JSON data saved to output.json")

				} else {
					fmt.Println("Received:", text)
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading from server: %v\n", err)
			}

		} else {
			fmt.Println("Invalid command. Use SUBSCRIBE <channel>, PUBLISH <channel> <message>, UNSUBSCRIBE <channel>, or SHOWLIST.")
			fmt.Println(" Use GET <fileName>, PUBLISH <channel> <message>, UNSUBSCRIBE <channel>, or SHOWLIST.")
		}
	}
	if err := inputScanner.Err(); err != nil {
		fmt.Printf("Error reading from input: %v\n", err)
	}
}

func isJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}
