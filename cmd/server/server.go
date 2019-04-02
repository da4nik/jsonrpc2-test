package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"

	"github.com/da4nik/jrpc2_try/internal/http"
	"github.com/da4nik/jrpc2_try/internal/log"
)

const defaultPort = 8080

var port int

func loadENV() {
	envPort, err := strconv.Atoi(os.Getenv("JRT_PORT"))
	if err != nil {
		port = defaultPort
		return
	}
	port = envPort
}

func parseFlags() {
	portPtr := flag.Int("port", port, "api port (env JRT_PORT)")

	flag.Parse()

	port = *portPtr
}

func main() {
	log.InitLogger()

	loadENV()
	parseFlags()

	httpServer, err := http.NewHTTPServer(port)
	if err != nil {
		log.Errorf("Error creating http server: %s", err.Error())
		os.Exit(1)
	}

	httpServer.Start()
	defer httpServer.Stop()

	log.Infof("Ctrl-C to interrupt")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
