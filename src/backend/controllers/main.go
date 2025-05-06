package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

// Recipe represents the scraped data model
type Recipe struct {
    Result     string   `json:"Result"`
    Components []string `json:"Components"`
}

// In-memory data store (Model)
var recipes []Recipe

// loadRecipes reads the JSON file into the recipes slice
func loadRecipes(path string) error {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, &recipes)
}

// Handlers (Controller)

// GetAllRecipes returns the full list of recipes
func GetAllRecipes(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(recipes)
}

// GetRecipe returns a single recipe by its Result field
func GetRecipe(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    target := vars["result"]
    for _, rec := range recipes {
        if rec.Result == target {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(rec)
            return
        }
    }
    http.Error(w, "Recipe not found", http.StatusNotFound)
}

// SearchRecipes filters recipes by a query parameter 'q'
func SearchRecipes(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")
    var matches []Recipe
    for _, rec := range recipes {
        if rec.Result == q {
            matches = append(matches, rec)
        }
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(matches)
}

func main() {
    // Initialize data
    if err := loadRecipes("recipes.json"); err != nil {
        log.Fatalf("Failed to load recipes: %v", err)
    }

    // Setup router (View + Controller wiring)
    router := mux.NewRouter()
    router.HandleFunc("/recipes", GetAllRecipes).Methods("GET")
    router.HandleFunc("/recipes/{result}", GetRecipe).Methods("GET")
    router.HandleFunc("/search", SearchRecipes).Methods("GET")

    // Start HTTP server
    log.Println("Backend server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
