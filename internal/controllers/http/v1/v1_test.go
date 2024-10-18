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

type ticketCreateReq struct {
	id       string
	Provider string `json:"provider"`
	FlyFrom  string `json:"flyFrom"`
	FlyTo    string `json:"flyTo"`
	FlyAt    string `json:"flyAt"`
	ArriveAt string `json:"arriveAt"`
}

func (s *Suite) TestCreateTicketPositive() {
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
				"POST",
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

func (s *Suite) TestCreateTicketNegative() {
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
		{
			key: "Mixed up fly dates",
			body: ticketCreateReq{
				Provider: "China Airlines",
				FlyFrom:  "Beijing",
				FlyTo:    "Moscow",
				FlyAt:    "2023-01-03T15:04:05+03:00",
				ArriveAt: "2023-01-02T15:04:05+08:00",
			},
		},
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
				"POST",
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
// func (s *Suite) TestUpdateTicketPositive() {
// 	t := s.T()
//
// 	tcs := []struct {
// 		key  string
// 		body ticketCreateReq
// 	}{
// 		{
// 			key: "1",
// 			body: ticketCreateReq{
// 				id:       s.utils.GetFirstTicketID(s.ctx),
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

func (s *Suite) TestDeleteTicket() {
	t := s.T()

	tcs := []struct {
		key  string
		id   string
		code int
	}{
		{
			key:  "Success",
			id:   s.utils.GetTicketByOffset(s.ctx, 0),
			code: http.StatusOK,
		},
		{
			key:  "No content",
			id:   s.utils.GetTicketByOffset(s.ctx, 0),
			code: http.StatusNoContent,
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				"DELETE",
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
