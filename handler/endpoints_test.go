package handler

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Register(t *testing.T) {
	t.Parallel()

	type Case struct {
		name     string
		request  generated.RegistrationRequest
		mock     func(repo *repository.MockRepositoryInterface, input generated.RegistrationRequest)
		expected int
	}
	var testCases = []Case{
		{
			name: "request with all valid parameter",
			request: generated.RegistrationRequest{
				FullName: "test case",
				Password: "T3stv@lid",
				Phone:    "+6281234567890",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.RegistrationRequest) {
				repo.EXPECT().
					FindByPhone(gomock.Any(), repository.FindByPhoneInput{Phone: input.Phone}).
					Return(repository.FindByPhoneOutput{}, sql.ErrNoRows)
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(repository.RegistrationOutput{Id: 1}, nil)
			},
			expected: 200,
		},
		{
			name: "request with invalid password Numeric parameter",
			request: generated.RegistrationRequest{
				FullName: "test case",
				Password: "Testv@lid",
				Phone:    "+6281234567890",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.RegistrationRequest) {
			},
			expected: 400,
		},
		{
			name: "request with invalid password Special Char parameter",
			request: generated.RegistrationRequest{
				FullName: "test case",
				Password: "T3stvalid",
				Phone:    "+6281234567890",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.RegistrationRequest) {
			},
			expected: 400,
		},
		{
			name: "request with invalid phone number parameter",
			request: generated.RegistrationRequest{
				FullName: "test case",
				Password: "T3stv@lid",
				Phone:    "81234567890",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.RegistrationRequest) {
			},
			expected: 400,
		},
		{
			name: "request with invalid phone number parameter",
			request: generated.RegistrationRequest{
				FullName: "test case",
				Password: "T3stv@lid",
				Phone:    "+628123",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.RegistrationRequest) {
			},
			expected: 400,
		},
		{
			name: "request with existing phone number parameter",
			request: generated.RegistrationRequest{
				FullName: "test case",
				Password: "T3stv@lid",
				Phone:    "+6281234567890",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.RegistrationRequest) {
				repo.EXPECT().
					FindByPhone(gomock.Any(), gomock.Any()).
					Return(repository.FindByPhoneOutput{
						Id:       1,
						Slug:     base64.RawStdEncoding.EncodeToString([]byte(input.Phone)),
						FullName: input.FullName,
						Phone:    input.Phone,
						Password: input.Password,
					}, nil)
			},
			expected: 409,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	e := echo.New()

	repo := repository.NewMockRepositoryInterface(ctrl)
	s := NewServer(NewServerOptions{Repository: repo})

	for _, cases := range testCases {
		t.Run(cases.name, func(t *testing.T) {
			b, _ := json.Marshal(cases.request)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			cases.mock(repo, cases.request)

			assert.NoError(t, s.Register(ctx))
			assert.Equal(t, cases.expected, rec.Result().StatusCode)
		})
	}
}

func TestServer_Login(t *testing.T) {
	t.Parallel()

	type Case struct {
		name     string
		request  generated.LoginRequest
		mock     func(repo *repository.MockRepositoryInterface, input generated.LoginRequest)
		expected int
	}
	var testCases = []Case{
		{
			name: "request with all valid parameter",
			request: generated.LoginRequest{
				Password: "secret",
				Phone:    "+6282213770600",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.LoginRequest) {
				p, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
				repo.EXPECT().
					FindByPhone(gomock.Any(), gomock.Any()).
					Return(repository.FindByPhoneOutput{
						Id:       1,
						Slug:     "any",
						FullName: "any",
						Phone:    input.Phone,
						Password: string(p),
					}, nil)
			},
			expected: 200,
		},
		{
			name: "request with invalid credentials",
			request: generated.LoginRequest{
				Password: "secret",
				Phone:    "+6282213770600",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.LoginRequest) {
				repo.EXPECT().
					FindByPhone(gomock.Any(), gomock.Any()).
					Return(repository.FindByPhoneOutput{}, sql.ErrNoRows)
			},
			expected: 404,
		},
		{
			name: "request with invalid phone number format",
			request: generated.LoginRequest{
				Password: "secret",
				Phone:    "82213770600",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.LoginRequest) {
			},
			expected: 400,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	e := echo.New()

	repo := repository.NewMockRepositoryInterface(ctrl)
	s := NewServer(NewServerOptions{Repository: repo})
	for _, cases := range testCases {
		t.Run(cases.name, func(t *testing.T) {
			b, _ := json.Marshal(cases.request)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			cases.mock(repo, cases.request)

			assert.NoError(t, s.Login(ctx))
			assert.Equal(t, cases.expected, rec.Result().StatusCode)
		})
	}
}

func TestServer_Profile(t *testing.T) {
	t.Parallel()

	type Case struct {
		name     string
		request  generated.LoginRequest
		mock     func(repo *repository.MockRepositoryInterface, input generated.LoginRequest)
		useAuth  bool
		expected int
	}
	var testCases = []Case{
		{
			name: "request with all valid parameter",
			request: generated.LoginRequest{
				Password: "secret",
				Phone:    "+6282213770600",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.LoginRequest) {
				p, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
				repo.EXPECT().
					FindByPhone(gomock.Any(), gomock.Any()).
					Return(repository.FindByPhoneOutput{
						Id:       1,
						Slug:     "slug",
						FullName: "full name",
						Phone:    input.Phone,
						Password: string(p),
					}, nil)
				repo.EXPECT().
					FindBySlug(gomock.Any(), gomock.Any()).
					Return(repository.FindBySlugOutput{
						Slug:     "slug",
						FullName: "full name",
						Phone:    input.Phone,
						Password: string(p),
					}, nil)
			},
			useAuth:  true,
			expected: 200,
		},
		{
			name: "request without authentication",
			request: generated.LoginRequest{
				Password: "secret",
				Phone:    "+6282213770600",
			},
			mock: func(repo *repository.MockRepositoryInterface, input generated.LoginRequest) {

			},
			useAuth:  false,
			expected: 403,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	e := echo.New()

	repo := repository.NewMockRepositoryInterface(ctrl)
	s := NewServer(NewServerOptions{Repository: repo})
	for _, cases := range testCases {
		t.Run(cases.name, func(t *testing.T) {
			var (
				req   *http.Request
				rec   = httptest.NewRecorder()
				token string
			)
			if cases.useAuth {
				b, _ := json.Marshal(cases.request)
				req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				ctx := e.NewContext(req, rec)

				cases.mock(repo, cases.request)

				assert.NoError(t, s.Login(ctx))

				var response generated.LoginResponse
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)

				token = response.Token
			}

			profileReq := httptest.NewRequest(http.MethodGet, "/profile", nil)
			if cases.useAuth {
				profileReq.Header.Set("Authorization", "Bearer "+token)
			}

			profileRec := httptest.NewRecorder()
			ctx := e.NewContext(profileReq, profileRec)
			profile := Middleware()(s.Profile)

			assert.NoError(t, profile(ctx))
			assert.Equal(t, cases.expected, profileRec.Code)
		})
	}
}
