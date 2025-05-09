package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "runtime"
    "strconv"
    "strings"
    "github.com/PuerkitoBio/goquery"
)

type ElementTier struct {
    Name string `json:"name"`
    URL  string `json:"url"`
    Tier int    `json:"tier"`
}

func getConfigDir() string {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        log.Fatal("Could not get current file path")
    }
    projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
    return filepath.Join(projectRoot, "backend", "configs")
}

func ScrapeTierMap() ([]ElementTier, error) {
    const url = "https://little-alchemy.fandom.com/wiki/Elements_%28Little_Alchemy_2%29"

    // --- FETCH ---
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("HTTP GET: %w", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status %d", res.StatusCode)
    }

    // --- PARSE ---
    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return nil, fmt.Errorf("parse HTML: %w", err)
    }

    var items []ElementTier

    // 1) find all <h3> under the main content
    sel := doc.Find("div.mw-parser-output > h3")

    sel.Each(func(i int, h3 *goquery.Selection) {
        span := h3.Find("span.mw-headline")
        hdr := strings.TrimSpace(span.Text())

        id, _ := span.Attr("id")
        if !strings.HasPrefix(id, "Tier_") {
            return
        }

        parts := strings.Fields(hdr)
        if len(parts) < 2 {
            return
        }
        tierNum, err := strconv.Atoi(parts[1])
        if err != nil {
            return
        }

        sib := h3.Next()
        walked := 0
        for sib.Length() > 0 && goquery.NodeName(sib) != "table" {
            walked++
            sib = sib.Next()
        }
        if sib.Length() == 0 {
            return
        }

        table := sib
        rows := table.Find("tr")

        // 3) iterate rows, skip <th>
        rows.Each(func(j int, row *goquery.Selection) {
            if row.Find("th").Length() > 0 {
                return
            }
            cell := row.Find("td").First()
            a := cell.Find("a[title]").First()
            name := strings.TrimSpace(a.Text())
            href, _ := a.Attr("href")

            if name == "" || !strings.HasPrefix(href, "/wiki/") {
                return
            }

            items = append(items, ElementTier{
                Name: name,
                URL:  "https://little-alchemy.fandom.com" + href,
                Tier: tierNum,
            })
        })
    })

    fmt.Printf("Total elements scraped = %d\n", len(items))
    if len(items) == 0 {
        return nil, fmt.Errorf("no elements scraped — check the above logs for where it failed")
    }
    return items, nil
}

func saveTierMap(path string, tierMap map[string]int) error {
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        return err
    }
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")
    return enc.Encode(tierMap)
}

func main() {
    elements, err := ScrapeTierMap()
    if err != nil {
        log.Fatalf("error scraping tiers: %v", err)
    }

    tierMap := make(map[string]int, len(elements))
    for _, e := range elements {
        tierMap[e.Name] = e.Tier
    }

    cfgDir := getConfigDir()
    outPath := filepath.Join(cfgDir, "tiers.json")

    fmt.Printf("⮩ Saving %d elements to %s\n", len(elements), outPath)
    if err := saveTierMap(outPath, tierMap); err != nil {
        log.Fatalf("failed saving tiers.json: %v", err)
    }
}
