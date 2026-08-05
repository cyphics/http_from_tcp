package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil { log.Fatalf("Error resolving UDP address: %s\n", err.Error()) }
	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil { log.Fatalf("Error establishing connection: %s\n", err.Error()) }
	defer connection.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil { fmt.Printf("Error reading from reader: %s\n", err.Error()) }
		b, err := connection.Write([]byte(line))
		if err != nil { fmt.Printf("Error writing line to connection: %s\n", err.Error()) }
		fmt.Println(b)
	}
}
