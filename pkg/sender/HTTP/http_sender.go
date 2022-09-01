package HTTP

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/whym9/receiving_service/pkg/metrics"
)

type HTTP_Handler struct {
	metrics.Metrics
	ch *chan []byte
}

func NewHTTPHandler(ch *chan []byte) HTTP_Handler {
	return HTTP_Handler{ch: ch}
}

func (h HTTP_Handler) StartServer(addr string) {

	name := string(<-*h.ch)
	mes, err := h.Upload(addr, "POST", name)
	if err != nil {
		log.Fatal(err)

	}
	*h.ch <- mes

}

func (h HTTP_Handler) Upload(urlPath, method, filename string) ([]byte, error) {
	h.RecordMetrics()
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("uploadFile", filename)
	if err != nil {

		return []byte{}, err
	}
	file, err := os.Open(filename)
	if err != nil {

		return []byte{}, err
	}
	_, err = io.Copy(fw, file)
	if err != nil {

		return []byte{}, err
	}
	writer.Close()
	req, err := http.NewRequest(method, urlPath, bytes.NewReader(body.Bytes()))
	if err != nil {
		fmt.Println(".request")
		return []byte{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)
	ans := make([]byte, 1024)
	rsp.Body.Read(ans)

	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)

		return []byte{}, errors.New("response error")
	}
	return ans, nil
}
