package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"
)

type CombinedResponse struct {
    Graphs []GraphResponse `json:"graphs"`
    Stats  StatsResponse   `json:"stats"`
}

type StatsResponse struct {
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
        maxStr := r.URL.Query().Get("max")
        if maxStr == "" {
            maxStr = "1"
        }
        numRecipes, err := strconv.Atoi(maxStr)
        if err != nil {
            http.Error(w, "Invalid max parameter", http.StatusBadRequest)
            return
        }

        startTime := time.Now()
        var chains []([]Recipe)
        var nodesVisited int
        
        var response CombinedResponse

        switch alg {
        case "bfs":
            chains, nodesVisited = bfs(root, idx, tiers, numRecipes)
            response.Graphs = buildMultipleTrees(root, chains)
        case "dfs":
            if mode == "shortest" {
                // Call findShortest function here
            } else {
                // Call findAll function here
            }
        }

        elapsedMs := int(time.Since(startTime).Milliseconds())

        response.Stats = StatsResponse{
            SearchTime:   elapsedMs,
            NodesVisited: nodesVisited,
            RecipesFound: len(chains),
        }

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