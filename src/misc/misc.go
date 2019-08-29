// misc contains some helper functions
package misc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ErrorCheck is a function to throw error to the log and exit the program
func ErrorCheck(msg error) {
	if msg != nil {
		log.Fatalf("terminated\n\nERROR --> %v\n\n", msg)
	}
}

// CheckDir is a function to check that a directory exists
func CheckDir(dir string) error {
	if dir == "" {
		return fmt.Errorf("no directory specified")
	}
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %v", dir)
		}
		return fmt.Errorf("can't access adirectory (check permissions): %v", dir)
	}
	return nil
}

// CheckFile is a function to check that a file can be read
func CheckFile(file string) error {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %v", file)
		}
		return fmt.Errorf("can't access file (check permissions): %v", file)
	}
	return nil
}

// CollectFiles is a function to find all files with the specified extension in a specified directory
func CollectFiles(inputDir string, requiredExtension string, recursive bool) ([]string, error) {
	if inputDir == "" {
		return nil, fmt.Errorf("no inputDir was specified for file collection")
	}
	if requiredExtension == "" {
		return nil, fmt.Errorf("no extension was specified for file collection")
	}
	filePaths := []string{}
	extGlob := "*." + requiredExtension

	// if recursive, find all EXT files in the input directory and its subdirectories
	if recursive == true {
		recursiveSketchGrabber := func(fp string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fi.IsDir() {
				return nil
			}
			matched, err := filepath.Match(extGlob, fi.Name())
			if err != nil {
				return err
			}
			// if a file is found, add the file path the list
			if matched {
				filePaths = append(filePaths, fp)
			}
			return nil
		}
		filepath.Walk(inputDir, recursiveSketchGrabber)
	} else {

		// otherwise, just find all EXT files in the supplied directory
		searchTerm := inputDir + "/" + extGlob
		var err error
		filePaths, err = filepath.Glob(searchTerm)
		if err != nil {
			return nil, err
		}
	}

	// check we got some files
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("no %v files found in supplied directory: %v", requiredExtension, inputDir)
	}

	// check the files are accessible
	for _, f := range filePaths {
		if err := CheckFile(f); err != nil {
			return nil, fmt.Errorf("could not access %v file: %v", requiredExtension, f)
		}
	}

	return filePaths, nil
}
