package config

type Database struct {
	Driver   string `name:"driver"   env:"DB_DRIVER,required"`
	Host     string `name:"host"     env:"DB_HOST,required"`
	Port     int    `name:"port"     env:"DB_PORT,required"`
	Name     string `name:"name"     env:"DB_NAME,required"`
	Username string `name:"username" env:"DB_USERNAME,required"`
	Password string `name:"password" env:"DB_PASSWORD"`
	SSLMode  string `name:"ssl_mode" env:"DB_SSL_MODE,default=disable"`
}
