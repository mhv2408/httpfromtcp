package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

func ProxyHandler(w *response.Writer, r *request.Request){
	// 
	resp, err := http.Get("https://httpbin.org/stream/100")
	if err!=nil{
		log.Fatalf("unable to get the response form http request: %s", err.Error())
	}
	w.WriteStatusLine(response.Success)
	defaultHeaders := response.GetDefaultHeaders(0)
	defaultHeaders.Remove("Content-Length")
	defaultHeaders.Override("Transfer-Encoding", "chunked")
	w.WriteHeaders(defaultHeaders)
	defer resp.Body.Close()
	data := make([]byte, 1024)
	for { 
		n, err := resp.Body.Read(data)
		fmt.Println("Read", n, "bytes")
		if n>0{
			_, err := w.WriteChunkedBody(data[:n])
			if err!=nil{	
				fmt.Println("Error writing chunks: ",err)
				break
			}

		}
		if err!=nil{	
			fmt.Println("Error reading response: ",err)
			break
		}
		if err == io.EOF{
			break
		}
	}
	
	_, err = w.WriteChunkedBodyDone()
	if err != nil{
		fmt.Println("error writing chunked body done", err)
	}
}

func Handler(w *response.Writer, r *request.Request){

	if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin"){
		ProxyHandler(w, r)
		return
	}
	
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