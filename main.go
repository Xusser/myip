package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
)

var path string
var port int64
var address string
var useForwardedIP bool

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
	log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	log.SetLevel(log.DebugLevel)
}

func handler(w http.ResponseWriter, req *http.Request) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	forwardedIP := req.Header.Get("X-Forwarded-For")
	if forwardedIP != "" && net.ParseIP(forwardedIP) != nil && useForwardedIP {
		ip = forwardedIP
	}

	fmt.Fprintf(w, "%s", ip)
}

func main() {
	flag.StringVar(&path, "s", "", "Sub path(default:\"\")")
	flag.Int64Var(&port, "p", 8000, "Listen port(default:8000)")
	flag.StringVar(&address, "b", "127.0.0.1", "Bind address(default:\"0.0.0.0\")")
	flag.BoolVar(&useForwardedIP, "f", false, "Use X-Forwarded-For given if exist and valid(default:false)")
	flag.Parse()

	if net.ParseIP(address) == nil {
		log.Fatal("Invalid address:", address)
	}
	if port < 1 || port > 65535 {
		log.Fatal("Invalid port:", port)
	}

	addr := fmt.Sprintf("%s:%d", address, port)
	pattern := fmt.Sprintf("/%s", path)
	http.HandleFunc(pattern, handler)

	log.Info("Serve at ", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
	log.Info("Stopped")
}
