// Unrated Coder t.me/Unrated_Coder

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bot/pkg/bot"
	"bot/pkg/config"
	"bot/pkg/database"
	"bot/pkg/streamer"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func main() {
	log.SetOutput(bot.GlobalLogCapturer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		log.Printf("Starting health check server on port %s...", port)
		if err := http.ListenAndServe("0.0.0.0:"+port, mux); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()
	cfg := config.LoadConfig()
	if cfg.APIID == 0 || cfg.APIHash == "" || cfg.BotToken == "" {
		log.Fatal("API_ID, API_HASH, and BOT_TOKEN are required environment variables")
	}

	db, err := database.NewDatabase(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping database to verify connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := db.Ping(pingCtx); err != nil {
		pingCancel()
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	pingCancel()
	log.Println("Connected to MongoDB successfully!")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		Logger: bot.NewZapLogger(),
	})
	if err != nil {
		log.Fatalf("Failed to create Telegram client: %v", err)
	}

	api := tg.NewClient(client)
	str := streamer.NewStreamer(api)
	str.Start("127.0.0.1", 8080)
	defer str.Stop()

	b := bot.NewBot(cfg, db, str)
	log.Println("Bot starting...")
	if err := b.Run(ctx); err != nil {
		log.Fatalf("Bot crashed: %v", err)
	}
}
