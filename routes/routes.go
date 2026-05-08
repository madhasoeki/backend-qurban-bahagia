package routes

import (
	"qurban/controllers"
	"qurban/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Public endpoints
	api.POST("/login", controllers.Login)
	api.GET("/stream", controllers.SSEHandler)
	api.GET("/dashboard/summary", controllers.GetDashboardSummary)
	api.GET("/dashboard/hewan", controllers.GetPublicHewan)

	// Authenticated: all roles can read hewan data
	allRoles := api.Group("/")
	allRoles.Use(middleware.AuthMiddleware(
		"admin", "pengawas", "jagal", "kulit",
		"cacah_daging", "cacah_tulang", "packing", "distribusi",
	))
	allRoles.GET("/hewan", controllers.GetHewan)

	// Operational: admin, koordinator, pengawas
	ops := api.Group("/")
	ops.Use(middleware.AuthMiddleware("admin", "koordinator_pengawas", "pengawas"))
	ops.PATCH("/hewan/:id/progress/:pos", controllers.UpdateProgressHewan)
	ops.PATCH("/hewan/:id/timbang", controllers.UpdateTimbangHewan)
	ops.PATCH("/hewan/:id/kelengkapan", controllers.UpdateKelengkapanHewan)

	// Packing station
	packing := api.Group("/")
	packing.Use(middleware.AuthMiddleware("admin", "packing"))
	packing.PATCH("/hewan/:id/packing", controllers.UpdatePackingHewan)

	// Distribution
	dist := api.Group("/")
	dist.Use(middleware.AuthMiddleware("admin", "distribusi"))
	dist.GET("/distribusi", controllers.GetAllDistribusi)
	dist.GET("/distribusi/:user_id", controllers.GetDistribusiUser)
	dist.PATCH("/distribusi/:user_id", controllers.UpdateDistribusi)

	// Admin only: user and hewan CRUD
	admin := api.Group("/")
	admin.Use(middleware.AuthMiddleware("admin"))
	admin.GET("/users", controllers.GetUsers)
	admin.POST("/users", controllers.CreateUser)
	admin.PUT("/users/:id", controllers.UpdateUser)
	admin.DELETE("/users/:id", controllers.DeleteUser)
	admin.POST("/hewan", controllers.CreateHewan)
	admin.PUT("/hewan/:id", controllers.UpdateHewan)
	admin.DELETE("/hewan/:id", controllers.DeleteHewan)
}
