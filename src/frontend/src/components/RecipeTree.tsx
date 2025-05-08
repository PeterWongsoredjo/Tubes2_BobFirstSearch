"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardHeader, CardTitle } from "@/components/ui/card"

export function RecipeTree() {
  const [zoom, setZoom] = useState(100)

  const handleZoomIn = () => {
    setZoom((prev) => Math.min(prev + 10, 150))
  }

  const handleZoomOut = () => {
    setZoom((prev) => Math.max(prev - 10, 50))
  }

  return (
    <Card className="w-full h-[500px] shadow-lg border-amber-800/50 flex flex-col bg-card/90 backdrop-blur-sm light-theme">
      <CardHeader className="bg-gradient-to-r from-amber-950/50 to-amber-900/50 rounded-t-lg border-b border-amber-800/50">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-xl text-amber-300  font-cinzel">Recipe Visualization</CardTitle>
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8 border-amber-800/50 bg-secondary/50"
              onClick={handleZoomOut}
            >
              <span className="text-amber-200">-</span>
            </Button>
            <span className="text-xs text-amber-300 w-12 text-center">{zoom}%</span>
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8 border-amber-800/50 bg-secondary/50"
              onClick={handleZoomIn}
            >
              <span className="text-amber-200">+</span>
            </Button>
          </div>
        </div>
      </CardHeader>
    </Card>
  )
}
