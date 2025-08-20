package request

import (
	"errors"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var NEWLINE = "\r\n"
var ERR_READING_REQUEST = errors.New("Error reading request body")
var ERR_MALFORMED_REQUEST_LINE = errors.New("Malformed request-line")

func parseRequestLine(ReqLine string) (*RequestLine, int, error) {
	i := strings.Index(ReqLine, NEWLINE)
	if i == -1 {
		return nil, 0, nil
	}
	firstLine := ReqLine[:i]
	splits := strings.Split(firstLine, " ")
	if len(splits) < 3 {
		return nil,
	}
	httpPart := strings.Split(splits[2], "/")
	if len(httpPart) != 2 || httpPart[0] != "HTTP" || httpPart[1] != "1.1" {
		return nil, ERR_MALFORMED_REQUEST_LINE
	}

	return &RequestLine{
		Method:        splits[0],
		RequestTarget: splits[1],
		HttpVersion:   httpPart[1],
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, ERR_READING_REQUEST
	}
	reqLine, readBytes, err := parseRequestLine(string(body))
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *reqLine,
	}, nil
}
