import { Button } from "@/components/ui/button"
import React, { useRef, useState } from "react"

const ScrapingButton = () => {
  const startSoundRef = useRef<HTMLAudioElement | null>(null)
  const successSoundRef = useRef<HTMLAudioElement | null>(null)
  const [isScraping, setIsScraping] = useState(false)

  const handleScrape = async () => {
    const userConfirmed = window.confirm("This will start scraping. Are you sure?")
    if (!userConfirmed) return

    setIsScraping(true)

    if (startSoundRef.current) {
      startSoundRef.current.loop = true
      startSoundRef.current.play().catch(e => console.error("Audio play failed:", e))
    }

    try {
      const response = await fetch("http://tubes2-bobfirstsearch.up.railway.app/api/scrape", {
        method: "POST",
      })

      const data = await response.json()

      startSoundRef.current?.pause()
      startSoundRef.current!.currentTime = 0 

      if (data.message === "Scraping completed successfully") {
        successSoundRef.current?.play()
        alert("Scraping completed successfully!")
      } else {
        alert("Scraping failed!")
      }
    } catch (error) {
      startSoundRef.current?.pause()
      startSoundRef.current!.currentTime = 0
      console.error("Error triggering scrape:", error)
      alert("An error occurred while starting the scraping process.")
    } finally {
      setIsScraping(false)
    }
  }

  return (
    <>
      {}
      <audio 
        ref={startSoundRef} 
        src="/sounds/BBtheme.mp3" 
        preload="auto"
      />
      <audio 
        ref={successSoundRef} 
        src="/sounds/YeahScience.mp3" 
        preload="auto" 
      />
      
      <Button 
        variant="scrape" 
        onClick={handleScrape}
        disabled={isScraping}
        className="w-full bg-gradient-to-r from-amber-600 to-yellow-500 hover:from-amber-500 hover:to-yellow-400 text-background font-medium shadow-[0_0_10px_rgba(217,119,6,0.3)]"
      >
        {isScraping ? "Scraping in progress..." : "Start Scraping"}
      </Button>
    </>
  )
}

export default ScrapingButton