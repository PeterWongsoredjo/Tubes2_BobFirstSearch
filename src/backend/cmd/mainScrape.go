package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "runtime"

    scrapeTier "github.com/BobKunanda/Tubes2_BobFirstSearch/src/backend/scrapeTier"
    "github.com/BobKunanda/Tubes2_BobFirstSearch/src/backend/scraping"
)

func getConfigDir() string {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        log.Fatal("Could not get current file path")
    }
    projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
    return filepath.Join(projectRoot, "backend", "configs")
}

func saveJSON(path string, v interface{}) {
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        log.Fatalf("mkdir %s: %v", path, err)
    }
    f, err := os.Create(path)
    if err != nil {
        log.Fatalf("create %s: %v", path, err)
    }
    defer f.Close()
    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")
    if err := enc.Encode(v); err != nil {
        log.Fatalf("encode %s: %v", path, err)
    }
}

func main() {
    cfgDir := getConfigDir()

    // 1) Scrape and save tiers.json
    fmt.Println("Scraping tiers...")
    tiers, err := scrapeTier.ScrapeTierMap()
    if err != nil {
        log.Fatalf("failed scraping tiers: %v", err)
    }
    tierMap := make(map[string]int, len(tiers))
    for _, e := range tiers {
        tierMap[e.Name] = e.Tier
    }
    tierFile := filepath.Join(cfgDir, "tiers.json")
    fmt.Printf("Saving %d tiers to %s\n", len(tierMap), tierFile)
    saveJSON(tierFile, tierMap)

    // 2) Scrape element list
    fmt.Println("Scraping element list...")
    elements, err := scraping.ScrapeElementList()
    if err != nil {
        log.Fatalf("failed scraping element list: %v", err)
    }
    fmt.Printf("Found %d elements\n", len(elements))

    // 3) Scrape all recipes
    fmt.Println("Scraping recipes...")
    recipes, err := scraping.ScrapeAllRecipes(elements)
    if err != nil {
        log.Fatalf("failed scraping recipes: %v", err)
    }
    recipeFile := filepath.Join(cfgDir, "recipes.json")
    fmt.Printf("Saving %d recipes to %s\n", len(recipes), recipeFile)
    saveJSON(recipeFile, recipes)

    fmt.Println("✅ Done. Configs written to", cfgDir)
}
