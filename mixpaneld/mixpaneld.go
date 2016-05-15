package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	retries := 10
	var conn net.Conn
	outAddr := os.Getevn("OUTSOCKET_ADDR")
	if outAddr == nil {
		outAddr == "localhost:4343"
	}
	var err error
	for i := 0; i < retries && conn == nil; i++ {
		if err != nil {
			log.Print("Failed to connect to", outAddr, " will retry in 5 seconds.")
			c := time.After(5 * time.Second)
			<-c
		}
		conn, err = net.Dial("tcp", outAddr)
	}
	if err != nil {
		panic(fmt.Sprintf("Could not connect after %d retries, because %s.", retries, err))
	}
	http.HandleFunc("/track/", func(w http.ResponseWriter, r *http.Request) {
		rawData := r.URL.Query().Get("data")
		data, err := base64.StdEncoding.DecodeString(rawData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		m := make(map[string]interface{})
		d := json.NewDecoder(bytes.NewReader(data))
		err = d.Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		buf := bytes.Buffer{}
		e := json.NewEncoder(&buf)
		err = e.Encode(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		conn.Write(buf.Bytes())
		client := http.Client{}
		_, err = client.Get(fmt.Sprintf("http://api.mixpanel.com/track/?data=%s", rawData))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
