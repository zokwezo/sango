# Sango tools

It is convenient to put all the Sango tools into one CLI.

## Syntax

```
sango help
```

## APPENDIX: Unicode glyphs used in Sango text

There are many good Unicode references online, including e.g.

- https://www.compart.com/en/unicode/ for individual runes
- https://util.unicode.org/UnicodeJsps/list-unicodeset.jsp when constructing regular expressions

For reference, all unicode glyphs (and their encodings) used in this library are listed below.

All glyphs can be represented in an NFD form, with base rune (NFD1 column), sometimes followed by combining mark rune (NFD2 column).

For many of these glyphs, there is also a single precomposed Unicode rune (NFC column), and this is always preferred when available.

### Punctuation

- Double quotes:
  - `"` ⟹ `“` (left double quote) and `”` (right double quote), alternating in a sentence.
  - `<<` ⟹ `«` (left angular brackets)
  - `>>` ⟹ `»` (right angular brackets)
- Connectors:
  - `...` ⟹ `…` (ellipsis)
  - `--` ⟹ `–` (n-dash)
  - `---` ⟹ `—` (m-dash)

### Upper case

| NFC UTF8 |  NFC   |  NFD1  |  NFD2  | Height  | Pitch   | ASCII |
| :------: | :----: | :----: | :----: | :------ | :------ | :---: |
|    A     | U+0041 | U+0041 |        | Close   | Low     |   A   |
|    Ä     | U+00C4 | U+0041 | U+0308 | Close   | Mid     |   A   |
|    Â     | U+00C2 | U+0041 | U+0302 | Close   | High    |   A   |
|    Ạ     | U+1EA0 | U+0041 | U+0323 | Close   | Unknown |   A   |
|    E     | U+0045 | U+0045 |        | Close   | Low     |   E   |
|    Ë     | U+00CB | U+0045 | U+0308 | Close   | Mid     |   E   |
|    Ê     | U+00CA | U+0045 | U+0302 | Close   | High    |   E   |
|    Ẹ     | U+1EB8 | U+0045 | U+0323 | Close   | Unknown |   E   |
|    Ɛ     | U+0190 | U+0190 |        | Open    | Low     |   E   |
|    Ɛ̈     |   ❌   | U+0190 | U+0308 | Open    | Mid     |   E   |
|    Ɛ̂     |   ❌   | U+0190 | U+0302 | Open    | High    |   E   |
|    Ɛ̣     |   ❌   | U+0190 | U+0323 | Open    | Unknown |   E   |
|    Ə     | U+018F | U+018F |        | Unknown | Low     |   E   |
|    Ə̈     |   ❌   | U+018F | U+0308 | Unknown | Mid     |   E   |
|    Ə̂     |   ❌   | U+018F | U+0302 | Unknown | High    |   E   |
|    Ə̣     |   ❌   | U+018F | U+0323 | Unknown | Unknown |   E   |
|    I     | U+0049 | U+0049 |        | Close   | Low     |   I   |
|    Ï     | U+00CF | U+0049 | U+0308 | Close   | Mid     |   I   |
|    Î     | U+00CE | U+0049 | U+0302 | Close   | High    |   I   |
|    Ị     | U+1ECA | U+0049 | U+0323 | Close   | Unknown |   I   |
|    O     | U+004F | U+004F |        | Close   | Low     |   O   |
|    Ö     | U+00D6 | U+004F | U+0308 | Close   | Mid     |   O   |
|    Ô     | U+00D4 | U+004F | U+0302 | Close   | High    |   O   |
|    Ọ     | U+1ECC | U+004F | U+0323 | Close   | Unknown |   O   |
|    Ɔ     | U+0186 | U+0186 |        | Open    | Low     |   O   |
|    Ɔ̈     |   ❌   | U+0186 | U+0308 | Open    | Mid     |   O   |
|    Ɔ̂     |   ❌   | U+0186 | U+0302 | Open    | High    |   O   |
|    Ɔ     |   ❌   | U+0186 | U+0323 | Open    | Unknown |   O   |
|    Ø     | U+00D8 | U+00D8 |        | Unknown | Low     |   O   |
|    Ø̈     |   ❌   | U+00D8 | U+0308 | Unknown | Mid     |   O   |
|    Ø̂     |   ❌   | U+00D8 | U+0302 | Unknown | High    |   O   |
|    Ø̣     |   ❌   | U+00D8 | U+0323 | Unknown | Unknown |   O   |
|    U     | U+0055 | U+0055 |        | Close   | Low     |   U   |
|    Ü     | U+00DC | U+0055 | U+0308 | Close   | Mid     |   U   |
|    Û     | U+00DB | U+0055 | U+0302 | Close   | High    |   U   |
|    Ụ     | U+1EE4 | U+0055 | U+0323 | Close   | Unknown |   U   |

### Lower case

| NFC UTF8 |  NFC   |  NFD1  |  NFD2  | Height  | Pitch   | ASCII |
| :------: | :----: | :----: | :----: | :------ | :------ | :---: |
|    a     | U+0061 | U+0061 |        | Close   | Low     |   a   |
|    ä     | U+00E4 | U+0061 | U+0308 | Close   | Mid     |   a   |
|    â     | U+00E2 | U+0061 | U+0302 | Close   | High    |   a   |
|    ạ     | U+1EA1 | U+0061 | U+0323 | Close   | Unknown |   a   |
|    e     | U+0065 | U+0065 |        | Close   | Low     |   e   |
|    ë     | U+00EB | U+0065 | U+0308 | Close   | Mid     |   e   |
|    ê     | U+00EA | U+0065 | U+0302 | Close   | High    |   e   |
|    ẹ     | U+1EB9 | U+0065 | U+0323 | Close   | Unknown |   e   |
|    ɛ     | U+025B | U+025B |        | Open    | Low     |   e   |
|    ɛ̈     |   ❌   | U+025B | U+0308 | Open    | Mid     |   e   |
|    ɛ̂     |   ❌   | U+025B | U+0302 | Open    | High    |   e   |
|    ɛ̣     |   ❌   | U+025B | U+0323 | Open    | Unknown |   e   |
|    ə     | U+0259 | U+0259 |        | Unknown | Low     |   e   |
|    ə̈     |   ❌   | U+0259 | U+0308 | Unknown | Mid     |   e   |
|    ə̂     |   ❌   | U+0259 | U+0302 | Unknown | High    |   e   |
|    ə̣     |   ❌   | U+0259 | U+0323 | Unknown | Unknown |   e   |
|    i     | U+0069 | U+0069 |        | Close   | Low     |   i   |
|    ï     | U+00EF | U+0069 | U+0308 | Close   | Mid     |   i   |
|    î     | U+00EE | U+0069 | U+0302 | Close   | High    |   i   |
|    ị     | U+1ECB | U+0069 | U+0323 | Close   | Unknown |   i   |
|    o     | U+006F | U+006F |        | Close   | Low     |   o   |
|    ö     | U+00F6 | U+006F | U+0308 | Close   | Mid     |   o   |
|    ô     | U+00F4 | U+006F | U+0302 | Close   | High    |   o   |
|    ọ     | U+1ECD | U+006F | U+0323 | Close   | Unknown |   o   |
|    ɔ     | U+0254 | U+0254 |        | Open    | Low     |   o   |
|    ɔ̈     |   ❌   | U+0254 | U+0308 | Open    | Mid     |   o   |
|    ɔ̂     |   ❌   | U+0254 | U+0302 | Open    | High    |   o   |
|    ɔ     |   ❌   | U+0254 | U+0323 | Open    | Unknown |   o   |
|    ø     | U+00F8 | U+00F8 |        | Unknown | Low     |   o   |
|    ø̈     |   ❌   | U+00F8 | U+0308 | Unknown | Mid     |   o   |
|    ø̂     |   ❌   | U+00F8 | U+0302 | Unknown | High    |   o   |
|    ø̣     |   ❌   | U+00F8 | U+0323 | Unknown | Unknown |   o   |
|    u     | U+0075 | U+0075 |        | Close   | Low     |   u   |
|    ü     | U+00FC | U+0075 | U+0308 | Close   | Mid     |   u   |
|    û     | U+00FB | U+0075 | U+0302 | Close   | High    |   u   |
|    ụ     | U+1EE5 | U+0075 | U+0323 | Close   | Unknown |   u   |
