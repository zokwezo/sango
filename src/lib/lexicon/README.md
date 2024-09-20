# Sango Lexicon

This directory contains a Sango-English lexicon and metadata, available in both row-major and column-major order, via a Go library.

This data can be extracted (if needed for other purposes) into CSV files with the bash shell commands:
```bash
echo '"Toneless","Sango","UDPos","UDFeature","Category","Frequency","EnglishTranslation","EnglishDefinition"' > /tmp/affixes.csv 
echo '"Toneless","Sango","UDPos","UDFeature","Category","Frequency","EnglishTranslation","EnglishDefinition"' > /tmp/lexicon.csv 
cat lexicon.go | sed -n "198,205p"  | sed -E 's/^\s*\{((.)*)\},/\1/' | sed -E 's/(["0-9],) /\1/g' >> /tmp/affixes.csv
cat lexicon.go | sed -n "210,2322p" | sed -E 's/^\s*\{((.)*)\},/\1/' | sed -E 's/(["0-9],) /\1/g' >> /tmp/lexicon.csv
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

| Name      | Type   | Description                                                                               |
| --------- | ------ | ----------------------------------------------------------------------------------------- |
| Toneless  | string | Lexeme after omitting pitch accent, vowel height, and punctuation                         |
| Sango     | string | Lexeme with accents and close/open vowel distinctions                                     |
| LexPos    | string | Lexical part-of-speech                                                                    |
| UDPos     | string | [Universal Dependency Part-of-speech](https://universaldependencies.org/u/pos/index.html) |
| UDFeature | string | [Universal Dependency Feature](https://universaldependencies.org/u/feat/)                 |
| Category  | string | Semantic cluster label                                                                    |
| Frequency | int    | Relative frequency (1=most common, 9=rare)                                                |
| English   | string | Brief English translation. Unrelated meanings are separated by semicolon                  |

Please note the following:

1. All columns have a non-empty value except UDFeature and Category which are empty if unclassified.
2. The first 6 columns (Toneless, Sango, LexPos, UDPos, UDFeature, Category) form a unique 6-tuple primary key, and rows
   are in strict ascending lexical key order to enable binary search and [lower, upper) bound interval lookups in the tables.
3. Hyphenation is used for clarity in separating Sango morphemes and is suitable for generation, but parsing
   should not depend on its presence as there is free variation in the use of punctuation in corpora.
4. The 6 productive affixes are separated out in their own table, being morphemes that can be prefixed or suffixed
   (depending on whether the LexPos is PREFIX or SUFFIX) via inner join on matching UDPos (e.g. NOUN or VERB)
   to any non-affix hyphenless Sango lexeme.
   - Note that the ngɔ̈ suffix enforces vowel harmony by changing all preceding pitch accents in the root lexeme
     (but not any other prefix or suffix) to circumflex (medium pitch), e.g.
     - **wa-** (one who) + **manda** (learn) + **-ngɔ̈** (-ing) + **kua** (work) = **wa-mändängɔ̈-kua** (apprentice).
5. **Pitch accent is always indicated in the official orthography**
   and is important to distinguish meanings and/or parts of speech, e.g.

   | Toneless | Sango | LexPos | UDPos | UDFeature   | Category | Frequency | English                         |
   | -------- | ----- | ------ | ----- | ----------- | -------- | :-------: | ------------------------------- |
   | iri      | îri   | VT     | VERB  | Subcat=Tran | INTERACT |     1     | call, name                      |
   | iri      | ïrï   | N      | NOUN  |             | INTERACT |     1     | name                            |
   | kua      | kua   | N      | NOUN  |             | CIVIL    |     2     | work, job, duty                 |
   | kua      | küä   | N      | NOUN  |             | BODY     |     3     | hair, fur, pelt, feathers, down |
   | kua      | kûâ   | N      | NOUN  |             | STATE    |     2     | death                           |

6. **Vowel height is not indicated in the official orthography** but is nonetheless important in aural understanding and therefore
   represented here in the Sango column explicitly.

   - Although easily restored from context in real time by native speakers when reading aloud, the open vowels Ɛ and Ɔ are used
     (leaving E and O to represent close vowels) in the Sango column to aid nonnative speakers and text-to-speech applications.
   - Conversion to the standard orthography is just a constant static many-to-one mapping and not worth persisting in the table.

   Consider the following distinctions in vowel height (which can arise with any pitch accent and are unfortunately neither
   productive nor easily predicted with transformation rules and must be cataloged explicitly as separate lexemes):

   | Toneless | Sango | LexPos | UDPos | UDFeature   | Category | Frequency | English                     |
   | -------- | ----- | ------ | ----- | ----------- | -------- | :-------: | --------------------------- |
   | de       | de    | VI     | VERB  | Subcat=Intr | BODY     |     3     | vomit                       |
   | de       | dɛ    | V      | VERB  |             | STATE    |     2     | remain                      |
   | de       | dë    | VI     | VERB  | Subcat=Intr | HOW      |     3     | be cold                     |
   | de       | dɛ̈    | VT     | VERB  | Subcat=Tran | ACT      |     2     | cut, slice; grow, cultivate |
   | de       | dê    | N      | NOUN  |             | HOW      |     3     | coldness, shade             |
   | de       | dɛ̈    | VT     | VERB  | Subcat=Tran | INTERACT |     2     | emit                        |

### Code Structure

The data is provided both as static comma-separated files [affixes.csv](affixes.csv) and
[lexicon.csv](lexicon.csv), as well as available to Go applications as a library, e.g.:

```go
package main

import (
  "fmt"
  "slices"
  "strings"

  l "github.com/zokwezo/sango/src/lib/lexicon"
)

func main() {
  type bytesField struct {
    name  string
    field *[][]byte
  }

  type runesField struct {
    name  string
    field *[][]rune
  }

  getBytesFields := func(d l.DictCols) []bytesField {
    return []bytesField{
      bytesField{"Toneless  ", &d.Toneless},
      bytesField{"SangoUTF8 ", &d.SangoUTF8},
      bytesField{"LexPos    ", &d.LexPos},
      bytesField{"UDPos     ", &d.UDPos},
      bytesField{"UDFeature ", &d.UDFeature},
      bytesField{"Category  ", &d.Category},
      bytesField{"English   ", &d.English},
    }
  }

  getRunesFields := func(d l.DictCols) []runesField {
    return []runesField{
      runesField{"Sango     ", &d.Sango},
    }
  }

  type Dict struct {
    name string
    rows l.DictRows
    cols l.DictCols
  }

  dicts := []Dict{
    {"AFFIXES", l.AffixesRows(), l.AffixesCols()},
    {"LEXICON", l.LexiconRows(), l.LexiconCols()},
  }

  for _, dict := range dicts {
    fmt.Printf("\n\n%v ROWS:\n%v\n", dict.name, dict.rows)
    fmt.Printf("\n\n%v COLS AS STRINGS:\n", dict.name)
    fmt.Printf("%v Cols[Frequency ] = %v\n", dict.name, dict.cols.Frequency)
    for _, bf := range getBytesFields(dict.cols) {
      for i, b := range *bf.field {
        fmt.Printf("%v Cols[%v][%v] = {%s}\n", dict.name, bf.name, i, string(b))
      }
    }
    for _, bf := range getRunesFields(dict.cols) {
      for i, b := range *bf.field {
        fmt.Printf("%v Cols[%v][%v] = {%s}\n", dict.name, bf.name, i, string(b))
      }
    }
  }

  // Look for first entry <= {Toneless: "butuma" Sango: "butuma"}
  cmpFunc := func(lhs, rhs l.DictRow) int {
    if cmp := strings.Compare(lhs.Toneless, rhs.Toneless); cmp != 0 {
      return cmp
    }
    return strings.Compare(lhs.Sango, rhs.Sango)
  }
  entry := l.DictRow{Toneless: "butuma", Sango: "butuma"}
  nBegin, found := slices.BinarySearchFunc(l.LexiconRows(), entry, cmpFunc)
  fmt.Printf("Looking for %v at entry[%v] (found = %v)\n", entry, nBegin, found)
  entry = l.DictRow{Toneless: "butuma", Sango: "butumb"}
  nEnd, _ := slices.BinarySearchFunc(l.LexiconRows(), entry, cmpFunc)
  for n := nBegin; n < nEnd; n++ {
    fmt.Printf("entry[%v] = %v\n", n, l.LexiconRows()[n])
  }
}
```
