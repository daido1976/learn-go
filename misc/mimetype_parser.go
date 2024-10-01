package main

import (
	"bufio"
	"regexp"
	"strings"
)

// Parse takes the output of the 'file --mime-type' command and returns a map of file names to their MIME types.
// See. https://note.com/knowledgework/n/nc4c0a24a9569
func parseMimeType(output string) map[string]string {
	// Create a map to store the parsed result
	result := make(map[string]string)

	// Regular expression to match MIME types (assuming they follow the given pattern)
	re := regexp.MustCompile(`[a-z0-9-_.]+/[a-z0-9-_.]+`)

	// Use a scanner to process the output line by line
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line at the first colon (":") followed by space, which separates the filename from the mime type
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			filename := parts[0]
			mimeType := parts[1]
			// Validate MIME type using the regular expression
			if re.MatchString(mimeType) {
				result[filename] = mimeType
			}
		}
	}

	return result
}
