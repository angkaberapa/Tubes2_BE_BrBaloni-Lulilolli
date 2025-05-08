// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// )

// // Element represents each element in the alchemy game
// type Element struct {
// 	Name      string     `json:"n"`
// 	Prime     *bool      `json:"prime,omitempty"` // nil if not present
// 	Combos    [][]string `json:"p"`               // Pairs that create this element
// 	CanCreate []string   `json:"c"`               // IDs this element helps create
// }

// // Combination represents a recipe: LeftID + RightID = ResultID
// type Combination struct {
// 	ResultID string
// 	LeftID   string
// 	RightID  string
// }

// func main() {
// 	// URL to fetch
// 	url := "https://unpkg.com/little-alchemy-2@0.0.1/dist/alchemy.json"

// 	// Step 1: Fetch the JSON
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to fetch JSON: %w", err))
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to read response: %w", err))
// 	}

// 	// Step 2: Parse JSON into map[string]Element
// 	elementMap := make(map[string]Element)
// 	if err := json.Unmarshal(body, &elementMap); err != nil {
// 		panic(fmt.Errorf("failed to unmarshal JSON: %w", err))
// 	}

// 	// Step 3: Build list of combinations
// 	var combinations []Combination
// 	for resultID, el := range elementMap {
// 		for _, pair := range el.Combos {
// 			if len(pair) == 2 {
// 				combinations = append(combinations, Combination{
// 					ResultID: resultID,
// 					LeftID:   pair[0],
// 					RightID:  pair[1],
// 				})
// 			}
// 		}
// 	}

// 	// Example: Print summary
// 	fmt.Printf("Total elements: %d\n", len(elementMap))
// 	fmt.Printf("Total combinations: %d\n\n", len(combinations))

// 	// Example: Print a few elements
// 	count := 0
// 	for id, el := range elementMap {
// 		isPrime := el.Prime != nil && *el.Prime
// 		fmt.Printf("ID: %s | Name: %-15s | Prime: %v | CanCreate: %d | CombosToCreateThis: %d\n",
// 			id, el.Name, isPrime, len(el.CanCreate), len(el.Combos))
// 		count++
// 		if count >= 10 {
// 			break
// 		}
// 	}
// }