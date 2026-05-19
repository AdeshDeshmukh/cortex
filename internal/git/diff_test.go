package git

import (
	"testing"
)

const sampleDiff = `diff --git a/main.go b/main.go
index 1234567..abcdefg 100644
--- a/main.go
+++ b/main.go
@@ -10,7 +10,8 @@ import (
 func main() {
-	fmt.Println("old code")
+	fmt.Println("new code")
+	fmt.Println("additional line")
 }
`

const multiFileDiff = `diff --git a/main.go b/main.go
index 1234567..abcdefg 100644
--- a/main.go
+++ b/main.go
@@ -10,3 +10,4 @@
-	old
+	new
diff --git a/server.go b/server.go
index 7654321..gfedcba 100644
--- a/server.go
+++ b/server.go
@@ -5,2 +5,3 @@
-	old server
+	new server
`

func TestDiffParser_Parse(t *testing.T) {
	parser := NewDiffParser(sampleDiff)
	changes, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}

	change := changes[0]

	if change.FilePath != "main.go" {
		t.Errorf("expected FilePath 'main.go', got '%s'", change.FilePath)
	}

	if change.FileType != "go" {
		t.Errorf("expected FileType 'go', got '%s'", change.FileType)
	}

	if len(change.AddedLines) != 2 {
		t.Errorf("expected 2 added lines, got %d", len(change.AddedLines))
	}

	if len(change.RemovedLines) != 1 {
		t.Errorf("expected 1 removed line, got %d", len(change.RemovedLines))
	}

	expectedAdded := "	fmt.Println(\"new code\")"
	if change.AddedLines[0] != expectedAdded {
		t.Errorf("expected first added line '%s', got '%s'", expectedAdded, change.AddedLines[0])
	}
}

func TestDiffParser_EmptyDiff(t *testing.T) {
	parser := NewDiffParser("")
	changes, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(changes) != 0 {
		t.Errorf("expected 0 changes for empty diff, got %d", len(changes))
	}
}

func TestDiffParser_MultipleFiles(t *testing.T) {
	parser := NewDiffParser(multiFileDiff)
	changes, err := parser.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}

	if changes[0].FilePath != "main.go" {
		t.Errorf("expected first file 'main.go', got '%s'", changes[0].FilePath)
	}

	if changes[1].FilePath != "server.go" {
		t.Errorf("expected second file 'server.go', got '%s'", changes[1].FilePath)
	}
}

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"script.py", "python"},
		{"app.js", "javascript"},
		{"component.ts", "typescript"},
		{"App.jsx", "javascript"},
		{"Component.tsx", "typescript"},
		{"test.java", "java"},
		{"program.cpp", "cpp"},
		{"lib.c", "c"},
		{"main.rs", "rust"},
		{"script.rb", "ruby"},
		{"index.php", "php"},
		{"config.yml", "unknown"},
		{"README.md", "unknown"},
		{"Makefile", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := detectFileType(tt.path)
			if result != tt.expected {
				t.Errorf("detectFileType(%s) = %s, want %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"42", 42},
		{"999", 999},
		{"12abc34", 1234},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseInt(tt.input)
			if result != tt.expected {
				t.Errorf("parseInt(%s) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}