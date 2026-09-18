package main

// main godoc
//
//	@title Sapasora
//	@version 1.0.0
//	@description Sapasora is gateway API for Whatsapp & Telegram
//
//	@license.name MIT
//	@license.url https://e-gitlab.prodak.id/nugrhrizki/sapasora/blob/main/LICENSE
//
//	@termsOfService https://e-gitlab.prodak.id/nugrhrizki/sapasora
//
//	@contact.name API Support
//	@contact.url https://e-gitlab.prodak.id/nugrhrizki/sapasora
//	@contact.email contact@mrrizkin.com
//
//	@host localhost:3000
//	@BasePath /
//	@schemes http https
//
//	@tag.name Account
//	@tag.description Account API
//
//	@tag.name APIKey
//	@tag.description APIKey API
//
//	@tag.name Device
//	@tag.description Device API
//
//	@tag.name DeviceToken
//	@tag.description DeviceToken API
//
//	@tag.name Gateway
//	@tag.description Gateway API
//
//	@tag.name Role
//	@tag.description Role API
//
//	@securityDefinitions.apikey X-API-KEY
//	@in header
//	@name Authorization
//	@description Provide a valid sk-dat-* or sk-dak-* token in the Authorization header; never put credentials in a URL.
func main() {
	serve()
}
