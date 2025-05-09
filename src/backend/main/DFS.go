package main

import (
    "fmt"
	"os"
	"strings"
)

var recipeMap map[string][][]string
var elementMap map[string]int

func buildRecipeMap(recipes []Recipe) {
    recipeMap = make(map[string][][]string)
    for _, r := range recipes {
        recipeMap[r.Result] = append(recipeMap[r.Result], r.Components)
    }
}

func buildElementMap(elements []Element) {
	elementMap = make(map[string]int)
	for _, e := range elements {
        elementMap[e.Name] = e.Tier
    }
}

func dfs(target string, depth int, visited map[string]bool) {
	if baseElements[target] {
        return
    }

    if visited[target] {
        return
    }

	
    visited[target] = true
    indent := strings.Repeat("  ", depth)

    if baseElements[recipeMap[target][0][0]] && baseElements[recipeMap[target][0][1]] {
        fmt.Println(indent + target + " = " + recipeMap[target][0][0] + " + " + recipeMap[target][0][1])
        return
    }

	index := 0
    if componentsList, ok := recipeMap[target]; ok {
		for i, components := range componentsList {
			if( elementMap[target] > elementMap[components[0]] && elementMap[target] > elementMap[components[1]]){
				index = i
				break
			}
		}
        fmt.Println(indent + target + " = " + componentsList[index][0] + " + " + componentsList[index][1])
        dfs(componentsList[index][0], depth+1, visited)
        dfs(componentsList[index][1], depth+1, visited)
    }
}

func main() {
	recipes, err := loadRecipes("configs/recipes.json")
	elements, err := LoadElements("configs/tiers.json")
	if err != nil {
        fmt.Fprintf(os.Stderr, "error loading recipes: %v\n", err)
        os.Exit(1)
    }
    buildRecipeMap(recipes)
	buildElementMap(elements)
    dfs("Sandwich", 0, make(map[string]bool))
}
