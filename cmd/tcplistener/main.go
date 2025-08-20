package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	stringchan := make(chan string, 1)

	go func() {
		var resultstr string
		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err != nil {
				break
			}
			i := bytes.IndexByte(data, '\n')
			if i != -1 {
				resultstr += string(data[:i])
				stringchan <- resultstr
				resultstr = ""
				if len(data) > i {
					resultstr += string(data[i+1:])
				}
			} else {
				resultstr += string(data)
			}
		}
		if len(resultstr) != 0 {
			stringchan <- resultstr
		}
		close(stringchan)
	}()

	return stringchan
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
}
