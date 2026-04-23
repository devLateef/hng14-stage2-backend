package handlers

import (
	"net/http"
	"strconv"

	"insight-api/config"
	"insight-api/models"
	"insight-api/services"

	"github.com/gin-gonic/gin"
)

// allowedSortFields restricts sort_by to safe column names.
var allowedSortFields = map[string]bool{
	"age":                true,
	"created_at":         true,
	"gender_probability": true,
}

func GetProfiles(c *gin.Context) {
	db := config.DB.Model(&models.Profile{})

	// --- Filters ---
	if gender := c.Query("gender"); gender != "" {
		if gender != "male" && gender != "female" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
			return
		}
		db = db.Where("gender = ?", gender)
	}

	if ageGroup := c.Query("age_group"); ageGroup != "" {
		valid := map[string]bool{"child": true, "teenager": true, "adult": true, "senior": true}
		if !valid[ageGroup] {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
			return
		}
		db = db.Where("age_group = ?", ageGroup)
	}

	if countryID := c.Query("country_id"); countryID != "" {
		db = db.Where("country_id = ?", countryID)
	}

	if minAgeStr := c.Query("min_age"); minAgeStr != "" {
		val, err := strconv.Atoi(minAgeStr)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
			return
		}
		db = db.Where("age >= ?", val)
	}

	if maxAgeStr := c.Query("max_age"); maxAgeStr != "" {
		val, err := strconv.Atoi(maxAgeStr)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
			return
		}
		db = db.Where("age <= ?", val)
	}

	if minGPStr := c.Query("min_gender_probability"); minGPStr != "" {
		val, err := strconv.ParseFloat(minGPStr, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
			return
		}
		db = db.Where("gender_probability >= ?", val)
	}

	if minCPStr := c.Query("min_country_probability"); minCPStr != "" {
		val, err := strconv.ParseFloat(minCPStr, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
			return
		}
		db = db.Where("country_probability >= ?", val)
	}

	// --- Sorting ---
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "asc")

	if !allowedSortFields[sortBy] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
		return
	}
	if order != "asc" && order != "desc" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
		return
	}

	db = db.Order(sortBy + " " + order)

	// --- Pagination ---
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Invalid query parameters"})
		return
	}
	if limit > 50 {
		limit = 50
	}

	var total int64
	db.Count(&total)

	offset := (page - 1) * limit
	var profiles []models.Profile
	db.Offset(offset).Limit(limit).Find(&profiles)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   profiles,
	})
}

func SearchProfiles(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing or empty parameter"})
		return
	}

	filters, err := services.ParseQuery(q)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Unable to interpret query"})
		return
	}

	db := config.DB.Model(&models.Profile{})

	if filters.Gender != "" {
		db = db.Where("gender = ?", filters.Gender)
	}
	if filters.AgeGroup != "" {
		db = db.Where("age_group = ?", filters.AgeGroup)
	}
	if filters.CountryID != "" {
		db = db.Where("country_id = ?", filters.CountryID)
	}
	if filters.MinAge != nil {
		db = db.Where("age >= ?", *filters.MinAge)
	}
	if filters.MaxAge != nil {
		db = db.Where("age <= ?", *filters.MaxAge)
	}

	// Pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var total int64
	db.Count(&total)

	offset := (page - 1) * limit
	var profiles []models.Profile
	db.Offset(offset).Limit(limit).Find(&profiles)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   page,
		"limit":  limit,
		"total":  total,
		"data":   profiles,
	})
}
