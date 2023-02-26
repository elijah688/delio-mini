package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/elijah688/delio-mini/src/domain/position"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slices"
)

func TestEvaluatePortfolio(t *testing.T) {
	// prepare inputs
	positionChan, errChan, fhSvc := make(chan position.Position), make(chan error), &MockFinnhubService{}

	// set up expectations
	testCases := []struct {
		name           string
		portfolio      map[string]int
		expectedResult map[string]position.Position
		expectedError  map[string]error
	}{
		{
			name:      "success",
			portfolio: map[string]int{AAPL: 10, MSFT: 10},
			expectedResult: map[string]position.Position{
				AAPL: position.New(AAPL, 1234.5),
				MSFT: position.New(MSFT, 6789.0),
			},

			expectedError: nil,
		},
		{
			name:           "empty portfolio",
			portfolio:      map[string]int{},
			expectedResult: map[string]position.Position{},
			expectedError:  nil,
		},
		{
			name:           "internal server error",
			portfolio:      map[string]int{AAPL: 10, MSFT: 10},
			expectedResult: nil,
			expectedError:  map[string]error{AAPL: errors.New("invalid symbol: INVALID"), MSFT: errors.New("invalid symbol: INVALID")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// evaluate portfolio
			for symbol, volume := range tc.portfolio {
				fhSvc.On("EvaluateHolding", context.Background(), symbol, volume).Return(tc.expectedResult[symbol], tc.expectedError[symbol]).Times(1)
			}

			evaluatePortfolio(fhSvc, &tc.portfolio, &positionChan, &errChan)
			actual, errs := make([]position.Position, 0, len(tc.portfolio)), 0
			for range tc.portfolio {
				select {
				case p := <-positionChan:
					actual = append(actual, p)
				case <-errChan:
					errs++
				}

			}

			// assert results
			assert.Equal(t, len(actual), len(tc.expectedResult))
			for _, p := range tc.expectedResult {
				assert.True(t, slices.Contains(actual, p))
			}

			// assert.Equal(t, actual, tc.expectedResult)
			assert.Equal(t, len(tc.expectedError), errs)

		})
	}
}

func TestRenderPortfolioEvaluation(t *testing.T) {
	// create a pipe for capturing stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// replace stdout with the write end of the pipe
	tmp := os.Stdout
	os.Stdout = w

	// create some test data
	positionChan, errChan := make(chan position.Position), make(chan error)
	go func() {
		positionChan <- position.New(AAPL, 100.0)
		positionChan <- position.New(MSFT, 200.0)
		errChan <- fmt.Errorf("something went wrong")
	}()
	portfolioSize := 3

	// call the function being tested
	renderPortfolioEvaluation(portfolioSize, &positionChan, &errChan)

	// restore stdout
	w.Close()
	os.Stdout = tmp

	// read the captured output from the read end of the pipe
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// assert against the captured output
	expectedOutput := "AAPL Changed(1D)	▲ Profit: 100.00$\nMSFT Changed(1D)	▲ Profit: 200.00$\nError: something went wrong\n"
	if buf.String() != expectedOutput {
		t.Errorf("Unexpected output. Expected: %s, got: %s", expectedOutput, buf.String())
	}
}

type MockFinnhubService struct {
	mock.Mock
}

func (self *MockFinnhubService) EvaluateHolding(ctx context.Context, symbol string, volume int) (position.Position, error) {
	args := self.Called(ctx, symbol, volume)
	p, err := args.Get(0), args.Error(1)
	if err != nil {
		return nil, err
	}
	return p.(position.Position), nil

}
