package ast

import "fmt"

//go:generate stringer -type=UnitType token.go unit.go
type UnitType int

const (
	UNIT_PX UnitType = iota
	UNIT_PT
	UNIT_EM
	UNIT_CM
	UNIT_MM
	UNIT_REM
	UNIT_DEG
	UNIT_PERCENT
)

func (unit UnitType) UnitString() string {
	switch unit {
	case UNIT_PX:
		return "px"
	case UNIT_PT:
		return "pt"
	case UNIT_EM:
		return "em"
	case UNIT_CM:
		return "cm"
	case UNIT_MM:
		return "mm"
	case UNIT_REM:
		return "rem"
	case UNIT_DEG:
		return "deg"
	case UNIT_PERCENT:
		return "%"
	default:
		panic(fmt.Errorf("Unsupported unit type: %s", unit))
	}
}

func ConvertTokenTypeToUnitType(tokenType TokenType) UnitType {
	switch tokenType {
	case T_UNIT_PX:
		return UNIT_PX
	case T_UNIT_PT:
		return UNIT_PT
	case T_UNIT_EM:
		return UNIT_EM
	case T_UNIT_CM:
		return UNIT_CM
	case T_UNIT_MM:
		return UNIT_MM
	case T_UNIT_REM:
		return UNIT_REM
	case T_UNIT_DEG:
		return UNIT_DEG
	case T_UNIT_PERCENT:
		return UNIT_PERCENT
	default:
		panic(fmt.Errorf("Unknown Token Type for converting unit type. Got '%s'", tokenType))
	}
}
