package kal

import (
	"fmt"
)

type OutputFormat string

const (
	OutputCLI   OutputFormat = "cli"
	OutputJSON  OutputFormat = "json"
	OutputNeo4j OutputFormat = "neo4j"
)

func (of *OutputFormat) String() string {
	return string(*of)
}

func (of *OutputFormat) Set(v string) error {
	switch v {
	case "cli", "json", "neo4j":
		*of = OutputFormat(v)
		return nil
	default:
		return fmt.Errorf(
			"must be one of: %s",
			[]OutputFormat{OutputNeo4j, OutputCLI, OutputJSON},
		)
	}
}

func (of *OutputFormat) Type() string {
	return "kal.OutputFormat"
}
