import type React from "react"
import type { Metadata } from "next"
import { Inter, Cinzel } from "next/font/google"
import "./globals.css"
import { ThemeProvider } from "@/components/theme-provider"

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" })
const cinzel = Cinzel({ subsets: ["latin"], variable: "--font-cinzel" })

export const metadata: Metadata = {
  title: "Little Alchemy 2 Recipe Search",
  description: "Find recipes and combinations for Little Alchemy 2 elements",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.variable} ${cinzel.variable} font-sans`}>
        <ThemeProvider attribute="class" defaultTheme="dark" enableSystem={false}>
          {children}
        </ThemeProvider>
      </body>
    </html>
  )
}
