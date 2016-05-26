package socketout

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/t15k/go-ingest"
)

const retries = 10

func init() {
	ingest.RegisterMod("socketout", func(cfg string) interface{} {
		return &Receiver{cfg, nil}
	})
}

// Receiver takes string and send to host. I.e localhost:4333.
type Receiver struct {
	host string
	conn net.Conn
}

func (r *Receiver) connect() (err error) {
	for i := 0; i < retries && r.conn == nil; i++ {
		if err != nil {
			log.Print("Failed to connect to ", r.host, " will retry in 5 seconds.")
			c := time.After(5 * time.Second)
			<-c
		}
		r.conn, err = net.Dial("tcp", r.host)
	}
	return
}

// Receive accepts a string and sends to host.
func (r *Receiver) Receive(v interface{}) {
	err := r.connect()
	if err != nil {
		panic(fmt.Sprintf("Could not connect after %d retries, because %s.", retries, err))
	}
	r.conn.Write(v.([]byte))
}
