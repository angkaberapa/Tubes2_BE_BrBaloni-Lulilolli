func stringToElement(elements map[string]*Element, name string) *scraper.Element {
	return elements[name]
}

func isBasicElement(element string) bool {
    return element == "Water" || element == "Earth" || element == "Fire" || element == "Air"
}

func isAllBasicElement(elements []string) bool {
	for _, element := range elements {
		if !isBasicElement(element) {
			return false
		}
	}
	return true
}

func findCombinationRouteBFS(elements map[string]*Element, target string) ([][]string, int) {
	targetElement := stringToElement(elements, target)
    if (isBasicElement(targetElement.Name)) {
        return [][]string{[]string{targetElement.Name, "-", "-"}}, 1
	}

    combinationsNumber := len(targetElement.Combinations) 
    combinationsChecked := 1
    base := [][]string{}
	checker := [][]string{}
    for i := 0; i < combinationsNumber; i++ {
		combinationsChecked += 2
        if (stringToElement(elements, targetElement.Combinations[i].Left).Tier < targetElement.Tier && stringToElement(elements, targetElement.Combinations[i].Right).Tier < targetElement.Tier) {
            base = append(base, []string{targetElement.Name, targetElement.Combinations[i].Left, targetElement.Combinations[i].Right})
			checker = append(checker, []string{targetElement.Combinations[i].Left, targetElement.Combinations[i].Right})
        } else {
            fmt.Println("kombinasi tidak valid")
        }
	}

	bfsFound := false
	for i := 0; i < len(checker); i++ {
		if (isAllBasicElement(checker[i])) {
			bfsFound = true
			break
		}
	}

	for !bfsFound {
		currentCheck := checker[0]
		checker = checker[1:]

		for i := 0; i < len(currentCheck); i++ {
			addToCheck := []string{}

			if (!isBasicElement(currentCheck[i])) {
				combinationsChecked += 2

				currentElement := stringToElement(elements, currentCheck[i])
				combinationsNumber = len(currentElement.Combinations)

				for i := 0; i < combinationsNumber; i++ {
					addToCheck = append(addToCheck, currentElement.Combinations[i].Left)
					addToCheck = append(addToCheck, currentElement.Combinations[i].Right)

					leftElement := stringToElement(elements, currentElement.Combinations[i].Left)
					rightElement := stringToElement(elements, currentElement.Combinations[i].Right)

					if (leftElement.Tier < currentElement.Tier && rightElement.Tier < currentElement.Tier) {
						toAdd := []string{currentElement.Name, currentElement.Combinations[i].Left, currentElement.Combinations[i].Right}
						base = append(base, toAdd)
					}
				}
			}

			checker = append(checker, addToCheck)
		}

		for i := 0; i < len(checker); i++ {
			if (isAllBasicElement(checker[i])) {
				bfsFound = true
				break
			}
		}

		if (bfsFound) {
			basicElements := [][]string{[]string{"Fire", "-", "-"}, []string{"Water", "-", "-"}, []string{"Earth", "-", "-"}, []string{"Air", "-", "-"}}
			for i := 0; i < len(basicElements); i++ {
				base = append(base, basicElements[i])
			}
		}
	}

	return base, combinationsChecked
}


func createFixList(list [][]string) ([][]string, []string) {
	fixList := [][]string{[]string{"Water", "-", "-"}, []string{"Earth", "-", "-"}, 
						  []string{"Fire", "-", "-"}, []string{"Air", "-", "-"}}
	exists := []string{"Water", "Earth", "Fire", "Air"}
	for i := len(list) - 4; i >= 0; i-- {
		toBeAdded1 := false
		toBeAdded2 := false
		alreadyExists := false
		for j := 0; j < len(exists); j++ {
			if (list[i][0] == exists[j]) {
				alreadyExists = true
			}
			if (list[i][1] == exists[j]) {
				toBeAdded1 = true
			}
			if (list[i][2] == exists[j]) {
				toBeAdded2 = true
			}
		}
		if (toBeAdded1 && toBeAdded2) {
			fixList = append(fixList, list[i])
			if (!alreadyExists) {
				exists = append(exists, list[i][0])
			}
		}
	}

	output := [][]string{}
	for i := len(fixList) - 1; i >= 0; i-- {
		output = append(output, fixList[i])
	}

	return output, exists
}

func createMap(list [][]string, mapList []string) map[string][][]string {
	m := make(map[string][][]string)
	for i := 0; i < len(mapList); i++ {
		if (!isBasicElement(mapList[i])) {
			m[mapList[i]] = [][]string{}
			for j := 0; j < len(list); j++ {
				if (list[j][0] == mapList[i]) {
					m[mapList[i]] = append(m[mapList[i]], []string{list[j][1], list[j][2]})
				}
			}
		}
	}

	return m
}

func expandNode(graph map[string][][]string, name string, parentID int, ID *int) ([]GraphNode, []GraphEdge) {
	nodesToAdd := []GraphNode{}
	edgesToAdd := []GraphEdge{}

	(*ID)++
	nodesToAdd = append(nodesToAdd, GraphNode{ID: *ID, Label: name})
	edgesToAdd = append(edgesToAdd, GraphEdge{From: parentID, To: *ID})

	// fmt.Println("toExpand:", name)
	if (!isBasicElement(name)) {
		edges := graph[name]
		parentID = *ID
		if (len(edges) == 1) {
			for j := 0; j < len(edges[0]); j++ {
				nodesToAddRecursive, edgesToAddRecursive := expandNode(graph, edges[0][j], parentID, ID)
				for k := 0; k < len(nodesToAddRecursive); k++ {
					nodesToAdd = append(nodesToAdd, nodesToAddRecursive[k])
				}
				for k := 0; k < len(edgesToAddRecursive); k++ {
					edgesToAdd = append(edgesToAdd, edgesToAddRecursive[k])
				}
			}
		} else {
			for i := 0; i < len(edges); i++ {
				(*ID)++
				nodesToAdd = append(nodesToAdd, GraphNode{ID: *ID, Label: "+"})
				edgesToAdd = append(edgesToAdd, GraphEdge{From: parentID, To: *ID})
				plusID := *ID
				for j := 0; j < len(edges[i]); j++ {
					nodesToAddRecursive, edgesToAddRecursive := expandNode(graph, edges[i][j], plusID, ID)
					for k := 0; k < len(nodesToAddRecursive); k++ {
						nodesToAdd = append(nodesToAdd, nodesToAddRecursive[k])
					}
					for k := 0; k < len(edgesToAddRecursive); k++ {
						edgesToAdd = append(edgesToAdd, edgesToAddRecursive[k])
					}
				}
			}
		}
	}

	fmt.Println("nodesToAdd", nodesToAdd)
	fmt.Println("edgesToAdd", edgesToAdd)
	return nodesToAdd, edgesToAdd
}

func createBFSNode(graph map[string][][]string, list []string) ([]GraphNode, []GraphEdge) {
	ID := -1
	nodes := []GraphNode{}
	edges := []GraphEdge{}
	nodesToAdd, edgesToAdd := expandNode(graph, list[len(list) - 1], ID, &ID)
	for j := 0; j < len(nodesToAdd); j++ {
		nodes = append(nodes, nodesToAdd[j])
	}
	for j := 0; j < len(edgesToAdd); j++ {
		edges = append(edges, edgesToAdd[j])
	}

	return nodes, edges[1:]
}

func findBFS(base [][]string) ([]graphNodes, []graphEdges) {
	fixList, existsList := createFixList(base)
	m := createMap(fixList, existsList)
	nodes, edges := createBFSNode(m, existsList)

	return nodes, edges
}
func BFS(target string, maxrecipe int) ([]graphNodes, []graphEdges, error, int) { // ganti return value nya
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err, -1
	}

	results, nodeCount := findCombinationRouteBFS(elements, target)
	// masukkan kode disini
	if results == nil {
		return nil, fmt.Errorf("no valid combination found"), nodeCount
	}

	nodes, edges := findBFS(results)
	return results, nil, nodeCount
}