package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
)

type RequestStatus int
const (
	Initialized RequestStatus = iota
	Done
	requestStateParsingHeaders
)
type Request struct{
	RequestLine RequestLine
	RequestStatus RequestStatus 
	Headers headers.Headers
}
type RequestLine struct{
	HttpVersion string
	RequestTarget string
	Method string
}
const bufferSize = 8
func (r *Request)parseSingle(data []byte)(int, error){
	switch r.RequestStatus{
	case Initialized:

		requestLine, parsedBytes, err := parseRequestLine(data)

		if err != nil{
			// something actually went wrong
			return 0, err
		}
		if parsedBytes == 0 {
			// need more bytes
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.RequestStatus = requestStateParsingHeaders
		fmt.Println("Parsed bytes from request line: ", parsedBytes)
		return parsedBytes, nil
	case requestStateParsingHeaders:
		fmt.Println("data inside headers: ", string(data))
		n, done, err := r.Headers.Parse(data)
		if err!=nil{
			return 0, err
		}
		if n==0{
			return 0, nil
		}
		if done{
			r.RequestStatus = Done
			return n, nil
		}
		return n, nil

	case Done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
func (r *Request) parse(data []byte) (int, error){
	totalBytesParsed := 0
	for r.RequestStatus != Done{
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil{
			return 0, err
		}
		if n == 0{
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
		fmt.Println("totalBytesParsed: ", n)
	}
	return totalBytesParsed, nil
}
const crlf = "\r\n"
func parseRequestLine(message_bytes []byte)(*RequestLine, int, error){

	idx := bytes.Index(message_bytes, []byte(crlf))

	if idx == -1{
		return nil, 0, nil
	}
	requestLine, err := requestLineFromString(string(message_bytes[:idx]))

	if err != nil{
		return nil, 0, err
	}

	return requestLine, idx+len(crlf), nil // total number of bytes consumed
}
func requestLineFromString(str string) (*RequestLine, error){
	request_line_parts := strings.Split(str, " ")
	if len(request_line_parts) != 3{
		return nil, fmt.Errorf("invalid number of parts in request line: %d parts received", len(request_line_parts))
	}

	method, requestTarget, httpName := request_line_parts[0], request_line_parts[1], request_line_parts[2]

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
	
	buf := make([]byte, bufferSize)
	readToIndex := 0
	newRequest := &Request{
		RequestStatus: Initialized,
		Headers: headers.NewHeaders(),
	}
	for newRequest.RequestStatus!=Done{
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
				if newRequest.RequestStatus != Done{
					return nil, fmt.Errorf("incomplete request")
				}
				break
			}
			return nil, err
		}
		readToIndex += bytes_read // number of bytes actually read
		// calling the parse function
		bytes_parsed, err := newRequest.parse(buf[:readToIndex])
		if err!=nil{
			return nil, err
		}
		copy(buf, buf[bytes_parsed:])
		readToIndex -= bytes_parsed // number of bytes parsed

	}
	return newRequest, nil
}