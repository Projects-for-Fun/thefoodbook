package integrationtests

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/auth"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlersTestSuite struct {
	suite.Suite

	router *chi.Mux
	ctx    context.Context
	db     neo4j.DriverWithContext
}

func TestUserHandlersTestSuite(t *testing.T) {
	if !runIntegrationTests() {
		t.Skip("Skipping integration test")
		return
	}

	userHandlingSuite := new(UserHandlersTestSuite)
	ctx, cancel := context.WithCancel(context.Background())
	userHandlingSuite.ctx = ctx
	defer cancel()

	suite.Run(t, userHandlingSuite)
}

func (suite *UserHandlersTestSuite) SetupSuite() {
	config, db, w, logger := setup(suite.T(), suite.ctx)
	suite.db = db

	suite.router = chi.NewRouter()

	suite.router.Use(LoggerTestMiddleware(logger))
	suite.router.Use(AuthTestMiddleware(config.JWTKey))

	suite.router.Post("/sign-up", w.HandleSignUp)
	suite.router.Post("/login", w.HandleLogin(config.JWTKey))
	suite.router.Post("/logout", w.HandleLogout)
}

func (suite *UserHandlersTestSuite) BeforeTest(_, _ string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := suite.db.NewSession(ctx, neo4j.SessionConfig{})

	defer func() {
		err := session.Close(ctx)
		if err != nil {
			suite.T().Fatal(err)
		}
	}()

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		cypher := `CREATE (u:User {
								id: randomUUID(),
								username: $username, 
								firstName: $firstName,  
								lastName: $lastName,  
								email: $email,  
								password: $password,
								createdAt: datetime()
							})
					   RETURN u.id`
		encryptedPass, err := auth.EncryptPassword("password")
		if err != nil {
			return 0, err
		}
		params := map[string]any{
			"username":  "username",
			"firstName": "First",
			"lastName":  "Last",
			"email":     "email@email.com",
			"password":  encryptedPass,
		}
		_, err = transaction.Run(ctx, cypher, params)
		if err != nil {
			return 0, err
		}
		return 0, err
	})

	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *UserHandlersTestSuite) AfterTest(_, _ string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := suite.db.NewSession(ctx, neo4j.SessionConfig{})

	defer func() {
		err := session.Close(ctx)
		if err != nil {
			suite.T().Fatal(err)
		}
	}()

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
		cypher := `MATCH (u:User) DELETE u`
		_, err := transaction.Run(ctx, cypher, map[string]any{})
		if err != nil {
			return 0, err
		}
		return 0, err
	})

	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *UserHandlersTestSuite) TestSignUpHandler() {
	tests := []struct {
		description string
		body        webservice.UserRequest
		status      int
	}{
		{
			description: "When trying to sign up with an existing username return 400",
			body: webservice.UserRequest{
				Username:  "username",
				FirstName: "First",
				LastName:  "Last",
				Email:     "email@random.com",
				Password:  "password",
			},
			status: http.StatusBadRequest,
		},
		{
			description: "When trying to sign up with an existing email return 400",
			body: webservice.UserRequest{
				Username:  "username_random",
				FirstName: "First",
				LastName:  "Last",
				Email:     "email@email.com",
				Password:  "password",
			},
			status: http.StatusBadRequest,
		},
		{
			description: "Happy case returns 201",
			body: webservice.UserRequest{
				Username:  "username_happycase",
				FirstName: "First",
				LastName:  "Last",
				Email:     "email@happycase.com",
				Password:  "password",
			},
			status: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.description, func(t *testing.T) {
			bodyJSON, err := json.Marshal(tt.body)
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

			assert.Equal(suite.T(), tt.status, res.StatusCode)
		})
	}
}

func (suite *UserHandlersTestSuite) TestLoginHandler() {
	tests := []struct {
		description   string
		hasAuthHeader bool
		username      string
		password      string
		cookieSet     bool
		status        int
	}{
		{
			description:   "When trying to sign up without Authorization header return 400",
			hasAuthHeader: false,
			cookieSet:     false,
			status:        http.StatusBadRequest,
		},
		{
			description:   "When trying to sign up with wrong username return 400",
			hasAuthHeader: true,
			username:      "",
			password:      "password",
			cookieSet:     false,
			status:        http.StatusUnauthorized,
		},
		{
			description:   "When trying to sign up with wrong password return 400",
			hasAuthHeader: true,
			username:      "username",
			password:      "",
			status:        http.StatusUnauthorized,
		},
		{
			description:   "Happy case return 200",
			hasAuthHeader: true,
			username:      "username",
			password:      "password1",
			cookieSet:     true,
			status:        http.StatusOK,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.description, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/login", nil)
			if err != nil {
				suite.T().Fatal(err)
			}

			if tt.hasAuthHeader {
				req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(tt.username+":"+tt.password)))
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

			assert.Equal(suite.T(), tt.status, res.StatusCode)
			assert.Equal(suite.T(), tt.cookieSet, res.Header.Get("Set-Cookie") != "")
		})
	}
}
