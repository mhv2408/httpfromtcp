package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const inputFile = "message.txt"

func getLinedChannel(c net.Conn) <-chan string {
	currentLineContents := ""
	res := make(chan string)

	go func() {
		defer close(res)
		for {
			buffer := make([]byte, 8, 8)
			n, err := c.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					res <- fmt.Sprintf("%s\n", currentLineContents)
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				res <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()

	return res
}

func main() {

	tcp_listner, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatal("Unable to open a TCP connection: ", err)
	}
	defer tcp_listner.Close()
	fmt.Println("Connected on Port:42069")

	for {
		conn, err := tcp_listner.Accept()
		if err != nil {
			log.Fatal("Cannot form a TCP connection: ", err)
		}
		fmt.Println("Connection has been accepted")
		tcp_channel := getLinedChannel(conn)

		for val := range tcp_channel {
			fmt.Printf("%s\n", val)
		}
		fmt.Println("Channel is Closed")
	}

}
