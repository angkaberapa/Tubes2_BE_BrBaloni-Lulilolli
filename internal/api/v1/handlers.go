package v1

import (
	"net/http"
	"time"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/search"
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

func SearchHandler(c *gin.Context) {
	target := "Aquarium"
	maxrecipe := 1

	// Call your search function here
	results, err, nodeCount := search.BFS(target, maxrecipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":   results,
		"nodeCount": nodeCount,
	})
}

func DFSHandler(c *gin.Context) {
	target := "Grilled cheese"
	maxrecipe := 1
	startTime := time.Now()
	// Call your search function here
	results, err, nodeCount := search.DFS(target, maxrecipe)

	endTime := time.Now()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":   results,
		"nodeCount": nodeCount,
		"duration":  endTime.Sub(startTime).String(),
	})
}

func DFSConcurrentHandler(c *gin.Context) {
	target := "Grilled cheese"
	maxrecipe := 1
	startTime := time.Now()
	// Call your search function here
	results, err, nodeCount := search.DFSConcurrent(target, maxrecipe)

	endTime := time.Now()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results":   results,
		"nodeCount": nodeCount,
		"duration":  endTime.Sub(startTime).String(),
	})
}
