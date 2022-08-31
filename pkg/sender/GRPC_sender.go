package sender

import (
	"context"
	"fmt"
	"log"
	"time"

	uploadpb "receiving_service/pkg/proto"

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

func (c Client) StartServer(addr string, ch *chan []byte) {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c = NewClient(conn)

	func() {
		var file []byte
		for {
			file = <-*ch

			name, err := c.Upload(file, context.Background())
			if err != nil {
				name = []byte("could not make statistics")
			}

			*ch <- name
		}

	}()
}

func (c Client) Upload(file []byte, con context.Context) ([]byte, error) {

	ctx, cancel := context.WithDeadline(con, time.Now().Add(10*time.Second))
	defer cancel()

	stream, err := c.client.Upload(ctx)
	if err != nil {

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
