package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type QueueItem struct {
    Elem  string
    Chain []Recipe
    Depth int
}

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

func findShortest(target string, idx map[string][][]string) ([]Recipe, bool) {
    queue := []QueueItem{{Elem: target, Chain: nil, Depth: 0}}

    for len(queue) > 0 {
        levelSize := len(queue)
        for i := 0; i < levelSize; i++ {
            item := queue[0]
            queue = queue[1:]

            recipes := idx[item.Elem]
            for _, comps := range recipes {
                c1, c2 := comps[0], comps[1]
                if baseElements[c1] && baseElements[c2] {
                    chain := append([]Recipe{}, item.Chain...)
                    chain = append(chain, Recipe{Result: item.Elem, Components: comps})
                    return chain, true
                }
                newChain := append([]Recipe{}, item.Chain...)
                newChain = append(newChain, Recipe{Result: item.Elem, Components: comps})

                if !baseElements[c1] {
                    queue = append(queue, QueueItem{Elem: c1, Chain: newChain, Depth: item.Depth + 1})
                } else {
                    newChainC1 := append([]Recipe{}, newChain...)
                    newChainC1 = append(newChainC1, Recipe{Result: c1, Components: []string{}})
                    queue = append(queue, QueueItem{Elem: c2, Chain: newChainC1, Depth: item.Depth + 1})
                }

                if !baseElements[c2] && !baseElements[c1] {
                    queue = append(queue, QueueItem{Elem: c2, Chain: newChain, Depth: item.Depth + 1})
                }
            }
        }
    }
    return nil, false
}

func printFullChain(elem string, idx map[string][][]string, printed map[string]bool, indent string) {
    if printed[elem] {
        return
    }
    chain, ok := findShortest(elem, idx)
    if !ok {
        fmt.Printf("%sno path found to %q\n", indent, elem)
        return
    }
    for _, r := range chain {
        if printed[r.Result] {
            continue
        }
        if len(r.Components) < 2 {
            fmt.Printf("%s%s is a base element\n", indent, r.Result)
        } else {
            fmt.Printf("%s%s = %s + %s\n", indent, r.Result, r.Components[0], r.Components[1])
        }
        printed[r.Result] = true
    }
    for _, r := range chain {
        for _, comp := range r.Components {
            if !baseElements[comp] {
                printFullChain(comp, idx, printed, indent+"  ")
            }
        }
    }
}

func maina() {
    // flags
    //mode := flag.String("mode", "shortest", "search mode: shortest or all")
    //maxR := flag.Int("max", 5, "max recipes to find when mode=all")
    //flag.Parse()

    recipes, err := loadRecipes("configs/recipes.json")
    if err != nil {
        fmt.Fprintf(os.Stderr, "error loading recipes: %v\n", err)
        os.Exit(1)
    }
    idx := buildIndex(recipes)

    fmt.Print("Enter target element: ")
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    target := strings.TrimSpace(input)

    printed := make(map[string]bool)
    fmt.Println("Shortest recipe chain:")
    printFullChain(target, idx, printed, "")

    /*
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
        */
}
