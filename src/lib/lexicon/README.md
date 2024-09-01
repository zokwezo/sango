# Sango lexicon

This module contains an open-source Sango-English lexicon, available to Go code in both row-major and column-major order
and as static comma-separated files [affixes.csv](affixes.csv) and [lexicon.csv](lexicon.csv).

# License

The code and data are released under [Apache License Version 2.0, January 2004](https://www.apache.org/licenses/LICENSE-2.0)
in order to facilitate and encourage development of Sango language understanding tools.

# Acknowledgements

This lexicon is only best effort and does not purport to meet the excellent quality of the two extant professional
Sango-French dictionaries:

- Koyt-Deballé, G. F. (2013). _Lexique illustré sängö-français - français-sängö_, Éditions universitaires européennes. 476 p. ISBN 978-6131592690.
- Bouquiaux, L. et al (1978). _Dictionnaire sango-français = Kété bàkàrī sāngō-fàránzì_, Société d'études linguistiques et anthropologiques de France. 663 p. ISBN : 2-85297-016-3.

> _It is not known whether a machine-readable version of either of the above exists, but professional or commercial users
> may wish to reach out to the publisher directly and (if warranted and possible) obtain a license to undertake to digitize one of these._

The English translation and linguistic annotations in this lexicon are the original creation of the author, based on analysis
of secular and religious texts, government publications, and personal knowledge and notes accumulated over two years by the author
in the Central African Republic in 1988-90 by the author, albeit before the importance of pitch or height accent was realized.
These accents, along with some less frequent or neologistic lexemes or meanings, have been manually restored in consultation
with the above two published references, but responsibility for any errata rests entirely with the author.

If you find this dictionary has been useful to you, I would love to hear about it. Please email westondan@zokwezo.net describing
your use case to help me gauge impact and motivate future work.

In publications, this work can be cited as:

- Weston, D. D. (2024). _Sango-English lexicon_. https://github.com/zokwezo/sango/tree/main/src/lib/lexicon.

# Data Structure

Each lexicon row contains the following columns:

| Name      | Type   | Description                                                                               |
| --------- | ------ | ----------------------------------------------------------------------------------------- |
| Toneless  | string | Lexeme after omitting pitch accent, vowel height, and punctuation                         |
| Sango     | string | Lexeme with accents and closed/open vowel distinctions                                    |
| LexPos    | string | Lexical part-of-speech                                                                    |
| UDPos     | string | [Universal Dependency Part-of-speech](https://universaldependencies.org/u/pos/index.html) |
| Frequency | int    | Relative frequency (1=most common, 9=rare)                                                |
| UDFeature | string | [Universal Dependency Feature](https://universaldependencies.org/u/feat/)                 |
| Category  | string | Semantic cluster label                                                                    |
| English   | string | Brief English translation. Unrelated meanings are separated by semicolon                  |

Note the following:

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
     - **wa-** (one who) + **manda** (learn) + **-ngɔ̈** (-ing) + **kua** (work) = **wa-mändängɔ̈-kua** (apprentice).
6. **Pitch accent is always indicated in the official orthography**
   and is important to distinguish meanings and/or parts of speech, e.g.

   | Toneless | Sango | LexPos | UDPos | UDFeature   | Category | Frequency | English                                     |
   | -------- | ----- | ------ | ----- | ----------- | -------- |:---------:| ------------------------------------------- |
   | iri      | îri   | VT     | VERB  | Subcat=Tran | INTERACT |     1     | call, name                                  |
   | iri      | ïrï   | N      | NOUN  |             | INTERACT |     1     | name                                        |
   | kua      | kua   | N      | NOUN  |             | CIVIL    |     2     | work, job, duty                             |
   | kua      | küä   | N      | NOUN  |             | BODY     |     3     | hair, fur, pelt, feathers, down             |
   | kua      | kûâ   | N      | NOUN  |             | STATE    |     2     | death                                       |

7. **Vowel height is not indicated in the official orthography** but is nonetheless important in aural understanding.
   Although easily restored from context in real time by native speakers when reading aloud, the open vowels Ɛ and Ɔ are used
   (leaving E and O to represent closed vowels) in the Sango column to aid nonnative speakers and text-to-speech applications.
   Consider the following distinctions in vowel height (which may be present in any pitch accent):

   | Toneless | Sango | LexPos | UDPos | UDFeature   | Category | Frequency | English                                     |
   | -------- | ----- | ------ | ----- | ----------- | -------- |:---------:| ------------------------------------------- |
   | de       | de    | VI     | VERB  | Subcat=Intr | BODY     |     3     | vomit                                       |
   | de       | dɛ    | V      | VERB  |             | STATE    |     2     | remain                                      |
   | de       | dë    | VI     | VERB  | Subcat=Intr | HOW      |     3     | be cold                                     |
   | de       | dɛ̈    | VT     | VERB  | Subcat=Tran | ACT      |     2     | cut, slice; grow, cultivate                 |
   | de       | dê    | N      | NOUN  |             | HOW      |     3     | coldness, shade                             |
   | de       | dɛ̈    | VT     | VERB  | Subcat=Tran | INTERACT |     2     | emit                                        |

# Usage

The data is provided both as static comma-separated files [affixes.csv](affixes.csv) and
[lexicon.csv](lexicon.csv), as well as available to Go applications as a library, e.g.:

```go
package main

import (
  "fmt"

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
}
```
