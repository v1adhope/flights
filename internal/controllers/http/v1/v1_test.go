package v1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "github.com/v1adhope/flights/internal/controllers/http/v1"
	"github.com/v1adhope/flights/internal/testhelpers"
	"github.com/v1adhope/flights/internal/usecases"
	"github.com/v1adhope/flights/internal/usecases/infrastructure/repository"
	"github.com/v1adhope/flights/pkg/logger"
	"github.com/v1adhope/flights/pkg/postgresql"
)

const (
	_pgMigrationsSourceUrl = "file://../../../../db/migrations"
	_loggerLevel           = "debug"
	_handlerMode           = gin.DebugMode
)

type Suite struct {
	suite.Suite
	ctx    context.Context
	pgC    *testhelpers.PostgresContainer
	router *gin.Engine
	utils  *testhelpers.Utils
}

func (s *Suite) SetupSuite() {
	t := s.T()

	s.ctx = context.Background()

	pgC, err := testhelpers.BuildContainer(s.ctx, _pgMigrationsSourceUrl)
	if err != nil {
		log.Fatalf("v1: v1_test: SetupSuite: BuildContainer: %v", err)
	}

	s.pgC = pgC

	if err := s.pgC.MigrateUp(); err != nil {
		log.Fatalf("v1: v1_test: SetupSuite: MigrateUp: %v", err)
	}

	pd, err := postgresql.Build(
		s.ctx,
		postgresql.WithConnStr(pgC.ConnStr),
	)
	if err != nil {
		log.Fatalf("v1: v1_test: SetupSuite: Build: %v", err)
	}
	t.Cleanup(func() {
		pd.Close()
	})

	repo := repository.New(pd)

	uc := usecases.New(repo)

	log := logger.New(
		logger.WithLevel(_loggerLevel),
	)

	v1.SetMode(_handlerMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	v1.Register(&v1.Router{
		Handler:  router,
		Usecases: uc,
		Log:      log,
	})

	s.router = router

	s.utils = testhelpers.NewUtils(pd)
}

func (s *Suite) TearDownSuite() {
	if err := s.pgC.Terminate(s.ctx); err != nil {
		log.Fatalf("v1: v1_test: TearDownSuite: Terminate: %v", err)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// INFO: tickets

type ticketCreateReq struct {
	id       string
	Provider string `json:"provider"`
	FlyFrom  string `json:"flyFrom"`
	FlyTo    string `json:"flyTo"`
	FlyAt    string `json:"flyAt"`
	ArriveAt string `json:"arriveAt"`
}

func (s *Suite) Test1CreateTicketPositive() {
	t := s.T()

	tcs := []struct {
		key  string
		body ticketCreateReq
	}{
		{
			key: "1",
			body: ticketCreateReq{
				Provider: "Emirates",
				FlyFrom:  "Moscow",
				FlyTo:    "Hanoi",
				FlyAt:    "2022-01-02T15:04:05+03:00",
				ArriveAt: "2022-01-03T15:04:05+07:00",
			},
		},
		{
			key: "The same as 1rd",
			body: ticketCreateReq{
				Provider: "Emirates",
				FlyFrom:  "Moscow",
				FlyTo:    "Hanoi",
				FlyAt:    "2022-01-02T15:04:05+03:00",
				ArriveAt: "2022-01-03T15:04:05+07:00",
			},
		},
		{
			key: "2",
			body: ticketCreateReq{
				Provider: "China Airlines",
				FlyFrom:  "Beijing",
				FlyTo:    "Moscow",
				FlyAt:    "2023-01-02T15:04:05+08:00",
				ArriveAt: "2023-01-03T15:04:05+03:00",
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPost,
				"/v1/tickets/",
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code, tc.key)
		}
	})
}

func (s *Suite) Test1CreateTicketNegative() {
	t := s.T()

	tcs := []struct {
		key  string
		body ticketCreateReq
	}{
		{
			key: "Required field miss",
			body: ticketCreateReq{
				Provider: "Emirates",
				FlyTo:    "Hanoi",
				FlyAt:    "2022-01-02T15:04:05+03:00",
				ArriveAt: "2022-01-03T15:04:05+07:00",
			},
		},
		// TODO: fix or remove
		// {
		// 	key: "Mixed up fly dates",
		// 	body: ticketCreateReq{
		// 		Provider: "China Airlines",
		// 		FlyFrom:  "Beijing",
		// 		FlyTo:    "Moscow",
		// 		FlyAt:    "2023-01-03T15:04:05+03:00",
		// 		ArriveAt: "2023-01-02T15:04:05+08:00",
		// 	},
		// },
		{
			key: "Spot overflow",
			body: ticketCreateReq{
				Provider: "China Airlines",
				FlyFrom:  strings.Repeat("A", 256),
				FlyTo:    "Moscow",
				FlyAt:    "2023-01-03T15:04:05+03:00",
				ArriveAt: "2023-01-02T15:04:05+08:00",
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPost,
				"/v1/tickets/",
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code, tc.key)
		}
	})
}

// TODO: error
// func (s *Suite) Test2ReplaceTicketPositive() {
// 	t := s.T()
//
// 	tcs := []struct {
// 		key  string
// 		body ticketCreateReq
// 	}{
// 		{
// 			key: "1",
// 			body: ticketCreateReq{
// 				id:       s.utils.GetTicketByOffset(s.ctx, 0),
// 				Provider: "China Airlines",
// 				FlyFrom:  "Beijing",
// 				FlyTo:    "Moscow",
// 				FlyAt:    "2023-01-02T15:04:05+08:00",
// 				ArriveAt: "2023-01-03T15:04:05+03:00",
// 			},
// 		},
// 	}
//
// 	t.Run("", func(t *testing.T) {
// 		for _, tc := range tcs {
// 			jsonData, err := json.Marshal(tc.body)
// 			assert.NoError(t, err, tc.key)
//
// 			req, err := http.NewRequest(
// 				"PUT",
// 				fmt.Sprintf("/v1/tickets/%s", tc.body.id),
// 				strings.NewReader(string(jsonData)),
// 			)
// 			assert.NoError(t, err, tc.key)
//
// 			w := httptest.NewRecorder()
//
// 			s.router.ServeHTTP(w, req)
//
// 			assert.Equal(t, http.StatusOK, w.Code, tc.key)
// 		}
// 	})
// }

func (s *Suite) Test3DeleteTicket() {
	t := s.T()

	tcs := []struct {
		key  string
		id   string
		code int
	}{
		{
			key:  "Success",
			id:   s.utils.GetTicketByOffset(s.ctx, 1),
			code: http.StatusOK,
		},
		{
			key:  "No content",
			id:   s.utils.GetTicketByOffset(s.ctx, 1),
			code: http.StatusNoContent,
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/v1/tickets/%s", tc.id),
				nil,
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, tc.code, w.Code, tc.key)
		}
	})
}

type ticket struct {
	Id        string `json:"id"`
	Provider  string `json:"provider"`
	FlyFrom   string `json:"flyFrom"`
	FlyTo     string `json:"flyTo"`
	FlyAt     string `json:"flyAt"`
	ArriveAt  string `json:"arriveAt"`
	CreatedAt string `json:"createdAt"`
}

func (s *Suite) Test4GetAllTicketsPositive() {
	t := s.T()

	tcs := []struct {
		key      string
		expected []ticket
	}{
		{
			key: "1",
			expected: []ticket{
				{
					Id:       s.utils.GetTicketByOffset(s.ctx, 0),
					Provider: "Emirates",
					FlyFrom:  "Moscow",
					FlyTo:    "Hanoi",
					FlyAt:    "2022-01-02T12:04:05Z",
					ArriveAt: "2022-01-03T08:04:05Z",
				},
				{
					Id:       s.utils.GetTicketByOffset(s.ctx, 1),
					Provider: "China Airlines",
					FlyFrom:  "Beijing",
					FlyTo:    "Moscow",
					FlyAt:    "2023-01-02T07:04:05Z",
					ArriveAt: "2023-01-03T12:04:05Z",
				},
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				http.MethodGet,
				"/v1/tickets/",
				nil,
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, tc.key)

			tickets := []ticket{}
			err = json.NewDecoder(w.Body).Decode(&tickets)
			assert.NoError(t, err, tc.key)
			for i := range tc.expected {
				flyAtRightFormat, err := time.Parse(time.RFC3339, tickets[i].FlyAt)
				assert.NoError(t, err, tc.key, i)
				arriveToRightFormat, err := time.Parse(time.RFC3339, tickets[i].ArriveAt)
				assert.NoError(t, err, tc.key, i)

				assert.Equal(t, tc.expected[i].Id, tickets[i].Id, tc.key, i)
				assert.Equal(t, tc.expected[i].Provider, tickets[i].Provider, tc.key, i)
				assert.Equal(t, tc.expected[i].FlyFrom, tickets[i].FlyFrom, tc.key, i)
				assert.Equal(t, tc.expected[i].FlyTo, tickets[i].FlyTo, tc.key, i)
				assert.Equal(t, tc.expected[i].FlyAt, flyAtRightFormat.UTC().Format(time.RFC3339), tc.key, i)
				assert.Equal(t, tc.expected[i].ArriveAt, arriveToRightFormat.UTC().Format(time.RFC3339), tc.key, i)
			}
		}
	})
}

// INFO: passangers

type passengerCreateReq struct {
	id         string
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

func (s *Suite) Test1CreatePassengerPositive() {
	t := s.T()

	tcs := []struct {
		key  string
		body passengerCreateReq
	}{
		{
			key: "1",
			body: passengerCreateReq{
				FirstName:  "Riley",
				LastName:   "Scott",
				MiddleName: "Reed",
			},
		},
		{
			key: "The same as 1rd",
			body: passengerCreateReq{
				FirstName:  "Riley",
				LastName:   "Scott",
				MiddleName: "Reed",
			},
		},

		{
			key: "2",
			body: passengerCreateReq{
				FirstName:  "Thomas",
				LastName:   "Langlois",
				MiddleName: "Floyd",
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPost,
				"/v1/passengers/",
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code, tc.key)
		}
	})
}

func (s *Suite) Test1CreatePassengerNegative() {
	t := s.T()

	tcs := []struct {
		key  string
		body passengerCreateReq
	}{
		{
			key: "Name overflow",
			body: passengerCreateReq{
				FirstName:  strings.Repeat("R", 256),
				LastName:   "Scott",
				MiddleName: "Reed",
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPost,
				"/v1/passengers/",
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code, tc.key)
		}
	})
}

// TODO: error
// func (s *Suite) Test2ReplacePassengerPositive() {
// 	t := s.T()
//
// 	tcs := []struct {
// 		key  string
// 		body passengerCreateReq
// 	}{
// 		{
// 			key: "1",
// 			body: passengerCreateReq{
// 				id:         s.utils.GetPassengerByOffset(s.ctx, 0),
// 				FirstName:  "",
// 				LastName:   "",
// 				MiddleName: "",
// 			},
// 		},
// 	}
//
// 	t.Run("", func(t *testing.T) {
// 		for _, tc := range tcs {
// 			jsonData, err := json.Marshal(tc.body)
// 			assert.NoError(t, err, tc.key)
//
// 			req, err := http.NewRequest(
// 				http.MethodPut,
// 				fmt.Sprintf("/v1/passengers/%s", tc.body.id),
// 				strings.NewReader(string(jsonData)),
// 			)
// 			assert.NoError(t, err, tc.key)
//
// 			w := httptest.NewRecorder()
//
// 			s.router.ServeHTTP(w, req)
//
// 			assert.Equal(t, http.StatusOK, w.Code, tc.key)
// 		}
// 	})
// }

func (s *Suite) Test3DeletePassenger() {
	t := s.T()

	tcs := []struct {
		key  string
		id   string
		code int
	}{
		{
			key:  "Success",
			id:   s.utils.GetPassengerByOffset(s.ctx, 1),
			code: http.StatusOK,
		},
		{
			key:  "No content",
			id:   s.utils.GetPassengerByOffset(s.ctx, 1),
			code: http.StatusNoContent,
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/v1/passengers/%s", tc.id),
				nil,
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, tc.code, w.Code, tc.key)
		}
	})
}

type passenger struct {
	Id         string `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

func (s *Suite) Test4GetAllPassengers() {
	t := s.T()

	tcs := []struct {
		key      string
		expected []passenger
	}{
		{
			key: "1",
			expected: []passenger{
				{
					Id:         s.utils.GetPassengerByOffset(s.ctx, 0),
					FirstName:  "Riley",
					LastName:   "Scott",
					MiddleName: "Reed",
				},
				{
					Id:         s.utils.GetPassengerByOffset(s.ctx, 1),
					FirstName:  "Thomas",
					LastName:   "Langlois",
					MiddleName: "Floyd",
				},
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				http.MethodGet,
				"/v1/passengers/",
				nil,
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			passengers := []passenger{}
			err = json.NewDecoder(w.Body).Decode(&passengers)
			assert.NoError(t, err, tc.key)

			assert.Equal(t, tc.expected, passengers)
		}
	})
}
