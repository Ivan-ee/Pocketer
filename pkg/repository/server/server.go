package server

import (
	"github.com/zhashkevych/go-pocket-sdk"
	"net/http"
	"strconv"
	"telegram-bot/pkg/repository"
)

type AuthServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectUrl     string
}

func NewAuthServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectUrl string) *AuthServer {
	return &AuthServer{redirectUrl: redirectUrl, pocketClient: pocketClient, tokenRepository: tokenRepository}
}

func (s *AuthServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIdParam := r.URL.Query().Get("chat_id")
	if chatIdParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := strconv.ParseInt(chatIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.Get(repository.RequestToken, chatId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authResp, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.tokenRepository.Save(repository.AccessToken, authResp.AccessToken, chatId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", s.redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}
