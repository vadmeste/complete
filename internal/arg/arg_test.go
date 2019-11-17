package arg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line string
		args []Arg
	}{
		{
			line: "a b",
			args: []Arg{{Text: "a", Completed: true}, {Text: "b", Completed: false}},
		},
		{
			line: " a b ",
			args: []Arg{{Text: "a", Completed: true}, {Text: "b", Completed: true}},
		},
		{
			line: "a  b",
			args: []Arg{{Text: "a", Completed: true}, {Text: "b", Completed: false}},
		},
		{
			line: " a ",
			args: []Arg{{Text: "a", Completed: true}},
		},
		{
			line: " a",
			args: []Arg{{Text: "a", Completed: false}},
		},
		{
			line: "  ",
			args: nil,
		},
		{
			line: "",
			args: nil,
		},
		{
			line: `\ a\ b c\ `,
			args: []Arg{{Text: `\ a\ b`, Completed: true}, {Text: `c\ `, Completed: false}},
		},
		{
			line: `"\"'\''" '"'`,
			args: []Arg{{Text: `"\"'\''"`, Completed: true}, {Text: `'"'`, Completed: false}},
		},
		{
			line: `"a b"`,
			args: []Arg{{Text: `"a b"`, Completed: false}},
		},
		{
			line: `"a b" `,
			args: []Arg{{Text: `"a b"`, Completed: true}},
		},
		{
			line: `"a b"c`,
			args: []Arg{{Text: `"a b"c`, Completed: false}},
		},
		{
			line: `"a b"c `,
			args: []Arg{{Text: `"a b"c`, Completed: true}},
		},
		{
			line: `"a b" c`,
			args: []Arg{{Text: `"a b"`, Completed: true}, {Text: "c", Completed: false}},
		},
		{
			line: `"a `,
			args: []Arg{{Text: `"a `, Completed: false}},
		},
		{
			line: `\"a b`,
			args: []Arg{{Text: `\"a`, Completed: true}, {Text: "b", Completed: false}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			args := Parse(tt.line)
			// Clear parsed part of the arguments. It is tested in the TestArgsParsed test.
			for i := range args {
				arg := args[i]
				arg.Parsed = Parsed{}
				args[i] = arg
			}
			assert.Equal(t, tt.args, args)
		})
	}
}

func TestArgsParsed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		text   string
		parsed Parsed
	}{
		{text: "-", parsed: Parsed{Dashes: "-", HasFlag: true}},
		{text: "--", parsed: Parsed{Dashes: "--", HasFlag: true}},
		{text: "---"}, // Forbidden.
		{text: "--="}, // Forbidden.
		{text: "-="},  // Forbidden.
		{text: "-a-b", parsed: Parsed{Dashes: "-", Flag: "a-b", HasFlag: true}},
		{text: "--a-b", parsed: Parsed{Dashes: "--", Flag: "a-b", HasFlag: true}},
		{text: "-a-b=c-d=e", parsed: Parsed{Dashes: "-", Flag: "a-b", HasFlag: true, Value: "c-d=e", HasValue: true}},
		{text: "--a-b=c-d=e", parsed: Parsed{Dashes: "--", Flag: "a-b", HasFlag: true, Value: "c-d=e", HasValue: true}},
		{text: "--a-b=", parsed: Parsed{Dashes: "--", Flag: "a-b", HasFlag: true, Value: "", HasValue: true}},
		{text: "a", parsed: Parsed{Value: "a", HasValue: true}},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			arg := Parse(tt.text)[0]
			assert.Equal(t, tt.parsed, arg.Parsed)
		})
	}
}
