"use client"
import { MainTitle } from "@/components/MainTitle"
import { RecipeExplorer } from "@/components/RecipeExplorer"
import ScrapingButton from "@/components/ScrapingButton" 

export default function Home() {
  return (
   <main className="min-h-screen celestial-background pt-8">
    <div className="absolute top-0 left-0 ml-4 mt-4"> 
        <ScrapingButton />
      </div>
      <div className="container mx-auto px-4 py-4 max-w-6xl">
        <MainTitle />
        <RecipeExplorer />
      </div>
    </main>
  )
}
