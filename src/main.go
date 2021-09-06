/*
Copyright 2021 Jung Bong-Hwa

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	yamlToMd := flag.Bool("ym", false, "YAML to Markdown")
	yamlToMdList := flag.Bool("yml", false, "YAML to Markdown list")
	mdToConfluence := flag.Bool("mc", false, "Markdown to Confluence")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Convert strings of clipboard.\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	clip, _ := clipboard.ReadAll()
	output := ""
	if *yamlToMd {
		output = convertYamlToMd(clip)
	} else if *yamlToMdList {
		output = convertYamlToMdList(clip)
	} else if *mdToConfluence {
		output = convertMarkdownToConfluence(clip)
	} else {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Println(output)
}

func convertYamlToMd(input string) string {
	var output bytes.Buffer

	block := false
	newBlock := false
	lastSpaces := ""
	blockSpaces := ""
	lineCnt := 0

	reTitleOnly := regexp.MustCompile(`^([^\s])(.*):\s*$`)
	reTitleContent := regexp.MustCompile(`^([^\s])(.*): (.+)`)
	reSubTitle := regexp.MustCompile(`^\s+(.+):\s*$`)
	reTitleBlock := regexp.MustCompile(`^([^\s])(.*): \|\s*$`)
	reBlock := regexp.MustCompile(`^\s*(.+): \|\s*$`)
	reSpace := regexp.MustCompile(`^[\s]+`)
	reHttp := regexp.MustCompile(`(http[^\s,\,]+)`)
	var reContent *regexp.Regexp = nil

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		spaces := reSpace.FindString(line)
		var s string

		if block {
			if newBlock {
				reContent = regexp.MustCompile(`^` + spaces)
				blockSpaces = spaces
				newBlock = false
			} else if len(line) > 0 &&
				strings.Compare(blockSpaces, spaces) > 0 {
				// Ends a block.
				output.WriteString("```\n")
				block = false
			}
		}

		if block {
			s = reContent.ReplaceAllString(line, "")
		} else {
			if reTitleOnly.MatchString(line) {
				s = reTitleOnly.ReplaceAllString(line, "\n## $1$2")
			} else if reTitleBlock.MatchString(line) {
				s = reBlock.ReplaceAllString(line, "\n## $1$2\n\n```txt")
				block = true
				newBlock = true
			} else if reBlock.MatchString(line) {
				s = reBlock.ReplaceAllString(line, "\n$1:\n\n```txt")
				block = true
				newBlock = true
			} else if reTitleContent.MatchString(line) {
				s = reTitleContent.ReplaceAllString(line, "\n## $1$2\n\n$3")
			} else if reSubTitle.MatchString(line) {
				s = reSubTitle.ReplaceAllString(line, "\n$1:")
			} else {
				s = strings.Trim(line, " ")
				s = reHttp.ReplaceAllString(s, "<$1>")
				if strings.Index(s, "- ") == 0 {
					s = "* " + s[2:]
				} else {
					s = "* " + s
				}
				if strings.Compare(lastSpaces, spaces) != 0 {
					s = "\n" + s
				}
			}

		}
		lastSpaces = spaces
		output.WriteString(s + "\n")
	}

	if block {
		output.WriteString("```\n")
	}
	return output.String()
}

func convertYamlToMdList(input string) string {
	var output bytes.Buffer

	reSpace := regexp.MustCompile(`^[\s]+`)
	reContent := regexp.MustCompile(`^[\s,-]*(.+)`)
	reHttp := regexp.MustCompile(`(http[^\s,\,]+)`)

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		spaces := reSpace.FindString(line)
		var s string

		s = strings.Trim(line, " ")
		s = reHttp.ReplaceAllString(s, "<$1>")
		s = reContent.ReplaceAllString(s, spaces+"* $1")
		output.WriteString(s + "\n")
	}

	return output.String()
}

func convertMarkdownToConfluence(input string) string {
	var output bytes.Buffer

	reHttp := regexp.MustCompile(`(.*)<(http[^\s,\,]+)>`)

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		var s string

		if reHttp.MatchString(line) {
			s = reHttp.ReplaceAllString(line, "$1$2")
		} else {
			s = line
		}
		output.WriteString(s + "\n")
	}

	return output.String()
}
