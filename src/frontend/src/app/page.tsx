"use client"
import { MainTitle } from "@/components/MainTitle"
import { RecipeExplorer } from "@/components/RecipeExplorer"

export default function Home() {
  return (
   <main className="min-h-screen celestial-background pt-8">
      <div className="container mx-auto px-4 py-4 max-w-6xl">
        <MainTitle />
        <RecipeExplorer />
      </div>
    </main>
  )
}
