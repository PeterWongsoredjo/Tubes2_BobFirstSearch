"use client"

import type React from "react"

import { useState } from "react"
import Image from "next/image"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Switch } from "@/components/ui/switch"
import { Card, CardContent } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

export function SearchPanel() {
  const [searchMode, setSearchMode] = useState("shortest")
  const [maxRecipes, setMaxRecipes] = useState("5")
  const [selectedElement, setSelectedElement] = useState("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
  }

  const handleElementSelect = (element: string) => {
    setSelectedElement(element)
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
                  onChange={(e) => handleElementSelect(e.target.value)}
                />

                {selectedElement && (
                  <div className="mt-4 text-sm text-amber-200/70">
                    <p>Selected element: {selectedElement}</p>
                    <div className="relative w-24 h-24 overflow-hidden flex items-center justify-center">
                      {selectedElement ? (
                        <Image
                          src={`/elements/${selectedElement}.svg?height=96&width=96&text=${encodeURIComponent(selectedElement)}`}
                          alt={selectedElement}
                          width={96}
                          height={96}
                          className="object-contain"
                        />
                      ) : (
                        <div className="text-amber-200/50 text-xs text-center p-2">
                          Select an element to see its image
                        </div>
                        )}
                    </div>
                  </div>
                )}
            </div>

            <div className="space-y-2 md:col-span-1">
              <Label htmlFor="algorithm" className="text-amber-200 mr-2 h-4 w-4">
                Search algorithm
              </Label>
              <Select defaultValue="bfs">
                <SelectTrigger id="algorithm" className="border-amber-800/50 bg-secondary/50">
                  <SelectValue placeholder="Select algorithm" />
                </SelectTrigger>
                <SelectContent className="bg-card border-amber-800/50">
                  <SelectItem value="bfs">Breadth-First Search (BFS)</SelectItem>
                  <SelectItem value="dfs">Depth-First Search (DFS)</SelectItem>
                  <SelectItem value="bidirectional">Bidirectional Search</SelectItem>
                </SelectContent>
              </Select>
              <div className="space-y-4 mt-4">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label className="text-amber-200">Search mode</Label>
                  <p className="text-sm text-amber-200/70">Choose between shortest path or multiple recipes</p>
                </div>
                <div className="flex items-center space-x-2">
                  <RadioGroup defaultValue="shortest" className="flex space-x-4" onValueChange={setSearchMode}>
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
          </div>

          <div className="space-y-2 md:col-span-1">
            {searchMode === "multiple" && (
                  <div className="flex items-center justify-between pl-4 pr-6 py-3 bg-secondary/30 rounded-lg border border-amber-800/30">
                    <Label htmlFor="maxRecipes" className="text-amber-200">
                      Maximum number of recipes
                    </Label>
                    <Select value={maxRecipes} onValueChange={setMaxRecipes}>
                      <SelectTrigger id="maxRecipes" className="w-20 border-amber-800/50 bg-secondary/50">
                        <SelectValue placeholder="Max" />
                      </SelectTrigger>
                      <SelectContent className="bg-card border-amber-800/50">
                        <SelectItem value="1">1</SelectItem>
                        <SelectItem value="2">2</SelectItem>
                        <SelectItem value="3">3</SelectItem>
                        <SelectItem value="4">4</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
              )}
              <div className="flex items-center justify-between pl-4 pr-6 py-3 bg-secondary/30 rounded-lg border border-amber-800/30">
                <div className="space-y-0.5">
                  <Label htmlFor="visualize" className="text-amber-200">
                    Visualize search process
                  </Label>
                  <p className="text-sm text-amber-200/70">Show algorithm steps in real-time</p>
                </div>
                <Switch id="visualize" />
              </div>
              <Button
                  type="submit"
                  className="w-full bg-gradient-to-r from-amber-600 to-yellow-500 hover:from-amber-500 hover:to-yellow-400 text-background font-medium shadow-[0_0_10px_rgba(217,119,6,0.3)]"
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
