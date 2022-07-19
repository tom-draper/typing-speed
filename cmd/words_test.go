package cmd

import "testing"

func TestWikiLinks(t *testing.T) {
	links := wikiLinks()
	stringSet := make(map[string]struct{})
	for _, link := range links {
		if _, ok := stringSet[link]; !ok {
			t.Errorf("fail found duplicate wiki link: %s", link)
		}
		stringSet[link] = struct{}{}
	}
}
