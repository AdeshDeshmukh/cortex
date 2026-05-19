package git

import (
	"bufio"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

var (
	diffHeaderRegex = regexp.MustCompile(`^diff --git a/(.*) b/(.*)$`)
	hunkHeaderRegex = regexp.MustCompile(`^@@ -(\d+),?\d* \+(\d+),?\d* @@`)
	fileTypeMap     = map[string]string{
		".go":   "go",
		".py":   "python",
		".js":   "javascript",
		".ts":   "typescript",
		".jsx":  "javascript",
		".tsx":  "typescript",
		".java": "java",
		".cpp":  "cpp",
		".c":    "c",
		".rs":   "rust",
		".rb":   "ruby",
		".php":  "php",
	}
)

type DiffParser struct {
	raw string
}

func NewDiffParser(rawDiff string) *DiffParser {
	return &DiffParser{raw: rawDiff}
}

func (p *DiffParser) Parse() ([]types.DiffChange, error) {
	if p.raw == "" {
		return nil, nil
	}

	var changes []types.DiffChange
	var currentChange *types.DiffChange

	scanner := bufio.NewScanner(strings.NewReader(p.raw))

	for scanner.Scan() {
		line := scanner.Text()

		if matches := diffHeaderRegex.FindStringSubmatch(line); matches != nil {
			if currentChange != nil {
				changes = append(changes, *currentChange)
			}

			filePath := matches[2]
			currentChange = &types.DiffChange{
				FilePath:   filePath,
				FileType:   detectFileType(filePath),
				ChangeType: "modified",
			}
			continue
		}

		if currentChange == nil {
			continue
		}

		if matches := hunkHeaderRegex.FindStringSubmatch(line); matches != nil {
			if len(matches) >= 3 {
				currentChange.StartLine = parseInt(matches[2])
			}
			continue
		}

		switch {
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			content := strings.TrimPrefix(line, "+")
			currentChange.AddedLines = append(currentChange.AddedLines, content)

		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			content := strings.TrimPrefix(line, "-")
			currentChange.RemovedLines = append(currentChange.RemovedLines, content)

		case strings.HasPrefix(line, " "):
			currentChange.Context += line + "\n"
		}
	}

	if currentChange != nil {
		changes = append(changes, *currentChange)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return changes, nil
}

func detectFileType(path string) string {
	ext := filepath.Ext(path)
	if fileType, ok := fileTypeMap[ext]; ok {
		return fileType
	}
	return "unknown"
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}