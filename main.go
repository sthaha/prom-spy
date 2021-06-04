package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type mux struct{}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	sb := strings.Builder{}

	for k, v := range r.Header {
		sb.WriteString("\t")
		sb.WriteString(k)
		sb.WriteString(" : ")
		sb.WriteString(strings.Join(v, ", "))
		sb.WriteString("\n")
	}

	fmt.Fprintf(os.Stderr, `
Time: %v
%s %s
Host: %s
Headers:
%s
_____________________________________________
	`, time.Now(), r.Method, r.URL, r.Host, sb.String())
	w.WriteHeader(200)
	w.Write([]byte("foobar"))
}

var _ http.Handler = (*mux)(nil)

func main() {
	log.Fatal(http.ListenAndServe(":8080", &mux{}))
}
