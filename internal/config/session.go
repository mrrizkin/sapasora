package config

type Session struct {
	CookieName string `name:"cookie_name" env:"SESSION_COOKIE_NAME,default=fiber"`
	Driver     string `name:"driver"      env:"SESSION_DRIVER,default=file"`
	HTTPOnly   bool   `name:"http_only"   env:"SESSION_HTTP_ONLY,default=true"`
	Secure     bool   `name:"secure"      env:"SESSION_SECURE,default=true"`
	SameSite   string `name:"same_site"   env:"SESSION_SAME_SITE,default=Lax"`
}
