package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	connect "connectrpc.com/connect"                                    // Request/Response/NewResponse
	generatorpb "quote_generator/generatorpb"                           // protobuf 消息类型
	generatorpbconnect "quote_generator/generatorpb/generatorpbconnect" // ConnectRPC Handler
)

// generatorService 实现 GeneratorServiceHandler 接口
type generatorService struct {
	quotes []string
}

// GetRandomNumber 返回指定范围内的随机整数
func (s *generatorService) GetRandomNumber(
	ctx context.Context,
	req *connect.Request[generatorpb.GetRandomNumberRequest],
) (*connect.Response[generatorpb.GetRandomNumberResponse], error) {
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

// GetRandomQuote 返回随机名言
func (s *generatorService) GetRandomQuote(
	ctx context.Context,
	req *connect.Request[generatorpb.GetRandomQuoteRequest],
) (*connect.Response[generatorpb.GetRandomQuoteResponse], error) {
	n := rand.Intn(len(s.quotes))
	return connect.NewResponse(&generatorpb.GetRandomQuoteResponse{
		Quote: s.quotes[n],
	}), nil
}

func main() {
	rand.Seed(0)

	// 初始化服务
	svc := &generatorService{
		quotes: []string{
			"Stay hungry, stay foolish.",
			"Life is what happens when you're busy making other plans.",
			"The only limit to our realization of tomorrow is our doubts of today.",
			"Do not take life too seriously. You will never get out of it alive.",
		},
	}

	// 创建 HTTP mux
	mux := http.NewServeMux()

	// 注册 ConnectRPC 服务
	path, handler := generatorpbconnect.NewGeneratorServiceHandler(svc)
	mux.Handle(path, handler)

	fmt.Println("Server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
