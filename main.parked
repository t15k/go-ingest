package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
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
		//buf := make([]byte, 0, 1024)
		buf := bytes.Buffer{}
		e := json.NewEncoder(&buf)
		err = e.Encode(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Print(string(buf.Bytes()))
		// TODO, we shouldn't just print it, we should ship it somewhere.
		// send to mix panel
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
