/**@type {import('eslint').Linter.Config} */
// eslint-disable-next-line no-undef
module.exports = {
	root: true,
	parser: '@typescript-eslint/parser',
	plugins: ['@typescript-eslint'],
	extends: [
		'eslint:recommended',
		'plugin:@typescript-eslint/recommended',
		'prettier',
		'prettier/@typescript-eslint',
	],
	rules: {
		'no-constant-condition': 0,
		'@typescript-eslint/no-explicit-any': 0,
		'@typescript-eslint/no-unused-vars': 2,
		'@typescript-eslint/no-non-null-assertion': 0,
		'@typescript-eslint/explicit-module-boundary-types': 2,
		'@typescript-eslint/explicit-function-return-type': 2,
	},
};
