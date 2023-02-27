package main

import (
	"context"
	"testing"

	"github.com/elijah688/delio-mini/src/domain/position"
)

func BenchmarkNetProfitCalaculationSerial(b *testing.B) {

	// initialize inputs
	fhSvc, portfolio, positionChan, errChan := &MockFinnhubService{}, map[string]int{AAPL: 10, MSFT: 10}, make(chan position.Position), make(chan error)
	fhSvc.On("EvaluateHolding", context.Background(), AAPL, portfolio[AAPL]).Return(position.New(AAPL, 100.0), nil)
	fhSvc.On("EvaluateHolding", context.Background(), MSFT, portfolio[MSFT]).Return(position.New(MSFT, 200.0), nil)

	for n := 0; n < b.N; n++ {
		evaluatePortfolio(fhSvc, &portfolio, &positionChan, &errChan)
		renderPortfolioEvaluation(len(portfolio), &positionChan, &errChan)
	}

}

func BenchmarkNetProfitCalaculationSerialParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		// initialize inputs
		fhSvc, portfolio, positionChan, errChan := &MockFinnhubService{}, map[string]int{AAPL: 10, MSFT: 10}, make(chan position.Position), make(chan error)
		fhSvc.On("EvaluateHolding", context.Background(), AAPL, portfolio[AAPL]).Return(position.New(AAPL, 100.0), nil)
		fhSvc.On("EvaluateHolding", context.Background(), MSFT, portfolio[MSFT]).Return(position.New(MSFT, 200.0), nil)

		b.ReportAllocs()
		b.ResetTimer()

		for pb.Next() {
			evaluatePortfolio(fhSvc, &portfolio, &positionChan, &errChan)
			renderPortfolioEvaluation(len(portfolio), &positionChan, &errChan)
		}
	})
}
