package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"bufio"
)

const (
	connHost = "192.168.1.9"
	connPort = "6666"
	connType = "tcp"
)

var output string = ""

func updateScreen(msg string) {
	output += string(msg)
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Printf("%s\n", output)
}

func main() {
	updateScreen("")

	conn, err := net.Dial(connType, connHost + ":" + connPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(-1)
	}

	defer conn.Close()
	go listener(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		msg, e := reader.ReadString('\n')
		if e != nil {
			fmt.Println("Failed to read message:", e.Error())
		}
		if len(msg) > 0 {
			updateScreen(msg)

			switch msg {
				case "/exit\n":
					fmt.Println("Disconnecting...")
					os.Exit(1)
				case "/clear\n":
					output = ""
					updateScreen("")
				default:
					_, e = conn.Write([]byte(msg))
					if e != nil {
						fmt.Println("Failed to send message:", e.Error())
					}
			}
		}
	}
}

func listener(conn net.Conn) {
	buf := make([]byte, 1024)

	for {
		n := 0
		var e error

		for n <= 0 {
			n, e = conn.Read(buf)
			if e != nil {
				if e.Error() == "EOF" {
					updateScreen("Server Disconnected.\n")
					os.Exit(0)
				}
			}
		}

		updateScreen(string(buf[:n]))
	}
}
