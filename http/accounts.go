package http

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/internal"
	"github.com/rpoletaev/huskyjam/pkg/auth"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// PassHashHelper generates and checks hash from pass
type PassHashHelper interface {
	Hash(pass string) (string, error)
	Check(pass, hash string) error
}

type accountHandler struct {
	Store internal.AccountsRepository
	PassHashHelper
	RefreshRepo internal.TokensRepository
	Auth        auth.Tokens
	Log         zerolog.Logger
}

func (h *accountHandler) logger(ctx context.Context) *zerolog.Logger {
	id, ok := hlog.IDFromCtx(ctx)
	var l zerolog.Logger
	if ok {
		l = h.Log.With().Str("req_id", id.String()).Logger()
		return &l
	}

	l = h.Log.With().Logger()
	return &l
}

type signupRequest struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

func (r *signupRequest) Validate() error {
	if len(r.Email) == 0 {
		return errors.Wrap(errValidateError, "login must be greater than 0")
	}

	if len(r.Pass) == 0 {
		return errors.Wrap(errValidateError, "password must be greater than 0")
	}
	return nil
}

func (h *accountHandler) Signup(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())

	req := &signupRequest{}
	if err := unmarshal(r, req); err != nil {
		logger.Error().Err(err).Msg("on get signup request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := h.PassHashHelper.Hash(req.Pass)
	if err != nil {
		logger.Error().Err(err).Msg("on hash password")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	acc := &internal.Account{
		Email: req.Email,
		Pass:  hash,
	}

	if err := h.Store.Create(acc); err != nil {
		logger.Error().Err(err).Msg("on create new account")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type signinRequest struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

func (r *signinRequest) Validate() error {
	if r == nil {
		return errors.New("recuest shouldn't be empty")
	}
	return nil
}

type signinResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (h *accountHandler) Signin(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())

	req := &signinRequest{}
	if err := unmarshal(r, req); err != nil {
		logger.Error().Err(err).Msg("on get signin request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := h.Store.GetByEmail(req.Email)
	if err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("on get stored account")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}

	claims := &auth.SystemClaims{
		ID:    acc.ID,
		Email: acc.Email,
	}

	h.signin(logger, w, claims)
}

func (h *accountHandler) signin(logger *zerolog.Logger, w http.ResponseWriter, claims *auth.SystemClaims) {
	token, err := h.RefreshRepo.New(claims)
	if err != nil {
		logger.Error().Err(err).Str("email", claims.Email).Msg("on create refresh token")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	jwt, err := h.Auth.SignToken(claims)
	if err != nil {
		logger.Error().Err(err).Str("email", claims.Email).Msg("on sign token")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := &signinResponse{
		Access:  jwt,
		Refresh: token,
	}

	if err := writeJSON(w, resp); err != nil {
		logger.Error().Err(err).Interface("response", resp).Msg("on sign token")
		return
	}
}

type refreshRequest struct {
	Token string `json:"token"`
}

func (r *refreshRequest) Validate() error {
	if r == nil || len(r.Token) == 0 {
		return errors.New("token sholdn't be empty")
	}
	return nil
}

func (h *accountHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	req := &refreshRequest{}
	if err := unmarshal(r, req); err != nil {
		logger.Error().Err(err).Msg("on unmarshal request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := h.RefreshRepo.Get(req.Token)
	if err != nil {
		logger.Error().Err(err).Str("token", req.Token).Msg("on getting saved refresh token")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}

	if err := h.RefreshRepo.Delete(req.Token); err != nil {
		logger.Warn().Err(err).Msg("on delete token")
	}

	h.signin(logger, w, claims)
}
