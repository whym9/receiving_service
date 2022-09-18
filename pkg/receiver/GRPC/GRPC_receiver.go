package GRPC

import (
	"fmt"
	"io"
	"log"
	"net"

	//"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promauto"
	uploadpb "github.com/whym9/receiving_service/pkg/GRPC_gen"
	"github.com/whym9/receiving_service/pkg/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// var (
// 	name   = "GRPC_receiver_processed_errors_total"
// 	help   = "The total number of receiver errors"
// 	errors prometheus.Counter
// )

type Server struct {
	uploadpb.UnimplementedUploadServiceServer
	metrics metrics.Metrics
	tr      chan []byte
}

func NewServer(m metrics.Metrics, ch chan []byte) Server {
	// errors = promauto.NewCounter(prometheus.CounterOpts{
	// 	Name: name,
	// 	Help: help,
	// })
	return Server{metrics: m, tr: ch}
}

func (s Server) StartServer(addr string) {

	lis, err := net.Listen("tcp", addr)
	fmt.Println("GRPC server has started")
	if err != nil {
		//errors.Inc()
		log.Fatal(err)
	}
	defer lis.Close()

	uplSrv := NewServer(s.metrics, s.tr)

	rpcSrv := grpc.NewServer()

	uploadpb.RegisterUploadServiceServer(rpcSrv, uplSrv)

	log.Fatal(rpcSrv.Serve(lis))
}

func (s Server) Upload(stream uploadpb.UploadService_UploadServer) error {
	s.metrics.RecordMetrics()
	fmt.Println("Got")
	chunk := []byte{}

	for {

		req, err := stream.Recv()

		if err == io.EOF {

			break
		}
		if err != nil {
			//errors.Inc()
			return status.Error(codes.Internal, err.Error())
		}

		bin := req.GetChunk()

		chunk = append(chunk, bin...)

		if err != nil {
			//errors.Inc()
			return status.Error(codes.Internal, err.Error())
		}

	}
	if s.tr == nil {
		//errors.Inc()
		return stream.SendAndClose(&uploadpb.UploadResponse{Name: "\n Error with statistics"})
	}

	s.tr <- chunk

	mes := string(<-s.tr)

	return stream.SendAndClose(&uploadpb.UploadResponse{Name: mes})
}
