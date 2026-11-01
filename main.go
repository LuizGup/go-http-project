package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	flagPort := flag.String("port", "9000", "Porta TCP para escutar")
	flag.Parse()
	tcpPort := ":" + *flagPort

	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("[SERVIDOR] Escutando na porta", tcpPort)
	for {
		conection, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go HandleClient(conection)

	}
}
