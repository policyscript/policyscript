export interface Position {
	/**
	 * 1-indexed line number.
	 */
	line: number;

	/**
	 * 0-indexed column number.
	 */
	column: number;

	/**
	 * 0-indexed offset in one-line string.
	 */
	offset: number;
}

export interface Range {
	start: Position;
	end: Position;
}
