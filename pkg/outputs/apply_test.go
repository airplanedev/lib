package outputs

import (
	"testing"

	"github.com/airplanedev/ojson"
	"github.com/stretchr/testify/require"
)

func TestApplyLegacy(tt *testing.T) {
	for _, test := range []struct {
		testName string
		name     string
		v        string
		initial  string
		final    string
	}{
		{
			"single string output",
			"output",
			`"testing"`,
			`{}`,
			`{"output":["testing"]}`,
		},
		{
			"single string null",
			"output",
			`"testing"`,
			`null`,
			`{"output":["testing"]}`,
		},
		{
			"append output",
			"output",
			`"testing"`,
			`{"output":["asdf"]}`,
			`{"output":["asdf","testing"]}`,
		},
		{
			"append output 2",
			"a.b.c",
			`"testing"`,
			`{"output":["asdf"]}`,
			`{"output":["asdf"],"a.b.c":["testing"]}`,
		},
		{
			"append object",
			"output",
			`{"a":"b"}`,
			`{"output":["asdf"]}`,
			`{"output":["asdf",{"a":"b"}]}`,
		},
	} {
		tt.Run(test.testName, func(t *testing.T) {
			require := require.New(t)
			cmd := &ParsedLine{
				Name:    test.name,
				Command: "",
				Value:   ojson.MustNewValueFromJSON(test.v),
			}
			v := ojson.MustNewValueFromJSON(test.initial)
			require.NoError(ApplyOutputCommand(cmd, &v))
			require.Equal(v, ojson.MustNewValueFromJSON(test.final))
		})
	}
}

func TestApplySet(tt *testing.T) {
	for _, test := range []struct {
		testName string
		jsonPath string
		v        string
		initial  string
		final    string
	}{
		{
			"root null set",
			"",
			`"testing"`,
			`null`,
			`"testing"`,
		},
		{
			"root object set",
			"",
			`{"c":"d"}`,
			`{"a":"b"}`,
			`{"c":"d"}`,
		},
		{
			"basic object set without existing key",
			"output2",
			`"testing"`,
			`{"output":["asdf"]}`,
			`{"output":["asdf"],"output2":"testing"}`,
		},
		{
			"basic object set with existing key",
			"output",
			`"testing"`,
			`{"output":["asdf"]}`,
			`{"output":"testing"}`,
		},
		{
			"empty key set",
			`[""]`,
			`"testing"`,
			`{}`,
			`{"":"testing"}`,
		},
		{
			"basic array set",
			"a[1]",
			`"testing"`,
			`{"a":[null,5,3]}`,
			`{"a":[null,"testing",3]}`,
		},
		{
			"nested object set",
			`a[0].b.c["\""]`,
			`"testing"`,
			`{"a":[{"b":{"c":{}}}]}`,
			`{"a":[{"b":{"c":{"\"":"testing"}}}]}`,
		},
		{
			"create nested objects",
			`a.b`,
			`"c"`,
			`null`,
			`{"a":{"b":"c"}}`,
		},
	} {
		tt.Run(test.testName, func(t *testing.T) {
			require := require.New(t)
			cmd := &ParsedLine{
				Command:  "set",
				JsonPath: test.jsonPath,
				Value:    ojson.MustNewValueFromJSON(test.v),
			}
			v := ojson.MustNewValueFromJSON(test.initial)
			require.NoError(ApplyOutputCommand(cmd, &v))
			require.Equal(v, ojson.MustNewValueFromJSON(test.final))
		})
	}
}

func TestApplyAppend(tt *testing.T) {
	for _, test := range []struct {
		testName string
		jsonPath string
		v        string
		initial  string
		final    string
	}{
		{
			"root null append",
			"",
			`"testing"`,
			`null`,
			`["testing"]`,
		},
		{
			"root array append",
			"",
			`{"c":"d"}`,
			`[true,1,null]`,
			`[true,1,null,{"c":"d"}]`,
		},
		{
			"basic object append without existing key",
			"output2",
			`"testing"`,
			`{"output":["asdf"]}`,
			`{"output":["asdf"],"output2":["testing"]}`,
		},
		{
			"basic object append with existing key",
			"output",
			`"testing"`,
			`{"output":["asdf"]}`,
			`{"output":["asdf","testing"]}`,
		},
		{
			"empty key append",
			`[""]`,
			`"testing"`,
			`{}`,
			`{"":["testing"]}`,
		},
		{
			"basic array append",
			"a[1]",
			`"testing"`,
			`{"a":[null,[],3]}`,
			`{"a":[null,["testing"],3]}`,
		},
		{
			"nested object append",
			`a[0].b.c["\""]`,
			`"testing"`,
			`{"a":[{"b":{"c":{}}}]}`,
			`{"a":[{"b":{"c":{"\"":["testing"]}}}]}`,
		},
		{
			"create nested objects",
			`a.b`,
			`"c"`,
			`null`,
			`{"a":{"b":["c"]}}`,
		},
	} {
		tt.Run(test.testName, func(t *testing.T) {
			require := require.New(t)
			cmd := &ParsedLine{
				Command:  "append",
				JsonPath: test.jsonPath,
				Value:    ojson.MustNewValueFromJSON(test.v),
			}
			v := ojson.MustNewValueFromJSON(test.initial)
			require.NoError(ApplyOutputCommand(cmd, &v))
			require.Equal(v, ojson.MustNewValueFromJSON(test.final))
		})
	}
}

func TestSetInvalid(tt *testing.T) {
	for _, test := range []struct {
		testName string
		jsonPath string
		v        string
		initial  string
	}{
		{
			"invalid path",
			"[",
			`"testing"`,
			"null",
		},
		{
			"array as object",
			"output",
			`"testing"`,
			`[]`,
		},
		{
			"object as array",
			`[0]`,
			`"testing"`,
			`{}`,
		},
		{
			"array too short",
			"[1]",
			`"testing"`,
			`[]`,
		},
	} {
		tt.Run(test.testName, func(t *testing.T) {
			require := require.New(t)
			cmd := &ParsedLine{
				Command:  "set",
				JsonPath: test.jsonPath,
				Value:    ojson.MustNewValueFromJSON(test.v),
			}
			v := ojson.MustNewValueFromJSON(test.initial)
			require.NotNil(ApplyOutputCommand(cmd, &v))
		})
	}
}

func TestAppendInvalid(tt *testing.T) {
	for _, test := range []struct {
		testName string
		jsonPath string
		v        string
		initial  string
	}{
		{
			"invalid path",
			"[",
			`"testing"`,
			"null",
		},
		{
			"array as object",
			"output",
			`"testing"`,
			`[]`,
		},
		{
			"object as array",
			`[0]`,
			`"testing"`,
			`{}`,
		},
		{
			"array too short",
			"[1]",
			`"testing"`,
			`[]`,
		},
		{
			"root not array",
			``,
			`"testing"`,
			`true`,
		},
		{
			"path not array (obj)",
			`a`,
			`""`,
			`{"a":{}}`,
		},
		{
			"path not array (array)",
			`[0]`,
			`""`,
			`[5,[]]`,
		},
	} {
		tt.Run(test.testName, func(t *testing.T) {
			require := require.New(t)
			cmd := &ParsedLine{
				Command:  "append",
				JsonPath: test.jsonPath,
				Value:    ojson.MustNewValueFromJSON(test.v),
			}
			v := ojson.MustNewValueFromJSON(test.initial)
			require.NotNil(ApplyOutputCommand(cmd, &v))
		})
	}
}
