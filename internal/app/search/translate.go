package search

import (
	"fmt"
)

// Struktur output graph (sama seperti sebelumnya)
type GraphNode struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

type GraphEdge struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// Fungsi utama untuk memulai translasi
func TranslateOutputPathToGraph(
	dfsTopLevelRecipe interface{}, // Hasil dari findCombinationRouteDFSConcurrentModified
	targetElementName string, // Nama elemen target awal
	// allElements map[string]*scraper.Element, // Tidak dibutuhkan lagi di sini jika nama sudah ada di output DFS
) (graphNodes []GraphNode, graphEdges []GraphEdge, err error) {

	var nodes []GraphNode
	var edges []GraphEdge
	nodeIDCounter := 0

	// Buat node untuk elemen target
	rootNodeID := nodeIDCounter
	nodeIDCounter++
	nodes = append(nodes, GraphNode{ID: rootNodeID, Label: targetElementName})

	// Mulai proses rekursif untuk resepnya
	err = convertModifiedToGraphRecursive(
		dfsTopLevelRecipe,
		rootNodeID, // ID dari elemen yang sedang dibuat (target)
		&nodes,
		&edges,
		&nodeIDCounter,
	)

	if err != nil {
		return nil, nil, err
	}
	return nodes, edges, nil
}

// Fungsi rekursif inti yang lebih sederhana
func convertModifiedToGraphRecursive(
	dfsRecipePart interface{}, // Bagian dari resep DFS yang sedang diproses
	parentGraphID int, // ID node graph dari elemen yang *dibuat* oleh resep ini
	nodes *[]GraphNode,
	edges *[]GraphEdge,
	nextNodeID *int,
	// allElements map[string]*scraper.Element, // Mungkin dibutuhkan lagi jika nama tidak selalu ada
) error {

	switch partData := dfsRecipePart.(type) {
	case string:
		// Ini adalah NAMA elemen dasar yang menjadi bahan.
		ingredientNodeID := *nextNodeID
		*nextNodeID++
		*nodes = append(*nodes, GraphNode{ID: ingredientNodeID, Label: partData})
		*edges = append(*edges, GraphEdge{From: parentGraphID, To: ingredientNodeID})
		return nil

	case []interface{}:
		// Ini adalah daftar bahan (seharusnya 2)
		if len(partData) != 2 {
			return fmt.Errorf("recipe part has %d items, expected 2: %v", len(partData), partData)
		}

		for _, item := range partData {
			// PERUBAHAN DI SINI: item bisa string atau map
			switch ingredientData := item.(type) {
			case string:
				// Ini adalah NAMA elemen dasar
				ingredientNodeID := *nextNodeID
				*nextNodeID++
				*nodes = append(*nodes, GraphNode{ID: ingredientNodeID, Label: ingredientData})
				*edges = append(*edges, GraphEdge{From: ingredientNodeID, To: parentGraphID})
				// Tidak ada resep lebih lanjut untuk diproses secara rekursif

			case map[string]interface{}:
				// Ini adalah elemen kompleks
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

				err := convertModifiedToGraphRecursive(ingredientRecipe, ingredientNodeID, nodes, edges, nextNodeID /*, allElements */)
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
// Fungsi utama untuk memulai translasi
// `targetInitialRecipes` adalah hasil dari findMultipleRouteDFS untuk targetElementName
// yaitu, []interface{} yang berisi resep-resep alternatif untuk target utama.
func TranslateMultiRecipeOutputToGraph(
	targetElementName string,
	targetInitialRecipes []interface{}, // Output dari findMultipleRouteDFS Anda
) (graphNodes []GraphNode, graphEdges []GraphEdge, err error) {

	var nodes []GraphNode
	var edges []GraphEdge
	nodeIDCounter := 0
	// Untuk scraper.isBasicElementByName jika diperlukan (opsional, bisa dihapus jika tidak dipakai)
	// allElements, _ := scraper.LoadElementsFromFile()

	// 1. Buat node untuk elemen target akhir
	targetNodeID := nodeIDCounter
	nodeIDCounter++
	nodes = append(nodes, GraphNode{ID: targetNodeID, Label: targetElementName})

	// 2. Proses resep-resep untuk target utama.
	err = processRecipeAlternativesInternal(
		targetInitialRecipes, // Ini adalah daftar resep alternatif untuk targetElementName
		targetElementName,    // Nama elemen yang sedang dibuat
		targetNodeID,         // ID node dari elemen yang sedang dibuat (menjadi TUJUAN)
		&nodes,
		&edges,
		&nodeIDCounter,
		targetElementName, // Konteks awal
		// allElements, // Teruskan jika isBasicElementByName membutuhkannya
	)

	if err != nil {
		return nil, nil, err
	}
	return nodes, edges, nil
}

// processRecipeAlternativesInternal menangani daftar resep alternatif untuk sebuah elemen.
// `idHasilCraftIni` adalah ID dari node elemen yang resepnya sedang kita proses (misal, "Aquarium" atau "Steam").
// Bahan-bahan akan memiliki edge 'From' mereka MENUNJUK ke 'idHasilCraftIni' (atau ke node '+' yang menunjuk ke sini).
func processRecipeAlternativesInternal(
	alternativeRecipesForElement []interface{}, // Daftar resep alternatif (output dari DFS untuk suatu elemen)
	nameHasilCraftIni string, // Nama elemen yang akan dihasilkan oleh alternatif-alternatif ini
	idHasilCraftIni int, // ID node dari `nameHasilCraftIni`
	nodes *[]GraphNode,
	edges *[]GraphEdge,
	nextNodeID *int,
	context string,
	// allElements map[string]*scraper.Element, // Opsional
) error {

	if len(alternativeRecipesForElement) == 0 {
		// Jika elemen ini dasar, output DFS Anda adalah `[]interface{}{"NamaDasar"}`.
		// Jika kompleks tapi tidak ada resep, ini adalah masalah.
		// Kita asumsikan DFS tidak akan mengembalikan list kosong untuk elemen kompleks yang punya resep.
		// Jika ini adalah elemen dasar, node-nya sudah dibuat oleh pemanggil.
		return nil
	}

	// Cek apakah ini adalah representasi elemen dasar dari DFS Anda (yaitu, slice berisi satu string)
	if len(alternativeRecipesForElement) == 1 {
		if _, ok := alternativeRecipesForElement[0].(string); ok {
			// Ini adalah elemen dasar. Node untuk `nameHasilCraftIni` sudah dibuat.
			// Tidak ada edge "bahan" yang perlu dibuat *dari sini* karena ini adalah hasil itu sendiri.
			// Bahan akan menunjuk ke sini ketika elemen ini digunakan sebagai HASIL dalam resep lain.
			return nil
		}
	}

	var destinationNodeForIngredients int // Ke mana bahan-bahan dari setiap resep akan membuat edge 'FROM'

	if len(alternativeRecipesForElement) > 1 {
		// Ada beberapa cara untuk membuat `nameHasilCraftIni`. Buat node "+".
		plusNodeID := *nextNodeID
		(*nextNodeID)++
		*nodes = append(*nodes, GraphNode{ID: plusNodeID, Label: "+"})
		// Edge DARI node "+" INI KE `idHasilCraftIni` (hasilnya).
		*edges = append(*edges, GraphEdge{From: plusNodeID, To: idHasilCraftIni})
		destinationNodeForIngredients = plusNodeID // Bahan akan menunjuk ke node "+" ini
	} else {
		// Hanya satu resep (dan itu bukan string dasar, sudah ditangani di atas),
		// Bahan akan langsung menunjuk ke `idHasilCraftIni`.
		destinationNodeForIngredients = idHasilCraftIni
	}

	// Untuk setiap resep alternatif (`oneSpecificRecipePath`) dalam `alternativeRecipesForElement`:
	for i, oneSpecificRecipePath := range alternativeRecipesForElement {
		// `oneSpecificRecipePath` adalah `[]interface{}{ partKiri, partKanan }`.
		recipePathActualIngredients, ok := oneSpecificRecipePath.([]interface{})
		if !ok {
			return fmt.Errorf("resep alternatif %d untuk '%s' bukan []interface{}: %T (konteks: %s)", i+1, nameHasilCraftIni, oneSpecificRecipePath, context)
		}

		// Fungsi ini akan memproses bahan-bahan dari `recipePathActualIngredients`
		// dan menghubungkannya (From) ke `destinationNodeForIngredients` (To).
		err := processSingleRecipePathIngredientsInternal(
			recipePathActualIngredients,   // Ini adalah [bahan1, bahan2]
			destinationNodeForIngredients, // Node yang dihasilkan oleh resep ini (bisa '+' atau elemen langsung)
			nodes,
			edges,
			nextNodeID,
			fmt.Sprintf("%s_alt%d", context, i+1),
			// allElements, // Opsional
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// processSingleRecipePathIngredientsInternal memproses bahan-bahan dari SATU resep spesifik.
// `idRecipeOutputNode` adalah ID dari node yang akan DIHASILKAN oleh `actualIngredientsList` ini
// (bisa node '+' atau node elemen sebenarnya). Bahan-bahan akan menunjuk (From) ke sini (To).
func processSingleRecipePathIngredientsInternal(
	actualIngredientsList []interface{}, // Berisi [partKiri, partKanan] dari SATU cara pembuatan
	idRecipeOutputNode int, // Node yang dihasilkan oleh resep ini
	nodes *[]GraphNode,
	edges *[]GraphEdge,
	nextNodeID *int,
	context string,
	// allElements map[string]*scraper.Element, // Opsional
) error {
	if len(actualIngredientsList) != 2 {
		return fmt.Errorf("daftar bahan untuk konteks '%s' tidak berjumlah 2: %d", context, len(actualIngredientsList))
	}

	for itemIdx, ingredientPart := range actualIngredientsList {
		// `ingredientPart` bisa string (nama elemen dasar)
		// atau map[string]interface{}{"name": "NamaBahanKompleks", "recipe": ResepBahanKompleks}

		switch actualIngredientData := ingredientPart.(type) {
		case string: // Bahan ini adalah elemen dasar
			ingredientName := actualIngredientData
			baseIngredientNodeID := *nextNodeID
			(*nextNodeID)++
			*nodes = append(*nodes, GraphNode{ID: baseIngredientNodeID, Label: ingredientName})
			// Edge DARI bahan dasar INI (From) KE node yang dihasilkan oleh resep (To: idRecipeOutputNode)
			*edges = append(*edges, GraphEdge{From: baseIngredientNodeID, To: idRecipeOutputNode})

		case map[string]interface{}: // Bahan ini adalah elemen kompleks
			ingredientMapName, nameOk := actualIngredientData["name"].(string)
			// `ingredientMapRecipe` adalah output DFS untuk `ingredientMapName`
			ingredientMapRecipe, recipeOk := actualIngredientData["recipe"]

			if !nameOk || !recipeOk {
				return fmt.Errorf("format map bahan tidak valid pada konteks '%s', item %d: %v", context, itemIdx, actualIngredientData)
			}

			// Buat node untuk bahan kompleks ini (ini adalah "hasil antara" yang menjadi bahan)
			nodeForThisIngredientID := *nextNodeID
			(*nextNodeID)++
			*nodes = append(*nodes, GraphNode{ID: nodeForThisIngredientID, Label: ingredientMapName})

			// Edge DARI bahan (hasil antara) INI (From) KE node yang dihasilkan oleh resep saat ini (To: idRecipeOutputNode)
			*edges = append(*edges, GraphEdge{From: nodeForThisIngredientID, To: idRecipeOutputNode})

			// Sekarang, proses bagaimana `nodeForThisIngredientID` (yaitu `ingredientMapName`) ini dibuat.
			// `ingredientMapRecipe` adalah hasil dari `findMultipleRouteDFS(ingredientMapName, ...)`
			// Jadi, ini bisa jadi `[]interface{}{"NamaDasarJikaDasar}` atau `[]interface{}{resepA, resepB, ...}`

			// Kita perlu memastikan `ingredientMapRecipe` adalah `[]interface{}` sebelum meneruskannya
			// ke `processRecipeAlternativesInternal`. Output DFS Anda untuk dasar adalah `[]interface{}{"NamaDasar"}`.
			// Output DFS Anda untuk kompleks adalah `[]interface{}{ resepAlt1, resepAlt2 ... }`.
			// Jadi, `ingredientMapRecipe` akan selalu `[]interface{}`.

			recipesForThisIngredient, ok := ingredientMapRecipe.([]interface{})
			if !ok {
				// Ini seharusnya tidak terjadi jika output DFS Anda konsisten.
				return fmt.Errorf("resep dalam map untuk bahan '%s' bukan []interface{}: %T (konteks %s)", ingredientMapName, ingredientMapRecipe, context)
			}

			// `recipesForThisIngredient` adalah daftar resep alternatif untuk `ingredientMapName`.
			// Panggil `processRecipeAlternativesInternal` untuknya.
			// Hasil dari alternatif-alternatif ini akan menjadi `nodeForThisIngredientID`.
			err := processRecipeAlternativesInternal(
				recipesForThisIngredient, // Daftar resep alternatif untuk bahan ini
				ingredientMapName,        // Nama elemen yang sedang dibuat (bahan perantara ini)
				nodeForThisIngredientID,  // ID node dari elemen yang sedang dibuat (bahan ini)
				nodes,
				edges,
				nextNodeID,
				fmt.Sprintf("%s_ing_%s", context, ingredientMapName),
				// allElements, // Opsional
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
