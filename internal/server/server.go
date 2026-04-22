package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
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

type Handler func (w io.Writer, r *request.Request) *HandlerError

type HandlerError struct{
	StatusCode response.StatusCode
	Message string 
}
func writeError(w io.Writer, he *HandlerError){
	response.WriteStatusLine(w, he.StatusCode)

	default_headers := response.GetDefaultHeaders(len(he.Message))

	response.WriteHeaders(w, default_headers)
	w.Write([]byte(he.Message))
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
		writeError(conn, &HandlerError{StatusCode: 500, Message: err.Error()})
	}
	// create a new bytes buffer
	writer := bytes.NewBuffer(nil)
	he := s.handler(writer, request)
	if he != nil{
		writeError(writer, he)
		conn.Write(writer.Bytes())
		return
	}
	b := writer.Bytes()
	response.WriteStatusLine(conn, response.Success)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)

	defer conn.Close()
}