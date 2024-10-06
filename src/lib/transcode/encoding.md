# Sango Syllabic Encoding

Internally, Sango is stored as syllables (which are the basic phonemic unit).

This has numerous advantages:

- Easy to convert into and out of UTF8, validate syllables, and avoid invalid Sango lexemes
- Easy to query different properties and mask unimportant properties
- Easy to iterate over without worrying about byte boundaries, and supports random access
- Compact notation with low entropy: suitable as vector embedding in machine learning algorithms
- Code switching is trivially easy with Language embedded into syllable
- Metadata can be easily embedded by setting Case to Hidden

## Encoding format

The 16-bit encoding divides up quasi-orthogonally into components (from MSB to LSB):

| # bits | Description                       |
| :----: | --------------------------------- |
|   1    | Component type                    |
|   2    | Language code                     |
|   2    | Typecase (capitalization)         |
|   5    | Consonants                        |
|   4    | Vowel (incl. height and nasality) |
|   2    | Pitch accent                      |

The bit vector will thus look like a 16-bit binary integer: `0bCLLTTCCCCCVVVVPP`.

Two big advantages to this specific encoding are that:

- It is easy to distinguish words from punctuation/whitespace just by inspecting the the high-order bit.
- If the high-order bit is 0, then the value is automatically a Unicode rune, easily
  converted to or from UTF8 using the [unicode/utf8](https://pkg.go.dev/unicode/utf8) library.

## Components

### Component type: 1 bit

| Bit F | 1                       |
| :---: | ----------------------- |
|   0   | Unicode (U+0000-U+7FFF) |
|   1   | Syllable                |

### Language code: 2 bits

| Bits E \\ D | 0       | 1       |
| :---------: | ------- | ------- |
|      0      | Unknown | Sango   |
|      1      | English | French  |

### Typecase: 2 bits

| Bits C \\ B | 0         | 1         |
| :---------: | --------- | --------- |
|      0      | Lowercase | Uppercase |
|      1      | Titlecase | Hidden    |

### Consonants: 5 bits

| Bits A98 \\ 76 | 00  | 01  | 10     | 11     |
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

| Bits 54 \\ 32 | 00     | 01     | 10       | 11 x     |
| :-----------: | ------ | ------ | -------- | -------- |
|      00       | a      | i      | o        | e        |
|      01       | an     | in     | on       | en       |
|      10       | un     | u      | ɔ \| ✖   | ɛ \| ∅   |
|      11       | ✖ \| à | ✖ \| ù | o/ɔ \| è | e/ɛ \| é |

### Pitch accent: 2 bits

| Bits 1 \\ 0 | 0                         | 1                                 |
| :---------: | ------------------------- | --------------------------------- |
|      0      | Unknown/None: (∅)         | Low: zero-breaking space (U+200b) |
|      1      | Mid: diaeresis (U+0308)   | High: circumflex (U+-302)         |
