package main

import (
	"bufio"
	//"encoding/binary"
	"fmt"
	"net"
	"os"
	//"strings"
	"time"
	//"strconv"
)

const (
	Head = 4
)

var (
	ClientMap map[int]net.Conn = make(map[int]net.Conn)
)

func main() {
	fmt.Println(os.Args[0])
	ip_port := "127.0.0.1:843"
	if len(os.Args) > 1 {
		ip_port = os.Args[1]
	}
	fmt.Println(
		"\nFlash AS 策略服务运行中...\n自动回应SocketXml端口(即843端口)的crossdomain.xml请求\n如需指定ip和端口，可以在程序启动时指定参数，格式如下\ngameserver843.exe 192.168.101.139:843",
		"\n当前正在侦听", ip_port,
		"\n请不要关闭此窗口...")
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
	fmt.Println("新用户连接, 来自: ", conn.RemoteAddr(), "index: ", index)
	fc := func() {
		time.Sleep(time.Second) //给客户端1秒的响应的时间，否则客户端有可能读不到数据就提前Close了
		conn.Close()
		delete(ClientMap, index)
		fmt.Println("移除序号为: ", index, "的客户端，断开客户端的连接")
		fmt.Println("=======================")
	}
	defer fc()
	sendFirstMsg(conn)

	//不必执行下列解析，收到请求可直接 sendFirstMsg(conn)
	/*
		reader := bufio.NewReader(conn)
		allLen := 23 //"<policy-file-request/>\0"长度正好23，剩下的就不必去读了，此函数退出时，会自动执行defer fc()关闭连接
		readedLen := 0
		bodySl := make([]byte, allLen)
		for readedLen < allLen {
			len, err := reader.Read(bodySl)
			if err != nil {
				fmt.Println("读取包体出错,: ", err.Error())
				break
			}
			readedLen += len
		}
		if strings.Contains(string(bodySl), "policy") {
			fmt.Println("收到策略文件请求：", string(bodySl))
			sendFirstMsg(conn)
		} else {
			fmt.Println("格式不正确:", string(bodySl))
		}
	*/
}
func sendFirstMsg(conn net.Conn) {
	str := `<?xml version="1.0" encoding="UTF-8"?>
            <!DOCTYPE cross-domain-policy SYSTEM "http://www.macromedia.com/xml/dtds/cross-domain-policy.dtd">
            <cross-domain-policy>
            <site-control permitted-cross-domain-policies="all" />
            <allow-access-from domain="*" to-ports="*" />
            </cross-domain-policy>\0`
	writer := bufio.NewWriter(conn)
	writer.WriteString(str)
	writer.Flush()
	fmt.Println("已经回应策略文件：crossdomain.xml")
}
func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
