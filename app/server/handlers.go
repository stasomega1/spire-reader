package server

import (
	"net/http"
	"spire-reader/app/model"
)

func (server *Server) Version() http.HandlerFunc {
	handlerName := "Version"
	return func(w http.ResponseWriter, r *http.Request) {
		server.Logger.Info("test")
		response, err := server.SpireService.Version()
		server.checkError(w, err, handlerName, r.Context().Value(CtxRequestIdKey).(string))
		server.RespondJson(w, http.StatusOK, &model.SimpleResponse{Code: model.SuccessCode, Message: response})
	}
}
