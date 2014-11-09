package testAPI

import (
	"bytes"
	"strings"
)

// YAMLContent is a simple test helper to make our tests look neater
func YAMLContent(content string) []byte {
	// remove the first blank lines if any
	contentLines := strings.Split(content, "\n")
	firstContentLine := 0
	for i, line := range contentLines {
		if strings.TrimSpace(line) == "" {
			// ignore the first empty lines
			continue
		} else {
			firstContentLine = i
			break
		}

	}
	contentLines = contentLines[firstContentLine:]

	// replace all tabs at the beginning of a line so indentation gets easier
	for i, line := range contentLines {
		modifiedLine := new(bytes.Buffer)
		for j, char := range line {
			if char != ' ' && char != '\t' {
				modifiedLine.WriteString(line[j:])
				break
			}
			if char == '\t' {
				modifiedLine.WriteString("    ")
			} else {
				modifiedLine.WriteRune(char)
			}
		}
		contentLines[i] = modifiedLine.String()
	}

	// find the indentation depth of the first line
	indentationDepth := 0
	for i, char := range contentLines[0] {
		if char != ' ' {
			indentationDepth = i
			break
		}
	}

	// remove indentation
	for i, line := range contentLines {
		if len(line) > indentationDepth {
			contentLines[i] = line[indentationDepth:]
		}
	}

	// clean uo blank lines at the end
	for i := len(contentLines) - 1; i >= 0; i-- {
		if strings.TrimSpace(contentLines[i]) == "" {
			contentLines[i] = ""
			continue
		} else {
			break
		}
	}

	return []byte(strings.Join(contentLines, "\n"))
}
