package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
)

type StatusCode int
const (
	Success StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500
)



func GetDefaultHeaders(contentLen int) headers.Headers{
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprint(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}


