package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	words := []string{
		"time", "person", "year", "way", "day", "thing", "man", "world", "life", "hand",
		"part", "child", "eye", "woman", "place", "work", "week", "case", "point", "government",
		"company", "number", "group", "problem", "fact", "be", "have", "do", "say", "get",
		"make", "go", "know", "take", "see", "come", "think", "look", "want", "give",
		"use", "find", "tell", "ask", "work", "seem", "feel", "try", "leave", "call",
	}

	// Build a test of 10 random words
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
	targetStr := strings.Join(words[:10], " ")

	// Init screen
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	defer s.Fini()
	s.Clear()

	// Styles
	base := tcell.StyleDefault.Foreground(tcell.ColorGray)                // light (untouched/backspaced)
	correct := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true) // darker/stronger for correct
	wrong := tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)     // wrong in red
	hud := tcell.StyleDefault.Foreground(tcell.ColorLightSlateGray)       // HUD text
	cursorStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Underline(true)

	targetRunes := []rune(targetStr)
	state := make([]int, len(targetRunes)) // 0: untyped, 1: correct, 2: wrong
	cursor := 0

	start := time.Now()

	render := func(done bool) {
		s.Clear()
		w, _ := s.Size()

		// Title
		title := "monkey-term (ESC to quit)"
		drawString(s, 0, 0, truncate(title, w), hud)

		// Text line (y=2)
		y := 2
		x := 0
		for i, r := range targetRunes {
			st := base
			switch state[i] {
			case 1:
				st = correct
			case 2:
				st = wrong
			default:
				st = base
			}
			// Show cursor underline at current position if not finished
			if i == cursor && !done {
				st = cursorStyle
			}
			s.SetContent(x, y, r, nil, st)
			x++
			// Wrap if necessary
			if x >= w {
				y++
				x = 0
			}
		}

		// HUD Footer
		elapsed := time.Since(start).Seconds()
		if elapsed < 0.001 {
			elapsed = 0.001
		}
		correctCount := 0
		for _, v := range state {
			if v == 1 {
				correctCount++
			}
		}
		// crude WPM: correct chars / 5 / minutes
		wpm := float64(correctCount) / 5.0 / (elapsed / 60.0)

		info := fmt.Sprintf("progress: %d/%d  |  accuracy: %d%%  |  wpm: %.1f   |  elapsed: %.0fs",
			cursor, len(targetRunes),
			int(accuracyPercent(state)*100+0.5),
			wpm,
			elapsed,
		)
		drawString(s, 0, y+2, truncate(info, w), hud)

		if done {
			msg := "Done! Press ENTER to exit or ESC to quit."
			drawString(s, 0, y+4, truncate(msg, w), hud)
		}

		s.Show()
	}

	done := false
	render(done)

	// Event loop
	for {
		ev := s.PollEvent()
		switch tev := ev.(type) {
		case *tcell.EventKey:
			switch tev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyEnter:
				// If finished, allow exit on Enter
				if done {
					return
				}
				// otherwise ignore
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if cursor > 0 && !done {
					cursor--
					state[cursor] = 0 // back to "light"
				}
			default:
				if done {
					// ignore keys after finish unless Enter/ESC
					continue
				}
				r := tev.Rune()
				if r == 0 { // non-printables (arrows, etc.)
					continue
				}
				if cursor < len(targetRunes) {
					if r == targetRunes[cursor] {
						state[cursor] = 1 // correct
					} else {
						state[cursor] = 2 // wrong
					}
					cursor++
					if cursor >= len(targetRunes) {
						done = true
					}
				}
			}
			render(done)
		case *tcell.EventResize:
			s.Sync()
			render(done)
		}
	}
}

// Helpers

func drawString(s tcell.Screen, x, y int, str string, st tcell.Style) {
	for i, r := range str {
		s.SetContent(x+i, y, r, nil, st)
	}
}

func truncate(s string, width int) string {
	if width <= 0 || len([]rune(s)) <= width {
		return s
	}
	rs := []rune(s)
	return string(rs[:width])
}

func accuracyPercent(state []int) float64 {
	total := 0
	correct := 0
	for _, v := range state {
		if v == 1 || v == 2 {
			total++
		}
		if v == 1 {
			correct++
		}
	}
	if total == 0 {
		return 1.0
	}
	return float64(correct) / float64(total)
}
