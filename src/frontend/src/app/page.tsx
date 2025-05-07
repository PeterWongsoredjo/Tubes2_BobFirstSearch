import { SearchPanel } from "@/components/SearchPanel"
import { RecipeTree } from "@/components/RecipeTree"
import { StatsPanel } from "@/components/StatsPanel"
import { Navbar } from "@/components/navbar"

export default function Home() {
  return (
    <main className="min-h-screen celestial-background">
      <Navbar />
      <div className="container mx-auto px-4 py-8 max-w-6xl">
        <div className="space-y-8">
          <SearchPanel />
          <div className="grid gap-8 md:grid-cols-3">
            <div className="md:col-span-2">
              <RecipeTree />
            </div>
            <div>
              <StatsPanel />
            </div>
          </div>
        </div>
      </div>
    </main>
  )
}
