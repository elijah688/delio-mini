package position

import (
	"testing"
)

func TestPositionToString(t *testing.T) {
	testCases := []struct {
		name         string
		netProfit    float32
		symbol       string
		expected     string
		expectedFail bool
	}{
		{
			name:      "Positive Net Profit",
			netProfit: 100.50,
			symbol:    "AAPL",
			expected:  "AAPL Changed(1D)\t▲ Profit: 100.50$",
		},
		{
			name:      "Negative Net Profit",
			netProfit: -50.25,
			symbol:    "AAPL",
			expected:  "AAPL Changed(1D)\t▼ Loss: -50.25$",
		},
		{
			name:      "Zero Net Profit",
			netProfit: 0,
			symbol:    "AAPL",
			expected:  "AAPL Changed(1D)\t▬ No Change: 0.00$",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a position
			p := New(tc.symbol, tc.netProfit)

			// Test the ToString method
			result := p.ToString()
			if result != tc.expected {
				t.Errorf("ToString() returned %s, expected %s", result, tc.expected)
			}
		})
	}
}
