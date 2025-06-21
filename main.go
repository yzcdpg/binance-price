package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

/*

	{
		"stream": "btcusdt@aggTrade",
		"data": {
			"e": "aggTrade",
			"E": 1750470348995,
			"a": 2764705375,
			"s": "BTCUSDT",
			"p": "103395.60",
			"q": "0.030",
			"f": 6415259597,
			"l": 6415259597,
			"T": 1750470348970,
			"m": false
		}
	}

*/

type ContractStreamData struct {
	Stream string `json:"stream"`
	Data   struct {
		E         string      `json:"e"`
		Timestamp json.Number `json:"E"`
		A         int64       `json:"a"`
		Symbol    string      `json:"s"`
		Price     json.Number `json:"p"`
		Q         string      `json:"q"`
		F         int64       `json:"f"`
		L         int64       `json:"l"`
		T         int64       `json:"T"`
		M         bool        `json:"m"`
	} `json:"data"`
}

/*

	{
		"stream": "ethçusdt@ticker",
		"data": {
			"e": "24hrTicker",
			"E": 1750412723024,
			"s": "ETHUSDT",
			"p": "21.50000000",
			"P": "0.850",
			"w": "2525.00764461",
			"x": "2530.40000000",
			"c": "2551.89000000",
			"Q": "0.20210000",
			"b": "2551.88000000",
			"B": "24.87050000",
			"a": "2551.89000000",
			"A": "114.39510000",
			"o": "2530.39000000",
			"h": "2569.00000000",
			"l": "2485.03000000",
			"v": "292887.90940000",
			"q": "739544210.24849600",
			"O": 1750326322968,
			"C": 1750412722968,
			"F": 2546217017,
			"L": 2547931514,
			"n": 1714498
		}
	}

*/

type SpotStreamData struct {
	Stream string `json:"stream"`
	Data   struct {
		E         string      `json:"e"`
		Timestamp json.Number `json:"E"`
		Symbol    string      `json:"s"`
		P         string      `json:"p"`
		P1        string      `json:"P"`
		W         string      `json:"w"`
		X         string      `json:"x"`
		Price     json.Number `json:"c"`
		Q         string      `json:"Q"`
		B         string      `json:"b"`
		B1        string      `json:"B"`
		A         string      `json:"a"`
		A1        string      `json:"A"`
		O         string      `json:"o"`
		H         string      `json:"h"`
		L         string      `json:"l"`
		V         string      `json:"v"`
		Q1        string      `json:"q"`
		O1        json.Number `json:"O"`
		C1        json.Number `json:"C"`
		F         json.Number `json:"F"`
		L1        json.Number `json:"L"`
		N         json.Number `json:"n"`
	} `json:"data"`
}

func main() {
	// 币安WebSocket URL
	contractUrl := "wss://fstream.binance.com/stream?streams=btcusdt@aggTrade/ethusdt@aggTrade/trxusdt@aggTrade/solusdt@aggTrade"
	spotUrl := "wss://stream.binance.com:9443/stream?streams=btcusdt@ticker/ethusdt@ticker/trxusdt@ticker/solusdt@ticker"

	// 使用WaitGroup确保程序持续运行
	var wg sync.WaitGroup
	wg.Add(2)

	//处理接收到的消息
	go func() {
		defer wg.Done()

		// 连接WebSocket
		contractConn, _, err := websocket.DefaultDialer.Dial(contractUrl, nil)
		if err != nil {
			log.Fatalf("连接合约ws错误: %v", err)
		}
		defer contractConn.Close()

		for {
			_, message, err := contractConn.ReadMessage()
			if err != nil {
				log.Printf("读取消息欧呜: %v", err)
				return
			}
			//fmt.Println(string(message))
			// 解析JSON消息
			var data ContractStreamData
			if err := json.Unmarshal(message, &data); err != nil {
				log.Printf("解析JSON错误: %v", err)
				continue
			}
			timestamp, err := data.Data.Timestamp.Int64()
			if err != nil {
				log.Printf("转换时间戳错误: %v", err)
				continue
			}
			price, err := data.Data.Price.Float64()
			if err != nil {
				log.Printf("转换价格错误: %v", err)
				continue
			}
			switch data.Data.Symbol {
			case "BTCUSDT":
				fmt.Printf("[c]交易对: %8s\t 最新价格: %16.08f\t 时间戳: %d\n", data.Data.Symbol, price, timestamp)
			case "ETHUSDT":
				fmt.Printf("%s[c]交易对: %8s\t 最新价格: %16.08f\t 时间戳: %d%s\n", Green, data.Data.Symbol, price, timestamp, Reset)
			case "TRXUSDT":
				fmt.Printf("%s[c]交易对: %8s\t 最新价格: %16.08f\t 时间戳: %d%s\n", Blue, data.Data.Symbol, price, timestamp, Reset)
			case "SOLUSDT":
				fmt.Printf("%s[c]交易对: %8s\t 最新价格: %16.08f\t 时间戳: %d%s\n", Red, data.Data.Symbol, price, timestamp, Reset)
			}
		}
	}()

	go func() {
		defer wg.Done()

		// 连接WebSocket
		spotConn, _, err := websocket.DefaultDialer.Dial(spotUrl, nil)
		if err != nil {
			log.Fatalf("连接现货ws错误: %v", err)
		}
		defer spotConn.Close()

		for {
			_, message, err := spotConn.ReadMessage()
			if err != nil {
				log.Printf("读取消息错误: %v", err)
				return
			}
			//fmt.Println(string(message))
			// 解析JSON消息
			var data SpotStreamData
			if err := json.Unmarshal(message, &data); err != nil {
				log.Printf("解析JSON错误: %v", err)
				continue
			}

			timestamp, err := data.Data.Timestamp.Int64()
			if err != nil {
				log.Printf("转换时间戳错误: %v", err)
				continue
			}
			price, err := data.Data.Price.Float64()
			if err != nil {
				log.Printf("转换价格错误: %v", err)
				continue
			}
			switch data.Data.Symbol {
			case "BTCUSDT":
				fmt.Printf("[s]交易对: %8s\t 最新价格: %16.08f\t 时间戳: %d \n", data.Data.Symbol, price, timestamp)
			case "ETHUSDT":
				fmt.Printf("%s[s]交易对: %8s\t 最新价格: %16.08f\t 时间戳: %d %s\n", Green, data.Data.Symbol, price, timestamp, Reset)
			case "TRXUSDT":
				fmt.Printf("%s[s]交易对: %8s\t 最新价格: %16.08f\t 时间戳: %d %s\n", Blue, data.Data.Symbol, price, timestamp, Reset)
			case "SOLUSDT":
				fmt.Printf("%s[s]交易对: %8s\t 最新价格: %16.08f\t 时间戳: %d %s\n", Red, data.Data.Symbol, price, timestamp, Reset)
			}
		}
	}()

	// 保持程序运行
	fmt.Println("已连接到币安WebSocket，开始订阅...")
	wg.Wait()
}
