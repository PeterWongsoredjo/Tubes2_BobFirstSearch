"use client"

import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

interface StatsProps {
  stats: {
    searchTime: number | null;
    nodesVisited: number | null;
    recipesFound: number | null;
  };
  alg?: "bfs" | "dfs";
  mode?: "shortest" | "multiple";
}

export function StatsPanel({ stats, alg, mode }: StatsProps) {
  const algorithmName = alg === "bfs" ? "Breadth-First Search" : "Depth-First Search";
  const modeName = mode === "shortest" ? "Shortest Path" : "Multiple Recipes";

  return (
    <Card className="w-full shadow-lg border-amber-800/50 bg-card/90 backdrop-blur-sm glow-border">
      <CardHeader className="bg-gradient-to-r from-amber-950/50 to-amber-900/50 rounded-t-lg border-b border-amber-800/50">
        <CardTitle className="text-xl text-amber-300 glow-text font-cinzel flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-2">
            <path d="M12 20v-6M6 20V10M18 20V4"></path>
          </svg>
          Search Results
        </CardTitle>
      </CardHeader>
      <CardContent className="pt-6">
        {!stats || stats.searchTime === null ? (
          <div className="text-center py-8 text-amber-300/70">
            <p>Run a search to see statistics</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-secondary/30 rounded-lg p-4 border border-amber-800/30">
              <div className="flex items-center gap-2 mb-1 text-amber-200">
                <span className="text-sm font-medium">Search Time</span>
              </div>
              <p className="text-xl font-semibold text-amber-300">{stats.searchTime} ms</p>
            </div>

            <div className="bg-secondary/30 rounded-lg p-4 border border-amber-800/30">
              <div className="flex items-center gap-2 mb-1 text-amber-200">
                <span className="text-sm font-medium">Nodes Visited</span>
              </div>
              <p className="text-xl font-semibold text-amber-300">{stats.nodesVisited}</p>
            </div>

            <div className="bg-secondary/30 rounded-lg p-4 border border-amber-800/30">
              <div className="flex items-center gap-2 mb-1 text-amber-200">
                <span className="text-sm font-medium">Recipes Found</span>
              </div>
              <p className="text-xl font-semibold text-amber-300">{stats.recipesFound}</p>
            </div>

            <div className="bg-secondary/30 rounded-lg p-4 border border-amber-800/30">
              <div className="flex flex-col h-full justify-between">
                <div>
                  <div className="flex justify-between text-sm">
                    <span className="text-amber-200/70">Algorithm</span>
                    <span className="font-medium text-amber-200">{algorithmName || "—"}</span>
                  </div>
                  
                  <div className="flex justify-between text-sm mt-1">
                    <span className="text-amber-200/70">Search Mode</span>
                    <span className="font-medium text-amber-200">{modeName || "—"}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
