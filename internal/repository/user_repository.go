package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/adapter"
	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateUserRepoFunc(db neo4j.DriverWithContext) adapter.CreateUserRepo {
	return func(ctx context.Context, user domain.User) (userID *uuid.UUID, err error) {
		session := db.NewSession(ctx, neo4j.SessionConfig{})

		defer func() {
			sessionError := session.Close(ctx)
			err = errors.Join(err, sessionError)
		}()

		result, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
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
			params := map[string]any{
				"username":  user.Username,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"password":  user.Password,
			}
			result, err := transaction.Run(ctx, cypher, params)
			if err != nil {
				return 0, err
			}

			record, err := result.Single(ctx)
			if err != nil {
				return 0, err
			}

			id, _ := record.Get("u.id")

			return id, nil
		})

		if neo4jError, ok := err.(*neo4j.Neo4jError); ok && neo4jError.Title() == "ConstraintValidationFailed" {
			return nil, domain.ErrUserExists
		}

		id, err := uuid.Parse(result.(string))
		return &id, err
	}
}

func LoginUserRepoFunc(db neo4j.DriverWithContext) adapter.ValidateLoginUserRepo {
	return func(ctx context.Context, username string) (loggedUser *domain.User, err error) {
		session := db.NewSession(ctx, neo4j.SessionConfig{})

		defer func() {
			sessionError := session.Close(ctx)
			err = errors.Join(err, sessionError)
		}()

		result, err := session.ExecuteRead(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
			cypher := `MATCH (u:User { username: $username }) RETURN u`
			params := map[string]any{
				"username": username,
			}

			result, err := transaction.Run(ctx, cypher, params)
			if err != nil {
				return nil, err
			}

			record, err := result.Single(ctx)
			if err != nil {
				return nil, domain.ErrInvalidUsernameOrPassword
			}

			u, _ := record.Get("u")

			return u, nil
		})

		if err != nil {
			return nil, err
		}

		userProps := result.(neo4j.Node).Props
		userID, err := uuid.Parse(userProps["id"].(string))
		if err != nil {
			return nil, err
		}

		userLogged := domain.User{
			ID:        userID,
			Username:  userProps["username"].(string),
			FirstName: userProps["firstName"].(string),
			LastName:  userProps["lastName"].(string),
			Email:     userProps["email"].(string),
			Password:  userProps["password"].(string),     // encrypted
			CreatedAt: userProps["createdAt"].(time.Time), // encrypted
		}

		return &userLogged, nil
	}
}

func SetLoginUserRepo(db neo4j.DriverWithContext) adapter.SetLoginUserRepo {
	return func(ctx context.Context, username string) (err error) {
		session := db.NewSession(ctx, neo4j.SessionConfig{})

		defer func() {
			sessionError := session.Close(ctx)
			err = errors.Join(err, sessionError)
		}()

		_, err = session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (interface{}, error) {
			cypher := `MERGE (u:User { username: $username })
					   ON MATCH SET u.modifiedAt = datetime()
					   RETURN u.id`
			params := map[string]any{
				"username": username,
			}
			result, err := transaction.Run(ctx, cypher, params)
			if err != nil {
				return 0, err
			}

			record, err := result.Single(ctx)
			if err != nil {
				return 0, err
			}

			id, _ := record.Get("u.id")

			return id, nil
		})

		return err
	}
}
