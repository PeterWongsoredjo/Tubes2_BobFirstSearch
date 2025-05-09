package main

import (
    "context"
    "sync"
)

func findAll(target string, idx map[string][][]string, maxRecipes int) [][]Recipe {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var (
        wg      sync.WaitGroup
        mu      sync.Mutex
        results [][]Recipe
    )

    var dfs func(elem string, chain []Recipe)
    dfs = func(elem string, chain []Recipe) {
        select {
        case <-ctx.Done():
            return
        default:
        }

        for _, comps := range idx[elem] {
            nextChain := append(append([]Recipe{}, chain...), Recipe{
                Result:     elem,
                Components: comps,
            })

            if baseElements[comps[0]] && baseElements[comps[1]] {
                mu.Lock()
                if len(results) < maxRecipes {
                    results = append(results, nextChain)
                    if len(results) >= maxRecipes {
                        cancel()
                    }
                }
                mu.Unlock()
                continue
            }

            wg.Add(1)
            go func(c []string, ch []Recipe) {
                defer wg.Done()
                if !baseElements[c[0]] {
                    dfs(c[0], ch)
                }
                if !baseElements[c[1]] {
                    dfs(c[1], ch)
                }
            }(comps, nextChain)
        }
    }

    wg.Add(1)
    go func() {
        defer wg.Done()
        dfs(target, nil)
    }()

    wg.Wait()
    return results
}
