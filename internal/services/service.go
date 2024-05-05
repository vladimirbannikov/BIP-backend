package services

import (
	"fmt"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/configer"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/logger"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Service struct {
	am AuthModelManager
	um UsersModelManager
}

func NewService(am AuthModelManager, um UsersModelManager) Service {
	return Service{am: am, um: um}
}

func (s Service) Launch() {
	implUsers := usersServer{s.um}
	implAuth := authServer{s.am}
	cfg, err := configer.GetConfig()
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: Launch: configer.GetConfig error: %s", err.Error()))
		return
	}
	muxx := http.NewServeMux()
	muxx.Handle("/", createRouter(implUsers, implAuth))
	serv := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Server.Port),
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           muxx,
	}
	if err = serv.ListenAndServe(); err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: Launch: serv.ListenAndServe error: %s", err.Error()))
	}
}

func createRouter(implUsers usersServer, implAuth authServer) *mux.Router {
	router := mux.NewRouter()
	router.Use(logMiddleware)
	router.Use(implAuth.authMiddleware)

	router.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implAuth.Register(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implAuth.Login(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodDelete:
			implAuth.Logout(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	return router
}

const authHeaderStr = "Auth"

func logMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		rawRequest, _ := httputil.DumpRequest(req, true)
		logger.Log(logger.InfoPrefix, fmt.Sprintf("%q", rawRequest))
		handler.ServeHTTP(writer, req)
	})
}

func (s *authServer) authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/register" || req.URL.Path == "/login" {
			handler.ServeHTTP(writer, req)
			return
		}
		tokenStr := req.Header.Get(authHeaderStr)
		if tokenStr == "" {
			logger.Log(logger.ErrPrefix, "Service: authServer: authMiddleware: No auth header found")
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		err := s.m.ValidateUserToken(req.Context(), tokenStr)
		if err != nil {
			logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: authServer: ValidateUserToken: %s", err))
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(writer, req)
	})
}
