package server

import (
	"fmt"
	"httpfromtcp/internal/response"
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
	/*
	response := "HTTP/1.1 200 OK" + "\r\n" +
"Content-Type: text/plain" + "\r\n" +
"Content-Length: 13" +
"\r\n\r\n"+
"Hello World!\n" */
	// Write the start line
	err := response.WriteStatusLine(conn, 200)
	if err != nil{
		log.Fatalf("error writing status line: %v", err)
	}
	// Write Body
	def_headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, def_headers)
	if err!=nil{
		log.Fatalf("error writing headers: %v", err)
	}
	defer conn.Close()
}