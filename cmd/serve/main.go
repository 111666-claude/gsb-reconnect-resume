// Command serve 用 httptest 跑固定场景，打印每个请求的状态码与响应。
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"example.com/reconnect-resume/session"
)

func post(srv *httptest.Server, path, body string) (int, string) {
	res, err := http.Post(srv.URL+path, "application/json", bytes.NewBufferString(body))
	if err != nil {
		return 0, err.Error()
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(bytes.TrimSpace(raw))
}

func main() {
	sample := flag.String("sample", "grace", "场景名")
	flag.Parse()

	m := &session.Manager{GraceMs: 1000}
	srv := httptest.NewServer(m.Handler())
	defer srv.Close()

	switch *sample {
	case "grace":
		code, body := post(srv, "/session/open", `{"player":"p1","atMs":0}`)
		fmt.Println(code, body)
		code, body = post(srv, "/session/resume", `{"requestId":"r1","token":"tok-p1-0","version":1,"atMs":5000}`)
		fmt.Println(code, body)
	case "version":
		post(srv, "/session/open", `{"player":"p1","atMs":0}`)
		code, body := post(srv, "/session/resume", `{"requestId":"r1","token":"tok-p1-0","version":99,"atMs":10}`)
		fmt.Println(code, body)
	case "idempotent":
		post(srv, "/session/open", `{"player":"p1","atMs":0}`)
		post(srv, "/session/resume", `{"requestId":"r1","token":"tok-p1-0","version":1,"atMs":10}`)
		code, body := post(srv, "/session/resume", `{"requestId":"r1","token":"tok-p1-0","version":1,"atMs":20}`)
		fmt.Println(code, body)
	case "work":
		for i := 0; i < 5000; i++ {
			post(srv, "/session/open", fmt.Sprintf(`{"player":"p%d","atMs":%d}`, i, i))
		}
		for i := 0; i < 5000; i++ {
			post(srv, "/session/resume",
				fmt.Sprintf(`{"requestId":"r%d","token":"tok-p%d-%d","version":1,"atMs":%d}`, i, i, i, i+1))
		}
		out, _ := json.Marshal(map[string]int{"resumes": m.Resumes, "scanned": m.Scanned})
		fmt.Println(string(out))
	default:
		fmt.Fprintf(os.Stderr, "未知场景 %q\n", *sample)
		os.Exit(2)
	}
}
