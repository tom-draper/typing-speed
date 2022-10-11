package cmd

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func request(url string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

type WikiMainPage struct {
	BatchComplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]struct {
			PageID int    `json:"pageid"`
			Ns     int    `json:"ns"`
			Title  string `json:"title"`
			Links  []struct {
				Ns    int    `json:"ns"`
				Title string `json:"title"`
			} `json:"links"`
		} `json:"pages"`
	} `json:"query"`
}

type WikiPage struct {
	BatchComplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]struct {
			PageID  int    `json:"pageid"`
			Ns      int    `json:"ns"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

func fetch(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func fetchMainPage(url string) WikiMainPage {
	body := fetch(url)

	var page WikiMainPage
	json.Unmarshal(body, &page) // Parse JSON

	return page
}

func fetchPage(url string) WikiPage {
	body := fetch(url)

	var page WikiPage
	json.Unmarshal(body, &page) // Parse JSON

	return page
}

func validLink(link string) bool {
	return link != "/wiki/" && link != "/wiki/Main_Page" && link != "/wiki/Wikipedia" && link != "/wiki/Free_content" && link != "/wiki/Encyclopedia" && link != "/wiki/English_language" && !strings.Contains(link, ".") && !strings.Contains(link, ":")
}

func getWikiLinksOld() []string {
	page := request("https://en.wikipedia.org/wiki/Main_Page")
	re := regexp.MustCompile(`/wiki/[^"]*`)
	matches := re.FindAll([]byte(page), -1)

	seen := make(map[string]struct{})
	var links []string
	for i := range matches {
		link := string(matches[i])
		if _, ok := seen[link]; !ok && validLink(link) {
			links = append(links, link)
			seen[link] = struct{}{}
		}
	}

	return links
}

func randomLink(links []string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(links))
	link := links[n]
	return link
}

func htmlDoc(url string) *goquery.Document {
	// Request the HTML page
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	return doc
}

func cleanParagraph(paragraph string) string {
	b_re := regexp.MustCompile(`\s?\[[^\]]*\]`)
	p_re := regexp.MustCompile(`\s?\([^\)]*\)`)
	nl_re := regexp.MustCompile(`\n`)
	nonascii_re := regexp.MustCompile(`[^\x00-\x7F]+`)
	ws_re := regexp.MustCompile(`\s{2,}`)
	// Remove references
	paragraph = b_re.ReplaceAllString(paragraph, "")
	// Remove text within parenthesis
	paragraph = p_re.ReplaceAllString(paragraph, "")
	// Remove newlines wordwrap will insert these
	paragraph = nl_re.ReplaceAllString(paragraph, "")
	// Remove any characters outside of the ascii set
	paragraph = nonascii_re.ReplaceAllString(paragraph, " ")
	// Remove any double+ spaces created by removals
	paragraph = ws_re.ReplaceAllString(paragraph, "")
	// Trim spaces at before and end of paragraph
	paragraph = strings.TrimSpace(paragraph)
	return paragraph
}

func extractParagraphs(doc *goquery.Document) string {
	// Find the review items
	var text strings.Builder

	max_words := 300
	n_words := 0
	doc.Find("p").Each(func(_ int, s *goquery.Selection) {
		paragraph := s.Text()
		// Check if paragraph is empty
		isPar, _ := regexp.Match(`[A-Za-z]`, []byte(paragraph))
		if isPar && n_words < max_words {
			paragraph = cleanParagraph(paragraph)
			words := strings.Split(paragraph, " ")
			for i := range words {
				text.WriteString(words[i])
				text.WriteString(" ")
				n_words++
				if n_words >= max_words {
					break
				}
			}
		}
	})

	return text.String()
}

func pageContentOld(link string) string {
	url := "https://en.wikipedia.org/" + link
	doc := htmlDoc(url)

	paragraphs := extractParagraphs(doc)
	return paragraphs
}
func pageContent(link string) string {
	url := "https://en.wikipedia.org/w/api.php?format=xml&action=query&prop=extracts&format=json&redirects=true&titles=" + strings.ReplaceAll(link, " ", "%20")

	page := fetchPage(url)

	var paragraphs string
	for pageID := range page.Query.Pages {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.Query.Pages[pageID].Extract))
		if err != nil {
			panic(err)
		}
		paragraphs = extractParagraphs(doc)
	}

	return paragraphs
}

func WikiWords2(config Config) string {
	links := config.wikiLinks
	if links == nil {
		// If haven't requested wiki links before
		links = getWikiLinksOld()
		config.wikiLinks = links
	}
	link := randomLink(links)
	text := pageContentOld(link)
	return text
}

func extractMainPageLinks(page WikiMainPage) []string {
	var links []string
	for _, link := range page.Query.Pages["15580374"].Links {
		links = append(links, strings.TrimSpace(link.Title))
	}
	return links
}

func validLink2(link string) bool {
	return !strings.Contains(link, "Wikipedia") && !strings.Contains(link, "Template") && !strings.Contains(link, "Help") && !strings.Contains(link, "Portal")
}

func filterLinks(links []string) []string {
	var filteredLinks []string
	for _, link := range links {
		if validLink2(link) {
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}

func wikiLinks() []string {
	url := "https://en.wikipedia.org/w/api.php?action=query&titles=Main_Page&prop=links&format=json&pllimit=max"
	page := fetchMainPage(url)
	links := extractMainPageLinks(page)
	links = filterLinks(links)
	return links
}

func WikiWords(config Config) string {
	links := config.wikiLinks
	if links == nil {
		// If haven't requesed wiki links before
		links = wikiLinks()
		config.wikiLinks = links
	}
	link := randomLink(links)
	text := pageContent(link)
	return text
}

func readWordsFile(path string) [1000]string {
	// Open the file.
	f, _ := os.Open(path)
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.
	var words [1000]string
	for i := 0; i < 1000; i++ {
		if scanner.Scan() {
			line := scanner.Text()
			words[i] = line
		}
	}

	return words
}

func shuffleWords(words [1000]string) [1000]string {
	loc_map := rand.Perm(len(words)) // Get a new location for each word

	var shuffled_words [1000]string
	for i, loc := range loc_map {
		shuffled_words[loc] = words[i]
	}
	return shuffled_words
}

func CommonWords(path string) string {
	words := readWordsFile(path)
	words = shuffleWords(words)
	text := strings.Join(words[:], " ")
	return text
}
