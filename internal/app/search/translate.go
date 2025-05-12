package search

import (
	"fmt"
)

type GraphNode struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type GraphEdge struct {
	From int `json:"from"`
	To   int `json:"to"`
}

func TranslateOutputPathToGraph(
	dfsTopLevelRecipe interface{},
	targetElementName string,
) (graphNodes []GraphNode, graphEdges []GraphEdge, err error) {

	var nodes []GraphNode
	var edges []GraphEdge
	nodeIDCounter := 0

	rootNodeID := nodeIDCounter
	nodeIDCounter++
	nodes = append(nodes, GraphNode{ID: rootNodeID, Label: targetElementName})

	err = convertModifiedToGraphRecursive(
		dfsTopLevelRecipe,
		rootNodeID,
		&nodes,
		&edges,
		&nodeIDCounter,
	)

	if err != nil {
		return nil, nil, err
	}
	return nodes, edges, nil
}

func convertModifiedToGraphRecursive(
	dfsRecipePart interface{},
	parentGraphID int,
	nodes *[]GraphNode,
	edges *[]GraphEdge,
	nextNodeID *int,
) error {

	switch partData := dfsRecipePart.(type) {
	case string:
		ingredientNodeID := *nextNodeID
		*nextNodeID++
		*nodes = append(*nodes, GraphNode{ID: ingredientNodeID, Label: partData})
		*edges = append(*edges, GraphEdge{From: parentGraphID, To: ingredientNodeID})
		return nil

	case []interface{}:
		if len(partData) != 2 {
			return fmt.Errorf("recipe part has %d items, expected 2: %v", len(partData), partData)
		}

		for _, item := range partData {
			switch ingredientData := item.(type) {
			case string:
				ingredientNodeID := *nextNodeID
				*nextNodeID++
				*nodes = append(*nodes, GraphNode{ID: ingredientNodeID, Label: ingredientData})
				*edges = append(*edges, GraphEdge{From: ingredientNodeID, To: parentGraphID})

			case map[string]interface{}:
				ingredientName, nameOk := ingredientData["name"].(string)
				if !nameOk {
					return fmt.Errorf("ingredient map missing or invalid 'name': %v", ingredientData)
				}

				ingredientRecipe, recipeOk := ingredientData["recipe"]
				if !recipeOk {
					return fmt.Errorf("ingredient map missing 'recipe' for '%s': %v", ingredientName, ingredientData)
				}

				ingredientNodeID := *nextNodeID
				*nextNodeID++
				*nodes = append(*nodes, GraphNode{ID: ingredientNodeID, Label: ingredientName})
				*edges = append(*edges, GraphEdge{From: ingredientNodeID, To: parentGraphID})

				err := convertModifiedToGraphRecursive(ingredientRecipe, ingredientNodeID, nodes, edges, nextNodeID)
				if err != nil {
					return fmt.Errorf("error processing recipe for ingredient '%s': %w", ingredientName, err)
				}
			default:
				return fmt.Errorf("unknown data type for ingredient item: %T: %v", item, item)
			}
		}
		return nil

	default:
		return fmt.Errorf("unknown data type in DFS recipe part: %T for parent ID %d", dfsRecipePart, parentGraphID)
	}
}

// Untuk multiple recipe
func TranslateMultiRecipeOutputToGraph(
	targetElementName string,
	targetInitialRecipes []interface{},
) (graphNodes []GraphNode, graphEdges []GraphEdge, err error) {

	var nodes []GraphNode
	var edges []GraphEdge
	nodeIDCounter := 0

	targetNodeID := nodeIDCounter
	nodeIDCounter++
	nodes = append(nodes, GraphNode{ID: targetNodeID, Label: targetElementName})

	err = processRecipeAlternativesInternal(
		targetInitialRecipes,
		targetElementName,
		targetNodeID,
		&nodes,
		&edges,
		&nodeIDCounter,
		targetElementName,
	)

	if err != nil {
		return nil, nil, err
	}
	return nodes, edges, nil
}

func processRecipeAlternativesInternal(
	alternativeRecipesForElement []interface{},
	nameHasilCraftIni string,
	idHasilCraftIni int,
	nodes *[]GraphNode,
	edges *[]GraphEdge,
	nextNodeID *int,
	context string,
) error {

	if len(alternativeRecipesForElement) == 0 {
		return nil
	}

	if len(alternativeRecipesForElement) == 1 {
		if _, ok := alternativeRecipesForElement[0].(string); ok {
			return nil
		}
	}

	var destinationNodeForIngredients int
	for i, oneSpecificRecipePath := range alternativeRecipesForElement {
		recipePathActualIngredients, ok := oneSpecificRecipePath.([]interface{})
		if !ok {
			return fmt.Errorf("resep alternatif %d untuk '%s' bukan []interface{}: %T (konteks: %s)", i+1, nameHasilCraftIni, oneSpecificRecipePath, context)
		}
		if len(alternativeRecipesForElement) > 1 {
			// buat node "+"
			plusNodeID := *nextNodeID
			(*nextNodeID)++
			*nodes = append(*nodes, GraphNode{ID: plusNodeID, Label: "+"})
			*edges = append(*edges, GraphEdge{From: plusNodeID, To: idHasilCraftIni})
			destinationNodeForIngredients = plusNodeID
		} else {
			destinationNodeForIngredients = idHasilCraftIni
		}

		err := processSingleRecipePathIngredientsInternal(
			recipePathActualIngredients,
			destinationNodeForIngredients,
			nodes,
			edges,
			nextNodeID,
			fmt.Sprintf("%s_alt%d", context, i+1),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func processSingleRecipePathIngredientsInternal(
	actualIngredientsList []interface{},
	idRecipeOutputNode int,
	nodes *[]GraphNode,
	edges *[]GraphEdge,
	nextNodeID *int,
	context string,
) error {
	if len(actualIngredientsList) != 2 {
		return fmt.Errorf("daftar bahan untuk konteks '%s' tidak berjumlah 2: %d", context, len(actualIngredientsList))
	}

	for itemIdx, ingredientPart := range actualIngredientsList {

		switch actualIngredientData := ingredientPart.(type) {
		case string:
			ingredientName := actualIngredientData
			baseIngredientNodeID := *nextNodeID
			(*nextNodeID)++
			*nodes = append(*nodes, GraphNode{ID: baseIngredientNodeID, Label: ingredientName})
			*edges = append(*edges, GraphEdge{From: baseIngredientNodeID, To: idRecipeOutputNode})

		case map[string]interface{}: // Bahan ini adalah elemen kompleks
			ingredientMapName, nameOk := actualIngredientData["name"].(string)
			ingredientMapRecipe, recipeOk := actualIngredientData["recipe"]

			if !nameOk || !recipeOk {
				return fmt.Errorf("format map bahan tidak valid pada konteks '%s', item %d: %v", context, itemIdx, actualIngredientData)
			}

			nodeForThisIngredientID := *nextNodeID
			(*nextNodeID)++
			*nodes = append(*nodes, GraphNode{ID: nodeForThisIngredientID, Label: ingredientMapName})

			*edges = append(*edges, GraphEdge{From: nodeForThisIngredientID, To: idRecipeOutputNode})

			recipesForThisIngredient, ok := ingredientMapRecipe.([]interface{})
			if !ok {
				return fmt.Errorf("resep dalam map untuk bahan '%s' bukan []interface{}: %T (konteks %s)", ingredientMapName, ingredientMapRecipe, context)
			}

			err := processRecipeAlternativesInternal(
				recipesForThisIngredient,
				ingredientMapName,
				nodeForThisIngredientID,
				nodes,
				edges,
				nextNodeID,
				fmt.Sprintf("%s_ing_%s", context, ingredientMapName),
			)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("tipe data tidak dikenal untuk bahan pada konteks '%s', item %d: %T", context, itemIdx, ingredientPart)
		}
	}
	return nil
}
