package main

import (
	"BakdauletKan/midka/midka_pb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type CalculatorService struct {

	midka_pb.UnimplementedCalculatorServiceServer
}

func (s *CalculatorService) PrimeNumberDecomposition(req *midka_pb.NumberRequest, stream midka_pb.CalculatorService_PrimeNumberDecompositionServer) error {
	num := req.GetNumber()
	primes := decomposeToPrimes(num)
	for i := 0; i < len(primes); i++ {
		response := &midka_pb.NumberResponse{Result: primes[i]}
		err := stream.Send(response)
		if err != nil {
			log.Fatal("Error while sending %v", err.Error())
		}
	}
	return nil
}

//Implemented from here https://www.geeksforgeeks.org/print-all-prime-factors-of-a-given-number/
func decomposeToPrimes(n int32) []int32 {
	var primeFactors []int32

	for n%2 == 0 {
		primeFactors = append(primeFactors, 2)
		n = n / 2
	}

	var i int32 = 0
	for i = 3; i*i <= n; i = i + 2 {
		for n%i == 0 {
			primeFactors = append(primeFactors, i)
			n = n / i
		}
	}

	if n > 2 {
		primeFactors = append(primeFactors, n)
	}

	return primeFactors
}

func (s *CalculatorService) ComputerAverage(stream midka_pb.CalculatorService_ComputerAverageServer) error {
	var streamSum int32 = 0
	var counter int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("%v", float64(streamSum/counter))
			response := &midka_pb.AverageResponse{Result: float64(streamSum) / float64(counter)}
			return stream.SendAndClose(response)
		}
		if err != nil {
			log.Fatal("Error with stream %v", err)
		}
		streamSum += req.GetNumber()
		counter++
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	midka_pb.RegisterCalculatorServiceServer(s, &CalculatorService{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
