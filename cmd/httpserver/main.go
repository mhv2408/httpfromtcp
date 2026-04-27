package main

import (
	"crypto/sha256"
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
func VideoHandler(w *response.Writer, r *request.Request){
	
	// read the entire video
	data, err := os.ReadFile("assets/vim.mp4")
	headers := response.GetDefaultHeaders(len(data))
	headers.Override("Content-Type", "video/mp4")
	if err!=nil{
		fmt.Println("cannot read the file", err)
	}
	err = w.WriteStatusLine(response.Success)
	if err!=nil{
		fmt.Println("error writing status line:", err)
	}
	err = w.WriteHeaders(headers)
	if err!=nil{
		fmt.Println("error writing headers:", err)
	}
	_, err = w.WriteBody(data)
	if err!=nil{
		fmt.Println("error writing body:", err)
	}
}

func ProxyHandler(w *response.Writer, r *request.Request){
	target := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")
	
	endPoint := "https://httpbin.org/" + target
	resp, err := http.Get(endPoint)
	if err != nil{
		log.Fatalf("unable to perform HTTP GET: %s", err.Error())
	}
	defer resp.Body.Close()

	defaultHeaders := response.GetDefaultHeaders(0)
	defaultHeaders.Remove("Content-Length")
	defaultHeaders.Set("Trailer","X-Content-Sha256")
	defaultHeaders.Set("Trailer","X-Content-Length")
	defaultHeaders.Set("Transfer-Encoding", "chunked")
	w.WriteStatusLine(response.Success)
	w.WriteHeaders(defaultHeaders)
	data := make([]byte, 1024)
	bytesRead := 0
	var fullBody []byte
	for {
		n, err := resp.Body.Read(data)
		if err!=nil{
			if err == io.EOF{
				break
			}
			fmt.Printf("unable to read from response: %s", err.Error())
			break
		}
		_, err = w.WriteChunkedBody(data[:n])
		
		if err != nil{
			fmt.Println("unable to write the chunked body response ", err)
			break
		}
		fullBody = append(fullBody, data[:n]...)
		bytesRead += n
	}
	
	shaSum := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	defaultHeaders.SetDefaultHeader("X-Content-SHA256", shaSum)
	defaultHeaders.SetDefaultHeader("X-Content-Length", fmt.Sprint(bytesRead))
	w.WriteChunkedBodyDone()
	w.WriteTrailers(defaultHeaders)	
	
}

func Handler(w *response.Writer, r *request.Request){

	if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin"){
		ProxyHandler(w, r)
		return
	}
	if strings.HasPrefix(r.RequestLine.RequestTarget, "/video"){
		VideoHandler(w, r)
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