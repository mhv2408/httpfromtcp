package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

const inputFilePath = "message.txt"


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
		request, err := request.RequestFromReader(conn)
		if err!=nil{
			log.Fatalf("Unable to Request from Reader: %s", err)
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", request.RequestLine.Method)
		fmt.Println("- Target:",request.RequestLine.RequestTarget)
		fmt.Println("- Version:",request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for headerName, headerValue := range request.Headers{
			fmt.Printf("- %s: %s\n", headerName, headerValue)
		}
	}
}
