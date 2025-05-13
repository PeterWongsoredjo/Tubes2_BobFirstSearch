package main

import (
	"time"
)

type ElementStatus int

const (
	Pending ElementStatus = iota
	Processing
	Resolved
)

type ElementNode struct {
	ID       int            // Unique identifier
	Element  string         // Element name
	Status   ElementStatus  // Current status
	Parent   *ElementNode   // Parent node (null for root)
	Children []*ElementNode // Child nodes (components)
	Recipe   Recipe         // Recipe used to create this element
}

type QueueItem struct {
	Node  *ElementNode
	Depth int
}

func splitbfs(target string, idx map[string][][]string, tiers map[string]int, limit int) ([]([]Recipe), int) {
	startTime := time.Now()
	maxExecutionTime := 60 * time.Second

	nextID := 1

	rootNode := &ElementNode{
		ID:       nextID,
		Element:  target,
		Status:   Processing,
		Children: []*ElementNode{},
	}
	nextID++

	allNodes := make(map[int]*ElementNode)
	allNodes[rootNode.ID] = rootNode

	queue := []QueueItem{{Node: rootNode, Depth: 0}}

	nodesVisited := 0

	for len(queue) > 0 && time.Since(startTime) < maxExecutionTime {
		item := queue[0]
		queue = queue[1:]
		node := item.Node
		depth := item.Depth

		nodesVisited++

		if node.Status == Resolved {
			continue
		}

		if baseElements[node.Element] {
			node.Status = Resolved

			if node.Parent != nil {
				checkAndUpdateParent(node.Parent)
			}

			continue
		}

		recipes := idx[node.Element]
		if len(recipes) == 0 {
			continue
		}

		for _, components := range recipes {
			c1, c2 := components[0], components[1]
			if tiers[c1] >= tiers[node.Element] || tiers[c2] >= tiers[node.Element] {
				continue
			}

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

			node.Recipe = Recipe{
				Result:     node.Element,
				Components: []string{c1, c2},
			}
			node.Children = []*ElementNode{leftChild, rightChild}

			if leftChild.Status == Resolved && rightChild.Status == Resolved {
				node.Status = Resolved
				if node.Parent != nil {
					checkAndUpdateParent(node.Parent)
				}
			} else {
				if leftChild.Status != Resolved {
					queue = append(queue, QueueItem{Node: leftChild, Depth: depth + 1})
				}

				if rightChild.Status != Resolved {
					queue = append(queue, QueueItem{Node: rightChild, Depth: depth + 1})
				}
			}
		}

		if limit == 1 && rootNode.Status == Resolved {
			break
		}
	}

	if rootNode.Status != Resolved {
		return []([]Recipe){}, nodesVisited
	}
	solutions := buildSolutions(rootNode, limit)

	processedChains := make([]([]Recipe), 0)

	for _, solution := range solutions {
		chain := processChainForTreeBuilding(solution)
		if len(chain) > 0 {
			processedChains = append(processedChains, chain)
		}
	}

	return processedChains, nodesVisited
}

func processChainForTreeBuilding(chain []Recipe) []Recipe {
	processed := reverseRecipeChainIfNeeded(chain)

	resultToRecipe := make(map[string]Recipe)
	for _, recipe := range processed {
		resultToRecipe[recipe.Result] = recipe
	}

	var orderedChain []Recipe

	visited := make(map[string]bool)
	toVisit := []string{processed[0].Result} 

	for len(toVisit) > 0 {
		element := toVisit[0]
		toVisit = toVisit[1:]

		if visited[element] {
			continue
		}
		visited[element] = true

		if recipe, exists := resultToRecipe[element]; exists {
			for _, comp := range recipe.Components {
				if !baseElements[comp] && !visited[comp] {
					toVisit = append(toVisit, comp)
				}
			}
			orderedChain = append(orderedChain, recipe)
		}
	}

	return orderedChain
}

func reverseRecipeChainIfNeeded(chain []Recipe) []Recipe {
	reversed := make([]Recipe, len(chain))
	for i, recipe := range chain {
		reversed[len(chain)-1-i] = recipe
	}
	return reversed
}

func checkAndUpdateParent(parent *ElementNode) {
	if parent.Status == Resolved {
		return
	}

	allResolved := true
	for _, child := range parent.Children {
		if child.Status != Resolved {
			allResolved = false
			break
		}
	}

	if allResolved && len(parent.Children) > 0 {
		parent.Status = Resolved
		if parent.Parent != nil {
			checkAndUpdateParent(parent.Parent)
		}
	}
}

func explorePath(node *ElementNode, currentChain []Recipe, solutions [][]Recipe, maxSolutions int) [][]Recipe {
	if len(solutions) >= maxSolutions || node == nil {
		return solutions
	}

	if !baseElements[node.Element] && len(node.Children) > 0 {
		currentChain = append(currentChain, node.Recipe)
	}

	if len(node.Children) == 0 || baseElements[node.Element] {
		if len(currentChain) > 0 {
			deduped := deduplicateChain(currentChain)
			if isFullyResolved(deduped, nil) {
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

func buildSolutions(root *ElementNode, limit int) []([]Recipe) {
	solutions := []([]Recipe){}

	if root.Status == Resolved {
		chain := []Recipe{}
		collectRecipes(root, &chain)

		if len(chain) > 0 {
			deduped := deduplicateChain(chain)
			if isFullyResolved(deduped, nil) {
				solutions = append(solutions, deduped)
			}
		}
	}

	if len(solutions) < limit {
		solutions = explorePath(root, []Recipe{}, solutions, limit)
	}

	return solutions
}

func collectRecipes(node *ElementNode, chain *[]Recipe) {
	if node == nil || baseElements[node.Element] {
		return
	}

	for _, child := range node.Children {
		collectRecipes(child, chain)
	}

	*chain = append(*chain, node.Recipe)
}

func resolveStatus(el string) ElementStatus {
	if baseElements[el] {
		return Resolved
	}
	return Processing
}


/*
func newbuildMultipleTrees(root string, chains [][]Recipe) []GraphResponse {
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


func collectEdgesFromSplitBFSChain(chain []Recipe) [][2]string {
	var pairs [][2]string

	for _, recipe := range chain {
		parent := recipe.Result
		for _, child := range recipe.Components {
			pair := [2]string{parent, child}
			pairs = append(pairs, pair)
		}
	}

	return pairs
}
*/
