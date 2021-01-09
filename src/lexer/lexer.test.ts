import {Token, TokenType} from '../token';
import {Lexer} from './lexer';

describe('lexer', () => {
	it.each(
		// prettier-ignore
		[
		[`Title`,                            `Title`],
		[`Title\n`,                          `Title`],
		[`Title\n\n`,                        `Title`],
		[`Title\nnext line`,                 `Title\nnext line`],
		[`Title\nnext line\n`,               `Title\nnext line`],
		[`Title\nnext line\n\n`,             `Title\nnext line`],
		[`Title\nnext line\n\nnot involved`, `Title\nnext line`],
		[`Title\n\tnext line`,               `Title\n\tnext line`],
		[`Title\n  next line`,               `Title\n  next line`],
		[`Title\n\nnext line`,               `Title`],
		],
	)('initial title: parses %j for title %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.TITLE,
			literal: expected,
		};
		expect(tokens[0]).toMatchObject(testToken);
	});

	it.each(
		// prettier-ignore
		[
		[`Title\n\n- (a) Title 2`,                `- (a) Title 2`],
		[`Title\n\n- (a) Title 2\n`,              `- (a) Title 2`],
		[`Title\n\n- (a) Title 2\n\n`,            `- (a) Title 2`],
		[`Title\n\n- (a) Title 2\nnext line`,     `- (a) Title 2\nnext line`],
		[`Title\n\n- (a) Title 2\nnext line\n`,   `- (a) Title 2\nnext line`],
		[`Title\n\n- (a) Title 2\nnext line\n\n`, `- (a) Title 2\nnext line`],
		[`Title\n\n- (a) A\nB\n\nC`,              `- (a) A\nB`],
		[`Title\n\n- (a) Title 2\n\tnext line`,   `- (a) Title 2\n\tnext line`],
		[`Title\n\n- (a) Title 2\n  next line`,   `- (a) Title 2\n  next line`],
		[`Title\n\n- (a) Title 2`,                `- (a) Title 2`],
		[`Title\n\n- - (a) Title 2`,              `- - (a) Title 2`],
		[`Title\n\n - - (a) Title 2`,             `- - (a) Title 2`],
		],
	)('secondary title: parses %j for title %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.TITLE,
			literal: expected,
		};
		expect(tokens[1]).toMatchObject(testToken);
	});

	it.each(
		// prettier-ignore
		[
		[`Title\n\nThis is a paragraph`,                  `This is a paragraph`],
		[`Title\n\nThis is a paragraph\n`,                `This is a paragraph`],
		[`Title\n\nThis is a paragraph\n\n`,              `This is a paragraph`],
		[`Title\n\nThis is a paragraph\nsecond line`,     `This is a paragraph\nsecond line`],
		[`Title\n\nThis is a paragraph\nsecond line\n`,   `This is a paragraph\nsecond line`],
		[`Title\n\nThis is a paragraph\nsecond line\n\n`, `This is a paragraph\nsecond line`],
		[`Title\n\nA\nB\n\nC`,                            `A\nB`],
		[`Title\n\nThis is a paragraph\n\tsecond line`,   `This is a paragraph\n\tsecond line`],
		[`Title\n\nThis is a paragraph\n  second line`,   `This is a paragraph\n  second line`],
		],
	)('paragraph: parses %j for paragraph %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.PARAGRAPH,
			literal: expected,
		};
		expect(tokens[1]).toMatchObject(testToken);
	});

	it.each(
		// prettier-ignore
		[
		[`Title\n\n#This is a comment`,                              `This is a comment`],
		[`Title\n\n#This is a comment\n#second line`,                `This is a comment\nsecond line`],
		[`Title\n\n# This is a comment`,                             ` This is a comment`],
		[`Title\n\n# This is a comment\n`,                           ` This is a comment`],
		[`Title\n\n# This is a comment\n\n`,                         ` This is a comment`],
		[`Title\n\n# This is a comment\nignore this, not a comment`, ` This is a comment`],
		[`Title\n\n# This is a comment\n# second line`,              ` This is a comment\n second line`],
		[`Title\n\n# This is a comment\n# second line\n`,            ` This is a comment\n second line`],
		[`Title\n\n# This is a comment\n# second line\n\n`,          ` This is a comment\n second line`],
		[`Title\n\n# A\n# B\n\n# C`,                                 ` A\n B`],
		[`Title\n\n# This is a comment\n# \tsecond line`,            ` This is a comment\n \tsecond line`],
		[`Title\n\n# This is a comment\n#   second line`,            ` This is a comment\n   second line`],
		],
	)('comment: parses %j for comment %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.COMMENT,
			literal: expected,
		};
		expect(tokens[1]).toMatchObject(testToken);
	});
});
