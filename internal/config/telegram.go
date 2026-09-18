package config

type Telegram struct {
	APIID   int    `name:"api_id"   env:"TELEGRAM_API_ID,required"`
	APIHash string `name:"api_hash" env:"TELEGRAM_API_HASH,required"`
}
