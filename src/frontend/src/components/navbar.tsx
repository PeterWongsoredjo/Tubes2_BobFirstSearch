import { Sparkles } from "lucide-react"
import Link from "next/link"

export function Navbar() {
  return (
    <header className="sticky top-0 z-10 bg-background/80 backdrop-blur-md border-b border-amber-800/30 shadow-md">
      <div className="container mx-auto px-4 py-3">
        <div className="flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2">
            <div className="bg-gradient-to-r from-amber-500 to-yellow-600 p-1.5 rounded-lg shadow-[0_0_10px_rgba(217,119,6,0.5)]">
              <Sparkles className="h-5 w-5 text-background" />
            </div>
            <span className="font-cinzel font-semibold text-xl bg-gradient-to-r from-amber-300 to-yellow-200 bg-clip-text text-transparent glow-text">
              Alchemy Search
            </span>
          </Link>
          <nav className="hidden md:flex items-center gap-6">
            <Link href="/" className="text-amber-200 hover:text-amber-100 transition-colors">
              Home
            </Link>
            <Link href="/about" className="text-amber-200 hover:text-amber-100 transition-colors">
              About
            </Link>
            <Link href="/elements" className="text-amber-200 hover:text-amber-100 transition-colors">
              Elements
            </Link>
          </nav>
        </div>
      </div>
    </header>
  )
}
