package exchange

import "github.com/pkg/errors"

var (
	InternalError        = errors.New("internal error")
	NotSupportedCurrency = errors.New("not supported currency for convert")
)

const cbrQuery = "http://www.cbr-xml-daily.ru/latest.js"

type Currency struct {
	Rates map[string]float64 `json:"rates"`
}

var validCurrency = []string{
	"AMD", "NOK", "TRY", "USD", "CAD", "CNY", "UAH", "CZK",
	"JPY", "GBP", "HUF", "MDL", "UZS", "AUD", "INR", "EUR",
	"KGS", "TMT", "ZAR", "BRL", "HKD", "KZT", "SEK", "CHF",
	"KRW", "BYN", "BGN", "PLN", "SGD", "AZN", "DKK", "RON",
	"XDR", "TJS"}
