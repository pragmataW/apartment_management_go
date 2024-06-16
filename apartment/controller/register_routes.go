package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pragmataW/apartment_management/middleware"
)

func (ctrl *controller) RegisterRoutes(app *fiber.App, jwtKey string) {
	app.Post("/admin/login", ctrl.LoginAdmin)
	app.Post("/user/login", ctrl.LoginUser)
	app.Post("/logout", ctrl.Logout)
	app.Post("/payment/callback", ctrl.PaymentCallback)

	adminMiddleware := middleware.JwtMiddleware(jwtKey, "admin")
	app.Post("/flat/:flatNo", adminMiddleware, ctrl.CreateFlat)
	app.Put("/flat", adminMiddleware, ctrl.UpdateFlatOwner)
	app.Delete("/flat/:flatNo", adminMiddleware, ctrl.DeleteFlat)
	app.Get("/flat/:flatNo", adminMiddleware, ctrl.GetAllInfoAboutFlat)
	app.Get("/flat", adminMiddleware, ctrl.GetAllInfoAboutAllFlat)
	app.Post("/flat/:flatNo/dues", adminMiddleware, ctrl.AddDues)
	app.Delete("/flat/:flatNo/dues", adminMiddleware, ctrl.DeleteDues)
	app.Put("/flat/dues/price", adminMiddleware, ctrl.ChangeDuesPrice)
	app.Put("/flat/payday", adminMiddleware, ctrl.ChangePayDay)
	app.Post("/announcement", adminMiddleware, ctrl.AddAnnouncement)
	app.Post("/sendmail", adminMiddleware, ctrl.SendMail)

	userMiddleware := middleware.JwtMiddleware(jwtKey, "user")
	app.Post("/payment/token", userMiddleware, ctrl.GetPaymentToken)

	app.Get("/config/dues/price", middleware.JwtMiddleware(jwtKey, "admin", "user"), ctrl.GetDuesPrice)
	app.Get("/config/payday", middleware.JwtMiddleware(jwtKey, "admin", "user"), ctrl.GetPayDay)
	app.Get("/announcement", middleware.JwtMiddleware(jwtKey, "admin", "user"), ctrl.GetAllAnnouncements)
}
