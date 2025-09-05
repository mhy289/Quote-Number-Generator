package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"quote_generator/generatorpb"

	"github.com/bufbuild/connect-go"
)

type generatorService struct {
	generatorpb.UnimplementedGeneratorServiceHandler
	quotes []string
}

func (s *generatorService) GetRandomNumber(ctx context.Context, req *connect.Request[generatorpb.GetRandomNumberRequest]) (*connect.Response[generatorpb.GetRandomNumberResponse], error) {
	min := req.Msg.Min
	max := req.Msg.Max
	if max < min {
		min, max = max, min
	}
	number := rand.Int31n(max-min+1) + min
	return connect.NewResponse(&generatorpb.GetRandomNumberResponse{
		Number: number,
	}), nil
}

func (s *generatorService) GetRandomQuote(ctx context.Context, req *connect.Request[generatorpb.GetRandomQuoteRequest]) (*connect.Response[generatorpb.GetRandomQuoteResponse], error) {
	n := rand.Intn(len(s.quotes))
	return connect.NewResponse(&generatorpb.GetRandomQuoteResponse{
		Quote: s.quotes[n],
	}), nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	quotes := []string{
		"Stay hungry, stay foolish.",
		"The only limit to our realization of tomorrow is our doubts of today.",
		"Life is what happens when you're busy making other plans.",
		"Do not take life too seriously. You will never get out of it alive.",
	}

	svc := &generatorService{quotes: quotes}

	mux := http.NewServeMux()
	mux.Handle(generatorpb.NewGeneratorServiceHandler(svc))

	fmt.Println("Server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
