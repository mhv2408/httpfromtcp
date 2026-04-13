package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)
const host = "localhost"
const port = "42069"
func main(){
	udp_addr, err := net.ResolveUDPAddr("udp",fmt.Sprintf("%s:%s",host, port))

	if err!= nil{
		log.Fatalf("Unable to resolve udp address: %s", err.Error())
	}
	fmt.Println("Address Resolved: ",udp_addr.AddrPort())
	// prepare the connection
	udp_conn, err := net.DialUDP(udp_addr.Network(), nil, udp_addr)

	if err!=nil{
		log.Fatalf("Unable to establish the udp connection: %s", err.Error())
	}
	defer udp_conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		message, err := reader.ReadString('\n')
		if err!=nil{
			log.Fatalf("Unable to read the string from the buffer: %s", err.Error())
		}
		_, err = udp_conn.Write([]byte(message))
		if err != nil{
			log.Fatalf("cannot write to udp connection: %s", err.Error())
		}
		fmt.Printf("Message sent: %s",message)
	}
}