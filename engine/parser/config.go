package parser

import "github.com/elias/axiom/engine/utils"

type setDirective struct {
	System    string
	Component string
	Value     string
}

type connectionDeclaration struct {
	SrcSystem  string
	SrcPort    string
	DestSystem string
	DestPort   string
	Throughput string
}

type ParserConfig struct {
	SubsystemDeclarations  map[utils.SubsystemName]utils.SubsystemType
	SetDirectives          []setDirective
	ConnectionDeclarations []connectionDeclaration
	Errors                 []parseError
}

func NewParserConfig() ParserConfig {
	return ParserConfig{
		SubsystemDeclarations: make(map[utils.SubsystemName]utils.SubsystemType),
	}
}
