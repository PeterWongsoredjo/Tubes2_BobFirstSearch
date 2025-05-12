package main

import (
	"encoding/json"
	"os"
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

type Element struct {
    Name string
    Tier int
}

// TreeNode is one node in our decomposition tree.
type TreeNode struct {
    ID       int
    Label    string
    Children []*TreeNode
}

type Node struct {
    ID    int    `json:"id"`
    Label string `json:"label"`
    Parent int   `json:"parent,omitempty"`
}

type Edge struct {
    From int `json:"from"`
    To   int `json:"to"`
}

type GraphResponse struct {
    Nodes []Node `json:"nodes"`
    Edges []Edge `json:"edges"`
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

func loadTiers(path string) (map[string]int, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var tiers map[string]int
    if err := json.NewDecoder(f).Decode(&tiers); err != nil {
        return nil, err
    }
    return tiers, nil
}

func buildIndex(recipes []Recipe) map[string][][]string {
    idx := make(map[string][][]string)
    seen := map[string]map[string]bool{}
    for _, r := range recipes {
        key := strings.Join(r.Components, "|")
        if seen[r.Result] == nil {
            seen[r.Result] = make(map[string]bool)
        }
        if seen[r.Result][key] {
            continue
        }
        seen[r.Result][key] = true
        idx[r.Result] = append(idx[r.Result], r.Components)
    }
    return idx
}

func copyMap(original map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}

// collectEdgesFromChain extracts parent-child relationships from a recipe chain
func collectEdgesFromChain(chain []Recipe) [][2]string {
	var pairs [][2]string
	for _, step := range chain {
		parent := step.Result
		for _, child := range step.Components {
			pairs = append(pairs, [2]string{parent, child})
		}
	}
	return pairs
}