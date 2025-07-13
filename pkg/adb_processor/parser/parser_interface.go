// Package parser provides interfaces and implementations for parsing the output of various ADB commands.
package parser

type (
	IParser interface {
		Parse(string) error
	}
)
