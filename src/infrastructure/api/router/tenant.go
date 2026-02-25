package router

import (
	"github.com/myproject/api/infrastructure/api/handler"
	"github.com/myproject/api/infrastructure/core"
)

func TenantRoutes(r core.IMux) {
	// Public routes
	r.Group(func(r core.IMux) {
		// Auth
		r.Post("/api/v1/auth/sign-in/", handler.SignInHandler)
		r.Post("/api/v1/auth/sign-up/", handler.SignUpHandler)
	})

	// Private routes
	r.Group(func(r core.IMux) {
		r.Use(core.Authorization)

		// File upload
		r.Post("/api/v1/file/upload/", handler.UploadFile)

		r.Get("/api/v1/me/", handler.MeHandler)
		r.Post("/api/v1/auth/sign-out/", handler.SignOutHandler)

		// User routes
		r.Get("/api/v1/users/list/", handler.ListUserHandler)
		r.Get("/api/v1/users/{id}/detail/", handler.DetailUserHandler)
		r.Post("/api/v1/users/create/", handler.CreateUserHandler)
		r.Put("/api/v1/users/{id}/update/", handler.UpdateUserHandler)
		r.Delete("/api/v1/users/{id}/delete/", handler.DeleteUserHandler)
		r.Put("/api/v1/users/{id}/change-password/", handler.ChangePasswordUserHandler)

		// Role routes
		r.Get("/api/v1/roles/list/", handler.ListRoleHandler)
		r.Get("/api/v1/roles/{id}/detail/", handler.DetailRoleHandler)
		r.Post("/api/v1/roles/create/", handler.CreateRoleHandler)
		r.Put("/api/v1/roles/{id}/update/", handler.UpdateRoleHandler)
		r.Delete("/api/v1/roles/{id}/delete/", handler.DeleteRoleHandler)

		// Permission routes
		r.Get("/api/v1/permissions/list/", handler.ListPermissionHandler)
	})
}
