package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

const (
	AAPL     = "AAPL"
	MSFT     = "MSFT"
	FH_TOKEN = "FH_TOKEN"
)

func main() {
	// Config FH Client
	cfg := finnhub.NewConfiguration()
	auth := os.Getenv(FH_TOKEN)
	if auth == "" {
		panic("finnhub authentication token not set")
	}
	cfg.AddDefaultHeader("X-Finnhub-Token", auth)
	client := finnhub.NewAPIClient(cfg).DefaultApi
	positionChan, errChan := make(chan *Position), make(chan error)
	fhClient := new(FHClient).Set(client, positionChan, errChan)

	// Evaluate portfolio
	portfolio := map[string]int{AAPL: 10, MSFT: 10}
	for s, v := range portfolio {
		symbol, volume := s, v
		go func() {
			fhClient.GetQuote(context.Background(), symbol, volume)
		}()
	}

	// Render portfolio evaluation
	for range portfolio {
		select {
		case position := <-positionChan:
			ppr(position)
		case err := <-errChan:
			panic(err)
		}
	}
}

type FHClient struct {
	client       *finnhub.DefaultApiService
	positionChan chan *Position
	errChan      chan error
}

func (self *FHClient) Set(
	client *finnhub.DefaultApiService,
	positionChan chan *Position,
	errChan chan error,
) *FHClient {
	self.client = client
	self.errChan = errChan
	self.positionChan = positionChan

	return self
}

type Position struct {
	NetProfit float32
	Symbol    string
}

func (p *Position) Set(symbol string, netProfit float32) *Position {
	p.NetProfit = netProfit
	p.Symbol = symbol

	return p
}

func (self *FHClient) GetQuote(ctx context.Context, symbol string, volume int) {
	if aapl, _, err := self.client.Quote(context.Background()).Symbol(symbol).Execute(); err != nil {
		self.errChan <- err
	} else {
		pc, c := aapl.GetC(), aapl.GetPc()
		netProfit := (c - pc) * float32(volume)
		self.positionChan <- new(Position).Set(symbol, netProfit)
	}
}

func handleErr(err error) {
	if err != nil {
		panic(fmt.Errorf("%w", err))
	}
}

func ppr(res *Position) {
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("%s ", res.Symbol))
	if res.NetProfit >= 0 {
		sb.WriteString(fmt.Sprintf("Changed(1D)\t▲ Profit: %.2f$", res.NetProfit))
	} else {
		sb.WriteString(fmt.Sprintf("Changed(1D)\t▼ Loss: %.2f$", res.NetProfit))
	}

	fmt.Println(sb.String())
}
func clearCurrentLine() {
	fmt.Print("\n\033[1A\033[K")
}
