package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "encoding/json"
	"strconv"
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

func loadTiers(path string) (map[string]int, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var tiers map[string]int
    if err := json.NewDecoder(f).Decode(&tiers); err != nil {
        return nil, err
    }
    return tiers, nil
}


func findShortest(target string, idx map[string][][]string, tiers map[string]int) ([]Recipe, bool) {
    queue := []QueueItem{{Elem: target, Chain: nil, Depth: 0}}
	visited := make(map[string]int)
	currentDepth := 0
	var bestSolution[]Recipe
	var foundSolution bool

    for len(queue) > 0 {
        levelSize := len(queue)
		//fmt.Printf("Processing depth %d with %d items\n", currentDepth, levelSize)
        for i := 0; i < levelSize; i++ {
            item := queue[0]
            queue = queue[1:]

			if prevDepth, seen := visited[item.Elem]; seen && prevDepth <= item.Depth {
				continue
			} 

			visited[item.Elem] = item.Depth

            recipes := idx[item.Elem]
            for _, comps := range recipes {
                c1, c2 := comps[0], comps[1]

                tTier := tiers[item.Elem]
                if tiers[c1] > tTier || tiers[c2] > tTier {
                    continue
                }
                
                newChain := append([]Recipe{}, item.Chain...)
                newChain = append(newChain, Recipe{Result: item.Elem, Components: comps})

				if baseElements[c1] && baseElements[c2] {
					bestSolution = newChain
					foundSolution = true
					continue
				}

				queue = append(queue, QueueItem{Elem: c1, Chain: newChain, Depth: item.Depth + 1})
				queue = append(queue, QueueItem{Elem: c2, Chain: newChain, Depth: item.Depth + 1})
            }
        }
		if foundSolution {
			return bestSolution, true
		}
		currentDepth++
    }
    return nil, false
}

func printFullChain(elem string, idx map[string][][]string, tiers map[string]int, printed map[string]bool, indent string) {
    if printed[elem] {
        return
    }
    chain, ok := findShortest(elem, idx, tiers)
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
                printFullChain(comp, idx, tiers, printed, indent+"  ")
            }
        }
    }
}

func maina() {
    reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter target element: ")
	input, _ := reader.ReadString('\n')
	target := strings.TrimSpace(input)

	recipes, err := loadRecipes("configs/recipes.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading recipes: %v\n", err)
		os.Exit(1)
	}

	tiers, err := loadTiers("configs/tiers.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading tiers: %v\n", err)
		os.Exit(1)
	}

	idx := buildIndex(recipes)

	printed := make(map[string]bool)
	fmt.Print("Do you want the shortest recipe chain or multiple distinct chains? (shortest/multiple): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	if choice == "shortest" {
		fmt.Println("\nShortest recipe chain:")
		printFullChain(target, idx, tiers, printed, "")
	} else if choice == "multiple" {
		fmt.Print("Enter the maximum number of distinct chains to find: ")
		maxInput, _ := reader.ReadString('\n')
		maxInput = strings.TrimSpace(maxInput)
		max, err := strconv.Atoi(maxInput)
		if err != nil || max <= 0 {
			fmt.Println("Invalid number. Exiting.")
			return
		}

	}else {
		fmt.Println("Invalid choice. Exiting.")
	}

}
