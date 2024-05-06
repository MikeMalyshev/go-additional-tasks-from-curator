package xmlparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	Separator         = ';'
	SequenceSeparator = '.'
	Equalize          = '='

	XMLIndent = "   "

	// Three types of Tags
	noTag    = 0
	openTag  = 1
	closeTag = 2
)

type tag struct {
	name     string
	tagInfo  int
	firstIdx int
	lastIdx  int
}

type hits struct {
	count               int
	idx                 int
	tagSequence         []string
	unavailableSequence []string
}

func TestTextToXmlParser() {
	input := "tasks/xmlparser/data/input.txt"
	output := "tasks/xmlparser/data/output.xml"
	fmt.Println("\n\tThis utility converts string comands like 'a.b.c=345' to xml-document.")
	fmt.Printf("\n\tInput expressions will be taken from \"%s\", output will be saved to \"%s\"\n", input, output)
	fmt.Printf("\nPress any key to continue ...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	saveOutput(output, parseInputFile(input))
}

func parseInputFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("\nfile reading error: %s\n", err)
		return ""
	}

	defer func() {
		if err = file.Close(); err != nil {
			fmt.Printf("\nerror while closing file: %s\n", err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("\nerror while reading file: %s\n", err)
		return ""
	}

	strData := string(data)
	expressionList := strings.Split(strData, string(Separator))

	var output string
	for _, expression := range expressionList {
		path, value, err := parseInputExpression(expression)
		if err != nil {
			fmt.Printf("error parsing input expression: %v", err)
			continue
		}

		output = addElement(output, path, value)
	}
	return output
}

func saveOutput(filepath, output string) {
	file, err := os.Create(filepath)

	if err != nil {
		fmt.Println("unable to create file:", err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			fmt.Printf("\nerror while closing file: %s\n", err)
		}
	}()

	file.WriteString(output)

	fmt.Println("Done.")
}

func parseInputExpression(expr string) ([]string, string, error) {
	expr = strings.TrimSpace(expr)
	pair := strings.Split(expr, string(Equalize))

	if len(pair) != 2 {
		err := fmt.Errorf("incorrect data")
		return []string{}, "", err
	}

	valuePath := strings.Split(pair[0], string(SequenceSeparator))
	value := pair[1]

	return valuePath, value, nil
}

// Add elements presented int 'path' as E1.E2.E2 = 'value' to 'data' string in XML format
func addElement(data string, targetTagSequence []string, value string) string {
	unavailableTagSequence, idx := findBestPath(data, targetTagSequence)

	additionalTags := createTags(unavailableTagSequence, value, len(targetTagSequence)-len(unavailableTagSequence))

	return data[:idx] + additionalTags + data[idx:]
}

// func createTags creates all necessary tags,  specified by 'tagSequence'. The last tag wil have a nested 'value'.
// The indent will be setted in according the 'level' value and const 'XMLIndent'.
func createTags(tagSequence []string, value string, level int) string {
	if len(tagSequence) < 1 {
		return ""
	}

	tagName := tagSequence[0]
	indent := strings.Repeat(XMLIndent, level)

	if len(tagSequence) == 1 {
		return indent + "<" + tagName + ">" + value + "</" + tagName + ">\n"
	}

	open := indent + "<" + tagName + ">" + "\n"
	close := indent + "</" + tagName + ">\n"

	return open + createTags(tagSequence[1:], value, level+1) + close
}

// Find longest existing sequence of tags in 'data', returns the unavailable part of sequence
// and index where it should be placed acording to target 'searchingSequence'
func findBestPath(data string, searchingSequence []string) ([]string, int) {
	var tagSequence []string

	maxHits, _, err := parseTags(data, 0, tagSequence, searchingSequence, hits{}, noTag)
	if err != nil {
		fmt.Println(err)
	}
	return maxHits.unavailableSequence, maxHits.idx
}

// Compares two sequences and return number of hits in a row
func compareSequences(seq1, seq2 []string) int {
	n := len(seq1)
	if len(seq1) > len(seq2) {
		n = len(seq2)
	}
	hits := 0
	for i := 0; i < n; i++ {
		if seq1[i] == seq2[i] {
			hits++
		} else {
			break
		}
	}
	return hits
}

// TagParser finds the longest existing way to 'searchingSequence'
// It also checks wheather transitional element is trying to be put into valued element
// or either if a valued element is trying to be put into transitional element
// The result is 'hits' structure, which contains estimated longest sequence and an unavailable sequence
// that should be created and the index of the last tag in existing sequence
func parseTags(data string, startIdx int, currentSequence, searchingSequence []string, maxHits hits, prevTagInfo int) (hits, int, error) {
	if len(strings.TrimSpace(data)) == 0 {
		return hits{0, 0, searchingSequence, searchingSequence}, noTag, nil
	}

	next, err := nextTag(data, startIdx)
	if err != nil {
		return hits{}, next.tagInfo, fmt.Errorf("parsing tags error: %w", err)
	}

	if next.tagInfo == noTag {
		return maxHits, next.tagInfo, nil
	}

	if next.tagInfo == openTag {
		currentSequence = append(currentSequence, next.name)

		// Checks whether valued element is trying to be put into transitional element
		hitsCount := compareSequences(currentSequence, searchingSequence)
		if hitsCount == len(searchingSequence) {
			return hits{}, next.tagInfo, fmt.Errorf("%s already exists", strings.Join(currentSequence, "."))
		}
	}

	if next.tagInfo == closeTag {
		lastIdx := len(currentSequence) - 1
		previousTagName := currentSequence[lastIdx]

		// Checks whether transitional element is trying to be put into valued element
		hitsCount := compareSequences(currentSequence, searchingSequence)
		if prevTagInfo == openTag && hitsCount == len(currentSequence) {
			return hits{}, next.tagInfo, fmt.Errorf("%s already exists", strings.Join(currentSequence, "."))
		}

		if next.name != previousTagName {
			return hits{}, next.tagInfo, fmt.Errorf("trying to close tag %s with a tag %s", previousTagName, next.name)
		}

		// Searching the best tag sequence
		hitsCount = compareSequences(currentSequence[:lastIdx], searchingSequence)
		if hitsCount >= maxHits.count {
			maxHits.count = hitsCount
			maxHits.idx = next.lastIdx + 1 // Next insertion should be after the last symbol (lastIdx + 1)
			maxHits.tagSequence = currentSequence[:lastIdx]

			if len(maxHits.tagSequence) < len(searchingSequence) {
				maxHits.unavailableSequence = searchingSequence[len(maxHits.tagSequence):]
			}
		}
		currentSequence = currentSequence[:lastIdx]

		// If matched tag sequence becomes shorted, than the result have already been found
		hitsCount = compareSequences(currentSequence[:lastIdx], searchingSequence)
		if hitsCount < maxHits.count {
			return maxHits, next.tagInfo, err
		}
	}
	return parseTags(data, next.lastIdx+1, currentSequence, searchingSequence, maxHits, next.tagInfo)
}

// Finds the next 'tag' in 'data', marks it as close or open and returns the index by which this tag ends
func nextTag(data string, startIdx int) (tag, error) {
	nameFirstIdx, nameLastIdx, tagFirstIdx, tagLastIdx := -1, -1, -1, -1
	info := noTag

	// Searching for '<'
	for i, s := range data[startIdx:] {
		if s == '<' {
			tagFirstIdx = startIdx + i
			nameFirstIdx = tagFirstIdx + 1
			info = openTag
			break
		}
	}

	// '<' have not been found
	if tagFirstIdx == -1 {
		return tag{}, nil
	}

	// Checking whether it is a close tag
	if data[nameFirstIdx] == '/' {
		info = closeTag
		nameFirstIdx++
	}

	// Searching for '>'
	for i, s := range data[nameFirstIdx:] {
		if s == '>' {
			nameLastIdx = nameFirstIdx + i
			tagLastIdx = nameLastIdx + 1
			break
		}
	}

	// '>' have not been found
	if tagLastIdx == -1 {
		return tag{}, fmt.Errorf("couldn't find end of a tag (idx:%d)", tagLastIdx)
	}

	name := strings.Split(data[nameFirstIdx:nameLastIdx], " ")[0]

	return tag{name, info, tagFirstIdx, tagLastIdx}, nil
}
