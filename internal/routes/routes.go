package routes

import (
	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/handlers"
	"job_solition/internal/middleware"
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
		auth.POST("/reset-password", authHandler.ResetPassword)
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

	authorized.POST("", companyHandler.CreateCompany)
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

	moderation.GET("/moderation/pending", reviewHandler.GetPendingReviews)
	moderation.PUT("/:id/approve", reviewHandler.ApproveReview)
	moderation.PUT("/:id/reject", reviewHandler.RejectReview)
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
}
