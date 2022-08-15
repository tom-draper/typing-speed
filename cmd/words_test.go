package cmd

import (
	"strings"
	"testing"
)

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

func TestCommonWords(t *testing.T) {
	commonWordsStr := CommonWords("../words/common_words.txt")
	commonWords := strings.Split(commonWordsStr, " ")
	if commonWords[0] == "" || commonWords[len(commonWords)-1] == "" {
		t.Errorf("common word is empty")
	}
	if len(commonWords) != 1000 {
		t.Errorf("common words not 1000 in length (%d)", len(commonWords))
	}
	seen := make(map[string]struct{})
	for _, word := range commonWords {
		if _, ok := seen[word]; ok {
			t.Errorf("common words contains duplicate: %s", word)
		}
	}
}

func TestShuffleWords(t *testing.T) {
	words := readWordsFile("../words/common_words.txt")
	shuffledWords := shuffleWords(words)
	if words == shuffledWords {
		t.Errorf("shuffled words equal original words")
	}

}
