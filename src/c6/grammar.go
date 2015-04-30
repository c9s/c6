package c6

/*
Statement := RuleSet | At-Rule | Mixin-Statement | FunctionStatement | VariableAssignment

VariableAssignment := Variable ':' Value ';'

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

Value := Map
       | List
	   | Expression
	   | Expression LiteralConcat Expression
	   | Keyword
	   | Url
	   | Color

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


Terminals
----------------

Unit := T_UNIT_PX | T_UNIT_PT | T_UNIT_EM | T_UNIT_PERCENT | T_UNIT_DEG

Color := T_HEX_COLOR

Variable := T_VARIABLE

Scalar := T_NUMBER | T_NUMBER Unit
*/
