package main

import (
	"context"
	"testing"

	connect "connectrpc.com/connect"
	generatorpb "quote_generator/generatorpb"
)

func TestGetRandomNumber(t *testing.T) {
	svc := &generatorService{}

	// 测试 min > max 应该返回错误
	req := connect.NewRequest(&generatorpb.GetRandomNumberRequest{Min: 10, Max: 5})
	_, err := svc.GetRandomNumber(context.Background(), req)
	if err == nil {
		t.Fatal("expected error when min > max")
	}

	// 测试正常范围
	req = connect.NewRequest(&generatorpb.GetRandomNumberRequest{Min: 1, Max: 10})
	resp, err := svc.GetRandomNumber(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Msg.Number < 1 || resp.Msg.Number > 10 {
		t.Fatalf("number %d out of range", resp.Msg.Number)
	}
}

func TestGetRandomQuote(t *testing.T) {
	quotes := []string{
		"Stay hungry, stay foolish.",
		"Life is what happens when you're busy making other plans.",
		"The only limit to our realization of tomorrow is our doubts of today.",
	}
	svc := &generatorService{quotes: quotes}

	req := connect.NewRequest(&generatorpb.GetRandomQuoteRequest{})
	resp, err := svc.GetRandomQuote(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 检查返回的名言是否在预设列表中
	found := false
	for _, q := range quotes {
		if resp.Msg.Quote == q {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("unexpected quote: %s", resp.Msg.Quote)
	}
}

func TestGetRandomNumberMultipleTimes(t *testing.T) {
	svc := &generatorService{}

	// 连续调用多次，确保随机数在范围内
	for i := 0; i < 100; i++ {
		req := connect.NewRequest(&generatorpb.GetRandomNumberRequest{Min: 5, Max: 15})
		resp, err := svc.GetRandomNumber(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Msg.Number < 5 || resp.Msg.Number > 15 {
			t.Fatalf("number %d out of range", resp.Msg.Number)
		}
	}
}
