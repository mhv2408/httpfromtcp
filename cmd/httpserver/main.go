package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

const badResponse = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
const internalServerError = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
const successResponse = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

func Handler(w *response.Writer, r *request.Request){
	
	switch r.RequestLine.RequestTarget{
	case "/yourproblem":
		
		err := w.WriteStatusLine(response.BadRequest)
		defaultHeaders := response.GetDefaultHeaders(len(badResponse))
		defaultHeaders.SetDefaultHeader("Content-Type", "text/html")

		if err != nil{
			fmt.Println(err)
		}
		err = w.WriteHeaders(defaultHeaders)
		w.WriteBody([]byte(badResponse))

	case "/myproblem":
		err := w.WriteStatusLine(response.InternalServerError)
		defaultHeaders := response.GetDefaultHeaders(len(internalServerError))
		defaultHeaders.SetDefaultHeader("Content-Type", "text/html")

		if err != nil{
			fmt.Println(err)
		}
		err = w.WriteHeaders(defaultHeaders)
		_, err = w.WriteBody([]byte(internalServerError))

	default:
		err := w.WriteStatusLine(response.Success)
		defaultHeaders := response.GetDefaultHeaders(len(successResponse))
		defaultHeaders.SetDefaultHeader("Content-Type", "text/html")

		if err != nil{
			fmt.Println(err)
		}
		err = w.WriteHeaders(defaultHeaders)
		w.WriteBody([]byte(successResponse))
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