package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

type Recipe struct {
    Result     string   `json:"result"`
    Components []string `json:"components"`
}

func loadRecipes(path string) ([]Recipe, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var recs []Recipe
    if err := json.NewDecoder(f).Decode(&recs); err != nil {
        return nil, err
    }
    return recs, nil
}

type TreeNode struct {
    Name     string
    Children []*TreeNode
}

// 1) Build the raw forest with cycle-pruning (as before).
func buildForest(recipes []Recipe) map[string]*TreeNode {
    // make one node per element
    nodes := map[string]*TreeNode{}
    for _, r := range recipes {
        if _, ok := nodes[r.Result]; !ok {
            nodes[r.Result] = &TreeNode{Name: r.Result}
        }
        for _, c := range r.Components {
            if _, ok := nodes[c]; !ok {
                nodes[c] = &TreeNode{Name: c}
            }
        }
    }

    // dedupe recipes
    uniq := []Recipe{}
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
        uniq = append(uniq, r)
    }

    // attach edges, pruning any back-edge
    added := map[string]map[string]bool{}
    var hasPath func(src, tgt *TreeNode, vis map[string]bool) bool
    hasPath = func(src, tgt *TreeNode, vis map[string]bool) bool {
        if src == tgt {
            return true
        }
        vis[src.Name] = true
        for _, ch := range src.Children {
            if !vis[ch.Name] && hasPath(ch, tgt, vis) {
                return true
            }
        }
        return false
    }

    for _, r := range uniq {
        parent := nodes[r.Result]
        if added[parent.Name] == nil {
            added[parent.Name] = map[string]bool{}
        }
        for _, c := range r.Components {
            child := nodes[c]
            if added[parent.Name][child.Name] {
                continue // duplicate
            }
            // prune cycles
            if hasPath(child, parent, map[string]bool{}) {
                continue
            }
            parent.Children = append(parent.Children, child)
            added[parent.Name][child.Name] = true
        }
    }

    return nodes
}

// 2) Find your tier-14 “roots”: those with no incoming edges.
func findRoots(forest map[string]*TreeNode) []*TreeNode {
    indeg := map[string]int{}
    for n := range forest {
        indeg[n] = 0
    }
    for _, node := range forest {
        for _, ch := range node.Children {
            indeg[ch.Name]++
        }
    }
    var roots []*TreeNode
    for name, node := range forest {
        if indeg[name] == 0 {
            roots = append(roots, node)
        }
    }
    return roots
}

// 3) BFS-build a true tree per root, never revisiting a node in that tree.
func buildBFSTrees(forest map[string]*TreeNode, roots []*TreeNode) []*TreeNode {
    var trees []*TreeNode
    for _, root := range roots {
        // copy root
        rootCopy := &TreeNode{Name: root.Name}
        queue := []*TreeNode{rootCopy}
        visited := map[string]bool{root.Name: true}
        // map from name→copy to link children
        copyMap := map[string]*TreeNode{root.Name: rootCopy}

        for i := 0; i < len(queue); i++ {
            curr := queue[i]
            orig := forest[curr.Name]
            for _, ch := range orig.Children {
                if visited[ch.Name] {
                    continue
                }
                visited[ch.Name] = true
                childCopy := &TreeNode{Name: ch.Name}
                curr.Children = append(curr.Children, childCopy)
                queue = append(queue, childCopy)
                copyMap[ch.Name] = childCopy
            }
        }
        trees = append(trees, rootCopy)
    }
    return trees
}

// 4) Print your BFS-built trees (now truly acyclic).
func printTree(node *TreeNode, level int) {
    fmt.Printf("%s%s\n", strings.Repeat("  ", level), node.Name)
    for _, ch := range node.Children {
        printTree(ch, level+1)
    }
}

func main() {
    recipes, err := loadRecipes("../recipes.json")
    if err != nil {
        panic(err)
    }

    forest := buildForest(recipes)
    roots := findRoots(forest)
    trees := buildBFSTrees(forest, roots)

    for _, t := range trees {
        printTree(t, 0)
    }
}
