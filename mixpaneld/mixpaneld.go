package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/t15k/go-ingest"
	_ "github.com/t15k/go-ingest/socketout"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("missing configuration")
		os.Exit(1)
	}
	receivers, err := ingest.Bootstrap(os.Args[1])
	if err != nil {
		panic(err)
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
		for _, r := range receivers {
			rr := r.(ingest.Receiver)
			rr.Receive(buf.Bytes())
		}
		/*client := http.Client{}
		_, err = client.Get(fmt.Sprintf("http://api.mixpanel.com/track/?data=%s", rawData))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}*/
	})
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
