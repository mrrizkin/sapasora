package config

import "time"

type Security struct {
	CSRF struct {
		Key        string        `name:"key"         env:"CSRF_KEY,default=X-CSRF-Token"`
		CookieName string        `name:"cookie_name" env:"CSRF_COOKIE_NAME,default=fiber_csrf_token"`
		SameSite   string        `name:"same_site"   env:"CSRF_SAME_SITE,default=Lax"`
		Secure     bool          `name:"secure"      env:"CSRF_SECURE,default=true"`
		Session    bool          `name:"session"     env:"CSRF_SESSION,default=true"`
		HTTPOnly   bool          `name:"http_only"   env:"CSRF_HTTP_ONLY,default=true"`
		Expiration time.Duration `name:"expiration"  env:"CSRF_EXPIRATION,default=3600s"`
	} `name:"csrf"`
}
