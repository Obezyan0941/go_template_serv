package main

import (
	"log"
	"os"
	"test_backend/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var port string = os.Getenv("HOST_PORT")

	srv := server.NewServer()

	log.Printf("Сервер запущен на http://localhost:%s", port)
	if err := srv.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
