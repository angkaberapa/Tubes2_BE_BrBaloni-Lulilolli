package search

import (
	"fmt"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)


func BFS(target string, maxrecipe int) (interface{}, error, int) { // ganti return value nya
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err
	}

	results, nodeCount := findCombinationRouteBFS(elements[target])
	// masukkan kode disini

	return results, nil, nodeCount
}

func findCombinationRouteBFS(target *Element) (interface{}, int) {
    if (isBasicElement(target)) {
        return target.Name, 1
    }

    combinationsNumber := len(target.Combinations) 
    combinationsChecked := 1
    base := []interface{}{}
    for i := 0; i < combinationsNumber; i++ {
        if (target.Combinations[i].Left.Tier < target.Tier && target.Combinations[i].Right.Tier < target.Tier) {
            base = append(base, combinationToElements(target.Combinations[i]))
        } else {
            // fmt.Println("kombinasi tidak valid")
        }
    }
    oldBase := base
    firstIteration := true
    // fmt.Println("Base:", printInterface(base))
    // fmt.Println("Base length:", len(base))

    for {
        for i, item := range base {
            // fmt.Println("Item:", printInterface(item))
            var itemLeft, itemRight interface{}
            if items, ok := item.([]*Element); ok && len(items) >= 2 {
                itemLeft = item.([]*Element)[0]
                itemRight = item.([]*Element)[1]
            }
            if items, ok := item.([]interface{}); ok && len(items) >= 2 {
                itemLeft = item.([]interface{})[0]
                itemRight = item.([]interface{})[1]
            }

            // fmt.Println("Item Left:", printInterface(itemLeft))
            // fmt.Println("Item Right:", printInterface(itemRight))

            resultLeftBool := false
            var resultLeft interface{}
            checkerLeft, countLeft, outputBaseLeft := checkInterface(itemLeft)
            if len(outputBaseLeft) == 1 {
                if x, ok := outputBaseLeft[0].(*Element); ok {
                    resultLeftBool = true
                    resultLeft = x
                }
            }
            combinationsChecked += countLeft // Hitung elemen yang diperiksa
        
            resultRightBool := false
            var resultRight interface{}
            checkerRight, countRight, outputBaseRight := checkInterface(itemRight)
            if len(outputBaseRight) == 1 {
                if x, ok := outputBaseRight[0].(*Element); ok {
                    resultRightBool = true
                    resultRight = x
                }
            }
            combinationsChecked += countRight // Hitung elemen yang diperiksa

            var outputBase []interface{}
            if (resultLeftBool && resultRightBool) {
                outputBase = []interface{}{resultLeft, resultRight}
            } else if (resultLeftBool && !resultRightBool) {
                outputBase = []interface{}{resultLeft, outputBaseRight}
            } else if (!resultLeftBool && resultRightBool) {
                outputBase = []interface{}{outputBaseLeft, resultRight}
            } else if (!resultLeftBool && !resultRightBool) {
                outputBase = []interface{}{outputBaseLeft, outputBaseRight}
            }
            
            if (!firstIteration) {
                // fmt.Println("Output Base:", printInterface(outputBase))
                // fmt.Println("oldBase[i]:", printInterface(oldBase[i]))
                // fmt.Println("Combinations Checked dikurangi:", compareInterfaceSlices(outputBase, oldBase[i]))
                combinationsChecked -= compareInterfaceSlices(outputBase, oldBase[i]) // Hitung elemen yang diperiksa
            }

            // fmt.Println("checkInterface:", printInterface(outputBase), "total checked:", combinationsChecked)
            if checkerLeft && checkerRight {
                // fmt.Println("Item:", printInterface(item))
                return outputBase, combinationsChecked
            }
        }

        if (firstIteration) {
            firstIteration = false
        }

        newBase := []interface{}{}
        for _, item := range base {
            expanded := expandInterface(item)
            // fmt.Println("Expanded:", expanded)

            // if expanded != nil {
            //     switch v := expanded.(type) {
            //     case []interface{}:
            //         for _, element := range v {
            //             newBase = append(newBase, element)
            //         }
            //     default:
            //         newBase = append(newBase, v)
            //     }
            // }

            if expanded != nil {
                newBase = append(newBase, expanded)
            }
            // fmt.Println("New Base:", printInterface(base))
        }
        oldBase = base
        base = newBase
    }
}

func checkInterface(a interface{}) (bool, int, []interface{}) {
    elementsChecked := 0
    switch res := a.(type) {
    case []interface{}:
        output := false
        outputInterface := []interface{}{}
        for _, v := range res {
            checker, count, outputToAdd := checkConsistsOfBasicElements(v)
            elementsChecked += count
            if checker {
                output = true
                outputInterface = append(outputInterface, outputToAdd)
            }
        }
        
        if output {
            return true, elementsChecked, outputInterface
        }
        return false, elementsChecked, res
    case interface{}:
        checker, count, outputToAdd := checkConsistsOfBasicElements(res)
        return checker, count, []interface{}{outputToAdd}
    default:
        return false, 0, nil
    }
}

func checkConsistsOfBasicElements(item interface{}) (bool, int, interface{}) {
    elementsChecked := 0
    switch v := item.(type) {
    case *Element:
        return isBasicElement(v), 1, v
    case []*Element:
        output := false
        outputInterface := []interface{}{}
        for _, el := range v {
            elementsChecked++
            if isBasicElement(el) {
                output = true
                outputInterface = append(outputInterface, el)
            }
        }

        if output {
            return true, elementsChecked, outputInterface
        }
        return false, elementsChecked, v
    case []interface{}:
        output := false
        outputInterface := []interface{}{}
        for _, x := range v {
            checker, count, outputToAdd := checkConsistsOfBasicElements(x)
            elementsChecked += count
            if checker {
                output = true
                outputInterface = append(outputInterface, outputToAdd)
            }
        }

        if output {
            return true, elementsChecked, outputInterface
        }
        return false, elementsChecked, v
    default:
        return false, 0, nil
    }
}

func expandInterface(item interface{}) interface{} {
    switch v := item.(type) {
    case []interface{}:
        result := []interface{}{}
        for _, el := range v {
            expanded := expandInterface(el)
            if expanded != nil {
                result = append(result, expanded)
            }
        }
        return result
    case []*Element:
        result := []interface{}{}
        for _, element := range v {
            expanded := expandInterface(element)
            if expanded != nil {
                result = append(result, expanded)
            }
        }
        return result
    case *Element:
        if isBasicElement(v) {
            return v
        } else {
            return convertToInterfaceSlice(expandElement(v))
        }

    default:
        return nil
    }
}

func convertToInterfaceSlice(data [][]*Element) []interface{} {
    result := make([]interface{}, len(data))
    for i, elementSlice := range data {
        result[i] = elementSlice
    }
    return result
}

func expandElement(element *Element) [][]*Element {
    combinations := elementToCombinations(element)
    // fmt.Println("Combinations Length:", len(combinations))

    var output [][]*Element
    for _, combination := range combinations {
        if combination.Left.Tier < element.Tier && combination.Right.Tier < element.Tier {
            output = append(output, combinationToElements(combination))
        } else {
            // fmt.Println("kombinasi tidak valid")
        }
    }

    // fmt.Println("Output: %v", output)
    return output;
}

// Fungsi untuk mengecek apakah elemen adalah elemen dasar
func isBasicElement(element *Element) bool {
    return element.Name == "Water" || element.Name == "Earth" || element.Name == "Fire" || element.Name == "Air"
}

// elementToCombinations mengubah array elemen menjadi [kombinasi1, kombinasi2, kombinasi3]
func elementToCombinations(element *Element) []*Combination {
    var combinations []*Combination
    for _, combo := range element.Combinations {
        combinations = append(combinations, combo)
    }
    return combinations
}

// combinationToElements mengubah kombinasi1 menjadi [elemenpenyusun1, elemenpenyusun2]
func combinationToElements(combination *Combination) []*Element {
    return []*Element{combination.Left, combination.Right}
}

func compareInterfaceSlices(slice1, slice2 interface{}) int {
    // fmt.Println("Slice1:", printInterface(slice1), "Slice2:", printInterface(slice2))
    output := 0
    if v1, ok := slice1.([]interface{}); ok {
        if v2, ok := slice2.([]interface{}); ok {
            for i, _ := range v1 {
                output += compareInterfaceSlices(v1[i], v2[i])
            }
        }
        if v2, ok := slice2.([]*Element); ok {
            for i, _ := range v1 {
                output += compareInterfaceSlices(v1[i], v2[i])
            }
        }
    }

    if v1, ok := slice1.([]*Element); ok {
        if v2, ok := slice2.([]interface{}); ok {
            for i, _ := range v1 {
                output += compareInterfaceSlices(v1[i], v2[i])
            }
        }
        if v2, ok := slice2.([]*Element); ok {
            for i, _ := range v1 {
                output += compareInterfaceSlices(v1[i], v2[i])
            }
        }
    }
    
    if v1, ok := slice1.(*Element); ok {
        if v2, ok := slice2.(*Element); ok {
            if v1.Name == v2.Name {
                output++
            }
        }
    }
    return output;
}

func printInterface(item interface{}) string {
    switch v := item.(type) {
    case *Element:
        return v.Name
    case []*Element:
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
    // case [][]*Element:
    //     names := "["
    //     for i, elementSlice := range v {
    //         names += printInterface(elementSlice)
    //         if i < len(v)-1 {
    //             names += ", "
    //         }
    //     }
    //     names += "]"
    //     return names
    default:
        return fmt.Sprintf("%v", v) // Handle other types
    }
}