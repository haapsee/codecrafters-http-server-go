package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
)

type Header struct {
    Header      string
    Content     string
}

type Request struct {
    Method      string
    Target      string
    HTTPVersion string
    Headers     map[string]string
    Body        string
}


func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

    for {
    	connection, err := l.Accept()
    	if err != nil {
    		fmt.Println("Error accepting connection: ", err.Error())
    		os.Exit(1)
    	}
    	defer connection.Close()
    	go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
    buffer := make([]byte, 2048)
    bytesRead, err := connection.Read(buffer)
    if err != nil {
        fmt.Println("Error while reading request", err.Error())
        os.Exit(1)
    }

    requeststr := string(buffer[:bytesRead])
    request := parseRequest(requeststr)

    if request.Target == "/" {
        connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
    } else if strings.HasPrefix(request.Target, "/echo/") {
        target := request.Target[6:]
        connection.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(target)) + "\r\n\r\n" + target))
    } else if request.Target == "/user-agent" {
        target := request.Headers["User-Agent"]
        connection.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(target)) + "\r\n\r\n" + target))
    } else {
        connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
    }

}

func parseRequest(requestStr string) Request {
    i := strings.Index(requestStr, " ")
    method := requestStr[:i]
    requestStr = requestStr[i+1:]

    i = strings.Index(requestStr, " ")
    target := requestStr[:i]
    requestStr = requestStr[i+1:]

    i = strings.Index(requestStr, "\r\n")
    httpVersion := requestStr[:i]
    requestStr = requestStr[i+2:]

    headers := parseHeaders(&requestStr)
    return Request{method, target, httpVersion, headers, requestStr}
}

func parseHeaders(requestStr *string) map[string]string {
    headers := make(map[string]string)
    for {
        i := strings.Index(*requestStr, "\r\n")
        if i == 0 {
            *requestStr = (*requestStr)[2:]
            break
        }
        header := (*requestStr)[:i]
        *requestStr = (*requestStr)[i+2:]
        j := strings.Index(header, ": ")
        headers[header[:j]] = header[j+2:]
    }
    return headers
}

