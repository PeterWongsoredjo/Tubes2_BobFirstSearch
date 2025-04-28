package scraping

import (
    "fmt"
    "net/http"
	"time"
	"strings"
    "github.com/PuerkitoBio/goquery"
)

type Recipe struct {
    Result    string   // Hasil kombinasi
    Components []string // Komponen yang harus dikombinasikan
}

type Element struct {
    Name string
    URL  string
}

func ScrapeMythAndMonstersElements() ([]string, error) {
	mythAndMonstersUrl := "https://little-alchemy.fandom.com/wiki/Category:Myths_and_Monsters"
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Get(mythAndMonstersUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Myth and Monsters page: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("page returned status: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse page: %w", err)
	}

	var elements []string

	doc.Find(".category-page__member-link").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		if name != "" {
			elements = append(elements, name)
		}
	})

	return elements, nil
}

func BuildMythAndMonstersBlacklist(elements []string) map[string]bool {
	blacklist := make(map[string]bool)
	for _, elem := range elements {
		blacklist[elem] = true
	}
	return blacklist
}

// ScrapeElementList meng-scrape semua element dari halaman Little Alchemy 2
func ScrapeElementList() ([]Element, error) {
    url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
    baseURL := "https://little-alchemy.fandom.com"

    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    res, err := client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch page: %w", err)
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("page returned status: %d", res.StatusCode)
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse HTML: %w", err)
    }

    var elements []Element
	seen := make(map[string]bool)
    
    myths, err := ScrapeMythAndMonstersElements()
    if err != nil {
        panic(err)
    }

    blacklist := BuildMythAndMonstersBlacklist(myths)

    doc.Find("a").Each(func(i int, s *goquery.Selection) {
        href, exists := s.Attr("href")
        title, titleExists := s.Attr("title")

        if exists && titleExists && strings.HasPrefix(href, "/wiki/") {
            if !strings.Contains(href, ":") { 
				if !blacklist[title] && !seen[title] && title != "Elements (Little Alchemy 1)" && title != "Elements (Little Alchemy 2)" && title != "Elements (Myths and Monsters)"{ // Cek kalau belum pernah
                    seen[title] = true
                    fullURL := baseURL + href
                    elements = append(elements, Element{
                        Name: title,
                        URL:  fullURL,
                    })
					fmt.Printf("Scraping element: %s\n", title)
                }
            }
        }
    })

    return elements, nil
}

func ScrapeElementPage(elem Element, validElements map[string]bool) ([]Recipe, error) {
    client := &http.Client{
        Timeout: 20 * time.Second,
    }

    res, err := client.Get(elem.URL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch element page %s: %w", elem.Name, err)
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("element page returned status: %d", res.StatusCode)
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse element page: %w", err)
    }

    var recipes []Recipe

    foundLA2 := false
    doc.Find("h3").EachWithBreak(func(i int, s *goquery.Selection) bool {
        headline := s.Find(".mw-headline").Text()
        if strings.TrimSpace(headline) == "Little Alchemy 2" {
            foundLA2 = true

            ul := s.NextAllFiltered("ul").First()
            if ul != nil {
                ul.Find("li").Each(func(j int, li *goquery.Selection) {
                    var components []string
                    li.Find("a").Each(func(k int, a *goquery.Selection) {
                        name := strings.TrimSpace(a.Text())
                        if name != "" {
                            components = append(components, name)
                        }
                    })

                    if len(components) == 2 {
						if validElements[components[0]] && validElements[components[1]] {
							recipes = append(recipes, Recipe{
								Result:     elem.Name,
								Components: components,
							})
						}
                    }
                })
            }
            return false
        }
        return true
    })

    if !foundLA2 {
        return nil, fmt.Errorf("no Little Alchemy 2 recipe found")
    }

    return recipes, nil
}

func ScrapeAllRecipes(elements []Element) ([]Recipe, error) {
    myths, err := ScrapeMythAndMonstersElements()
    if err != nil {
        panic(err)
    }

    blacklist := BuildMythAndMonstersBlacklist(myths)

    validElements := make(map[string]bool)
    for _, elem := range elements {
        validElements[elem.Name] = true
    }

    for elem := range blacklist {
        validElements[elem] = true
    }

    var allRecipes []Recipe

    for _, elem := range elements {
        recipes, err := ScrapeElementPage(elem, validElements)
        if err != nil {
            fmt.Printf("Warning: gagal scrape %s (%v), lanjutkan...\n", elem.Name, err)
            continue
        }
        allRecipes = append(allRecipes, recipes...)
    }

    return allRecipes, nil
}

