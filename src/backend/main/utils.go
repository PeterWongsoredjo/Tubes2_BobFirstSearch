package main

import (
	"encoding/json"
	"os"
	"strings"
    "sort"
    "fmt"
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