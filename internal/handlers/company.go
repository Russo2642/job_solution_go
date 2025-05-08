package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"job_solition/internal/config"
	"job_solition/internal/db"
	"job_solition/internal/middleware"
	"job_solition/internal/models"
	"job_solition/internal/repository"
	"job_solition/internal/utils"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewCompanyHandler(postgres *db.PostgreSQL, cfg *config.Config) *CompanyHandler {
	repo := repository.NewRepository(postgres)
	return &CompanyHandler{
		repo: repo,
		cfg:  cfg,
	}
}

// GetCompanies возвращает список компаний с фильтрацией и пагинацией
// @Summary Список компаний
// @Description Возвращает список компаний с возможностью фильтрации и пагинации
// @Tags companies
// @Accept json
// @Produce json
// @Param search query string false "Поисковый запрос"
// @Param industries query []int false "Фильтр по индустриям (может содержать несколько ID индустрий через запятую, например: industries=1,2,3)"
// @Param size query string false "Фильтр по размеру компании" Enums(small, medium, large, enterprise)
// @Param city query string false "Фильтр по названию города"
// @Param city_id query int false "Фильтр по ID города"
// @Param sort_by query string false "Поле для сортировки (name, rating, reviews_count, created_at)"
// @Param sort_order query string false "Порядок сортировки (asc, desc)"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /companies [get]
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	var filter models.CompanyFilter

	if sizeStr := c.Query("size"); sizeStr != "" {
		filter.Size = sizeStr
	}

	if cityIDStr := c.Query("city_id"); cityIDStr != "" {
		cityID, err := strconv.Atoi(cityIDStr)
		if err == nil && cityID > 0 {
			filter.CityID = &cityID
		}
	}

	if cityStr := c.Query("city"); cityStr != "" {
		filter.City = cityStr
	}

	if searchStr := c.Query("search"); searchStr != "" {
		filter.Search = searchStr
	}

	if sortByStr := c.Query("sort_by"); sortByStr != "" {
		filter.SortBy = sortByStr
	}

	if sortOrderStr := c.Query("sort_order"); sortOrderStr != "" {
		filter.SortOrder = sortOrderStr
	}

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			filter.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	industriesStr := c.Query("industries")
	if industriesStr != "" {
		industries, err := parseIndustriesParam(industriesStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации параметров", err)
			return
		}
		filter.Industries = industries
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.SortBy == "" {
		filter.SortBy = "rating"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	fmt.Printf("Применяемые фильтры: Size=%s, CityID=%v, Industries=%v\n",
		filter.Size, filter.CityID, filter.Industries)

	companies, total, err := h.repo.Companies.GetAll(c, filter)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении компаний", err)
		return
	}
	
	companySizes := make(map[string]string)
	for k, v := range models.CompanySizes {
		companySizes[k] = v
	}

	utils.Response(c, http.StatusOK, gin.H{
		"companies":     companies,
		"company_sizes": companySizes,
		"pagination": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
			"pages": (total + filter.Limit - 1) / filter.Limit,
		},
	})
}

// GetCompany возвращает информацию о компании по ID или slug
// @Summary Информация о компании
// @Description Возвращает детальную информацию о компании по её ID или slug
// @Tags companies
// @Accept json
// @Produce json
// @Param id_or_slug path string true "ID или slug компании"
// @Success 200 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 404 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /companies/{id_or_slug} [get]
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	idOrSlug := c.Param("id")

	var company *models.CompanyWithRatings
	var err error

	id, err := strconv.Atoi(idOrSlug)
	if err == nil {
		company, err = h.repo.Companies.GetByID(c, id)
	} else {
		company, err = h.repo.Companies.GetBySlug(c, idOrSlug)
	}

	if err != nil {
		if err.Error() == "компания не найдена" {
			utils.ErrorResponse(c, http.StatusNotFound, "Компания не найдена", nil)
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении компании", err)
		}
		return
	}

	utils.Response(c, http.StatusOK, company)
}

// CreateCompany создает новую компанию
// @Summary Создание компании
// @Description Создает новую компанию
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.CompanyInput true "Данные компании (name, size, city_id, industries обязательны)"
// @Success 201 {object} utils.ResponseDTO
// @Failure 400 {object} utils.ErrorResponseDTO
// @Failure 401 {object} utils.ErrorResponseDTO
// @Failure 500 {object} utils.ErrorResponseDTO
// @Router /companies [post]
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	fmt.Printf("CreateCompany: проверка авторизации\n")
	fmt.Printf("IsAuthenticatedKey=%v\n", c.GetBool(middleware.IsAuthenticatedKey))

	authHeader := c.GetHeader("Authorization")
	fmt.Printf("CreateCompany: Authorization header=%s\n", authHeader)

	userID, userIDExists := c.Get(middleware.UserIDKey)
	role, roleExists := c.Get(middleware.RoleKey)
	fmt.Printf("CreateCompany: userID=%v (exists=%v), role=%v (exists=%v)\n",
		userID, userIDExists, role, roleExists)

	roleValue, exists := c.Get(middleware.RoleKey)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Требуется авторизация", nil)
		return
	}

	if roleValue.(models.UserRole) != models.RoleAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "Недостаточно прав для создания компании", nil)
		return
	}

	var input models.CompanyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ошибка валидации", err)
		return
	}

	existingCompany, err := h.repo.Companies.GetByName(c, input.Name)
	if err == nil && existingCompany != nil {
		utils.ErrorResponse(c, http.StatusConflict, "Компания с таким названием уже существует", nil)
		return
	} else if err != nil && err.Error() != "компания не найдена" {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке существующей компании", err)
		return
	}

	industries, err := h.repo.Industries.GetByIDs(c, input.Industries)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке отраслей", err)
		return
	}

	if len(industries) != len(input.Industries) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Одна или несколько указанных отраслей не существуют", nil)
		return
	}

	if input.CityID != nil {
		city, err := h.repo.Cities.GetByID(c, *input.CityID)
		if err != nil {
			if err.Error() == "город не найден" {
				utils.ErrorResponse(c, http.StatusBadRequest, "Указанный город не существует", nil)
			} else {
				utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при проверке города", err)
			}
			return
		}
		if city == nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Указанный город не существует", nil)
			return
		}
	}

	company := models.NewCompany(input)

	id, err := h.repo.Companies.Create(c, company)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при сохранении компании", err)
		return
	}

	company.ID = id

	company.Slug = utils.GenerateUniqueSlug(company.Name, id)
	if err := h.repo.Companies.Update(c, company); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при обновлении slug компании", err)
		return
	}

	for _, industryID := range input.Industries {
		err = h.repo.Industries.AddCompanyIndustry(c, id, industryID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при добавлении отрасли к компании", err)
			return
		}
	}

	companyWithDetails, err := h.repo.Companies.GetByID(c, id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Ошибка при получении информации о созданной компании", err)
		return
	}

	utils.Response(c, http.StatusCreated, companyWithDetails)
}

func parseIndustriesParam(industriesStr string) ([]int, error) {
	var industries []int

	if strings.Contains(industriesStr, ",") {
		industriesArr := strings.Split(industriesStr, ",")
		for _, idStr := range industriesArr {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return nil, err
			}
			if id > 0 {
				industries = append(industries, id)
			}
		}
	} else {
		id, err := strconv.Atoi(industriesStr)
		if err != nil {
			return nil, err
		}
		if id > 0 {
			industries = append(industries, id)
		}
	}

	return industries, nil
}
