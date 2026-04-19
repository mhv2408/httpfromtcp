package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string] string
var specialChars = map[rune]bool{'!':true, '#':true, '$':true, '%':true, '&':true, '|':true, '~':true, '^':true, '*':true, '`':true, '_':true, '+':true, '\'':true, '-':true, '.':true,}

// !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~

// constructor to make headers
func NewHeaders() Headers {
	return make(Headers)
}
const crlf = "\r\n"
func (h Headers) Parse(data []byte) (n int, done bool, err error){
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1{
		// not enough data need more
		return 0, false, fmt.Errorf("cannot find crlf: Need more data")
	}
	// found the crlf
	if idx == 0{ // no headers
		return 2, true, nil
	}
	// find the colon (:)
	colon_idx := bytes.Index(data[:idx], []byte(":"))
	if colon_idx==-1{
		return 0, false, fmt.Errorf("can't find seperator colon ':'")
	}
	// split by colon(:)
	parts := bytes.SplitN(data[:idx], []byte(":"), 2)

	rawFieldName, rawFieldValue := string(parts[0]), string(parts[1])
	// check the spacing constraint for field-name
	if strings.TrimSpace(rawFieldName) != rawFieldName{
		return 0, false, fmt.Errorf("spaces before fieldName or betwixt colon")
	}
	fieldName, fieldValue := strings.TrimSpace(rawFieldName), strings.TrimSpace(rawFieldValue)
	// check if fieldName comtains any special characters
	for _, c := range fieldName{
		if !unicode.IsNumber(c) && !unicode.IsLetter(c) && !specialChars[c]{
			fmt.Println("inside error unicoe")
			return 0, false, fmt.Errorf("Invalid filed name: %c", c)
		}
	}
	
	h.Set(fieldName, fieldValue)
	return idx + 2, false, nil

}
func (h Headers)Set(key, value string){
	key = strings.ToLower(key)
	if h.Get(key)==""{
		h[key] = value
		return
	}
	h[key] = h.Get(key) + ", " + value
}

func (h Headers)Get(key string) string{
	key = strings.ToLower((key))
	
	return h[key]
}

