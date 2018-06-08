package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/fatih/color"
	server "github.com/tendermint/abci/server"
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	tdmLog "github.com/tendermint/tmlibs/log"
	"github.com/watcharaphat/tendermint-benchmark/abci/did"
)

func main() {
	IPRegister()
	runABCIServer(os.Args)
}

// Add ip to ip list
func IPRegister() {
	ip := GetOutboundIP().String()
	url := `http://` + os.Getenv("DISCOVERY_HOSTNAME") + `:` + os.Getenv("DISCOVERY_PORT") + `/abci-ip/add/` + string(ip)
	fmt.Println(`curl ` + url)

	resp := getRequest(url)

	defer resp.Body.Close()
}

func getRequest(url string) *http.Response {
	resp, err := http.Get(os.ExpandEnv(url))
	if err != nil {
		fmt.Println("Trying to do IP registration again.")
		getRequest(url)
	}

	return resp
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func runABCIServer(args []string) {
	address := "tcp://127.0.0.1:46658"
	if len(args) > 1 {
		address = args[1]
		fmt.Println(address)
	}

	logger := tdmLog.NewTMLogger(tdmLog.NewSyncWriter(os.Stdout))

	var app types.Application
	app = did.NewDIDApplication()

	// Start the listener
	srv, err := server.NewServer(address, "socket", app)
	if err != nil {
		color.Red("%s", err)
	}
	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		color.Red("%s", err)
	}

	// Create a web server on port 8100

	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

	// Wait forever
	cmn.TrapSignal(func() {
		srv.Stop()
	})
}
