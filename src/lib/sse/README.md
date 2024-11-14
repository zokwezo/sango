# Sango Syllabic Encoding (SSE)

[Sango phonology](./phonology.md) is a rigid C?V format and can be efficiently encoded as
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

| MSB\\LSB | 00  | 01  | 10  | 11  |
| :------: | --- | --- | --- | --- |
|   000    |     | b   | d   | f   |
|   001    | g   | gb  | h   | k   |
|   010    | kp  | l   | m   | mb  |
|   011    | mp  | mv  | n   | nd  |
|   100    | ng  | ngb | ny  | nz  |
|   101    | p   | r   | s   | t   |
|   110    | v   | w   | y   | z   | 

### Vowel

| MSB\\LSB | 00 | 01 | 10 | 11 |
| :------: | -- | -- | -- | -- |
|    00    |    | a  | añ | ə  |
|    01    | ɛ  | e  | eñ | i  |
|    10    | iñ | ø  | ɔ  | o  |
|    11    | oñ | u  | uñ | —— |

* The following stand-in vowels are not found in normal Sango text and are used internally to indicate that
  the vowel height is unknown and is to be replaced by the appropriate open or close vowel once known:
  - **ə** ⟹ **e** or **ɛ**
  - **ø** ⟹ **o** or **ɔ**
