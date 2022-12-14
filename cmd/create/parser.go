package main

// Code generated by peg -inline -switch -output cmd/create/parser.go grammar/createfile.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleExpression
	ruleAssignmentExpression
	ruleUnaryExpression
	rulePostfixExpression
	rulePrimaryExpression
	ruleArgumentExpressionList
	ruleIdentifier
	ruleIdNondigit
	ruleIdChar
	ruleKeyword
	ruleStringLiteral
	ruleStringChar
	ruleConstant
	ruleIntegerConstant
	ruleIntegerSuffix
	ruleLsuffix
	ruleDecimalConstant
	ruleFloatConstant
	ruleEscape
	ruleSimpleEscape
	ruleCharacterConstant
	ruleChar
	ruleSpacing
	ruleWhiteSpace
	ruleLPAR
	ruleRPAR
	ruleEQU
	ruleINC
	ruleDEC
	ruleLBRK
	ruleRBRK
	ruleCOMMA
	ruleSEMICOL
)

var rul3s = [...]string{
	"Unknown",
	"Expression",
	"AssignmentExpression",
	"UnaryExpression",
	"PostfixExpression",
	"PrimaryExpression",
	"ArgumentExpressionList",
	"Identifier",
	"IdNondigit",
	"IdChar",
	"Keyword",
	"StringLiteral",
	"StringChar",
	"Constant",
	"IntegerConstant",
	"IntegerSuffix",
	"Lsuffix",
	"DecimalConstant",
	"FloatConstant",
	"Escape",
	"SimpleEscape",
	"CharacterConstant",
	"Char",
	"Spacing",
	"WhiteSpace",
	"LPAR",
	"RPAR",
	"EQU",
	"INC",
	"DEC",
	"LBRK",
	"RBRK",
	"COMMA",
	"SEMICOL",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type Createfile struct {
	Buffer string
	buffer []rune
	rules  [34]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *Createfile) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *Createfile) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Createfile
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *Createfile) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *Createfile) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *Createfile) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func Pretty(pretty bool) func(*Createfile) error {
	return func(p *Createfile) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*Createfile) error {
	return func(p *Createfile) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *Createfile) Init(options ...func(*Createfile) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Expression <- <(AssignmentExpression SEMICOL)*> */
		func() bool {
			{
				position1 := position
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					if !_rules[ruleAssignmentExpression]() {
						goto l3
					}
					{
						position4 := position
						if buffer[position] != rune(';') {
							goto l3
						}
						position++
						if !_rules[ruleSpacing]() {
							goto l3
						}
						add(ruleSEMICOL, position4)
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				add(ruleExpression, position1)
			}
			return true
		},
		/* 1 AssignmentExpression <- <(UnaryExpression EQU AssignmentExpression)> */
		func() bool {
			position5, tokenIndex5 := position, tokenIndex
			{
				position6 := position
				if !_rules[ruleUnaryExpression]() {
					goto l5
				}
				{
					position7 := position
					if buffer[position] != rune('=') {
						goto l5
					}
					position++
					if !_rules[ruleSpacing]() {
						goto l5
					}
					add(ruleEQU, position7)
				}
				if !_rules[ruleAssignmentExpression]() {
					goto l5
				}
				add(ruleAssignmentExpression, position6)
			}
			return true
		l5:
			position, tokenIndex = position5, tokenIndex5
			return false
		},
		/* 2 UnaryExpression <- <(PostfixExpression / (INC UnaryExpression) / (DEC UnaryExpression))> */
		func() bool {
			position8, tokenIndex8 := position, tokenIndex
			{
				position9 := position
				{
					position10, tokenIndex10 := position, tokenIndex
					{
						position12 := position
					l13:
						{
							position14, tokenIndex14 := position, tokenIndex
							{
								position15, tokenIndex15 := position, tokenIndex
								{
									position17 := position
									{
										position18, tokenIndex18 := position, tokenIndex
										{
											position20 := position
											{
												position21, tokenIndex21 := position, tokenIndex
												{
													position22 := position
													{
														position23, tokenIndex23 := position, tokenIndex
														if buffer[position] != rune('s') {
															goto l24
														}
														position++
														if buffer[position] != rune('h') {
															goto l24
														}
														position++
														if buffer[position] != rune('e') {
															goto l24
														}
														position++
														if buffer[position] != rune('l') {
															goto l24
														}
														position++
														if buffer[position] != rune('l') {
															goto l24
														}
														position++
														goto l23
													l24:
														position, tokenIndex = position23, tokenIndex23
														if buffer[position] != rune('o') {
															goto l21
														}
														position++
														if buffer[position] != rune('u') {
															goto l21
														}
														position++
														if buffer[position] != rune('t') {
															goto l21
														}
														position++
														if buffer[position] != rune('p') {
															goto l21
														}
														position++
														if buffer[position] != rune('u') {
															goto l21
														}
														position++
														if buffer[position] != rune('t') {
															goto l21
														}
														position++
													}
												l23:
													{
														position25, tokenIndex25 := position, tokenIndex
														if !_rules[ruleIdChar]() {
															goto l25
														}
														goto l21
													l25:
														position, tokenIndex = position25, tokenIndex25
													}
													add(ruleKeyword, position22)
												}
												goto l19
											l21:
												position, tokenIndex = position21, tokenIndex21
											}
											{
												position26 := position
												{
													switch buffer[position] {
													case '_':
														if buffer[position] != rune('_') {
															goto l19
														}
														position++
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l19
														}
														position++
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l19
														}
														position++
													}
												}

												add(ruleIdNondigit, position26)
											}
										l28:
											{
												position29, tokenIndex29 := position, tokenIndex
												if !_rules[ruleIdChar]() {
													goto l29
												}
												goto l28
											l29:
												position, tokenIndex = position29, tokenIndex29
											}
											if !_rules[ruleSpacing]() {
												goto l19
											}
											add(ruleIdentifier, position20)
										}
										goto l18
									l19:
										position, tokenIndex = position18, tokenIndex18
										{
											switch buffer[position] {
											case '(':
												if !_rules[ruleLPAR]() {
													goto l16
												}
												if !_rules[ruleExpression]() {
													goto l16
												}
												if !_rules[ruleRPAR]() {
													goto l16
												}
											case '"':
												{
													position31 := position
													if buffer[position] != rune('"') {
														goto l16
													}
													position++
												l34:
													{
														position35, tokenIndex35 := position, tokenIndex
														{
															position36 := position
															{
																position37, tokenIndex37 := position, tokenIndex
																if !_rules[ruleEscape]() {
																	goto l38
																}
																goto l37
															l38:
																position, tokenIndex = position37, tokenIndex37
																{
																	position39, tokenIndex39 := position, tokenIndex
																	{
																		switch buffer[position] {
																		case '\\':
																			if buffer[position] != rune('\\') {
																				goto l39
																			}
																			position++
																		case '\n':
																			if buffer[position] != rune('\n') {
																				goto l39
																			}
																			position++
																		default:
																			if buffer[position] != rune('"') {
																				goto l39
																			}
																			position++
																		}
																	}

																	goto l35
																l39:
																	position, tokenIndex = position39, tokenIndex39
																}
																if !matchDot() {
																	goto l35
																}
															}
														l37:
															add(ruleStringChar, position36)
														}
														goto l34
													l35:
														position, tokenIndex = position35, tokenIndex35
													}
													if buffer[position] != rune('"') {
														goto l16
													}
													position++
													if !_rules[ruleSpacing]() {
														goto l16
													}
												l32:
													{
														position33, tokenIndex33 := position, tokenIndex
														if buffer[position] != rune('"') {
															goto l33
														}
														position++
													l41:
														{
															position42, tokenIndex42 := position, tokenIndex
															{
																position43 := position
																{
																	position44, tokenIndex44 := position, tokenIndex
																	if !_rules[ruleEscape]() {
																		goto l45
																	}
																	goto l44
																l45:
																	position, tokenIndex = position44, tokenIndex44
																	{
																		position46, tokenIndex46 := position, tokenIndex
																		{
																			switch buffer[position] {
																			case '\\':
																				if buffer[position] != rune('\\') {
																					goto l46
																				}
																				position++
																			case '\n':
																				if buffer[position] != rune('\n') {
																					goto l46
																				}
																				position++
																			default:
																				if buffer[position] != rune('"') {
																					goto l46
																				}
																				position++
																			}
																		}

																		goto l42
																	l46:
																		position, tokenIndex = position46, tokenIndex46
																	}
																	if !matchDot() {
																		goto l42
																	}
																}
															l44:
																add(ruleStringChar, position43)
															}
															goto l41
														l42:
															position, tokenIndex = position42, tokenIndex42
														}
														if buffer[position] != rune('"') {
															goto l33
														}
														position++
														if !_rules[ruleSpacing]() {
															goto l33
														}
														goto l32
													l33:
														position, tokenIndex = position33, tokenIndex33
													}
													add(ruleStringLiteral, position31)
												}
											default:
												{
													position48 := position
													{
														position49, tokenIndex49 := position, tokenIndex
														{
															position51 := position
															{
																position52, tokenIndex52 := position, tokenIndex
																if c := buffer[position]; c < rune('1') || c > rune('9') {
																	goto l53
																}
																position++
															l54:
																{
																	position55, tokenIndex55 := position, tokenIndex
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l55
																	}
																	position++
																	goto l54
																l55:
																	position, tokenIndex = position55, tokenIndex55
																}
																if buffer[position] != rune('.') {
																	goto l53
																}
																position++
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l53
																}
																position++
															l56:
																{
																	position57, tokenIndex57 := position, tokenIndex
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l57
																	}
																	position++
																	goto l56
																l57:
																	position, tokenIndex = position57, tokenIndex57
																}
																goto l52
															l53:
																position, tokenIndex = position52, tokenIndex52
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l50
																}
																position++
															l58:
																{
																	position59, tokenIndex59 := position, tokenIndex
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l59
																	}
																	position++
																	goto l58
																l59:
																	position, tokenIndex = position59, tokenIndex59
																}
																if buffer[position] != rune('.') {
																	goto l50
																}
																position++
															}
														l52:
															add(ruleFloatConstant, position51)
														}
														goto l49
													l50:
														position, tokenIndex = position49, tokenIndex49
														{
															position61 := position
															{
																position62 := position
																if c := buffer[position]; c < rune('1') || c > rune('9') {
																	goto l60
																}
																position++
															l63:
																{
																	position64, tokenIndex64 := position, tokenIndex
																	if c := buffer[position]; c < rune('0') || c > rune('9') {
																		goto l64
																	}
																	position++
																	goto l63
																l64:
																	position, tokenIndex = position64, tokenIndex64
																}
																add(ruleDecimalConstant, position62)
															}
															{
																position65, tokenIndex65 := position, tokenIndex
																{
																	position67 := position
																	{
																		position68, tokenIndex68 := position, tokenIndex
																		{
																			position70, tokenIndex70 := position, tokenIndex
																			if buffer[position] != rune('u') {
																				goto l71
																			}
																			position++
																			goto l70
																		l71:
																			position, tokenIndex = position70, tokenIndex70
																			if buffer[position] != rune('U') {
																				goto l69
																			}
																			position++
																		}
																	l70:
																		{
																			position72, tokenIndex72 := position, tokenIndex
																			if !_rules[ruleLsuffix]() {
																				goto l72
																			}
																			goto l73
																		l72:
																			position, tokenIndex = position72, tokenIndex72
																		}
																	l73:
																		goto l68
																	l69:
																		position, tokenIndex = position68, tokenIndex68
																		if !_rules[ruleLsuffix]() {
																			goto l65
																		}
																		{
																			position74, tokenIndex74 := position, tokenIndex
																			{
																				position76, tokenIndex76 := position, tokenIndex
																				if buffer[position] != rune('u') {
																					goto l77
																				}
																				position++
																				goto l76
																			l77:
																				position, tokenIndex = position76, tokenIndex76
																				if buffer[position] != rune('U') {
																					goto l74
																				}
																				position++
																			}
																		l76:
																			goto l75
																		l74:
																			position, tokenIndex = position74, tokenIndex74
																		}
																	l75:
																	}
																l68:
																	add(ruleIntegerSuffix, position67)
																}
																goto l66
															l65:
																position, tokenIndex = position65, tokenIndex65
															}
														l66:
															if !_rules[ruleSpacing]() {
																goto l60
															}
															add(ruleIntegerConstant, position61)
														}
														goto l49
													l60:
														position, tokenIndex = position49, tokenIndex49
														{
															position78 := position
															{
																position79, tokenIndex79 := position, tokenIndex
																if buffer[position] != rune('L') {
																	goto l79
																}
																position++
																goto l80
															l79:
																position, tokenIndex = position79, tokenIndex79
															}
														l80:
															if buffer[position] != rune('\'') {
																goto l16
															}
															position++
														l81:
															{
																position82, tokenIndex82 := position, tokenIndex
																{
																	position83 := position
																	{
																		position84, tokenIndex84 := position, tokenIndex
																		if !_rules[ruleEscape]() {
																			goto l85
																		}
																		goto l84
																	l85:
																		position, tokenIndex = position84, tokenIndex84
																		{
																			position86, tokenIndex86 := position, tokenIndex
																			{
																				switch buffer[position] {
																				case '\\':
																					if buffer[position] != rune('\\') {
																						goto l86
																					}
																					position++
																				case '\n':
																					if buffer[position] != rune('\n') {
																						goto l86
																					}
																					position++
																				default:
																					if buffer[position] != rune('\'') {
																						goto l86
																					}
																					position++
																				}
																			}

																			goto l82
																		l86:
																			position, tokenIndex = position86, tokenIndex86
																		}
																		if !matchDot() {
																			goto l82
																		}
																	}
																l84:
																	add(ruleChar, position83)
																}
																goto l81
															l82:
																position, tokenIndex = position82, tokenIndex82
															}
															if buffer[position] != rune('\'') {
																goto l16
															}
															position++
															if !_rules[ruleSpacing]() {
																goto l16
															}
															add(ruleCharacterConstant, position78)
														}
													}
												l49:
													add(ruleConstant, position48)
												}
											}
										}

									}
								l18:
									add(rulePrimaryExpression, position17)
								}
								goto l15
							l16:
								position, tokenIndex = position15, tokenIndex15
								{
									switch buffer[position] {
									case '-':
										if !_rules[ruleDEC]() {
											goto l14
										}
									case '+':
										if !_rules[ruleINC]() {
											goto l14
										}
									case '(':
										if !_rules[ruleLPAR]() {
											goto l14
										}
										{
											position89, tokenIndex89 := position, tokenIndex
											{
												position91 := position
												if !_rules[ruleAssignmentExpression]() {
													goto l89
												}
											l92:
												{
													position93, tokenIndex93 := position, tokenIndex
													{
														position94 := position
														if buffer[position] != rune(',') {
															goto l93
														}
														position++
														if !_rules[ruleSpacing]() {
															goto l93
														}
														add(ruleCOMMA, position94)
													}
													if !_rules[ruleAssignmentExpression]() {
														goto l93
													}
													goto l92
												l93:
													position, tokenIndex = position93, tokenIndex93
												}
												add(ruleArgumentExpressionList, position91)
											}
											goto l90
										l89:
											position, tokenIndex = position89, tokenIndex89
										}
									l90:
										if !_rules[ruleRPAR]() {
											goto l14
										}
									default:
										{
											position95 := position
											if buffer[position] != rune('[') {
												goto l14
											}
											position++
											if !_rules[ruleSpacing]() {
												goto l14
											}
											add(ruleLBRK, position95)
										}
										if !_rules[ruleExpression]() {
											goto l14
										}
										{
											position96 := position
											if buffer[position] != rune(']') {
												goto l14
											}
											position++
											if !_rules[ruleSpacing]() {
												goto l14
											}
											add(ruleRBRK, position96)
										}
									}
								}

							}
						l15:
							goto l13
						l14:
							position, tokenIndex = position14, tokenIndex14
						}
						add(rulePostfixExpression, position12)
					}
					goto l10

					position, tokenIndex = position10, tokenIndex10
					if !_rules[ruleINC]() {
						goto l97
					}
					if !_rules[ruleUnaryExpression]() {
						goto l97
					}
					goto l10
				l97:
					position, tokenIndex = position10, tokenIndex10
					if !_rules[ruleDEC]() {
						goto l8
					}
					if !_rules[ruleUnaryExpression]() {
						goto l8
					}
				}
			l10:
				add(ruleUnaryExpression, position9)
			}
			return true
		l8:
			position, tokenIndex = position8, tokenIndex8
			return false
		},
		/* 3 PostfixExpression <- <(PrimaryExpression / ((&('-') DEC) | (&('+') INC) | (&('(') (LPAR ArgumentExpressionList? RPAR)) | (&('[') (LBRK Expression RBRK))))*> */
		nil,
		/* 4 PrimaryExpression <- <(Identifier / ((&('(') (LPAR Expression RPAR)) | (&('"') StringLiteral) | (&('\'' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'L') Constant)))> */
		nil,
		/* 5 ArgumentExpressionList <- <(AssignmentExpression (COMMA AssignmentExpression)*)> */
		nil,
		/* 6 Identifier <- <(!Keyword IdNondigit IdChar* Spacing)> */
		nil,
		/* 7 IdNondigit <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		nil,
		/* 8 IdChar <- <((&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position103, tokenIndex103 := position, tokenIndex
			{
				position104 := position
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l103
						}
						position++
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l103
						}
						position++
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l103
						}
						position++
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l103
						}
						position++
					}
				}

				add(ruleIdChar, position104)
			}
			return true
		l103:
			position, tokenIndex = position103, tokenIndex103
			return false
		},
		/* 9 Keyword <- <((('s' 'h' 'e' 'l' 'l') / ('o' 'u' 't' 'p' 'u' 't')) !IdChar)> */
		nil,
		/* 10 StringLiteral <- <('"' StringChar* '"' Spacing)+> */
		nil,
		/* 11 StringChar <- <(Escape / (!((&('\\') '\\') | (&('\n') '\n') | (&('"') '"')) .))> */
		nil,
		/* 12 Constant <- <(FloatConstant / IntegerConstant / CharacterConstant)> */
		nil,
		/* 13 IntegerConstant <- <(DecimalConstant IntegerSuffix? Spacing)> */
		nil,
		/* 14 IntegerSuffix <- <((('u' / 'U') Lsuffix?) / (Lsuffix ('u' / 'U')?))> */
		nil,
		/* 15 Lsuffix <- <(('l' 'l') / ('L' 'L') / ('l' / 'L'))> */
		func() bool {
			position112, tokenIndex112 := position, tokenIndex
			{
				position113 := position
				{
					position114, tokenIndex114 := position, tokenIndex
					if buffer[position] != rune('l') {
						goto l115
					}
					position++
					if buffer[position] != rune('l') {
						goto l115
					}
					position++
					goto l114
				l115:
					position, tokenIndex = position114, tokenIndex114
					if buffer[position] != rune('L') {
						goto l116
					}
					position++
					if buffer[position] != rune('L') {
						goto l116
					}
					position++
					goto l114
				l116:
					position, tokenIndex = position114, tokenIndex114
					{
						position117, tokenIndex117 := position, tokenIndex
						if buffer[position] != rune('l') {
							goto l118
						}
						position++
						goto l117
					l118:
						position, tokenIndex = position117, tokenIndex117
						if buffer[position] != rune('L') {
							goto l112
						}
						position++
					}
				l117:
				}
			l114:
				add(ruleLsuffix, position113)
			}
			return true
		l112:
			position, tokenIndex = position112, tokenIndex112
			return false
		},
		/* 16 DecimalConstant <- <([1-9] [0-9]*)> */
		nil,
		/* 17 FloatConstant <- <(([1-9] [0-9]* '.' [0-9]+) / ([0-9]+ '.'))> */
		nil,
		/* 18 Escape <- <SimpleEscape> */
		func() bool {
			position121, tokenIndex121 := position, tokenIndex
			{
				position122 := position
				{
					position123 := position
					if buffer[position] != rune('\\') {
						goto l121
					}
					position++
					{
						switch buffer[position] {
						case 'v':
							if buffer[position] != rune('v') {
								goto l121
							}
							position++
						case 't':
							if buffer[position] != rune('t') {
								goto l121
							}
							position++
						case 'r':
							if buffer[position] != rune('r') {
								goto l121
							}
							position++
						case 'n':
							if buffer[position] != rune('n') {
								goto l121
							}
							position++
						case 'f':
							if buffer[position] != rune('f') {
								goto l121
							}
							position++
						case 'b':
							if buffer[position] != rune('b') {
								goto l121
							}
							position++
						case 'a':
							if buffer[position] != rune('a') {
								goto l121
							}
							position++
						case '%':
							if buffer[position] != rune('%') {
								goto l121
							}
							position++
						case '\\':
							if buffer[position] != rune('\\') {
								goto l121
							}
							position++
						case '?':
							if buffer[position] != rune('?') {
								goto l121
							}
							position++
						case '"':
							if buffer[position] != rune('"') {
								goto l121
							}
							position++
						default:
							if buffer[position] != rune('\'') {
								goto l121
							}
							position++
						}
					}

					add(ruleSimpleEscape, position123)
				}
				add(ruleEscape, position122)
			}
			return true
		l121:
			position, tokenIndex = position121, tokenIndex121
			return false
		},
		/* 19 SimpleEscape <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('%') '%') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 20 CharacterConstant <- <('L'? '\'' Char* '\'' Spacing)> */
		nil,
		/* 21 Char <- <(Escape / (!((&('\\') '\\') | (&('\n') '\n') | (&('\'') '\'')) .))> */
		nil,
		/* 22 Spacing <- <WhiteSpace*> */
		func() bool {
			{
				position129 := position
			l130:
				{
					position131, tokenIndex131 := position, tokenIndex
					{
						position132 := position
						{
							switch buffer[position] {
							case '\t':
								if buffer[position] != rune('\t') {
									goto l131
								}
								position++
							case '\r':
								if buffer[position] != rune('\r') {
									goto l131
								}
								position++
							case '\n':
								if buffer[position] != rune('\n') {
									goto l131
								}
								position++
							default:
								if buffer[position] != rune(' ') {
									goto l131
								}
								position++
							}
						}

						add(ruleWhiteSpace, position132)
					}
					goto l130
				l131:
					position, tokenIndex = position131, tokenIndex131
				}
				add(ruleSpacing, position129)
			}
			return true
		},
		/* 23 WhiteSpace <- <((&('\t') '\t') | (&('\r') '\r') | (&('\n') '\n') | (&(' ') ' '))> */
		nil,
		/* 24 LPAR <- <('(' Spacing)> */
		func() bool {
			position135, tokenIndex135 := position, tokenIndex
			{
				position136 := position
				if buffer[position] != rune('(') {
					goto l135
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l135
				}
				add(ruleLPAR, position136)
			}
			return true
		l135:
			position, tokenIndex = position135, tokenIndex135
			return false
		},
		/* 25 RPAR <- <(')' Spacing)> */
		func() bool {
			position137, tokenIndex137 := position, tokenIndex
			{
				position138 := position
				if buffer[position] != rune(')') {
					goto l137
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l137
				}
				add(ruleRPAR, position138)
			}
			return true
		l137:
			position, tokenIndex = position137, tokenIndex137
			return false
		},
		/* 26 EQU <- <('=' Spacing)> */
		nil,
		/* 27 INC <- <('+' '+' Spacing)> */
		func() bool {
			position140, tokenIndex140 := position, tokenIndex
			{
				position141 := position
				if buffer[position] != rune('+') {
					goto l140
				}
				position++
				if buffer[position] != rune('+') {
					goto l140
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l140
				}
				add(ruleINC, position141)
			}
			return true
		l140:
			position, tokenIndex = position140, tokenIndex140
			return false
		},
		/* 28 DEC <- <('-' '-' Spacing)> */
		func() bool {
			position142, tokenIndex142 := position, tokenIndex
			{
				position143 := position
				if buffer[position] != rune('-') {
					goto l142
				}
				position++
				if buffer[position] != rune('-') {
					goto l142
				}
				position++
				if !_rules[ruleSpacing]() {
					goto l142
				}
				add(ruleDEC, position143)
			}
			return true
		l142:
			position, tokenIndex = position142, tokenIndex142
			return false
		},
		/* 29 LBRK <- <('[' Spacing)> */
		nil,
		/* 30 RBRK <- <(']' Spacing)> */
		nil,
		/* 31 COMMA <- <(',' Spacing)> */
		nil,
		/* 32 SEMICOL <- <(';' Spacing)> */
		nil,
	}
	p.rules = _rules
	return nil
}
