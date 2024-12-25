package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	fmt.Println("server start.")

	// 指定されたアドレスをTCP通信で利用できる形式に解決する
	// IPv4を使うのでtcp4を指定する
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:8080")
	errHandler(err, "resolve tcp addr")

	// 指定したアドレスとポート番号でTCP接続を待ち受けるソケットを作成する
	// 低レイヤー的な補足として、
	// - OSのsocket()システムコールを内部で呼び出し、待受用のソケットを生成する
	// - ソケットは「未接続状態」のままリスニングモードで待機する
	listener, err := net.ListenTCP("tcp", tcpAddr)
	errHandler(err, "listen tcp")

	for {
		// リッスンソケットが接続要求を受けて、新しいソケット(クライアントごと)を生成する
		// accept()システムコールを内部で呼び出し、3-way handshake(TCP接続確立)を完了させる
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("listener accept failed: %s", err)
			continue
		}

		// リクエスト内容を読み取る
		data := make([]byte, 1024)
		rLen, err := conn.Read(data)
		errHandler(err, "read from socket")
		fmt.Printf("request len: %d", rLen)

		// リクエスト内容を書き込むファイルを書き込み権限付きで開く
		f, err := os.Create("./server/request_body.data")
		errHandler(err, "open file")

		// リクエスト内容をファイルに書き込む
		_, err = f.Write(data)
		errHandler(err, "write to file")

		// タイムゾーンをGMTに指定
		gmt, _ := time.LoadLocation("GMT")

		// レスポンス生成
		resBody := "<html><body><h1>It works!</h1></body></html>"
		resLine := "HTTP/1.1 200 OK\r\n"
		resHeader := ""
		resHeader += fmt.Sprintf("Date: %s\r\n", time.Now().In(gmt).Format(time.RFC1123))
		resHeader += "Host: HenaServer/0.1\r\n"
		resHeader += fmt.Sprintf("Content-Length: %d\r\n", len(resBody))
		resHeader += "Connection: Close\r\n"
		resHeader += "Content-Type: text/html\r\n"
		res := resLine + resHeader + "\r\n" + resBody

		conn.Write([]byte(res))

		err = conn.Close()
		errHandler(err, "close")
	}
}

func errHandler(err error, msg string) {
	if err != nil {
		fmt.Printf("%s failed: %s\n", msg, err)
		os.Exit(1)
	}
}
