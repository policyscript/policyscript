import {Position, Range} from '../shared';
import {Token, TokenType} from '../token';

const whitespaceChars: (undefined | string)[] = [' ', '\t', '\r'];
const whitespaceAndBreak: (undefined | string)[] = [...whitespaceChars, '\n'];

export class Lexer {
	private readonly _input: string;
	private _line: number;
	private _column: number;
	private _offset: number;
	private _ch: string | undefined;

	constructor(input: string) {
		this._input = input;
		this._line = 1; // Line is 1-indexed.
		this._column = 0;
		this._offset = 0;
		this._ch = input[0];
	}

	scan(): Token[] {
		const tokens: Token[] = [];

		// Read the top heading.
		tokens.push(this.readHeading());

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
					return this.readHeading();
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
		let literal = '';

		while (true) {
			// Consume `#`.
			this.advance();

			const segmentPosition = this._offset;

			// Advance until new line.
			while (!this.isAtEnd() && this._ch !== '\n') {
				this.advance();
			}

			const end = this.getPositions();
			literal += this._input.slice(segmentPosition, end.offset);
			const range: Range = {start, end};

			// If end of file, return token.
			if (this.isAtEnd()) {
				return {type: TokenType.COMMENT, literal, range};
			}

			// Consume `\n`.
			this.advance();

			// If end of file, return token.
			if (this.isAtEnd()) {
				return {type: TokenType.COMMENT, literal, range};
			}

			// Advance until non-whitespace.
			while (!this.isAtEnd() && whitespaceChars.includes(this._ch)) {
				this.advance();
			}

			// If char is not `#`, this is the end of the comment.
			if (this._ch !== '#') {
				return {type: TokenType.COMMENT, literal, range};
			}

			// Otherwise, add in the line break we skipped.
			literal += '\n';
		}
	}

	private readHeading(): Token {
		return this.readUntilDoubleLineBreak(TokenType.HEADING);
	}

	private readParagraph(): Token {
		return this.readUntilDoubleLineBreak(TokenType.PARAGRAPH);
	}

	/**
	 * Will collect literal until it hits a double line break, comment or
	 * EOF.
	 * @param type The token type to create.
	 * @param prevStart Provide this if there is a previous start point to
	 * collect from, otherwise will start at current position.
	 */
	private readUntilDoubleLineBreak(type: TokenType, prevStart?: Position): Token {
		this.skipWhitespace();

		let start = this.getPositions();
		let literal = '';

		// If previous start is defined, get literal until current
		// position, then set start to previous start.
		if (prevStart) {
			literal += this._input.slice(prevStart.offset, start.offset);
			start = prevStart;
		}

		let end = this.getPositions();
		let range: Range = {start, end};

		while (true) {
			const previousOffset = this._offset;

			let whitespaceOnlyLine = true;

			// Advance until new line.
			while (!this.isAtEnd() && this._ch !== '\n') {
				// Return early if only whitespace chars and
				// a comment char appears.
				if (whitespaceOnlyLine && this._ch === '#') {
					return {type, literal, range};
				}

				// If char is non-whitespace, set bool to false.
				if (!whitespaceChars.includes(this._ch)) {
					whitespaceOnlyLine = false;
				}
				this.advance();
			}

			// Update position if non whitespace only line.
			// Otherwise we use positions from previous end line.
			if (!whitespaceOnlyLine) {
				end = this.getPositions();

				// Prefix with line break if literal is empty
				// string, since it is a new iteration.
				const prefix = literal ? '\n' : '';
				literal += prefix + this._input.slice(previousOffset, end.offset);
				range = {start, end};
			}

			// If end of file, return token.
			if (this.isAtEnd()) {
				return {type, literal, range};
			}

			// Consume `\n`.
			this.advance();

			// Return since we had a whitespace only line aka a
			// double line break.
			if (whitespaceOnlyLine) {
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
			this._column = -1;
		}

		if (this._offset + 1 >= this._input.length) {
			this._ch = undefined;
		} else {
			this._ch = this._input[this._offset + 1];
		}

		this._offset++;
		this._column++;
	}

	private peek(): string | undefined {
		if (this._offset + 1 >= this._input.length) {
			return undefined;
		}
		return this._input[this._offset + 1];
	}

	private isAtEnd(): boolean {
		return this._offset >= this._input.length;
	}

	private getPositions(): Position {
		return {
			line: this._line,
			column: this._column,
			offset: this._offset,
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
