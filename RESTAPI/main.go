package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"encoding/json"
	"log"

	"github.com/gorilla/mux"
)

type payload struct {
	Payload string `json:"message"`
}

func SendAll(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p payload
	_ = decoder.Decode(&p)
	defer r.Body.Close()

	params := mux.Vars(r)
	seq := params["seq"]

	//url := `http://` + os.Getenv("TENDERMINT_IP") + `:` + os.Getenv("TENDERMINT_PORT") + `/broadcast_tx_async?tx="` + seq + `=` + p.Payload + `"`
    url := `http://localhost:46657/broadcast_tx_async?tx="` + seq + `=` + p.Payload + `"`

	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
}

func SendIdp(w http.ResponseWriter, r *http.Request) {
	/*
		var t test_struct
		defer r.Body.Close()
		payload, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(payload, &t)
		println("OK" + t.Test)
	*/ ///////////////
	decoder := json.NewDecoder(r.Body)
	var p payload
	_ = decoder.Decode(&p)
	defer r.Body.Close()
	//log.Println(p.Payload)

	params := mux.Vars(r)
	seq := params["seq"]

	url := `http://localhost:46657/broadcast_tx_async?tx="` + seq + `=` + p.Payload + `"`
	// url := `http://` + os.Getenv("TENDERMINT_IP") + `:` + os.Getenv("TENDERMINT_PORT") + `/broadcast_tx_async?tx="` + seq + `=` + p.Payload + `"`

	fmt.Println(string(url))
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
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

// Add ip to ip list
func IPRegister() {
	ip := GetOutboundIP().String()
	url := `http://` + os.Getenv("DISCOVERY_HOSTNAME") + `:` + os.Getenv("DISCOVERY_PORT") + `/ip/add/` + string(ip)
	fmt.Println(`curl ` + url)

	resp, err := http.Get(os.ExpandEnv(url))
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
}

func main() {
	IPRegister()
	router := mux.NewRouter()
	router.HandleFunc("/send_all/{seq}", SendAll).Methods("POST")
	router.HandleFunc("/send_idp/{seq}", SendIdp).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
