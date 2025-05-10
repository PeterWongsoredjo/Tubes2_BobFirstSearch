// src/app/page.tsx
"use client"
import { MainTitle } from "@/components/MainTitle"
import { RecipeExplorer } from "@/components/RecipeExplorer"
import { StatsPanel } from "@/components/StatsPanel"

export default function Home() {
  return (
    <main className="min-h-screen bg-black/80">
      <div className="container mx-auto p-4 space-y-8">
        <MainTitle />
        <RecipeExplorer />
        <StatsPanel />
      </div>
    </main>
  )
}
