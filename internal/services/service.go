package services

import (
	"errors"
	"fmt"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/configer"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/logger"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"time"

	"github.com/vladimirbannikov/BIP-backend/internal/models"

	"github.com/gorilla/mux"
)

type Service struct {
	am AuthModelManager
	um UsersModelManager
	tm TestsModelManager
}

func NewService(am AuthModelManager, um UsersModelManager, tm TestsModelManager) Service {
	return Service{am: am, um: um, tm: tm}
}

const loginKey = "login"
const testIdKey = "id"

func (s Service) Launch() {
	implUsers := usersServer{s.um}
	implAuth := authServer{s.am}
	implTests := testsServer{s.tm}
	cfg, err := configer.GetConfig()
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: Launch: configer.GetConfig error: %s", err.Error()))
		return
	}
	muxx := http.NewServeMux()
	muxx.Handle("/", createRouter(implUsers, implAuth, implTests))
	serv := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Server.Port),
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           muxx,
	}
	if err = serv.ListenAndServe(); err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: Launch: serv.ListenAndServe error: %s", err.Error()))
	}
}

func createRouter(implUsers usersServer, implAuth authServer, implTests testsServer) *mux.Router {
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

	router.HandleFunc("/user2fa", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implAuth.User2FA(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc(fmt.Sprintf("/user/{%s:[a-zA-Z0-9_]+}", loginKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implUsers.GetUserProfile(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc(fmt.Sprintf("/userMy"), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implUsers.GetUserProfileOwn(w, req)
		case http.MethodPut:
			implUsers.UpdateUserProfile(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc(fmt.Sprintf("/tests"), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implTests.GetAllTests(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc(fmt.Sprintf("/tests/{%s:[0-9]+}", testIdKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implTests.GetTest(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc(fmt.Sprintf("/tests/{%s:[0-9]+}/commit", testIdKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implTests.GetScore(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc(fmt.Sprintf("/rating"), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implTests.GetRating(w, req)
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
		if req.URL.Path == "/register" {
			handler.ServeHTTP(writer, req)
		}
		logger.Log(logger.InfoPrefix, fmt.Sprintf("%q", rawRequest))
		handler.ServeHTTP(writer, req)
	})
}

func (s *authServer) authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		userPathReg, _ := regexp.Match(`/user/.*`, []byte(req.URL.Path))
		if req.URL.Path == "/register" || req.URL.Path == "/login" || req.URL.Path == "/user2fa" ||
			req.URL.Path == "/rating" || req.URL.Path == "/tests" || userPathReg {
			handler.ServeHTTP(writer, req)
			// после успешного логина клиенту нужно запросить 2fa
			// после регистрации клиенту выдается qr код
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

		err = s.m.CheckUser2FA(req.Context(), tokenStr)
		if err != nil {
			if errors.Is(err, models.ErrNo2FAOrExpired) {
				writer.WriteHeader(http.StatusUnauthorized)
				// клиент должен дергать логин
				return
			}
			logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: authServer: CheckUser2FA: %s", err))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(writer, req)
	})
}
