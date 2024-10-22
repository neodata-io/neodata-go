package entities

import (
	"github.com/golang-jwt/jwt/v5"
)

// Define a struct for JWT claims
type Claims struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Abilities []Ability `json:"abilities"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	jwt.RegisteredClaims
}

type Ability struct {
	Action  string `json:"action"`
	Subject string `json:"subject"`
}

/* Example
{
	"user_id": "123456",
	"username": "john_doe",
	"email": "john.doe@example.com",
	"roles": ["Admin", "User"],
	"scopes": ["read:resource1", "write:resource2"],
	"first_name": "John",
	"last_name": "Doe",
	"iss": "myapp.com",
	"aud": "myapp-users",
	"iat": 1638360000,
	"exp": 1638370000
  }
*/
