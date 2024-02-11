package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	ipAddress := os.Args[1]
	port := os.Args[2]

	listenAddress := ipAddress + ":" + port

	broadcastAddress := "255.255.255.255:" + port

	fmt.Print("Enter your nickname: ")
	reader := bufio.NewReader(os.Stdin)
	nickname, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading nickname:", err)
		os.Exit(1)
	}

	nickname = strings.TrimSpace(nickname)

	go listenUDP(listenAddress, nickname)

	go sendUDP(broadcastAddress, nickname)

	select {}
}

func listenUDP(address, nickname string) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening for UDP packets:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Listening for UDP packets on", address)

	buffer := make([]byte, 1024)

	// Принимаем входящие датаграммы
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}

		fmt.Printf("Received from %s: %s\n", addr.IP.String(), string(buffer[:n]))
	}
}

func sendUDP(address, nickname string) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	// Создаем UDP соединение для отправки
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error connecting to UDP server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// Бесконечный цикл для чтения ввода пользователя и отправки датаграмм
	for {
		fmt.Print("Enter message: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		if len(message) > 1000 {
			fmt.Println("Error: Message exceeds 1000 bytes.")
			continue
		}

		fullMessage := fmt.Sprintf("%s: %s", nickname, message)

		// Отправляем датаграмму
		_, err = conn.Write([]byte(fullMessage))
		if err != nil {
			fmt.Println("Error sending UDP packet:", err)
			continue
		}
	}
}
