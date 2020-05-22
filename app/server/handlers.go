package server

import (
	"net/http"
	"spire-reader/app/model"
)

func (server *Server) Version() http.HandlerFunc {
	handlerName := "Version"
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := server.SpireService.Version()
		server.checkError(w, err, handlerName, r.Context().Value(CtxRequestIdKey).(string))
		server.RespondJson(w, http.StatusOK, &model.SimpleResponse{Code: model.SuccessCode, Message: response})
	}
}

func (server *Server) GetExampleRunData() http.HandlerFunc {
	handlerName := "GetExampleData"
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := server.SpireService.GetExampleRunData()
		server.checkError(w, err, handlerName, r.Context().Value(CtxRequestIdKey).(string))
		server.RespondJson(w, http.StatusOK, response)
	}
}

func (server *Server) TestFunc() http.HandlerFunc {
	handlerName := "TestFunc"
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := server.SpireService.TestFunc("./exampledata/jorbs")
		server.checkError(w, err, handlerName, r.Context().Value(CtxRequestIdKey).(string))
		server.RespondJson(w, http.StatusOK, response)
	}
}
