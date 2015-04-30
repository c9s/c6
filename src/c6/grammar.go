package c6

/*
Statement := RuleSet | At-Rule | Mixin-Statement | FunctionStatement | VariableAssignment

VariableAssignment := Variable ':' Value

At-Rule := '@' T_IDENT ';'

RuleSet := Rule | RuleSet

SelectorList := Selector | Selector ',' SelectorList

Rule := SelectorList '{'
			RuleSet
		'}'

Property := PropertyName ':' PropertyValue

PropertyName :=  PropertyNameToken LiteralConcat PropertyNameToken
			 |   PropertyNameToken

PropertyNameToken := Ident | Interpolation

PropertyValue: List

List := CommaSeparatoredList | '(' CommaSeparatoredList ')'

CommaSeparatoredList := Value T_COMMA CommaSeparatoredList | SpaceSeparatoredList

SpaceSeparatoredList := Value SpaceSeparatoredList

Value		   := Expression LiteralConcat Expression
				| Expression
				| Keyword
				| Url
				| List

Url := T_URL '(' T_QQ_STRING ')'

Expression := Interpolation
			| Term '+' Term
	        | Term '-' Term
			| Term

Interpolation := "#{" Expression "}"

Term := Factor '*' Factor
		Factor '/' Factor

Factor := Number
       | Variable
	   | '(' Expression ')'

Variable := T_VARIABLE

Scalar := T_NUMBER | T_NUMBER Unit

Unit := T_UNIT_PX | T_UNIT_PT | T_UNIT_EM | T_UNIT_PERCENT | T_UNIT_DEG

Color := T_HEX_COLOR

Value := Color
       | '(' List ')'
	   | Map

Value := Map
       | List
	   | Expression

List := '(' CommaSep List ')'
      | CommaSep List

CommaSepList := SpaceSepList ',' CommaSepList
			  | SpaceSepList

SpaceSepList := Value SpaceSepList
              | Value

Map := '(' MapPair ')'

MapPairList := MapPair ',' MapPairList
             | MapPair

MapPair := Expression ':' Value






*/
