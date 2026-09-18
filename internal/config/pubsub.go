package config

type Pubsub struct {
	Size int `name:"size" env:"PUBSUB_SIZE,default=100"`
}
