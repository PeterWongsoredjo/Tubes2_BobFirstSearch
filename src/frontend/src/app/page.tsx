import { SearchPanel } from "@/components/SearchPanel"
import { RecipeTree } from "@/components/RecipeTree"
import { StatsPanel } from "@/components/StatsPanel"
import { MainTitle } from "@/components/MainTitle"

export default function Home() {
  return (
    <main className="min-h-screen celestial-background pt-8">
      <div className="container mx-auto px-4 py-4 max-w-6xl">
        <MainTitle />
        <div className="space-y-8 mt-6">
          <SearchPanel />
          <RecipeTree />
          <StatsPanel />
        </div>
      </div>
    </main>
  )
}
