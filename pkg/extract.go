package pkg

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// # MiniYaml
//
//	```md
//	- Each line is a node.
//	- Nodes can be:
//	  - key-only.
//	  - key,value, and optionally a comment.
//	  - comment-only.
//	  - empty.
//	- Node indentation is either 1-tab-per-level or 4-spaces-per-level.
//	- Comments start with `#` and span the remainder of the line.
//	```
//
// Ref: https://www.openra.net/book/modding/miniyaml/index.html
var _ struct{}

const (
	DefaultKeyRegex string = `^[^:]+:[^.]+[.](Tooltip[.]Name|Buildable[.]Description|TooltipExtras(@[^.]+)?[.][^.]+)$`
)

func ExtractStringsFromFile(ctx context.Context, filenames []string, output, keyRegex string) (err error) {
	if keyRegex == "" {
		keyRegex = DefaultKeyRegex
	}
	regex, err := regexp.Compile(keyRegex)
	if err != nil {
		err = fmt.Errorf("failed to compile regexp: %w", err)
		return
	}

	outputFile := os.Stdout
	if output != "" && output != "-" {
		outputFile, err = os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o777)
		if err != nil {
			err = fmt.Errorf("failed to open output file: %w", err)
			return
		}
		defer outputFile.Close()
	}

	patchMap := map[string]string{}
	for _, filename := range filenames {
		var file *os.File
		file, err = os.Open(filename)
		if err != nil {
			err = fmt.Errorf("failed to open file: %w", err)
			return
		}
		defer file.Close()
		basename := filepath.Base(file.Name())

		var stack []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			indent, key, value := parseLine(scanner.Text())
			if len(key) == 0 {
				continue
			}

			for len(stack) < indent {
				stack = append(stack, "")
			}
			stack = append(stack[:indent], key)

			if len(value) == 0 {
				continue
			}

			keyPath := fmt.Sprintf("%s:%s", basename, strings.Join(stack, "."))
			if regex.MatchString(keyPath) {
				patchMap[keyPath] = value
			}
		}

		file.Close()
	}

	// encode patch map into file
	encoder := yaml.NewEncoder(outputFile)
	defer encoder.Close()
	err = encoder.Encode(patchMap)
	if err != nil {
		err = fmt.Errorf("failed to save patch map: %w", err)
		return
	}

	return
}

func PatchStringsInFile(ctx context.Context, filenames []string, patchFilename, outputDir string) (err error) {
	patchFile := os.Stdin
	// open patch file
	if patchFilename != "" && patchFilename != "-" {
		patchFile, err = os.Open(patchFilename)
		if err != nil {
			err = fmt.Errorf("failed to open patch file: %w", err)
			return
		}
		defer patchFile.Close()
	}

	patchMap := map[string]string{}
	// decode patch map data
	decoder := yaml.NewDecoder(patchFile)
	err = decoder.Decode(&patchMap)
	if err != nil {
		err = fmt.Errorf("failed to parse patch map: %w", err)
		return
	}

	// ensure output directory
	err = os.MkdirAll(outputDir, 0o777)
	if err != nil {
		err = fmt.Errorf("failed to create directory: %w", err)
		return
	}

	for _, filename := range filenames {
		var file, outputFile *os.File
		// open input file
		file, err = os.Open(filename)
		if err != nil {
			err = fmt.Errorf("failed to open file: %w", err)
			return
		}
		defer file.Close()
		basename := filepath.Base(file.Name())

		// open output file
		outputFile, err = os.OpenFile(filepath.Join(outputDir, basename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o777)
		if err != nil {
			err = fmt.Errorf("failed to open output file: %w", err)
			return
		}
		defer outputFile.Close()

		var stack []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			indent, key, value := parseLine(scanner.Text())
			if len(key) == 0 {
				fmt.Fprintln(outputFile)
				continue
			}

			for len(stack) < indent {
				stack = append(stack, "")
			}
			stack = append(stack[:indent], key)

			keyPath := fmt.Sprintf("%s:%s", basename, strings.Join(stack, "."))
			if patch := patchMap[keyPath]; patch != "" && value != patch {
				value = patch
			}
			fmt.Fprintf(outputFile, "%s%s: %s\n", strings.Repeat("\t", indent), key, value)
		}

		file.Close()
		outputFile.Close()
	}

	return
}

func parseLine(line string) (indent int, key, value string) {
	var tabs, spaces int
loop:
	for i, r := range line {
		switch r {
		case ' ':
			spaces++
		case '\t':
			tabs++
		default:
			// trim left spaces and trail comments
			line = strings.SplitN(line[i:], "#", 2)[0]
			break loop
		}
	}
	indent = tabs + spaces/4

	// split keys and values
	splitted := strings.SplitN(line, ":", 2)
	key = splitted[0]
	if len(splitted) > 1 {
		value = strings.TrimSpace(splitted[1])
	}

	return
}
