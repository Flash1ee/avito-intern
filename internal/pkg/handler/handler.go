package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

type BaseHandler struct {
	router *mux.Router
	HelpHandlers
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
