package main

import (
	"fmt"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
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
		defer conn.Close()
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

		/*
			リクエスト内容の解析
		*/
		reqStr := string(data)
		fmt.Printf("request: %s\n", reqStr)
		// リクエストラインを取得
		reqLine, reqRest, _ := strings.Cut(reqStr, "\r\n")
		fmt.Printf("request line: %s\n", reqLine)
		fmt.Printf("rest: %s\n", reqRest)
		// リクエストヘッダとボディを取得
		reqRestSplit := strings.Split(reqRest, "\r\n\r\n")
		fmt.Printf("request header: %s\n", reqRestSplit[0])
		fmt.Printf("request body: %s\n", reqRestSplit[1])
		// リクエストラインをパースする(-> [GET / HTTP/1.1])
		splitReqLine := strings.Split(reqLine, " ")
		reqPath := splitReqLine[1]
		if reqPath == "/" {
			reqPath = "/index.html"
		}

		/*
			レスポンス生成
		*/
		staticDir, err := filepath.Abs("./server/static")
		fmt.Printf("static dir: %s\n", staticDir)
		fmt.Printf("req path: %s\n", staticDir+reqPath)

		var resLine, resHeader, resBody string

		// 静的資材が置いてあるstaticディレクトリからリクエストパスに対応するファイルを取得
		staticFile, err := os.Open(staticDir + reqPath)
		defer staticFile.Close()
		if err != nil {
			resBody = "<html><body><h1>404 Not Found</h1></body></html>"
			resLine = "HTTP/1.1 404 Not Found\r\n"
		} else {
			staticFileInfo, _ := staticFile.Stat()
			staticFileContent := make([]byte, staticFileInfo.Size())
			_, err = staticFile.Read(staticFileContent)
			resBody = string(staticFileContent)
			resLine = "HTTP/1.1 200 OK\r\n"
		}

		resHeader = ""
		// Dateヘッダー用の時刻生成のためタイムゾーン:GMTの現在時刻を取得
		gmt, _ := time.LoadLocation("GMT")
		resHeader += fmt.Sprintf("Date: %s\r\n", time.Now().In(gmt).Format(time.RFC1123))
		resHeader += "Host: HenaServer/0.1\r\n"
		resHeader += fmt.Sprintf("Content-Length: %d\r\n", len(resBody))
		resHeader += "Connection: Close\r\n"
		// リクエストパスから拡張子を取得し、Content-Typeヘッダーを設定
		ext := filepath.Ext(reqPath)
		contentType := mime.TypeByExtension(ext)
		resHeader += fmt.Sprintf("Content-Type: %s\r\n", contentType)

		/*
			レスポンス送信
		*/
		res := resLine + resHeader + "\r\n" + resBody
		conn.Write([]byte(res))

		continue
	}
}

func errHandler(err error, msg string) {
	if err != nil {
		fmt.Printf("%s failed: %s\n", msg, err)
		// os.Exit(1)
	}
}
