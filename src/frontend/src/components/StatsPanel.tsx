"use client"

import { Clock, Hash, Search, Zap } from "lucide-react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export function StatsPanel() {
  // This would be populated with real data after a search
  const stats = {
    searchTime: null as number | null,
    nodesVisited: null as number | null,
    recipesFound: null as number | null,
    algorithm: null as string | null,
    mode: null as string | null,
  }

  return (
    <Card className="w-full h-full shadow-lg border-amber-800/50 bg-card/90 backdrop-blur-sm glow-border">
      <CardHeader className="bg-gradient-to-r from-amber-950/50 to-amber-900/50 rounded-t-lg border-b border-amber-800/50">
        <CardTitle className="text-xl text-amber-300 glow-text font-cinzel">Search Results</CardTitle>
        <CardDescription className="text-amber-200/70">Statistics from your search</CardDescription>
      </CardHeader>
      <CardContent className="pt-6">
        {stats.searchTime === null ? (
          <div className="text-center py-8 text-amber-300/70">
            <Search className="mx-auto h-8 w-8 mb-3 opacity-50" />
            <p>Run a search to see statistics</p>
          </div>
        ) : (
          <div className="space-y-6">
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-secondary/30 rounded-lg p-4 border border-amber-800/30">
                <div className="flex items-center gap-2 mb-1 text-amber-200">
                  <Clock className="h-4 w-4" />
                  <span className="text-sm font-medium">Search Time</span>
                </div>
                <p className="text-xl font-semibold text-amber-300">{stats.searchTime} ms</p>
              </div>
              <div className="bg-secondary/30 rounded-lg p-4 border border-amber-800/30">
                <div className="flex items-center gap-2 mb-1 text-amber-200">
                  <Hash className="h-4 w-4" />
                  <span className="text-sm font-medium">Nodes Visited</span>
                </div>
                <p className="text-xl font-semibold text-amber-300">{stats.nodesVisited}</p>
              </div>
            </div>

            <div className="bg-secondary/30 rounded-lg p-4 border border-amber-800/30">
              <div className="flex items-center gap-2 mb-1 text-amber-200">
                <Zap className="h-4 w-4" />
                <span className="text-sm font-medium">Recipes Found</span>
              </div>
              <p className="text-xl font-semibold text-amber-300">{stats.recipesFound}</p>
            </div>

            <div className="space-y-3">
              <div className="flex justify-between text-sm">
                <span className="text-amber-200/70">Algorithm</span>
                <span className="font-medium text-amber-200">{stats.algorithm || "—"}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-amber-200/70">Search Mode</span>
                <span className="font-medium text-amber-200">{stats.mode || "—"}</span>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
