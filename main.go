package main

import (
	"main/server"
	"net/http"
)

func main() {
	r := server.InitRouter()
	http.ListenAndServe(":8000", r)
}
