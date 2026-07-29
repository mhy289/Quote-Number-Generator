package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/rs/cors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	// 命令行参数定义
	mode := flag.String("mode", "all", "启动模式: http / console / all")
	port := flag.Int("port", 8080, "HTTP 服务端口")
	flag.Parse()

	// 从外部文件加载名言（如果存在）
	quotes := loadQuotes("data/quotes.json")

	// 初始化服务
	svc := &generatorService{
		quotes: quotes,
	}

	switch *mode {
	case "http":
		startHTTPService(svc, *port)
	case "console":
		runConsoleMenu(svc)
	case "all":
		startHTTPService(svc, *port)
		runConsoleMenu(svc)
	default:
		fmt.Printf("❌ 未知模式: %s，可用模式: http / console / all\n", *mode)
		os.Exit(1)
	}
}

// startHTTPService 启动 HTTP 服务
func startHTTPService(svc *generatorService, port int) {
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

	addr := fmt.Sprintf(":%d", port)
	go func() {
		fmt.Printf("🌐 HTTP 服务已启动: http://localhost%s\n", addr)
		log.Fatal(http.ListenAndServe(addr, c.Handler(mux)))
	}()
}

// runConsoleMenu 提供交互式控制台菜单
func runConsoleMenu(svc *generatorService) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n========================================")
	fmt.Println("  🎯 QuoteGenerator 控制台交互模式")
	fmt.Println("========================================")
	fmt.Println("  1. 生成随机数")
	fmt.Println("  2. 获取随机名言")
	fmt.Println("  0. 退出")
	fmt.Println("========================================")

	for {
		fmt.Print("\n请输入数字选择操作: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("❌ 请输入有效的数字")
			continue
		}

		switch choice {
		case 1:
			handleConsoleRandomNumber(scanner)
		case 2:
			handleConsoleQuote(svc)
		case 0:
			fmt.Println("👋 已退出，服务仍在后台运行")
			return
		default:
			fmt.Println("❌ 无效选项，请输入 0-2")
		}
	}
}

// handleConsoleRandomNumber 控制台交互：生成随机数
func handleConsoleRandomNumber(scanner *bufio.Scanner) {
	fmt.Print("  请输入最小值: ")
	if !scanner.Scan() {
		return
	}
	minStr := strings.TrimSpace(scanner.Text())
	min, err := strconv.Atoi(minStr)
	if err != nil {
		fmt.Println("❌ 无效的最小值")
		return
	}

	fmt.Print("  请输入最大值: ")
	if !scanner.Scan() {
		return
	}
	maxStr := strings.TrimSpace(scanner.Text())
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		fmt.Println("❌ 无效的最大值")
		return
	}

	if max < min {
		fmt.Println("❌ 最小值不能大于最大值")
		return
	}

	number := rand.Int31n(int32(max-min+1)) + int32(min)
	fmt.Printf("✅ 随机数: %d (范围: %d ~ %d)\n", number, min, max)
}

// handleConsoleQuote 控制台交互：获取随机名言
func handleConsoleQuote(svc *generatorService) {
	n := rand.Intn(len(svc.quotes))
	fmt.Printf("💬 名言: \"%s\"\n", svc.quotes[n])
}
