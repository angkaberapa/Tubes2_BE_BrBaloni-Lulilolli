package v1

import (
	"net/http"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
	"github.com/gin-gonic/gin"
)

func ScrapeHandler(c *gin.Context) {
	elements, err := scraper.ScrapeElements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, elements)
}
