import { useState } from "react"
import { RecipeTree } from "./RecipeTree"
import { SearchPanel } from "./SearchPanel"

// RecipeExplorer.tsx
export function RecipeExplorer() {
  const [root, setRoot] = useState<string>("")
  const [alg,  setAlg]  = useState<"bfs"|"dfs">("bfs")
  const [mode, setMode] = useState<"shortest"|"multiple">("shortest")
  const [max,  setMax]  = useState<number>(5)
  const [nonce, setNonce] = useState(0)

  const handleSearch = (
    element: string,
    algorithm: "bfs"|"dfs",
    mode: "shortest"|"multiple",
    maxRecipes: number
  ) => {
    setRoot(element)
    setAlg(algorithm)
    setMode(mode)
    setMax(maxRecipes)
    setNonce(n => n+1)
  }

  return (
    <>
      <SearchPanel onSearch={handleSearch} />
      {root && (
        <RecipeTree
          key={`${root}-${alg}-${mode}-${max}-${nonce}`}
          root={root}
          alg={alg}
          mode={mode}
          maxRecipes={max}
        />
      )}
    </>
  )
}
