package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "sort"
    "strings"
)

var baseElements = map[string]bool{
    "Water": true,
    "Fire":  true,
    "Air":   true,
    "Earth": true,
	"Time":  true,
}

type Recipe struct {
    Result     string   `json:"result"`
    Components []string `json:"components"`
}

type TreeNode struct {
    Name     string
    Children []*TreeNode
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

// buildComponentMap: element → list of unique component-pairs
func buildComponentMap(recipes []Recipe) map[string][][]string {
    dedup := make(map[string]map[string]bool)
    compMap := make(map[string][][]string)
    for _, r := range recipes {
        key := strings.Join(r.Components, "|")
        if dedup[r.Result] == nil {
            dedup[r.Result] = make(map[string]bool)
        }
        if dedup[r.Result][key] {
            continue
        }
        dedup[r.Result][key] = true
        compMap[r.Result] = append(compMap[r.Result], r.Components)
    }
    return compMap
}

// scorePair counts how many base elements are in this pair
func scorePair(pair []string) int {
    cnt := 0
    for _, c := range pair {
        if baseElements[c] {
            cnt++
        }
    }
    return cnt
}

// chooseBestPair picks the recipe with highest base-element score,
// tiebreak by lex order of the joined string.
func chooseBestPair(pairs [][]string) []string {
    best := pairs[0]
    bestScore := scorePair(best)
    bestKey := strings.Join(best, "|")

    for _, p := range pairs[1:] {
        s := scorePair(p)
        key := strings.Join(p, "|")
        if s > bestScore || (s == bestScore && key < bestKey) {
            best, bestScore, bestKey = p, s, key
        }
    }
    // sort so output order is consistent
    sort.Strings(best)
    return best
}

// buildTreeBFS starts from startName, expands by BFS,
// prunes repeats, stops at baseElements, and picks one recipe per node.
func buildTreeBFS(startName string, compMap map[string][][]string) *TreeNode {
    root := &TreeNode{Name: startName}
    queue := []*TreeNode{root}
    visited := map[string]bool{startName: true}

    for i := 0; i < len(queue); i++ {
        node := queue[i]

        // STOP if this is a base element
        if baseElements[node.Name] {
            continue
        }

        pairs := compMap[node.Name]
        if len(pairs) == 0 {
            continue
        }

        // pick the single best recipe
        chosen := chooseBestPair(pairs)
        for _, childName := range chosen {
            if visited[childName] {
                continue
            }
            visited[childName] = true
            child := &TreeNode{Name: childName}
            node.Children = append(node.Children, child)
            queue = append(queue, child)
        }
    }
    return root
}

func printTree(n *TreeNode, level int) {
    fmt.Printf("%s%s\n", strings.Repeat("  ", level), n.Name)
    for _, ch := range n.Children {
        printTree(ch, level+1)
    }
}

func main() {
    recipes, err := loadRecipes("../configs/recipes.json")
    if err != nil {
        fmt.Fprintf(os.Stderr, "error loading recipes: %v\n", err)
        os.Exit(1)
    }

    compMap := buildComponentMap(recipes)

    fmt.Print("Enter an element: ")
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    name := strings.TrimSpace(input)

    if _, ok := compMap[name]; !ok && !baseElements[name] {
        fmt.Printf("no recipes found for element %q\n", name)
        os.Exit(0)
    }

    tree := buildTreeBFS(name, compMap)
    printTree(tree, 0)
}
