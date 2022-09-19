package HTTP

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/whym9/receiving_service/pkg/metrics"
)

type HTTP_Handler struct {
	metrics  metrics.Metrics
	transfer chan []byte
}

var (
	name   = "HTTP_receiver_processed_errors_total"
	help   = "The total number of receiver errors"
	errors prometheus.Counter
)

func NewHTTPHandler(m metrics.Metrics, ch chan []byte) HTTP_Handler {
	return HTTP_Handler{m, ch}
}

func (h HTTP_Handler) StartServer() {
	addr := os.Getenv("HTTP_RECEIVER")

	errors = promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
	http.HandleFunc("/", h.Receive)
	fmt.Println("HTTP server has started")
	http.ListenAndServe(addr, nil)
}

func (h HTTP_Handler) Receive(w http.ResponseWriter, r *http.Request) {

	h.metrics.RecordMetrics()
	if r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong request method"))
		errors.Inc()
		return
	}

	if err := r.ParseMultipartForm(100 * 1024 * 1024); err != nil {
		fmt.Printf("could not parse multipart form: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("CANT_PARSE_FORM"))
		errors.Inc()
		return
	}

	file, fileHeader, err := r.FormFile("uploadFile")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Couldn't convert"))
		errors.Inc()
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		errors.Inc()
		return
	}
	defer file.Close()

	fileSize := fileHeader.Size

	fileContent, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		errors.Inc()
		return
	}

	fileType := http.DetectContentType(fileContent)
	if fileType != "application/octet-stream" {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Wrong file type!"))
		errors.Inc()
		if err != nil {
			log.Printf("We have this error: %v", err)
			return

		}
		return
	}

	fmt.Printf("FileType: %s, File: %s\n", fileType, fileHeader.Filename)
	fmt.Printf("File size (bytes): %v\n", fileSize)

	h.transfer <- fileContent
	res := <-h.transfer

	w.WriteHeader(http.StatusOK)
	w.Write(res)

	return
}
