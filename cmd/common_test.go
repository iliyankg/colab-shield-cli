package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
)

func TestBuildContext(t *testing.T) {
	// Setup & Action
	ctx, cancel := buildContext("projectId", "userId")

	// Assert
	if ctx == nil {
		t.Error("Expected context to not be nil")
	}

	if cancel == nil {
		t.Error("Expected cancel to not be nil")
	}

	metadata, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Error("Expected metadata to be present in context")
	}

	if len(metadata) != 2 {
		t.Errorf("Expected metadata to have length 2, got %v", len(metadata))
	}

	if metadata.Get("projectId")[0] != "projectId" {
		t.Errorf("Expected projectId to be projectId, got %v", metadata["projectId"][0])
	}

	if metadata.Get("userId")[0] != "userId" {
		t.Errorf("Expected userId to be userId, got %v", metadata["userId"][0])
	}
}

func TestStringToHashMap(t *testing.T) {
	// Setup
	strings := []string{"a", "b", "c", "d"}

	// Action
	hashMap := strSliceToHashMap(strings)

	// Assert
	if len(hashMap) != 4 {
		t.Errorf("Expected hashmap to have length 4, got %v", len(hashMap))
	}

	for _, str := range strings {
		if _, ok := hashMap[str]; !ok {
			t.Errorf("Expected hashmap to contain %v", str)
		}
	}
}

func TestMatchPathToExcludedPaths(t *testing.T) {
	// TODO: Consider more thorough testing but at some point we'd just be 
	// testing the path match.

	// Setup
	excludedPaths := []string{
		"file1.go",
		"file2.go",
		"file3.go",
		"file4.go",
	}

	t.Run("Match Excluded Path", func(t *testing.T) {
		filePath := "file1.go"

		// Action
		match, err := matchPathToExcludedPaths(filePath, excludedPaths)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Assert
		if !match {
			t.Errorf("Expected match to be true, got %v", match)
		}
	})

	t.Run("Match Unexcluded Path", func(t *testing.T) {
		filePath := "file5.go"

		// Action
		match, err := matchPathToExcludedPaths(filePath, excludedPaths)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Assert
		if match {
			t.Errorf("Expected match to be false, got %v", match)
		}
	})
}

func TestFilterToFilesOfInterest(t *testing.T) {
	// Setup
	files := []string{
		"file1.go",
		"file2.go",
		"file3.go",
		"file4.go",
		"file5.go",
		"file6.go",
	}

	viper.Set("extensions", []string{".go"})
	viper.Set("ignore", []string{"file1.go", "file2.go"})

	// Action
	filteredFiles, err := filterToFilesOfInterest(files)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Assert
	if len(filteredFiles) != 4 {
		t.Errorf("Expected filtered files to have length 4, got %v", len(filteredFiles))
	}

	for i, file := range filteredFiles {
		if file != files[i+2] {
			t.Errorf("Expected file to be %v, got %v", files[i+2], file)
		}
	}
}
