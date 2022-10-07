grammar createfile;

file_
   : (local | module | variable | step)+ EOF
   ;

step
   : stepname blockbody
   ;

stepname
   : STRING
   ;

local
  : 'locals' blockbody
  ;

module
  : 'module' name blockbody
  ;

variable
   : VARIABLE name blockbody
   ;

block
   : blocktype label* blockbody
   ;

blocktype
   : IDENTIFIER
   ;

resourcetype
   : STRING
   ;

name
   : STRING
   ;

label
   : STRING
   ;

blockbody
   : LCURL (argument | block)* RCURL
   ;

argument
   : identifier '=' expression
   ;

identifier
   : (('local' | 'data' | 'var' | 'module') DOT)? identifierchain
   ;

identifierchain
   : (IDENTIFIER | IN | VARIABLE | PROVIDER) index? (DOT identifierchain)*
   | STAR (DOT identifierchain)*
   | inline_index (DOT identifierchain)*
   ;

inline_index
   : NATURAL_NUMBER
   ;

expression
   : section
   | expression operator_ expression
   | LPAREN expression RPAREN
   | expression '?' expression ':' expression
   | forloop
   ;

forloop
   : 'for' identifier IN expression ':' expression
   ;

section
   : list_
   | map_
   | val
   ;

val
   : NULL_
   | signed_number
   | string
   | identifier
   | BOOL
   | DESCRIPTION
   | filedecl
   | functioncall
   | EOF_
   ;

functioncall
   : functionname LPAREN functionarguments RPAREN
   | 'jsonencode' LPAREN (.)*? RPAREN
   ;

functionname
   : IDENTIFIER
   ;

functionarguments
   : //no arguments
   | expression (',' expression)*
   ;

index
   : '[' expression ']'
   ;

filedecl
   : 'file' '(' expression ')'
   ;

list_
   : '[' (expression (',' expression)* ','?)? ']'
   ;

map_
   : LCURL (argument ','?)* RCURL
   ;

string
   : STRING
   | MULTILINESTRING
   ;

fragment DIGIT
   : [0-9]
   ;

signed_number
   : ('+' | '-')? number
   ;

VARIABLE
   : 'variable'
   ;

PROVIDER
   : 'provider'
   ;

IN
   : 'in'
   ;

STAR
   : '*'
   ;

DOT
   : '.'
   ;

operator_
   : '/'
   | STAR
   | '%'
   | '+'
   | '-'
   | '>'
   | '>='
   | '<'
   | '<='
   | '=='
   | '!='
   | '&&'
   | '||'
   ;

LCURL
   : '{'
   ;

RCURL
   : '}'
   ;

LPAREN
   : '('
   ;

RPAREN
   : ')'
   ;

EOF_
   : '<<EOF' .*? 'EOF'
   ;

NULL_
   : 'nul'
   ;

NATURAL_NUMBER
   : DIGIT+
   ;

number
   : NATURAL_NUMBER (DOT NATURAL_NUMBER)?
   ;

BOOL
   : 'true'
   | 'false'
   ;

DESCRIPTION
   : '<<DESCRIPTION' .*? 'DESCRIPTION'
   ;

MULTILINESTRING
   : '<<-EOF' .*? 'EOF'
   ;

STRING
   : '"' ( '\\"' | ~["\r\n] )* '"'
   ;
IDENTIFIER
   : [a-zA-Z] ([a-zA-Z0-9_-])*
   ;
COMMENT
  : ('#' | '//') ~ [\r\n]* -> channel(HIDDEN)
  ;

BLOCKCOMMENT
  : '/*' .*? '*/' -> channel(HIDDEN)
  ;

WS
   : [ \r\n\t]+ -> skip
   ;
