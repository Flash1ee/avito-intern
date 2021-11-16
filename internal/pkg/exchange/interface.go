package exchange

//go:generate mockgen -destination=mocks/exchange.go -package=mock_exchange -mock_names=Interface=Exchange . Interface

type Interface interface {
	// Exchange Errors:
	//		NotSupportedCurrency
	// 		app.GeneralError with Errors
	// 			InternalError
	Exchange(amount float64, currency string) (float64, error)
}
