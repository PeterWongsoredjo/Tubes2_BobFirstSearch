package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

// Recipe represents one Little Alchemy combination result
// and the two component names used to create it.
type Recipe struct {
    Result     string   `json:"result"`
    Components []string `json:"components"`
}

// loadRecipes reads a JSON file at path and decodes it into []Recipe.
func loadRecipes(path string) ([]Recipe, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var recipes []Recipe
    if err := json.NewDecoder(f).Decode(&recipes); err != nil {
        return nil, err
    }
    return recipes, nil
}

// TreeNode represents one element and its recipe-components.
type TreeNode struct {
    Name     string
    Children []*TreeNode
}

// makeNodeMap builds a map[name]*TreeNode so that every
// element (result or component) gets exactly one TreeNode.
func makeNodeMap(recipes []Recipe) map[string]*TreeNode {
    nodes := make(map[string]*TreeNode, len(recipes))
    for _, rec := range recipes {
        if _, exists := nodes[rec.Result]; !exists {
            nodes[rec.Result] = &TreeNode{Name: rec.Result}
        }
        for _, comp := range rec.Components {
            if _, exists := nodes[comp]; !exists {
                nodes[comp] = &TreeNode{Name: comp}
            }
        }
    }
    return nodes
}

// hasPath checks whether 'target' is reachable from 'src' by following Children links.
// It uses a visited map to avoid infinite loops.
func hasPath(src, target *TreeNode, visited map[string]bool) bool {
    if src == target {
        return true
    }
    visited[src.Name] = true

    for _, child := range src.Children {
        if !visited[child.Name] {
            if hasPath(child, target, visited) {
                return true
            }
        }
    }
    return false
}

// buildCompositionForest links every Recipe’s result node to its two component nodes,
// but skips (prunes) any link that would create a cycle.
func buildCompositionForest(recipes []Recipe) map[string]*TreeNode {
    nodes := makeNodeMap(recipes)

    for _, rec := range recipes {
        parent := nodes[rec.Result]
        for _, compName := range rec.Components {
            child := nodes[compName]

            // cycle-pruning: if 'child' already reaches 'parent', skip
            if hasPath(child, parent, make(map[string]bool)) {
                continue
            }
            parent.Children = append(parent.Children, child)
        }
    }
    return nodes
}

// findRoots finds all nodes with zero incoming edges.
// It first tallies indegrees, then returns those with indegree==0.
func findRoots(forest map[string]*TreeNode) []*TreeNode {
    indegree := make(map[string]int, len(forest))
    for name := range forest {
        indegree[name] = 0
    }
    for _, node := range forest {
        for _, child := range node.Children {
            indegree[child.Name]++
        }
    }

    var roots []*TreeNode
    for name, node := range forest {
        if indegree[name] == 0 {
            roots = append(roots, node)
        }
    }
    return roots
}

// printTree recursively prints each node and its children, indenting by level.
func printTree(node *TreeNode, level int) {
    fmt.Printf("%s%s\n", strings.Repeat("  ", level), node.Name)
    for _, child := range node.Children {
        printTree(child, level+1)
    }
}

func main() {
    // 1. Load recipes from JSON
    recipes, err := loadRecipes("../recipes.json")
    if err != nil {
        panic(err)
    }
    // 2. Build forest with cycle-pruning
    forest := buildCompositionForest(recipes)

    // 3. Find root/base elements (no incoming edges)
    roots := findRoots(forest)

    // 4. Print each tree
    for _, root := range roots {
        printTree(root, 0)
    }
}
