# Sango Syllabic Encoding (SSE)

Sango phonology is a rigid CVP (Consonant + Vowel + Pitch) syllabic format,
and each Sango word can fit inside a single **uint64** integer,
where it is much easier to manipulate and store internally. Decoding into ASCII
is done only on output. This allows the software to easily:

- Compactify notation with low entropy: suitable as vector embedding in machine
  learning algorithms
- Convert into and out of UTF8 and avoid invalid Sango phonemes
- Iterate by Unicode glyphs, Sango glyphs, or entire Sango words without
  worrying about byte boundaries
- Code switch between Sango, human languages, and symbols in the same document
- Query different properties and mask or filter on unimportant ones
- Record inline metadata by setting Shift to Hidden
- Isolate use of hyphens and spaces, which in Sango are neither syntactically
  standardized nor semantically important

## Kind of encoding

The 64-bit encoding divides up into components, from Most Significant Bit
(MSB = B63) to Least (LSB = B0) as follows:

| B63 | Contents encoded  |
| :-: | ----------------- |
|  0  | Unicode substring |
|  1  | Sango word        |

### Kind 0: Unicode substrings

For a Unicode substring, the 4 bytes comprise up to 4 Unicode code points
(U+0001 through U+FFFF, ignoring any U+0000 bytes). The first byte must be
less than U+8000, otherwise set it to zero and continue with the next byte.

#### UNICODE EXAMPLE

| Format    | Value                                                                            |
| --------- | -------------------------------------------------------------------------------- |
| SSEs      | [0x65E5672C8A9E306F, 0x000096E330573044, 0x0021002000A70000, 0xF062BE5451320000] |
| Canonical | "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7"                   |
| UTF8      | "日本語は難しい! §"                                                                |

### Kind 1: Sango words

For a Sango word, the rest of the bits encode a single Sango word of up to
5 syllables (there are no Sango words of more than 5 syllables).
Words are implicitly and by default separated by spaces, but a compound
word may be connected by hyphens by specifying a hyphen word separator.

The 64-bit encoding divides up into 4 global bits (1 hex digit) and up to
five groups of 12 bits (3 hex digits) followed by all ONES in unused LSB.

The first 4 bits **1PCC** are global to the word encode the `P`refix and `C`ase:

#### Word separator

| B62 | Prefix |
| :-: | ------ |
|  0  | None   |
|  1  | Space  |

#### Shift

| B61 \\ B60 |     0         |     1     |
| :--------: | :-----------: | :-------: |
|      0     | Invisible (#) | lower ()  |
|      1     | Title (~)     | UPPER (=) |

The shift is applied to UTF8 directly, and the value in parentheses is
prepended to canonical format. In both cases, invisible shift is unrendered.

#### Sango syllable(s)

The next 60 bits comprise 5 groups (one per syllable) of 12 bits **CCCCCVVVVPP**
(where b0 = LSB), with unused syllables indicated by setting all ONES.
No two Sango words are placed in one **uint64**.

Each Sango syllable is phonemically strictly CVP (one consonant cluster +
one vowel + a pitch tone), although the standard orthography obscures this.
The consonants **L** and **R** have their own encoding although they are mostly
allophonic in Sango and might be merged internally as **l**.

The encoding also represents a missing consonant as an unaspirated **h**
(written here in lowercase) as most Sango **h** letters are nearly silent anyway,
and there is no well accepted convention on when to write or omit
initial unaspirated **h**. When semantically important, aspirated **H**
(almost exclusively, word initial) is written herein in uppercase.

##### Syllable separator

Hyphens internal to a word are often (but not always) used in compound
words and this usage is not entirely standardized in the lexicon.
Where important, a compound word can be specify a hyphen syllable infix
to separate word constituents.

| b11 | Infix  |
| :-: | ------ |
|  0  | None   |
|  1  | Hyphen |

##### Consonant cluster

Consonants are arranged in quasi-lexical CV sorting order, except that
allophones and near-allophones are grouped:
* unaspirated and aspirated H immediately follow no consonant
* Nasal consonants and vowels immediately follow their non-nasal counterpart
* Open vowels immediately precede their closed counterpart

###### Canonical

Consonants are ordered so that natural allophones and
phonetically similar phonemes are adjacent:

| b10-b8 \\ b7-b6 | 00  | 01  | 10  | 11  |
| :-------------: | :-: | :-: | :-: | :-: |
|        000      |  —  |  —  |  h  |  b  |
|        001      |  d  |  f  |  g  |  q  |
|        010      |  H  |  k  |  K  |  l  |
|        011      |  m  |  B  |  P  |  V  |
|        100      |  n  |  D  |  G  |  Q  |
|        101      |  Y  |  Z  |  p  |  r  |
|        110      |  s  |  t  |  v  |  w  |
|        111      |  y  |  z  |  —  |  —  |

###### UTF8

| b10-b8 \\ b7-b6 | 00  | 01  | 10  | 11  |
| :-------------: | :-: | :-: | :-: | :-: |
|        000      |  —  |  —  |     |  b  |
|        001      |  d  |  f  |  g  | gb  |
|        010      |  h  |  k  | kp  |  l  |
|        011      |  m  | mb  | mp  | mv  |
|        100      |  n  | nd  | ng  | ngb |
|        101      | ny  | nz  |  p  |  r  |
|        110      |  s  |  t  |  v  |  w  |
|        111      |  y  |  z  |  —  |  —  |

Syllables with consonant codes marked by a — are ignored entirely.
Unaspirated h (`0b00010`) is omitted entirely when outputting UTF8.

##### Vowel

###### Canonical

| b5-b4 \\ b3-b2 | 00  | 01  | 10  | 11  |
| :------------: | :-: | :-: | :-: | :-: |
|       00       |  —  |  —  |  a  |  A  |
|       01       |  X  |  x  |  e  |  E  |
|       10       |  i  |  I  |  C  |  c  |
|       11       |  o  |  O  |  u  |  U  |

###### UTF8

| b5-b4 \\ b3-b2 | 00  | 01  | 10  | 11  |
| :------------: | :-: | :-: | :-: | :-: |
|       00       |  —  |  —  |  a  |  añ |
|       01       |  x  |  ɛ  |  e  |  eñ |
|       10       |  i  |  iñ |  c  |  ɔ  |
|       11       |  o  |  oñ |  u  |  uñ |

Syllables with vowel codes marked by a — are ignored entirely.

Externally, nasal vowels are followed by an **n** with no tilde. To resolve
ambiguity with a following syllable starting with an **n** or omitted unaspirated
**h**, an apostrophy or hyphen is used to separate the syllables.

Internally, **E** and **O** respectively represent an **e**/**ɛ** or **o**/**ɔ** vowel of
unknown height. The letters **x** and **c** are not used in Sango and
therefore used internally to represent the Unicode glyphs **ɛ** and **ɔ**.

In the standard orthography, vowel height (which is a meaningful distinction
only for **E** and **O**) is not expressed in writing because it is readily restored
orally by native speakers, but it is phonemically stable and semantically
important and therefore encoded internally when known, and restored where
possible when not known.

##### Pitch

All Sango vowels have exactly one of three pitch tones (Low, Mid, High) or else
unknown pitch. These are represented in the standard orthography by diacritics, respectively
none (**o**), dieresis (**ö**), or circumflex (**ô**). This project nonstandardly represents
unknown pitch with a dot below (**ọ**). Internally, these are represented by
vowel suffixes **_**, **:**, **^**, and **=** respectively, for ease of typing
and use in code.

| b1 \\ b0 |    0    |   1   |
| :------: | :-----: | :---: |
|     0    | UNKNOWN | Low   |
|     1    | Mid     | High  |

#### SANGO EXAMPLES

##### lowercase

| Format    | Value                   |
| --------- | ----------------------- |
| SSE       | 0x9_0D7_A6E_362_655_000 |
| Canonical | "bx^-kc:Bi:tx_"         |
| UTF8      | "bɛ̂-kɔ̈mbïtɛ"            |


##### UPPERCASE with space prefix

| Format    | Value                   |
| --------- | ----------------------- |
| SSEs      | 0xF_0D7_A6E_362_655_000 |
| Canonical | " =bx^-=kc:=Bi:=tx_"    |
| UTF8      | " BƐ̂-KƆ̈MBÏTƐ"           |
