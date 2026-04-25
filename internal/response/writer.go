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
		return fmt.Errorf("incorrect writer state: %d", w.writerState)
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
		return fmt.Errorf("incorrect writer state: %d", w.writerState)
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
		return 0, fmt.Errorf("incorrect writer state: %d", w.writerState)
	}
	w.Body = append(w.Body, p...)
	return len(p), nil

}

func (w *Writer) WriteChunkedBody(p []byte) (int, error){
	if w.writerState != Body{
		return 0, fmt.Errorf("incorrect writer state: %d", w.writerState)
	}
	chunk_size := len(p)
	hex_n := fmt.Sprintf("%x", chunk_size)
	chunk_body := []byte(hex_n + "\r\n")
	w.Body = append(w.Body, chunk_body...)
	w.Body = append(w.Body, p...)
	w.Body = append(w.Body, []byte("\r\n")...)
	return len(chunk_body) + chunk_size + 2, nil // 2 for the last "\r\n"
}
func (w *Writer) WriteChunkedBodyDone() (int, error){
	if w.writerState != Body{
		return 0, fmt.Errorf("incorrect writer state: %d", w.writerState)
	}
	final_body := []byte("0\r\n\r\n")
	w.Body = append(w.Body, final_body...)
	return len(final_body), nil
}