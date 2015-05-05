package ast

import "fmt"

// import "strings"

//go:generate stringer -type=UnitType token.go unit.go
type UnitType int

const (
	UNIT_NONE UnitType = iota

	/*
		Length Unit
		@see https://developer.mozilla.org/en-US/docs/Web/CSS/length
	*/
	UNIT_EM
	UNIT_EX
	UNIT_CH
	UNIT_REM

	// Absolute length
	UNIT_CM
	UNIT_IN
	UNIT_MM
	UNIT_PC
	UNIT_PT
	UNIT_PX

	// Viewport-percentage lengths
	UNIT_VH
	UNIT_VW
	UNIT_VMIN
	UNIT_VMAX

	/*
		Angle
		@see https://developer.mozilla.org/en-US/docs/Web/CSS/angle
	*/
	UNIT_DEG
	UNIT_GRAD
	UNIT_RAD
	UNIT_TURN

	UNIT_PERCENT

	/*
		Time Unit
		@see https://developer.mozilla.org/zh-TW/docs/Web/CSS/time
	*/
	UNIT_SECOND
	UNIT_MILLISECOND

	/*
		Resolution Unit
		@see https://developer.mozilla.org/en-US/docs/Web/CSS/resolution
	*/
	UNIT_DPI
	UNIT_DPPX
	UNIT_DPCM

	/*
		Frequency unit
		@see https://developer.mozilla.org/en-US/docs/Web/CSS/frequency
	*/
	UNIT_HZ
	UNIT_KHZ
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
	case UNIT_SECOND:
		return "s"
	case UNIT_MILLISECOND:
		return "ms"
	case UNIT_DPI:
		return "dpi"
	case UNIT_DPPX:
		return "dppx"
	case UNIT_DPCM:
		return "dpcm"
	case UNIT_NONE:
		return ""
	default:
		panic(fmt.Errorf("Unsupported unit type: %s", unit))
		/*
			// For undefined unit, convert the unit name dynamically
			var str = unit.String()
			return strings.ToLower(str[len("UNIT_"):])
		*/
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
	case T_UNIT_SECOND:
		return UNIT_SECOND
	case T_UNIT_MILLISECOND:
		return UNIT_MILLISECOND
	case T_UNIT_PERCENT:
		return UNIT_PERCENT
	case T_UNIT_DPI:
		return UNIT_DPI
	case T_UNIT_DPPX:
		return UNIT_DPPX
	case T_UNIT_DPCM:
		return UNIT_DPCM
	default:
		panic(fmt.Errorf("Unknown Token Type for converting unit type. Got '%s'", tokenType))
	}
}
