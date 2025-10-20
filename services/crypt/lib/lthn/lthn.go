package lthn

// keyMap is the default character-swapping map used for the quasi-salting process.
var keyMap = map[rune]rune{
	'o': '0',
	'l': '1',
	'e': '3',
	'a': '4',
	's': 'z',
	't': '7',
	'0': 'o',
	'1': 'l',
	'3': 'e',
	'4': 'a',
	'7': 't',
}
