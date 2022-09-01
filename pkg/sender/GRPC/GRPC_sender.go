package GRPC

import (
	"context"
	"fmt"
	"log"
	"time"

	uploadpb "github.com/whym9/receiving_service/pkg/GRPC_gen"

	"google.golang.org/grpc"
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
	ch *chan []byte
}

func NewGRPCHandler(ch *chan []byte) Handler {
	return Handler{ch}
}

func (h Handler) StartServer(addr string) {
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

			name, err := cl.Upload(file, context.Background())
			if err != nil {
				name = []byte("could not make statistics")
			}

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
	//fmt.Println(res.GetName())

	return []byte(res.GetName()), nil
}
