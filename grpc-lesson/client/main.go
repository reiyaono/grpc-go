package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// 通信が暗号化されないため本番では使わない
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	// unaryClientとして使う場合は以下のように呼び出す
	// callListFiles(client)

	// serverStreamClientとして使う場合は以下のように呼び出す
	// callDownload(client)

	// clientStreamClientとして使う場合は以下のように呼び出す
	CallUpload(client)
}

func callListFiles(client pb.FileServiceClient) {
	req := &pb.ListFilesRequest{}
	res, err := client.ListFiles(context.Background(), req)
	if err != nil {
		log.Fatalf("could not list files: %v", err)
	}

	fmt.Println((res.GetFilenames()))
}

func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{
		Filename: "name.txt",
	}

	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalf("could not download file: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// すべてのデータを受信した場合
			break
		}
		if err != nil {
			log.Fatalf("error receiving file data: %v", err)
		}
		log.Printf("Response from Download(bytes): %v", res.GetData())
		log.Printf("Response from Download(string): %v", string(res.GetData()))
	}
}

func CallUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/018018s/dev/go-grpc/grpc-lesson/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("could not upload file: %v", err)
	}

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			// すべてのデータを送信した場合
			break
		}
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}

		req := &pb.UploadRequest{
			Data: buf[:n],
		}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalf("error sending data: %v", sendErr)
		}

		time.Sleep(1 * time.Second) // Simulate delay for streaming
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving upload response: %v", err)
	}
	log.Printf("Upload completed, total size: %d bytes", res.GetSize())

}
