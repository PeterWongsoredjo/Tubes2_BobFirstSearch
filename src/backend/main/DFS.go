package main

import (
    "fmt"
	"os"
	"strings"
)

var recipeMap map[string][][]string
var elementMap map[string]int
var count int
type Step struct {
	Target string
	Left   string
	Right  string
}

var allPaths [][]Step // to store N recipe trees
var neededStarts[][]Step


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

func dfs(target string, depth int) {
	if baseElements[target] {
		return
	}

	indent := strings.Repeat("  ", depth)

	if baseElements[recipeMap[target][0][0]] && baseElements[recipeMap[target][0][1]] {
		fmt.Println(indent + target + " = " + recipeMap[target][0][0] + " + " + recipeMap[target][0][1])
		return
	}

	if componentsList, ok := recipeMap[target]; ok {
		for _, components := range componentsList {
			if elementMap[target] > elementMap[components[0]] && elementMap[target] > elementMap[components[1]] {
				fmt.Println(indent + target + " = " + components[0] + " + " + components[1])
                if(components[0] == components[1]){
                    dfs(components[0], depth+1)
                } else{
                    dfs(components[0], depth+1)
                    dfs(components[1], depth+1)
                }
                break
			}
		}
	}
}


func getDepthNeeded(target string, count int, currentDepth int) int {
    // Base case: if we've reached a base element or count is satisfied
    if baseElements[target] || count <= 1 {
        return currentDepth
    }

    maxDepth := currentDepth
    for _, comp := range recipeMap[target] {
        if elementMap[target] > elementMap[comp[0]] && 
           elementMap[target] > elementMap[comp[1]] {
            // Calculate depth for both components
            leftDepth := getDepthNeeded(comp[0], count/2, currentDepth+1)
            rightDepth := getDepthNeeded(comp[1], count/2, currentDepth+1)
            
            // Track the maximum depth needed
            if leftDepth > maxDepth {
                maxDepth = leftDepth
            }
            if rightDepth > maxDepth {
                maxDepth = rightDepth
            }
        }
    }
    return maxDepth
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
	//visitedMap := new(map[string][][]string)
    dfs("Sandwich", 0)
}
