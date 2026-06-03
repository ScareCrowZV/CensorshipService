package censor

import (
	"strings"
	"sync"
)

type Censor struct {
	bannedWords []string
	mu          sync.RWMutex
}

func NewCensor() *Censor {
	return &Censor{
		bannedWords: []string{"qwerty", "йцукен", "zxvbnm"},
	}
}

func (c *Censor) CheckText(text string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	lowerText := strings.ToLower(text)
	for _, word := range c.bannedWords {
		if strings.Contains(lowerText, word) {
			return false
		}
	}
	return true
}

func (c *Censor) AddBannedWord(word string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bannedWords = append(c.bannedWords, strings.ToLower(word))
}
