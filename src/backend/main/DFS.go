package main

import (
	//"encoding/json"
	"fmt"
	//"os"
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

	built[target] = true
	return allChains, nodesVisited
}

func chainSignature(chain []Step) string {
	s := ""
	for _, step := range chain {
		s += fmt.Sprintf("%s>%s;", step.Target, step.Parent)
	}
	return s
}

func buildTrueTreeFromDFS(root string, steps []Step, idx int) GraphResponse {
	type nodeKey struct {
		Label  string
		Parent string
	}

	nextID := (1000 * idx) + 1
	nodes := []Node{}
	edges := []Edge{}

	nodeIDs := map[string]int{}
	parentOf := map[string]string{}

	for _, step := range steps {
		if _, exists := nodeIDs[step.Target]; !exists {
			nodeIDs[step.Target] = nextID
			nextID++
		}
		if step.Parent != "" {
			if _, exists := nodeIDs[step.Parent]; !exists {
				nodeIDs[step.Parent] = nextID
				nextID++
			}
			parentOf[step.Target] = step.Parent
		}
	}

	for label, id := range nodeIDs {
		node := Node{ID: id, Label: label}
		if parentLabel, ok := parentOf[label]; ok {
			node.Parent = nodeIDs[parentLabel]
			edges = append(edges, Edge{
				From: node.Parent,
				To:   id,
			})
		}
		nodes = append(nodes, node)
	}

	return GraphResponse{Nodes: nodes, Edges: edges}
}

/*
func maien() {
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
	solutions := dfsAll("Human", built, 1, "", )

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
*/