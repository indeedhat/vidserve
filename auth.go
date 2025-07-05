package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ApiAccessMiddleware verifies that an api token header is passed that
// matches the api token pulled from the environment variable
//
// Token must be passed via the X-Api-Token header
func ApiAccessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		passedToken := r.Header.Get("X-Api-Token")
		if passedToken == "" || passedToken != ApiToken.Get() {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(rw, r)
	}
}

// FileAccessMiddleware verifies that a jwt is passed that claims access to the requeted file
//
// JWT can be passed either as a bearer token in the Authorization header or via the query string
// under the token key
// the header vaule will be prefered over the query string if both are present
func FileAccessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "" {
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		token := r.Header.Get("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		}
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := verifyJwt(token)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims.File != path {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		next(rw, r)
	}
}

// JwtClaims provides custom claims for the file access JWT
type JwtClaims struct {
	jwt.RegisteredClaims

	File string `json:"file"`
}

// verifyJwt and extract claims for later use
func verifyJwt(jwtString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(jwtString, &JwtClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != "HS256" {
			return nil, errors.New("invalid algo")
		}

		return []byte(JwtSecret.Get()), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt")
	}

	return claims, nil
}
