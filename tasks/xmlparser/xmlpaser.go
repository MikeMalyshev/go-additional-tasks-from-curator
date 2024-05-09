package xmlparser

import (
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
)

// Three types of Tags for tag structure
const (
	noTag = iota
	openTag
	closeTag
)

// Struct to store the tag information, whitch is returned by nextTag()
type tag struct {
	name     string
	content  string
	tagInfo  int
	firstIdx int
	lastIdx  int
}

// Struct to store the number of hits, while comparing tag sequences
type hits struct {
	count               int
	idx                 int
	tagSequence         []string
	unavailableSequence []string
}

// Testing function
func TestTextToXmlParser() {
	inputXml := "tasks/xmlparser/data/input.xml"
	inputTxt := "tasks/xmlparser/data/input.txt"

	outputXml := "tasks/xmlparser/data/output.xml"
	outputTxt := "tasks/xmlparser/data/output.txt"

	fmt.Printf("\nPrint the number of the test to begin (1 to perform xmlParser, 2 to perform textParser, 0 to exit):")

	var idx int
	_, err := fmt.Fscan(os.Stdin, &idx)
	if err != nil || idx < 0 {
		fmt.Println("\nError: Incorrect input")
	}

	switch idx {
	case 0:
		fmt.Println("Exit")
	case 1:
		output, err := parseInputXml(inputXml)

		if err != nil {
			fmt.Printf("An error have occured while parsing file %s: %s", inputXml, err)
			return
		}
		saveOutput(outputTxt, output)
		fmt.Printf("\n\tDone: Input have been taken from \"%s\", output have been be saved to \"%s\"\n", inputXml, outputTxt)
	case 2:
		output, err := parseInputTxt(inputTxt)

		if err != nil {
			fmt.Printf("An error have occured while parsing file %s: %s", inputXml, err)
			return
		}
		saveOutput(outputXml, output)
		fmt.Printf("\n\tDone: Input have been taken from \"%s\", output have been be saved to \"%s\"\n", inputTxt, outputXml)
	default:
		fmt.Printf("\nError: %d is not configured yet\n\n", idx)
	}
}

// Parser for text files. Elements presented as "E1.E2.E2 = Val" are converted to XML format
func parseInputTxt(filePath string) (string, error) {
	input, err := readFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	expressionList := strings.Split(input, string(Separator))

	var output string
	for _, expression := range expressionList {
		path, value, err := parseInputTextExpression(expression)
		if err != nil {
			err = fmt.Errorf("error parsing input expression: %v", err)
			return "", err
		}

		output, err = addElement(output, path, value)
		if err != nil {
			err = fmt.Errorf("error adding path %s to existing xml: %w", expression, err)
			return "", err
		}
	}
	return output, nil
}

// Parser for Xml files. Element from XML format are converted to text format presented as "E1.E2.E2 = Val"
func parseInputXml(filePath string) (string, error) {
	input, err := readFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}

	output, err := parseXml(input, 0, []string{}, tag{}, []string{})
	if err != nil {
		return "", fmt.Errorf("error parsing xml:  %w", err)
	}

	return strings.Join(output, "; "), nil
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("file reading error: %w", err)
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("error while reading file: %w", err)
		return "", err
	}

	return string(data), nil
}

func saveOutput(filepath, output string) {
	file, err := os.Create(filepath)

	if err != nil {
		fmt.Println("unable to create file:", err)
	}
	defer file.Close()

	file.WriteString(output)
}

// Parser for expressions in text file, splited by ";"
func parseInputTextExpression(expr string) ([]string, string, error) {
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

// Add a 'value' to 'data' string using XML format. Destination to a value must be given in 'tagSequence'
func addElement(data string, targetTagSequence []string, value string) (string, error) {
	maxHits, _, err := parseInsertedTags(data, 0, []string{}, targetTagSequence, hits{}, noTag)
	if err != nil {
		return "", err
	}

	additionalTags := createTags(maxHits.unavailableSequence, value, len(targetTagSequence)-len(maxHits.unavailableSequence))

	return data[:maxHits.idx] + additionalTags + data[maxHits.idx:], nil
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

// parseInsertedTags finds the longest existing way to 'searchingSequence'
// It also checks wheather transitional element is trying to be put into valued element
// or either if a valued element is trying to be put into transitional element
// The result is 'hits' structure, which contains estimated longest sequence and an unavailable sequence
// that should be created and the index of the last tag in existing sequence
func parseInsertedTags(data string, startIdx int, currentSequence, searchingSequence []string, maxHits hits, prevTagInfo int) (hits, int, error) {
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
	return parseInsertedTags(data, next.lastIdx+1, currentSequence, searchingSequence, maxHits, next.tagInfo)
}

// Iteratively parses xml document 'data' starting from index 'startIdx'. Current path (tagSequence) stores in 'elements'.
// 'previousTag' is used for finding a closed elements and to check for the posible errors in xml format
// 'finiteElements' contains the output and returns at the end of 'data'
func parseXml(data string, startIdx int, elements []string, previousTag tag, finiteElements []string) ([]string, error) {
	next, err := nextTag(data, startIdx)
	if err != nil {
		return []string{}, fmt.Errorf("parsing tags error: %w", err)
	}

	if next.tagInfo == noTag {
		return finiteElements, nil
	}

	if next.tagInfo == openTag {
		elements = append(elements, next.content)
	}

	if next.tagInfo == closeTag {
		closingTagName := elements[len(elements)-1]
		if next.name != closingTagName {
			return []string{}, fmt.Errorf("trying to close tag %s with a tag %s", closingTagName, next.name)
		}

		if previousTag.name == closingTagName {
			value := data[previousTag.lastIdx:next.firstIdx]
			finiteElements = append(finiteElements, strings.Join(elements, string(SequenceSeparator))+string(Equalize)+value)
		}
		elements = elements[:(len(elements))-1]
	}

	return parseXml(data, next.lastIdx+1, elements, next, finiteElements)
}

// Finds the next 'tag' in 'data', marks it as close or open and returns the index by which this tag ends
func nextTag(data string, startIdx int) (tag, error) {
	contentFirstIdx, contentLastIdx, tagFirstIdx, tagLastIdx := -1, -1, -1, -1
	info := noTag

	// Searching for '<'
	for i, s := range data[startIdx:] {
		if s == '<' {
			tagFirstIdx = startIdx + i
			contentFirstIdx = tagFirstIdx + 1
			info = openTag
			break
		}
	}

	// '<' have not been found
	if tagFirstIdx == -1 {
		return tag{}, nil
	}

	// Checking whether it is a close tag
	if data[contentFirstIdx] == '/' {
		info = closeTag
		contentFirstIdx++
	}

	// Searching for '>'
	for i, s := range data[contentFirstIdx:] {
		if s == '>' {
			contentLastIdx = contentFirstIdx + i
			tagLastIdx = contentLastIdx + 1
			break
		}
	}

	// '>' have not been found
	if tagLastIdx == -1 {
		return tag{}, fmt.Errorf("couldn't find end of a tag (idx:%d)", tagLastIdx)
	}

	content := data[contentFirstIdx:contentLastIdx]
	name := strings.Split(content, " ")[0]

	return tag{name, content, info, tagFirstIdx, tagLastIdx}, nil
}
