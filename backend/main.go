package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"connectrpc.com/connect"
)

// ===== 请求/响应结构体（JSON 字段小写，匹配前端） =====
type GetRandomNumberRequest struct {
	Min int32 `json:"min"`
	Max int32 `json:"max"`
}
type GetRandomNumberResponse struct {
	Number int32 `json:"number"`
}

type GetRandomQuoteRequest struct{}
type GetRandomQuoteResponse struct {
	Quote string `json:"quote"`
}

// ===== handler 实现（与 connect.NewUnaryHandler 的泛型签名匹配） =====
func handleGetRandomNumber(ctx context.Context, req *connect.Request[GetRandomNumberRequest]) (*connect.Response[GetRandomNumberResponse], error) {
	min, max := req.Msg.Min, req.Msg.Max
	if max < min {
		max = min
	}
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := src.Int31n(max-min+1) + min
	return connect.NewResponse(&GetRandomNumberResponse{Number: n}), nil
}

func handleGetRandomQuote(ctx context.Context, req *connect.Request[GetRandomQuoteRequest]) (*connect.Response[GetRandomQuoteResponse], error) {
	quotes := []string{
		"Stay hungry, stay foolish.",
		"Life is what happens when you're busy making other plans.",
		"Be yourself; everyone else is already taken.",
	}
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	q := quotes[src.Intn(len(quotes))]
	return connect.NewResponse(&GetRandomQuoteResponse{Quote: q}), nil
}

// 简单 CORS 包装（允许来自 http://localhost:3000 的前端）
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	// v1.18.1: NewUnaryHandler 返回 *connect.Handler（单返回值）
	numHandler := connect.NewUnaryHandler[GetRandomNumberRequest, GetRandomNumberResponse](
		"/generator.v1.GeneratorService/GetRandomNumber",
		handleGetRandomNumber,
	)

	quoteHandler := connect.NewUnaryHandler[GetRandomQuoteRequest, GetRandomQuoteResponse](
		"/generator.v1.GeneratorService/GetRandomQuote",
		handleGetRandomQuote,
	)

	// 注册到 mux（路径与 NewUnaryHandler 中的第一个参数保持一致）
	mux.Handle("/generator.v1.GeneratorService/GetRandomNumber", withCORS(numHandler))
	mux.Handle("/generator.v1.GeneratorService/GetRandomQuote", withCORS(quoteHandler))

	log.Println("ConnectRPC (JSON) server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
