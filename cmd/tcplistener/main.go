package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const inputFilePath = "message.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {
	
	line_channel := make(chan string)
	

	go func() {
		defer f.Close()
		defer close(line_channel)
		currentLine := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err!=nil{
				if currentLine != ""{
					line_channel <- currentLine
				}
				if errors.Is(err, io.EOF){
					break
				}
				fmt.Printf("error: %s\n", err.Error())
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")

			for i :=0; i<len(parts)-1; i++{
				line_channel <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]

		}
		
	}()
	return line_channel

}

const port = ":42069"
func main() {
	
	// start a tcp connection
	ln, err := net.Listen("tcp", port)
	if err!=nil{
		// handle error
		log.Fatalf("cannot listen to the tcp connection: %s", err)
	}
	defer ln.Close()
	fmt.Println("Listening to TCP traffic on", port)
	for {
		conn, err := ln.Accept()
		
		if err!=nil{
			log.Fatalf("Unable to accept the connection: %s", err)
		}
		fmt.Println("Accepting connection from", conn.RemoteAddr())
		result := getLinesChannel(conn)
		for line := range result{
			fmt.Println(line)

		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed!")
		for line := range result{
			fmt.Println(line) // remainig data after cloising the connection
		}
	}
	
	
	
}
