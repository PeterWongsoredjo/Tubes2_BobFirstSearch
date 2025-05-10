"use client"

import { useState, useEffect, useRef } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardHeader, CardTitle } from "@/components/ui/card"

import type { Network as NetworkType, Node, Edge } from "vis-network"
import { DataSet } from "vis-data"

export function RecipeTree({
    root,
    alg,
    mode,
    maxRecipes,
}: {
    root: string
    alg: "bfs"|"dfs"
    mode: "shortest"|"multiple"
    maxRecipes: number
  }) {
    const [zoom, setZoom] = useState(1)
    const containerRef = useRef<HTMLDivElement>(null)
    const networkRef = useRef<NetworkType | null>(null)

    const handleZoomIn = () => setZoom(z => Math.min(z + 0.1, 2))
    const handleZoomOut = () => setZoom(z => Math.max(z - 0.1, 0.5))

    useEffect(() => {
      console.log("🌳 RecipeTree mounted with props:", { root, alg, mode, maxRecipes })
      if (!root) return
      if (!containerRef.current) return

      const url = new URL("http://localhost:8080/api/tree")
      url.searchParams.set("root", root)
      url.searchParams.set("alg", alg)
      url.searchParams.set("mode", mode)
      if (mode === "multiple") {
        url.searchParams.set("max", String(maxRecipes))
      }

      fetch(url.toString())
        .then(r => {
          if (!r.ok) throw new Error(`HTTP ${r.status}`)
          return r.json()
        })
        .then(async (graph: { nodes: Node[]; edges: Edge[] }) => {
          const nodes = new DataSet<Node>(graph.nodes)
          const edges = new DataSet<Edge>(graph.edges)

          const options = {
            layout: { hierarchical: { direction: "UD", sortMethod: "directed" } },
            interaction: { zoomView: true, dragView: true },
            physics: false,
          }

          const { Network } = await import("vis-network")

          if (!networkRef.current) {
            networkRef.current = new Network(
              containerRef.current!,
              { nodes, edges },
              options
            )
          } else {
            networkRef.current.setData({ nodes, edges })
          }

          networkRef.current.moveTo({ scale: zoom })
        })
        .catch(err => {
          console.error("Failed to load tree:", err)
        })
    }, [root, alg, mode, maxRecipes])

    useEffect(() => {
      networkRef.current?.moveTo({ scale: zoom })
    }, [zoom])

    return (
      <Card className="w-full h-[500px] shadow-lg border-amber-800/50 flex flex-col bg-card/90 backdrop-blur-sm">
        <CardHeader className="flex items-center justify-between bg-amber-900/50 p-2 rounded-t-lg">
          <CardTitle className="text-amber-300">Recipe Visualization</CardTitle>
        </CardHeader>
        <div ref={containerRef} className="flex-1" />
      </Card>
    )
  }
