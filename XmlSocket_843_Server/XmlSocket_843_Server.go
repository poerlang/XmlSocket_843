package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	Head = 4
)

var (
	ClientMap map[int]net.Conn = make(map[int]net.Conn)
)

func main() {
	fmt.Println(os.Args[0])
	ip_port := ":843"
	if len(os.Args) > 1 {
		ip_port = os.Args[1]
	}
	fmt.Println(
		"\nFlash AS SocketXml 843 crossdomain.xml\nhow to use:\nXmlSocket_843_Server.exe 192.168.101.139:843",
		"\nlisten:", ip_port,
		"\ndo not close this window...")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip_port)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	clientIndex := 0

	for {
		clientIndex++
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, clientIndex)
	}
}

func handleClient(conn net.Conn, index int) {
	ClientMap[index] = conn

	fmt.Println("")
	fmt.Println("=======================")
	fmt.Println("client from: ", conn.RemoteAddr(), "index: ", index)
	fc := func() {
		time.Sleep(time.Second) //给客户端1秒的响应的时间，否则客户端有可能读不到数据就提前Close了
		conn.Close()
		delete(ClientMap, index)
		fmt.Println("remove client: ", index)
		fmt.Println("=======================")
	}
	defer fc()
	sendFirstMsg(conn)
}
func sendFirstMsg(conn net.Conn) {
	str := `<?xml version="1.0"?>
			<!DOCTYPE cross-domain-policy SYSTEM "/xml/dtds/cross-domain-policy.dtd">
			<cross-domain-policy>
				<site-control permitted-cross-domain-policies="master-only"/>
				<allow-access-from domain="*" to-ports="*" />
			</cross-domain-policy>`
	writer := bufio.NewWriter(conn)
	writer.WriteString(str)
	writer.Flush()
	fmt.Println("done: crossdomain.xml")
}
func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
