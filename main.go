package main

import (
	"log"
	"main/data"
	"main/server"
	"net/http"
)

func main() {
	db := data.connectToDB()
	defer db.Close()
	s := server.createServer(db)
	s.makeEndPoints()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
