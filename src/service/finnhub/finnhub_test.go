package finnhub

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/elijah688/delio-mini/src/domain/position"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	INVALID_SYMBOL = "INVALID_SYMBOL"
)

type ErrInvalidSymbol struct{}

func (*ErrInvalidSymbol) Error() string {
	return INVALID_SYMBOL
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestFinnhubService_EvaluateHolding(t *testing.T) {

	// override the http.DefaultClient with the mock server client
	cfg := finnhub.NewConfiguration()
	mockClient := &MockHTTPClient{}

	// Use the mock client with an http.Client instance
	cfg.HTTPClient = &http.Client{
		Transport: mockClient,
	}

	cfg.Servers = finnhub.ServerConfigurations{
		{
			URL: "finnhub.mock.io",
		},
	}

	// create the finnhub client and call the API
	fhSvc := New(cfg)

	testCases := map[string]struct {
		symbol, url, quote string
		volume, status     int
		result             position.Position
		err                error
	}{
		"valid input": {
			symbol: "AAPL",
			volume: 10,
			status: http.StatusOK,
			url:    "finnhub.mock.io/quote?symbol=AAPL",
			quote:  `{"c": 2.00, "pc": 0.07}`,
			result: position.New("AAPL", 19.3),
			err:    nil,
		},
		"error response": {
			symbol: "ERROR",
			status: http.StatusInternalServerError,
			url:    "finnhub.mock.io/quote?symbol=ERROR",
			volume: 10,
			result: nil,
			err:    new(ErrInvalidSymbol),
		},
	}

	// assert that the response matches the mock response
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			url, err := url.Parse(tc.url)
			if err != nil {
				panic(err)
			}
			err, body, res := tc.err, tc.quote, &http.Request{
				Method:     "GET",
				URL:        url,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Accept": []string{"application/json"}, "User-Agent": []string{"OpenAPI-Generator/2.0.15/go"}},
			}
			resCtx := res.WithContext(context.Background())
			mockClient.On("RoundTrip", resCtx).Return(&http.Response{
				StatusCode: tc.status,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       ioutil.NopCloser(strings.NewReader(body)),
			}, err)
			result, err := fhSvc.EvaluateHolding(context.Background(), tc.symbol, tc.volume)
			if tc.err != nil {
				assert.Contains(t, err.Error(), tc.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.result, result)
		})
	}
}
