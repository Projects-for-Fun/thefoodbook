package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlersTestSuite struct {
	suite.Suite
	router *chi.Mux
	ctx    context.Context
	cancel context.CancelFunc
}

func TestUserHandlersTestSuite(t *testing.T) {
	if !runIntegrationTests() {
		logger.Warn().Msg("Not running this test")
		return
	}

	userHandlingSuite := new(UserHandlersTestSuite)
	ctx, cancel := context.WithCancel(context.Background())
	userHandlingSuite.ctx = ctx
	userHandlingSuite.cancel = cancel
	//defer cancel()

	suite.Run(t, userHandlingSuite)
}

func (suite *UserHandlersTestSuite) SetupSuite() {

	setup(suite.T(), suite.ctx)

	fmt.Printf("%+v", suite.ctx)

	suite.router = chi.NewRouter()

	suite.router.Use(LoggerTestMiddleware(logger))
	suite.router.Post("/sign-up", w.HandleSignUp)
}

func (suite *UserHandlersTestSuite) TestHappyCase() {
	body := webservice.UserRequest{
		Username:  "username",
		FirstName: "First",
		LastName:  "Last",
		Email:     "email@abc.com",
		Password:  "password",
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		suite.T().Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/sign-up", bytes.NewBuffer(bodyJSON))
	if err != nil {
		suite.T().Fatal(err)
	}

	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)

	res := recorder.Result()
	defer func() {
		err := res.Body.Close()
		if err != nil {
			suite.T().Fatal(err)
		}
	}()

	assert.Equal(suite.T(), 201, res.StatusCode)
}
