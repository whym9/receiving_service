package GRPC

import (
	"context"
	"fmt"
	"log"
	"time"

	uploadpb "github.com/whym9/receiving_service/pkg/GRPC_gen"

	"github.com/whym9/receiving_service/pkg/metrics"

	"google.golang.org/grpc"
)

var (
	name1 = "GPRC_sent_processed_opts_total"
	help1 = "The total number of sending requsets"

	name2 = "GRPC_sending_processed_errors_total"
	help2 = "The total number of sender errors"

	key1 = "sendmetrics"
	key2 = "errormetrics"
)

type Client struct {
	client uploadpb.UploadServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return Client{
		client: uploadpb.NewUploadServiceClient(conn),
	}
}

type Handler struct {
	metrics metrics.Metrics
	ch      *chan []byte
}

func NewGRPCHandler(m metrics.Metrics, ch *chan []byte) Handler {
	return Handler{metrics: m, ch: ch}
}

func (h Handler) StartServer(addr string) {
	h.metrics.AddMetrics(name1, help1, key1)
	h.metrics.AddMetrics(name2, help2, key2)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	cl := NewClient(conn)

	func() {
		var file []byte
		for {
			file = <-*h.ch
			h.metrics.RecordMetrics()
			name, err := cl.Upload(file, context.Background())
			if err != nil {
				h.metrics.Count(key2)
				name = []byte("could not make statistics")
			}
			h.metrics.Count(key1)
			*h.ch <- name
		}

	}()
}

func (c Client) Upload(file []byte, con context.Context) ([]byte, error) {

	ctx, cancel := context.WithDeadline(con, time.Now().Add(10*time.Second))
	defer cancel()

	stream, err := c.client.Upload(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return []byte{}, err
	}
	be := 0
	en := 1024

	for {

		if en > len(file) {
			if err := stream.Send(&uploadpb.UploadRequest{Chunk: file[be:]}); err != nil {
				return []byte{}, err
			}
			break
		}

		if err := stream.Send(&uploadpb.UploadRequest{Chunk: file[be:en]}); err != nil {
			return []byte{}, err
		}

		be = en
		en += 1024
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return []byte{}, err
	}

	fmt.Println("stopped sending")

	return []byte(res.GetName()), nil
}
