package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct{
	RequestLine RequestLine
}
type RequestLine struct{
	HttpVersion string
	RequestTarget string
	Method string
}
const crlf = "\r\n"
func parseRequestLine(message_bytes []byte)(*RequestLine, error){

	idx := bytes.Index(message_bytes, []byte(crlf))

	if idx == -1{
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}
	requestLineText := string(message_bytes[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	
	if err != nil{
		return nil, err
	}

	return requestLine, nil
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
		return nil, fmt.Errorf("unrecognized HTTP version: %s", versionParts[0])
	}
	return &RequestLine{
		Method: method,
		RequestTarget: requestTarget,
		HttpVersion: versionParts[1],
	}, nil
}
func RequestFromReader(reader io.Reader) (*Request, error){
	
	message_bytes, err := io.ReadAll(reader)

	if err!=nil{
		return nil, fmt.Errorf("cannot read the message from reader: %s", err.Error())
	}
	requestLine, err := parseRequestLine(message_bytes)
	if err != nil{
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}