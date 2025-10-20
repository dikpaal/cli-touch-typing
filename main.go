package main

import (
	"fmt"
	"math/rand"
)

func main() {

	words := []string{
		"time", "person", "year", "way", "day", "thing", "man", "world", "life", "hand",
		"part", "child", "eye", "woman", "place", "work", "week", "case", "point", "government",
		"company", "number", "group", "problem", "fact", "be", "have", "do", "say", "get",
		"make", "go", "know", "take", "see", "come", "think", "look", "want", "give",
		"use", "find", "tell", "ask", "work", "seem", "feel", "try", "leave", "call",
	}

	// shuffle the slice
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })

	randomWords := words[:10]

	fmt.Println(randomWords)
}
