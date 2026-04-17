package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type RequestStatus int
const (
	Initialized RequestStatus = iota
	Done
)
type Request struct{
	RequestLine RequestLine
	RequestStatus RequestStatus 
}
type RequestLine struct{
	HttpVersion string
	RequestTarget string
	Method string
}
const bufferSize = 8
func (r *Request) parse(data []byte) (int, error){

	if r.RequestStatus==Done{
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}

	if r.RequestStatus != Initialized{
		return 0, fmt.Errorf("error: unknown state")
	}
	
	n, err := parseRequestLine(data)
	if n==0 { // more bytes required to parse
		return n, err
	}
	requestLine, err := requestLineFromString(string(data[:n]))
	if err!=nil{
		fmt.Println("err is not nil: ", err)
		return 0, nil
	}
	fmt.Println(requestLine.HttpVersion, requestLine.Method, requestLine.RequestTarget)
	r.RequestLine = *requestLine
	r.RequestStatus = Done
	return n, nil
}
const crlf = "\r\n"
func parseRequestLine(message_bytes []byte)(int, error){

	idx := bytes.Index(message_bytes, []byte(crlf))

	if idx == -1{
		return 0, nil
	}

	return idx+len(crlf), nil // total number of bytes consumed
}
func requestLineFromString(str string) (*RequestLine, error){
	request_line_parts := strings.Split(str, " ")
	if len(request_line_parts) != 3{
		return nil, fmt.Errorf("invalid number of parts in request line: %d parts received", len(request_line_parts))
	}

	method, requestTarget, httpName := request_line_parts[0], request_line_parts[1], request_line_parts[2]

	// verify method

	for _, c := range method{
		if c<'A' || c>'Z'{
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	// verify http Name
	versionParts := strings.Split(httpName, "/")
	if len(versionParts) !=2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}
	if versionParts[0] != "HTTP"{
		return nil, fmt.Errorf("unrecognized HTTP version: %s", versionParts[0])
	}
	if versionParts[1] != "1.1"{
		return nil, fmt.Errorf("unrecognized HTTP version: %s", versionParts[1])
	}
	return &RequestLine{
		Method: method,
		RequestTarget: requestTarget,
		HttpVersion: versionParts[1],
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error){
	
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	newRequest := &Request{
		RequestStatus: Initialized,
	}
	for {
		// check if the buffer is full
		if readToIndex == cap(buf) {
			// grow the buffer 2 fold
			curr_buff_size := len(buf)
			new_buf := make([]byte,curr_buff_size*2)
			copy(new_buf,buf)
			buf = new_buf
		}
		bytes_read, err := reader.Read(buf[readToIndex:])
		if err != nil{
			if errors.Is(err, io.EOF){
				newRequest.RequestStatus = Done
				break
			}
			log.Fatalf("unable to read from buffer: %s", err.Error())
		}
		readToIndex += bytes_read // number of bytes actually read
		// calling the parse function
		bytes_parsed, err := newRequest.parse(buf[:readToIndex])
		if err!=nil{
			fmt.Printf("Error while parsing the request: %s",err.Error())
			os.Exit(1)
		}
		copy(buf, buf[bytes_parsed:])
		readToIndex -= bytes_parsed // number of bytes parsed

	}
	return newRequest, nil
}


func ProcessOnce(reader io.Reader, buf []byte )([]byte, int, error) {
	// check if the buffer is full
	newRequest := &Request{
		RequestStatus: Initialized,
	}
	readToIndex := 0
		if len(buf) == cap(buf) {
			// grow the buffer 2 fold
			curr_buff_size := len(buf)
			new_buf := make([]byte,curr_buff_size*2)
			copy(new_buf,buf)
			buf = new_buf
		}
		bytes_read, err := reader.Read(buf[readToIndex:])
		if err != nil{
			if errors.Is(err, io.EOF){
				newRequest.RequestStatus = Done
				return buf, readToIndex, nil
			}
			return nil, 0, err 
		}
		readToIndex += bytes_read // number of bytes actually read
		// calling the parse function
		bytes_parsed, err := newRequest.parse(buf)
		if err!=nil{
			return nil, 0, err
		}
		copy(buf, buf[bytes_parsed:])
		readToIndex -= bytes_parsed // number of bytes parsed

		return buf, readToIndex, nil
}