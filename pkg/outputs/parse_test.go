package outputs

import (
	"strings"
	"testing"

	"github.com/airplanedev/ojson"
	"github.com/stretchr/testify/require"
)

func TestParseOutput(tt *testing.T) {
	for _, test := range []struct {
		name          string
		log           string
		expectedName  string
		expectedValue interface{}
		opts          ParseOptions
	}{
		{
			name:          "default no colon",
			log:           "airplane_output hello",
			expectedName:  "output",
			expectedValue: "hello",
		},
		{
			name:          "default with colon",
			log:           "airplane_output: true",
			expectedName:  "output",
			expectedValue: true,
		},
		{
			name:          "named",
			log:           "airplane_output:named [1, 2, 3]",
			expectedName:  "named",
			expectedValue: []interface{}{float64(1), float64(2), float64(3)},
		},
		{
			name:          "quoted string",
			log:           "airplane_output \"hello world\"",
			expectedName:  "output",
			expectedValue: "hello world",
		},
		{
			name:          "named with extra spaces",
			log:           "airplane_output:my_output   hello world  ",
			expectedName:  "my_output",
			expectedValue: "hello world",
		},
		{
			name:          "named with tabs",
			log:           "airplane_output:tabs \thello\tworld",
			expectedName:  "tabs",
			expectedValue: "hello\tworld",
		},
		{
			name:          "empty value with colon and space",
			log:           "airplane_output: ",
			expectedName:  "output",
			expectedValue: "",
		},
		{
			name:          "empty value no colon and space",
			log:           "airplane_output ",
			expectedName:  "output",
			expectedValue: "",
		},
		{
			name:          "named and quoted",
			log:           "airplane_output:\"output\" hello",
			expectedName:  "output",
			expectedValue: "hello",
		},
		{
			name:          "named and quoted with spaces",
			log:           "airplane_output:\"output value\" hello",
			expectedName:  "output value",
			expectedValue: "hello",
		},
		{
			name:          "named and quoted with quoted value",
			log:           "airplane_output:\"output\" \"hello\"",
			expectedName:  "output",
			expectedValue: "hello",
		},
		{
			name:          "named and quoted with spaces with quoted value",
			log:           "airplane_output:\"output value\" \"hello\"",
			expectedName:  "output value",
			expectedValue: "hello",
		},
		{
			name:          "empty quoted name",
			log:           "airplane_output:\"\" hello",
			expectedName:  "output",
			expectedValue: "hello",
		},
		{
			name:          "named and single quoted",
			log:           "airplane_output:'output' hello",
			expectedName:  "output",
			expectedValue: "hello",
		},
		{
			name:          "named and single quoted with spaces",
			log:           "airplane_output:'output value' hello",
			expectedName:  "output value",
			expectedValue: "hello",
		},
		{
			name:          "named and single quoted with quoted value",
			log:           "airplane_output:'output' \"hello\"",
			expectedName:  "output",
			expectedValue: "hello",
		},
		{
			name:          "named and single quoted with spaces with quoted value",
			log:           "airplane_output:'output value' \"hello\"",
			expectedName:  "output value",
			expectedValue: "hello",
		},
		{
			name:          "empty single quoted name",
			log:           "airplane_output:'' hello",
			expectedName:  "output",
			expectedValue: "hello",
		},
		{
			name:          "malformed output",
			log:           "airplane_output:''' hello",
			expectedName:  "'''",
			expectedValue: "hello",
		},
	} {
		tt.Run(test.name, func(t *testing.T) {
			m := make(map[string]*strings.Builder)
			output, err := Parse(m, test.log, test.opts)
			require.NoError(t, err)
			require.Equal(t, test.expectedName, output.Name)
			require.Equal(t, "", output.Command)
			require.Equal(t, len(test.log), output.Size)
			expectedJSON := ojson.Value{V: test.expectedValue}
			require.Equal(t, expectedJSON, output.Value)
		})
	}
}

func TestParseOutputV2(tt *testing.T) {
	for _, test := range []struct {
		name             string
		log              string
		expectedJsonPath string
		expectedValue    interface{}
		opts             ParseOptions
	}{
		{
			name:             "default no colon",
			log:              `${COMMAND_PLACEHOLDER} "hello"`,
			expectedJsonPath: "",
			expectedValue:    "hello",
		},
		{
			name:             "default colon",
			log:              `${COMMAND_PLACEHOLDER}: true`,
			expectedJsonPath: "",
			expectedValue:    true,
		},
		{
			name:             "named",
			log:              "${COMMAND_PLACEHOLDER}:named [1, 2, 3]",
			expectedJsonPath: "named",
			expectedValue:    []interface{}{float64(1), float64(2), float64(3)},
		},
		{
			name:             "named with extra spaces",
			log:              `${COMMAND_PLACEHOLDER}:my_output   "hello world"  `,
			expectedJsonPath: "my_output",
			expectedValue:    "hello world",
		},
		{
			name:             "named with tabs",
			log:              "${COMMAND_PLACEHOLDER}:tabs \t\"hello world\"",
			expectedJsonPath: "tabs",
			expectedValue:    "hello world",
		},
		{
			name:             "complex path",
			log:              `${COMMAND_PLACEHOLDER}:a.b.c[5][4]["asdf"]["\"]"] "test"`,
			expectedJsonPath: `a.b.c[5][4]["asdf"]["\"]"]`,
			expectedValue:    "test",
		},
		{
			name:             "json test",
			log:              `${COMMAND_PLACEHOLDER}:["json[\""] {"b":[],"a":true,"\"] ":3}`,
			expectedJsonPath: `["json[\""]`,
			expectedValue: ojson.NewObject().
				SetAndReturn("b", []interface{}{}).
				SetAndReturn("a", true).
				SetAndReturn("\"] ", float64(3)),
		},
		{
			name:             "json test with output limit",
			log:              `${COMMAND_PLACEHOLDER}:["json[\""] {"b":[],"a":true,"\"] ":3}`,
			expectedJsonPath: `["json[\""]`,
			expectedValue: ojson.NewObject().
				SetAndReturn("b", []interface{}{}).
				SetAndReturn("a", true).
				SetAndReturn("\"] ", float64(3)),
			opts: ParseOptions{
				OutputLineMaxBytes: 100,
			},
		},
	} {
		tt.Run(test.name+" set", func(t *testing.T) {
			logText := strings.Replace(test.log, "${COMMAND_PLACEHOLDER}", "airplane_output_set", 1)
			m := make(map[string]*strings.Builder)
			output, err := Parse(m, logText, test.opts)
			require.NoError(t, err)
			require.Equal(t, test.expectedJsonPath, output.JsonPath)
			require.Equal(t, "set", output.Command)
			expectedJSON := ojson.Value{V: test.expectedValue}
			require.Equal(t, expectedJSON, output.Value)
		})

		tt.Run(test.name+" append", func(t *testing.T) {
			logText := strings.Replace(test.log, "${COMMAND_PLACEHOLDER}", "airplane_output_append", 1)
			m := make(map[string]*strings.Builder)
			output, err := Parse(m, logText, test.opts)
			require.NoError(t, err)
			require.Equal(t, test.expectedJsonPath, output.JsonPath)
			require.Equal(t, "append", output.Command)
			expectedJSON := ojson.Value{V: test.expectedValue}
			require.Equal(t, expectedJSON, output.Value)
		})
	}
}

func TestParseOutputChunks(tt *testing.T) {
	for _, test := range []struct {
		name    string
		logs    []string
		outputs []ParsedLine
		opts    ParseOptions
	}{
		{
			name: "simple chunk",
			logs: []string{
				`airplane_chunk:asdf airplane_output:asdf `,
				`airplane_chunk:asdf hello world`,
				`airplane_chunk_end:asdf`,
			},
			outputs: []ParsedLine{
				{
					Name:  "asdf",
					Value: ojson.Value{V: "hello world"},
				},
			},
		},
		{
			name: "interlaced chunks",
			logs: []string{
				`airplane_chunk:asdf airplane_output:asdf `,
				`airplane_chunk:ghjkl airplane_output_set:ghjkl `,
				`airplane_chunk:asdf hello`,
				`airplane_chunk:ghjkl ["ghjkl"]`,
				`airplane_chunk:asdf  world`,
				`airplane_chunk_end:asdf`,
				`airplane_chunk_end:ghjkl`,
			},
			outputs: []ParsedLine{
				{
					Name:  "asdf",
					Value: ojson.Value{V: "hello world"},
				},
				{
					Command:  "set",
					JsonPath: "ghjkl",
					Value:    ojson.Value{V: []interface{}{"ghjkl"}},
				},
			},
			opts: ParseOptions{
				OutputLineMaxBytes: 100,
			},
		},
	} {
		tt.Run(test.name, func(t *testing.T) {
			m := make(map[string]*strings.Builder)
			var outputs []ParsedLine
			for _, logText := range test.logs {
				output, err := Parse(m, logText, test.opts)
				require.NoError(t, err)
				if output != nil {
					outputs = append(outputs, *output)
				}
			}
			require.Equal(t, len(test.outputs), len(outputs))
			for i := 0; i < len(outputs); i++ {
				require.Equal(t, test.outputs[i].Name, outputs[i].Name)
				require.Equal(t, test.outputs[i].Command, outputs[i].Command)
				require.Equal(t, test.outputs[i].JsonPath, outputs[i].JsonPath)
				require.Equal(t, test.outputs[i].Value, outputs[i].Value)
			}
		})
	}
}

func TestLongOutput(tt *testing.T) {
	for _, test := range []struct {
		name string
		logs []string
	}{
		{
			name: "output chunks",
			logs: []string{
				`airplane_chunk:asdf airplane_output:asdf `,
				`airplane_chunk:asdf hello world`,
				`airplane_chunk_end:asdf`,
			},
		},
		{
			name: "output set",
			logs: []string{
				`airplane_output_set "abcdefghijklmnopqrstuvwxyz"`,
			},
		},
		{
			name: "output append",
			logs: []string{
				`airplane_output_append "abcdefghijklmnopqrstuvwxyz"`,
			},
		},
	} {
		tt.Run(test.name, func(t *testing.T) {
			require := require.New(t)
			m := make(map[string]*strings.Builder)
			for i, logText := range test.logs {
				output, err := Parse(m, logText, ParseOptions{
					OutputLineMaxBytes: 10,
				})
				require.Nil(output)
				if i == len(test.logs)-1 {
					require.Equal(err, ErrOutputLineTooLong)
				} else {
					require.NoError(err)
				}
			}
		})
	}
}
