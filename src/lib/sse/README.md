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
| SSEs      | [0x65E5672C8A9E306F, 0x000096E330573044, 0x0021002000A70000, 0xF062BE5451320FFF] |
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

| B61 \\ B60 |     0     |       1       |
| :--------: | :-------: | :-----------: |
|      0     | lower ()  | Title (~)     |
|      1     | UPPER (=) | Invisible (#) |

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
* R immediately follows L
* TODO: Nasal consonants immediately follow their non-nasal counterpart
* Nasal vowels immediately follow their non-nasal counterpart
* Open vowels immediately precede their closed counterpart

###### Canonical

Consonants are ordered so that natural allophones are adjacent:

| b10-b8 \\ b7-b6 | 00  | 01  | 10  | 11  |
| :-------------: | :-: | :-: | :-: | :-: |
|        000      |  h  |  H  |  b  |  d  |
|        001      |  f  |  g  |  q  |  k  |
|        010      |  K  |  l  |  r  |  m  |
|        011      |  B  |  P  |  V  |  n  |
|        100      |  D  |  G  |  Q  |  Y  |
|        101      |  Z  |  p  |  s  |  t  |
|        110      |  v  |  w  |  y  |  z  |

###### UTF8

| b10-b8 \\ b7-b6 | 00  | 01  | 10  | 11  |
| :-------------: | :-: | :-: | :-: | :-: |
|        000      |     |  h  |  b  |  d  |
|        001      |  f  |  g  | gb  |  k  |
|        010      | kp  |  l  |  r  |  m  |
|        011      | mb  | mp  | mv  |  n  |
|        100      | nd  | ng  | ngb | ny  |
|        101      | nz  |  p  |  s  |  t  |
|        110      |  v  |  w  |  y  |  z  |

<!-- There is an alternate arrangement (not used herein) of consonant clusters that
maximally encodes phonetic symmetries (voiced, nasal, aspirated, yotated) in
bit patterns, subject to having a dense compact encoding scheme, as shown below.

###### Canonical

| b10-b8 \\ b7-b6 | 00  | 01  | 10  | 11  |
| :-------------: | :-: | :-: | :-: | :-: |
|        000      |  h  |  H  |  b  |  v  |
|        001      |  y  |  d  |  z  |  q  |
|        010      |  g  |  p  |  f  |  l  |
|        011      |  r  |  t  |  s  |  K  |
|        100      |  k  |  w  |  B  |  V  |
|        101      |  Y  |  D  |  Z  |  Q  |
|        110      |  G  |  n  |  P  |  m  |

###### UTF8

| b10-b8 \\ b7-b6 | 00  | 01  | 10  | 11  |
| :-------------: | :-: | :-: | :-: | :-: |
|        000      |     |  b  |  v  |  y  |
|        001      |  d  |  z  |  gb |  g  |
|        010      |  h  |  p  |  f  |  l  |
|        011      |  t  |  s  |  kp |  k  |
|        100      |  w  |  mb |  mv |  ny |
|        101      |  nd |  nz | ngb |  ng |
|        110      |  n  |  mp |  m  |  r  |
-->

Syllables with consonant codes starting with **111** are ignored entirely.

##### Vowel

###### Canonical

| b5-b4 \\ b3-b2 | 00  | 01  | 10  | 11  |
| :------------: | :-: | :-: | :-: | :-: |
|       00       |  a  |  A  |  X  |  x  |
|       01       |  e  |  E  |  i  |  I  |
|       10       |  C  |  c  |  o  |  O  |
|       11       |  u  |  U  |  —— |  —— |

###### UTF8

| b5-b4 \\ b3-b2 | 00  | 01  | 10  | 11  |
| :------------: | :-: | :-: | :-: | :-: |
|       00       |  a  |  añ |  ∉ |  ɛ  |
|       01       |  e  |  eñ |  i  |  iñ |
|       10       |  ∅ |  ɔ  |  o  |  oñ |
|       11       |  u  |  uñ |  —— |  —— |

<!-- There is also an alternate arrangement (not used herein) of vowels that
maximally encodes phonetic symmetries (height, nasal) in bit patterns,
subject to having a dense compact encoding scheme, as shown below.

###### Canonical

| b5-b4 \\ b3-b2 | 00  | 01  | 10  | 11  |
| :------------: | :-: | :-: | :-: | :-: |
|       00       |  a  |  A  |  e  |  E  |
|       01       |  i  |  I  |  o  |  O  |
|       10       |  x  |  c  |  u  |  U  |
|       11       |  X  |  C  |  —— |  —— |

###### UTF8

| b5-b4 \\ b3-b2 | 00  | 01  | 10  | 11  |
| :------------: | :-: | :-: | :-: | :-: |
|       00       |  a  |  añ |  e  |  eñ |
|       01       |  i  |  iñ |  o  |  oñ |
|       10       |  ɛ  |  ɔ  |  u  |  uñ |
|       11       |  E  |  O  |  —— |  —— |
-->

Syllables with vowel codes starting with **11** are ignored entirely.

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

| b1 \\ b0 |   0   |   1     |
| :------: | :---: | :-----: |
|     0    | Low   | Mid     |
|     1    | High  | UNKNOWN |

#### SANGO EXAMPLES

##### lowercase

| Format    | Value                   |
| --------- | ----------------------- |
| SSE       | 0x8_08E_9E5_319_5CC_FFF |
| Canonical | "bx^-kc:Bi:tx_"         |
| UTF8      | "bɛ̂-kɔ̈mbïtɛ"            |


##### UPPERCASE with space prefix

| Format    | Value                   |
| --------- | ----------------------- |
| SSEs      | 0xE_08E_9E5_319_5CC_FFF |
| Canonical | " =bx^-=kc:=Bi:=tx_"    |
| UTF8      | " BƐ̂-KƆ̈MBÏTƐ"           |
