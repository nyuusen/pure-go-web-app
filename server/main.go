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
	fmt.Println("===サーバーを起動します===")

	// サーバーを起動する
	listener, err := serve("localhost", 8080)
	errHandler(err, "server")

	for {
		fmt.Println("===クライアントからの接続を待機===")
		// リッスンソケットが接続要求を受けて、新しいソケット(クライアントごと)を生成する
		// accept()システムコールを内部で呼び出し、3-way handshake(TCP接続確立)を完了させる
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Printf("listener accept failed: %s", err)
			continue
		}

		go handleConnection(conn)
	}
}

func serve(host string, port int) (*net.TCPListener, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	// 指定されたアドレスをTCP通信で利用できる形式に解決する
	// IPv4を使うのでtcp4を指定する
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}

	// 指定したアドレスとポート番号でTCP接続を待ち受けるソケットを作成する
	// 低レイヤー的な補足として、
	// - OSのsocket()システムコールを内部で呼び出し、待受用のソケットを生成する
	// - ソケットは「未接続状態」のままリスニングモードで待機する
	return net.ListenTCP("tcp", tcpAddr)
}

func handleConnection(conn net.Conn) {
	fmt.Println("===クライアントとの接続を開始===")
	// リクエスト内容を読み取る
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	errHandler(err, "read from socket")

	method, path, header, body, err := parseHttpRequest(string(data))

	/*
	   レスポンス生成
	*/
	staticDir, err := filepath.Abs("./server/static")

	var resLine, resHeader, resBody string

	if path == "/now" {
		resBody = fmt.Sprintf("<html><body><h1>Now: %s</h1></body></html>", time.Now().Format(time.RFC3339))
		resLine = "HTTP/1.1 200 OK\r\n"
	} else if path == "/show_request" {
		resBody = fmt.Sprintf(`
		<html>
			<body>
				<h1>Request Line</h1>
				<p>%s %s</p>
				<h1>Request Header</h1>
				<p>%s</p>
				<h1>Request Body</h1>
				<p>%s</p>
			</body>
		</html>`, method, path, header, body)
		resLine = "HTTP/1.1 200 OK\r\n"
	} else {
		// 静的資材が置いてあるstaticディレクトリからリクエストパスに対応するファイルを取得
		staticFile, err := os.Open(staticDir + path)
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
	}

	resHeader = ""
	// Dateヘッダー用の時刻生成のためタイムゾーン:GMTの現在時刻を取得
	gmt, _ := time.LoadLocation("GMT")
	resHeader += fmt.Sprintf("Date: %s\r\n", time.Now().In(gmt).Format(time.RFC1123))
	resHeader += "Host: HenaServer/0.1\r\n"
	resHeader += fmt.Sprintf("Content-Length: %d\r\n", len(resBody))
	resHeader += "Connection: Close\r\n"
	// リクエストパスから拡張子を取得し、Content-Typeヘッダーを設定
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	resHeader += fmt.Sprintf("Content-Type: %s\r\n", contentType)

	/*
	   レスポンス送信
	*/
	res := resLine + resHeader + "\r\n" + resBody
	conn.Write([]byte(res))
	fmt.Println("===クライアントとの接続を完了===")
}

func parseHttpRequest(reqStr string) (method, path, header, body string, err error) {
	// リクエストラインを取得
	reqLine, reqRest, _ := strings.Cut(reqStr, "\r\n")
	// リクエストラインをパースする(-> [GET / HTTP/1.1])
	splitReqLine := strings.Split(reqLine, " ")
	fmt.Println(splitReqLine)
	reqPath := splitReqLine[1]
	if reqPath == "/" {
		reqPath = "/index.html"
	}
	// リクエストヘッダとボディを取得
	reqRestSplit := strings.Split(reqRest, "\r\n\r\n")
	// TODO: HTTPメソッドを動的に返却する
	return "GET", reqPath, reqRestSplit[0], reqRestSplit[1], nil
}

func errHandler(err error, msg string) {
	if err != nil {
		fmt.Printf("%s failed: %s\n", msg, err)
	}
}
