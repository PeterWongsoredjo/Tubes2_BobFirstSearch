package main

import (
	//"encoding/json"
	//"fmt"
	//"log"
	//"os"
	"fmt"
	"strings"
	"sort"
)

type ChainNode struct {
	Recipe Recipe
	Parent *ChainNode
}

type QueueItem struct {
	Elems   []string
	Chain   *ChainNode
	Depth   int
	Visited map[string]bool
}

func copyMap(original map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}

func deduplicateChain(chain []Recipe) []Recipe {
	seen := make(map[string]bool)
	deduped := []Recipe{}

	for _, r := range chain {
		if seen[r.Result] {
			continue
		}
		seen[r.Result] = true
		deduped = append(deduped, r)
	}
	return deduped
}

func buildChainList(node *ChainNode) []Recipe {
	var chain []Recipe
	for node != nil {
		chain = append([]Recipe{node.Recipe}, chain...)
		node = node.Parent
	}
	return chain
}

func chainKey(chain []Recipe) string {
	var parts []string
	for _, r := range chain {
		parts = append(parts, fmt.Sprintf("%s=%s+%s", r.Result, r.Components[0], r.Components[1]))
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

func bfs(target string, idx map[string][][]string, tiers map[string]int, limit int) []([]Recipe) {
	queue := []QueueItem{{
		Elems:   []string{target},
		Chain:   nil,
		Depth:   0,
		Visited: map[string]bool{},
	}}
	solutions := []([]Recipe){}
	seenChains := map[string]bool{}

	for depth := 0; len(queue) > 0 && depth < 30; depth++ {
		nextQueue := []QueueItem{}

		for len(queue) > 0 {
			item := queue[0]
			queue = queue[1:]

			if len(item.Elems) == 0 {
				chainList := buildChainList(item.Chain)
				if isFullyResolved(chainList, nil) {
					key := chainKey(chainList)
					if !seenChains[key] {
						seenChains[key] = true
						solutions = append(solutions, deduplicateChain(chainList))
						if len(solutions) >= limit {
							return solutions
						}
					}
				}
				continue
			}

			elem := item.Elems[0]
			rest := item.Elems[1:]

			if baseElements[elem] || item.Visited[elem] {
				queue = append(queue, QueueItem{
					Elems:   rest,
					Chain:   item.Chain,
					Depth:   item.Depth,
					Visited: copyMap(item.Visited),
				})
				continue
			}

			recipes := idx[elem]
			for _, comps := range recipes {
				c1, c2 := comps[0], comps[1]
				tTier := tiers[elem]
				if tiers[c1] > tTier || tiers[c2] > tTier {
					continue
				}

				node := &ChainNode{
					Recipe: Recipe{
						Result:     elem,
						Components: []string{c1, c2},
					},
					Parent: item.Chain,
				}

				newElems := append([]string{}, rest...)
				if !item.Visited[c1] && !baseElements[c1] {
					newElems = append(newElems, c1)
				}
				if !item.Visited[c2] && !baseElements[c2] {
					newElems = append(newElems, c2)
				}

				newVisited := copyMap(item.Visited)
				newVisited[elem] = true

				nextQueue = append(nextQueue, QueueItem{
					Elems:   newElems,
					Chain:   node,
					Depth:   depth + 1,
					Visited: newVisited,
				})
			}
		}

		queue = nextQueue
	}
	return solutions
}

func isFullyResolved(chain []Recipe, _ map[string]bool) bool {
	fmt.Println(chain)
	resolved := make(map[string]bool)
	for _, r := range chain {
		resolved[r.Result] = true
	}

	for _, r := range chain {
		for _, comp := range r.Components {
			if baseElements[comp] {
				continue
			}
			if resolved[comp] {
				continue
			}
			return false
		}
	}
	return true
}

func collectEdgesFromChain(chain []Recipe) [][2]string {
	var pairs [][2]string
	for _, step := range chain {
		parent := step.Result
		for _, child := range step.Components {
			pairs = append(pairs, [2]string{parent, child})
		}
	}
	return pairs
}

func buildTrueTree(root string, pairs [][2]string) GraphResponse {
	idOf := map[string]int{}
	parentOf := map[string]int{}
	var nodes []Node
	var edges []Edge
	nextID := 1

	elementInstances := make(map[string][]int)

	elementUseCount := make(map[string]int)

	idOf[root] = nextID
	nextID++

	for _, p := range pairs {
		parentLabel := p[0]
		childLabel := p[1]

		parentID := idOf[parentLabel]
		if parentLabel != root {
			if elementUseCount[parentLabel] > 0 && len(elementInstances[parentLabel]) > 0 {
				instances := elementInstances[parentLabel]
				parentID = instances[len(instances)-1]
			}
		}

		elementUseCount[childLabel]++
		if baseElements[childLabel] || (elementUseCount[childLabel] > 1 && childLabel != root) {
			childID := nextID
			nextID++

			elementInstances[childLabel] = append(elementInstances[childLabel], childID)

			edges = append(edges, Edge{
				From: parentID,
				To:   childID,
			})

			instanceKey := fmt.Sprintf("%s_%d", childLabel, childID)
			parentOf[instanceKey] = parentID
		} else {
			if _, ok := idOf[childLabel]; !ok {
				idOf[childLabel] = nextID
				nextID++
				parentOf[childLabel] = parentID

				elementInstances[childLabel] = append(elementInstances[childLabel], idOf[childLabel])
			}

			edges = append(edges, Edge{
				From: parentID,
				To:   idOf[childLabel],
			})
		}
	}

	for label, id := range idOf {
		node := Node{ID: id, Label: label}
		if pid, ok := parentOf[label]; ok {
			node.Parent = pid
		}
		nodes = append(nodes, node)
	}

	for elemLabel, instances := range elementInstances {
		for i, id := range instances {
			if i == 0 && idOf[elemLabel] == id {
				continue
			}

			instanceKey := fmt.Sprintf("%s_%d", elemLabel, id)
			node := Node{
				ID:     id,
				Label:  elemLabel,
				Parent: parentOf[instanceKey],
			}
			nodes = append(nodes, node)
		}
	}

	return GraphResponse{Nodes: nodes, Edges: edges}
}

/*func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <RootElement>\n", os.Args[0])
        os.Exit(1)
    }
    root := os.Args[1]

    // load data
    recipes, err := loadRecipes("configs/recipes.json")
    if err != nil {
        log.Fatalf("loadRecipes: %v", err)
    }
    tiers, err := loadTiers("configs/tiers.json")
    if err != nil {
        log.Fatalf("loadTiers: %v", err)
    }
    idx := buildIndex(recipes)

    // find & print the BFS-shortest chain
    // find & print the BFS-shortest chain
        chains := bfs(root, idx, tiers, 1)
    if len(chains) == 0 {
        fmt.Printf("No BFS-shortest chain found for %q\n", root)
        return
    }

    chain := chains[0]
    fmt.Printf("BFS-shortest chain for %q:\n", root)
    for _, step := range chain {
        fmt.Printf("  %s = %s + %s\n", step.Result, step.Components[0], step.Components[1])
    }

    // build and print the full tree JSON from one BFS chain
    pairs := collectEdgesFromChain(chain)
    tree := buildTrueTree(root, pairs)
    b, _ := json.MarshalIndent(tree, "", "  ")
    fmt.Println("\nTree JSON:")
    fmt.Println(string(b))

}*/
