"use client"

import { useState, useEffect, useRef, useCallback } from "react"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"

import type { Network as NetworkType, Node, Edge } from "vis-network"
import { DataSet } from "vis-data"
///mport Image from "next/image"

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
    const [graphs, setGraphs] = useState<{ nodes: Node[]; edges: Edge[] }[]>([])
    const networkRefs = useRef<(NetworkType | null)[]>([])
    const containerRefs = useRef<(HTMLDivElement | null)[]>([])
    const setContainerRef = useCallback((index: number) => (node: HTMLDivElement | null) => {
      containerRefs.current[index] = node;
    }, []);
    
    useEffect(() => {
      console.log("🌳 RecipeTree mounted with props:", { root, alg, mode, maxRecipes })
      if (!root) return

      const url = new URL("http://localhost:8080/api/tree")
      url.searchParams.set("root", root)
      url.searchParams.set("alg", alg)
      url.searchParams.set("mode", mode)
      if (mode === "multiple") {
        url.searchParams.set("maxRecipes", String(maxRecipes))
      }

      fetch(url.toString())
        .then(r => {
          if (!r.ok) throw new Error(`HTTP ${r.status}`)
          return r.json()
        })
        .then(async (data: { graphs: { nodes: Node[]; edges: Edge[] }[] }) => {
          console.log("Received data:", data)
          if (!data.graphs || data.graphs.length === 0) {
            console.warn("No graphs received from API")
            return
          }
          
          setGraphs(data.graphs)
          
          // Initialize network refs array
          networkRefs.current = new Array(data.graphs.length).fill(null)
          containerRefs.current = new Array(data.graphs.length).fill(null)
          
          const { Network } = await import("vis-network")
          setTimeout(() => {
            data.graphs.forEach((graph, index) => {
              const container = containerRefs.current[index]
              if (!container) return
              
              // Enhance nodes with image information
              const enhancedNodes = graph.nodes.map(node => {
                const label = node.label || 'unknown';
                return {
                  ...node,
                  label,
                  shape: 'circularImage',
                  image: `/elements/${label}.svg`,
                  size: 30,
                  font: {
                    color: '#f5ce85',
                    size: 14,
                    face: 'Inter',
                    background: 'rgba(59, 40, 18, 0.8)',
                    strokeWidth: 2,
                    strokeColor: '#3b2812'
                  },
                  borderWidth: 2,
                  borderWidthSelected: 3,
                  color: {
                    border: '#8b5a28',
                    background: '#8b5a28',
                    highlight: {
                      border: '#f5ce85',
                      background: '#c18a36'
                    }
                  }
                };
              });
              
              const nodes = new DataSet<Node>(enhancedNodes)
              const edges = new DataSet<Edge>(graph.edges)
              
              const options = {
                layout: { 
                  hierarchical: { 
                    direction: "UD", 
                    sortMethod: "directed",
                    levelSeparation: 100,
                    nodeSpacing: 150 
                  } 
                },
                interaction: { zoomView: true, dragView: true },
                physics: false,                edges: {
                  color: {
                    color: '#8b5a28',
                    highlight: '#f5ce85'
                  },
                  width: 2,
                  smooth: {
                    enabled: true,
                    type: 'cubicBezier',
                    roundness: 0.5
                  }
                }
              };
              
              networkRefs.current[index] = new Network(
                container,
                { nodes, edges },
                options
              )
            })
          }, 100)
        })
        .catch(err => {
          console.error("Failed to load tree:", err)
        })
    }, [root, alg, mode, maxRecipes]);
    
    return (
      <Card className="w-full shadow-lg border-amber-800/50 flex flex-col bg-card/90 backdrop-blur-sm">
        <CardHeader className="bg-gradient-to-r from-amber-950/80 to-amber-900/50 p-2 rounded-t-lg border-b border-amber-700/30">
          <CardTitle className="text-amber-300 flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="mr-2">
              <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
              <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
            </svg>
            Recipe Visualization
          </CardTitle>
        </CardHeader>
        
        <CardContent className="p-4">
          {graphs.length === 0 ? (
            <div className="flex items-center justify-center h-[400px] text-amber-700">
              <div className="text-center">
                <svg className="animate-spin h-10 w-10 text-amber-600 mb-4 mx-auto" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <p>Looking for recipe chains...</p>
              </div>
            </div>
          ) : graphs.length === 1 ? (
            <div 
              ref={setContainerRef(0)} 
              className="h-[500px] w-full border border-amber-800/30 rounded-md bg-amber-950/30 shadow-inner"
            />
          ) : (
            // Multiple trees view - grid
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              {graphs.map((_, index) => (
                <div key={index} className="relative">
                  <div className="absolute top-2 left-2 bg-gradient-to-r from-amber-900/90 to-amber-800/90 text-amber-300 px-3 py-1 rounded-md text-xs shadow-md z-10 border border-amber-700/50">
                    Recipe Path {index + 1}
                  </div>
                  <div 
                    ref={setContainerRef(index)}
                    className="h-[350px] border border-amber-800/30 rounded-md bg-amber-950/30 shadow-inner"
                  />
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    )
  }
