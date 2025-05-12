package main

/*import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run TestDriver.go <target_element> [max_recipes]")
		os.Exit(1)
	}

	// Parse command-line arguments
	target := os.Args[1]
	maxRecipes := 1
	if len(os.Args) > 2 {
		fmt.Sscanf(os.Args[2], "%d", &maxRecipes)
	}

	// Load recipes and tiers
	fmt.Println("Loading configuration files...")
	recipes, err := loadRecipes("configs/recipes.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading recipes: %v\n", err)
		os.Exit(1)
	}

	tiers, err := loadTiers("configs/tiers.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tiers: %v\n", err)
		os.Exit(1)
	}

	// Build the recipe index
	idx := buildIndex(recipes)

	// Run both algorithms
	fmt.Printf("Finding recipes for %s (limit: %d)...\n", target, maxRecipes)
	fmt.Println("=== Running BFS algorithm ===")
	bfsChains, bfsNodesVisited := bfs(target, idx, tiers, maxRecipes)
	fmt.Printf("BFS Found %d solutions, visited %d nodes\n", len(bfsChains), bfsNodesVisited)

	// Create visualization trees
	fmt.Println("=== Building visualization trees ===")
	bfsTrees := buildMultipleTrees(target, bfsChains)

	// Output BFS results
	if len(bfsTrees) > 0 {
		fmt.Println("=== BFS Solution Trees ===")
		for i, tree := range bfsTrees {
			fmt.Printf("BFS Tree #%d:\n", i+1)
			printTreeSummary(tree)

			// Output full tree to file
			jsonData, err := json.MarshalIndent(tree, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
				continue
			}

			outFile := fmt.Sprintf("bfs_tree_%d.json", i+1)
			err = os.WriteFile(outFile, jsonData, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to file %s: %v\n", outFile, err)
				continue
			}
			fmt.Printf("Complete tree saved to %s\n", outFile)
		}
	} else {
		fmt.Println("No BFS solutions found!")
	}
}

// Helper function to print a summary of a tree
func printTreeSummary(tree GraphResponse) {
	fmt.Printf("  - Nodes: %d\n", len(tree.Nodes))
	fmt.Printf("  - Edges: %d\n", len(tree.Edges))

	// Print the first few nodes and their connections
	maxNodes := 5
	if len(tree.Nodes) < maxNodes {
		maxNodes = len(tree.Nodes)
	}

	fmt.Printf("  - First %d nodes:\n", maxNodes)
	for i := 0; i < maxNodes; i++ {
		node := tree.Nodes[i]
		fmt.Printf("    Node %d: %s\n", node.ID, node.Label)

		// Find connections
		for _, edge := range tree.Edges {
			if edge.From == node.ID {
				// Find the target node label
				targetLabel := "unknown"
				for _, target := range tree.Nodes {
					if target.ID == edge.To {
						targetLabel = target.Label
						break
					}
				}
				fmt.Printf("      -> %s (Node %d)\n", targetLabel, edge.To)
			}
		}
	}
}
*/