package main

import (
	//"encoding/json"
	//"fmt"
	"fmt"
	//"log"
	//"os"
	"sync"
	"sync/atomic"
)

type ChainNode struct {
	Recipe Recipe
	Parent *ChainNode
}

type QueueItemO struct {
	Elems   []string
	Chain   *ChainNode
	Depth   int
	Visited map[string]bool
}

func buildChainList(node *ChainNode) []Recipe {
	var chain []Recipe
	for node != nil {
		chain = append([]Recipe{node.Recipe}, chain...)
		node = node.Parent
	}
	return chain
}

func bfs(target string, idx map[string][][]string, tiers map[string]int, limit int) ([]([]Recipe), int) {
	queue := []QueueItemO{{
		Elems:   []string{target},
		Chain:   nil,
		Depth:   0,
		Visited: map[string]bool{},
	}}
	solutions := []([]Recipe){}
	seenChains := map[string]bool{}
	nodesVisited := 0

	var mu sync.Mutex
	var limitReached int32 = 0 // atomic flag

	for depth := 0; len(queue) > 0 && depth < 50; depth++ {
		if atomic.LoadInt32(&limitReached) == 1 {
			break
		}

		nextQueue := []QueueItemO{}
		var wg sync.WaitGroup
		var nextQueueMu sync.Mutex

		currentQueue := make([]QueueItemO, len(queue))
		copy(currentQueue, queue)
		queue = nil

		for _, item := range currentQueue {
			if atomic.LoadInt32(&limitReached) == 1 {
				continue
			}

			wg.Add(1)
			go func(item QueueItemO) {
				defer wg.Done()

				if atomic.LoadInt32(&limitReached) == 1 {
					return
				}

				mu.Lock()
				nodesVisited++
				mu.Unlock()

				if len(item.Elems) == 0 {
					chainList := buildChainList(item.Chain)
					if isFullyResolved(chainList, nil) {
						key := chainKey(chainList)
						mu.Lock()
						if !seenChains[key] {
							seenChains[key] = true
							solutions = append(solutions, deduplicateChain(chainList))

							if len(solutions) >= limit {
								atomic.StoreInt32(&limitReached, 1)
							}
						}
						mu.Unlock()
					}
					return
				}

				if atomic.LoadInt32(&limitReached) == 1 {
					return
				}

				localNextQueue := []QueueItemO{}

				elem := item.Elems[0]
				rest := item.Elems[1:]

				if baseElements[elem] || item.Visited[elem] {
					localNextQueue = append(localNextQueue, QueueItemO{
						Elems:   rest,
						Chain:   item.Chain,
						Depth:   item.Depth,
						Visited: copyMap(item.Visited),
					})
				} else {
					recipes := idx[elem]
					for _, comps := range recipes {
						c1, c2 := comps[0], comps[1]
						tTier := tiers[elem]
						if tiers[c1] >= tTier || tiers[c2] >= tTier {
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

						localNextQueue = append(localNextQueue, QueueItemO{
							Elems:   newElems,
							Chain:   node,
							Depth:   depth + 1,
							Visited: newVisited,
						})
					}
				}

				if len(localNextQueue) > 0 && atomic.LoadInt32(&limitReached) == 0 {
					nextQueueMu.Lock()
					nextQueue = append(nextQueue, localNextQueue...)
					nextQueueMu.Unlock()
				}
			}(item)
		}

		wg.Wait()
		mu.Lock()
		if len(solutions) >= limit {
			mu.Unlock()
			return solutions[:limit], nodesVisited 
		}
		mu.Unlock()

		queue = nextQueue
	}

	return solutions, nodesVisited
}

func buildMultipleTrees(root string, chains [][]Recipe) []GraphResponse {
	var trees []GraphResponse
	for idx, chain := range chains {
		pairs := collectEdgesFromChain(chain)
		tree := buildTrueTree(root, pairs, idx)
		trees = append(trees, tree)
	}
	return trees
}

func buildTrueTree(root string, pairs [][2]string, idx int) GraphResponse {
	idOf := map[string]int{}
	parentOf := map[string]int{}
	var nodes []Node
	var edges []Edge
	nextID := (1000 * idx) + 1

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

/*
func mainaf() {
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

	// Set a limit for how many solutions to find (e.g., 5)
	limit := 5
	chains := bfs(root, idx, tiers, limit)
	if len(chains) == 0 {
		fmt.Printf("No chains found for %q\n", root)
		return
	}

	fmt.Printf("Found %d chains for %q:\n", len(chains), root)

	// Print all chains
	for i, chain := range chains {
		fmt.Printf("\nChain %d:\n", i+1)
		for _, step := range chain {
			fmt.Printf("  %s = %s + %s\n", step.Result, step.Components[0], step.Components[1])
		}
	}

	// Build multiple trees
	trees := buildMultipleTrees(root, chains)

	// Print all trees as JSON
	fmt.Println("\nTrees JSON:")
	for i, tree := range trees {
		fmt.Printf("\nTree %d:\n", i+1)
		b, _ := json.MarshalIndent(tree, "", "  ")
		fmt.Println(string(b))
	}
}
*/
