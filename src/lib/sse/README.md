# Sango Syllabic Encoding (SSE)

Sango phonology is a rigid CVP (Consonant + Vowel + Pitch) syllabic format,
and each Sango word can fit inside a single `uint64` integer,
where it is much easier to manipulate and store internally. Decoding into ASCII
is done only on output. This allows the software to easily:

- Compactify notation with low entropy: suitable as vector embedding in machine
  learning algorithms
- Convert into and out of UTF8 and avoid invalid Sango phonemes
- Iterate by Unicode glyphs, Sango glyphs, or entire Sango words without
  worrying about byte boundaries
- Code switch between Sango, human languages, and symbols in the same document
- Query different properties and mask or filter on unimportant ones
- Record inline metadata by setting Case to Hidden
- Isolate use of hyphens and spaces, which in Sango are neither syntactically
  standardized nor semantically important

## Word type encoded

The 64-bit encoding divides up into components, from Most Significant Bit
(MSB = B63) to Least (LSB = B0) as follows:

| B63 | Word type         |
| --- | ----------------- |
| `0` | Unicode substring |
| `1` | Sango word        |

### Type 0: Unicode substrings

For a Unicode substring, the 4 bytes comprise up to 4 Unicode code points
(U+0001 through U+FFFF, ignoring any U+0000 bytes). The first byte must be
less than U+8000, otherwise set it to zero and continue with the next byte.

#### UNICODE EXAMPLE

| Format      | Value                                                               |
| ----------- | ------------------------------------------------------------------- |
| Unicode     | `日本語は難しい!`                                                     |
| UTF8        | `0xE697A5_E69CAC_E8AA9E_E381AF_E99BA3_E38197_E38184_21`             |
| Code points | `U+65E5 U+672C U+8A9E u+306F U+96E3 U+3057 U+3044 U+0021`           |
| SSEs        | `0x65E5_672C_8A9E_306F 0x0000_96E3_3057_3044 0x0021_0000_0000_0000` |

### Type 1: Sango words

For a Sango word, the rest of the bits encode a single Sango word of up to
5 syllables (there are no Sango words of more than 5 syllables).
Words are implicitly and by default separated by spaces, but a compound
word may be connected by hyphens by specifying a hyphen word separator.

The 64-bit encoding divides up into 4 global bits (1 hex digit) and up to
five groups of 12 bits (3 hex digits) followed by all ONES in unused LSB.

#### Word separator

| B62 | Prefix |
| --- | ------ |
| `0` | None   |
| `1` | Space  |

#### Case

The first 4 bits `1PCC` are global to the word encode the `P`refix and `C`ase:

| B61 \\ B60 |       `0`       |       `1`       |
| :--------: | :-------------: | :-------------: |
|     `0`    | lowercase ()    | Titlecase (`~`) |
|     `1`    | UPPERCASE (`$`) | Invisible       |

The case is applied to UTF8 directly, and the value in parentheses is
prepended to canonical format. In both cases, invisible case is unrendered.

#### Sango syllable(s)

The next 60 bits comprise 5 groups (one per syllable) of 12 bits `CCCCCVVVVPP`
(designated b11 = MSB, b0 = LSB), with unused syllables indicated by setting all ONES.
No two Sango words are placed in one `uint64`.

Each Sango syllable is phonemically strictly CVP (one consonant cluster +
one vowel + a pitch tone), although the standard orthography obscures this.
This encoding recognizes that `L` and `R` are mostly allophonic and might be
merged internally as `l`, though when phonemically well established a separate
`r` is available.

The encoding also represents a missing consonant as an unaspirated `h`
(written here in lowercase) as most Sango `h` letters are nearly silent anyway,
and there is no well accepted convention on when to write or omit
initial unaspirated `h`. When semantically important, aspirated `H`
(almost exclusively, word initial) is written herein in uppercase.

##### Syllable separator

Hyphens internal to a word are often (but not always) used in compound
words and this usage is not entirely standardized in the lexicon.
Where important, a compound word can be specify a hyphen syllable prefix
to separate word constituents.

| b11 | Prefix |
| --- | ------ |
| `0` | None   |
| `1` | Hyphen |

##### Consonant cluster

The arrangement of consonant clusters is meant to maximally encode phonetic
symmetries (voiced, nasal, aspirated, yotated) in bit patterns, subject to
having a dense compact encoding scheme.

###### Canonical

| b10-b8 \\ b7-b6 | `00` | `01` | `10` | `11` |
| :-------------: | :--: | :--: | :--: | :--: |
|       `000`     |  h   |  b   |  v   |  y   |
|       `001`     |  d   |  z   |  q   |  g   |
|       `010`     |  H   |  p   |  f   |  l   |
|       `011`     |  t   |  s   |  K   |  k   |
|       `100`     |  w   |  B   |  V   |  Y   |
|       `101`     |  N   |  Z   |  Q   |  G   |
|       `110`     |  n   |  P   |  m   |  r   |

###### UTF8

| b10-b8 \\ b7-b6 | `00` | `01` | `10` | `11` |
| :-------------: | :--: | :--: | :--: | :--: |
|       `000`     |      |  b   |  v   |  y   |
|       `001`     |  d   |  z   |  gb  |  g   |
|       `010`     |  h   |  p   |  f   |  l   |
|       `011`     |  t   |  s   |  kp  |  k   |
|       `100`     |  w   |  mb  |  mv  |  ny  |
|       `101`     |  nd  |  nz  | ngb  |  ng  |
|       `110`     |  n   |  mp  |  m   |  r   |

Syllables starting with `111` are ignored entirely, irrespective of the
value of the remaining 9 bits.

##### Vowel

The arrangement of vowel clusters is also meant to maximally encode phonetic
symmetries (height, nasal) in bit patterns, subject to having a dense compact
encoding scheme.

###### Canonical

| b5-b4 \\ b3-b2 | `00` | `01` | `10` | `11` |
| :------------: | :--: | :--: | :--: | :--: |
|      `00`      |  a   |  A   |  e   |  E   |
|      `01`      |  i   |  I   |  o   |  O   |
|      `10`      |  x   |  c   |  u   |  U   |
|      `11`      |  X   |  C   |  ——  |  ——  |

###### UTF8

| b5-b4 \\ b3-b2 | `00` | `01` | `10` | `11` |
| :------------: | :--: | :--: | :--: | :--: |
|      `00`      |  a   |  añ  |  e   |  eñ  |
|      `01`      |  i   |  iñ  |  o   |  oñ  |
|      `10`      |  ɛ   |  ɔ   |  u   |  uñ  |
|      `11`      |  E   |  O   |  ——  |  ——  |

Syllables whose vowels start with `111` are ignored entirely, irrespective of
the value of the other bits.

Externally, nasal vowels are followed by an `n` with no tilde. To resolve
ambiguity with a following syllable starting with an `n` or omitted unaspirated
`h`, an apostrophy or hyphen is used to separate the syllables.

Internally, `E` and `O` respectively represent an e/ɛ or o/ɔ vowel of
unknown height. The letters `x` and `c` are not used in Sango and
therefore used internally to represent the Unicode glyphs `ɛ` and `ɔ`.

In the standard orthography, vowel height (which is a meaningful distinction
only for `E` and `O`) is not expressed in writing because it is readily restored
orally by native speakers, but it is phonemically stable and semantically
important and therefore encoded internally when known, and restored where
possible when not known.

##### Pitch

All Sango vowels have exactly one of three pitch tones (Low, Mid, High) or else
unknown pitch. These are represented in the standard orthography by diacritics, respectively
none (`o`), dieresis (`ö`), or circumflex (`ô`). This project nonstandardly represents
unknown pitch with a dot below (`ọ`). Internally, these are represented by
vowel suffixes `_`, `:`, `^`, and `=` respectively, for ease of typing
and use in code.

| b1 \\ b0 |   `0`   |   `1`   |
| :------: | :-----: | :-----: |
|    `0`   |   Low   | Mid     |
|    `1`   |   High  | UNKNOWN |

#### SANGO EXAMPLE

| Format    | Value                                                                   |
| :-------: | ----------------------------------------------------------------------- |
| Sango     | `Bɛ̂-kɔ̈mbïtɛ`                                                            |
| Canonical | `bx^-kc:Bi:tx_`                                                         |
| Binary    | `1001_000001100010_101111100101_010001010001_001100100000_111111111111` |
| Hex       | `9_062_BE5_451_320_FFF`                                                 |
| uint64    | `10404087358528032767`                                                  |
