package main

import (
	"encoding/json"
	"os"
)
var baseElements = map[string]bool{
    "Water": 	true,
    "Fire":  	true,
    "Air":   	true,
    "Earth": 	true,
    "Time": 	true,
}

type Recipe struct {
    Result     	string   `json:"result"`
    Components 	[]string `json:"components"`
}

type Element struct {
	Name		string
	Tier		int
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

func LoadElements(path string) ([]Element, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var raw map[string]int
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, err
	}

	var elements []Element
	for name, tier := range raw {
		elements = append(elements, Element{
			Name: name,
			Tier: tier,
		})
	}

	return elements, nil
}