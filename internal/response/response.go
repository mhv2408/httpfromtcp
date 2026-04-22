package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type StatusCode int
const (
	Success StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error{
	response := ""
	switch statusCode{
	
	case Success:
		response = "HTTP/1.1 200 OK\r\n"
		_, err := w.Write([]byte(response))
		if err != nil{
			return err
		}
	case BadRequest:
		response = "HTTP/1.1 400 Bad Request\r\n"
		_, err := w.Write([]byte(response))
		if err != nil{
			return err
		}
	case InternalServerError:
		response = "HTTP/1.1 500 Internal Server Error\r\n"
		_, err := w.Write([]byte(response))
		if err != nil{
			return err
		}
	default:
		response = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
		_, err := w.Write([]byte(response))
		if err != nil{
			return err
		}
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers{
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprint(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
func WriteHeaders(w io.Writer, headers headers.Headers)error{
	for key, val := range headers{
		resp := fmt.Sprintf("%s: %s\r\n", key, val)
		_, err := w.Write([]byte(resp))
		if err != nil{
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}

