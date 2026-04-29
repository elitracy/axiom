package parser

import (
	"fmt"
	"strings"
)

const (
	keySystem  = "system"
	keySet     = "set"
	keyConnect = "connect"
	keyType    = "type"
	keyArrow   = "->"
	keyComment = "//"
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

func (e parseError) String() string {
	return fmt.Sprintf("line %d: %s", e.line, e.msg)
}

func newParseError(line int, msg string) parseError {
	return parseError{line, msg}
}

type Parser struct {
	Config ParserConfig
	tokens []token
}

func NewParser(config ParserConfig) *Parser {
	return &Parser{
		Config: config,
	}
}

func (p *Parser) Parse(content []byte) error {
	var errors []parseError

	lines := strings.Split(string(content), "\n")

	for row, line := range lines {
		tokens := strings.Fields(line)
		if len(tokens) == 0 {
			continue
		}

		switch tokens[0] {
		case keyComment:
			continue
		case keySystem:
			if len(tokens) != 3 {
				errors = append(errors, newParseError(row, "invalid system declaration"))
				continue
			}

			name := tokens[1]

			typeDec := strings.Split(tokens[2], "=")
			if len(typeDec) != 2 {
				errors = append(errors, newParseError(row, "invalid type declaration"))
				continue
			}

			systemType := typeDec[1]

			p.Config.SubsystemDeclarations[name] = systemType

		case keySet:
			if len(tokens) != 3 {
				errors = append(errors, newParseError(row, "invalid set directive"))
				continue
			}

			systemComponent := strings.Split(tokens[1], ".")
			if len(systemComponent) != 2 {
				errors = append(errors, newParseError(row, "invalid set directive"))
				continue
			}

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
				errors = append(errors, newParseError(row, "invalid connection declaration"))
				continue
			}

			srcNamePort := strings.Split(tokens[1], ".")
			if len(srcNamePort) != 3 {
				errors = append(errors, newParseError(row, "invalid connection source declaration"))
				continue
			}

			destNamePort := strings.Split(tokens[3], ".")
			if len(destNamePort) != 3 {
				errors = append(errors, newParseError(row, "invalid connection destination declaration"))
				continue
			}

			srcSystem := srcNamePort[0]
			srcPort := fmt.Sprintf("%s.%s", srcNamePort[1], srcNamePort[2])

			destSystem := destNamePort[0]
			destPort := fmt.Sprintf("%s.%s", destNamePort[1], destNamePort[2])

			throughput := tokens[4]

			connection := connectionDeclaration{
				srcSystem,
				srcPort,
				destSystem,
				destPort,
				throughput,
			}

			p.Config.ConnectionDeclarations = append(p.Config.ConnectionDeclarations, connection)
		default:
			error := fmt.Sprintf("Unknown symbol: %v", tokens[0])
			errors = append(errors, newParseError(row, error))
		}
	}

	p.Config.Errors = errors
	return nil

}
