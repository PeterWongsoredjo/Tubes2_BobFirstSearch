package main

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

func printUsage() {
    fmt.Println("BFS Algorithm Driver")
    fmt.Println("Usage: go run BFSDriver.go <target_element> [limit]")
    fmt.Println("Example: go run BFSDriver.go Human 5")
}

func printRecipeChain(chain []Recipe) {
    fmt.Println("Recipe Chain:")
    fmt.Println("-------------")
    for i, recipe := range chain {
        fmt.Printf("%d. %s = %s + %s\n", i+1, recipe.Result, recipe.Components[0], recipe.Components[1])
    }
    fmt.Println("-------------")
}

func debugBFS(target string, limit int) {
    // Load necessary data
    fmt.Printf("Loading recipes and tiers...\n")
    
    recipes, err := loadRecipes("configs/recipes.json")
    if err != nil {
        fmt.Printf("Error loading recipes: %v\n", err)
        return
    }
    fmt.Printf("Loaded %d recipes.\n", len(recipes))
    
    tiers, err := loadTiers("configs/tiers.json")
    if err != nil {
        fmt.Printf("Error loading tiers: %v\n", err)
        return
    }
    fmt.Printf("Loaded tiers for %d elements.\n", len(tiers))
    
    // Build recipe index
    fmt.Printf("Building recipe index...\n")
    idx := buildIndex(recipes)
    fmt.Printf("Recipe index built with %d elements.\n", len(idx))
    
    // Check if target exists and get its tier
    tier, exists := tiers[target]
    if !exists {
        fmt.Printf("Error: Target element '%s' not found in tiers data.\n", target)
        return
    }
    fmt.Printf("Target element: %s (Tier %d)\n", target, tier)
    
    // Print recipe count for target
    if recipes, ok := idx[target]; ok {
        fmt.Printf("Found %d recipes for %s\n", len(recipes), target)
    } else {
        fmt.Printf("Warning: No recipes found for %s\n", target)
    }
    
    // Start BFS
    fmt.Printf("\nStarting BFS search for %s (limit: %d recipes)...\n", target, limit)
    startTime := time.Now()
    
    chains, nodesVisited := bfs(target, idx, tiers, limit)
    
    elapsedTime := time.Since(startTime)
    
    // Print results
    fmt.Printf("\nBFS Search Results:\n")
    fmt.Printf("-------------------\n")
    fmt.Printf("Time taken: %v\n", elapsedTime)
    fmt.Printf("Nodes visited: %d\n", nodesVisited)
    fmt.Printf("Solutions found: %d\n\n", len(chains))
    
    // Print each solution
    if len(chains) == 0 {
        fmt.Printf("No solutions found for %s\n", target)
    } else {
        for i, chain := range chains {
            fmt.Printf("Solution %d/%d:\n", i+1, len(chains))
            printRecipeChain(chain)
            
            // Verify solution
            if isFullyResolved(chain, nil) {
                fmt.Printf("✓ Solution is fully resolved.\n")
            } else {
                fmt.Printf("✗ WARNING: Solution is NOT fully resolved!\n")
            }
            
            fmt.Println()
        }
    }
    
    // Print debug tree stats
    printTreeStats(target, chains)
}

func printTreeStats(target string, chains [][]Recipe) {
    if len(chains) == 0 {
        return
    }
    
    // Count unique elements in solution
    elementsInSolution := make(map[string]bool)
    for _, recipe := range chains[0] {
        elementsInSolution[recipe.Result] = true
        for _, comp := range recipe.Components {
            elementsInSolution[comp] = true
        }
    }
    
    fmt.Printf("Solution Stats:\n")
    fmt.Printf("---------------\n")
    fmt.Printf("Total elements in solution: %d\n", len(elementsInSolution))
    fmt.Printf("Total recipes in solution: %d\n", len(chains[0]))
    
    // Count base elements used
    baseElementsUsed := 0
    for e := range elementsInSolution {
        if baseElements[e] {
            baseElementsUsed++
        }
    }
    fmt.Printf("Base elements used: %d\n", baseElementsUsed)
    
    // Print base elements
    fmt.Printf("Base elements: ")
    for e := range baseElements {
        if elementsInSolution[e] {
            fmt.Printf("%s ", e)
        }
    }
    fmt.Println()
}

func maina() {
    
    target := "Human"
    limit := 2 // Default limit
    
    if len(os.Args) > 2 {
        var err error
        limit, err = strconv.Atoi(os.Args[2])
        if err != nil || limit < 1 {
            fmt.Println("Invalid limit value. Using default limit of 1.")
            limit = 1
        }
    }
    
    // Run BFS with debug output
    debugBFS(target, limit)
}