package finnhub

import (
	"context"

	"github.com/elijah688/delio-mini/src/domain/position"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type FinnhubService interface {
	EvaluateHolding(ctx context.Context, symbol string, volume int) (position.Position, error)
}

type finnhubService struct {
	inner *finnhub.DefaultApiService
}

func New(cfg *finnhub.Configuration) FinnhubService {

	return &finnhubService{
		inner: finnhub.NewAPIClient(cfg).DefaultApi,
	}
}

func (self *finnhubService) EvaluateHolding(ctx context.Context, symbol string, volume int) (position.Position, error) {
	// Get Symbol Quote
	res, _, err := self.inner.Quote(ctx).Symbol(symbol).Execute()
	if err != nil {
		return nil, err
	}

	// Calculate Net Profit
	c, pc := res.GetC(), res.GetPc()
	netProfit := (c - pc) * float32(volume)
	return position.New(symbol, netProfit), nil
}
