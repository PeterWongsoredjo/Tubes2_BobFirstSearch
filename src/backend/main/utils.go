package main

import (
    "encoding/json"
    "os"
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
