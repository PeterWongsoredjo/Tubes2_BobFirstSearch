package main

import (
	"encoding/json"
	//"fmt"
	//"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MultipleGraphsResponse struct {
	Graphs []GraphResponse `json:"graphs"`
	Stats  StatsRespone    `json:"stats"`
}

type StatsRespone struct {
	SearchTime   int `json:"searchTime"`
	NodesVisited int `json:"nodesVisited"`
	RecipesFound int `json:"recipesFound"`
}

func treeHandler(idx map[string][][]string, tiers map[string]int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		root := r.URL.Query().Get("root")
		if root == "" {
			http.Error(w, "missing ?root=...", 400)
			return
		}
		mode := r.URL.Query().Get("mode")
		if mode == "" {
			http.Error(w, "missing ?mode=...", 400)
			return
		}
		alg := r.URL.Query().Get("alg")
		maxRecipes := r.URL.Query().Get("maxRecipes")
		if maxRecipes == "" {
			maxRecipes = "1"
		}
		numRecipes, err := strconv.Atoi(maxRecipes)
		if err != nil {
			http.Error(w, "Invalid maxRecipes parameter", http.StatusBadRequest)
			return
		}

		startTime := time.Now()
		var response MultipleGraphsResponse
		var stats StatsRespone

		switch alg {
		case "bfs":
			chains, nodesVisited := bfs(root, idx, tiers, numRecipes)
			stats.NodesVisited = nodesVisited
			response.Graphs = buildMultipleTrees(root, chains)
		case "dfs":
			recipes, err := loadRecipes("configs/recipes.json")
			if err != nil {
				http.Error(w, "Error loading recipes", http.StatusInternalServerError)
				return
			}
			buildRecipeMap(recipes)
			elementMap = tiers
			built := map[string]bool{}
			solutions, nodesVisited := dfsAll(root, built, numRecipes, "")
			stats.NodesVisited = nodesVisited

			var graphs []GraphResponse
			for i, solution := range solutions {
				if i >= numRecipes {
					break
				}
				tree := buildTrueTreeFromDFS(root, solution, i)
				graphs = append(graphs, tree)
			}

			response.Graphs = graphs
		}

		elapsedMs := int(time.Since(startTime).Milliseconds())

		stats.SearchTime = elapsedMs
		stats.RecipesFound = len(response.Graphs)

		response.Stats = stats

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Encode error: %v\n", err)
		}
	}
}

func main() {
	recipes, err := loadRecipes("configs/recipes.json")
	if err != nil {
		log.Fatalf("loadRecipes: %v", err)
	}
	tiers, err := loadTiers("configs/tiers.json")
	if err != nil {
		log.Fatalf("loadTiers: %v", err)
	}
	idx := buildIndex(recipes)
	http.HandleFunc("/api/tree", treeHandler(idx, tiers))
	addr := ":8080"
	log.Printf("Listening on %s…\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
