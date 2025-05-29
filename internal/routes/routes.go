package routes

import (
	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/handlers"
	"job_solition/internal/middleware"
	"job_solition/internal/models"
	"job_solition/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	repo := repository.NewRepository(postgres)
	authHandler := handlers.NewAuthHandler(repo, cfg)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
	}
}

func SetupUserRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	userHandler := handlers.NewUserHandler(postgres, cfg)

	users := router.Group("/users")

	authorized := users.Group("")
	authorized.Use(middleware.OptionalAuth(cfg))
	authorized.Use(middleware.RequireAuth())

	authorized.GET("/me", userHandler.GetProfile)
	authorized.PUT("/me", userHandler.UpdateProfile)
	authorized.GET("/me/reviews", userHandler.GetUserReviews)
}

func SetupCompanyRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	companyHandler := handlers.NewCompanyHandler(postgres, cfg)

	companies := router.Group("/companies")

	companies.GET("", companyHandler.GetCompanies)
	companies.GET("/:id", companyHandler.GetCompany)

	authorized := companies.Group("")
	authorized.Use(middleware.OptionalAuth(cfg))
	authorized.Use(middleware.RequireAuth())
}

func SetupReviewRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	reviewHandler := handlers.NewReviewHandler(postgres, cfg)

	reviews := router.Group("/reviews")

	optionalAuth := reviews.Group("")
	optionalAuth.Use(middleware.OptionalAuth(cfg))

	optionalAuth.GET("/:id", reviewHandler.GetReview)
	optionalAuth.GET("/company/:companyId", reviewHandler.GetCompanyReviews)

	authorized := reviews.Group("")
	authorized.Use(middleware.OptionalAuth(cfg))
	authorized.Use(middleware.RequireAuth())

	authorized.POST("", reviewHandler.CreateReview)
	authorized.POST("/:id/useful", reviewHandler.MarkReviewAsUseful)
	authorized.DELETE("/:id/useful", reviewHandler.RemoveUsefulMark)

	moderation := reviews.Group("")
	moderation.Use(middleware.OptionalAuth(cfg))
	moderation.Use(middleware.RequireAuth())
}

func SetupCityRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	cityHandler := handlers.NewCityHandler(postgres, cfg)

	cities := router.Group("/cities")
	{
		cities.GET("", cityHandler.GetCities)
		cities.GET("/search", cityHandler.SearchCities)
	}
}

func SetupIndustryRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	industryHandler := handlers.NewIndustryHandler(postgres, cfg)

	industries := router.Group("/industries")
	{
		industries.GET("", industryHandler.GetIndustries)
		industries.GET("/company/:id", industryHandler.GetCompanyIndustries)

		authorized := industries.Group("")
		authorized.Use(middleware.OptionalAuth(cfg))
		authorized.Use(middleware.RequireAuth())
		authorized.PUT("/:id/color", industryHandler.UpdateIndustryColor)
	}
}

func SetupRatingCategoryRoutes(router *gin.RouterGroup, repo *repository.Repository) {
	ratingCategoryHandler := handlers.NewRatingCategoryHandler(repo)

	ratingCategories := router.Group("/rating-categories")
	{
		ratingCategories.GET("", ratingCategoryHandler.GetAll)
		ratingCategories.GET("/:id", ratingCategoryHandler.GetByID)
	}
}

func SetupBenefitTypeRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	benefitTypeHandler := handlers.NewBenefitTypeHandler(postgres, cfg)

	benefitTypes := router.Group("/benefit-types")
	{
		benefitTypes.GET("", benefitTypeHandler.GetAll)
		benefitTypes.GET("/:id", benefitTypeHandler.GetByID)
	}
}

func SetupEmploymentPeriodRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	employmentPeriodHandler := handlers.NewEmploymentPeriodHandler(postgres, cfg)

	employmentPeriods := router.Group("/employment-periods")
	{
		employmentPeriods.GET("", employmentPeriodHandler.GetAll)
		employmentPeriods.GET("/:id", employmentPeriodHandler.GetByID)
	}
}

func SetupEmploymentTypeRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	employmentTypeHandler := handlers.NewEmploymentTypeHandler(postgres, cfg)

	employmentTypes := router.Group("/employment-types")
	{
		employmentTypes.GET("", employmentTypeHandler.GetAll)
		employmentTypes.GET("/:id", employmentTypeHandler.GetByID)
	}
}

func SetupAdminRoutes(router *gin.RouterGroup, postgres *db.PostgreSQL, cfg *config.Config) {
	repo := repository.NewRepository(postgres)
	adminHandler := handlers.NewAdminHandler(repo, cfg)
	companyHandler := handlers.NewCompanyHandler(postgres, cfg)
	reviewHandler := handlers.NewReviewHandler(postgres, cfg)

	admin := router.Group("/admin")
	admin.Use(middleware.OptionalAuth(cfg))
	admin.Use(middleware.RequireAuth())
	admin.Use(middleware.RequireRoleMiddleware(models.RoleAdmin))

	admin.GET("/statistics", adminHandler.GetStatistics)

	admin.POST("/companies", companyHandler.CreateCompany)
	admin.PUT("/companies/:id", companyHandler.UpdateCompany)
	admin.DELETE("/companies/:id", companyHandler.DeleteCompany)

	admin.GET("/users", adminHandler.GetUsers)
	admin.GET("/users/:id", adminHandler.GetUser)
	admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
	admin.DELETE("/users/:id", adminHandler.DeleteUser)

	admin.POST("/rating-categories", adminHandler.CreateRatingCategory)
	admin.PUT("/rating-categories/:id", adminHandler.UpdateRatingCategory)
	admin.DELETE("/rating-categories/:id", adminHandler.DeleteRatingCategory)

	admin.PUT("/reviews/:id", adminHandler.UpdateReview)
	admin.DELETE("/reviews/:id", adminHandler.DeleteReview)
	admin.GET("/reviews/moderation/pending", reviewHandler.GetPendingReviews)
	admin.GET("/reviews/moderation/approved", reviewHandler.GetApprovedReviews)
	admin.GET("/reviews/moderation/rejected", reviewHandler.GetRejectedReviews)
	admin.PUT("/reviews/:id/approve", reviewHandler.ApproveReview)
	admin.PUT("/reviews/:id/reject", reviewHandler.RejectReview)

	admin.POST("/cities", adminHandler.CreateCity)
	admin.PUT("/cities/:id", adminHandler.UpdateCity)
	admin.DELETE("/cities/:id", adminHandler.DeleteCity)

	admin.POST("/industries", adminHandler.CreateIndustry)
	admin.PUT("/industries/:id", adminHandler.UpdateIndustry)
	admin.DELETE("/industries/:id", adminHandler.DeleteIndustry)

	admin.POST("/benefit-types", adminHandler.CreateBenefitType)
	admin.PUT("/benefit-types/:id", adminHandler.UpdateBenefitType)
	admin.DELETE("/benefit-types/:id", adminHandler.DeleteBenefitType)

	admin.POST("/employment-periods", adminHandler.CreateEmploymentPeriod)
	admin.PUT("/employment-periods/:id", adminHandler.UpdateEmploymentPeriod)
	admin.DELETE("/employment-periods/:id", adminHandler.DeleteEmploymentPeriod)

	admin.POST("/employment-types", adminHandler.CreateEmploymentType)
	admin.PUT("/employment-types/:id", adminHandler.UpdateEmploymentType)
	admin.DELETE("/employment-types/:id", adminHandler.DeleteEmploymentType)
}

func SetupSuggestionRoutes(router *gin.RouterGroup, repo *repository.Repository, cfg *config.Config) {
	suggestionHandler := handlers.NewSuggestionHandlers(repo)

	suggestions := router.Group("/suggestions")

	suggestions.POST("", suggestionHandler.CreateSuggestion)

	adminSuggestions := suggestions.Group("")
	adminSuggestions.Use(middleware.OptionalAuth(cfg))
	adminSuggestions.Use(middleware.RequireAuth())
	adminSuggestions.Use(middleware.RequireRoleMiddleware(models.RoleAdmin))

	adminSuggestions.GET("", suggestionHandler.GetAllSuggestions)
	adminSuggestions.DELETE("/:id", suggestionHandler.DeleteSuggestion)
}

func SetupAllRoutes(router *gin.Engine, postgres *db.PostgreSQL, cfg *config.Config) {
	repo := repository.NewRepository(postgres)

	api := router.Group("/api")

	SetupAuthRoutes(api, postgres, cfg)
	SetupUserRoutes(api, postgres, cfg)
	SetupCompanyRoutes(api, postgres, cfg)
	SetupReviewRoutes(api, postgres, cfg)
	SetupCityRoutes(api, postgres, cfg)
	SetupIndustryRoutes(api, postgres, cfg)
	SetupRatingCategoryRoutes(api, repo)
	SetupBenefitTypeRoutes(api, postgres, cfg)
	SetupEmploymentPeriodRoutes(api, postgres, cfg)
	SetupEmploymentTypeRoutes(api, postgres, cfg)
	SetupAdminRoutes(api, postgres, cfg)
	SetupSuggestionRoutes(api, repo, cfg)

	apiV1 := router.Group("/api/v1")

	SetupAuthRoutes(apiV1, postgres, cfg)
	SetupUserRoutes(apiV1, postgres, cfg)
	SetupCompanyRoutes(apiV1, postgres, cfg)
	SetupReviewRoutes(apiV1, postgres, cfg)
	SetupCityRoutes(apiV1, postgres, cfg)
	SetupIndustryRoutes(apiV1, postgres, cfg)
	SetupRatingCategoryRoutes(apiV1, repo)
	SetupBenefitTypeRoutes(apiV1, postgres, cfg)
	SetupEmploymentPeriodRoutes(apiV1, postgres, cfg)
	SetupEmploymentTypeRoutes(apiV1, postgres, cfg)
	SetupAdminRoutes(apiV1, postgres, cfg)
	SetupSuggestionRoutes(apiV1, repo, cfg)
}
