package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("client start.")

	port := flag.Int("port", 80, "destination server port(int)")
	flag.Parse()

	fmtPort := fmt.Sprintf(":%d", *port)
	fmt.Printf("destination server port is %s\n", fmtPort)

	// 指定されたアドレスをTCP通信で利用できる形式に解決する
	// TCP通信にはソケットが必要で、そのソケットを作成するには接続先のIPとポートが必要となる
	// IPv4を使うのでtcp4を指定する
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmtPort)
	errHandler(err, "resolve tcp addr")

	// サーバーにTCP接続を確立する
	// ソケット作成し、3ウェイハンドシェイクを経て接続を確立
	// 第2引数(local address)はnilを指定(nil指定で自動割り当て)
	// 第3引数(remote address)はサーバー側のアドレスを指定
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	errHandler(err, "dial")

	// ファイルからリクエスト内容を読み込む
	// ファイルパスはgoコマンドを実行するディレクトリからの相対パスになる
	// NOTE: この実装では、Apacheにリクエストを送ることができないため、コメントアウトしている
	// 発生するエラー：「Your browser sent a request that this server could not understand.」
	// f, err := os.Open("./client/client_send.data")
	// errHandler(err, "open client_send.data")
	// data := make([]byte, 1000)
	// cnt, err := f.Read(data)
	// fmt.Printf("read %d byte\n", cnt)

	req := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	reqBytes := []byte(req)

	// ソケットにリクエストデータを書き込む
	_, err = conn.Write(reqBytes)
	errHandler(err, "write to socket")

	// ソケットからレスポンスデータを読み込む
	res := make([]byte, 1024)
	len, err := conn.Read(res)
	errHandler(err, "read from socket")

	// レスポンスデータをファイルに書き込む
	f2, err := os.Create("./client/client_recv.data")
	errHandler(err, "open client_recv.data")

	_, err = f2.Write(res)

	fmt.Printf("response: %s\n", string(res[:len]))

	// 接続を切断
	conn.Close()
}

func errHandler(err error, msg string) {
	if err != nil {
		fmt.Printf("%s failed: %s\n", msg, err)
		os.Exit(1)
	}
}
