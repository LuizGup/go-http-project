package main

import (
	"fmt"
	"net"
)

func HandleClient(connection net.Conn) {
	defer connection.Close()
	defer fmt.Println("[CONEXÃO] Encerrada ", connection.RemoteAddr())

	fmt.Println("[CONEXÃO] Nova conexão de ", connection.RemoteAddr())
	buffer := make([]byte, 1024)

	n, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao ler dados do cliente: ", err)
		return
	}

	fmt.Printf("[DADOS RECEBIDOS] %d bytes: %s\n", n, string(buffer[:n]))

	_, err = connection.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, World!"))
	if err != nil {
		fmt.Println("Erro ao enviar dados para o cliente: ", err)
	}

}
