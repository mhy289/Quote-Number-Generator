package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/cors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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
		// 返回 ConnectRPC 参数错误
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("最小值不能大于最大值"))
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

// loadQuotes 从外部 JSON 文件加载名言列表，文件不存在则使用默认值
func loadQuotes(path string) []string {
	data, err := os.ReadFile(path)
	if err == nil {
		var quotes []string
		if json.Unmarshal(data, &quotes) == nil && len(quotes) > 0 {
			return quotes
		}
	}
	// 默认名言列表
	return []string{
		"Stay hungry, stay foolish.",
		"Life is what happens when you're busy making other plans.",
		"The only limit to our realization of tomorrow is our doubts of today.",
		"Do not take life too seriously. You will never get out of it alive.",
	}
}

func main() {
	// 使用时间戳作为随机种子，确保每次启动序列不同
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 从外部文件加载名言（如果存在）
	quotes := loadQuotes("data/quotes.json")

	// 初始化服务
	svc := &generatorService{
		quotes: quotes,
	}

	// 创建 HTTP mux
	mux := http.NewServeMux()

	// 注册 ConnectRPC 服务
	path, handler := generatorpbconnect.NewGeneratorServiceHandler(svc)
	mux.Handle(path, handler)

	// CORS 支持
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // 前端地址
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	fmt.Println("Server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(mux)))
}
