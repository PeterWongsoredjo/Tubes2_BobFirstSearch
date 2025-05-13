package main

import (
	//"encoding/json"
	"fmt"
	//"os"
	"sync"
)

var recipeMap map[string][][]string
var elementMap map[string]int

type Step struct {
	Target string
	Parent string
}

func buildRecipeMap(recipes []Recipe) {
	recipeMap = make(map[string][][]string)
	for _, r := range recipes {
		recipeMap[r.Result] = append(recipeMap[r.Result], r.Components)
	}
}


func dfsAll(target string, built map[string]bool, limit int, parent string) ([][]Step, int) {
	nodesVisited := 1 // node ori

	if baseElements[target] {
		return [][]Step{{Step{Target: target, Parent: parent}}}, nodesVisited
	}

	if built[target] {
		return nil, nodesVisited
	}

	var allChains [][]Step
	seenChains := map[string]bool{}

	var wg sync.WaitGroup        
	var mu sync.Mutex            
	var leftChains, rightChains [][]Step
	var leftVisited, rightVisited int

	if componentsList, ok := recipeMap[target]; ok {
		for _, components := range componentsList {
			if len(allChains) >= limit {
				break
			}

			left, right := components[0], components[1]

			if elementMap[target] > elementMap[left] && elementMap[target] > elementMap[right] {
				if left == right && !baseElements[left] {
					leftChains, leftVisited = dfsAll(left, copyMap(built), limit, target)
					nodesVisited += leftVisited

					for _, l := range leftChains {
						if len(allChains) >= limit {
							break
						}

						newChain := append([]Step{}, l...)
						newChain = append(newChain, Step{Target: target, Parent: parent})
						allChains = append(allChains, newChain)
					}
				} else {
					wg.Add(2) 

					go func() {
						defer wg.Done() 
						leftChains, leftVisited = dfsAll(left, copyMap(built), limit, target)
						nodesVisited += leftVisited
					}()

					go func() {
						defer wg.Done()
						rightChains, rightVisited = dfsAll(right, copyMap(built), limit, target)
						nodesVisited += rightVisited
					}()

					wg.Wait()

					for _, l := range leftChains {
						for _, r := range rightChains {
							if len(allChains) >= limit {
								break
							}

							newChain := append([]Step{}, l...)
							newChain = append(newChain, r...)
							newChain = append(newChain, Step{Target: target, Parent: parent})

							sig := chainSignature(newChain)
							mu.Lock()
							if !seenChains[sig] {
								allChains = append(allChains, newChain)
								seenChains[sig] = true
							}
							mu.Unlock()
						}
					}
				}
			}
		}
	}

	// Ruby chan haiiiiiii
	built[target] = true //nani ga suki
	return allChains, nodesVisited //🍆🍆💦
}
//JOKO MINTOOOOOOOOOOOO

func buildTrueTreeFromDFS(root string, steps []Step, idx int) GraphResponse {
	reverseChain(steps) 
	if len(steps) > 0{
		steps = steps[1:]
	}
	
	nextID := (1000 * idx) + 1
	nodes := []Node{}
	edges := []Edge{}

	nodeIDs := map[string]int{}
	parentOf := map[string]string{}
	childCount := map[string]int{}

	nodeIDs[root] = nextID
	nextID++
	nodes = append(nodes, Node{ID: nodeIDs[root], Label: root}) 

	for _, step := range steps {
		nodeIDs[step.Target] = nextID
			nextID++
			nodes = append(nodes, Node{ID: nodeIDs[step.Target], Label: step.Target})

		if step.Parent != "" {
			if _, exists := nodeIDs[step.Parent]; !exists {
				nodeIDs[step.Parent] = nextID
				nextID++
				nodes = append(nodes, Node{ID: nodeIDs[step.Parent], Label: step.Parent})
			}

			parentOf[step.Target] = step.Parent
			childCount[step.Parent]++

			edges = append(edges, Edge{
				From: nodeIDs[step.Parent],
				To:   nodeIDs[step.Target],
			})
		}
	}

	for parent, count := range childCount {
        if count == 1 {
            for target, p := range parentOf {
                if p == parent {
                    duplicateID := nextID
                    nextID++
                    nodes = append(nodes, Node{ID: duplicateID, Label: target})
                    edges = append(edges, Edge{
                        From: nodeIDs[parent],
                        To:   duplicateID,
                    })
                    break
                }
            }
        }
    }

	return GraphResponse{Nodes: nodes, Edges: edges}
}

func reverseChain(chain []Step) {
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
}

func chainSignature(chain []Step) string {
	s := ""
	for _, step := range chain {
		s += fmt.Sprintf("%s>%s;", step.Target, step.Parent)
	}
	return s
}

/*

func dfs(target string, depth int, built map[string]bool) ([]Step, bool) {
	if baseElements[target] {
		return nil, true
	}

	if built[target] {
		return nil, true
	}

	if componentsList, ok := recipeMap[target]; ok {
		for _, components := range componentsList {
			left, right := components[0], components[1]
			if elementMap[target] > elementMap[left] && elementMap[target] > elementMap[right] {
				var leftChain, rightChain []Step
				var ok1, ok2 bool

				if left == right {
					leftChain, ok1 = dfs(left, depth+1, built)
					rightChain = []Step{}
					ok2 = ok1
				} else {
					leftChain, ok1 = dfs(left, depth+1, built)
					rightChain, ok2 = dfs(right, depth+1, built)
				}

				if ok1 && ok2 {
					built[target] = true

					chain := append([]Step{}, leftChain...)
					chain = append(chain, rightChain...)
					chain = append(chain, Step{Target: target})
					return chain, true
				}
			}
		}
	}
	return nil, false
}


func maidn() {
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

	buildRecipeMap(recipes)
	elementMap = tiers
	built := map[string]bool{}
	solutions, _ := dfsAll("Human", built, 1, "")



	for _, chain := range solutions {
		fmt.Println(chain)
	}

	if len(solutions) == 0 {
		fmt.Println("No valid DFS chain found.")
		return
	}

	tree := buildTrueTreeFromDFS("Human", solutions[0], 0)

	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling JSON: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Tree JSON for DFS Solution:")
	fmt.Println(string(jsonData))

}
*/