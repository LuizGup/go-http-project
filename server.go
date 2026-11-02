package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type HttpRequest struct {
	Method       string
	Path         string
	ProtoVersion string
	Headers      map[string]string
	ContentLength int
	BodyReader   *bufio.Reader
}

func HandleClient(connection net.Conn) {
	defer connection.Close()
	defer fmt.Println("[CONEXÃO] Encerrada", connection.RemoteAddr())

	fmt.Println("[CONEXÃO] Nova conexão de", connection.RemoteAddr())

	req := HttpRequest{
		ContentLength: 0,
		Headers: make(map[string]string),
		BodyReader: bufio.NewReader(connection),
	}

	line, err := req.BodyReader.ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler dados do cliente:", err)
		return
	}

	fmt.Printf("Linha de Requisição: %s", line)

	lineSplit := strings.Split(line, " ")
	if len(lineSplit) != 3 {
		fmt.Println("Linha de requisição inválida:", line)
		return
	}

	req.Method = lineSplit[0]
	pathRequest, _, found := strings.Cut(lineSplit[1], "?")
	if !found {
		req.Path = lineSplit[1]
	} else {
		req.Path = pathRequest
	}

	req.ProtoVersion = strings.TrimSpace(lineSplit[2])

	fmt.Println("metodo de requisição:", req.Method)
	fmt.Println("caminho de requisição:", req.Path)
	fmt.Println("versão do protocolo:", req.ProtoVersion)

	if req.Path == "/" {
		req.Path = "/index.html"
	}

	for {
		line, err := req.BodyReader.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler dados do cliente:", err)
			return
		}
		if line == "\r\n" || line == "\n" {
			break
		}

		fmt.Printf("Header: %s", line)

		splitLine := strings.SplitN(line, ":", 2)
		if len(splitLine) == 2 {
			headerName := strings.ToLower(strings.TrimSpace(splitLine[0]))
			headerValue := strings.TrimSpace(splitLine[1])

			req.Headers[headerName] = headerValue
		}

	}

	req.ContentLength = 0
	if val, ok := req.Headers["content-length"]; ok {
		req.ContentLength, err = strconv.Atoi(val)
		if err != nil {
			fmt.Println("Erro ao converter Content-Length:", err)
			return
		}
		if req.ContentLength < 0 {
			fmt.Println("Content-Length inválido:", req.ContentLength)
			return
		}
	}

	fmt.Printf("[REQUEST] Método: %s | Path: %s | Content-Length: %d\n", req.Method, req.Path, req.ContentLength)
	_, err = connection.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, World!"))
	if err != nil {
		fmt.Println("Erro ao enviar dados para o cliente:", err)
	}
}
