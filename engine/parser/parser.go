package parser

import (
	"fmt"
	"os"
	"strings"
)

const (
	keySystem  = "system"
	keySet     = "set"
	keyConnect = "connect"
	keyType    = "type"
	keyArrow   = "->"
)

type token struct {
	line    int
	content string
}

func newToken(line int, content string) token {
	return token{line, content}
}

type parseError struct {
	line int
	msg  string
}

func newParseError(line int, msg string) parseError {
	return parseError{line, msg}
}

type Parser struct {
	Config ParserConfig
	tokens []token
	errors []parseError
}

func NewParser(config ParserConfig) Parser {
	return Parser{
		Config: config,
	}
}

func (p *Parser) Parse(content []byte) error {
	lines := strings.Split(string(content), "\n")

	for row, line := range lines {
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		switch tokens[0] {
		case keySystem:
			if len(tokens) != 3 {
				p.errors = append(p.errors, newParseError(row, "invalid system declaration"))
				continue
			}

			name := tokens[1]

			typeDec := strings.Split(tokens[2], "=")
			systemType := typeDec[1]

			p.Config.SubsystemDeclarations[name] = systemType

		case keySet:
			if len(tokens) != 3 {
				p.errors = append(p.errors, newParseError(row, "invalid set directive"))
				continue
			}

			systemComponent := strings.Split(tokens[1], ".")

			system := systemComponent[0]
			component := systemComponent[1]
			value := tokens[2]

			p.Config.SetDirectives = append(p.Config.SetDirectives,
				setDirective{
					system,
					component,
					value,
				})

		case keyConnect:
			if len(tokens) != 5 {
				p.errors = append(p.errors, newParseError(row, "invalid connection declaration"))
				continue
			}

			srcSystemPort := strings.Split(tokens[1], ".")
			srcSystem := srcSystemPort[0]
			srcPort := srcSystemPort[1]

			destSystemPort := strings.Split(tokens[3], ".")
			destSystem := destSystemPort[0]
			destPort := destSystemPort[1]

			throughput := tokens[4]

			connection := connectionDeclaration{
				srcSystem,
				srcPort,
				destSystem,
				destPort,
				throughput,
			}

			p.Config.ConnectionDeclarations = append(p.Config.ConnectionDeclarations, connection)
		}
	}

	return nil

}

func (p *Parser) ReadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Could not read file: %s", path)
	}
	return file, nil

}
