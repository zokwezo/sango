# Sango phonology

Sango words consist of one or more syllables, with the following lexicon histogram:

| # syllables | # lemmas |
|:-----------:|---------:|
|      1      |    254   |
|      2      |    903   |
|      3      |    378   |
|      4      |    100   |
|      5      |      5   |
|      6      |      1   |

Each syllable consists of a consonant cluster followed by a vowel. Hyphens, spaces, or fusion may connect
closely related morphemes, prefixes, suffixes, and compound words, and usage is not standardized.

## Consonants

The consonants comprise the closed set

_{*ø*, **b, d, f, g, gb, h, k, kp, l, m, mb, mv, n, nd, ng, ngb, ny, nz, p, r, s, t, v, w, y, z**}_

where *ø* indicates a zero consonant, i.e. a syllable with only a vowel.

> NOTE: Including a zero consonant is purely a housekeeping device to simplify tokenization (to yield a strict CV syllable encoding) which simplifies many algorithms, as when enumerating the set of all possible syllables. This convention in no way implies any underlying historic linguistic mechanism, in the way that silent letters may imply in other languages.

Words imported from other languages are remapped into the closest Sango phoneme, e.g.:
* Afrique (Africa)  ⟹ **Afirîka**
* Allemagne (Germany)  ⟹ **Zalamäa**
* [Mpoko](https://en.wikipedia.org/wiki/Mpoko_River) ⟹ **Pökö**
* Tchad (Chad)  ⟹ **Sâde**

> Sometimes foreign words (especially French!) are inserted untranslated or untransliterated, especially in cities or among professionals, if the corresponding Sango word does not exist, is unknown to the speaker, or to signal erudition. This is treated internally as code switching (i.e. mostly ignored) and not as Sango text.

## Vowels

The vowels comprise the closed set

__{**a, an, e, en, ɛ, i, in, o, on, ɔ, u, un**}__

Every syllable has exactly one vowel. When a zero consonant occurs, two adjacent vowels are left adjoining, and
these are not diphthongized but retain their full individual pitch, height, and length.
There is a clear audible distinction between the words **ba bâ bâa**, **lâ laâ**, **da dä daä**.

The following _VV_ combinations are **not** found in Sango: **ei, ɛ**_V_**, ie, iɛ, ou, ɔu, uo, uɔ**

### Vowel height

Spoken Sango strictly distinguishes open (**ɛ** and **ɔ**) and close (**e** and **o**) middle vowels.

Written Sango usually does not note open vowels separately (only because there are no keys on a typewriter for them).

> Since technology has greatly improved since the orthography was standardized in 1984 and many applications
> (including text-to-speech and language learning) as well as human users would benefit from making this distinction,
> it is strongly recommended to maintain this distinction in any intermediate tools, discarding it only on final output if desired.

As in Hungarian or Finnish, Sango enforces vowel harmony, where these open and close vowels do not coexist in the same word root.
- Out of 450 roots in the lexicon with at least two middle vowels, there are only four exceptions violating vowel harmony:
  - **lêngbêtɔ̂rɔ̂** (soybean)
  - **mɔkondö** (holy, sainthood)
  - **môlɛngɛ̂** (child)
  - **omɛnë** (six)
- Prefixes and suffixes (whether hyphenated or not) are not part of the word root and so do not participate in vowel harmony.
  - Cf. the gerund suffix **ngɔ̈**: *dëngɔ̈*, *kpëngbängɔ̈*, *mböngɔ̈*, *töngɔ̈*, *wököngɔ̈*

> Internally, the symbols **ə** and **ø** are used as placeholders when vowel height is unknown, but this convention is not used externally.

### Nasal vowels

Vowels ending in _n_ are pronounced as nasals (_in_ sounds like the French _hein_).

The open vowels **ɛ** and **ɔ** do not have nasal variants distinct from their close vowels **e** and **o**,
and both are canonically written using the close vowels **en** and **on** (even though arguably they are
closer in pronunciation to /ɛn/ and /ɔn/).

Note that **n** can be either vowel-final or consonant-initial but not both simultaneously,
and there is no ambiguity:
1. If followed by a vowel or by **d**, **g**, **y**, or **z**,
   the **n** is consonantal and starts the following syllable.
   - In the rare cases where this not intended, a hyphen is infixed, e.g. **hôn-gɛrɛ̂** ("nose-foot" = caterpillar).
2. Otherwise _V**n**_ is vowel-final (with **ɛn** and **ɔn** corrected to canonical **en** and **on**),
   e.g **hôntï** (= **hôn-tï** = "nose-arm" = wrist) with hyphen completely optional.
   - In most cases, a nasal vowel is also word-final, which makes parsing easy in practice.
3. Other cases are spelling errors.

Note that a double **-nn-** is possible: **fünngɔ̈** ("odor") has syllables **fün** followed by **ngɔ̈**.

### Pitch accent

The pitch accents comprise the closed set

_{**¨** (mid), **^** (high)}_

* Either may applied to any vowel.
* Unmarked syllables are low pitch (which is the most common vowel pitch).
* Unknown pitch is marked internally with a dot below the vowel.

## Lexical ordering

There is no universally agreed-upon lexical ordering of Sango words.
Different references have adopted different systems.

The most intuitive, simple, and still efficient scheme (in the world of computers)
is the one adopted in this project, namely the lexical UTF byte ordering (after NFKC normalization)
of the space-separated ordered pair (without, with) vowel height and pitch,
e.g. **hɔ̂ndɛ < hôntï** because **"honde hɔ̂ndɛ" < "honti hôntï"**.
This is the ordering of the rows in the [lexicon](../lexicon/lexicon.csv).

> One quirk of UTF8 is that circumflex is lexically less than diaeresis, so that **fa < fâ < fä**.