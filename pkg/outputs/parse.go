package outputs

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/airplanedev/ojson"
	"github.com/airplanedev/path"
	"github.com/pkg/errors"
)

const outputPrefix = "airplane_output"
const defaultOutputName = "output"

var outputRegexp = regexp.MustCompile((`^airplane_output(?::(?:("[^"]*")|('[^']*')|([^ ]+))?)? (.*)$`))
var jsonPathRegexp = regexp.MustCompile((`^airplane_output(_set|_append)(:| )(.*)$`))

type ParsedLine struct {
	Command  string
	Name     string
	JsonPath string
	Value    ojson.Value
	Size     int
}

var ErrOutputLineTooLong = errors.New("output line too long")

func parseOutputValue(matches []string) ojson.Value {
	vs := strings.TrimSpace(matches[4])

	var target ojson.Value
	if err := json.Unmarshal([]byte(vs), &target); err != nil {
		// Interpret this output as a string
		target = ojson.Value{V: vs}
	}
	return target
}

func parseOutputName(matches []string) string {
	var outputName string
	if matches[1] != "" {
		outputName = strings.Trim(matches[1], "\"")
	} else if matches[2] != "" {
		outputName = strings.Trim(matches[2], "'")
	} else if matches[3] != "" {
		outputName = matches[3]
	}
	outputName = strings.TrimSpace(outputName)
	if outputName != "" {
		return outputName
	}
	return defaultOutputName
}

func parseOutputLegacy(s string) *ParsedLine {
	if matches := outputRegexp.FindStringSubmatch(s); matches != nil {
		name := parseOutputName(matches)
		value := parseOutputValue(matches)

		return &ParsedLine{
			Command: "",
			Name:    name,
			Value:   value,
		}
	}
	return nil
}

func parseOutputV2(s string) (*ParsedLine, error) {
	if matches := jsonPathRegexp.FindStringSubmatch(s); matches != nil {
		command := matches[1][1:]
		var jsonPath string
		var valueStr string
		if matches[2] == ":" {
			_, idx, err := path.FromJSPartial(matches[3])
			if err != nil {
				return nil, err
			}
			if len(matches[3]) <= idx || matches[3][idx] != ' ' {
				return nil, errors.New("invalid output line")
			}
			jsonPath = matches[3][:idx]
			valueStr = strings.TrimSpace(matches[3][idx+1:])
		} else {
			valueStr = matches[3]
		}

		var value ojson.Value
		if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
			return nil, err
		}

		return &ParsedLine{
			Command:  command,
			JsonPath: jsonPath,
			Value:    value,
		}, nil
	}
	return nil, nil
}

const outputChunkPrefix = "airplane_chunk"

var outputChunkRegexp = regexp.MustCompile((`^airplane_chunk(|_end):([^ ]*)(?: (.+)|)$`))

// parseLogLineWithChunks converts a log line's text into the "effective" log
// text, by handling merging individual chunks into full lines. The chunking
// protocol works as follows:
//
// Chunks have a "chunk key". Lines starting with `airplane_chunk:<chunk_key>`
// will accumulate (returning empty strings in the meantime), and the
// concatenation of all the chunks will be returned when a
// `airplane_chunk_end:<chunk_key>` line is encountered. Non-chunk-related
// lines will just return their original value as normal.
//
// For example, feeding the following lines (in backticks) into
// parseLogLineWithChunks in succession will produce:
// `abcdef` => "abcdef"
// `airplane_chunk:key1 abc` => ""
// `airplane_chunk:key2 123` => ""
// `airplane_chunk:key1 def` => ""
// `airplane_chunk:key2 456` => ""
// `airplane_chunk_end:key1` => "abcdef"
// `airplane_chunk_end:key2` => "123456"
func parseLogLineWithChunks(text string, chunks map[string]*strings.Builder) (string, error) {
	if strings.HasPrefix(text, outputChunkPrefix) {
		if matches := outputChunkRegexp.FindStringSubmatch(text); matches != nil {
			chunkKey := matches[2]
			chunk, ok := chunks[chunkKey]
			if !ok {
				chunks[chunkKey] = &strings.Builder{}
				chunk = chunks[chunkKey]
			}
			if matches[1] == "_end" {
				delete(chunks, chunkKey)
				return chunk.String(), nil
			} else {
				if len(matches[3]) > 0 {
					chunk.WriteString(matches[3])
				}
				// If this is not an ending chunk, then we treat the resulting text as
				// an empty line.
				return "", nil
			}
		} else {
			return "", errors.Errorf("line started with airplane_chunk but was not a valid chunk: %s", text)
		}
	}
	return text, nil
}

type ParseOptions struct {
	// OutputLineMaxBytes is the maximum number of bytes a single line of output (across
	// all output chunks) can contain. Any lines larger than this will return an error.
	//
	// Disabled if <= 0 (i.e. no limit is applied on the size of each output line).
	OutputLineMaxBytes int
}

func Parse(chunks map[string]*strings.Builder, logText string, opts ParseOptions) (*ParsedLine, error) {
	outputText, err := parseLogLineWithChunks(logText, chunks)
	if err != nil {
		return nil, err
	}

	if opts.OutputLineMaxBytes > 0 && opts.OutputLineMaxBytes < len(outputText) {
		return nil, ErrOutputLineTooLong
	}

	// Extract subset of logs as output
	if !strings.HasPrefix(outputText, outputPrefix) {
		return nil, nil
	}

	var line *ParsedLine
	line = parseOutputLegacy(outputText)
	if line == nil {
		line, err = parseOutputV2(outputText)
	}
	if line == nil && err != nil {
		err = errors.New("line does not match any known airplane_output format")
	}

	if line == nil {
		// This catch-all clause is for backwards compatibility with the
		// original airplane_output behavior.
		line = &ParsedLine{
			Command: "",
			Name:    defaultOutputName,
			Value:   ojson.Value{V: ""},
		}
	} else {
		line.Size = len(outputText)
	}

	return line, err
}
