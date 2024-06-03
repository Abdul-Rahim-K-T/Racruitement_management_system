package routes

import (
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/controllers"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/middleware"
	"github.com/gin-gonic/gin"
)

// func LoadUserRoutes(router *gin.Engine) {
// 	user := router.Group("/user")
// 	{
// 		user.GET("/jobs", controllers.ViewJobs)
// 		user.POST("/jobs/apply", middleware.AuthMiddleware,controllers.ApplyJob)
// 		// user.POST("/uploadResume", controllers.UploadResume)
// 	}
// }

// func LoadAdminRoutes(router *gin.Engine) {
// 	admin := router.Group("/admin")
// 	{
// 		admin.POST("/job", middleware.AuthMiddleware,controllers.CreateJob)
// 		admin.GET("/job/:job_id", controllers.ViewJob)
// 		admin.GET("/applicants", middleware.AuthMiddleware,controllers.ViewApplicants)
// 		admin.GET("/applicant/:applicant_id", middleware.AuthMiddleware,controllers.ViewApplicantData)
// 	}
// }

// func LoadAuthRoutes(router *gin.Engine) {
// 	router.POST("/signup", controllers.Signup)
// 	router.POST("/login", controllers.Login)
// }

func LoadAuthRoutes(router *gin.Engine) {
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/logout", middleware.AuthMiddleware, controllers.Logout)
	router.GET("/user/jobs", middleware.AuthMiddleware, controllers.ViewJobs)
	router.POST("/user/jobs/apply", middleware.AuthMiddleware, controllers.ApplyJob)
	router.POST("/user/uploadResume", middleware.AuthMiddleware, controllers.UploadResume)
	router.POST("/admin/job", middleware.AuthMiddleware, controllers.CreateJob)
	router.GET("/admin/job/:job_id", middleware.AuthMiddleware, controllers.ViewJob)
	router.GET("/admin/applicants", middleware.AuthMiddleware, controllers.ViewApplicants)
	router.GET("/admin/applicant/:applicant_id", middleware.AuthMiddleware, controllers.ViewApplicantData)
}
