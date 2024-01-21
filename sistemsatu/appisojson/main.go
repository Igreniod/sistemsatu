package main

import (
	"fmt"
	"net"
)

func main() {
	//-----Listen di port 8082-----//
	listener, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8082")

	for {
		// Menerima koneksi yang tersambung
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	//-----Membuat buffer untuk membaca dan menyimpan data di buffer-----//
	buffer := make([]byte, 1024)

	for {
		//-----Membaca data dari pengirim pesan-----//
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		//-----Memproses dan menggunakan data-----//
		receivedData := buffer[:n]
		// fmt.Printf("Received: %s\n", receivedData)

		datarsp := IsoHandler(string(receivedData))

		//-----Mengirim respond ke pengirim pesan-----//
		response := []byte(string(datarsp))
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error responding to client:", err)
			return
		}
	}
}
