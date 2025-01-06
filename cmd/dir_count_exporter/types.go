package main

import (
	"os"
	"path/filepath"
)

type dirIdentifier struct {
	Glob string `prometheus_label:"glob"`
	Path string `prometheus_label:"path"`
}

type DirStats struct {
	GlobErrors map[string]int `prometheus_map:"glob_errors" prometheus_map_key:"glob"`
	GlobPathsIgnored map[string]int `prometheus_map:"glob_paths_ignored" prometheus_map_key:"glob"`
	ReadDirErrors map[dirIdentifier]int `prometheus_map:"read_dir_errors"`
	DirEntries map[dirIdentifier]int `prometheus_map:"dir_entries"`
}

func genGlobs(globs []string) (map[string][]string, map[string]error) {
	matches := make(map[string][]string)
	globErrors := make(map[string]error)

	for _, glob := range globs {
		ms, err := filepath.Glob(glob)
		if err != nil {
			globErrors[glob] = err
		} else {
			globErrors[glob] = nil
		}
				
		matches[glob] = ms
	}
	
	return matches, globErrors
}

func gatherStats(globMatches map[string][]string, xGlobErrors map[string]error) DirStats {
	globErrors := make(map[string]int)
	globPathsIgnored := make(map[string]int)
	readDirErrors := make(map[dirIdentifier]int)
	dirEntries := make(map[dirIdentifier]int)
	
	// glob errors first
	for glob, err := range xGlobErrors {
		if err != nil {
			globErrors[glob] = 1
		} else {
			globErrors[glob] = 0
		}
	}
	
	// Now iterate
	for glob, matches := range globMatches {
		globPathsIgnored[glob] = 0
		
		for _, match := range matches {
			id := dirIdentifier{
				Glob: glob,
				Path: match,
			}
			readDirErrors[id] = 0
			dirEntries[id] = 0			
			
			// Is it a directory?
			fi, err := os.Stat(match)
			if err != nil {
				readDirErrors[id]++
				continue
			}
			
			if !fi.IsDir() {
				globPathsIgnored[glob]++
				delete(readDirErrors, id)
				delete(dirEntries, id)
				continue
			}
			
			entries, err := os.ReadDir(match)
			if err != nil {
				readDirErrors[id]++
				continue
			}
			
			dirEntries[id] = len(entries)
		}
	}
	
	return DirStats{
		GlobErrors: globErrors,
		GlobPathsIgnored: globPathsIgnored,
		ReadDirErrors: readDirErrors,
		DirEntries: dirEntries,
	}
}