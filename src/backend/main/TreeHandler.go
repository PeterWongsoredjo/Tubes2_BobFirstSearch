package main

import (
	"net/http"
	//"strconv"
	"encoding/json"
)

func treeHandler(idx map[string][][]string, tiers map[string]int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		root := r.URL.Query().Get("root")
		alg := r.URL.Query().Get("alg")   // "bfs" or "dfs"
		mode := r.URL.Query().Get("mode") // "shortest" or "multiple"
		//maxN, _ := strconv.Atoi(r.URL.Query().Get("max"))

		var pairs [][2]string

		switch alg {
		case "bfs":
			if mode == "shortest" {
				// same as collectEdges
				collectEdges(root, idx, tiers, map[string]bool{}, &pairs)
			} /*else {
				chains := findAll(root, idx, maxN)
				for _, ch := range chains {
					for _, step := range ch {
						for _, child := range step.Components {
							pairs = append(pairs, [2]string{step.Result, child})
						}
					}
				}
			}

		case "dfs":
			if mode == "shortest" {
				collectDFS(root, idx, tiers, map[string]bool{}, &pairs)
			} else {
				// you’d have a findAllDFS variant
				chains := findAllDFS(root, idx, maxN)
				// same flattening as above…
			}*/

		default:
			http.Error(w, "unknown alg", 400)
			return
		}

		graph := buildGraphFromPairs(root, pairs)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(graph)
	}
}

func buildGraphFromPairs(root string, pairs [][2]string) GraphResponse {
    // map labels to integer IDs
    idOf := map[string]int{}
    nodes := make([]Node, 0, len(pairs)*2)
    edges := make([]Edge, 0, len(pairs))

    nextID := 1
    // ensure root is present first
    idOf[root] = nextID
    nodes = append(nodes, Node{ID: nextID, Label: root})
    nextID++

    // assign IDs to every label we see in pairs
    for _, p := range pairs {
        parent, child := p[0], p[1]
        if _, ok := idOf[parent]; !ok {
            idOf[parent] = nextID
            nodes = append(nodes, Node{ID: nextID, Label: parent})
            nextID++
        }
        if _, ok := idOf[child]; !ok {
            idOf[child] = nextID
            nodes = append(nodes, Node{ID: nextID, Label: child})
            nextID++
        }
        // record the edge
        edges = append(edges, Edge{
            From: idOf[parent],
            To:   idOf[child],
        })
    }

    return GraphResponse{Nodes: nodes, Edges: edges}
}