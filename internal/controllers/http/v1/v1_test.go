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

// INFO: general

type id struct {
	Id string `json:"id"`
}

func convertTime(target string) (string, error) {
	timeT, err := time.Parse(time.RFC3339, target)
	if err != nil {
		return "", err
	}

	return timeT.UTC().Format(time.RFC3339), nil
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

func (s *Suite) Test1aCreateTicketPositive() {
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

func (s *Suite) Test1bCreateTicketNegative() {
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
// func (s *Suite) Test1cReplaceTicketPositive() {
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
// 				http.MethodPut,
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

func (s *Suite) Test1dDeleteTicket() {
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
		{
			key:  "Wrong id",
			id:   "01ef8e28-1800-6569-ac7c",
			code: http.StatusUnprocessableEntity,
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
	Id       string `json:"id"`
	Provider string `json:"provider"`
	FlyFrom  string `json:"flyFrom"`
	FlyTo    string `json:"flyTo"`
}

type timeFieldsTicket struct {
	FlyAt     string `json:"flyAt"`
	ArriveAt  string `json:"arriveAt"`
	CreatedAt string `json:"createdAt"`
}

func (s *Suite) Test1dGetTicketsPositive() {
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
				},
				{
					Id:       s.utils.GetTicketByOffset(s.ctx, 1),
					Provider: "China Airlines",
					FlyFrom:  "Beijing",
					FlyTo:    "Moscow",
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
			err = json.NewDecoder(w.Result().Body).Decode(&tickets)
			assert.NoError(t, err, tc.key)

			assert.Equal(t, tc.expected, tickets, tc.key)

			timeFields := []timeFieldsTicket{}
			err = json.NewDecoder(w.Body).Decode(&timeFields)
			assert.NoError(t, err, tc.key)

			for _, tf := range timeFields {
				_, err = convertTime(tf.FlyAt)
				assert.NoError(t, err, tc.key)
				_, err = convertTime(tf.ArriveAt)
				assert.NoError(t, err, tc.key)
				_, err = convertTime(tf.CreatedAt)
				assert.NoError(t, err, tc.key)
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

func (s *Suite) Test1eCreatePassengerPositive() {
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

			resp := id{}
			err = json.NewDecoder(w.Body).Decode(&resp)
			assert.NoError(t, err, tc.key)
		}
	})
}

func (s *Suite) Test1fCreatePassengerNegative() {
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

func (s *Suite) Test1gReplacePassengerPositive() {
	t := s.T()

	tcs := []struct {
		key  string
		body passengerCreateReq
	}{
		{
			key: "1",
			body: passengerCreateReq{
				id:         s.utils.GetPassengerByOffset(s.ctx, 0),
				FirstName:  "Randy",
				LastName:   "Shelby",
				MiddleName: "McGrath",
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/v1/passengers/%s", tc.body.id),
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, tc.key)
		}
	})
}

func (s *Suite) Test1hDeletePassenger() {
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
		{
			key:  "Wrong id",
			id:   "01ef8e28-1800-6569-ac7c",
			code: http.StatusUnprocessableEntity,
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

func (s *Suite) Test1iGetPassengers() {
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
					FirstName:  "Randy",
					LastName:   "Shelby",
					MiddleName: "McGrath",
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

			assert.Equal(t, tc.expected, passengers, tc.key)
		}
	})
}

// INFO: documents

type documentCreateReq struct {
	Id          string
	Type        string `json:"type"`
	Number      string `json:"number"`
	PassengerId string `json:"passengerId"`
}

func (s *Suite) Test1jCreateDocumentPositive() {
	t := s.T()

	tcs := []struct {
		key  string
		body documentCreateReq
	}{
		{
			key: "1",
			body: documentCreateReq{
				Type:        "Passport",
				Number:      "5555666777",
				PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
			},
		},
		{
			key: "2",
			body: documentCreateReq{
				Type:        "Id card",
				Number:      "5555666777",
				PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
			},
		},
		{
			key: "1",
			body: documentCreateReq{
				Type:        "International passport",
				Number:      "3333888000",
				PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPost,
				"/v1/documents/",
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code, tc.key)

			resp := id{}
			err = json.NewDecoder(w.Body).Decode(&resp)
			assert.NoError(t, err, tc.key)
		}
	})
}

func (s *Suite) Test1kCreateDocumentNegative() {
	t := s.T()

	tcs := []struct {
		key  string
		body documentCreateReq
		code int
	}{
		{
			key: "Has already exists",
			body: documentCreateReq{
				Type:        "Passport",
				Number:      "5555666777",
				PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
			},
			code: http.StatusConflict,
		},
		{
			key: "Not allowed document type",
			body: documentCreateReq{
				Type:        "Birthday cert",
				Number:      "5555666777",
				PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			key: "Number overflow",
			body: documentCreateReq{
				Type:        "International passport",
				Number:      strings.Repeat("1", 256),
				PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
			},
			code: http.StatusUnprocessableEntity,
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPost,
				"/v1/documents/",
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, tc.code, w.Code, tc.key)
		}
	})
}

func (s *Suite) Test1lReplaceDocument() {
	t := s.T()

	tcs := []struct {
		key  string
		body documentCreateReq
	}{
		{
			key: "Has already exists",
			body: documentCreateReq{
				Id:          s.utils.GetDocumentByOffset(s.ctx, 0),
				Type:        "International passport",
				Number:      "5555666888",
				PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/v1/documents/%s", tc.body.Id),
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, tc.key)
		}
	})

}

func (s *Suite) Test1mDeleteDocument() {
	t := s.T()

	tcs := []struct {
		key  string
		id   string
		code int
	}{
		{
			key:  "Success",
			id:   s.utils.GetDocumentByOffset(s.ctx, 1),
			code: http.StatusOK,
		},
		{
			key:  "No content",
			id:   s.utils.GetDocumentByOffset(s.ctx, 1),
			code: http.StatusNoContent,
		},
		{
			key:  "Wrong id",
			id:   "01ef8e28-1800-6569-ac7c",
			code: http.StatusUnprocessableEntity,
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/v1/documents/%s", tc.id),
				nil,
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, tc.code, w.Code, tc.key)
		}
	})
}

type document struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Number      string `json:"number"`
	PassengerId string `json:"passengerId"`
}

func (s *Suite) Test1nGetDocumentsByPassengerId() {
	t := s.T()

	tcs := []struct {
		key      string
		id       string
		expected []document
	}{
		{
			key: "1",
			id:  s.utils.GetPassengerByOffset(s.ctx, 0),
			expected: []document{
				{
					Id:          s.utils.GetDocumentByOffset(s.ctx, 0),
					Type:        "Id card",
					Number:      "5555666777",
					PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
				},
				{
					Id:          s.utils.GetDocumentByOffset(s.ctx, 1),
					Type:        "International passport",
					Number:      "5555666888",
					PassengerId: s.utils.GetPassengerByOffset(s.ctx, 0),
				},
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/v1/documents/by-passenger/%s", tc.id),
				nil,
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, tc.key)

			documents := []document{}
			err = json.NewDecoder(w.Body).Decode(&documents)
			assert.NoError(t, err, tc.key)

			assert.Equal(t, tc.expected, documents, tc.key)
		}
	})
}

// INFO: passanger,  ticket

type passengerBoundingTicketReq struct {
	Id       string `json:"id"`
	TicketId string `json:"ticketId"`
}

func (s *Suite) Test1oBoundPassengerToTicket() {
	t := s.T()

	tcs := []struct {
		key  string
		body passengerBoundingTicketReq
		code int
	}{
		{
			key: "1",
			body: passengerBoundingTicketReq{
				Id:       s.utils.GetPassengerByOffset(s.ctx, 0),
				TicketId: s.utils.GetTicketByOffset(s.ctx, 0),
			},
			code: http.StatusCreated,
		},
		{
			key: "1",
			body: passengerBoundingTicketReq{
				Id:       s.utils.GetPassengerByOffset(s.ctx, 1),
				TicketId: s.utils.GetTicketByOffset(s.ctx, 0),
			},
			code: http.StatusCreated,
		},
		{
			key: "Has already bounded",
			body: passengerBoundingTicketReq{
				Id:       s.utils.GetPassengerByOffset(s.ctx, 0),
				TicketId: s.utils.GetTicketByOffset(s.ctx, 0),
			},
			code: http.StatusConflict,
		},
		{
			key: "Wrong id",
			body: passengerBoundingTicketReq{
				Id:       "01ef8e24-55c9-6316-ac7c",
				TicketId: s.utils.GetTicketByOffset(s.ctx, 0),
			},
			code: http.StatusUnprocessableEntity,
		},
		{
			key: "Passenger by id not found",
			body: passengerBoundingTicketReq{
				Id:       "01ef8e24-55c9-6316-ac7c-0242ac120003",
				TicketId: s.utils.GetTicketByOffset(s.ctx, 0),
			},
			code: http.StatusConflict,
		},
		{
			key: "Ticket by id not found",
			body: passengerBoundingTicketReq{
				Id:       s.utils.GetPassengerByOffset(s.ctx, 0),
				TicketId: "01ef8e24-55c9-6316-ac7c-0242ac120003",
			},
			code: http.StatusConflict,
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPost,
				"/v1/passengers/bound-to-ticket/",
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, tc.code, w.Code, tc.key)
		}
	})
}

func (s *Suite) Test1qForbidenDeleteTicketIfPassengersOnBoarding() {
	t := s.T()

	tcs := []struct {
		key  string
		id   string
		code int
	}{
		{
			key:  "1",
			id:   s.utils.GetTicketByOffset(s.ctx, 0),
			code: http.StatusForbidden,
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

func (s *Suite) Test1rUnboundPassengerToTicket() {
	t := s.T()

	tcs := []struct {
		key  string
		body passengerBoundingTicketReq
		code int
	}{
		{
			key: "1",
			body: passengerBoundingTicketReq{
				Id:       s.utils.GetPassengerByOffset(s.ctx, 1),
				TicketId: s.utils.GetTicketByOffset(s.ctx, 0),
			},
			code: http.StatusOK,
		},
		{
			key: "Nothing to unbound",
			body: passengerBoundingTicketReq{
				Id:       s.utils.GetPassengerByOffset(s.ctx, 1),
				TicketId: s.utils.GetTicketByOffset(s.ctx, 0),
			},
			code: http.StatusNoContent,
		},
		{
			key: "Wrong id",
			body: passengerBoundingTicketReq{
				Id:       "01ef8e24-55c9-6316-ac7c",
				TicketId: s.utils.GetTicketByOffset(s.ctx, 0),
			},
			code: http.StatusUnprocessableEntity,
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			jsonData, err := json.Marshal(tc.body)
			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				http.MethodPost,
				"/v1/passengers/unbound-from-ticket/",
				strings.NewReader(string(jsonData)),
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, tc.code, w.Code, tc.key)
		}
	})
}

func (s *Suite) Test1sGetPassengersByTicketIdPositive() {
	t := s.T()

	tcs := []struct {
		key      string
		id       string
		expected []passenger
	}{
		{
			key: "1",
			id:  s.utils.GetTicketByOffset(s.ctx, 0),
			expected: []passenger{
				{
					Id:         s.utils.GetPassengerByOffset(s.ctx, 0),
					FirstName:  "Riley",
					LastName:   "Scott",
					MiddleName: "Reed",
				},
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/v1/passengers/by-ticket-id/%s", tc.id),
				nil,
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, tc.key)

			passengers := []passenger{}
			err = json.NewDecoder(w.Body).Decode(&passengers)
			assert.NoError(t, err, tc.key)

			assert.Equal(t, tc.expected, passengers, tc.key)
		}
	})
}

// INFO: passanger,  ticket, document

type ticketWholeInfo struct {
	Id         string                     `json:"id"`
	Provider   string                     `json:"provider"`
	FlyFrom    string                     `json:"flyFrom"`
	FlyTo      string                     `json:"flyTo"`
	Passengers []passengerTicketWholeInfo `json:"passengers,omitempty"`
}

type passengerTicketWholeInfo struct {
	Id         string                    `json:"id"`
	FirstName  string                    `json:"firstName"`
	LastName   string                    `json:"lastName"`
	MiddleName string                    `json:"middleName"`
	Documents  []documentTicketWholeInfo `json:"documents,omitempty"`
}

type documentTicketWholeInfo struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Number string `json:"number"`
}

type timeFieldsTicketWholeInfo struct {
	FlyAt     string `json:"flyAt"`
	ArriveAt  string `json:"arriveAt"`
	CreatedAt string `json:"createdAt"`
}

func (s *Suite) Test1tGetTicketWholeInfo() {
	t := s.T()

	tcs := []struct {
		key      string
		id       string
		expected ticketWholeInfo
	}{
		{
			key: "1",
			id:  s.utils.GetTicketByOffset(s.ctx, 0),
			expected: ticketWholeInfo{
				Id:       s.utils.GetTicketByOffset(s.ctx, 0),
				Provider: "Emirates",
				FlyFrom:  "Moscow",
				FlyTo:    "Hanoi",
				Passengers: []passengerTicketWholeInfo{
					{
						Id:         s.utils.GetPassengerByOffset(s.ctx, 0),
						FirstName:  "Riley",
						LastName:   "Scott",
						MiddleName: "Reed",
						Documents: []documentTicketWholeInfo{
							{
								Id:     s.utils.GetDocumentByOffset(s.ctx, 0),
								Type:   "Id card",
								Number: "5555666777",
							},
							{
								Id:     s.utils.GetDocumentByOffset(s.ctx, 1),
								Type:   "International passport",
								Number: "5555666888",
							},
						},
					},
				},
			},
		},
	}

	t.Run("", func(t *testing.T) {
		for _, tc := range tcs {
			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/v1/tickets/whole-info/%s", tc.id),
				nil,
			)
			assert.NoError(t, err, tc.key)

			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, tc.key)

			ticket := ticketWholeInfo{}
			err = json.NewDecoder(w.Result().Body).Decode(&ticket)
			assert.NoError(t, err, tc.key)

			assert.Equal(t, tc.expected, ticket, tc.key)

			timeFields := timeFieldsTicketWholeInfo{}
			err = json.NewDecoder(w.Body).Decode(&timeFields)
			assert.NoError(t, err, tc.key)

			_, err = convertTime(timeFields.FlyAt)
			assert.NoError(t, err, tc.key)
			_, err = convertTime(timeFields.ArriveAt)
			assert.NoError(t, err, tc.key)
			_, err = convertTime(timeFields.CreatedAt)
			assert.NoError(t, err, tc.key)
		}
	})
}
