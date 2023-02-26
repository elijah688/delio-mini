package main

import (
	"context"
	"fmt"

	"github.com/elijah688/delio-mini/src/config"
	"github.com/elijah688/delio-mini/src/domain/position"
	"github.com/elijah688/delio-mini/src/service/finnhub"
)

const (
	AAPL = "AAPL"
	MSFT = "MSFT"
)

func main() {
	// initialize FH Config and comm channels
	cfg, portfolio, positionChan, errChan := config.New(), map[string]int{AAPL: 10, MSFT: 10}, make(chan position.Position), make(chan error)

	// config FH Service
	fhSvc := finnhub.New(cfg)

	// evaluate portfolio and render evaluation
	evaluatePortfolio(fhSvc, &portfolio, &positionChan, &errChan)
	renderPortfolioEvaluation(len(portfolio), &positionChan, &errChan)
}

func evaluatePortfolio(fhSvc finnhub.FinnhubService, portfolio *map[string]int, positionChan *chan position.Position, errChan *chan error) {
	for s, v := range *portfolio {
		symbol, volume := s, v
		go func(pc *chan position.Position, ec *chan error) {
			if res, err := fhSvc.EvaluateHolding(context.Background(), symbol, volume); err != nil {
				*ec <- err
			} else {
				*pc <- res
			}
		}(positionChan, errChan)
	}
}

func renderPortfolioEvaluation(portfolioSize int, positionChan *chan position.Position, errChan *chan error) {
	for i := 0; i < portfolioSize; i++ {
		select {
		case p := <-*positionChan:
			fmt.Println(p.ToString())
		case err := <-*errChan:
			fmt.Println(fmt.Errorf("Error: %w", err))
		}
	}
}
