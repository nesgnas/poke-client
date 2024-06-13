package connectionWorld

import (
	"bufio"
	"bytes"
	"client/bluePrint"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var LocalAddress string
var (
	Conn             net.Conn
	PlayerX, PlayerY int
	CHANNEL_HI       bool = false
)

func ConnectWorld(conn net.Conn) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	remoteAddr := conn.RemoteAddr().String()

	fmt.Println("Local address:", LocalAddress)
	fmt.Println("Remote address:", remoteAddr)

	initiateConnectionWorld(conn)

	go listenForMessages(conn)

	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		line := inputScanner.Text()
		switch {
		case strings.HasPrefix(line, "SUBSCRIBE"), strings.HasPrefix(line, "PUBLISH"), strings.HasPrefix(line, "UNSUBSCRIBE"), strings.HasPrefix(line, "SHOWLIST"):
			SendMessage(conn, line)
		case strings.HasPrefix(line, "EXIT"):
			SendMessage(conn, line)
			return
		case strings.HasPrefix(line, "GET"):
			handleGetCommand(conn, line)
		default:
			fmt.Println("Invalid command. Use SUBSCRIBE <channel>, PUBLISH <channel> <message>, UNSUBSCRIBE <channel>, SHOWLIST, or GET <fileName>.")
		}
	}
	if err := inputScanner.Err(); err != nil {
		fmt.Printf("Error reading from input: %v\n", err)
	}
}

func listenForMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		receivedMessage := scanner.Text()
		fmt.Println("Received:", receivedMessage)

		switch {
		case strings.HasPrefix(receivedMessage, "REPEAT "):
			messageToRepeat := strings.TrimPrefix(receivedMessage, "REPEAT ")

			switch {
			case strings.HasPrefix(messageToRepeat, "GET "):
				SendMessage(conn, messageToRepeat)
				waitForJSONResponse(conn, messageToRepeat)
				bluePrint.ReadUser("clients.json")
			default:
				SendMessage(conn, messageToRepeat)
			}
		case strings.HasPrefix(receivedMessage, "PUBLISH"):
			checkChannelAction(receivedMessage)
		default:
			// Handle other cases if needed
			//fmt.Println("Unhandled message:", receivedMessage)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
	}
}

func SendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Printf("Error sending to server: %v\n", err)
	}
}

func handleGetCommand(conn net.Conn, line string) {
	SendMessage(conn, line)
	fmt.Println("HANNDEL COMMAND --- ###")
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

			parts := strings.SplitN(line, " ", 2)
			err := ioutil.WriteFile(parts[1], prettyJSON.Bytes(), 0644)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
			fmt.Println("JSON data saved to", parts[1])
		} else {
			fmt.Println("Received:", text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
	}
}

func isJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

func waitForJSONResponse(conn net.Conn, fileName string) {
	scanner := bufio.NewScanner(conn)
	fmt.Println("WAITFOR JSON --- ###")
	for scanner.Scan() {
		text := scanner.Text()
		if isJSON(text) {
			fmt.Println("Received JSON:")
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(text), "", "  "); err != nil {
				fmt.Println("Error parsing JSON:", err)
				continue
			}

			parts := strings.SplitN(fileName, " ", 2)
			err := ioutil.WriteFile(parts[1], prettyJSON.Bytes(), 0644)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return
			}
			fmt.Println("JSON data saved to", parts[1])
			break
		} else {
			fmt.Println("Received:", text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
	}
}
