# Challenges and Solutions in Working with Sango Text

## Glyph UTF representations are not fixed-width: **Segment input and use glyph representation internally**

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

### SOLUTION

All input is first converted to [NFKD](https://unicode.org/reports/tr15/#Norm_Forms)
for ease of handling accents separately, then reconverts to NFC on output.

Although the Go standard library (surprisingly) does not provide this functionality, there is
a third-party [Unicode Text Segmentation library](https://github.com/rivo/uniseg/blob/master/README.md) that does.

## Segmentation is hard: **Parse sequence of syllables instead**

### CHALLENGES

1. Non-literal translation of single words is difficult
2. Nonspace vs hyphen vs space segmentation of morphemes is not standardized and subject to significant variation.
3. A single English word might require a multiword phrase in Sango due to the latter's impoverished vocabulary, e.g.
   - **mafüta tî ngû tî mɛ tî bâgara** = "grease of water of teat of cow" = butter
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

## Unicode is inconvenient for humans: **Use ASCII encodings instead**

### CHALLENGES

UTF8 is a neverending source of bugs, results in complex code, requires the use of specialty libraries,
and makes it difficult to have random access to substrings or even determine string length.

1. It is tedious to input Sango text with accents and open vowels ɛ and ɔ from a keyboard.
2. It is hard for the visually challenged to distingish between circumflex and diaeresis accents.
3. The width of text hard to predict or assess when dealing with text layout.

### REQUIREMENTS

An ASCII-only representation is much easier for keyboard input, canonicalization, internal processing,
column alignment, and reading (especially aloud) in smaller font. The encoding should:

1. be one byte per glyph, for random access and predictable string length
2. consist only of lowercase letters, so that fingers rarely need to leave the home row of a keyboard
3. be easily stripped of vowel pitch and height through trivial transformations
4. be human readable, and facilitate correct pronunciation by nonnative speakers

### SOLUTION

Sango uses only 22 letters, leaving the other 4 letters (`x`, `c`, `j`, and `q`) (along with various punctuation) available for other uses.

In particular, the phonemic rigidity of Sango allows for a way to encode vowel pitch and height using uppercase:

|   Encodes   | ASCII                    |
| :---------: | :----------------------- |
|  Low pitch  | LC whole syllable        |
|   High ^    | UC whole syllable        |
|    Mid ¨    | LC consonants + UC vowel |
| Escape next | q                        |
| Raise pitch | j                        |
|      ɔ      | c                        |
|      ɛ      | x                        |
|      “      | ``                       |
|      ”      | ''                       |
|      ‘      | \`                       |
|      ’      | '                        |
|      «      | <<                       |
|      »      | >>                       |
|  Start UC   | <                        |
|   End UC    | >                        |
| Start Annot | {                        |
|  End Annot  | }                        |

Note that:

1. A `q` (or `Q`) escapes, i.e. disables the decoding for any single ASCII character after it
2. Pitch is preferably encoded via uppercase, but `j` (or `J`) is available as an alternative.
   Each `j` (cyclically) raises one level of pitch, e.g. `bja` = **bä**, `jja` = **bâ**, and `bjjja` = **ba**.
   - This is useful mainly to facilitate input from a keyboard, since it is much easier to type `j` than press and hold the Shift key.
   - Internally, the encoding is canonicalized whatever its original format.
3. Uppercase is used for pitch encoding and not available for capitalization.
   Instead, any encoded string within angle brackets (which may be nested) is uppercased.
4. Any string within braces (which may be nested) is syntactically ignored, and may
   be used for comments, semantic annotations, translations, or metadata.
5. Mid pitch syllables-initial vowels are not directly representable.
   Instead, high pitch syllable-initial vowels are implicitly lowered to mid pitch:
   - in a gerund (words ending in **-ngɔ̈** except for the word **îngɔ̈**)
   - for a closed known small fixed set of lexemes: `Apx` = **äpɛ**, not **âpɛ**.
6. Enclosing in single angular brackets is used to render uppercase,
   but unneed at the start of a sentence (or after final punctuation) where it is done automatically:
   - `balaɔ̂. ïrï tî mbï <kɔ̂lïngbâ> <k>ɔ̈sï.` = **Balaɔ̂. Ïrï tî mbï KƆ̂LÏNGBÂ Kɔ̈sï.**

## Code switching is frequent: **Classify using phonemics and lexicon**

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

In case of ambiguity:

- lexical lookup can be used as an allowlist of valid Sango lexemes
- language annotations enclosed in braces can be added
