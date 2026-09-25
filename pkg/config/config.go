// Unrated Coder t.me/Unrated_Coder

package config

import (
	"os"
	"strconv"
)

type Config struct {
	APIID    int
	APIHash  string
	BotToken string
	MongoURI string
	OwnerID  int64
}

func LoadConfig() *Config {
	apiIDStr := os.Getenv("API_ID")
	if apiIDStr == "" {
		apiIDStr = os.Getenv("APP_ID")
	}
	apiID, _ := strconv.Atoi(apiIDStr)

	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		apiHash = os.Getenv("APP_HASH")
	}

	if apiIDStr != "" {
		_ = os.Setenv("API_ID", apiIDStr)
		_ = os.Setenv("APP_ID", apiIDStr)
	}
	if apiHash != "" {
		_ = os.Setenv("API_HASH", apiHash)
		_ = os.Setenv("APP_HASH", apiHash)
	}

	ownerIDStr := os.Getenv("OWNER_ID")
	ownerID, _ := strconv.ParseInt(ownerIDStr, 10, 64)

	return &Config{
		APIID:    apiID,
		APIHash:  apiHash,
		BotToken: os.Getenv("BOT_TOKEN"),
		MongoURI: os.Getenv("MONGO_URI"),
		OwnerID:  ownerID,
	}
}
