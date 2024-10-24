# Sango Syllabic Encoding (SSE)

[Sango phonology](./phonology.csv) is a rigid C?V format and can be efficiently encoded as
`uint16` tokens are used for simpler coding and manipulation, making it easy to:

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

| Binary bit pattern | Description                      |
| ------------------ | -------------------------------- |
| `00UUUUUUUUUUUUUU` | Unicode rune                     |
| `01LNNNNNAAAAAAAA` | ASCII letter (English or French) |
| `1SSXXCCCCCVVVVPP` | Syllable (Sango)                 |

where the bit substrings are fixed-length binary numerals when masked and shifted:

| Bit | Description                                   |
| :-: | --------------------------------------------- |
| `U` | Unicode rune (U+0000 - U+3FFF)                |
| `L` | Language (0=English, 1=French)                |
| `N` | `min(31,n)` where `n` = # letters remaining   |
| `A` | ASCII character (U+00 - U+FF)                 |
| `S` | `min(3,m)` where `m` = # syllables remaining  |
| `X` | Case                                          |
| `C` | Consonant cluster                             |
| `V` | Vowel                                         |
| `P` | Pitch                                         |

The Sango syllable encoding is defined as follows:

### Case

| MSB\\LSB |        0        |     1     |
| :------: | :-------------- | :-------: |
|    0     | lowercase       | Titlecase |
|    1     | hyphen-prefixed | UPPERCASE |

### Pitch

| MSB\\LSB |       0      |       1       |
| :------: | :----------: | :-----------: |
|    0     | Unknown  (ọ) | Low  tone (o) |
|    1     | Mid tone (ö) | High tone (ô) |

### Consonant cluster

| MSB \\ LSB | 00           | 01  | 10   | 11     |
| :--------: | ------------ | --- | ---- | ------ |
|    000     | *missing*    | f   | r    | k      |
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
* These are stand-in vowels when the vowel height is unknown,
  to be replaced by the appropriate open or close vowel once known:
  - **ə** ⟹ **e** or **ɛ**
  - **ø** ⟹ **o** or **ɔ**

## Examples

| Text   | Tokens                                                 |
| ------ | ------------------------------------------------------ |
| Hi     | `[0b_01_0_00001_01001000,   0b_01_0_00000_01101001]`   |
| Bɛ̂-bïn | `[0b_1_01_01_01101_0011_11, 0b_1_00_10_01101_1101_10]` |
| bə̣bị   | `[0b_1_01_11_01101_1011_00, 0b_1_00_11_01101_0101_00]` |
