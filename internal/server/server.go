package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)
type ServerState int
const (
	running ServerState=iota
	stopped
)
type Server struct{
	listner net.Listener
	clossed atomic.Bool
}

func Serve(port int)(*Server, error){
	lsnr, err := net.Listen("tcp",  fmt.Sprintf(":%d",port))
	if err != nil{
		return nil, err
	}
	
	newServer := &Server{
		listner: lsnr,
	}

	go newServer.listen()
	
	return newServer, nil
}

func (s *Server) Close() error{
	s.clossed.Store(true)
	if s.listner != nil{
		return s.listner.Close()
	}
	return nil
}

func (s *Server) listen(){
	for {
		conn, err := s.listner.Accept()
		
		if err!=nil{
			if s.clossed.Load(){
				return
			}
			log.Printf("unable to accept the connection: %v", err)
			continue
		}
		go s.Handle(conn)

	}
	
}

func (s *Server) Handle(conn net.Conn){
	response := "HTTP/1.1 200 OK" + "\r\n" +
"Content-Type: text/plain" + "\r\n" +
"Content-Length: 13" +
"\r\n\r\n"+
"Hello World!\n"
	_, err := conn.Write([]byte(response))
	if err != nil{
		log.Fatalf("unable to write response from conn: %v", err)
	}

	defer conn.Close()
}