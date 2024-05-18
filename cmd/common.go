package cmd

import (
	"context"
	"errors"
	"path"
	"time"

	"github.com/iliyankg/colab-shield/cli/config"
	"google.golang.org/grpc/metadata"
)

var (
	ErrFileToHashMissmatch = errors.New("files and hashes must be of the same length")
)

func buildContext(projectId string, userId string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	metaInfo := metadata.Pairs(
		"projectId", projectId,
		"userId", userId,
	)

	return metadata.NewOutgoingContext(ctx, metaInfo), cancel
}

// filterToFilesOfInterest filters the given files to only those that are of interest basedon
// the extensions and ignore paths in the config.
func filterToFilesOfInterest(files []string) ([]string, error) {
	mappedExtensions := strSliceToHashMap(config.Extensions())
	excludedPaths := config.IgnorePaths()

	toReturn := make([]string, 0)
	for _, file := range files {
		extension := path.Ext(file)
		if _, ok := mappedExtensions[extension]; !ok {
			continue
		}

		if match, err := matchPathToExcludedPaths(file, excludedPaths); err != nil {
			return nil, err
		} else if match {
			continue
		}

		toReturn = append(toReturn, file)
	}

	return toReturn, nil
}

func matchPathToExcludedPaths(filePath string, excludedPaths []string) (bool, error) {
	for _, excludedPath := range excludedPaths {
		if match, err := path.Match(excludedPath, filePath); err != nil {
			return false, err
		} else if match {
			return true, nil
		}
	}
	return false, nil
}

func strSliceToHashMap(slice []string) map[string]any {
	hashMap := make(map[string]any)
	for _, s := range slice {
		hashMap[s] = nil
	}
	return hashMap
}
