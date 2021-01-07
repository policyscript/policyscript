export interface Position {
	/**
	 * 0-based line index.
	 */
	line: number;

	/**
	 * 0-based line offset.
	 */
	character: number;

	/**
	 * Absolute position in input string.
	 */
	position: number;
}

export interface Range {
	start: Position;
	end: Position;
}
