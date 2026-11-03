package main

import (
	"bufio"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type HttpRequest struct {
	Method        string
	Path          string
	ProtoVersion  string
	Headers       map[string]string
	ContentLength int
	BodyReader    *bufio.Reader
}

var ErrAccessDenied = errors.New("acesso negado")

func HandleClient(connection net.Conn) {
	defer connection.Close()
	defer fmt.Println("[CONEXÃO] Encerrada", connection.RemoteAddr())

	fmt.Println("[CONEXÃO] Nova conexão de", connection.RemoteAddr())

	reader := bufio.NewReader(connection)

	req, err := parseRequest(reader)
	if err != nil {
		fmt.Println("Erro ao parsear:", err)
		basicErrorResponse(connection, "400 Bad Request", "")
		return
	}

	switch req.Method {
	case "GET":
		err := serveFile(connection, req)
		if err != nil {
			switch {
			case errors.Is(err, os.ErrNotExist):
				fmt.Println("[404] Recurso não encontrado:", req.Path)
				pagina_404, err := os.ReadFile("./public/404.html")
				if err != nil {
					fmt.Println(err)
					basicErrorResponse(connection, "500 Internal Server Error", "")
					return
				}
				basicErrorResponse(connection, "404 Not Found", string(pagina_404))

			case errors.Is(err, ErrAccessDenied) || errors.Is(err, os.ErrPermission):
				fmt.Println("Acesso negado:", req.Path)
				basicErrorResponse(connection, "403 Forbidden", "")

			default:
				fmt.Println("Erro ao servir arquivo:", err)
				basicErrorResponse(connection, "500 Internal Server Error", "")
			}
			return
		}

	case "POST":
		// Implementação do POST
		body := make([]byte, req.ContentLength)
		n, err := io.ReadFull(req.BodyReader, body)
		if err != nil {
			fmt.Println("Erro ao ler corpo da requisição:", err)
			basicErrorResponse(connection, "400 Bad Request", "")
			return
		}
		fmt.Printf("[POST] Body bruto (%d bytes): %s\n", n, string(body))
		var pagina strings.Builder
		pagina.WriteString(`<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>Dados Recebidos</title>
</head>
<body>
    <h1>Dados Recebidos</h1>
    <ul>
`)
		pares := strings.Split(string(body), "&")
		for par := range pares {
			if pares[par] == "" {
				continue
			}
			nome, valor, _ := strings.Cut(pares[par], "=")
			valor_decodificado, err := url.QueryUnescape(valor)
			if err != nil {
				fmt.Println("Erro ao decodificar valor:", err)
				basicErrorResponse(connection, "400 Bad Request", "")
				return
			}
			nome_decodificado, err := url.QueryUnescape(nome)
			if err != nil {
				fmt.Println("Erro ao decodificar nome:", err)
				basicErrorResponse(connection, "400 Bad Request", "")
				return
			}

			fmt.Println("[POST] Campo:", nome_decodificado, "=", valor_decodificado)
			fmt.Fprintf(&pagina, "        <li><strong>%s:</strong> %s</li>\n",
				html.EscapeString(nome_decodificado), html.EscapeString(valor_decodificado))
		}
		pagina.WriteString(`    </ul>
    <a href="/">Voltar para a página inicial</a>
</body>
</html>
`)
		resposta := pagina.String()
		if err := basicErrorResponse(connection, "200 OK", resposta); err != nil {
			return
		}
		fmt.Printf("[RESPONSE] Resposta dinâmica enviada (text/html, %d bytes)\n", len(resposta))
	default:
		fmt.Println("Método não suportado:", req.Method)
		basicErrorResponse(connection, "405 Method Not Allowed", "")
	}
}

func parseRequest(reader *bufio.Reader) (HttpRequest, error) {
	req := HttpRequest{
		ContentLength: 0,
		Headers:       make(map[string]string),
		BodyReader:    reader,
	}

	line, err := req.BodyReader.ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler dados do cliente:", err)
		return req, err
	}

	fmt.Printf("Linha de Requisição: %s", line)

	lineSplit := strings.Split(line, " ")
	if len(lineSplit) != 3 {
		fmt.Println("Linha de requisição inválida:", line)
		return req, fmt.Errorf("linha de requisição inválida")
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
			return req, err
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
			return req, err
		}
		if req.ContentLength < 0 {
			fmt.Println("Content-Length inválido:", req.ContentLength)
			return req, fmt.Errorf("Content-Length inválido")
		}
	}
	fmt.Printf("[REQUEST] Método: %s | Path: %s | Content-Length: %d\n", req.Method, req.Path, req.ContentLength)
	return req, nil
}

func serveFile(connection net.Conn, req HttpRequest) error {
	lowerPath := strings.ToLower(req.Path)
	contentType := ""
	switch {
	case strings.HasSuffix(lowerPath, ".html"):
		contentType = "text/html"
	case strings.HasSuffix(lowerPath, ".css"):
		contentType = "text/css"
	case strings.HasSuffix(lowerPath, ".png"):
		contentType = "image/png"
	case strings.HasSuffix(lowerPath, ".jpg") || strings.HasSuffix(lowerPath, ".jpeg"):
		contentType = "image/jpeg"
	default:
		contentType = "application/octet-stream"
	}

	publicAbs, err := filepath.Abs("./public")
	if err != nil {
		fmt.Println("Erro ao obter caminho absoluto do diretório public:", err)
		return err
	}

	cleanPath := filepath.Join("./public", req.Path)
	fileAbs, err := filepath.Abs(cleanPath)
	if err != nil {
		fmt.Println("Erro ao obter caminho absoluto do arquivo:", err)
		return err
	}
	publicAbs = publicAbs + string(filepath.Separator)

	boolFilePath := strings.HasPrefix(fileAbs, publicAbs)
	if !boolFilePath {
		fmt.Println("Tentativa de acesso a arquivo fora do diretório public:", fileAbs)
		return ErrAccessDenied
	}

	println("booleano de verificação de caminho - Debbug: " + strconv.FormatBool(boolFilePath))
	println("as 3 Variaveis de caminho - Debbug: " + "\npublicAbs: " + publicAbs + "\nfileAbs: " + fileAbs + "\ncleanPath: " + cleanPath)
	fileBytes, err := os.ReadFile(fileAbs)
	if err != nil {
		fmt.Println("Erro ao ler arquivo:", err)
		return err
	}

	fileLength := strconv.Itoa(len(fileBytes))

	message := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %s\r\n\r\n", contentType, fileLength)
	_, err = connection.Write([]byte(message))
	if err != nil {
		fmt.Println("Erro ao enviar dados para o cliente:", err)
		return err
	}

	_, err = connection.Write(fileBytes)
	if err != nil {
		fmt.Println("Erro ao enviar dados para o cliente:", err)
		return err
	} else {
		fmt.Printf("[ARQUIVO] Servindo %s (%s, %s bytes)\n", req.Path, contentType, fileLength)
	}
	return nil
}

func basicErrorResponse(connection net.Conn, status string, body string) error {
	message := fmt.Sprintf("HTTP/1.1 %s\r\nContent-Type: text/html\r\nContent-Length: %d\r\n\r\n%s", status, len(body), body)
	n, err := connection.Write([]byte(message))
	if err == nil && n != len(message) {
		err = io.ErrShortWrite
	}
	if err != nil {
		fmt.Println("Erro ao enviar resposta de erro para o cliente:", err)
	}
	return err
}
