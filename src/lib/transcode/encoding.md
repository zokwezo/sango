# Sango Syllabic Encoding

Internally, tokens are used for simpler coding and manipulation, making it easy to:

- Compactify notation with low entropy: suitable as vector embedding in machine learning algorithms
- Easy to convert into and out of UTF8, validate syllables, and avoid invalid Sango phonemes
- Iterate by symbol, letter (English/French), syllable (Sango), and word (both) without worrying about byte boundaries
- Distinguish language and punctuation/whitespace just by inspecting the the high-order bits
- Use interlinear code switching
- Query different properties and mask or filter on unimportant ones
- Record inline metadata by setting Case to Hidden
- Isolate use of a hyphen, which in Sango is neither syntactically standardized nor semantically important

## Encoding format

The 16-bit encoding divides up quasi-orthogonally into components:

| Binary bit pattern | Description                     |
| ------------------ | ------------------------------- |
| `00UUUUUUUUUUUUUU` | Unicode rune                    |
| `010LLLLLAAAAAAAA` | ASCII letter in an English word |
| `011LLLLLAAAAAAAA` | ASCII letter in a French word   |
| `1SSXXPPCCCCCVVVV` | Syllable in a Sango word        |

where the bit substrings are fixed-length binary numerals that encode various components.
Each code is trivially convertible to a unicode rune or (after masking) ASCII character
if the most-significant-bit (MSB) is `0`.

| Bit | Description                                   |
| :-: | --------------------------------------------- |
| `U` | Unicode rune (U+0000 - U+3FFF)                |
| `A` | ASCII character (U+00 - U+FF)                 |
| `L` | `min(31,n)` where `n` = # letters remaining   |
| `S` | `min( 3,m)` where `m` = # syllables remaining |
| `X` | Case                                          |
| `P` | Pitch                                         |
| `C` | Consonant cluster                             |
| `V` | Vowel                                         |

### Case

| MSB\\LSB |        0        |     1     |
| :------: | :-------------: | :-------: |
|    0     |     Hidden      | lowercase |
|    1     | hyphen-prefixed | Titlecase |

### Pitch

| MSB\\LSB |    0     |     1     |
| :------: | :------: | :-------: |
|    0     | Unknown  | Low tone  |
|    1     | Mid tone | High tone |

### Consonant cluster

| MSB \\ LSB | 00           | 01           | 10           | 11           |
| :--------: | ------------ | ------------ | ------------ | ------------ |
|    000     | ∅            | f            | r            | k            |
|    001     | mv           | v            | ng           | g            |
|    010     | m            | p            | l            | kp           |
|    011     | mb           | b            | ngb          | gb           |
|    100     | **reserved** | s            | y            | h            |
|    101     | nz           | z            | ny           | w            |
|    110     | n            | t            | nd           | d            |
|    111     | **reserved** | **reserved** | **reserved** | **reserved** |

### Vowel

| MSB\\LSB | 00           | 01  | 10  | 11  |
| :------: | ------------ | --- | --- | --- |
|    00    | **reserved** | u   | ɔ   | ɛ   |
|    01    | a            | i   | o   | e   |
|    10    | **reserved** | uñ  | ɔ/o | ɛ/e |
|    11    | añ           | iñ  | oñ  | eñ  |

## Examples

| Text      | Tokens                                                                     |
| --------- | -------------------------------------------------------------------------- |
| Hello!    | `[0x4448, 0x4365, 0x426c, 0x416c, 0x406f, 0x0021]`                         |
| c'est ça… | `[0x6463, 0x6327, 0x6265, 0x6173, 0x6074, 0x0020, 0x61e7, 0x6061, 0x2026]` |
| Bɛ̂-bï     | `[0xa4b3, 0xa4d5]` _(visible, known vowel pitch/height)_                   |
| _bebi_    | `[0xa0db, 0x80d5]` _(hidden, unknown vowel pitch/height)_                  |
