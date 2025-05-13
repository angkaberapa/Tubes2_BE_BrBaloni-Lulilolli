package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func ImageListHandler(c *gin.Context) {
	file, err := os.Open("data/element_images.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file element_images.json"})
		return
	}
	defer file.Close()

	var images []scraper.ElementImage
	err = json.NewDecoder(file).Decode(&images)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca isi JSON"})
		return
	}

	c.JSON(http.StatusOK, images)
}

func SearchHandler(c *gin.Context) {
	algo := c.DefaultQuery("algo", "DFS")
	target := c.DefaultQuery("target", "Aquarium")
	totalrecipe := c.DefaultQuery("totalrecipe", "1")

	maxrecipe, err := strconv.Atoi(totalrecipe)
	if err != nil || maxrecipe <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid totalrecipe parameter"})
		return
	}

	if algo != "DFS" && algo != "BFS" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid algo parameter"})
		return
	}
	var result interface{}
	var nodeCount int
	var nodes []search.GraphNode
	var edges []search.GraphEdge
	var totalRecipe int
	if maxrecipe == 1 {
		totalRecipe = 1
	}

	startTime := time.Now()
	switch algo {
	case "DFS":
		switch maxrecipe {
		case 1:
			result, err, nodeCount = search.DFS(target, maxrecipe)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menjalankan DFS", "details": err.Error()})
				return
			}
			nodes, edges, err = search.TranslateOutputPathToGraph(result, target)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mentranslasi data graph", "details": err.Error()})
				return
			}

		default:
			recipesOutputFromDFS, errDfs, combinationsExplored, recipes := search.DFSMultipleRecipe(target, maxrecipe)
			nodeCount = combinationsExplored
			totalRecipe = recipes
			if errDfs != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menjalankan DFS", "details": errDfs.Error()})
				return
			}
			if recipesOutputFromDFS == nil {
				c.JSON(http.StatusOK, gin.H{
					"message":  fmt.Sprintf("Tidak ada resep yang ditemukan untuk %s", target),
					"nodes":    []search.GraphNode{}, // Kirim array kosong
					"edges":    []search.GraphEdge{}, // Kirim array kosong
					"duration": time.Since(startTime).String(),
				})
				return
			}

			// Panggil fungsi Translasi Anda
			nodes, edges, err = search.TranslateMultiRecipeOutputToGraph(target, recipesOutputFromDFS)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mentranslasi data graph", "details": err.Error()})
				return
			}
		}
	case "BFS":
		nodes, edges, nodeCount, err = search.BFS(target, maxrecipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menjalankan BFS", "details": err.Error()})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	endTime := time.Now()
	// 5. Response JSON
	c.JSON(http.StatusOK, gin.H{
		"target":       target,
		"nodes":        nodes,
		"edges":        edges,
		"duration":     fmt.Sprintf("%.2f ms", endTime.Sub(startTime).Seconds()*1000),
		"totalNodes":   nodeCount,
		"message":      fmt.Sprintf("Graph berhasil dibuat untuk %s", target),
		"totalRecipes": totalRecipe,
	})
}
