package TCP

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type TCP_Handler struct {
	ch *chan []byte
}

func NewTCPHandler(ch *chan []byte) TCP_Handler {
	return TCP_Handler{ch}
}

func (t TCP_Handler) StartServer(addr string) {
	connect, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer connect.Close()
	func() {
		var file []byte
		for {
			file = <-*t.ch
			fmt.Println(len(file))
			name, err := t.Upload(file, connect)
			if err != nil {
				name = []byte("could not make statistics.")
				*t.ch <- name
				break
			}

			*t.ch <- name
		}

	}()

}

func (t TCP_Handler) Upload(file []byte, connect net.Conn) ([]byte, error) {
	be := 0
	en := 1024

	for {

		if en > len(file) {
			bin := make([]byte, 8)
			binary.BigEndian.PutUint64(bin, uint64(len(file)-be))

			if _, err := connect.Write(bin); err != nil {

				return []byte{}, err
			}

			if _, err := connect.Write(file[be:]); err != nil {

				return []byte{}, err
			}
			bin = make([]byte, 8)

			binary.BigEndian.PutUint64(bin, uint64(4))

			if _, err := connect.Write(bin); err != nil {

				return []byte{}, err
			}

			if _, err := connect.Write([]byte("STOP")); err != nil {

				return []byte{}, err
			}

			break
		}
		bin := make([]byte, 8)
		binary.BigEndian.PutUint64(bin, uint64(1024))

		if _, err := connect.Write(bin); err != nil {

			return []byte{}, err
		}

		if _, err := connect.Write(file[be:en]); err != nil {

			return []byte{}, err
		}

		be = en
		en += 1024
	}

	read := make([]byte, 1024)

	_, err := connect.Read(read)

	if err != nil {

		return []byte{}, err
	}

	connect.Close()
	return read, nil
}
