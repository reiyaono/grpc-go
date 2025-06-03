package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"log"

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
	callDownload(client)
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
