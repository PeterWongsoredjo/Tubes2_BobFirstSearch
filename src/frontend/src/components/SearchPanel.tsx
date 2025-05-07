"use client"

import type React from "react"

import { useState } from "react"
import { Search } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Switch } from "@/components/ui/switch"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

export function SearchPanel() {
  const [searchMode, setSearchMode] = useState("shortest")
  const [maxRecipes, setMaxRecipes] = useState("5")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // Handle search submission
  }

  return (
    <Card className="w-full shadow-lg border-amber-800/50 bg-card/90 backdrop-blur-sm glow-border">
      <CardHeader className="bg-gradient-to-r from-amber-950/50 to-amber-900/50 rounded-t-lg border-b border-amber-800/50">
        <CardTitle className="text-2xl text-amber-300 glow-text font-cinzel">Recipe Search</CardTitle>
        <CardDescription className="text-amber-200/70">
          Find the perfect recipe combinations for any element
        </CardDescription>
      </CardHeader>
      <CardContent className="pt-6">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid gap-6 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="element" className="text-amber-200">
                Element to search
              </Label>
              <div className="relative">
                <Search className="absolute left-3 top-2.5 h-4 w-4 text-amber-300" />
                <Input
                  id="element"
                  placeholder="Enter element name..."
                  className="pl-9 border-amber-800/50 bg-secondary/50"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="algorithm" className="text-amber-200">
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
            </div>
          </div>

          <div className="space-y-4">
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
                    <SelectItem value="3">3</SelectItem>
                    <SelectItem value="5">5</SelectItem>
                    <SelectItem value="10">10</SelectItem>
                    <SelectItem value="20">20</SelectItem>
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
          </div>

          <Button
            type="submit"
            className="w-full bg-gradient-to-r from-amber-600 to-yellow-500 hover:from-amber-500 hover:to-yellow-400 text-background font-medium shadow-[0_0_10px_rgba(217,119,6,0.3)]"
          >
            <Search className="mr-2 h-4 w-4" />
            Search Recipes
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
