// package routes

// import (
//     "github.com/prabhatKr-1/lib-man-sys/backend/controllers"
//     "github.com/prabhatKr-1/lib-man-sys/backend/middleware"

//     "github.com/gin-gonic/gin"
// )

// func SetupRoutes(r *gin.Engine) {

// 	auth := r.Group("v1/auth/")
// 	{
// 		auth.POST("/signup", controllers.Signup)
// 		auth.POST("/login", controllers.Login)
// 	}
	
// 	owner := r.Group("v1/owner/").Use(middleware.AuthMiddleware("Owner"))
// 	{
// 		owner.POST("/password", controllers.UpdatePassword)
// 		owner.POST("/create-admin", controllers.CreateAdminUser)
// 		owner.GET("/logout", controllers.Logout)
// 	}
	
// 	admin := r.Group("v1/admin/").Use(middleware.AuthMiddleware("Admin", "Owner"))
// 	{
// 		admin.POST("/password", controllers.UpdatePassword)
// 		admin.POST("/create-reader", controllers.CreateReaderUser)
// 		admin.GET("/books/search", controllers.SearchBook)
// 		admin.POST("/books/add", controllers.AddBook)
// 		admin.PATCH("/books/:isbn", controllers.UpdateBook)
// 		admin.DELETE("/books/:isbn", controllers.DeleteBook)
// 		admin.GET("/requests/all", controllers.ListRequests)
// 		admin.POST("/requests/process", controllers.ProcessRequest)
// 		admin.GET("/logout", controllers.Logout)
// 	}
	
// 	reader := r.Group("v1/reader/").Use(middleware.AuthMiddleware("Reader"))
// 	{
// 		reader.POST("/password", controllers.UpdatePassword)
// 		reader.GET("/books/search", controllers.SearchBook)
// 		reader.POST("/books/requests", controllers.RaiseBookRequest)
// 		reader.GET("/logout", controllers.Logout)
// 	}
// }

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prabhatKr-1/lib-man-sys/backend/controllers"
	"github.com/prabhatKr-1/lib-man-sys/backend/middleware"
	"github.com/prabhatKr-1/lib-man-sys/backend/utils"
)

func SetupRoutes(r *gin.Engine) {
	auth := r.Group("v1/auth/")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
	}

	owner := r.Group("v1/owner/").Use(middleware.AuthMiddleware(utils.ValidateJWT, "Owner"))
	{
		owner.POST("/password", controllers.UpdatePassword)
		owner.POST("/create-admin", controllers.CreateAdminUser)
		owner.GET("/logout", controllers.Logout)
	}

	admin := r.Group("v1/admin/").Use(middleware.AuthMiddleware(utils.ValidateJWT, "Admin", "Owner"))
	{
		admin.POST("/password", controllers.UpdatePassword)
		admin.POST("/create-reader", controllers.CreateReaderUser)
		admin.GET("/books/search", controllers.SearchBook)
		admin.POST("/books/add", controllers.AddBook)
		admin.PATCH("/books/:isbn", controllers.UpdateBook)
		admin.DELETE("/books/:isbn", controllers.DeleteBook)
		admin.GET("/requests/all", controllers.ListRequests)
		admin.POST("/requests/process", controllers.ProcessRequest)
		admin.GET("/logout", controllers.Logout)
	}

	reader := r.Group("v1/reader/").Use(middleware.AuthMiddleware(utils.ValidateJWT, "Reader"))
	{
		reader.POST("/password", controllers.UpdatePassword)
		reader.GET("/books/search", controllers.SearchBook)
		reader.POST("/books/requests", controllers.RaiseBookRequest)
		reader.GET("/logout", controllers.Logout)
	}
}
