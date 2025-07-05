package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/indeedhat/dotenv"
	_ "github.com/indeedhat/dotenv/autoload"
)

const (
	ApiToken dotenv.String = "API_TOKEN"

	JwtSecret dotenv.String = "JWT_SECRET"
	JwtTtl    dotenv.Int    = "JWT_TTL"
)

func main() {
	r := NewRouter()
	{
		r.HandleFunc("GET /api/request_token/:file", requestToken, ApiAccessMiddleware)
		r.Handle("GET /files/",
			http.StripPrefix("/files/", http.FileServer(http.Dir("files"))),
			FileAccessMiddleware,
		)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	svr := &http.Server{
		Addr:    ":8087",
		Handler: r.mux,
	}

	go func() {
		log.Printf("ListenAndServer: %v", svr.ListenAndServe())
	}()

	<-quit
	log.Print("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Print("Server forced to shutdown after timeout")
	}
}

// requestToken generates a new jwt that can later be used to request access to a file
func requestToken(rw http.ResponseWriter, r *http.Request) {
	file := r.PathValue("file")

	claims := JwtClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(JwtTtl.Get()) * time.Second),
			),
		},
		file,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).
		SignedString(JwtSecret.Get())

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write([]byte(token))
}
