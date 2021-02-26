package main

import (
	"BakdauletKan/midka/midka_pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer connection.Close()

	calculatorService := midka_pb.NewCalculatorServiceClient(connection)
	runPrimeNumberDecomposition(120, calculatorService)
	runComputeAverage([]int32{1, 2, 3, 4}, calculatorService)
}

func runPrimeNumberDecomposition(n int32, calcService midka_pb.CalculatorServiceClient) {
	req := &midka_pb.NumberRequest{Number: n}
	stream, err := calcService.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatal("Error 500")
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("Error with response from server stream RPC %v", err)
		}
		log.Printf(fmt.Sprint(res.GetResult(), " "))
	}
}

func runComputeAverage(numbersStream []int32, calcService midka_pb.CalculatorServiceClient) {
	stream, err := calcService.ComputerAverage(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to server")
	}
	for _, n := range numbersStream {
		stream.Send(&midka_pb.NumberRequest{Number: n})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error recieving response")
	}

	fmt.Printf("Result: %f", res.GetResult())
}
