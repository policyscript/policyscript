import {Range} from './shared';

export interface Token {
	readonly type: TokenType;
	readonly literal: string;
	readonly range: Range;
}

export enum TokenType {
	ILLEGAL = 'ILLEGAL',
	EOF = 'EOF',
	NEWLINE = 'NEWLINE',
	COMMENT = 'COMMENT',

	// Identifiers
	IDENT = 'IDENT',

	// Literals.
	INTEGER = 'integer',
	DECIMAL = 'decimal',
	MONEY = 'money',
	PERIOD = 'period',
	TEXT = 'text',
	DATE = 'date',
	TIME = 'time',

	// Documentation literals.
	TITLE = 'TITLE',
	PARAGRAPH = 'PARAGRAPH',

	// Operators.
	EQ = '=',
	NOT_EQ = '!=',
	PLUS = '+',
	MINUS = '-',
	MULT = '*',
	DIV = '/',
	LT = '<',
	GT = '>',
	GT_EQ = '>=',
	LT_EQ = '<=',
	TYPE = ':',

	// Delimiters.
	LPAREN = '(',
	RPAREN = ')',
	LBRACE = '{',
	RBRACE = '}',

	// Keywords.
	TRUE = 'true',
	FALSE = 'false',
	IF = 'if',
	FOR = 'for',
	IN = 'in',
	SET = 'set',
	TO = 'to',

	// Blocks.
	META = '@meta',
	DEFINE = '@define',
	ENUM = '@enum',
	INPUTS = '@inputs',
	OUTPUTS = '@outputs',
	LOCALS = '@locals',
	CODE = '@code',
}

export const keywordToTokenType = {
	true: TokenType.TRUE,
	false: TokenType.FALSE,
	if: TokenType.IF,
	for: TokenType.FOR,
	in: TokenType.IN,
	set: TokenType.SET,
	to: TokenType.TO,
} as const;

export const blockToTokenType = {
	'@meta': TokenType.META,
	'@define': TokenType.DEFINE,
	'@enum': TokenType.ENUM,
	'@inputs': TokenType.INPUTS,
	'@outputs': TokenType.OUTPUTS,
	'@locals': TokenType.LOCALS,
	'@code': TokenType.CODE,
} as const;
