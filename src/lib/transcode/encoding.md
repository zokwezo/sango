# Sango Syllabic Encoding

Internally, Sango is stored as syllables (which are the basic phonemic unit).

This has numerous advantages:

- Easy to convert into and out of UTF8, validate syllables, and hard to inadvertently create invalid Sango lexemes
- Easy to query different properties and mask unimportant properties
- Easy to iterate over without worrying about byte boundaries, and supports random access
- Serves as its own compact and semantically meaningful vector embedding in machine learning algorithms
- Compact notation with low entropy: common important cases are small, but full Unicode expressibility is available when needed
- Code switching is trivially easy: easy to ignore non-Sango or convert row format to column format by
  masking on language for interlingual translations
- Metadata can be easily embedded by setting Case to Hidden, allowing for preserving annotations inline as documents move through a pipleline

## Encoding format

The 16-bit encoding divides up quasi-orthogonally into components (from LSB to MSB):

| # bits | Description                       |
| :----: | --------------------------------- |
|   5    | Consonants                        |
|   4    | Vowel (incl. height and nasality) |
|   2    | Pitch accent                      |
|   2    | Case (upper, lower)               |
|   2    | Language code                     |
|   1    | Type                              |

The bit vector will thus look like a 16-bit binary integer: `0bTLLCCPPVVVVCCCCC`.

Two big advantages to this specific encoding are that:

- It is easy to distinguish words from punctuation/whitespace just by inspecting the the high-order bit.
- If the high-order bit is 0, then the value is automatically a Unicode rune, easily
  converted to or from UTF8 using the [unicode/utf8](https://pkg.go.dev/unicode/utf8) library.

## Components

In the table below:

- where there is avertical bar, the value before it is for Sango, the value after for non-Sango
- ∅ indicates that the component is intentionally absent
- ✖ indicates that the value is invalid (for that Language code)
- In non-Sango syllables, circumflex, diaeresis, and macron are used.

### Consonants: 5 bits

| Bits 432 \\ 10 | 00  | 01  | 10     | 11     |
| :------------: | --- | --- | ------ | ------ |
|      000       | ∅   | f   | r      | k      |
|      001       | mv  | v   | ng     | g      |
|      010       | m   | p   | l      | kp     |
|      011       | mb  | b   | ngb    | gb     |
|      100       | ç   | s   | y      | h      |
|      101       | nz  | z   | ny     | w      |
|      110       | n   | t   | ✖ \| c | ✖ \| x |
|      111       | nd  | d   | ✖ \| j | ✖ \| q |

### Vowel: 4 bits

| Bits 87 \\ 65 | 00     | 01     | 10       | 11 x     |
| :-----------: | ------ | ------ | -------- | -------- |
|      00       | a      | i      | o        | e        |
|      01       | an     | in     | on       | en       |
|      10       | un     | u      | ɔ \| ✖   | ɛ \| ∅   |
|      11       | ✖ \| à | ✖ \| ù | o/ɔ \| è | e/ɛ \| é |

### Pitch accent: 2 bits

| Bits A \\ 9 | 0                    | 1                      |
| :---------: | -------------------- | ---------------------- |
|      0      | Low \| ∅             | High \| circumflex (^) |
|      1      | Mid \| diaeresis (¨) | ?? \| macron (¯)       |

### Case: 2 bits

| Bits C \\ B | 0         | 1         |
| :---------: | --------- | --------- |
|      0      | Lowercase | Uppercase |
|      1      | Titlecase | Hidden    |

### Language code: 2 bits

| Bits E \\ D | 0       | 1       |
| :---------: | ------- | ------- |
|      0      | Sango   | French  |
|      1      | English | Unknown |

### Component type: 1 bit

| Bit F | 1                       |
| :---: | ----------------------- |
|   0   | Unicode (U+0000-U+7FFF) |
|   1   | Syllable                |
