package main

import (
	"binance-price/sub"
	"context"
	"fmt"
	"sync"
)

func main() {
	// 币安WebSocket URL
	contractUrl := "wss://fstream.binance.com/stream?streams=btcusdt@aggTrade/ethusdt@aggTrade/trxusdt@aggTrade/solusdt@aggTrade"
	spotUrl := "wss://stream.binance.com:9443/stream?streams=btcusdt@ticker/ethusdt@ticker/trxusdt@ticker/solusdt@ticker"

	ctx := context.Background()

	// 使用WaitGroup确保程序持续运行
	var wg sync.WaitGroup
	wg.Add(2)

	// 订阅合约
	go sub.SubscribeContract(ctx, contractUrl, &wg)
	// 订阅现货
	go sub.SubscribeSpot(ctx, spotUrl, &wg)

	// 保持程序运行
	fmt.Println("已连接到币安WebSocket，开始订阅...")
	wg.Wait()
}
