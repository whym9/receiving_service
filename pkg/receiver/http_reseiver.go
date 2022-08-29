package receiver

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"receiving_service/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HTTP_Handler struct {
	transfer *chan []byte
}

func NewHTTPHandler(ch *chan []byte) HTTP_Handler {
	return HTTP_Handler{ch}
}

func (h HTTP_Handler) StartServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", h.Receive)
	fmt.Println("HTTP server has started")
	http.ListenAndServe(addr, nil)
}

func (h HTTP_Handler) Receive(w http.ResponseWriter, r *http.Request) {
	metrics.PromoHandler{}.RecordMetrics()
	if r.Method == "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong request method"))
		return
	}

	if err := r.ParseMultipartForm(100 * 1024 * 1024); err != nil {
		fmt.Printf("could not parse multipart form: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("CANT_PARSE_FORM"))
		return
	}

	file, fileHeader, err := r.FormFile("uploadFile")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Couldn't convert"))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		return
	}
	defer file.Close()

	fileSize := fileHeader.Size

	fileContent, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		return
	}

	fileType := http.DetectContentType(fileContent)
	if fileType != "application/octet-stream" {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Wrong file type!"))
		if err != nil {
			log.Fatal(err)
			return

		}
		return
	}

	fmt.Printf("FileType: %s, File: %s\n", fileType, fileHeader.Filename)
	fmt.Printf("File size (bytes): %v\n", fileSize)

	*h.transfer <- fileContent
	res := <-*h.transfer

	w.WriteHeader(http.StatusOK)
	w.Write(res)

	return
}
