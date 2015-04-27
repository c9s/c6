package c6

/*
Statement := RuleSet | At-Rule | Mixin-Statement | FunctionStatement

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

PropertyValue	:= Expression LiteralConcat Expression
				| Expression Expression
				| Expression LiteralConcat(',') Expression
				| Expression
				| Keyword
				| Url

Url := T_URL '(' T_QQ_STRING ')'

Expression := Interpolation
			| Term '+' Term
	        | Term '-' Term
			| Term

Interpolation := "#{" Expression "}"

Term := Factor '*' Factor
		Factor '/' Factor

Factor := Number | Variable

Variable := T_VARIABLE

Scalar := T_NUMBER | T_NUMBER Unit

Unit := T_UNIT_PX | T_UNIT_PT | T_UNIT_EM | T_UNIT_PERCENT | T_UNIT_DEG
*/
