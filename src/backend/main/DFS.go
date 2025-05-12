package main

import (
	"encoding/json"
	"fmt"
	"os"
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

	if componentsList, ok := recipeMap[target]; ok {
		for _, components := range componentsList {
			if len(allChains) >= limit {
				break
			}

			left, right := components[0], components[1]

			if elementMap[target] > elementMap[left] && elementMap[target] > elementMap[right] {
				if(left == right && !baseElements[left]) {
					leftChains, leftVisited := dfsAll(left, copyMap(built), limit, target)
					nodesVisited += leftVisited

					for _, l := range leftChains {
						if len(allChains) >= limit {
							break
						}

						newChain := []Step{}

						for _, step := range l {
							newChain = append(newChain, step)
						}
						newChain = append(newChain, Step{
							Target: target,
							Parent: parent,
						})
						allChains = append(allChains, newChain)
						sig := chainSignature(newChain)
						if !seenChains[sig] {
							allChains = append(allChains, newChain)
							seenChains[sig] = true
						}
					
					}
				} else {
					leftChains, leftVisited := dfsAll(left, copyMap(built), limit, target)
					rightChains, rightVisited := dfsAll(right, copyMap(built), limit, target)
					nodesVisited += leftVisited + rightVisited

					for _, l := range leftChains {
						for _, r := range rightChains {
							if len(allChains) >= limit {
								break
							}

							newChain := []Step{}

							for _, step := range l {
								newChain = append(newChain, step)
							}

							for _, step := range r {
								newChain = append(newChain, step)
							}

							newChain = append(newChain, Step{
								Target: target,
								Parent: parent,
							})

							sig := chainSignature(newChain)
							if !seenChains[sig] {
								allChains = append(allChains, newChain)
								seenChains[sig] = true
							}
						}
					}
				}
			}
		}
	}

	built[target] = true
	return allChains, nodesVisited
}

func buildTrueTreeFromDFS(root string, steps []Step, idx int) GraphResponse {
	// Reverse the steps to build the tree from the root.
	reverseChain(steps) // Ensure the last step is the root (Human)
	//fmt.Println(steps)
	if len(steps) > 0{
		steps = steps[1:]
	}
	

	nextID := (1000 * idx) + 1
	nodes := []Node{}
	edges := []Edge{}

	nodeIDs := map[string]int{}
	parentOf := map[string]string{}
	childCount := map[string]int{}

	// Initialize the root node first
	nodeIDs[root] = nextID
	nextID++
	nodes = append(nodes, Node{ID: nodeIDs[root], Label: root}) // Add root node

	// Create a map to hold the parent-child relationships
	for _, step := range steps {
		// Ensure that each child gets added
		fmt.Println("Node: ", step.Target, ", Parent: ", step.Parent)
		nodeIDs[step.Target] = nextID
			nextID++
			nodes = append(nodes, Node{ID: nodeIDs[step.Target], Label: step.Target})

		// Now link the parent with the child
		if step.Parent != "" {
			// Make sure the parent node exists in the tree
			if _, exists := nodeIDs[step.Parent]; !exists {
				nodeIDs[step.Parent] = nextID
				nextID++
				nodes = append(nodes, Node{ID: nodeIDs[step.Parent], Label: step.Parent}) // Add parent node if it doesn't exist
			}

			// Add the parent-child relationship
			parentOf[step.Target] = step.Parent
			childCount[step.Parent]++

			// Create an edge from parent to the target (child)
			edges = append(edges, Edge{
				From: nodeIDs[step.Parent],
				To:   nodeIDs[step.Target],
			})
		}
	}

	for parent, count := range childCount {
        if count == 1 {
            // Find the existing child of this parent
            for target, p := range parentOf {
                if p == parent {
                    // Duplicate the child node
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

	// Return the GraphResponse with nodes and edges
	return GraphResponse{Nodes: nodes, Edges: edges}
}




// Function to reverse the chain of steps (to build the tree in the correct order).
func reverseChain(chain []Step) {
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
}

// Helper function to check for duplicate chains.
func isDuplicateChain(chain []Step, seenChains map[string]bool) bool {
	sig := chainSignature(chain)
	if seenChains[sig] {
		return true
	}
	seenChains[sig] = true
	return false
}

// Build a unique signature for each chain.
func chainSignature(chain []Step) string {
	s := ""
	for _, step := range chain {
		s += fmt.Sprintf("%s>%s;", step.Target, step.Parent)
	}
	return s
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

	// Output Tree JSON
	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling JSON: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Tree JSON for DFS Solution:")
	fmt.Println(string(jsonData))

}

// Helper function
func countUsesInChain(element string, chain []Step) int {
	count := 0
	for _, step := range chain {
		if step.Target == element {
			count++
		}
	}
	return count
}

