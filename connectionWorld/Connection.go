package connectionWorld

import (
	"bufio"
	"bytes"
	"client/bluePrint"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

var LocalAddress string = "[::1]:58065"
var (
	PlayerX, PlayerY int
)

func ConnecWorld() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	LocalAddress = conn.LocalAddr().String()
	remoteAddr := conn.RemoteAddr().String()

	fmt.Println("Local address:", conn.LocalAddr().String())
	fmt.Println("Remote address:", remoteAddr)

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			receivedMessage := scanner.Text()
			fmt.Println("Received:", receivedMessage)

			// Switch based on whether the message starts with "REPEAT "
			switch {
			case strings.HasPrefix(receivedMessage, "REPEAT "):
				// Extract the message after "REPEAT "
				messageToRepeat := strings.TrimPrefix(receivedMessage, "REPEAT ")

				// Switch based on the command type
				switch {
				case strings.HasPrefix(messageToRepeat, "GET "):
					_, err := conn.Write([]byte(messageToRepeat + "\n"))
					if err != nil {
						fmt.Printf("Error sending to server: %v\n", err)
						return
					}
					// If the message is a "GET" command, wait for JSON response
					fmt.Println(messageToRepeat)
					fmt.Println("hehehehe")
					waitForJSONResponse(conn, messageToRepeat)
					bluePrint.ReadUser("clients.json")
				default:
					// For other messages, directly send them back to the server
					_, err := conn.Write([]byte(messageToRepeat + "\n"))
					if err != nil {
						fmt.Printf("Error sending to server: %v\n", err)
						return
					}
				}
			}
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

func waitForJSONResponse(conn net.Conn, fileName string) {
	fmt.Println("doheheeh")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		if isJSON(text) {
			fmt.Println("Received JSON:")
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(text), "", "  "); err != nil {
				fmt.Println("Error parsing JSON:", err)
				continue
			}

			parts := strings.SplitN(fileName, " ", 3)

			err := ioutil.WriteFile(parts[1], prettyJSON.Bytes(), 0644)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
			fmt.Println("JSON data saved to" + parts[1])

		} else {
			fmt.Println("Received:", text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
	}
}
