"use client"

import React, { useState } from "react"
import Image from "next/image"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Switch } from "@/components/ui/switch"
import { Card, CardContent } from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

type Props = {
  onSearch: (
    element: string,
    alg: "bfs" | "dfs",
    mode: "shortest" | "multiple",
    maxRecipes: number
  ) => void
}

export function SearchPanel({ onSearch }: Props) {
  const [element, setElement] = useState("")
  const [alg, setAlg] = useState<"bfs" | "dfs">("bfs")
  const [mode, setMode] = useState<"shortest" | "multiple">("shortest")
  const [maxRecipes, setMaxRecipes] = useState("5")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    console.log("SearchPanel.handleSubmit, element =", element, "alg=", alg, "mode=", mode, "maxRecipes=", maxRecipes)
    onSearch(element, alg, mode, parseInt(maxRecipes, 10))
  }

  const handleElementChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setElement(e.target.value)
  }

  return (
    <Card className="w-full shadow-lg border-amber-800/50 bg-card/90 backdrop-blur-sm glow-border">
      <CardContent className="pt-6">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid gap-6 md:grid-cols-3">
            <div className="space-y-2 md:col-span-1">
              <Label htmlFor="element" className="text-amber-200">
                Element to search
              </Label>
              <Input
                id="element"
                placeholder="Enter element name..."
                className="mr-2 text-amber-300"
                value={element}
                onChange={handleElementChange}
              />
              {element && (
                <div className="mt-4 text-sm text-amber-200/70">
                  <p>Selected element: {element}</p>
                  <div className="relative w-24 h-24 overflow-hidden flex items-center justify-center">
                    <Image
                      src={`/elements/${element}.svg?height=96&width=96&text=${encodeURIComponent(
                        element
                      )}`}
                      alt={element}
                      width={96}
                      height={96}
                      className="object-contain"
                    />
                  </div>
                </div>
              )}
            </div>

            <div className="space-y-2 md:col-span-1">
              <Label htmlFor="algorithm" className="text-amber-200">
                Search algorithm
              </Label>
              <Select value={alg} onValueChange={(val: "bfs" | "dfs") => setAlg(val)}>
                <SelectTrigger id="algorithm" className="border-amber-800/50 bg-secondary/50">
                  <SelectValue placeholder="Select algorithm" />
                </SelectTrigger>
                <SelectContent className="bg-card border-amber-800/50">
                  <SelectItem value="bfs">Breadth-First Search (BFS)</SelectItem>
                  <SelectItem value="dfs">Depth-First Search (DFS)</SelectItem>
                </SelectContent>
              </Select>

              <div className="space-y-4 mt-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label className="text-amber-200">Search mode</Label>
                    <p className="text-sm text-amber-200/70">
                      Shortest path or multiple recipes
                    </p>
                  </div>
                  <RadioGroup
                    defaultValue={mode}
                    onValueChange={(val: "shortest" | "multiple") => setMode(val)}
                    className="flex space-x-4">
                      <div className="flex items-center space-x-2">
                        <RadioGroupItem value="shortest" id="shortest" />
                        <Label htmlFor="shortest" className="cursor-pointer text-amber-200">
                          Shortest
                        </Label>
                      </div>
                      <div className="flex items-center space-x-2">
                        <RadioGroupItem value="multiple" id="multiple" />
                        <Label htmlFor="multiple" className="cursor-pointer text-amber-200">
                          Multiple
                        </Label>
                      </div>
                  </RadioGroup>
                </div>
              </div>
            </div>

            <div className="space-y-2 md:col-span-1">
              {mode === "multiple" && (
                <div className="flex items-center justify-between pl-4 pr-6 py-3 bg-secondary/30 rounded-lg border border-amber-800/30">
                  <Label htmlFor="maxRecipes" className="text-amber-200">
                    # of Recipes
                  </Label>
                  <Select
                    value={maxRecipes}
                    onValueChange={(val) => setMaxRecipes(val)}
                  >
                    <SelectTrigger className="border-amber-800/50 bg-secondary/50 w-[120px]">
                      <SelectValue placeholder="Select max" />
                    </SelectTrigger>
                    <SelectContent className="bg-card border-amber-800/50">
                      {Array.from({ length: 10 }, (_, i) => (
                        <SelectItem key={i + 1} value={`${i + 1}`}>
                          {i + 1}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )}

              <div className="flex items-center justify-between pl-4 pr-6 py-3 bg-secondary/30 rounded-lg border border-amber-800/30">
                <div className="space-y-0.5">
                  <Label htmlFor="visualize" className="text-amber-200">
                    Visualize search process
                  </Label>
                  <p className="text-sm text-amber-200/70">
                    Show algorithm steps in real-time
                  </p>
                </div>
                <Switch id="visualize" />
              </div>

              <Button
                type="submit"
                disabled={!element.trim()}
                className="w-full bg-gradient-to-r from-amber-600 to-yellow-500 hover:from-amber-500 hover:to-yellow-400 text-background font-medium shadow-[0_0_10px_rgba(217,119,6,0.3)]"
                onClick={() => new Audio('sounds/Jesse, We Have To Cook - Sound Effect (Breaking Bad) (mp3cut.net).mp3').play()}
              >
                Search Recipes
              </Button>
            </div>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
