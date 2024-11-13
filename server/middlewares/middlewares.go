package middlewares

import (
	"log"
	"net/http"
)

func SomeMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("before middleware")
	next(rw, r)
	log.Println("after middleware")
}
