# Sango Lexicon

This directory contains a Sango-English lexicon and metadata, available in both row-major and column-major order, via a Go library.

The lexicon data can be extracted from code (if needed for other purposes) into a standalone CSV file with the following bash shell commands:

```bash
outfile="/tmp/lexicon.csv"
echo '"Toneless","Heightless","Lemma","UDPos","UDFeature","Category","Frequency","EnglishTranslation","EnglishDefinition"' >  "${outfile}"
cat lexicon.go | grep -E '^\s*{"[a-z]*",' | head -n -9 | sed -E 's/^\s*\{((.)*)\},/\1/' | sed -E 's/(["0-9],) /\1/g' >> "${outfile}"
```

### Background

The English translations and linguistic annotations in this lexicon are based on my analysis of the limited written corpus I could
locate (secular and religious texts, government publications), combined with personal knowledge and notes and merged into a lexicon
I accumulated over two years in the Central African Republic in 1988-90.

Most of this source material was not in digital form and dates to before the orthography was standardized in 1984 (after which
high and mid vowel pitch were represented in writing with circumflex and diaeresis, respectively). Also, in retrospect,
it is important to distinguishing between open (Ɛ and Ɔ) and close (E and O) vowel height, and to encode these in writing
to assist nonnative speakers and to support text-to-speech. Consequently, I have since manually restored these
in consultation with the three published references below, but responsibility for any errata rests entirely with me.

Please email me (Dan Weston <westondan@zokwezo.net>) if you:

- find this dictionary has been useful to you or your project
- wish to collaborate/co-contribute to this project
- have any questions, suggestions, errata, or bug reports

In publications, this work can be cited as:

- Weston, D. D. (2024). _Sango-English lexicon_. https://github.com/zokwezo/sango/tree/main/src/lib/lexicon.

### External Resources

This lexicon is only best effort and does not purport to meet the excellent quality of the (only) two published professional
Sango-French dictionaries and orthography manual:

- Koyt-Deballé, G. F. (2013). _Lexique illustré sängö-français - français-sängö_, Éditions universitaires européennes. 476 p. ISBN 978-6131592690.
- Bouquiaux, L. et al (1978). _Dictionnaire sango-français = Kété bàkàrī sāngō-fàránzì_, Société d'études linguistiques et anthropologiques de France. 663 p. ISBN 2-85297-016-3.
- Diki-Dikiri, M. (1977). _Le Sango s'écrit aussi… Esquisse Linguistique du Sango, Langue Nationale de l'Empire Centrafricain_, Selaf-Paris. 187 p. ISBN 2-85297-057-0.

I strongly encourage users of this project to also acquire a personal copy of the first (or better, all three!) of these to serve as source of truth
(and of course to support the authors and publishers!). In particular, _Lexique illustré sängö-français - français-sängö_, which reflects changes to the
language in the subsequent 35 years (including computer terms) and exhaustively catalogs the Sango words (along with the scientific names) of
indigenous Central African flora and fauna, of which only the most basic are found therein.

### Copyright and License

Copyright 2024 Daniel D. Weston

Except as otherwise noted, all code and data in this repository are the original work
of the copyright holder, or else derived from fair use of reference materials generally
available to the public.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

### Data Structure

Each lexicon row contains the following columns, intended i.a. to establish a reference for a future Sango treebank:

| Name                | Type   | Description                                                                               |
| ------------------- | ------ | ----------------------------------------------------------------------------------------- |
| Toneless            | string | = Heightless column, but with whitespace, hyphens, and pitch removed                      |
| Heightless          | string | = Lemma column, but with case and height removed                                          |
| Lemma               | string | Lemma (possibly multiword) with accents and close/open vowel distinctions                 |
| UDPos               | string | [Universal Dependency Part-of-speech](https://universaldependencies.org/u/pos/index.html) |
| UDFeature           | string | [Universal Dependency Feature](https://universaldependencies.org/u/feat/)                 |
| Category            | string | Semantic cluster label                                                                    |
| Frequency           | int    | Relative frequency (1=most common, 9=rare)                                                |
| English Translation | string | English translation (one or two words, omitting less frequent usage)                      |
| English Definition  | string | Full English description (definition, etymology, usage notes)                             |

Please note the following:

1. The first lexicon row is a copyright and license statement, which must remain with the data (whether stored or in memory).
   It is the only row with empty values for Toneless, Heightless, and Lemma columns.
2. The columns (`Lemma`, `UDPos`, `UDFeature`, `Category`, `Frequency`) form a minimal 5-tuple primary key and the rows
   are sorted in strict ascending order to enable binary search and [lower, upper) bound interval lookups in the tables.
3. Hyphenation is used for clarity in separating Sango morphemes and is suitable for generation, but parsing
   should not depend on its presence as there is free variation in the use of punctuation in corpora.
4. **Pitch accent is always indicated in the official orthography**
   and is important to distinguish meanings and/or parts of speech, e.g.

   | Toneless | Lemma | UDPos | UDFeature   | Category | English Translation             |
   | -------- | ----- | ----- | ----------- | -------- | ------------------------------- |
   | iri      | îri   | VERB  | Subcat=Tran | INTERACT | call, name                      |
   | iri      | ïrï   | NOUN  |             | INTERACT | name                            |
   | kua      | kua   | NOUN  |             | CIVIL    | work, job, duty                 |
   | kua      | kûâ   | NOUN  |             | STATE    | death                           |
   | kua      | küä   | NOUN  |             | BODY     | hair, fur, pelt, feathers, down |

5. **Vowel height is not indicated in the official orthography** but is nonetheless important in aural understanding and therefore
   represented here in the Lemma column explicitly.

   - Although easily restored from context in real time by native speakers when reading aloud, the open vowels **Ɛ** and **Ɔ** are used
     (leaving **E** and **O** to represent close vowels) in the Lemma column to aid nonnative speakers and text-to-speech applications.
   - Conversion to the standard orthography is just a constant static many-to-one mapping and not worth persisting in the table.
   - There is a strong tendency towards internal vowel harmony, so that open and close vowels do not co-occur in the same lemma.

   Consider the following distinctions in vowel height (which can arise with any pitch accent and are unfortunately neither
   productive nor easily predicted with transformation rules and must be cataloged explicitly as separate lexemes):

   | Lemma | UDPos | UDFeature   | Category | English Translation         |
   | ----- | ----- | ----------- | -------- | --------------------------- |
   | de    | VERB  | Subcat=Intr | BODY     | vomit                       |
   | dɛ    | VERB  |             | STATE    | remain                      |
   | dë    | VERB  | Subcat=Intr | HOW      | be cold                     |
   | dɛ̈    | VERB  | Subcat=Tran | ACT      | cut, slice; grow, cultivate |
   | dê    | NOUN  |             | HOW      | coldness, shade             |
   | dɛ̈    | VERB  | Subcat=Tran | INTERACT | emit                        |

6. There are 8 affixes sufficiently productive that they are automatically affixed to all compatible lemmas
   (governed by UDFeature `CanPrefix=`_UDPos_) on startup to generate derived lemmas. These are:

   | Lemma | UDPos | UDFeature                              | Category | English Definition  |
   | ----- | ----- | -------------------------------------- | -------- | ------------------- |
   | a     | VERB  | CanPrefix=VERB\|Person=3\|VerbForm=Fin | WHO      | subject marker      |
   | â     | ADJ   | CanPrefix=ADJ\|Number=Plur             | NUM      | plural marker       |
   | â     | NOUN  | CanPrefix=NOUN\|Number=Plur            | NUM      | plural marker       |
   | bâ    | NOUN  | CanPrefix=NOUN                         | WHERE    | canonical place for |
   | nga   | VERB  | Aspect=Iter\|CanSuffix=VERB            | HOW      | periodic action     |
   | ngbi  | VERB  | Aspect=Imp\|CanSuffix=VERB\|Reflex=Yes | HOW      | synchronic action   |
   | ngɔ̈   | VERB  | CanSuffix=VERB\|VerbForm=Vnoun         | HOW      | gerund ("-ing")     |
   | wa    | NOUN  | CanPrefix=VERB                         | WHO      | agent ("one who")   |

   - This attempts to generate the [convex set](https://en.wikipedia.org/wiki/Convex_set) of all possible lemmas,
     not a minimal plausible set, and will generate many lemmas not found in native speech, and is suitable for
     language understanding but not language generation (which would require a curated outer product).
   - The **ngɔ̈** suffix enforces vowel _pitch_ harmony by changing all preceding pitch accents in the root lexeme
     (but not any other prefix or suffix) to circumflex (medium pitch), e.g.
     - **wa-** (one who) + **manda** (learn) + **-ngɔ̈** (-ing) + **kua** (work) = **wa-mändängɔ̈-kua** (apprentice).
   - The **ngɔ̈** suffix does NOT enforce vowel _height_ harmony, and leaves close root vowels closed.
   - The **ngɔ̈** suffix, when affixed to stative verbs, functions as both gerund (action) and noun (state),
     as in **nɛ** = "weigh, to be heavy" ⇒ **nɛ̈ngɔ̈** = "weighing, weight".

7. There are other affixes that had at one time been productive (esp. in progenitor tribal languages such as Ngbandi) but
   no longer sufficiently productive in Sango to generate automatically, and are listed as explicit lexemes in the lexicon:

   - initial syllable (or word) reduplication

     | Lemma      | UDPos | UDFeature                   | Category | English Definition                                         |
     | ---------- | ----- | --------------------------- | -------- | ---------------------------------------------------------- |
     | fadë       | ADV   |                             | WHEN     | [postverbal]: right now; [preverbal]: will, shall          |
     | fafadësô   | ADV   |                             | WHEN     | immediately                                                |
     | fadëfadësô | ADV   |                             | WHEN     | immediately                                                |
     | dɔdɔ       | VERB  | Subcat=Tran                 | BODY     | dance                                                      |
     | dɔngɔ̈ dɔ̈dɔ̈ | VERB  | Subcat=Intr\|VerbForm=Vnoun | BODY     | dancing [the expected gerund form *dɔ̈dɔ̈ngɔ̈?* is not found] |
     | tɛnɛ       | VERB  | Subcat=Tran                 | INTERACT | say, tell                                                  |
     | tɛ̈nɛ̈       | NOUN  |                             | INTERACT | problem, quarrel; speech, talk, words, tale                |
     | tɛnɛ tɛ̈nɛ̈  | VERB  | Subcat=Intr                 | INTERACT | talk                                                       |

   - final syllable reduplication (often presenting as **-ra**) = similar to **-nga** but as a diminuitive repetitive action, cf.

     | Lemma  | UDPos | UDFeature   | Category | English Definition                               |
     | ------ | ----- | ----------- | -------- | ------------------------------------------------ |
     | fâa    | VERB  | Subcat=Tran | ACT      | cross, traverse; cut, strike, break; wound, kill |
     | fângbi | VERB  | Subcat=Tran | ACT      | cut up, dissect (into large pieces)              |
     | fâra   | VERB  | Subcat=Tran | ACT      | chop up, shred (into small pieces)               |

   - high pitch morphology to indicate irrealis mood (moribund, now subsumed by the subordinating conjunction **töngana**)

     | Lemma           | UDPos | UDFeature                                     | Category | English Definition                  |
     | --------------- | ----- | --------------------------------------------- | -------- | ----------------------------------- |
     | atɛnɛ           | VERB  | Mood=Ind\|Person=3\|Subcat=Tran\|VerbForm=Fin | INTERACT | one says/tells                      |
     | âtɛnɛ [rare]    | VERB  | Mood=Irr\|Person=3\|Subcat=Tran\|VerbForm=Fin | INTERACT | if one had said/told                |
     | töngana lo tɛnɛ | VERB  | Mood=Cnd\|Person=3\|Subcat=Tran\|VerbForm=Fin | INTERACT | if he/she says/tells, had said/told |
