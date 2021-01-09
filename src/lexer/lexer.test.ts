import {Token, TokenType} from '../token';
import {Lexer} from './lexer';

describe('lexer', () => {
	it.each([
		[``, ``],
		[`\n`, ``],
		[`\n\n`, ``],
		[`\nC`, ``],
		[`A`, `A`],
		[`A\n`, `A`],
		[`A\n\n`, `A`],
		[`A\nB`, `A\nB`],
		[`A\nB\n`, `A\nB`],
		[`A\nB\n\n`, `A\nB`],
		[`A\nB\n\nC`, `A\nB`],
		[`A\n\tB`, `A\n\tB`],
		[`A\n  B`, `A\n  B`],
		[`A\n\nB`, `A`],
	])('initial heading: parses %j for heading %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.HEADING,
			literal: expected,
		};
		console.log(tokens[0]);
		expect(tokens[0]).toMatchObject(testToken);
	});

	it.each([
		[`\n- A`, `- A`],
		[`\n- A\n`, `- A`],
		[`\n- A\n\n`, `- A`],
		[`\n- A\nB`, `- A\nB`],
		[`\n- A\nB\n`, `- A\nB`],
		[`\n- A\nB\n\n`, `- A\nB`],
		[`\n- A\nB\n\nC`, `- A\nB`],
		[`\n- A\n\tB`, `- A\n\tB`],
		[`\n- A\n  B`, `- A\n  B`],
		[`\n- A`, `- A`],
		[`\n- -A`, `- -A`],
		[`\n - - A`, `- - A`],
	])('secondary heading: parses %j for heading %j', (input, expected) => {
		const lexer = new Lexer(input);
		const tokens = lexer.scan();

		const testToken: Partial<Token> = {
			type: TokenType.HEADING,
			literal: expected,
		};
		expect(tokens[1]).toMatchObject(testToken);
	});

	it.each([
		[`\nA`, `A`],
		[`\nA\n`, `A`],
		[`\nA\n\n`, `A`],
		[`\nA\nB`, `A\nB`],
		[`\nA\nB\n`, `A\nB`],
		[`\nA\nB\n\n`, `A\nB`],
		[`\nA\nB\n\nC`, `A\nB`],
		[`\nA\n\tB`, `A\n\tB`],
		[`\nA\n  B`, `A\n  B`],
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
		[`\n#A`, `A`],
		[`\n#A\n#B`, `A\nB`],
		[`\n# A`, ` A`],
		[`\n# A\n`, ` A`],
		[`\n# A\n\n`, ` A`],
		[`\n# A\nC`, ` A`],
		[`\n# A\n# B`, ` A\n B`],
		[`\n# A\n# B\n`, ` A\n B`],
		[`\n# A\n# B\n\n`, ` A\n B`],
		[`\n# A\n# B\n\n# C`, ` A\n B`],
		[`\n# A\n# \tB`, ` A\n \tB`],
		[`\n# A\n#   B`, ` A\n   B`],
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
