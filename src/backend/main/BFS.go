package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
)

// Path is a sequence of element names from a base up to some element
// Combined via recipes, recorded as pairs.
type Path struct {
    Steps []Recipe
}

// index recipes by output element
func buildIndex(recipes []Recipe) map[string][][]string {
    idx := make(map[string][][]string)
    seen := map[string]map[string]bool{}
    for _, r := range recipes {
        key := strings.Join(r.Components, "|")
        if seen[r.Result] == nil {
            seen[r.Result] = map[string]bool{}
        }
        if seen[r.Result][key] {
            continue
        }
        seen[r.Result][key] = true
        idx[r.Result] = append(idx[r.Result], r.Components)
    }
    return idx
}

// BFS for shortest single recipe chain
func findShortest(target string, idx map[string][][]string) ([]Recipe, bool) {
    type queueItem struct {
        Current string
        Chain   []Recipe
    }
    visited := map[string]bool{target: true}
    queue := []queueItem{{Current: target, Chain: nil}}

    for len(queue) > 0 {
        item := queue[0]
        queue = queue[1:]

        if baseElements[item.Current] {
            // reached a base element: return the chain reversed
            // chain holds recipes from target downwards
            // reverse order to show from base -> target
            for i, j := 0, len(item.Chain)-1; i < j; i, j = i+1, j-1 {
                item.Chain[i], item.Chain[j] = item.Chain[j], item.Chain[i]
            }
            return item.Chain, true
        }

        for _, comps := range idx[item.Current] {
            // expand each recipe
            for _, c := range comps {
                if visited[c] {
                    continue
                }
            }
            // enqueue this recipe
            newChain := append([]Recipe{}, item.Chain...)
            newChain = append(newChain, Recipe{Result: item.Current, Components: comps})
            for _, comp := range comps {
                if !visited[comp] {
                    visited[comp] = true
                    queue = append(queue, queueItem{Current: comp, Chain: newChain})
                }
            }
        }
    }
    return nil, false
}

// BFS for multiple recipe chains
func findAll(target string, idx map[string][][]string, max int) [][]Recipe {
    type queueItem struct {
        Current string
        Chain   []Recipe
    }
    results := make([][]Recipe, 0, max)
    queue := []queueItem{{Current: target, Chain: nil}}

    for len(queue) > 0 && len(results) < max {
        item := queue[0]
        queue = queue[1:]

        if baseElements[item.Current] {
            // found a full chain
            chainCopy := append([]Recipe{}, item.Chain...)
            // reverse to base->target
            for i, j := 0, len(chainCopy)-1; i < j; i, j = i+1, j-1 {
                chainCopy[i], chainCopy[j] = chainCopy[j], chainCopy[i]
            }
            results = append(results, chainCopy)
            continue
        }

        for _, comps := range idx[item.Current] {
            newChain := append([]Recipe{}, item.Chain...)
            newChain = append(newChain, Recipe{Result: item.Current, Components: comps})
            for _, comp := range comps {
                queue = append(queue, queueItem{Current: comp, Chain: newChain})
            }
        }
    }
    return results
}

func maina() {
    // flags
    mode := flag.String("mode", "shortest", "search mode: shortest or all")
    maxR := flag.Int("max", 5, "max recipes to find when mode=all")
    flag.Parse()

    recipes, err := loadRecipes("../configs/recipes.json")
    if err != nil {
        fmt.Fprintf(os.Stderr, "error loading recipes: %v\n", err)
        os.Exit(1)
    }
    idx := buildIndex(recipes)

    fmt.Print("Enter target element: ")
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    target := strings.TrimSpace(input)

    if *mode == "shortest" {
        chain, ok := findShortest(target, idx)
        if !ok {
            fmt.Printf("no path found to %q\n", target)
            os.Exit(0)
        }
        fmt.Println("Shortest recipe chain:")
        for _, r := range chain {
            fmt.Printf("%s = %s + %s\n", r.Result, r.Components[0], r.Components[1])
        }
    } else {
        chains := findAll(target, idx, *maxR)
        if len(chains) == 0 {
            fmt.Printf("no paths found to %q\n", target)
            os.Exit(0)
        }
        fmt.Printf("Found %d recipe chains (max %d):\n", len(chains), *maxR)
        for i, chain := range chains {
            fmt.Printf("Chain #%d:\n", i+1)
            for _, r := range chain {
                fmt.Printf("  %s = %s + %s\n", r.Result, r.Components[0], r.Components[1])
            }
        }
    }
}
