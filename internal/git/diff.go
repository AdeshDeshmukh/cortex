package git

import (
	"strings"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

// ParseDiff parses git diff output and returns a slice of DiffChange
func ParseDiff(diffOutput string) []types.DiffChange {
	var changes []types.DiffChange
	lines := strings.Split(diffOutput, "\n")

	var currentChange *types.DiffChange
	var inHunk bool

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			if currentChange != nil {
				changes = append(changes, *currentChange)
			}
			currentChange = &types.DiffChange{}
			inHunk = false
		} else if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
			if currentChange != nil {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					currentChange.FilePath = strings.TrimPrefix(parts[1], "a/")
				}
			}
		} else if strings.HasPrefix(line, "@@") {
			inHunk = true
		} else if inHunk {
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				if currentChange != nil {
					currentChange.AddedLines = append(currentChange.AddedLines, strings.TrimPrefix(line, "+"))
				}
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				if currentChange != nil {
					currentChange.RemovedLines = append(currentChange.RemovedLines, strings.TrimPrefix(line, "-"))
				}
			}
		}
	}

	if currentChange != nil {
		changes = append(changes, *currentChange)
	}

	return changes
}
