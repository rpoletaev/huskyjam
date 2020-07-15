package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/internal"
	"github.com/rpoletaev/huskyjam/mock"
	"github.com/rpoletaev/huskyjam/pkg/auth"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/gavv/httpexpect.v2"
)

// var svr *httptest.Server
// var api *Api

// func TestMain(m *testing.M) {
// 	api = &Api{
// 		Config: &Config{
// 			Port: ":3000",
// 		},
// 		Accounts: &AccountHandler{},
// 		Logger:   log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
// 	}

// 	svr = httptest.NewServer(api.Server().Handler)
// 	exitStatus := m.Run()
// 	os.Exit(exitStatus)
// }

func TestSucceessfullSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	passHelper := mock.NewMockPassHashHelper(ctrl)
	passHelper.EXPECT().Hash(gomock.Any()).Return("", nil)

	store := mock.NewMockAccountsRepository(ctrl)
	store.EXPECT().Create(gomock.Any()).Return(nil)

	api := &Api{
		Config: &Config{
			Port: ":3000",
		},
		Accounts: &AccountHandler{
			Store:          store,
			PassHashHelper: passHelper,
			Log:            log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
		},
		Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}

	svr := httptest.NewServer(api.Server().Handler)
	e := httpexpect.New(t, svr.URL)
	req := &signupRequest{
		Email: "test@email.zone",
		Pass:  "supperpass",
	}

	e.POST("/auth/signup").
		WithJSON(req).
		Expect().
		Status(http.StatusCreated)
}

func TestSignupWithEmptyEmail(t *testing.T) {

	api := &Api{
		Config: &Config{
			Port: ":3000",
		},
		Accounts: &AccountHandler{
			Log: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
		},
		Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}

	svr := httptest.NewServer(api.Server().Handler)
	e := httpexpect.New(t, svr.URL)
	req := &signupRequest{
		Email: "",
		Pass:  "supperpass",
	}

	e.POST("/auth/signup").
		WithJSON(req).
		Expect().
		Status(http.StatusBadRequest).
		Body().Contains("login must be greater than 0: validation error")
}

func TestSignupWithEmptyPassword(t *testing.T) {

	api := &Api{
		Config: &Config{
			Port: ":3000",
		},
		Accounts: &AccountHandler{
			Log: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
		},
		Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}

	svr := httptest.NewServer(api.Server().Handler)
	e := httpexpect.New(t, svr.URL)
	req := &signupRequest{
		Email: "test@email.zone",
		Pass:  "",
	}

	e.POST("/auth/signup").
		WithJSON(req).
		Expect().
		Status(http.StatusBadRequest).
		Body().Contains("password must be greater than 0: validation error")
}

func TestSuccessfullSignin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	acc := &internal.Account{
		ID:    1,
		Email: "test@email.zone",
		Pass:  "supperpasshash",
	}
	store := mock.NewMockAccountsRepository(ctrl)
	store.EXPECT().GetByEmail(gomock.Any()).Return(acc, nil)

	passHelper := mock.NewMockPassHashHelper(ctrl)
	passHelper.EXPECT().Check("supperpass", "supperpasshash").Return(nil)

	claims := &auth.SystemClaims{
		ID:    acc.ID,
		Email: acc.Email,
	}
	refreshRepo := mock.NewMockTokensRepository(ctrl)
	refreshRepo.EXPECT().New(claims).Return("refresh-token", nil)

	jwtHelper := mock.NewMockTokens(ctrl)
	jwtHelper.EXPECT().SignToken(claims).Return("access-token", nil)

	api := &Api{
		Config: &Config{
			Port: ":3000",
		},
		Accounts: &AccountHandler{
			Store:          store,
			PassHashHelper: passHelper,
			RefreshRepo:    refreshRepo,
			Auth:           jwtHelper,
			Log:            log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
		},
		Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}

	svr := httptest.NewServer(api.Server().Handler)
	e := httpexpect.New(t, svr.URL)
	req := &signupRequest{
		Email: "test@email.zone",
		Pass:  "supperpass",
	}

	obj := e.POST("/auth/signin").
		WithJSON(req).
		Expect().
		Status(http.StatusOK).JSON().Object()

	obj.Value("access").String().Equal("access-token")
	obj.Value("refresh").String().Equal("refresh-token")
}

func TestSigninWithNonExistingUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockAccountsRepository(ctrl)
	store.EXPECT().GetByEmail("test@email.zone").Return(nil, internal.ErrNotFound)

	api := &Api{
		Config: &Config{
			Port: ":3000",
		},
		Accounts: &AccountHandler{
			Store: store,
			Log:   log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
		},
		Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}

	svr := httptest.NewServer(api.Server().Handler)
	e := httpexpect.New(t, svr.URL)
	req := &signupRequest{
		Email: "test@email.zone",
		Pass:  "supperpass",
	}

	e.POST("/auth/signin").
		WithJSON(req).
		Expect().
		Status(http.StatusNotFound).Body().Contains("not found")
}

func TestSigninWithWrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	acc := &internal.Account{
		ID:    1,
		Email: "test@email.zone",
		Pass:  "supperpasshash",
	}

	store := mock.NewMockAccountsRepository(ctrl)
	store.EXPECT().GetByEmail("test@email.zone").Return(acc, nil)

	passHelper := mock.NewMockPassHashHelper(ctrl)
	passHelper.EXPECT().Check("supperpass", "supperpasshash").Return(errors.New("password not match"))

	api := &Api{
		Config: &Config{
			Port: ":3000",
		},
		Accounts: &AccountHandler{
			Store:          store,
			PassHashHelper: passHelper,
			Log:            log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
		},
		Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}

	svr := httptest.NewServer(api.Server().Handler)
	e := httpexpect.New(t, svr.URL)
	req := &signupRequest{
		Email: "test@email.zone",
		Pass:  "supperpass",
	}

	e.POST("/auth/signin").
		WithJSON(req).
		Expect().
		Status(http.StatusUnauthorized).Body().Contains("wrong password")
}
