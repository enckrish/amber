package main

import (
	"fmt"
	"github.com/satyrius/gonx"
	"plugin"
)

type ParserPlugin interface {
	parseLog(format string, log string) (string, error)
}

func ParseLog(format string, log string) (string, error) {
	parser := gonx.NewParser(format)
	res, err := parser.ParseString(log)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return fmt.Sprintf("%+v", res), nil
}

func ParseLogUsingPy(string, string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

type TParserFn func(string, string) (string, error)

// GetPluginParser Returns parser function from Go Plugin, format argument is set to empty string
func GetPluginParser(pluginPath string) (TParserFn, error) {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	symParser, err := plug.Lookup("parseLog")
	if err != nil {
		return nil, err
	}

	parser, ok := symParser.(ParserPlugin)
	if !ok {
		return nil, fmt.Errorf("unexpected type from module symbol")
	}

	return parser.parseLog, nil
}

func ParseUsingFn(parserFn TParserFn, format string, log string) (string, error) {
	return parserFn(format, log)
}
