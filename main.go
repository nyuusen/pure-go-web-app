package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("server start.")

	// socketを作成し、localhostの8080番に割り当てる
	// IPv4を使うのでtcp4を指定する
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:8080")
	errHandler(err)

	// 外部からの接続を待機
	// ListenTCPはTCPListener構造体を返却する
	listener, err := net.ListenTCP("tcp", tcpAddr)
	errHandler(err)

	// 接続があったらコネクションを確立する
	for {
		// 接続を受信
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("listener accept failed: %s", err)
			continue
		}

		// リクエスト内容を読み取る
		req := make([]byte, 1024)
		rLen, err := conn.Read(req)
		errHandler(err)
		fmt.Printf("request len: %d", rLen)

		// リクエスト内容を書き込むファイルを書き込み権限付きで開く
		f, err := os.Create("./request-body.txt")
		errHandler(err)

		// リクエスト内容をファイルに書き込む
		_, err = f.Write(req)
		errHandler(err)

		conn.Write([]byte("success"))

		err = conn.Close()
		errHandler(err)
	}

	fmt.Println("server stop.")

}

func errHandler(err error) {
	if err != nil {
		fmt.Printf("error occurred: %s", err)
		os.Exit(1)
	}
}
