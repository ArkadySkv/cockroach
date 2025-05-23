// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package faker

import (
	"math/rand/v2"
	"strings"
	"unicode"
)

type loremFaker struct {
	words *weightedEntries
}

// Words returns the requested number of random en_US words.
func (f *loremFaker) Words(rng *rand.Rand, num int) []string {
	w := make([]string, num)
	for i := range w {
		w[i] = f.words.Rand(rng).(string)
	}
	return w
}

// Sentences returns the requested number of random en_US sentences.
func (f *loremFaker) Sentences(rng *rand.Rand, num int) []string {
	s := make([]string, num)
	for i := range s {
		var b strings.Builder
		numWords := randInt(rng, 4, 8)
		for j := 0; j < numWords; j++ {
			word := f.words.Rand(rng).(string)
			if j == 0 {
				word = firstToUpper(word)
			}
			b.WriteString(word)
			if j == numWords-1 {
				b.WriteString(`.`)
			} else {
				b.WriteString(` `)
			}
		}
		s[i] = b.String()
	}
	return s
}

func firstToUpper(s string) string {
	isFirst := true
	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToUpper(r)
		}
		return r
	}, s)
}

// Paragraph returns a random en_US paragraph.
func (f *loremFaker) Paragraph(rng *rand.Rand) string {
	return strings.Join(f.Sentences(rng, randInt(rng, 1, 5)), ` `)
}

func newLoremFaker() loremFaker {
	f := loremFaker{}
	f.words = words()
	return f
}
