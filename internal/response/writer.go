package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
)
type WriterState int
const (
	StatusLine WriterState = iota
	Headers 
	Body
)

type Writer struct{
	StatusCode StatusCode
	Headers headers.Headers
	Body []byte
	writerState WriterState
}


func (w *Writer)WriteStatusLine(statusCode StatusCode) error{
	if w.writerState != StatusLine{
		return fmt.Errorf("incorrect writer state: %s", w.writerState)
	}
	response := ""
	switch statusCode{
	case Success:
		response = "HTTP/1.1 200 OK\r\n"
		w.Body = append(w.Body, []byte(response)...)
		
	case BadRequest:
		response = "HTTP/1.1 400 Bad Request\r\n"
		w.Body = append(w.Body, []byte(response)...)
	case InternalServerError:
		response = "HTTP/1.1 500 Internal Server Error\r\n"
		w.Body = append(w.Body, []byte(response)...)
	default:
		response = fmt.Sprintf("HTTP/1.1 %d\r\n", statusCode)
		w.Body = append(w.Body, []byte(response)...)
	}
	w.writerState = Headers
	return nil
}

func (w *Writer)WriteHeaders(headers headers.Headers)error{
	if w.writerState != Headers{
		return fmt.Errorf("incorrect writer state: %s", w.writerState)
	}
	for key, val := range headers{
		resp := fmt.Sprintf("%s: %s\r\n", key, val)
		w.Body = append(w.Body, []byte(resp)...)
	}
	w.Body = append(w.Body, []byte("\r\n")...)
	w.writerState = Body
	return nil
}

func (w *Writer)WriteBody(p []byte) (int, error){
	if w.writerState != Body{
		return 0, fmt.Errorf("incorrect writer state: %s", w.writerState)
	}
	w.Body = append(w.Body, p...)
	return len(p), nil

}