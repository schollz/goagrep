package main

import (
	"net"
	"os"
)

func main() {
	usage := "Echos the message to the TCP host specified.\nusage: tcpecho [message] [destination] [port]"
	if len(os.Args) != 4 {
		println(usage)
		os.Exit(1)
	}
	strEcho := os.Args[1] + "\n"
	servAddr := os.Args[2] + ":" + os.Args[3]
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println(usage)
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	// start := time.Now()
	// for i := 0; i < 1; i++ {

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(strEcho))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	println(string(reply))

	conn.Close()

	// }
	// elapsed := time.Since(start)
	// log.Printf("Binomial took %s", elapsed)
}
