package search

import (
	"fmt"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)

// func BFS(target string, maxrecipe int) (interface{}, error, int) { // ganti return value nya
// 	elements, err := scraper.LoadElementsFromFile()
// 	if err != nil {
// 		return nil, err, -1
// 	}

// 	results, nodeCount := findCombinationRouteBFS(elements[target], maxrecipe)
// 	// masukkan kode disini
// 	if results == nil {
// 		return nil, fmt.Errorf("no valid combination found"), nodeCount
// 	}
// 	fmt.Println("Results:", printInterface(results))
// 	return results, nil, nodeCount
// }

// func findCombinationRouteBFS(target *scraper.Element, recipeTarget int) (interface{}, int, []interface{}) {
// 	if isBasicElement(target) {
// 		return target.Name, 1, []interface{}{target}

// 	}
// 	elements, err := scraper.LoadElementsFromFile()
// 	if err != nil {
// 		return nil, 0, nil
// 	}

// 	combinationsNumber := len(target.Combinations)
// 	combinationsChecked := 1
// 	base := []interface{}{}
// 	for i := 0; i < combinationsNumber; i++ {
// 		if elements[target.Combinations[i].LeftName].Tier < target.Tier && elements[target.Combinations[i].RightName].Tier < target.Tier {
// 			base = append(base, combinationToElements(target.Combinations[i]))
// 		} else {
// 			// fmt.Println("kombinasi tidak valid")
// 		}
// 	}
// 	oldBase := base
// 	firstIteration := true
// 	listOfElements := []interface{}{target}
// 	// fmt.Println("Base:", printInterface(base))
// 	// fmt.Println("Base length:", len(base))

// 	for {
// 		numberOfRecipesObtained := 0

// 		for i, item := range base {
// 			// fmt.Println("Item:", printInterface(item))
// 			var itemLeft, itemRight interface{}
// 			if items, ok := item.([]*scraper.Element); ok && len(items) >= 2 {
// 				itemLeft = item.([]*scraper.Element)[0]
// 				itemRight = item.([]*scraper.Element)[1]
// 			}
// 			if items, ok := item.([]interface{}); ok && len(items) >= 2 {
// 				itemLeft = item.([]interface{})[0]
// 				itemRight = item.([]interface{})[1]
// 			}

// 			// fmt.Println("Item Left:", printInterface(itemLeft))
// 			// fmt.Println("Item Right:", printInterface(itemRight))

// 			resultLeftBool := false
// 			var resultLeft interface{}
// 			checkerLeft, countLeft, outputBaseLeft := checkInterface(itemLeft)
// 			if len(outputBaseLeft) == 1 {
// 				if x, ok := outputBaseLeft[0].(*scraper.Element); ok {
// 					resultLeftBool = true
// 					resultLeft = x
// 				}
// 			}
// 			combinationsChecked += countLeft // Hitung elemen yang diperiksa

// 			resultRightBool := false
// 			var resultRight interface{}
// 			checkerRight, countRight, outputBaseRight := checkInterface(itemRight)
// 			if len(outputBaseRight) == 1 {
// 				if x, ok := outputBaseRight[0].(*scraper.Element); ok {
// 					resultRightBool = true
// 					resultRight = x
// 				}
// 			}
// 			combinationsChecked += countRight // Hitung elemen yang diperiksa

// 			var outputBase []interface{}
// 			if resultLeftBool && resultRightBool {
// 				outputBase = []interface{}{resultLeft, resultRight}
// 			} else if resultLeftBool && !resultRightBool {
// 				outputBase = []interface{}{resultLeft, outputBaseRight}
// 			} else if !resultLeftBool && resultRightBool {
// 				outputBase = []interface{}{outputBaseLeft, resultRight}
// 			} else if !resultLeftBool && !resultRightBool {
// 				outputBase = []interface{}{outputBaseLeft, outputBaseRight}
// 			}

// 			if !firstIteration {
// 				// fmt.Println("Output Base:", printInterface(outputBase))
// 				// fmt.Println("oldBase[i]:", printInterface(oldBase[i]))
// 				// fmt.Println("Combinations Checked dikurangi:", compareInterfaceSlices(outputBase, oldBase[i]))
// 				combinationsChecked -= compareInterfaceSlices(outputBase, oldBase[i]) // Hitung elemen yang diperiksa
// 			}

// 			// fmt.Println("checkInterface:", printInterface(outputBase), "total checked:", combinationsChecked)
// 			if checkerLeft && checkerRight {
// 				// fmt.Println("Item:", printInterface(item))
// 				// return outputBase, combinationsChecked
// 			}
// 		}

// 		if firstIteration {
// 			firstIteration = false
// 		}

// 		newBase := []interface{}{}
// 		for _, item := range base {
// 			// fmt.Println("item:", printInterface(item))
// 			expanded, resultIterated := expandInterface(item)
// 			// fmt.Println("resultIterated:", resultIterated)
// 			listOfElements = addToIterated(listOfElements, resultIterated)
// 			// fmt.Println("Expanded:", expanded)

// 			// if expanded != nil {
// 			//     switch v := expanded.(type) {
// 			//     case []interface{}:
// 			//         for _, scraper.Element := range v {
// 			//             newBase = append(newBase, scraper.Element)
// 			//         }
// 			//     default:
// 			//         newBase = append(newBase, v)
// 			//     }
// 			// }

// 			if expanded != nil {
// 				newBase = append(newBase, expanded)
// 			}
// 			// fmt.Println("New Base:", printInterface(base))
// 		}
// 		oldBase = base
// 		base = newBase

// 		output := []interface{}{}
// 		for _, item := range oldBase {
// 			numberOfRecipesObtained += countRecipesObtained(item)
// 			output = append(output, item)
// 			if numberOfRecipesObtained > recipeTarget {
// 				return output, combinationsChecked, listOfElements

// 			}
// 		}
// 	}
// }

// func countRecipesObtained(item interface{}) int {
// 	count := 0
// 	switch v := item.(type) {
// 	case *scraper.Element:
// 		if isBasicElement(v) {
// 			return 1
// 		} else {
// 			return 0
// 		}
// 	case []*scraper.Element:
// 		if isBasicElement(v[0]) && isBasicElement(v[1]) {
// 			count++
// 		} else {
// 			return 0
// 		}
// 	case []interface{}:
// 		for _, el := range v {
// 			if el_1, ok := el.([]*scraper.Element); ok {
// 				count += countRecipesObtained(el_1[0]) * countRecipesObtained(el_1[1])
// 			} else if el_2, ok := el.([]interface{}); ok {
// 				for i, _ := range el_2 {
// 					count += countRecipesObtained(el_2[i])
// 				}
// 			}
// 		}
// 	default:
// 		return 0
// 	}
// 	return count
// }

// func checkInterface(a interface{}) (bool, int, []interface{}) {
// 	elementsChecked := 0
// 	switch res := a.(type) {
// 	case []interface{}:
// 		output := false
// 		outputInterface := []interface{}{}
// 		for _, v := range res {
// 			checker, count, outputToAdd := checkConsistsOfBasicElements(v)
// 			elementsChecked += count
// 			if checker {
// 				output = true
// 				outputInterface = append(outputInterface, outputToAdd)
// 				// fmt.Println("Outputinterface:", printInterface(outputInterface))
// 			}
// 		}

// 		if output {
// 			return true, elementsChecked, outputInterface
// 		}
// 		return false, elementsChecked, res
// 	case interface{}:
// 		checker, count, outputToAdd := checkConsistsOfBasicElements(res)
// 		return checker, count, []interface{}{outputToAdd}
// 	default:
// 		return false, 0, nil
// 	}
// }

// func checkConsistsOfBasicElements(item interface{}) (bool, int, interface{}) {
// 	elementsChecked := 0
// 	switch v := item.(type) {
// 	case *scraper.Element:
// 		return isBasicElement(v), 1, v
// 	case []*scraper.Element:
// 		output := true
// 		outputInterface := []interface{}{}
// 		for _, el := range v {
// 			elementsChecked++
// 			if isBasicElement(el) {
// 				outputInterface = append(outputInterface, el)
// 			} else {
// 				output = false
// 				outputInterface = append(outputInterface, convertToInterfaceSlice(expandElement(el)))
// 			}
// 		}

// 		// if output {
// 		// 	return true, elementsChecked, outputInterface
// 		// }
// 		return output, elementsChecked, outputInterface
// 	case []interface{}:
// 		output := true
// 		outputInterface := []interface{}{}
// 		for _, x := range v {
// 			checker, count, outputToAdd := checkConsistsOfBasicElements(x)
// 			elementsChecked += count
// 			if !checker {
// 				output = false
// 			}
// 			outputInterface = append(outputInterface, outputToAdd)
// 		}

// 		if output {
// 			return true, elementsChecked, outputInterface
// 		}
// 		return false, elementsChecked, outputInterface
// 	default:
// 		return false, 0, nil
// 	}
// }

// func addToIterated(iterated []interface{}, resultIterated interface{}) []interface{} {

// 	if resultIterated1, ok := resultIterated.([]interface{}); ok {
// 		// fmt.Println("resultIterated1:", printInterface(resultIterated1))
// 		// fmt.Println("iterated:", printInterface(iterated))
// 		for i, _ := range resultIterated1 {
// 			exists := false
// 			for j, _ := range iterated {
// 				if resultIterated1[i] == iterated[j] {
// 					exists = true
// 					break
// 				}
// 			}

// 			if !exists {
// 				iterated = append(iterated, resultIterated1[i])
// 			}
// 		}

// 		return iterated
// 	}

// 	return nil
// }

// func expandInterface(item interface{}) ([]interface{}, []interface{}) {
// 	switch v := item.(type) {
// 	case []interface{}:
// 		result := []interface{}{}
// 		iterated := []interface{}{}
// 		for _, el := range v {
// 			expanded, resultIterated := expandInterface(el)
// 			if expanded != nil {
// 				result = append(result, expanded)
// 			}
// 			for i, _ := range resultIterated {
// 				iterated = addToIterated(iterated, resultIterated[i])
// 			}
// 		}
// 		return result, iterated
// 	case []*scraper.Element:
// 		result := []interface{}{}
// 		for _, element := range v {
// 			expanded, _ := expandInterface(element)
// 			if expanded != nil {
// 				result = append(result, expanded)
// 			}
// 		}
// 		return result, []interface{}{v[0], v[1]}
// 	case *scraper.Element:
// 		if isBasicElement(v) {
// 			return []interface{}{v}, []interface{}{v}
// 		} else {
// 			return convertToInterfaceSlice(expandElement(v)), []interface{}{v}
// 		}

// 	default:
// 		return nil, nil
// 	}
// }

// func convertToInterfaceSlice(data []interface{}) []interface{} {
// 	result := make([]interface{}, len(data))
// 	for i, elementSlice := range data {
// 		result[i] = elementSlice
// 	}
// 	return result
// }

// func expandElement(element *scraper.Element) []interface{} {
// 	combinations := elementToCombinations(element)
// 	// fmt.Println("Combinations Length:", len(combinations))
// 	elements, err := scraper.LoadElementsFromFile()
// 	if err != nil {
// 		return nil
// 	}
// 	var output []interface{}
// 	for _, combination := range combinations {
// 		if elements[combination.LeftName].Tier < element.Tier && elements[combination.RightName].Tier < element.Tier {
// 			output = append(output, combinationToElements(combination))
// 		} else {
// 			// fmt.Println("kombinasi tidak valid")
// 		}
// 	}

// 	// fmt.Println("Output: %v", output)
// 	return output
// }

// Fungsi untuk mengecek apakah elemen adalah elemen dasar
func isBasicElement(element *scraper.Element) bool {
	return element.Name == "Water" || element.Name == "Earth" || element.Name == "Fire" || element.Name == "Air"
}

// // elementToCombinations mengubah array elemen menjadi [kombinasi1, kombinasi2, kombinasi3]
// func elementToCombinations(element *scraper.Element) []*scraper.Combination {
// 	var combinations []*scraper.Combination
// 	for _, combo := range element.Combinations {
// 		combinations = append(combinations, combo)
// 	}
// 	return combinations
// }

// // combinationToElements mengubah kombinasi1 menjadi [elemenpenyusun1, elemenpenyusun2]
// func combinationToElements(combination *scraper.Combination) []*scraper.Element {
// 	elements, err := scraper.LoadElementsFromFile()
// 	if err != nil {
// 		return nil
// 	}
// 	return []*scraper.Element{elements[combination.LeftName], elements[combination.RightName]}
// }

// func compareInterfaceSlices(slice1, slice2 interface{}) int {
// 	// fmt.Println("Slice1:", printInterface(slice1), "Slice2:", printInterface(slice2))
// 	output := 0
// 	if v1, ok := slice1.([]interface{}); ok {
// 		if v2, ok := slice2.([]interface{}); ok {
// 			for i, _ := range v1 {
// 				output += compareInterfaceSlices(v1[i], v2[i])
// 			}
// 		}
// 		if v2, ok := slice2.([]*scraper.Element); ok {
// 			for i, _ := range v1 {
// 				output += compareInterfaceSlices(v1[i], v2[i])
// 			}
// 		}
// 	}

// 	if v1, ok := slice1.([]*scraper.Element); ok {
// 		if v2, ok := slice2.([]interface{}); ok {
// 			for i, _ := range v1 {
// 				output += compareInterfaceSlices(v1[i], v2[i])
// 			}
// 		}
// 		if v2, ok := slice2.([]*scraper.Element); ok {
// 			for i, _ := range v1 {
// 				output += compareInterfaceSlices(v1[i], v2[i])
// 			}
// 		}
// 	}

// 	if v1, ok := slice1.(*scraper.Element); ok {
// 		if v2, ok := slice2.(*scraper.Element); ok {
// 			if v1.Name == v2.Name {
// 				output++
// 			}
// 		}
// 	}
// 	return output
// }

func printInterface(item interface{}) string {
	switch v := item.(type) {
	case *scraper.Element:
		return v.Name
	case []*scraper.Element:
		names := "["
		for i, element := range v {
			names += element.Name
			if i < len(v)-1 {
				names += ", "
			}
		}
		names += "]"
		return names
	case []interface{}:
		names := "["
		for i, element := range v {
			names += printInterface(element)
			if i < len(v)-1 {
				names += ", "
			}
		}
		names += "]"
		return names
	case [][]*scraper.Element:
		names := "["
		for i, elementSlice := range v {
			names += printInterface(elementSlice)
			if i < len(v)-1 {
				names += ", "
			}
		}
		names += "]"
		return names
	default:
		return fmt.Sprintf("%v", v) // Handle other types
	}
}
