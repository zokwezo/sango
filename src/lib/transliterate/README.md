# Input ASCII encoding of UTF8

Entering UTF8 from the keyboard is cumbersome. This library and CLI tool allow for using an ASCII equivalent:

| UTF8 | ASCII | Description               |
| :--: | :---: | :------------------------ |
|  “   |  ``   | left double quotes        |
|  ”   |  ''   | right double quotes       |
|  ‘   |  \`   | left single quotes        |
|  ’   |   '   | right single quotes       |
|  «   |  <<   | left double angle quotes  |
|  »   |  >>   | right double angle quotes |
|  …   | . . . | ellipses                  |
|  -   |   -   | hyphen                    |
|  –   |  - -  | en-dash                   |
|  —   | - - - | em-dash                   |
|  ɛ   |   x   | open e                    |
|  ɔ   |   c   | open o                    |
|  ç   |  ,c   | c cedilla                 |
|  ə   |  .e   | e or ɛ (height unknown)   |
|  ø   |  .o   | o or ɔ (height unknown)   |
|  ạ   |   a   | vowel of unknown pitch    |
|  ạ   |  aJ   | vowel of low pitch        |
|  â   |  aj   | vowel with high pitch     |
|  ä   |  aq   | vowel with mid pitch      |

and repeating the last 4 rows for all vowels, both upper and lower case.
