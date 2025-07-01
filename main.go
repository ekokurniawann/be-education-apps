package main

import (
	"be-education/config"
	"be-education/db"
	"be-education/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan atau tidak dapat dimuat. Menggunakan variabel lingkungan sistem.")
	}

	cfg := config.LoadConfig()

	gin.SetMode(cfg.Server.Mode)

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer func() {
		if err := db.Close(dbConn); err != nil {
			log.Printf("Gagal menutup koneksi database: %v", err)
		}
	}()

	r := router.InitRouter(dbConn, cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Memulai server di %s", cfg.Server.BaseURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gagal memulai server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Menerima sinyal shutdown. Mematikan server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server gagal mati dengan graceful: %v", err)
	}

	log.Println("Server dihentikan.")
}
