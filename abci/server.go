package main

import (
	"fmt"
	"os"
    "flag"
    "net/http"
    "log"
	"github.com/fatih/color"
	"github.com/oatsaysai/tendermint-benchmark/abci/did"
	server "github.com/tendermint/abci/server"
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	tdmLog "github.com/tendermint/tmlibs/log"
)

func main() {
runABCIServer(os.Args)



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
