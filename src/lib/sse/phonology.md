# Sango phonology

Sango is a strictly _C?V_ language (i.e. each syllable is either _V_ or _CV_).

There are no diphthongs: two juxtaposed vowels are pronounced individually.
This means that **bâa** is pronounced as two syllables (without introducing
any stop between the two vowels), twice as long as either **ba** or **bâ**.

## Consonants

The consonants comprise the closed set

_{**b, d, f, g, gb, h, k, kp, l, m, mb, mv, n, nd, ng, ngb, ny, nz, p, r, s, t, v, w, y, z**}_

## Vowels

The vowels comprise the closed set

__{**a, an, e, en, ɛ, i, in, o, on, ɔ, u, un**}__

The following _VV_ combinations are **not** found in Sango: **ei, ɛ**_V_**, ie, iɛ, ou, ɔu, uo, uɔ**

### Vowel height

Spoken Sango strictly distinguishes close and open middle vowels (**e**/**ɛ** and **o**/**ɔ**).

> Internally, the vowels (**ə**/**ø**) are used to indicate unknown vowel height,
> but this convention is not used externally.

Very unfortunately, written Sango usually does not (only because there are no keys on the typewriter for them).

Because technology has since improved, and many applications (including text-to-speech and language learning)
require making this distinction, and because it can only help the human reader, it is strongly recommended to
render this distinction wherever the technology allows for it. This also future-proofs applications, as it is
possible that future language standards may reincorporate vowel height distinction in written Sango.

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