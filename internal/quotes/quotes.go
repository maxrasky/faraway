package quotes

import (
	_ "embed"
	"encoding/json"
	"math/rand"
)

//go:embed quotes.json
var quotes []byte

type QuoteResponder struct {
	quotes []string
}

func New() *QuoteResponder {
	var q []string
	if err := json.Unmarshal(quotes, &q); err != nil {
		panic(err)
	}

	return &QuoteResponder{
		quotes: q,
	}
}

func (r *QuoteResponder) GetQuote() string {
	return r.quotes[rand.Intn(len(r.quotes))]
}
