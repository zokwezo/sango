# Transcoding text

There are various ways to represent text in different contexts, each with its own advantages and disadvantages.

_Transcoding_ is the process of switching from one representation to another (with no loss of information or meaning), and code in this directory provides functions to do so.

This is distinguished from _transliteration_, which is a purely syntactic remapping that may not preserve the semantics of the text.

The various schemes are described in detail below.

## Unicode and UTF8

A _glyph_ (or _grapheme_) is any printable text that fits within a single-spaced column. Interlinear text can be most easily vertically aligned when encoded as sequences of glyphs.

Each glyph is composed of one or more _runes_ (or _codepoints_) which are 32-bit indexes into Unicode table entries.

- In this library, all runes can be addressed with 14 bits, and can represent any codepoint in the range [U+0000..U+3FFF].
- Each accented vowel may be represented with a base rune and a combining accent rune using an [NFD format](https://unicode.org/reports/tr15/#Norm_Forms).
- For some of these, there is also a precomposed single-rune NFC format, and this is always preferred where available.
- _Sequences of runes arise when converting between glyphs and UTF8, and are otherwise not of much use on their own._

Each rune can be serialized into a [UTF8](https://en.wikipedia.org/wiki/UTF-8) byte sequence of variable length, and these can be concatenated into strings,
which is the internet interchange standard. Well-formed UTF8 can be uniquely decomposed back into runes when needed.

One unfortunate side effect is that different glyph representations have varying lengths, depending on the glyph:

| Glyph | # Glyphs | # Runes | # Bytes |
| :---: | :------: | :-----: | :-----: |
|   e   |    1     |    1    |    1    |
|   ɛ   |    1     |    1    |    2    |
|   ê   |    1     |    2    |    3    |
|   ɛ̂   |    1     |    2    |    4    |

This means that whereas a sequence of glyphs in a string can be randomly accessed via index and its length determined in constant time,
this is not possible when represented as a sequence of Unicode runes or as a UTF8 string, and these must be iterated over.

> Go provides native functionality to iterate over runes, but this causes problems with multirune glyphs (those with combining marks).\
> Fortunately, there is a third-party [Unicode Text Segmentation library](https://github.com/rivo/uniseg/blob/master/README.md) that
> does correctly iterate over glyphs, whatever the representation.

## Alternate representations

For particular use cases (described below), other representations are more useful.

### Internal representation

For efficiency, text is stored internally as a sequence of 16-bit tokens (called _SSE tokens_), each of which may represent one of:

- a single Unicode rune (U+0000..U+3FFF)
  - For convenience, the token is defined such that when representing a Unicode rune, the token has the same numeric value as the rune value.
- an ASCII character and its location within an English or French word
- a Sango syllable and its syntactic properties and location within a Sango word

This makes it easier to:

- work parametrically across languages
- access text by index without iteration
- orthogonally query or set text properties
- pass data through middleware without worrying about escaping or conventions
- use directly as a dense vector embedding for use in machine learning algorithms.

See the documentation on [Sango Syllabic Encoding](./encoding.md) for full details.

### Output representations

Sango consonants are all ASCII, but the vowels need to reflect pitch and height and therefore require non-ASCII runes.

The standard orthography is insufficient to fully represent vowel height (this being supplied pragmatically by speakers
based on context and lexicon), and also cannot encode uncertainty in pitch or height. Therefore, the following enriched
output format is used to represent:

- 4 cases: UPPERCASE, Titlecase, lowercase, and -case (pronounced "hyphencase")
- 7 vowels of varying height (and two of unknown height), indicated in the unaccented glyph:
  - 5 close (long) vowels written in ASCII: `a`, `e`, `i`, `o`, `u`.
  - 2 open (short) vowels written in [IPA symbols](https://en.wikipedia.org/wiki/International_Phonetic_Alphabet): `ɛ`, `ɔ`.
    - _This is readily understood notation, but in the standard orthography `e` and `o` represent both close and open vowels._
  - 2 vowels of unknown height: `ə`, `ø`.
    - _This is completely nonstandard notation, for use (and useful) only internal to this library for text in an intermediate state
      which eventually will be resolved into `e`|`ɛ` or `o`|`ɔ`, respectively._
- 3 pitch levels (and one of unknown pitch), indicated using accents on the vowel:
  - Low, unmarked (e.g. `a`)
  - Medium, marked with diaeresis (e.g. `ä`)
  - High, marked with circumflex (e.g. `â`)
  - Unknown, marked with dot below (e.g. `ạ`)
    - _This is completely nonstandard notation, for use (and useful) only internal to this library for text in an intermediate state
      which eventually will be resolved into one of the three pitches above once known._

Some vowel/accent combinations cannot be expressed as a single Unicode rune and require a second combining rune to form a single glyph.

This enriched output format can be easily post-processed into Sango standard orthography, but as this is a lossy mapping, it is done only just before export for external use.

#### Harnessing case as a redundant pitch signal

Accent marks are regrettably small and difficult to see when reading written Sango text, and especially difficult both for nonnative speakers and the visually challenged.

Conversely, uppercase provides little signal beyond what punctuation already provides.

Consequently, pitch can be made more obvious by using case as a redundant pitch signal[^1]:

1. Convert all text to lowercase
2. Convert syllable-final nasal `n` to uppercase.
3. Convert all mid or high pitch vowels to uppercase
4. For high pitch vowels, also convert any preceding consonants to uppercase.

This format may be particularly useful for nonnative speakers when reading Sango text aloud. Anecdotally[^2], mapping visual letter size to voice pitch seems to happen at a lower level in the brain than mapping diacritics and (unlike the latter) can be readily mastered after minimal training even by those with no Sango knowledge at all.

[^1]: This format is novel and therefore completely nonstandard, but if it is found useful and adopted externally, a link back to [this page](https://github.com/zokwezo/sango/src/lib/transcode/README.md) would be appreciated.
[^2]: This hypothesis has not been empirically verified nor the effect quantified, but I invite anyone interested to investigate this further and send me a link to your published results.

Compare reading the following two sentences aloud quickly using the correct pitch pattern ˧ ˥ ˥ ˩ ˧ ˥:

- `Mbï yê nî ahön kûɛ̂.`
- `mbÏ YÊ NÎ ahÖN KÛƐ̂.`

### Input representations

It is tedious to input non-ASCII glyphs on a keyboard.

For convenience, in addition to using UTF8 directly, the ASCII characters `\`, `c`, `j`, `q`, and `x` (which are not otherwise used in Sango)
may be used to encode unicode in any word where the result is consistent with Sango phonemics:

1. A backslash can be used to escape any Unicode rune
   - `\x`_HH_ for hex digits _HH_, encodes U+00HH (where `x` is lowercase)
   - `\X`_HHHH_ for hex digits _HHHH_, encodes U+HHHH (where `X` is uppercase)
2. Vowels
   - If the paragraph does not have a prior open vowel, automatically demote close vowel to unknown height:
     - `E` ⟹ `Ə`
     - `e` ⟹ `ə`
     - `O` ⟹ `Ø`
     - `o` ⟹ `ø`
   - Transliterate open vowels:
     - `X` ⟹ `Ɛ`
     - `x` ⟹ `ɛ`
     - `C` ⟹ `Ɔ`
     - `c` ⟹ `ɔ`
   - After any vowel _v_, transliterate mid (`q`) or high (`j`) pitch modifier to a diacritic:
     - _v_`q` ⟹ v̈
     - _v_`j` ⟹ v̂
   - If the paragraph does not have a prior open vowel, automatically demote low pitch to unknown pitch:
     - _v_ ⟹ ṿ
3. Punctuation
   - Double quotes:
     - `"` ⟹ `“` (left double quote) and `”` (right double quote), alternating in a sentence.
     - `<<` ⟹ `«` (left angular brackets)
     - `>>` ⟹ `»` (right angular brackets)
   - Connectors:
     - `...` ⟹ `…` (ellipsis)
     - `--` ⟹ `–` (n-dash)
     - `---` ⟹ `—` (m-dash)

## Metadata

Sometimes, there is a need to provide metadata inline that is not intended as literal text. Use cases include:

- semantic annotations
- parsing/control directives
- changes to default behavior in the current parsing context
- language tags or annotations, e.g. to mark code switching
- inline translations or alternative spellings
- comments

### Syntax

All braces in the text must correctly nest. This is to ensure that the start and end of any metadata block can be detected lexically without backtracking.

- When scanning left to right, metadata starts with the first `{` and ending with its matching `}`, including both braces.

### Semantics

- Brace literals can be specified with `{LEFTBRACE}` and `{RIGHTBRACE}` respectively.
  - These do not need to nest or be paired.
  - The intended semantics is to restore `{` and `}` literals to the text on output
- The semantics of all other metadata may be content and context dependent, and not further specified here.
  - Applications that do not recognize metadata are expected to retain it as an inert payload and otherwise ignore it.
  - Text sinks may excise all metadata before exporting the text.

## APPENDIX: Unicode glyphs used in Sango text

There are many good Unicode references online, including e.g.

- https://www.compart.com/en/unicode/ for individual runes
- https://util.unicode.org/UnicodeJsps/list-unicodeset.jsp when constructing regular expressions

For reference, all unicode glyphs (and their encodings) used in this library are listed below.

All glyphs can be represented in an NFD form, with base rune (NFD1 column), sometimes followed by combining mark rune (NFD2 column).

For many of these glyphs, there is also a single precomposed Unicode rune (NFC column), and this is always preferred when available.

### Upper case

| NFC UTF8 |  NFC   |  NFD1  |  NFD2  | Height  | Pitch   | ASCII |
| :------: | :----: | :----: | :----: | :------ | :------ | :---: |
|    A     | U+0041 | U+0041 |        | Close   | Low     |   A   |
|    Ä     | U+00C4 | U+0041 | U+0308 | Close   | Mid     |   A   |
|    Â     | U+00C2 | U+0041 | U+0302 | Close   | High    |   A   |
|    Ạ     | U+1EA0 | U+0041 | U+0323 | Close   | Unknown |   A   |
|    E     | U+0045 | U+0045 |        | Close   | Low     |   E   |
|    Ë     | U+00CB | U+0045 | U+0308 | Close   | Mid     |   E   |
|    Ê     | U+00CA | U+0045 | U+0302 | Close   | High    |   E   |
|    Ẹ     | U+1EB8 | U+0045 | U+0323 | Close   | Unknown |   E   |
|    Ɛ     | U+0190 | U+0190 |        | Open    | Low     |   E   |
|    Ɛ̈     |   ❌   | U+0190 | U+0308 | Open    | Mid     |   E   |
|    Ɛ̂     |   ❌   | U+0190 | U+0302 | Open    | High    |   E   |
|    Ɛ̣     |   ❌   | U+0190 | U+0323 | Open    | Unknown |   E   |
|    Ə     | U+018F | U+018F |        | Unknown | Low     |   E   |
|    Ə̈     |   ❌   | U+018F | U+0308 | Unknown | Mid     |   E   |
|    Ə̂     |   ❌   | U+018F | U+0302 | Unknown | High    |   E   |
|    Ə̣     |   ❌   | U+018F | U+0323 | Unknown | Unknown |   E   |
|    I     | U+0049 | U+0049 |        | Close   | Low     |   I   |
|    Ï     | U+00CF | U+0049 | U+0308 | Close   | Mid     |   I   |
|    Î     | U+00CE | U+0049 | U+0302 | Close   | High    |   I   |
|    Ị     | U+1ECA | U+0049 | U+0323 | Close   | Unknown |   I   |
|    O     | U+004F | U+004F |        | Close   | Low     |   O   |
|    Ö     | U+00D6 | U+004F | U+0308 | Close   | Mid     |   O   |
|    Ô     | U+00D4 | U+004F | U+0302 | Close   | High    |   O   |
|    Ọ     | U+1ECC | U+004F | U+0323 | Close   | Unknown |   O   |
|    Ɔ     | U+0186 | U+0186 |        | Open    | Low     |   O   |
|    Ɔ̈     |   ❌   | U+0186 | U+0308 | Open    | Mid     |   O   |
|    Ɔ̂     |   ❌   | U+0186 | U+0302 | Open    | High    |   O   |
|    Ɔ     |   ❌   | U+0186 | U+0323 | Open    | Unknown |   O   |
|    Ø     | U+00D8 | U+00D8 |        | Unknown | Low     |   O   |
|    Ø̈     |   ❌   | U+00D8 | U+0308 | Unknown | Mid     |   O   |
|    Ø̂     |   ❌   | U+00D8 | U+0302 | Unknown | High    |   O   |
|    Ø̣     |   ❌   | U+00D8 | U+0323 | Unknown | Unknown |   O   |
|    U     | U+0055 | U+0055 |        | Close   | Low     |   U   |
|    Ü     | U+00DC | U+0055 | U+0308 | Close   | Mid     |   U   |
|    Û     | U+00DB | U+0055 | U+0302 | Close   | High    |   U   |
|    Ụ     | U+1EE4 | U+0055 | U+0323 | Close   | Unknown |   U   |

### Lower case

| NFC UTF8 |  NFC   |  NFD1  |  NFD2  | Height  | Pitch   | ASCII |
| :------: | :----: | :----: | :----: | :------ | :------ | :---: |
|    a     | U+0061 | U+0061 |        | Close   | Low     |   a   |
|    ä     | U+00E4 | U+0061 | U+0308 | Close   | Mid     |   a   |
|    â     | U+00E2 | U+0061 | U+0302 | Close   | High    |   a   |
|    ạ     | U+1EA1 | U+0061 | U+0323 | Close   | Unknown |   a   |
|    e     | U+0065 | U+0065 |        | Close   | Low     |   e   |
|    ë     | U+00EB | U+0065 | U+0308 | Close   | Mid     |   e   |
|    ê     | U+00EA | U+0065 | U+0302 | Close   | High    |   e   |
|    ẹ     | U+1EB9 | U+0065 | U+0323 | Close   | Unknown |   e   |
|    ɛ     | U+025B | U+025B |        | Open    | Low     |   e   |
|    ɛ̈     |   ❌   | U+025B | U+0308 | Open    | Mid     |   e   |
|    ɛ̂     |   ❌   | U+025B | U+0302 | Open    | High    |   e   |
|    ɛ̣     |   ❌   | U+025B | U+0323 | Open    | Unknown |   e   |
|    ə     | U+0259 | U+0259 |        | Unknown | Low     |   e   |
|    ə̈     |   ❌   | U+0259 | U+0308 | Unknown | Mid     |   e   |
|    ə̂     |   ❌   | U+0259 | U+0302 | Unknown | High    |   e   |
|    ə̣     |   ❌   | U+0259 | U+0323 | Unknown | Unknown |   e   |
|    i     | U+0069 | U+0069 |        | Close   | Low     |   i   |
|    ï     | U+00EF | U+0069 | U+0308 | Close   | Mid     |   i   |
|    î     | U+00EE | U+0069 | U+0302 | Close   | High    |   i   |
|    ị     | U+1ECB | U+0069 | U+0323 | Close   | Unknown |   i   |
|    o     | U+006F | U+006F |        | Close   | Low     |   o   |
|    ö     | U+00F6 | U+006F | U+0308 | Close   | Mid     |   o   |
|    ô     | U+00F4 | U+006F | U+0302 | Close   | High    |   o   |
|    ọ     | U+1ECD | U+006F | U+0323 | Close   | Unknown |   o   |
|    ɔ     | U+0254 | U+0254 |        | Open    | Low     |   o   |
|    ɔ̈     |   ❌   | U+0254 | U+0308 | Open    | Mid     |   o   |
|    ɔ̂     |   ❌   | U+0254 | U+0302 | Open    | High    |   o   |
|    ɔ     |   ❌   | U+0254 | U+0323 | Open    | Unknown |   o   |
|    ø     | U+00F8 | U+00F8 |        | Unknown | Low     |   o   |
|    ø̈     |   ❌   | U+00F8 | U+0308 | Unknown | Mid     |   o   |
|    ø̂     |   ❌   | U+00F8 | U+0302 | Unknown | High    |   o   |
|    ø̣     |   ❌   | U+00F8 | U+0323 | Unknown | Unknown |   o   |
|    u     | U+0075 | U+0075 |        | Close   | Low     |   u   |
|    ü     | U+00FC | U+0075 | U+0308 | Close   | Mid     |   u   |
|    û     | U+00FB | U+0075 | U+0302 | Close   | High    |   u   |
|    ụ     | U+1EE5 | U+0075 | U+0323 | Close   | Unknown |   u   |
