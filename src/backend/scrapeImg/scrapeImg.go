package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
	"unicode"
    "github.com/PuerkitoBio/goquery"
)

func main() {
    url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
    fmt.Println("Fetching:", url)

    // Fetch the webpage
    resp, err := http.Get(url)
    if err != nil {
        log.Fatalf("Failed to fetch page: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        log.Fatalf("Status code error: %d %s", resp.StatusCode, resp.Status)
    }

    // Load the HTML document
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    // Create folder for images
    os.MkdirAll("images", os.ModePerm)

	downloadedImg := make(map[string]bool)

    // Select and download each image
    doc.Find("a.mw-file-description").Each(func(i int, s *goquery.Selection) {
        src, exists := s.Attr("href")
        if exists && strings.HasPrefix(src, "https") && !downloadedImg[src]{
            fmt.Printf("Downloading image %d: %s\n", i, src)
			downloadedImg[src] = true
            downloadImage(src)
        }
    })
	
}

// downloadImage downloads an image and saves it to the "images" folder
func downloadImage(url string) {
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Failed to download: %v", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Bad status: %d", resp.StatusCode)
        return
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Failed to read body: %v", err)
        return
    }

    // Try to extract the <title> tag for the correct name
    title := extractSVGTitle(string(body))
    if title == "" {
        title = "unknown"
    }

    filename := "images/" + title + ".svg"
    os.MkdirAll("images", os.ModePerm)

    err = os.WriteFile(filename, body, 0644)
    if err != nil {
        log.Printf("Failed to write file: %v", err)
    } else {
        log.Printf("Saved: %s", filename)
    }
}

func extractSVGTitle(svgContent string) string {
    start := strings.Index(svgContent, "<title>")
    end := strings.Index(svgContent, "</title>")
    if start != -1 && end != -1 && end > start {
        raw := strings.TrimSpace(svgContent[start+len("<title>") : end])
		return capitalizeWords(raw)
    }
    return ""
}

func capitalizeWords(s string) string {
    words := strings.Fields(s)
    for i, word := range words {
        if len(word) > 0 {
            r := []rune(word)
            r[0] = unicode.ToUpper(r[0])
            words[i] = string(r)
        }
    }
    return strings.Join(words, " ")
}