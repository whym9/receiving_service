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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/whym9/receiving_service/pkg/metrics"
)

var (
	name1 = "HTTP_sent_processed_opts_total"
	help1 = "The total number of sending requsets"

	name2 = "HTTP_sending_processed_errors_total"
	help2 = "The total number of sender errors"

	sent prometheus.Counter
	errs prometheus.Counter
)

type HTTP_Handler struct {
	metrics metrics.Metrics
	ch      chan []byte
}

func NewHTTPHandler(m metrics.Metrics, ch chan []byte) HTTP_Handler {
	return HTTP_Handler{metrics: m, ch: ch}
}

func (h HTTP_Handler) StartServer(addr string) {
	sent = promauto.NewCounter(prometheus.CounterOpts{
		Name: name1,
		Help: help1,
	})

	errs = promauto.NewCounter(prometheus.CounterOpts{
		Name: name2,
		Help: help2,
	})
	name := string(<-h.ch)
	mes, err := h.Upload(addr, "POST", name)
	if err != nil {
		log.Fatal(err)

	}
	h.ch <- mes

}

func (h HTTP_Handler) Upload(urlPath, method, filename string) ([]byte, error) {
	sent.Inc()
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("uploadFile", filename)
	if err != nil {
		errs.Inc()
		return []byte{}, err
	}
	file, err := os.Open(filename)
	if err != nil {
		errs.Inc()
		return []byte{}, err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		errs.Inc()
		return []byte{}, err
	}
	writer.Close()
	req, err := http.NewRequest(method, urlPath, bytes.NewReader(body.Bytes()))
	if err != nil {
		errs.Inc()
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
