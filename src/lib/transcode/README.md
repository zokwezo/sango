# Challenges and Solutions in Working with Sango Text

## Glyph UTF representations are not fixed-width: **Segment input and use syllable encoding internally**

### CHALLENGES

All Sango input is assumed to be in UTF8 format, but may not be in normal form, and lack of byte equality
can mess up lexicon lookup and be impossible to diagnose visually.

Even after normalizing input to [NFKD](https://unicode.org/reports/tr15/#Norm_Forms),
the number of bytes and runes varies with vowel, as seen in the following table:

| Graphemes | # Runes | # Bytes |
| :-------: | :-----: | :-----: |
|    e o    |    1    |    1    |
|    ɛ ɔ    |    1    |    2    |
|  ë ê ö ô  |    2    |    3    |
|  ɛ̈ ɛ̂ ɔ̈ ɔ̂  |    2    |    4    |

This means that graphemes in a Sango phrase cannot be randomly accessed, nor phrase length be
immediately calculated from its byte or rune representations when e.g. aligning interlinearly
with translations into columns.

Also, combining marks are separate runes that attach to their immediately preceding rune and can
be easily separated from them during segmentation or string reversal.

### SOLUTION

#### Input and Output
All input is first converted to [NFKD](https://unicode.org/reports/tr15/#Norm_Forms)
for ease of handling accents separately, then reconverts to NFC on output.

> _NOTE: Although the Go standard library (surprisingly) does not provide native functionality to iterate over whole glyphs (including combining marks), there is a third-party [Unicode Text Segmentation library](https://github.com/rivo/uniseg/blob/master/README.md) that does. However, its use increases complexity and where possible an ASCII encoding is used internally (see below)._

#### Internally

The data is stored internally encoded as sequences of syllables in `[]uint16`. This is also an optimal encoding for
use in machine learning algorithms. See the documentation on [Sango Syllabic Encoding](encoding.md) for full details.

## Segmentation is hard: **Parse sequence of syllables instead**

### CHALLENGES

1. Non-literal translation of single words is difficult
2. Nonspace vs hyphen vs space segmentation of morphemes is not standardized and subject to significant variation.
3. A single English word might require a multiword phrase in Sango due to the latter's impoverished vocabulary, e.g.
   - **mafüta tî ngû tî mɛ tî bâgara** = _oil of water of teat of cow_ = butter
4. A period (**.**) may indicate an abbreviation in English, but this is very uncommon in Sango.

### SOLUTION

1. When parsing, consider most likely sequence of syllables, not words
2. Prefer longest phrase found in lexicon
3. Use hyphens and word breaks only for ranking parses, not filtering them
   4 In Sango, a period can always be considered a sentence final marker.
   - Abbreviations such as _M_, _Mme_, or _Melle_ should be rendered without period.
4. In English translations, it is good practice to disambiguate abbreviation with a nonintrusive convention:
   - a period followed by a single space and a (possibly uppercase) letter always signals abbreviation
   - use two spaces to signal sentence final.

## Inputting Unicode from the keyboard is inconvenient: **Use an ASCII encoding**

### CHALLENGES

UTF8 is a neverending source of bugs, results in complex code, requires the use of specialty libraries,
and makes it difficult to have random access to substrings or even determine string length.

1. It is tedious to input Sango text with accents and open vowels ɛ and ɔ from a keyboard.
2. It is hard for the visually challenged to distingish between circumflex and diaeresis accents.
3. The width of text hard to predict or assess when dealing with text layout.

### REQUIREMENTS

An ASCII-only representation is much easier for keyboard input, canonicalization, internal processing,
column alignment, and reading (especially aloud) in smaller font. Desirable encoding properties include:

1. use only ASCII accessible on a keyboard, ideally only lowercase letters
   that can be typed without needing to leave the home row of a keyboard
2. be easily stripped of vowel pitch and/or height through trivial transformations
3. be one byte per glyph, for random access and predictable string length
4. be human readable, and facilitate correct pronunciation by nonnative speakers

### SOLUTION

#### Punctuation

It is convenient to use ASCII to encode punctuation common in Sango text that require UTF8 glyphs:

|     UTF8     | ASCII |
| :----------: | :---: |
|      “       |  ``   |
|      ”       |  ''   |
|      ‘       |  \`   |
|      ’       |   '   |
|      «       |  <<   |
|      »       |  >>   |
| ellipses (…) | . . . |
|  hyphen (-)  |   -   |
| en-dash (–)  |  - -  |
| em-dash (—)  | - - - |

> _NOTE: the spaces within the ASCII encodings are for display purposes only and not to be added in text._

#### Vowel height: SEXC encoding

> _NOTE: SEXC (pronounced "sexy-jay") is short for ‘Sango Encoding X and C’._

The consonants **x** and **c** are not used in Sango, and are repurposed to represent open vowels instead:

| UTF8 | ASCII |
| :--: | :---: |
|  ɛ   |   x   |
|  ɔ   |   c   |

#### Vowel pitch

There are two different schemes to encode vowel pitch, each with its advantages and disadvantages:

##### SEXC-J encoding

> _NOTE: SEXC-J (pronounced "sexy-jay") is short for ‘Sango Encoding X and C using J’._

The consonants **j** and **q** (both case-insensitive) do not occur in Sango and are repurposed to encode vowel pitch.
By default, vowels have low pitch. When immediately following a vowel, `j` indicates high pitch and `q` indicates
medium pitch.

* PRO: Trivial to convert bijectively to and from UTF8 with simple rune remapping.
* PRO: Case invariant
* CON: String length differs from glyph length (no random access)

##### SEXC-U encoding

> _NOTE: SEXC-U (pronounced "sexy-you") is short for ‘Sango Encoding X and C using Uppercase’._

Vowel pitch is important when humans speak or read Sango text, and humans are not good at transliterating `j`
on the fly, nor can easily distinguish visually the small circumflex and diaeresis accent marks.

Therefore, as an output format to faciliate human comprehension, uppercase is repurposed to encode vowel pitch:

1. Whole syllable is lowercase ⟹ Low pitch
   - **ngba** ⟹ `ngba`
2. Whole syllable is uppercase ⟹ High pitch
   - **ngbâ** ⟹ `NGBA`
3. Consonants in lowercase, Vowel in uppercase ⟹ Medium pitch
   - **ngbä** ⟹ `ngbA`

Note that this encoding is lossy:

1. Mid pitch vowel-only syllables (no consonants) are not directly representable.
   Instead, high pitch syllable-initial vowels are implicitly lowered to mid pitch:
   - in a gerund (words ending in **-ngɔ̈** except for the word **îngɔ̈**)
   - for a closed fixed set of lexemes: `Apx` = **äpɛ**, not **âpɛ**.
2. Case is NOT preserved in this encoding (having been repurposed for pitch),
   and not usually needed since this encoding is not intended for subsequent processing.
	 If needed to preserve a lossless conversion, case might be preserved as metadata.

#### Metadata

Syntactically, metadata

1. starts with an opening brace
2. ends with a matching closing brace
3. may contain metadata of its own, whose braces must correctly nest.
4. all text not within braces represents literal syntax (in an encoding defined by its current context).

The interpretation of metadata may be exogenous or context-dependent, and not further discussed here.
Code which cannot interpret metadata must ignore it completely.

Specific use cases include:

- encoding a UTF8 literal (e.g. the glyph **x** within a SEXC encoding where `x` represents **ɛ**)
- case (e.g. that the following word is in uppercase, in SEXC-U encoding where uppercase represents a high pitch tone)
- segmentation into tokens
- semantic annotations
- set or modify the current context
- language tags, to mark code switching
- parsing/control directives
- translations
- comments

## Code switching is frequent: **Classify using phonemics and lexicon, record as metadata**

### CHALLENGES

French is commonly injected into Sango, and increases with the competence of the speaker in French. This occurs because:

- Sango is an impoverished language, and may not have a suitable replacement
- The speaker is more fluent in French than Sango
- As a signaling mechanism of erudition
- In urban environments, where code switching occurs frequenty in French
  - In villages, code switching with the local tribal language is more common
  - However, rural Sango is rarely found in written form anyway, so this is less of a problem

### REQUIREMENTS

Whereas spoken Sango often uses French loan words after retrofitting them into Sango phonemes,
written Sango prefers to leave the French word in its original French orthography.
Consequently, language parsing needs to recognize these, persist but otherwise
ignore them during parsing, then reproduce them during generation.

### SOLUTION

Phonemic solutions are fastest and preferred:

- If a syllable is not Sango, any other syllables juxtaposed or connected by hyphen are also not Sango.
- Non-initial capitalized words are always either proper nouns or loan words, not Sango lexemes.
- Sango has a rigid _C?V_ phonology (see [phonology.md](phonology.md) for details).
  Any nonconforming spelling should be considered non-Sango.
  - The acute and grave accents are never found in Sango
    - NOTE: circumflex and diaeresis accents occur over all vowels in both French and Sango, so their presence is not dispositive.
  - The letters **c**, **j**, **q**, and **x** and the diphthongs **ei**, **ie**, and **ou** are never found in Sango.
  - The letters **ɛ** and **ɔ** are not found in French (nor in the standard Sango orthography) but
    can be assumed to be open **e** and **o** in Sango.
- In case of ambiguity:
  - lexical lookup can be used as an allowlist of valid Sango lexemes
  - language annotations enclosed in braces can be added
