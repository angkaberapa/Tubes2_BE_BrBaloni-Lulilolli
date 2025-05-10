package search

import (
	"fmt"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)


func BFS(target string, maxrecipe int) (<returnvalue>, error) { // ganti return value nya
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err
	}
	// masukkan kode disini

	
	return results, nil
}
