"use client"

import { useState } from "react"
import { ZoomIn, ZoomOut, Maximize2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export function RecipeTree() {
  const [zoom, setZoom] = useState(100)

  const handleZoomIn = () => {
    setZoom((prev) => Math.min(prev + 10, 150))
  }

  const handleZoomOut = () => {
    setZoom((prev) => Math.max(prev - 10, 50))
  }

  return (
    <Card className="w-full h-[500px] shadow-lg border-amber-800/50 flex flex-col bg-card/90 backdrop-blur-sm glow-border light-theme">
      <CardHeader className="bg-gradient-to-r from-amber-950/50 to-amber-900/50 rounded-t-lg border-b border-amber-800/50">
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-xl text-amber-300 glow-text font-cinzel">Recipe Visualization</CardTitle>
            <CardDescription className="text-amber-200/70">
              Visual representation of element combinations
            </CardDescription>
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8 border-amber-800/50 bg-secondary/50"
              onClick={handleZoomOut}
            >
              <ZoomOut className="h-4 w-4 text-amber-300" />
            </Button>
            <span className="text-xs text-amber-300 w-12 text-center">{zoom}%</span>
            <Button
              variant="outline"
              size="icon"
              className="h-8 w-8 border-amber-800/50 bg-secondary/50"
              onClick={handleZoomIn}
            >
              <ZoomIn className="h-4 w-4 text-amber-300" />
            </Button>
            <Button variant="outline" size="icon" className="h-8 w-8 border-amber-800/50 bg-secondary/50">
              <Maximize2 className="h-4 w-4 text-amber-300" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex-1 p-0 relative overflow-hidden">
        <div className="absolute inset-0 flex items-center justify-center bg-amber-50/95 p-6 overflow-auto">
          <div className="w-full h-full flex items-center justify-center" style={{ transform: `scale(${zoom / 100})` }}>
            <div className="text-center p-8 rounded-lg border-2 border-dashed border-amber-200">
              <p className="text-amber-900">Search for an element to visualize its recipe tree</p>
              {/* Tree visualization will be rendered here */}
              <div className="mt-4 opacity-50">
                <svg width="200" height="150" viewBox="0 0 200 150" className="mx-auto">
                  <g transform="translate(100,20)">
                    <circle cx="0" cy="0" r="15" fill="#d97706" />
                    <text x="0" y="4" textAnchor="middle" fill="white" fontSize="10">
                      Fire
                    </text>

                    <line x1="-30" y1="30" x2="-5" y2="5" stroke="#92400e" strokeWidth="1.5" />
                    <circle cx="-40" cy="40" r="12" fill="#b45309" />
                    <text x="-40" y="44" textAnchor="middle" fill="white" fontSize="8">
                      Water
                    </text>

                    <line x1="30" y1="30" x2="5" y2="5" stroke="#92400e" strokeWidth="1.5" />
                    <circle cx="40" cy="40" r="12" fill="#b45309" />
                    <text x="40" y="44" textAnchor="middle" fill="white" fontSize="8">
                      Earth
                    </text>

                    <line x1="-40" y1="60" x2="-40" y2="52" stroke="#92400e" strokeWidth="1.5" />
                    <circle cx="-40" cy="70" r="10" fill="#d97706" opacity="0.8" />
                    <text x="-40" y="73" textAnchor="middle" fill="white" fontSize="6">
                      Steam
                    </text>

                    <line x1="40" y1="60" x2="40" y2="52" stroke="#92400e" strokeWidth="1.5" />
                    <circle cx="40" cy="70" r="10" fill="#d97706" opacity="0.8" />
                    <text x="40" y="73" textAnchor="middle" fill="white" fontSize="6">
                      Lava
                    </text>
                  </g>
                </svg>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
