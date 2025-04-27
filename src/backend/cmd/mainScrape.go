package main

import (
    "encoding/json"
    "fmt"
    "os"
    "backend/scraping"
)

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

    err = saveRecipesToFile("../recipes.json", recipes)
    if err != nil {
        panic(fmt.Errorf("failed to save recipes to file: %w", err))
    }

    fmt.Println("Done! Recipes saved to recipes.json ✅")
}

// saveRecipesToFile saves the recipes to a JSON file
func saveRecipesToFile(filename string, recipes []scraping.Recipe) error {
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
// go run src/backend/mainScrape.go