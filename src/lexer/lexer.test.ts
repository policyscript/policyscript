import {Token, TokenType} from '../token';
import {Lexer} from './lexer';

describe('lexer', () => {
	it.each([
		[`Title`, `Title`],
		[`Title\n`, `Title`],
		[`Title\n\n`, `Title`],
		[`Title\nnext line`, `Title next line`],
		[`Title\nnext line\n`, `Title next line`],
		[`Title\nnext line\n\n`, `Title next line`],
		[`Title\n\nnext line`, `Title`],
	])('initial title: parses %j for title %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.TITLE,
			literal: expected,
		};
		expect(tokens[0]).toMatchObject(testToken);
	});

	it.each([
		[`Title\n\n- (a) Title 2`, `- (a) Title 2`],
		[`Title\n\n- (a) Title 2\n`, `- (a) Title 2`],
		[`Title\n\n- (a) Title 2\n\n`, `- (a) Title 2`],
		[`Title\n\n- (a) Title 2\nnext line`, `- (a) Title 2 next line`],
		[`Title\n\n- (a) Title 2\nnext line\n`, `- (a) Title 2 next line`],
		[`Title\n\n- (a) Title 2\nnext line\n\n`, `- (a) Title 2 next line`],
		[`Title\n\n- (a) Title 2`, `- (a) Title 2`],
		[`Title\n\n- - (a) Title 2`, `- - (a) Title 2`],
		[`Title\n\n - - (a) Title 2`, `- - (a) Title 2`],
	])('secondary title: parses %j for title %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.TITLE,
			literal: expected,
		};
		expect(tokens[1]).toMatchObject(testToken);
	});

	it.each([
		[`Title\n\nThis is a paragraph`, `This is a paragraph`],
		[`Title\n\nThis is a paragraph\n`, `This is a paragraph`],
		[`Title\n\nThis is a paragraph\n\n`, `This is a paragraph`],
		[`Title\n\nThis is a paragraph\nsecond line`, `This is a paragraph second line`],
		[`Title\n\nThis is a paragraph\nsecond line\n`, `This is a paragraph second line`],
		[`Title\n\nThis is a paragraph\nsecond line\n\n`, `This is a paragraph second line`],
	])('paragraph: parses %j for paragraph %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.PARAGRAPH,
			literal: expected,
		};
		expect(tokens[1]).toMatchObject(testToken);
	});

	it.each([
		[`Title\n\n# This is a comment`, `This is a comment`],
		[`Title\n\n# This is a comment\n`, `This is a comment`],
		[`Title\n\n# This is a comment\n\n`, `This is a comment`],
		[`Title\n\n# This is a comment\nignore this, not a comment`, `This is a comment`],
		[`Title\n\n# This is a comment\n# second line`, `This is a comment second line`],
		[`Title\n\n# This is a comment\n# second line\n`, `This is a comment second line`],
		[`Title\n\n# This is a comment\n# second line\n\n`, `This is a comment second line`],
	])('comment: parses %j for comment %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.COMMENT,
			literal: expected,
		};
		expect(tokens[1]).toMatchObject(testToken);
	});
});
