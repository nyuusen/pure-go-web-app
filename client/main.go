package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("client start.")

	localSrvPort := ":8080"
	// 指定されたアドレスをTCP通信で利用できる形式に解決する
	// TCP通信にはソケットが必要で、そのソケットを作成するには接続先のIPとポートが必要となる
	// IPv4を使うのでtcp4を指定する
	tcpAddr, err := net.ResolveTCPAddr("tcp4", localSrvPort)
	errHandler(err, "resolve tcp addr")

	// サーバーにTCP接続を確立する
	// ソケット作成し、3ウェイハンドシェイクを経て接続を確立
	// 第2引数(local address)はnilを指定(nil指定で自動割り当て)
	// 第3引数(remote address)はサーバー側のアドレスを指定
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	errHandler(err, "dial")

	// ソケットにデータを書き込む
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	errHandler(err, "write to socket")

	// ソケットからデータを読み込む
	res := make([]byte, 1024)
	len, err := conn.Read(res)
	errHandler(err, "read from socket")
	fmt.Println("response: ", string(res[:len]))

	// 接続を切断
	conn.Close()
}

func errHandler(err error, msg string) {
	if err != nil {
		fmt.Printf("%s failed: %s\n", msg, err)
		os.Exit(1)
	}
}
