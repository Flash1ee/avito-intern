package handler

import (
	"avito-intern/internal/app"
	"avito-intern/internal/pkg/utilits"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Code  int
	Error error
	Level logrus.Level
}

type CodeMap map[error]ResponseError

type RequestBody interface{}

type HelpHandlers struct {
	utilits.Responder
}

func (h *HelpHandlers) UsecaseError(w http.ResponseWriter, r *http.Request, usecaseErr error, codeByErr CodeMap) {
	var generalError *app.GeneralError
	orginalError := usecaseErr
	if errors.As(usecaseErr, &generalError) {
		usecaseErr = errors.Cause(usecaseErr).(*app.GeneralError).Err
	}

	respond := ResponseError{http.StatusServiceUnavailable,
		app.UnknownError, logrus.ErrorLevel}

	for err, respondErr := range codeByErr {
		if errors.Is(usecaseErr, err) {
			respond = respondErr
			break
		}
	}

	h.Log(r).Logf(respond.Level, "Gotted error: %v", orginalError)
	h.Error(w, r, respond.Code, respond.Error)
}

// GetInt64FromParam HTTPErrors
//		Status 400 handler.InvalidParameters
func (h *HelpHandlers) GetInt64FromParam(w http.ResponseWriter, r *http.Request, name string) (int64, bool) {
	vars := mux.Vars(r)
	number, ok := vars[name]
	numberInt, err := strconv.ParseInt(number, 10, 64)
	if !ok || err != nil {
		h.Log(r).Infof("can'not get parametrs %s, was got %v)", name, vars)
		h.Error(w, r, http.StatusBadRequest, InvalidParameters)
		return app.InvalidInt, false
	}
	return numberInt, true
}

func (h *HelpHandlers) GetRequestBody(_ http.ResponseWriter, r *http.Request,
	reqStruct RequestBody) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(reqStruct); err != nil {
		return err
	}

	return nil
}
