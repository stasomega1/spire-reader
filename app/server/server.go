package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"spire-reader/app/model"
	"spire-reader/app/server/servererrors"
	"spire-reader/app/services"
)

type Server struct {
	Router       *mux.Router
	Logger       *logrus.Logger
	SpireService *services.SpireService
	//JWT
	Auth                bool
	JwtKey              []byte
	JwtTokenLiveMinutes int
}

func NewServer(service *services.SpireService, logger *logrus.Logger, jwtKey []byte, jwtTokenLiveMinutes int) *Server {
	server := &Server{
		Router:              mux.NewRouter(),
		Logger:              logger,
		SpireService:        service,
		Auth:                false,
		JwtKey:              jwtKey,
		JwtTokenLiveMinutes: jwtTokenLiveMinutes,
	}
	server.configureRouterMiddleware()
	server.configureRouter()
	return server
}

func (server *Server) Start(port int) error {
	if server.Auth && (server.JwtKey == nil || server.JwtTokenLiveMinutes == 0) {
		return servererrors.New(servererrors.ARGISNIL, "For authenticateUser middleware JwtKey and JwtTokenLiveMinutes must not be nil", "")
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), server)
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.Router.ServeHTTP(w, r)
}

//Endpoints
func (server *Server) configureRouter() {
	server.Router.HandleFunc("/version", server.Version()).Methods(http.MethodGet, http.MethodOptions)
	server.Router.HandleFunc("/test", server.TestFunc()).Methods(http.MethodGet, http.MethodOptions)
}

func (server *Server) logRestHandler(level model.LogLevel, handlerName, info string) {
	switch level {
	case model.LogLevelDebug:
		server.Logger.Debugf("HandlerName: %s, information: %s", handlerName, info)
	case model.LogLevelErr:
		server.Logger.Errorf("HandlerName: %s, information: %s", handlerName, info)
	case model.LogLevelTrace:
		server.Logger.Tracef("HandlerName: %s, information: %s", handlerName, info)
	case model.LogLevelInfo:
		server.Logger.Infof("HandlerName: %s, information: %s", handlerName, info)
	case model.LogLevelWarn:
		server.Logger.Warnf("HandlerName: %s, information: %s", handlerName, info)
	}
}

func (server *Server) checkError(w http.ResponseWriter, err error, serviceName string, requestId string) {
	if err == nil {
		return
	}
	switch err := errors.Cause(err).(type) {
	case *servererrors.ErrorStruct:
		if err.Code == servererrors.INVALID {
			server.log(model.LogLevelWarn, serviceName, requestId, err.Message)
			server.RespondJson(w, http.StatusBadRequest, &model.SimpleResponse{Code: model.ErrCode, Message: err.MessageRu})
		} else if err.Code == servererrors.EMPTY {
			server.log(model.LogLevelWarn, serviceName, requestId, err.Message)
			server.RespondJson(w, http.StatusBadRequest, &model.SimpleResponse{Code: model.EmptyCode, Message: err.MessageRu})
		} else if err.Code == servererrors.INTERNAL {
			server.log(model.LogLevelErr, serviceName, requestId, err.Message)
			server.RespondJson(w, http.StatusInternalServerError, &model.SimpleResponse{Code: model.ErrCode, Message: err.MessageRu})
		} else if err.Code == servererrors.ARGISNIL {
			server.log(model.LogLevelErr, serviceName, requestId, err.Message)
			server.RespondJson(w, http.StatusBadRequest, &model.SimpleResponse{Code: model.ErrCode, Message: err.MessageRu})
		} else {
			server.log(model.LogLevelErr, serviceName, requestId, err.Message)
			server.RespondJson(w, http.StatusInternalServerError, &model.SimpleResponse{Code: model.ErrCode, Message: err.MessageRu})
		}
	default:
		server.log(model.LogLevelErr, serviceName, requestId, fmt.Sprintf("Unknown err: %v", err))
		server.Error(w, http.StatusInternalServerError, err)
	}
}

func (server *Server) log(level model.LogLevel, serviceName, requestId, info string) {
	switch level {
	case model.LogLevelDebug:
		server.Logger.Debugf("ServiceName: %s, requestId: %s, information: %s", serviceName, requestId, info)
	case model.LogLevelErr:
		server.Logger.Errorf("ServiceName: %s, requestId: %s, information: %s", serviceName, requestId, info)
	case model.LogLevelTrace:
		server.Logger.Tracef("ServiceName: %s, requestId: %s, information: %s", serviceName, requestId, info)
	case model.LogLevelInfo:
		server.Logger.Infof("ServiceName: %s, requestId: %s, information: %s", serviceName, requestId, info)
	case model.LogLevelWarn:
		server.Logger.Warnf("ServiceName: %s, requestId: %s, information: %s", serviceName, requestId, info)
	}
}

func (server *Server) configureRouterMiddleware() {
	server.Router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodOptions}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Content-Disposition"}),
		handlers.ExposedHeaders([]string{"Authorization", "Set-Cookie", "Content-Disposition"})))
	//Middleware
	server.Router.Use(server.setRequestId)
	server.Router.Use(server.loggingRequests)
	server.Router.Use(server.authenticateUser)
}

func (server *Server) RespondJson(w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}

func (server *Server) Error(w http.ResponseWriter, status int, err error) {
	server.RespondJson(w, status, &model.SimpleResponse{Code: model.ErrCode, Message: err.Error()})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rs *responseWriter) WriteHeader(statusCode int) {
	rs.statusCode = statusCode
	rs.ResponseWriter.WriteHeader(statusCode)
}
