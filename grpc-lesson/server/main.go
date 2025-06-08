package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles called")

	dir := "/Users/018018s/dev/go-grpc/grpc-lesson/storage"

	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Filenames: filenames,
	}
	return res, nil
}

func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download called")
	filename := req.GetFilename()
	filePath := fmt.Sprintf("/Users/018018s/dev/go-grpc/grpc-lesson/storage/%s", filename)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	buffer := make([]byte, 5)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		res := &pb.DownloadResponse{
			Data: buffer[:n],
		}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second) // Simulate delay for streaming
	}
	return nil
}

func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload called")

	var buf bytes.Buffer
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// すべてのデータを受信した場合
			res := &pb.UploadResponse{Size: int32(buf.Len())}
			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("Received data: %v", data)
		log.Printf("Received data: %v", string(data))
		buf.Write(data)
	}
}

func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress called")
	size := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("Received data: %v", data)
		log.Printf("Received data: %v", string(data))
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("Received %d bytes", size),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}

func myLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("Request - Method:%s, Request:%v", info.FullMethod, req)
		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("Error - Method:%s, Error:%v", info.FullMethod, err)
			return nil, err
		}
		log.Printf("Response - Method:%s, Response:%v, Error:%v", info.FullMethod, resp, err)
		return resp, nil
	}
}

func authorize(ctx context.Context) (context.Context, error) {
	// MD = Metadata
	// ここでは認証のためのメタデータを取得します。
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %v", err)
	}

	if token != "test-token" {
		return nil, errors.New("unauthorized: invalid token")
	}
	return ctx, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			myLogging(),
			grpc_auth.UnaryServerInterceptor(authorize))))
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
