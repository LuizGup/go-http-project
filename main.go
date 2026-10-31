package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":7676")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor: ", err)
		return
	}
	defer listener.Close()
	
	fmt.Println("[SERVIDOR] Escutando na porta :7676...")
	for {
		conection, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão: ", err)
			continue
		}
		fmt.Println("vou lidar com o cliente...")
		go HandleClient(conection)

	}
}