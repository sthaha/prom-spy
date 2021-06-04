package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type mw struct {
	next http.Handler
}

type responseSpy struct {
	wrap http.ResponseWriter
	body string
}

var _ http.ResponseWriter = (*responseSpy)(nil)

func (spy *responseSpy) Write(b []byte) (int, error) {
	spy.body = spy.body + string(b)
	return spy.wrap.Write(b)
}

func (spy *responseSpy) Header() http.Header {
	return spy.wrap.Header()
}

func (spy *responseSpy) WriteHeader(s int) {
	spy.wrap.WriteHeader(s)
}

func (m mw) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	sb := strings.Builder{}

	for k, v := range r.Header {
		sb.WriteString("\t")
		sb.WriteString(k)
		sb.WriteString(" : ")
		sb.WriteString(strings.Join(v, ", "))
		sb.WriteString("\n")
	}

	ts := time.Now()

	r.Header["Accept-Encoding"] = []string{"deflate"}
	spy := &responseSpy{wrap: w}

	m.next.ServeHTTP(spy, r)

	fmt.Fprintf(os.Stderr, `
Time: %v
%s %s
Host: %s
Headers:
%s
		------------------------------
%s
_____________________________________________
	`, ts, r.Method, r.URL, r.Host, sb.String(), spy.body)
}

var _ http.Handler = (*mw)(nil)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", &mw{promhttp.Handler()})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
