package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func valid_link(link string) bool {
	return link != "/wiki/" && link != "/wiki/Main_Page" && link != "/wiki/Wikipedia" && link != "/wiki/Free_content" && link != "/wiki/Encyclopedia" && link != "/wiki/English_language" && !strings.Contains(link, ".") && !strings.Contains(link, ":")
}

func wiki_links() []string {
	page := request("https://en.wikipedia.org/wiki/Main_Page")
	re := regexp.MustCompile(`/wiki/[^"]*`)
	matches := re.FindAll([]byte(page), -1)

	var links []string
	for i := range matches {
		link := string(matches[i])
		if valid_link(link) {
			links = append(links, link)
		}
	}

	return links
}

func random_link(links []string) string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(links))
	fmt.Println(n, len(links))
	link := links[n]
	return link
}

func html_doc(url string) *goquery.Document {
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

func extract_paragraphs(doc *goquery.Document) string {
	// Find the review items
	var text strings.Builder

	b_re := regexp.MustCompile(`\s?\[[^\]]*\]`)
	p_re := regexp.MustCompile(`\s?\([^\)]*\)`)
	nl_re := regexp.MustCompile(`\n`)

	max_words := 300
	n_words := 0
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		paragraph := s.Text()
		isPar, _ := regexp.Match(`[A-Za-z]`, []byte(paragraph))
		if isPar && n_words < max_words {
			paragraph = b_re.ReplaceAllString(paragraph, "")
			paragraph = p_re.ReplaceAllString(paragraph, "")
			paragraph = nl_re.ReplaceAllString(paragraph, "")
			paragraph = strings.TrimSpace(paragraph)
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

func page_content(link string) string {
	url := "https://en.wikipedia.org/" + link
	doc := html_doc(url)

	paragraphs := extract_paragraphs(doc)
	return paragraphs
}

func wiki_words() string {
	links := wiki_links()
	link := random_link(links)
	text := page_content(link)

	return text
}