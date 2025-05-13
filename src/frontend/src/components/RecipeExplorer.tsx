import { useState, useCallback } from "react"
import { RecipeTree } from "./RecipeTree"
import { SearchPanel } from "./SearchPanel"
import { StatsPanel } from "./StatsPanel"

export function RecipeExplorer() {
  const [root, setRoot] = useState<string>("")
  const [alg, setAlg] = useState<"bfs"|"dfs"|"splitbfs">("bfs")
  const [mode, setMode] = useState<"shortest"|"multiple">("shortest")
  const [max, setMax] = useState<number>(5)
  const [nonce, setNonce] = useState(0)
  const [stats, setStats] = useState({
    searchTime: null as number | null,
    nodesVisited: null as number | null,
    recipesFound: null as number | null
  })

  const handleSearch = (
    element: string,
    algorithm: "bfs"|"dfs"|"splitbfs",
    mode: "shortest"|"multiple",
    maxRecipes: number
  ) => {
    console.log("RecipeExplorer.handleSearch got:", { element, algorithm, mode, maxRecipes })
    setRoot(element)
    setAlg(algorithm)
    setMode(mode)
    setMax(maxRecipes)
    setNonce(n => n+1)
    setStats({
      searchTime: null,
      nodesVisited: null,
      recipesFound: null
    })
  }

  const handleStatsUpdate = useCallback((newStats: { searchTime: number | null; nodesVisited: number | null; recipesFound: number | null }) => {
    console.log("📊 Received stats update:", newStats)
    setStats(newStats)
  }, [])

  return (
    <div className="space-y-6">
      <SearchPanel onSearch={handleSearch} />
      
      {root && (
        <>
          <RecipeTree
            key={`${root}-${alg}-${mode}-${max}-${nonce}`}
            root={root}
            alg={alg}
            mode={mode}
            maxRecipes={max}
            onStatsUpdate={handleStatsUpdate}
          />
          
          <div className="mt-6">
            <StatsPanel 
              stats={stats}
              alg={alg}
              mode={mode}
            />
          </div>
        </>
      )}
    </div>
  )
}
