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

| MSB \\ LSB | 00           | 01  | 10   | 11     |
| :--------: | ------------ | --- | ---- | ------ |
|    000     | missing      | f   | r    | k      |
|    001     | mv           | v   | ng   | g      |
|    010     | m            | p   | l    | kp     |
|    011     | mb           | b   | ngb  | gb     |
|    100     | **invalid**  | s   | y    | h      |
|    101     | nz           | z   | ny   | w      |
|    110     | n            | t   | nd   | d      |

### Vowel

| MSB \\ LSB | 00           | 01  |  10   |  11   |
| :--------: | ------------ | --- | ----- | ----- |
|     00     | **missing**  | u   |   ɔ   |   ɛ   |
|     01     | a            | i   |   o   |   e   |
|     10     | **invalid**  | uñ  | **ø** | **ə** |
|     11     | añ           | iñ  |  oñ   |  eñ   |

* Bold entries are not found in normal Sango text.
* **ə** is a stand-in for either **e** or **ɛ** when vowel height is unknown. On output, all three should be replace by **e**.
* **ø** is a stand-in for either **o** or **ɔ** when vowel height is unknown. On output, all three should be replace by **o**.


## Examples

| Text      | Tokens                                                     |
| --------- | ---------------------------------------------------------- |
| "Hello"   | `[0x4548, 0x4465, 0x436c, 0x426c, 0x416f]` _ASCII English_ |
| "Bɛ̂-bïn"  | `[0xbed3, 0x94dd]` _(visible, known vowel pitch/height)_   |
| "bebi"    | `[0xa8db, 0x88d5]` _(hidden, unknown vowel pitch/height)_  |
