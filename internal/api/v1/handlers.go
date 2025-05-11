package v1

import (
	"fmt"
	"net/http"
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

func SearchHandler(c *gin.Context) {
	target := "Livestock"
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
	target := "Aquarium"
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

func DFSMultipleRecipeHandler(c *gin.Context) {
	// Dapatkan target dari query param, default ke "Aquarium" jika tidak ada
	targetElementName := c.DefaultQuery("element", "Aquarium")
	// Dapatkan maxrecipe dari query param, default ke 0 (tanpa batas)
	// Anda mungkin ingin validasi input ini lebih lanjut.
	maxRecipeQuery := c.DefaultQuery("maxrecipe", "1")
	maxRecipeCap, err := strconv.Atoi(maxRecipeQuery)
	if err != nil {
		maxRecipeCap = 1 // Default jika konversi gagal
	}

	startTime := time.Now()

	// Panggil fungsi DFS Anda
	// Asumsi: search.DFSMultipleRecipe(targetName string, maxRecipeCap int, allElements map[string]*scraper.Element)
	//                                 (recipesOutput []interface{}, totalUniquePaths int, err error)
	// Jika signature Anda berbeda, sesuaikan.
	// Dari kode Anda: results, err, nodeCount := search.DFSMultipleRecipe(target, maxrecipe)
	// Kita asumsikan `results` adalah `recipesOutput` dan `nodeCount` adalah `totalUniquePaths`.
	// `target` harus berupa `string` nama, bukan `*scraper.Element`.

	// recipesOutput, totalUniquePaths, combinationsExplored, errDfs := search.GetCorrectedMultipleRecipes(targetElementName, maxRecipeCap, allElements)
	// ATAU jika menggunakan signature dari kode Anda:
	recipesOutputFromDFS, errDfs, combinationsExplored, totalRecipe := search.DFSMultipleRecipe(targetElementName, maxRecipeCap)

	if errDfs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menjalankan DFS", "details": errDfs.Error()})
		return
	}
	if recipesOutputFromDFS == nil {
		// Ini mungkin bukan error, tapi tidak ada resep yang ditemukan.
		c.JSON(http.StatusOK, gin.H{
			"message":  fmt.Sprintf("Tidak ada resep yang ditemukan untuk %s", targetElementName),
			"nodes":    []search.GraphNode{}, // Kirim array kosong
			"edges":    []search.GraphEdge{}, // Kirim array kosong
			"duration": time.Since(startTime).String(),
		})
		return
	}

	// Panggil fungsi Translasi Anda
	graphNodes, graphEdges, errTranslate := search.TranslateMultiRecipeOutputToGraph(targetElementName, recipesOutputFromDFS)
	if errTranslate != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mentranslasi data graph", "details": errTranslate.Error()})
		return
	}

	endTime := time.Now()

	c.JSON(http.StatusOK, gin.H{
		"target":       targetElementName,
		"nodes":        graphNodes,
		"edges":        graphEdges,
		"duration":     endTime.Sub(startTime).String(),
		"totalNodes":   combinationsExplored, // Jika DFS Anda mengembalikan ini
		"message":      fmt.Sprintf("Graph berhasil dibuat untuk %s", targetElementName),
		"totalRecipes": totalRecipe,
	})
}

// target := "Aquarium"
// maxrecipe := 100
// // elements, err := scraper.LoadElementsFromFile()
// // if err != nil {
// // 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load elements"})
// // 	return
// // }

// startTime := time.Now()
// // Call your search function here
// results, err, nodeCount := search.DFSMultipleRecipe(target, maxrecipe)
// endTime := time.Now()
// if err != nil {
// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	return
// }

// c.JSON(http.StatusOK, gin.H{
// 	"results":   results,
// 	"nodeCount": nodeCount,
// 	"duration":  endTime.Sub(startTime).String(),
// })
// }

// dari frontend:
// async function fetchGraphData(elementName) {
//     try {
//         const response = await fetch(`http://localhost:8080/api/graph?element=${encodeURIComponent(elementName)}`);
//         if (!response.ok) {
//             const errorText = await response.text();
//             throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
//         }
//         const data = await response.json(); // Ini akan berisi { nodes: [...], edges: [...] }
//         console.log("Graph data diterima:", data);

//         // Sekarang Anda bisa menggunakan data.nodes dan data.edges
//         // untuk merender graph menggunakan library seperti Vis.js, Cytoscape.js, D3.js, dll.
//         // renderGraph(data.nodes, data.edges);

//     } catch (error) {
//         console.error("Gagal mengambil data graph:", error);
//     }
// }

// // Panggil fungsi untuk elemen tertentu
// // fetchGraphData("Plant"); // Ganti dengan elemen yang Anda inginkan
