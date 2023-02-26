package position

import (
	"fmt"
	"strings"
)

type Position interface {
	ToString() string
}

type position struct {
	NetProfit float32
	Symbol    string
}

func New(symbol string, netProfit float32) Position {
	return &position{
		NetProfit: netProfit,
		Symbol:    symbol,
	}
}

func (self *position) ToString() string {
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("%s ", self.Symbol))

	switch np := self.NetProfit; {
	case np > 0:
		sb.WriteString(fmt.Sprintf("Changed(1D)\t▲ Profit: %.2f$", self.NetProfit))
	case np < 0:
		sb.WriteString(fmt.Sprintf("Changed(1D)\t▼ Loss: %.2f$", self.NetProfit))
	default:
		sb.WriteString(fmt.Sprintf("Changed(1D)\t▬ No Change: %.2f$", self.NetProfit))
	}

	return sb.String()
}
