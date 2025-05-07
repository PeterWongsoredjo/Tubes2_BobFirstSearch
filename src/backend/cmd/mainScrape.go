package main

import (
    "log"
	"path/filepath"
	"runtime"
    "encoding/json"
    "fmt"
    "os"
    "github.com/BobKunanda/Tubes2_BobFirstSearch/src/backend/scraping"
)

func getProjectRoot() string {
	// Get the path to this current file (main.go)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get current file path")
	}

	// Go up 3 levels from cmd/ to reach project root
	// src/backend/cmd -> src/backend -> src -> project root
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	return filepath.Join(projectRoot, "backend", "configs")
}

func main() {
    fmt.Println("Scraping element list...")

    elements, err := scraping.ScrapeElementList()
    if err != nil {
        panic(fmt.Errorf("failed to scrape element list: %w", err))
    }

    fmt.Printf("Found %d elements.\n", len(elements))

    fmt.Println("Scraping recipes for each element...")

    recipes, err := scraping.ScrapeAllRecipes(elements)
    if err != nil {
        panic(fmt.Errorf("failed to scrape recipes: %w", err))
    }

    fmt.Printf("Successfully scraped %d recipes!\n", len(recipes))

    fmt.Println("Saving recipes to file...")

	configsDir := getProjectRoot()
	outputPath := filepath.Join(configsDir, "recipes.json")

	err = saveRecipesToFile(outputPath, recipes)
	if err != nil {
		panic(fmt.Errorf("failed to save recipes to file: %w", err))
	}

    fmt.Println("Done! Recipes saved to recipes.json ✅")
}

// saveRecipesToFile saves the recipes to a JSON file
func saveRecipesToFile(filename string, recipes []scraping.Recipe) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(recipes)
}

// how to run this code:
// go run cmd/mainScrape.go (tergantung dari mana run nya sama harus dalam 1 folder sm go.mod)