package main

import (
    "encoding/json"
    "log"
    "net/http"
)

func treeHandler(idx map[string][][]string, tiers map[string]int) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        root := r.URL.Query().Get("root")
        if root == "" {
            http.Error(w, "missing ?root=...", 400)
            return
        }

        // collect BFS-shortest edges
        var pairs [][2]string
        collectEdges(root, idx, tiers, map[string]bool{}, &pairs)
        log.Printf("Collected %d edges for %q\n", len(pairs), root)

        graph := buildTrueTree(root, pairs)
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        if err := json.NewEncoder(w).Encode(graph); err != nil {
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
