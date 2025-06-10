package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Box represents a terminal UI box with a title and content.
type Box struct {
	Title      string
	TopLine    string
	MiddleLine string
	BottomLine string
	Lines      []string
}

// splitLines splits a string into lines of a specified maximum length.
func splitLines(text string, maxLength int) []string {
	var lines []string
	for len(text) > maxLength {
		line := text[:maxLength]
		lines = append(lines, line)
		text = text[maxLength:]
	}
	if len(text) > 0 {
		lines = append(lines, text)
	}
	return lines
}

// CreateBox creates a Box with a title and content, ensuring the title and
// content are padded to the same width.
func CreateBox(boxData BoxData) Box {
	// TODO Configur the min wideth
	minWidth := 22
	maxLength := len(boxData.Title)
	for _, line := range boxData.Content {
		if len(line) > maxLength {
			maxLength = utf8.RuneCountInString(line)
		}
	}
	if maxLength < minWidth {
		maxLength = minWidth
	}
	var padded []string
	for _, line := range boxData.Content {
		paddedLine := fmt.Sprintf("║ %-*s ║", maxLength, line)
		padded = append(padded, paddedLine)
	}
	width := maxLength + 2 // 2 for 1 space on each side
	border := strings.Repeat("═", width)
	return Box{
		Title:      fmt.Sprintf("║ %-*s ║", maxLength, boxData.Title),
		TopLine:    fmt.Sprintf("╔%s╗", border),
		MiddleLine: fmt.Sprintf("╠%s╣", border),
		BottomLine: fmt.Sprintf("╚%s╝", border),
		Lines:      padded,
	}
}

// CombineBoxesSideBySide combines two Boxes side by side, ensuring they are
// padded to the same height and width.
func CombineBoxesSideBySide(left, right Box) Box {
	gap := "  " // Two spaces between the boxes
	maxLines := len(left.Lines)
	if len(right.Lines) > maxLines {
		maxLines = len(right.Lines)
	}
	paddedLeft := append([]string{}, left.Lines...)
	paddedRight := append([]string{}, right.Lines...)
	for len(paddedLeft) < maxLines {
		paddedLeft = append(
			paddedLeft,
			left.BottomLine,
		)
		left.BottomLine = strings.Repeat(
			" ",
			utf8.RuneCountInString(paddedLeft[0]),
		)
	}
	for len(paddedRight) < maxLines {
		paddedRight = append(
			paddedRight,
			right.BottomLine,
		)
		right.BottomLine = strings.Repeat(
			" ",
			utf8.RuneCountInString(paddedRight[0]),
		)
	}
	combinedLines := make([]string, maxLines)
	for i := 0; i < maxLines; i++ {
		combinedLines[i] = paddedLeft[i] + gap + paddedRight[i]
	}
	combinedBox := Box{
		Title:      left.Title + gap + right.Title,
		TopLine:    left.TopLine + gap + right.TopLine,
		MiddleLine: left.MiddleLine + gap + right.MiddleLine,
		BottomLine: left.BottomLine + gap + right.BottomLine,
		Lines:      combinedLines,
	}
	return combinedBox
}
