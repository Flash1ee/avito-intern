package exchange

import (
	"avito-intern/internal/app"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type ExchangeService struct {
	Api string
}

func NewExchangeService(apiUrl string) *ExchangeService {
	return &ExchangeService{
		Api: apiUrl,
	}
}

// Exchange Errors:
//		NotSupportedCurrency
// 		app.GeneralError with Errors
// 			InternalError
func (ex *ExchangeService) Exchange(amount float64, currency string) (float64, error) {
	if amount == 0 {
		return amount, nil
	}

	if !InArray(validCurrency, currency) {
		return app.InvalidFloat, NotSupportedCurrency
	}

	req, err := http.NewRequest("GET", ex.Api, nil)
	if err != nil {
		return app.InvalidFloat, app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return app.InvalidFloat, app.GeneralError{
			Err:         InternalError,
			ExternalErr: errors.Wrapf(err, "error send Exchange request"),
		}
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var record Currency

	if err = json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return app.InvalidFloat, app.GeneralError{
			Err:         InternalError,
			ExternalErr: err,
		}
	}

	retValue := record.Rates[currency] * amount
	return retValue, nil
}
