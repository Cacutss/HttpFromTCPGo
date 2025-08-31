package request

import (
	"bytes"
	"errors"
	"io"
)

type RequestLine struct {
	HttpVersion   []byte
	RequestTarget []byte
	Method        []byte
}

type Request struct {
	RequestLine RequestLine
	Status      int
}

const (
	INITIALIZED = iota
	DONE
)

func (r *Request) Parse(data []byte) (int, error) {
	numOfBytesRead := 0
	if r.Status == INITIALIZED {
		reqline, read, err := parseRequestLine(data)
		if err != nil {
			return numOfBytesRead, err
		}
		numOfBytesRead += read
		if reqline != nil {
			r.Status = DONE
			r.RequestLine = *reqline
		}
	}
	return numOfBytesRead, nil
}

func (r *Request) Done() bool {
	return r.Status == DONE
}

var NEWLINE = "\r\n"
var HTTP = "HTTP"
var HTTPVERSION = "1.1"
var ERR_READING_REQUEST = errors.New("Error reading request body")
var ERR_MALFORMED_REQUEST_LINE = errors.New("Malformed request-line")
var ERR_NO_HTTP_VERSION = errors.New("No http version or no / provided in it")

func NewRequest() *Request {
	return &Request{
		Status: INITIALIZED,
	}
}

func parseRequestLine(ReqLine []byte) (*RequestLine, int, error) {
	consumed := len(ReqLine)
	i := bytes.Index(ReqLine, []byte(NEWLINE))
	if i == -1 {
		return nil, 0, nil
	}
	splits := bytes.Split(ReqLine[:i], []byte(" "))
	if len(splits) < 3 {
		return nil, consumed, ERR_MALFORMED_REQUEST_LINE
	}
	httpPart := bytes.Split(splits[2], []byte("/"))
	if len(httpPart) != 2 || string(httpPart[0]) != HTTP || string(httpPart[1]) != HTTPVERSION {
		return nil, consumed, ERR_NO_HTTP_VERSION
	}
	return &RequestLine{
		Method:        splits[0],
		RequestTarget: splits[1],
		HttpVersion:   httpPart[1],
	}, consumed, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, 32)
	readData := 0
	parsedData := 0
	request := NewRequest()
	for !request.Done() {
		readBytes, err := reader.Read(buff[readData:])
		if err != nil {
			return nil, err
		}
		readData += readBytes
		parsedBytes, err := request.Parse(buff[parsedData:readData])
		if err != nil {
			return nil, err
		}
		parsedData += parsedBytes
	}
	return request, nil
}
