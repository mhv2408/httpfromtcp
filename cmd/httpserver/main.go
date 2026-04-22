package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func Handler(w io.Writer, r *request.Request)*server.HandlerError{
	switch r.RequestLine.RequestTarget{
	case "/yourproblem":

		return &server.HandlerError{
			StatusCode:response.BadRequest,
			Message:"Your problem is not my problem\n",
		}
	case "/myproblem":

		return &server.HandlerError{
			StatusCode: response.InternalServerError,
			Message:  "Woopsie, my bad\n",
		}
	default:
		body := []byte("All good, frfr\n")
		_, err := w.Write(body)
		if err!=nil{
			return &server.HandlerError{
				StatusCode: 500,
				Message: "error in writing the body message",
			}

		}
		return nil
	}

}

func main(){
	server, err := server.Serve(port, Handler)
	if err != nil{
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started in port:", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server stopped gracefully")
}