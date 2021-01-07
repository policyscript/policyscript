import {Position, Range} from '../shared';
import {Token, TokenType} from '../token';

const whitespaceChars = [' ', '\t', '\r'];
const whitespaceAndBreak = [...whitespaceChars, '\n'];

type Char = string | undefined;

export class Lexer {
	private readonly _input: string;
	private _line: number;
	private _character: number;
	private _absolutePosition: number;
	private _ch: Char;

	constructor(input: string) {
		this._input = input;
		this._line = 0;
		this._character = 0;
		this._absolutePosition = 0;
		this._ch = input[0];
	}

	scan(): Token[] {
		const tokens: Token[] = [];

		// Read the top title.
		tokens.push(this.readTitle());

		// Continue parsing tokens until we reach the end of the file.
		while (!this.isAtEnd()) {
			tokens.push(this.nextToken());
		}

		// Get EOF token.
		const position = this.getPositions();
		tokens.push({type: TokenType.EOF, literal: '', range: {start: position, end: position}});

		return tokens;
	}

	private nextToken(): Token {
		this.skipWhitespaceAndBreaks();

		switch (this._ch) {
			case '-':
				if (this.peek() === ' ') {
					return this.readTitle();
				} else {
					return this.readParagraph();
				}
			case '#':
				return this.readComment();
			case '@':
				return this.blockToken();
			default:
				return this.readParagraph();
		}
	}

	private blockToken(): Token {
		let tok: Token;

		switch (this._ch) {
			// case '+':
			// 	tok = this.newToken(TokenType.PLUS, this._ch);
			// 	break;
			// case '-':
			// 	tok = this.newToken(TokenType.MINUS, this._ch);
			// 	break;
			// case '*':
			// 	tok = this.newToken(TokenType.MULT, this._ch);
			// 	break;
			// case '/':
			// 	tok = this.newToken(TokenType.DEFINE, this._ch);
			// 	break;
			// case '(':
			// 	tok = this.newToken(TokenType.LPAREN, this._ch);
			// 	break;
			// case ')':
			// 	tok = this.newToken(TokenType.RPAREN, this._ch);
			// 	break;
			// case '{':
			// 	tok = this.newToken(TokenType.LBRACE, this._ch);
			// 	break;
			// case '}':
			// 	tok = this.newToken(TokenType.RBRACE, this._ch);
			// 	break;
			default:
				tok = this.singleCharToken(TokenType.ILLEGAL);
		}

		this.advance();
		return tok;
	}

	private readComment(): Token {
		const start = this.getPositions();
		const literalSegments: string[] = [];

		while (true) {
			// Consume `#`.
			this.advance();

			const segmentPosition = this._absolutePosition;

			// Advance until new line.
			while (!this.isAtEnd() && this._ch !== '\n') {
				this.advance();
			}

			const end = this.getPositions();
			literalSegments.push(this._input.slice(segmentPosition, end.position).trim());
			const literal = literalSegments.join(' ');
			const range: Range = {start, end};

			// If end of file, return title token.
			if (this.isAtEnd()) {
				return {type: TokenType.COMMENT, literal, range};
			}

			// Consume `\n`, skip through whitespace.
			this.advance();
			this.skipWhitespace();

			// If end of file, return title token.
			if (this.isAtEnd()) {
				return {type: TokenType.COMMENT, literal, range};
			}

			// If char is not `#`, this is the end of the comment
			if (this._ch !== '#') {
				return {type: TokenType.COMMENT, literal, range};
			}
		}
	}

	private readTitle(): Token {
		return this.readUntilDoubleLineBreak(TokenType.TITLE);
	}

	private readParagraph(): Token {
		return this.readUntilDoubleLineBreak(TokenType.PARAGRAPH);
	}

	private readUntilDoubleLineBreak(type: TokenType): Token {
		this.skipWhitespace();

		const start = this.getPositions();
		const literalSegments: string[] = [];

		while (true) {
			const segmentPosition = this._absolutePosition;

			// Advance until new line.
			while (!this.isAtEnd() && this._ch !== '\n') {
				this.advance();
			}

			const end = this.getPositions();
			literalSegments.push(this._input.slice(segmentPosition, end.position).trim());
			const literal = literalSegments.join(' ');
			const range: Range = {start, end};

			// If end of file, return title token.
			if (this.isAtEnd()) {
				return {type, literal, range};
			}

			// Consume `\n`, skip through whitespace.
			this.advance();
			this.skipWhitespace();

			// If end of file, return title token.
			if (this.isAtEnd()) {
				return {type, literal, range};
			}

			// If char is `\n`, this is a double line break and we
			// end the title.
			if (this._ch === '\n') {
				this.advance();
				return {type, literal, range};
			}
		}
	}

	private skipWhitespace(): void {
		while (!this.isAtEnd() && whitespaceChars.includes(this._ch!)) {
			this.advance();
		}
	}

	private skipWhitespaceAndBreaks(): void {
		while (!this.isAtEnd() && whitespaceAndBreak.includes(this._ch!)) {
			this.advance();
		}
	}

	private advance(): void {
		if (this._ch === '\n') {
			this._line++;
			this._character = -1;
		}

		if (this._absolutePosition + 1 >= this._input.length) {
			this._ch = undefined;
		} else {
			this._ch = this._input[this._absolutePosition + 1];
		}

		this._absolutePosition++;
		this._character++;
	}

	private peek(): Char {
		if (this._absolutePosition + 1 >= this._input.length) {
			return undefined;
		}
		return this._input[this._absolutePosition + 1];
	}

	private isAtEnd(): boolean {
		return this._absolutePosition >= this._input.length;
	}

	private getPositions(): Position {
		return {
			line: this._line,
			character: this._character,
			position: this._absolutePosition,
		};
	}

	private singleCharToken(tokenType: TokenType): Token {
		const start = this.getPositions();
		const literal = this._ch;
		this.advance();
		const end = this.getPositions();

		return {
			type: tokenType,
			literal: literal || '',
			range: {start, end},
		};
	}
}
