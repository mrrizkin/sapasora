package config

type App struct {
	Name    string `name:"name"    env:"APP_NAME,required"`
	Env     string `name:"env"     env:"APP_ENV,required"`
	URL     string `name:"url"     env:"APP_URL,required"`
	Port    int    `name:"port"    env:"APP_PORT,required"`
	Prefork bool   `name:"prefork" env:"APP_PREFORK,default=false"`
	Debug   bool   `name:"debug"   env:"APP_DEBUG,default=false"`

	Log struct {
		Name      string `name:"name"       env:"LOG_NAME,default=application"`
		Level     string `name:"level"      env:"LOG_LEVEL,default=debug"`
		Console   bool   `name:"console"    env:"LOG_CONSOLE,default=true"`
		File      bool   `name:"file"       env:"LOG_FILE,default=true"`
		Dir       string `name:"dir"        env:"LOG_DIR"`
		MaxSize   int    `name:"max_size"   env:"LOG_MAX_SIZE,default=50"`
		MaxAge    int    `name:"max_age"    env:"LOG_MAX_AGE,default=7"`
		MaxBackup int    `name:"max_backup" env:"LOG_MAX_BACKUP,default=20"`
		JSON      bool   `name:"json"       env:"LOG_JSON,default=true"`
	} `name:"log"`

	Swagger struct {
		Path string `name:"path" env:"SWAGGER_PATH,default=/docs/v3/openapi.json"`
	} `name:"swagger"`
}
