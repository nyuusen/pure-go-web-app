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

	// ファイルからリクエスト内容を読み込む
	// ファイルパスはgoコマンドを実行するディレクトリからの相対パスになる
	f, err := os.Open("./client/client-send.data")
	errHandler(err, "open client-send.data")

	data := make([]byte, 1000)
	cnt, err := f.Read(data)
	fmt.Printf("read %d byte", cnt)

	// ソケットにリクエストデータを書き込む
	_, err = conn.Write(data)
	errHandler(err, "write to socket")

	// ソケットからレスポンスデータを読み込む
	res := make([]byte, 1024)
	len, err := conn.Read(res)
	errHandler(err, "read from socket")

	// レスポンスデータをファイルに書き込む
	f2, err := os.Create("./client/client-received.data")
	errHandler(err, "open client-received.data")

	_, err = f2.Write(res)

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
