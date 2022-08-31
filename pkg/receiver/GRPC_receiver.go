package receiver

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/whym9/receiving_service/pkg/metrics"
	uploadpb "github.com/whym9/receiving_service/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	uploadpb.UnimplementedUploadServiceServer
	tr *chan []byte
}

func NewServer(ch *chan []byte) Server {

	return Server{tr: ch}
}

func (s Server) StartServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	fmt.Println("GRPC server has started")
	if err != nil {

		log.Fatal(err)
	}
	defer lis.Close()

	uplSrv := NewServer(s.tr)

	rpcSrv := grpc.NewServer()

	uploadpb.RegisterUploadServiceServer(rpcSrv, uplSrv)

	log.Fatal(rpcSrv.Serve(lis))
}

func (s Server) Upload(stream uploadpb.UploadService_UploadServer) error {
	metrics.PromoHandler{}.RecordMetrics()
	fmt.Println("Got")
	chunk := []byte{}

	for {

		req, err := stream.Recv()

		if err == io.EOF {

			break
		}
		if err != nil {

			return status.Error(codes.Internal, err.Error())
		}

		bin := req.GetChunk()

		chunk = append(chunk, bin...)

		if err != nil {

			return status.Error(codes.Internal, err.Error())
		}

	}
	if s.tr == nil {
		return stream.SendAndClose(&uploadpb.UploadResponse{Name: "\n Error with statistics"})
	}

	*s.tr <- chunk

	mes := string(<-*s.tr)

	return stream.SendAndClose(&uploadpb.UploadResponse{Name: mes})
}
