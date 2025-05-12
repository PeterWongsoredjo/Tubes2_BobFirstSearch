package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ElementStatus represents the resolution status of an element
type ElementStatus int

const (
	Pending ElementStatus = iota
	Processing
	Resolved
)

// ElementNode represents an instance of an element in our recipe tree
type ElementNode struct {
	ID       int            // Unique identifier for this instance
	Element  string         // Element name
	Status   ElementStatus  // Current status
	Parent   *ElementNode   // Parent node (null for root)
	Children []*ElementNode // Child nodes (components)
	Recipe   Recipe         // Recipe used to create this element
}

// QueueItem represents an item in our preprocessing queue
type QueueItem struct {
	Node  *ElementNode
	Depth int
}

func bfs(target string, idx map[string][][]string, tiers map[string]int, limit int) ([]([]Recipe), int) {
	startTime := time.Now()
	maxExecutionTime := 60 * time.Second

	// Initialize node ID counter
	nextID := 1

	// Create the root node (target element)
	rootNode := &ElementNode{
		ID:       nextID,
		Element:  target,
		Status:   Processing,
		Children: []*ElementNode{},
	}
	nextID++

	// Track all nodes by ID
	allNodes := make(map[int]*ElementNode)
	allNodes[rootNode.ID] = rootNode

	// Initialize preprocessing queue with target element
	queue := []QueueItem{{Node: rootNode, Depth: 0}}

	// Track nodes visited for stats
	nodesVisited := 0

	// Process elements until solution is found or timeout occurs
	for len(queue) > 0 && time.Since(startTime) < maxExecutionTime {
		// Get next item
		item := queue[0]
		queue = queue[1:]
		node := item.Node
		depth := item.Depth

		nodesVisited++

		// Skip if this specific node instance is already resolved
		if node.Status == Resolved {
			continue
		}

		// Handle base elements
		if baseElements[node.Element] {
			node.Status = Resolved

			// Update parent if this is a child node
			if node.Parent != nil {
				checkAndUpdateParent(node.Parent)
			}

			continue
		}

		// Get recipes for this element
		recipes := idx[node.Element]
		if len(recipes) == 0 {
			// No recipes found, element can't be resolved
			continue
		}

		// Try each recipe
		for _, components := range recipes {
			c1, c2 := components[0], components[1]
			// Skip invalid recipes (higher tier components)
			if tiers[c1] >= tiers[node.Element] || tiers[c2] >= tiers[node.Element] {
				continue
			}

			// Create child nodes (always create new instances)
			leftChild := &ElementNode{
				ID:       nextID,
				Element:  c1,
				Status:   resolveStatus(c1),
				Parent:   node,
				Children: []*ElementNode{},
			}
			nextID++
			allNodes[leftChild.ID] = leftChild

			rightChild := &ElementNode{
				ID:       nextID,
				Element:  c2,
				Status:   resolveStatus(c2),
				Parent:   node,
				Children: []*ElementNode{},
			}
			nextID++
			allNodes[rightChild.ID] = rightChild

			// Store recipe and link children
			node.Recipe = Recipe{
				Result:     node.Element,
				Components: []string{c1, c2},
			}
			node.Children = []*ElementNode{leftChild, rightChild}

			// Check if children are resolved immediately (base elements)
			if leftChild.Status == Resolved && rightChild.Status == Resolved {
				node.Status = Resolved

				// Update parent's parent recursively
				if node.Parent != nil {
					checkAndUpdateParent(node.Parent)
				}
			} else {
				// Add children to queue if they're not already resolved
				if leftChild.Status != Resolved {
					queue = append(queue, QueueItem{Node: leftChild, Depth: depth + 1})
				}

				if rightChild.Status != Resolved {
					queue = append(queue, QueueItem{Node: rightChild, Depth: depth + 1})
				}
			}

			// We used to stop here with "if rootNode.Status == Resolved { break }"
			// but that made our BFS behave like DFS - we need to explore all recipes
			// This ensures we find ALL possible solutions breadth-first
		}

		// If we need to find multiple solutions, we'll do a post-processing phase
		// If we only need one solution and the root is resolved, we can stop early
		if limit == 1 && rootNode.Status == Resolved {
			break
		}
	}

	// Check if target was resolved
	if rootNode.Status != Resolved {
		return []([]Recipe){}, nodesVisited
	}

	// Build solution chains
	solutions := buildSolutions(rootNode, limit)

	return solutions, nodesVisited
}

// checkAndUpdateParent checks if a parent can be resolved
func checkAndUpdateParent(parent *ElementNode) {
	// Skip if parent is already resolved
	if parent.Status == Resolved {
		return
	}

	// Check if all children are resolved
	allResolved := true
	for _, child := range parent.Children {
		if child.Status != Resolved {
			allResolved = false
			break
		}
	}

	// If all children are resolved, mark parent as resolved
	if allResolved && len(parent.Children) > 0 {
		parent.Status = Resolved

		// Update parent's parent recursively
		if parent.Parent != nil {
			checkAndUpdateParent(parent.Parent)
		}
	}
}

// explorePath is a recursive helper function that explores the recipe tree
// to find all possible solution paths
func explorePath(node *ElementNode, currentChain []Recipe, solutions [][]Recipe, maxSolutions int) [][]Recipe {
	// Base case: we have enough solutions or node is nil
	if len(solutions) >= maxSolutions || node == nil {
		return solutions
	}

	// Add current recipe to chain if it's not a base element and has children
	if !baseElements[node.Element] && len(node.Children) > 0 {
		currentChain = append(currentChain, node.Recipe)
	}

	// If this is a leaf node or we've reached a base element
	if len(node.Children) == 0 || baseElements[node.Element] {
		if len(currentChain) > 0 {
			deduped := deduplicateChain(currentChain)
			if isFullyResolved(deduped, nil) {
				// Check if this solution is already in our list
				key := chainKey(deduped)
				isNew := true
				for _, existing := range solutions {
					if chainKey(existing) == key {
						isNew = false
						break
					}
				}
				if isNew {
					solutions = append(solutions, deduped)
				}
			}
		}
		return solutions
	}

	// Explore paths recursively with a copy of the current chain
	for _, child := range node.Children {
		chainCopy := make([]Recipe, len(currentChain))
		copy(chainCopy, currentChain)
		solutions = explorePath(child, chainCopy, solutions, maxSolutions)
		if len(solutions) >= maxSolutions {
			break
		}
	}

	return solutions
}

// buildSolutions constructs solution chains from the resolved tree
func buildSolutions(root *ElementNode, limit int) []([]Recipe) {
	solutions := []([]Recipe){}

	// Build a solution by traversing the resolved tree
	if root.Status == Resolved {
		chain := []Recipe{}
		collectRecipes(root, &chain)

		// Deduplicate and verify
		if len(chain) > 0 {
			deduped := deduplicateChain(chain)
			if isFullyResolved(deduped, nil) {
				solutions = append(solutions, deduped)
			}
		}
	}

	// Find more solutions if needed by exploring alternative paths
	if len(solutions) < limit {
		// Start exploring from root with empty chain
		solutions = explorePath(root, []Recipe{}, solutions, limit)
	}

	return solutions
}

// collectRecipes recursively collects recipes from the tree
func collectRecipes(node *ElementNode, chain *[]Recipe) {
	if node == nil || baseElements[node.Element] {
		return
	}

	// Process children first (bottom-up)
	for _, child := range node.Children {
		collectRecipes(child, chain)
	}

	// Add this node's recipe
	*chain = append(*chain, node.Recipe)
}

// Helper functions remain the same
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

func chainKey(chain []Recipe) string {
	var parts []string
	for _, r := range chain {
		parts = append(parts, fmt.Sprintf("%s=%s+%s", r.Result, r.Components[0], r.Components[1]))
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

func isFullyResolved(chain []Recipe, _ map[string]bool) bool {
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

func resolveStatus(el string) ElementStatus {
	if baseElements[el] {
		return Resolved
	}
	return Processing
}

// buildTrueTree builds a visualization tree directly from a recipe chain
// buildMultipleTrees constructs visualization trees from recipe chains
func buildMultipleTrees(root string, chains [][]Recipe) []GraphResponse {
    var trees []GraphResponse
    
    for i, chain := range chains {
        if len(chain) == 0 {
            continue
        }
        
        tree := buildTreeFromRecipeChain(root, chain, i)
        trees = append(trees, tree)
    }
    
    return trees
}

// buildTreeFromRecipeChain builds a tree directly from a recipe chain
func buildTreeFromRecipeChain(root string, chain []Recipe, idx int) GraphResponse {
    // Start with a high base ID to avoid conflicts
    nextID := (1000 * idx) + 1
    
    // Track elements and their node instances
    elementNodes := make(map[string][]int)
    nodeLabels := make(map[int]string)
    nodeParents := make(map[int]int)
    
    var nodes []Node
    var edges []Edge
    
    // Create root node
    rootID := nextID
    nextID++
    nodeLabels[rootID] = root
    elementNodes[root] = []int{rootID}
    
    // Build a map of recipes for quick lookup
    recipeMap := make(map[string]Recipe)
    for _, recipe := range chain {
        recipeMap[recipe.Result] = recipe
    }
    
    // Process the elements recursively, starting from root
    var processElement func(element string) int
    processElement = func(element string) int {
        // If we've seen this element before, return its node ID
        // For non-base elements, always use the first instance to maintain tree structure
        if ids, exists := elementNodes[element]; exists && !baseElements[element] {
            return ids[0]
        }
        
        // For base elements, create a new instance each time
        if baseElements[element] {
            if _, exists := elementNodes[element]; exists {
                nodeID := nextID
                nextID++
                nodeLabels[nodeID] = element
                elementNodes[element] = append(elementNodes[element], nodeID)
                return nodeID
            }
        }
        
        // Create a new node for this element
        nodeID := nextID
        nextID++
        nodeLabels[nodeID] = element
        elementNodes[element] = append(elementNodes[element], nodeID)
        
        // Check if this element has a recipe in our chain
        recipe, hasRecipe := recipeMap[element]
        if hasRecipe {
            // Process its components
            comp1 := recipe.Components[0]
            comp2 := recipe.Components[1]
            
            comp1ID := processElement(comp1)
            comp2ID := processElement(comp2)
            
            // Set parent relationships
            nodeParents[comp1ID] = nodeID
            nodeParents[comp2ID] = nodeID
            
            // Add edges
            edges = append(edges, Edge{From: nodeID, To: comp1ID})
            edges = append(edges, Edge{From: nodeID, To: comp2ID})
        }
        
        return nodeID
    }
    
    // Start processing from the root
    processElement(root)
    
    // Build nodes array from our maps
    for id, label := range nodeLabels {
        node := Node{
            ID:    id,
            Label: label,
        }
        if parent, hasParent := nodeParents[id]; hasParent {
            node.Parent = parent
        }
        nodes = append(nodes, node)
    }
    
    return GraphResponse{Nodes: nodes, Edges: edges}
}

// Helper function to extract recipes from parent-child pairs
func extractRecipesFromPairs(pairs [][2]string) map[string]Recipe {
    recipeMap := make(map[string]Recipe)
    
    // Group pairs by parent
    parentToChildren := make(map[string][]string)
    for _, pair := range pairs {
        parent := pair[0]
        child := pair[1]
        parentToChildren[parent] = append(parentToChildren[parent], child)
    }
    
    // Create recipes
    for parent, children := range parentToChildren {
        if len(children) == 2 {
            recipeMap[parent] = Recipe{
                Result:     parent,
                Components: []string{children[0], children[1]},
            }
        } else if len(children) == 1 {
            // Handle case where one component appears twice
            recipeMap[parent] = Recipe{
                Result:     parent,
                Components: []string{children[0], children[0]},
            }
        }
    }
    
    return recipeMap
}
