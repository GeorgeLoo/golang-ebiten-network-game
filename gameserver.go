/* FTP Server
 */
package main

import (
	"fmt"
	"net"
	"os"
)

const (
	DIR       = "DIR"
	CD        = "CD"
	PWD       = "PWD"
	kShutdown = "SHUTDOWN" // gloo 5.3.2018
	kSend     = "SEND"
)

var (
	gSendBuff string
)

func main() {

	fmt.Println("Server running ")

	service := "0.0.0.0:1202"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {

	fmt.Println("handleClient ")

	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			conn.Close()
			return
		}

		s := string(buf[0:n])
		fmt.Println("s ", s)
		// decode request
		if s[0:2] == CD {
			chdir(conn, s[3:])
		} else if s[0:3] == DIR {
			dirList(conn)
		} else if s[0:3] == PWD {
			//pwd(conn)
			sendmessage(conn) // gloo 25.3.2018
		} else if s[0:8] == kShutdown {
			fmt.Println("shut down ")
			shutdown(conn)
			os.Exit(0)
		} else if s[0:4] == kSend {
			fmt.Println("send ")
			gSendBuff = s[4:]
			fmt.Println("gSendBuff ", gSendBuff)
		}

	}
}

func sendmessage(conn net.Conn) {
	s := gSendBuff
	conn.Write([]byte(s))
}

func shutdown(conn net.Conn) {
	s := "OK"
	conn.Write([]byte(s))
}

func chdir(conn net.Conn, s string) {
	if os.Chdir(s) == nil {
		conn.Write([]byte("OK"))
	} else {
		conn.Write([]byte("ERROR"))
	}
}

func pwd(conn net.Conn) {
	s, err := os.Getwd()
	if err != nil {
		conn.Write([]byte(""))
		return
	}
	conn.Write([]byte(s))
}

func dirList(conn net.Conn) {
	defer conn.Write([]byte("\r\n"))

	dir, err := os.Open(".")
	if err != nil {
		return
	}

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, nm := range names {
		conn.Write([]byte(nm + "\r\n"))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
