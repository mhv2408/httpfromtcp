package server

import (
	"fmt"
	"httpfromtcp/internal/request"
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
	closed atomic.Bool
	handler Handler
}

type Handler func (w *response.Writer, r *request.Request)

type HandlerError struct{
	StatusCode response.StatusCode
	Message string 
}

func Serve(port int, handler Handler)(*Server, error){
	lsnr, err := net.Listen("tcp",  fmt.Sprintf(":%d",port))
	if err != nil{
		return nil, err
	}
	
	newServer := &Server{
		listner: lsnr,
		handler: handler,
	}

	go newServer.listen()
	
	return newServer, nil
}

func (s *Server) Close() error{
	s.closed.Store(true)
	if s.listner != nil{
		return s.listner.Close()
	}
	return nil
}

func (s *Server) listen(){
	for {
		conn, err := s.listner.Accept()
		
		if err!=nil{
			if s.closed.Load(){
				return
			}
			log.Printf("unable to accept the connection: %v", err)
			continue
		}
		go s.handle(conn)

	}
	
}


func (s *Server) handle(conn net.Conn){

	// Parse the request
	request, err := request.RequestFromReader(conn)
	if err != nil{
		log.Fatalf("unable to request from reader: %s", err)
	}
	// create a new bytes buffer
	writer := &response.Writer{}
	s.handler(writer, request)
	conn.Write(writer.Body)
	defer conn.Close()
}