"use client"

import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

type Stats = {
  searchTime: number | null
  nodesVisited: number | null
  recipesFound: number | null
}

type StatsPanelProps = {
  stats: Stats
  alg: "bfs" | "dfs"
  mode: "shortest" | "multiple"
}

export function StatsPanel({ stats, alg, mode }: StatsPanelProps) {
  const { searchTime, nodesVisited, recipesFound } = stats

  const formatTime = (ms: number | null) => {
    if (ms === null) return "N/A"
    if (ms < 1000) return `${ms} ms`
    return `${(ms / 1000).toFixed(2)} s`
  }
  
  const algName = alg === "bfs" ? "Breadth-First Search" : "Depth-First Search"
  const modeText = mode === "shortest" ? "Shortest Path" : "Multiple Paths"

  return (
    <Card className="w-full shadow-lg border-amber-800/50 bg-card/90 backdrop-blur-sm">
      <CardHeader className="bg-gradient-to-r from-amber-950/80 to-amber-900/50 p-3 rounded-t-lg border-b border-amber-700/30">
        <CardTitle className="text-amber-300 flex items-center text-lg">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-2">
            <path d="M12 20v-6M6 20V10M18 20V4"></path>
          </svg>
          Algorithm Performance Stats
        </CardTitle>
      </CardHeader>
      <CardContent className="p-4">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-amber-950/30 p-3 rounded-md border border-amber-800/30">
            <div className="text-sm text-amber-200/70 mb-1">Algorithm</div>
            <div className="text-lg text-amber-100 font-medium">{algName}</div>
            <div className="text-xs text-amber-300/80 mt-1">Mode: {modeText}</div>
          </div>
          
          <div className="bg-amber-950/30 p-3 rounded-md border border-amber-800/30">
            <div className="text-sm text-amber-200/70 mb-1">Search Time</div>
            <div className="text-lg text-amber-100 font-medium">
              {searchTime === null ? (
                <span className="text-amber-500/50">Calculating...</span>
              ) : (
                formatTime(searchTime)
              )}
            </div>
            <div className="text-xs text-amber-300/80 mt-1">Total processing time</div>
          </div>
          
          <div className="bg-amber-950/30 p-3 rounded-md border border-amber-800/30">
            <div className="text-sm text-amber-200/70 mb-1">Algorithm Efficiency</div>
            <div className="flex justify-between">
              <div>
                <div className="text-sm text-amber-100">
                  {nodesVisited === null ? (
                    <span className="text-amber-500/50">...</span>
                  ) : (
                    `${nodesVisited} nodes visited`
                  )}
                </div>
                <div className="text-sm text-amber-100">
                  {recipesFound === null ? (
                    <span className="text-amber-500/50">...</span>
                  ) : (
                    `${recipesFound} recipes found`
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
