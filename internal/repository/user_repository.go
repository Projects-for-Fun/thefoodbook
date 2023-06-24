package repository

import (
	"context"
	"errors"

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
			//			`CREATE (u:User {id: randomUUID()})
			//			 SET
			//					u.usernam = $username,
			//					u.firstName = $firstName,
			//					u.lastName = $lastName,
			//					u.email = $email,
			//					u.password = $password,
			//					u.createdAt = datetime()
			//
			//		     RETURN u.id
			//`
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
