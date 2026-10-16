package parse

import (
	"strings"
	"testing"
)

func TestParseShortCanonical(t *testing.T) {
	s := ` § ~ha_HO:Do:ni^U+003F 日本語はU+96E3しい,`
	parseMatcherForGroupName := map[string]ParseMatcher{
		"consonant":         Literal,
		"hexcode":           Literal,
		"pitch":             Literal,
		"sentence":          CanonicalWord,
		"sentenceUnmatched": Literal,
		"shift":             Literal,
		"syllable":          Literal,
		"syllables":         CanonicalSyllable,
		"symbol":            Literal,
		"vowelNasal":        Literal,
		"whole":             Sentence,
		"word":              Literal,
	}
	actual := Parse(parseMatcherForGroupName, s, ParseTree{Begin: 0, End: len(s), Name: "whole"})
	var actualBuilder strings.Builder = actual.Annotate(s, true, map[string]bool{
		"closeVowel":        true,
		"closeVowelNasal":   true,
		"consonant":         true,
		"hexcode":           true,
		"nasal":             true,
		"openVowel":         true,
		"pitch":             true,
		"sentence":          true,
		"sentenceUnmatched": true,
		"shift":             true,
		"syllable":          true,
		"syllables":         true,
		"symbol":            true,
		"vowelNasal":        true,
		"whole":             true,
		"word":              true,
	})
	actualStr := actualBuilder.String()
	expectStr := `
whole                [  0: 49] =  § ~ha_HO:Do:ni^U+003F 日本語はU+96E3しい,
sentence             [  0: 24] =  § ~ha_HO:Do:ni^U+003F 
symbol               [  0:  1] =  
symbol               [  1:  3] =  §
symbol               [  3:  4] =    
word                 [  4: 17] =    ~ha_HO:Do:ni^
shift                [  4:  5] =    ~
syllables            [  5: 17] =     ha_HO:Do:ni^
syllable             [  5:  8] =     ha_
consonant            [  5:  6] =     h
vowelNasal           [  6:  7] =      a
pitch                [  7:  8] =       _
syllable             [  8: 11] =        HO:
consonant            [  8:  9] =        H
vowelNasal           [  9: 10] =         O
pitch                [ 10: 11] =          :
syllable             [ 11: 14] =           Do:
consonant            [ 11: 12] =           D
vowelNasal           [ 12: 13] =            o
pitch                [ 13: 14] =             :
syllable             [ 14: 17] =              ni^
consonant            [ 14: 15] =              n
vowelNasal           [ 15: 16] =               i
pitch                [ 16: 17] =                ^
hexcode              [ 19: 23] =                   003F
symbol               [ 23: 24] =                        
sentence             [ 24: 49] =                        日本語はU+96E3しい,
symbol               [ 24: 27] =                        日
symbol               [ 27: 30] =                          本
symbol               [ 30: 33] =                            語
symbol               [ 33: 36] =                              は
hexcode              [ 38: 42] =                                  96E3
symbol               [ 42: 45] =                                      し
symbol               [ 45: 48] =                                        い
symbol               [ 48: 49] =                                          ,
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}

func TestParseLongCanonical(t *testing.T) {
	s := `U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7U+0020~ha_HO:U+002D~Do:ni^U+0020=ha_HO:U+002D=Do:ni^U+0020~ha^Dx_U+0020~ba^ha_U+002D~mo_U+002D~tx_nx_U+0020~bx^U+002D~kc:Bi:tx_bx^U+002D~kc:Bi:tx_U+96E3=bx^U+002D=kc:=Bi:=tx_U+0020~bx^U+002D#kc:Bi:tx_U+0020~bx^U+002D@kc:Bi:tx_U+0020=bx^U+002D=kc:=Bi:=tx_U+0020#ha_HO:Do:ni^`
	parseMatcherForGroupName := map[string]ParseMatcher{
		"consonant":  Literal,
		"hexcode":    Literal,
		"pitch":      Literal,
		"sentence":   CanonicalWord,
		"shift":      Literal,
		"syllables":  CanonicalSyllable,
		"symbol":     Literal,
		"vowelNasal": Literal,
		"whole":      Sentence,
	}
	actual := Parse(parseMatcherForGroupName, s, ParseTree{Begin: 0, End: len(s), Name: "whole"})
	var actualBuilder strings.Builder = actual.Annotate(s, false, map[string]bool{
		"consonant":  true,
		"hexcode":    true,
		"pitch":      true,
		"sentence":   true,
		"shift":      true,
		"syllables":  true,
		"symbol":     true,
		"vowelNasal": true,
		"whole":      true,
	})
	actualStr := actualBuilder.String()
	expectStr := `
whole                [  0:333] = "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020U+00A7U+0020~ha_HO:U+002D~Do:ni^U+0020=ha_HO:U+002D=Do:ni^U+0020~ha^Dx_U+0020~ba^ha_U+002D~mo_U+002D~tx_nx_U+0020~bx^U+002D~kc:Bi:tx_bx^U+002D~kc:Bi:tx_U+96E3=bx^U+002D=kc:=Bi:=tx_U+0020~bx^U+002D#kc:Bi:tx_U+0020~bx^U+002D@kc:Bi:tx_U+0020=bx^U+002D=kc:=Bi:=tx_U+0020#ha_HO:Do:ni^"
sentence             [  0: 54] = "U+65E5U+672CU+8A9EU+306FU+96E3U+3057U+3044U+0021U+0020"
hexcode              [  2:  6] = "65E5"
hexcode              [  8: 12] = "672C"
hexcode              [ 14: 18] = "8A9E"
hexcode              [ 20: 24] = "306F"
hexcode              [ 26: 30] = "96E3"
hexcode              [ 32: 36] = "3057"
hexcode              [ 38: 42] = "3044"
hexcode              [ 44: 48] = "0021"
hexcode              [ 50: 54] = "0020"
sentence             [ 54:333] = "U+00A7U+0020~ha_HO:U+002D~Do:ni^U+0020=ha_HO:U+002D=Do:ni^U+0020~ha^Dx_U+0020~ba^ha_U+002D~mo_U+002D~tx_nx_U+0020~bx^U+002D~kc:Bi:tx_bx^U+002D~kc:Bi:tx_U+96E3=bx^U+002D=kc:=Bi:=tx_U+0020~bx^U+002D#kc:Bi:tx_U+0020~bx^U+002D@kc:Bi:tx_U+0020=bx^U+002D=kc:=Bi:=tx_U+0020#ha_HO:Do:ni^"
hexcode              [ 56: 60] = "00A7"
hexcode              [ 62: 66] = "0020"
shift                [ 66: 67] = "~"
syllables            [ 67: 73] = "ha_HO:"
consonant            [ 67: 68] = "h"
vowelNasal           [ 68: 69] = "a"
pitch                [ 69: 70] = "_"
consonant            [ 70: 71] = "H"
vowelNasal           [ 71: 72] = "O"
pitch                [ 72: 73] = ":"
hexcode              [ 75: 79] = "002D"
shift                [ 79: 80] = "~"
syllables            [ 80: 86] = "Do:ni^"
consonant            [ 80: 81] = "D"
vowelNasal           [ 81: 82] = "o"
pitch                [ 82: 83] = ":"
consonant            [ 83: 84] = "n"
vowelNasal           [ 84: 85] = "i"
pitch                [ 85: 86] = "^"
hexcode              [ 88: 92] = "0020"
shift                [ 92: 93] = "="
syllables            [ 93: 99] = "ha_HO:"
consonant            [ 93: 94] = "h"
vowelNasal           [ 94: 95] = "a"
pitch                [ 95: 96] = "_"
consonant            [ 96: 97] = "H"
vowelNasal           [ 97: 98] = "O"
pitch                [ 98: 99] = ":"
hexcode              [101:105] = "002D"
shift                [105:106] = "="
syllables            [106:112] = "Do:ni^"
consonant            [106:107] = "D"
vowelNasal           [107:108] = "o"
pitch                [108:109] = ":"
consonant            [109:110] = "n"
vowelNasal           [110:111] = "i"
pitch                [111:112] = "^"
hexcode              [114:118] = "0020"
shift                [118:119] = "~"
syllables            [119:125] = "ha^Dx_"
consonant            [119:120] = "h"
vowelNasal           [120:121] = "a"
pitch                [121:122] = "^"
consonant            [122:123] = "D"
vowelNasal           [123:124] = "x"
pitch                [124:125] = "_"
hexcode              [127:131] = "0020"
shift                [131:132] = "~"
syllables            [132:138] = "ba^ha_"
consonant            [132:133] = "b"
vowelNasal           [133:134] = "a"
pitch                [134:135] = "^"
consonant            [135:136] = "h"
vowelNasal           [136:137] = "a"
pitch                [137:138] = "_"
hexcode              [140:144] = "002D"
shift                [144:145] = "~"
syllables            [145:148] = "mo_"
consonant            [145:146] = "m"
vowelNasal           [146:147] = "o"
pitch                [147:148] = "_"
hexcode              [150:154] = "002D"
shift                [154:155] = "~"
syllables            [155:161] = "tx_nx_"
consonant            [155:156] = "t"
vowelNasal           [156:157] = "x"
pitch                [157:158] = "_"
consonant            [158:159] = "n"
vowelNasal           [159:160] = "x"
pitch                [160:161] = "_"
hexcode              [163:167] = "0020"
shift                [167:168] = "~"
syllables            [168:171] = "bx^"
consonant            [168:169] = "b"
vowelNasal           [169:170] = "x"
pitch                [170:171] = "^"
hexcode              [173:177] = "002D"
shift                [177:178] = "~"
syllables            [178:190] = "kc:Bi:tx_bx^"
consonant            [178:179] = "k"
vowelNasal           [179:180] = "c"
pitch                [180:181] = ":"
consonant            [181:182] = "B"
vowelNasal           [182:183] = "i"
pitch                [183:184] = ":"
consonant            [184:185] = "t"
vowelNasal           [185:186] = "x"
pitch                [186:187] = "_"
consonant            [187:188] = "b"
vowelNasal           [188:189] = "x"
pitch                [189:190] = "^"
hexcode              [192:196] = "002D"
shift                [196:197] = "~"
syllables            [197:206] = "kc:Bi:tx_"
consonant            [197:198] = "k"
vowelNasal           [198:199] = "c"
pitch                [199:200] = ":"
consonant            [200:201] = "B"
vowelNasal           [201:202] = "i"
pitch                [202:203] = ":"
consonant            [203:204] = "t"
vowelNasal           [204:205] = "x"
pitch                [205:206] = "_"
hexcode              [208:212] = "96E3"
shift                [212:213] = "="
syllables            [213:216] = "bx^"
consonant            [213:214] = "b"
vowelNasal           [214:215] = "x"
pitch                [215:216] = "^"
hexcode              [218:222] = "002D"
shift                [222:223] = "="
syllables            [223:226] = "kc:"
consonant            [223:224] = "k"
vowelNasal           [224:225] = "c"
pitch                [225:226] = ":"
shift                [226:227] = "="
syllables            [227:230] = "Bi:"
consonant            [227:228] = "B"
vowelNasal           [228:229] = "i"
pitch                [229:230] = ":"
shift                [230:231] = "="
syllables            [231:234] = "tx_"
consonant            [231:232] = "t"
vowelNasal           [232:233] = "x"
pitch                [233:234] = "_"
hexcode              [236:240] = "0020"
shift                [240:241] = "~"
syllables            [241:244] = "bx^"
consonant            [241:242] = "b"
vowelNasal           [242:243] = "x"
pitch                [243:244] = "^"
hexcode              [246:250] = "002D"
shift                [250:251] = "#"
syllables            [251:260] = "kc:Bi:tx_"
consonant            [251:252] = "k"
vowelNasal           [252:253] = "c"
pitch                [253:254] = ":"
consonant            [254:255] = "B"
vowelNasal           [255:256] = "i"
pitch                [256:257] = ":"
consonant            [257:258] = "t"
vowelNasal           [258:259] = "x"
pitch                [259:260] = "_"
hexcode              [262:266] = "0020"
shift                [266:267] = "~"
syllables            [267:270] = "bx^"
consonant            [267:268] = "b"
vowelNasal           [268:269] = "x"
pitch                [269:270] = "^"
hexcode              [272:276] = "002D"
shift                [276:277] = "@"
syllables            [277:286] = "kc:Bi:tx_"
consonant            [277:278] = "k"
vowelNasal           [278:279] = "c"
pitch                [279:280] = ":"
consonant            [280:281] = "B"
vowelNasal           [281:282] = "i"
pitch                [282:283] = ":"
consonant            [283:284] = "t"
vowelNasal           [284:285] = "x"
pitch                [285:286] = "_"
hexcode              [288:292] = "0020"
shift                [292:293] = "="
syllables            [293:296] = "bx^"
consonant            [293:294] = "b"
vowelNasal           [294:295] = "x"
pitch                [295:296] = "^"
hexcode              [298:302] = "002D"
shift                [302:303] = "="
syllables            [303:306] = "kc:"
consonant            [303:304] = "k"
vowelNasal           [304:305] = "c"
pitch                [305:306] = ":"
shift                [306:307] = "="
syllables            [307:310] = "Bi:"
consonant            [307:308] = "B"
vowelNasal           [308:309] = "i"
pitch                [309:310] = ":"
shift                [310:311] = "="
syllables            [311:314] = "tx_"
consonant            [311:312] = "t"
vowelNasal           [312:313] = "x"
pitch                [313:314] = "_"
hexcode              [316:320] = "0020"
shift                [320:321] = "#"
syllables            [321:333] = "ha_HO:Do:ni^"
consonant            [321:322] = "h"
vowelNasal           [322:323] = "a"
pitch                [323:324] = "_"
consonant            [324:325] = "H"
vowelNasal           [325:326] = "O"
pitch                [326:327] = ":"
consonant            [327:328] = "D"
vowelNasal           [328:329] = "o"
pitch                [329:330] = ":"
consonant            [330:331] = "n"
vowelNasal           [331:332] = "i"
pitch                [332:333] = "^"
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}

func TestParseSango(t *testing.T) {
	s := `Ahöñ-ndönî AHÖÑ-NDÖNÎ, ândɛ bâa-mo-tɛnɛ? BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ!`
	parseMatcherForGroupName := map[string]ParseMatcher{
		"closeVowel":                 Literal,
		"closeVowelNasal":            Literal,
		"closeVowelNasalUnmatched":   Literal,
		"closeVowelUnmatched":        Literal,
		"consonant":                  Literal,
		"consonantUnmatched":         Literal,
		"nasal":                      Literal,
		"nasalUnmatched":             Literal,
		"openVowel":                  Literal,
		"openVowelUnmatched":         Literal,
		"pitch":                      Literal,
		"pitchUnmatched":             Literal,
		"sentence":                   SangoWord,
		"sentenceUnmatched":          Literal,
		"sentenceUnmatchedUnmatched": Literal,
		"shift":                      Literal,
		"shiftUnmatched":             Literal,
		"syllable":                   Literal,
		"syllableUnmatched":          Literal,
		"syllables":                  SangoSyllable,
		"syllablesUnmatched":         Literal,
		"vowelNasal":                 Literal,
		"vowelNasalUnmatched":        Literal,
		"whole":                      Sentence,
		"wholeUnmatched":             Literal,
		"word":                       Literal,
		"wordUnmatched":              Literal,
	}
	actual := Parse(parseMatcherForGroupName, s, ParseTree{Begin: 0, End: len(s), Name: "whole"})
	var actualBuilder strings.Builder = actual.Annotate(s, true, map[string]bool{
		"closeVowel":        true,
		"closeVowelNasal":   true,
		"consonant":         true,
		"nasal":             true,
		"openVowel":         true,
		"pitch":             true,
		"sentence":          true,
		"sentenceUnmatched": true,
		"shift":             true,
		"syllable":          true,
		"syllables":         true,
		"unmatched":         true,
		"vowelNasal":        true,
		"whole":             true,
		"word":              true,
	})
	actualStr := actualBuilder.String()
	expectStr := `
whole                [  0:112] = Ahöñ-ndönî AHÖÑ-NDÖNÎ, ândɛ bâa-mo-tɛnɛ? BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ!
sentence             [  0: 54] = Ahöñ-ndönî AHÖÑ-NDÖNÎ, ândɛ bâa-mo-tɛnɛ? 
word                 [  0:  6] = Ahöñ
syllables            [  0:  6] = Ahöñ
syllable             [  0:  1] = A
vowelNasal           [  0:  1] = A
closeVowelNasal      [  0:  1] = A
closeVowel           [  0:  1] = A
syllable             [  1:  6] =  höñ
consonant            [  1:  2] =  h
vowelNasal           [  2:  6] =   öñ
closeVowelNasal      [  2:  6] =   öñ
closeVowel           [  2:  4] =   ö
nasal                [  4:  6] =    ñ
sentenceUnmatched    [  6:  7] =     -
word                 [  7: 14] =      ndönî
syllables            [  7: 14] =      ndönî
syllable             [  7: 11] =      ndö
consonant            [  7:  9] =      nd
vowelNasal           [  9: 11] =        ö
closeVowelNasal      [  9: 11] =        ö
closeVowel           [  9: 11] =        ö
syllable             [ 11: 14] =         nî
consonant            [ 11: 12] =         n
vowelNasal           [ 12: 14] =          î
closeVowelNasal      [ 12: 14] =          î
closeVowel           [ 12: 14] =          î
sentenceUnmatched    [ 14: 15] =            
word                 [ 15: 21] =            AHÖÑ
syllables            [ 15: 21] =            AHÖÑ
syllable             [ 15: 16] =            A
vowelNasal           [ 15: 16] =            A
closeVowelNasal      [ 15: 16] =            A
closeVowel           [ 15: 16] =            A
syllable             [ 16: 21] =             HÖÑ
consonant            [ 16: 17] =             H
vowelNasal           [ 17: 21] =              ÖÑ
closeVowelNasal      [ 17: 21] =              ÖÑ
closeVowel           [ 17: 19] =              Ö
nasal                [ 19: 21] =               Ñ
sentenceUnmatched    [ 21: 22] =                -
word                 [ 22: 29] =                 NDÖNÎ
syllables            [ 22: 29] =                 NDÖNÎ
syllable             [ 22: 26] =                 NDÖ
consonant            [ 22: 24] =                 ND
vowelNasal           [ 24: 26] =                   Ö
closeVowelNasal      [ 24: 26] =                   Ö
closeVowel           [ 24: 26] =                   Ö
syllable             [ 26: 29] =                    NÎ
consonant            [ 26: 27] =                    N
vowelNasal           [ 27: 29] =                     Î
closeVowelNasal      [ 27: 29] =                     Î
closeVowel           [ 27: 29] =                     Î
sentenceUnmatched    [ 29: 31] =                      , 
word                 [ 31: 37] =                        ândɛ
syllables            [ 31: 37] =                        ândɛ
syllable             [ 31: 33] =                        â
vowelNasal           [ 31: 33] =                        â
closeVowelNasal      [ 31: 33] =                        â
closeVowel           [ 31: 33] =                        â
syllable             [ 33: 37] =                         ndɛ
consonant            [ 33: 35] =                         nd
vowelNasal           [ 35: 37] =                           ɛ
openVowel            [ 35: 37] =                           ɛ
sentenceUnmatched    [ 37: 38] =                             
word                 [ 38: 42] =                             bâa
syllables            [ 38: 42] =                             bâa
syllable             [ 38: 41] =                             bâ
consonant            [ 38: 39] =                             b
vowelNasal           [ 39: 41] =                              â
closeVowelNasal      [ 39: 41] =                              â
closeVowel           [ 39: 41] =                              â
syllable             [ 41: 42] =                               a
vowelNasal           [ 41: 42] =                               a
closeVowelNasal      [ 41: 42] =                               a
closeVowel           [ 41: 42] =                               a
sentenceUnmatched    [ 42: 43] =                                -
word                 [ 43: 45] =                                 mo
syllables            [ 43: 45] =                                 mo
syllable             [ 43: 45] =                                 mo
consonant            [ 43: 44] =                                 m
vowelNasal           [ 44: 45] =                                  o
closeVowelNasal      [ 44: 45] =                                  o
closeVowel           [ 44: 45] =                                  o
sentenceUnmatched    [ 45: 46] =                                   -
word                 [ 46: 52] =                                    tɛnɛ
syllables            [ 46: 52] =                                    tɛnɛ
syllable             [ 46: 49] =                                    tɛ
consonant            [ 46: 47] =                                    t
vowelNasal           [ 47: 49] =                                     ɛ
openVowel            [ 47: 49] =                                     ɛ
syllable             [ 49: 52] =                                      nɛ
consonant            [ 49: 50] =                                      n
vowelNasal           [ 50: 52] =                                       ɛ
openVowel            [ 50: 52] =                                       ɛ
sentenceUnmatched    [ 52: 54] =                                        ? 
sentence             [ 54:112] =                                          BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ!
word                 [ 54: 59] =                                          BƐ̂
syllables            [ 54: 59] =                                          BƐ̂
syllable             [ 54: 59] =                                          BƐ̂
consonant            [ 54: 55] =                                          B
vowelNasal           [ 55: 59] =                                           Ɛ̂
openVowel            [ 55: 59] =                                           Ɛ̂
sentenceUnmatched    [ 59: 60] =                                             -
word                 [ 60: 77] =                                              kɔ̈mbïtɛbɛ̂
syllables            [ 60: 77] =                                              kɔ̈mbïtɛbɛ̂
syllable             [ 60: 65] =                                              kɔ̈
consonant            [ 60: 61] =                                              k
vowelNasal           [ 61: 65] =                                               ɔ̈
openVowel            [ 61: 65] =                                               ɔ̈
syllable             [ 65: 69] =                                                 mbï
consonant            [ 65: 67] =                                                 mb
vowelNasal           [ 67: 69] =                                                   ï
closeVowelNasal      [ 67: 69] =                                                   ï
closeVowel           [ 67: 69] =                                                   ï
syllable             [ 69: 72] =                                                    tɛ
consonant            [ 69: 70] =                                                    t
vowelNasal           [ 70: 72] =                                                     ɛ
openVowel            [ 70: 72] =                                                     ɛ
syllable             [ 72: 77] =                                                      bɛ̂
consonant            [ 72: 73] =                                                      b
vowelNasal           [ 73: 77] =                                                       ɛ̂
openVowel            [ 73: 77] =                                                       ɛ̂
sentenceUnmatched    [ 77: 99] =                                                         -kɔ̈mbïtɛ難BƐ̂-
word                 [ 99:111] =                                                                        KƆ̈MBÏTƐ
syllables            [ 99:111] =                                                                        KƆ̈MBÏTƐ
syllable             [ 99:104] =                                                                        KƆ̈
consonant            [ 99:100] =                                                                        K
vowelNasal           [100:104] =                                                                         Ɔ̈
openVowel            [100:104] =                                                                         Ɔ̈
syllable             [104:108] =                                                                           MBÏ
consonant            [104:106] =                                                                           MB
vowelNasal           [106:108] =                                                                             Ï
closeVowelNasal      [106:108] =                                                                             Ï
closeVowel           [106:108] =                                                                             Ï
syllable             [108:111] =                                                                              TƐ
consonant            [108:109] =                                                                              T
vowelNasal           [109:111] =                                                                               Ɛ
openVowel            [109:111] =                                                                               Ɛ
sentenceUnmatched    [111:112] =                                                                                !
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}

func TestParseShortSangoText(t *testing.T) {
	s := `
Na yâ tî takisï na Bangî

Benoît Kongbo

«… na pöpö tî âzo otâ sô na pêko tî takisï nî sô mbênî wâlï akono gï töngana bâgara.
Tîtene mbï dutï nzönî, mbï pûsu lo kêtê…»

`
	parseMatcherForGroupName := map[string]ParseMatcher{
		"closeVowel":                 Literal,
		"closeVowelNasal":            Literal,
		"closeVowelNasalUnmatched":   Literal,
		"closeVowelUnmatched":        Literal,
		"consonant":                  Literal,
		"consonantUnmatched":         Literal,
		"nasal":                      Literal,
		"nasalUnmatched":             Literal,
		"openVowel":                  Literal,
		"openVowelUnmatched":         Literal,
		"pitch":                      Literal,
		"pitchUnmatched":             Literal,
		"sentence":                   SangoWord,
		"sentenceUnmatched":          Literal,
		"sentenceUnmatchedUnmatched": Literal,
		"shift":                      Literal,
		"shiftUnmatched":             Literal,
		"syllable":                   Literal,
		"syllableUnmatched":          Literal,
		"syllables":                  SangoSyllable,
		"syllablesUnmatched":         Literal,
		"vowelNasal":                 Literal,
		"vowelNasalUnmatched":        Literal,
		"whole":                      Sentence,
		"wholeUnmatched":             Literal,
		"word":                       Literal,
		"wordUnmatched":              Literal,
	}
	actual := Parse(parseMatcherForGroupName, s, ParseTree{Begin: 0, End: len(s), Name: "whole"})
	var actualBuilder strings.Builder = actual.Annotate(s, false, map[string]bool{
		"closeVowel":                 true,
		"closeVowelNasal":            true,
		"closeVowelNasalUnmatched":   true,
		"closeVowelUnmatched":        true,
		"consonant":                  true,
		"consonantUnmatched":         true,
		"nasal":                      true,
		"nasalUnmatched":             true,
		"openVowel":                  true,
		"openVowelUnmatched":         true,
		"pitch":                      true,
		"pitchUnmatched":             true,
		"sentence":                   true,
		"sentenceUnmatched":          true,
		"sentenceUnmatchedUnmatched": true,
		"shift":                      true,
		"shiftUnmatched":             true,
		"syllable":                   true,
		"syllableUnmatched":          true,
		"syllables":                  true,
		"syllablesUnmatched":         true,
		"vowelNasal":                 true,
		"vowelNasalUnmatched":        true,
		"whole":                      true,
		"wholeUnmatched":             true,
		"word":                       true,
		"wordUnmatched":              true,
	})
	actualStr := actualBuilder.String()
	expectStr := `
whole                [  0:208] = "\nNa yâ tî takisï na Bangî\n\nBenoît Kongbo\n\n«… na pöpö tî âzo otâ sô na pêko tî takisï nî sô mbênî wâlï akono gï töngana bâgara.\nTîtene mbï dutï nzönî, mbï pûsu lo kêtê…»\n\n"
sentence             [  0: 31] = "\nNa yâ tî takisï na Bangî\n\n"
sentenceUnmatched    [  0:  1] = "\n"
word                 [  1:  3] = "Na"
syllables            [  1:  3] = "Na"
syllable             [  1:  3] = "Na"
consonant            [  1:  2] = "N"
vowelNasal           [  2:  3] = "a"
closeVowelNasal      [  2:  3] = "a"
closeVowel           [  2:  3] = "a"
sentenceUnmatched    [  3:  4] = " "
word                 [  4:  7] = "yâ"
syllables            [  4:  7] = "yâ"
syllable             [  4:  7] = "yâ"
consonant            [  4:  5] = "y"
vowelNasal           [  5:  7] = "â"
closeVowelNasal      [  5:  7] = "â"
closeVowel           [  5:  7] = "â"
sentenceUnmatched    [  7:  8] = " "
word                 [  8: 11] = "tî"
syllables            [  8: 11] = "tî"
syllable             [  8: 11] = "tî"
consonant            [  8:  9] = "t"
vowelNasal           [  9: 11] = "î"
closeVowelNasal      [  9: 11] = "î"
closeVowel           [  9: 11] = "î"
sentenceUnmatched    [ 11: 12] = " "
word                 [ 12: 19] = "takisï"
syllables            [ 12: 19] = "takisï"
syllable             [ 12: 14] = "ta"
consonant            [ 12: 13] = "t"
vowelNasal           [ 13: 14] = "a"
closeVowelNasal      [ 13: 14] = "a"
closeVowel           [ 13: 14] = "a"
syllable             [ 14: 16] = "ki"
consonant            [ 14: 15] = "k"
vowelNasal           [ 15: 16] = "i"
closeVowelNasal      [ 15: 16] = "i"
closeVowel           [ 15: 16] = "i"
syllable             [ 16: 19] = "sï"
consonant            [ 16: 17] = "s"
vowelNasal           [ 17: 19] = "ï"
closeVowelNasal      [ 17: 19] = "ï"
closeVowel           [ 17: 19] = "ï"
sentenceUnmatched    [ 19: 20] = " "
word                 [ 20: 22] = "na"
syllables            [ 20: 22] = "na"
syllable             [ 20: 22] = "na"
consonant            [ 20: 21] = "n"
vowelNasal           [ 21: 22] = "a"
closeVowelNasal      [ 21: 22] = "a"
closeVowel           [ 21: 22] = "a"
sentenceUnmatched    [ 22: 23] = " "
word                 [ 23: 29] = "Bangî"
syllables            [ 23: 29] = "Bangî"
syllable             [ 23: 25] = "Ba"
consonant            [ 23: 24] = "B"
vowelNasal           [ 24: 25] = "a"
closeVowelNasal      [ 24: 25] = "a"
closeVowel           [ 24: 25] = "a"
syllable             [ 25: 29] = "ngî"
consonant            [ 25: 27] = "ng"
vowelNasal           [ 27: 29] = "î"
closeVowelNasal      [ 27: 29] = "î"
closeVowel           [ 27: 29] = "î"
sentenceUnmatched    [ 29: 31] = "\n\n"
sentence             [ 31: 47] = "Benoît Kongbo\n\n"
sentenceUnmatched    [ 31: 39] = "Benoît "
word                 [ 39: 45] = "Kongbo"
syllables            [ 39: 45] = "Kongbo"
syllable             [ 39: 41] = "Ko"
consonant            [ 39: 40] = "K"
vowelNasal           [ 40: 41] = "o"
closeVowelNasal      [ 40: 41] = "o"
closeVowel           [ 40: 41] = "o"
syllable             [ 41: 45] = "ngbo"
consonant            [ 41: 44] = "ngb"
vowelNasal           [ 44: 45] = "o"
closeVowelNasal      [ 44: 45] = "o"
closeVowel           [ 44: 45] = "o"
sentenceUnmatched    [ 45: 47] = "\n\n"
sentence             [ 47:153] = "«… na pöpö tî âzo otâ sô na pêko tî takisï nî sô mbênî wâlï akono gï töngana bâgara.\n"
sentenceUnmatched    [ 47: 53] = "«… "
word                 [ 53: 55] = "na"
syllables            [ 53: 55] = "na"
syllable             [ 53: 55] = "na"
consonant            [ 53: 54] = "n"
vowelNasal           [ 54: 55] = "a"
closeVowelNasal      [ 54: 55] = "a"
closeVowel           [ 54: 55] = "a"
sentenceUnmatched    [ 55: 56] = " "
word                 [ 56: 62] = "pöpö"
syllables            [ 56: 62] = "pöpö"
syllable             [ 56: 59] = "pö"
consonant            [ 56: 57] = "p"
vowelNasal           [ 57: 59] = "ö"
closeVowelNasal      [ 57: 59] = "ö"
closeVowel           [ 57: 59] = "ö"
syllable             [ 59: 62] = "pö"
consonant            [ 59: 60] = "p"
vowelNasal           [ 60: 62] = "ö"
closeVowelNasal      [ 60: 62] = "ö"
closeVowel           [ 60: 62] = "ö"
sentenceUnmatched    [ 62: 63] = " "
word                 [ 63: 66] = "tî"
syllables            [ 63: 66] = "tî"
syllable             [ 63: 66] = "tî"
consonant            [ 63: 64] = "t"
vowelNasal           [ 64: 66] = "î"
closeVowelNasal      [ 64: 66] = "î"
closeVowel           [ 64: 66] = "î"
sentenceUnmatched    [ 66: 67] = " "
word                 [ 67: 71] = "âzo"
syllables            [ 67: 71] = "âzo"
syllable             [ 67: 69] = "â"
vowelNasal           [ 67: 69] = "â"
closeVowelNasal      [ 67: 69] = "â"
closeVowel           [ 67: 69] = "â"
syllable             [ 69: 71] = "zo"
consonant            [ 69: 70] = "z"
vowelNasal           [ 70: 71] = "o"
closeVowelNasal      [ 70: 71] = "o"
closeVowel           [ 70: 71] = "o"
sentenceUnmatched    [ 71: 72] = " "
word                 [ 72: 76] = "otâ"
syllables            [ 72: 76] = "otâ"
syllable             [ 72: 73] = "o"
vowelNasal           [ 72: 73] = "o"
closeVowelNasal      [ 72: 73] = "o"
closeVowel           [ 72: 73] = "o"
syllable             [ 73: 76] = "tâ"
consonant            [ 73: 74] = "t"
vowelNasal           [ 74: 76] = "â"
closeVowelNasal      [ 74: 76] = "â"
closeVowel           [ 74: 76] = "â"
sentenceUnmatched    [ 76: 77] = " "
word                 [ 77: 80] = "sô"
syllables            [ 77: 80] = "sô"
syllable             [ 77: 80] = "sô"
consonant            [ 77: 78] = "s"
vowelNasal           [ 78: 80] = "ô"
closeVowelNasal      [ 78: 80] = "ô"
closeVowel           [ 78: 80] = "ô"
sentenceUnmatched    [ 80: 81] = " "
word                 [ 81: 83] = "na"
syllables            [ 81: 83] = "na"
syllable             [ 81: 83] = "na"
consonant            [ 81: 82] = "n"
vowelNasal           [ 82: 83] = "a"
closeVowelNasal      [ 82: 83] = "a"
closeVowel           [ 82: 83] = "a"
sentenceUnmatched    [ 83: 84] = " "
word                 [ 84: 89] = "pêko"
syllables            [ 84: 89] = "pêko"
syllable             [ 84: 87] = "pê"
consonant            [ 84: 85] = "p"
vowelNasal           [ 85: 87] = "ê"
closeVowelNasal      [ 85: 87] = "ê"
closeVowel           [ 85: 87] = "ê"
syllable             [ 87: 89] = "ko"
consonant            [ 87: 88] = "k"
vowelNasal           [ 88: 89] = "o"
closeVowelNasal      [ 88: 89] = "o"
closeVowel           [ 88: 89] = "o"
sentenceUnmatched    [ 89: 90] = " "
word                 [ 90: 93] = "tî"
syllables            [ 90: 93] = "tî"
syllable             [ 90: 93] = "tî"
consonant            [ 90: 91] = "t"
vowelNasal           [ 91: 93] = "î"
closeVowelNasal      [ 91: 93] = "î"
closeVowel           [ 91: 93] = "î"
sentenceUnmatched    [ 93: 94] = " "
word                 [ 94:101] = "takisï"
syllables            [ 94:101] = "takisï"
syllable             [ 94: 96] = "ta"
consonant            [ 94: 95] = "t"
vowelNasal           [ 95: 96] = "a"
closeVowelNasal      [ 95: 96] = "a"
closeVowel           [ 95: 96] = "a"
syllable             [ 96: 98] = "ki"
consonant            [ 96: 97] = "k"
vowelNasal           [ 97: 98] = "i"
closeVowelNasal      [ 97: 98] = "i"
closeVowel           [ 97: 98] = "i"
syllable             [ 98:101] = "sï"
consonant            [ 98: 99] = "s"
vowelNasal           [ 99:101] = "ï"
closeVowelNasal      [ 99:101] = "ï"
closeVowel           [ 99:101] = "ï"
sentenceUnmatched    [101:102] = " "
word                 [102:105] = "nî"
syllables            [102:105] = "nî"
syllable             [102:105] = "nî"
consonant            [102:103] = "n"
vowelNasal           [103:105] = "î"
closeVowelNasal      [103:105] = "î"
closeVowel           [103:105] = "î"
sentenceUnmatched    [105:106] = " "
word                 [106:109] = "sô"
syllables            [106:109] = "sô"
syllable             [106:109] = "sô"
consonant            [106:107] = "s"
vowelNasal           [107:109] = "ô"
closeVowelNasal      [107:109] = "ô"
closeVowel           [107:109] = "ô"
sentenceUnmatched    [109:110] = " "
word                 [110:117] = "mbênî"
syllables            [110:117] = "mbênî"
syllable             [110:114] = "mbê"
consonant            [110:112] = "mb"
vowelNasal           [112:114] = "ê"
closeVowelNasal      [112:114] = "ê"
closeVowel           [112:114] = "ê"
syllable             [114:117] = "nî"
consonant            [114:115] = "n"
vowelNasal           [115:117] = "î"
closeVowelNasal      [115:117] = "î"
closeVowel           [115:117] = "î"
sentenceUnmatched    [117:118] = " "
word                 [118:124] = "wâlï"
syllables            [118:124] = "wâlï"
syllable             [118:121] = "wâ"
consonant            [118:119] = "w"
vowelNasal           [119:121] = "â"
closeVowelNasal      [119:121] = "â"
closeVowel           [119:121] = "â"
syllable             [121:124] = "lï"
consonant            [121:122] = "l"
vowelNasal           [122:124] = "ï"
closeVowelNasal      [122:124] = "ï"
closeVowel           [122:124] = "ï"
sentenceUnmatched    [124:125] = " "
word                 [125:130] = "akono"
syllables            [125:130] = "akono"
syllable             [125:126] = "a"
vowelNasal           [125:126] = "a"
closeVowelNasal      [125:126] = "a"
closeVowel           [125:126] = "a"
syllable             [126:128] = "ko"
consonant            [126:127] = "k"
vowelNasal           [127:128] = "o"
closeVowelNasal      [127:128] = "o"
closeVowel           [127:128] = "o"
syllable             [128:130] = "no"
consonant            [128:129] = "n"
vowelNasal           [129:130] = "o"
closeVowelNasal      [129:130] = "o"
closeVowel           [129:130] = "o"
sentenceUnmatched    [130:131] = " "
word                 [131:134] = "gï"
syllables            [131:134] = "gï"
syllable             [131:134] = "gï"
consonant            [131:132] = "g"
vowelNasal           [132:134] = "ï"
closeVowelNasal      [132:134] = "ï"
closeVowel           [132:134] = "ï"
sentenceUnmatched    [134:135] = " "
word                 [135:143] = "töngana"
syllables            [135:143] = "töngana"
syllable             [135:138] = "tö"
consonant            [135:136] = "t"
vowelNasal           [136:138] = "ö"
closeVowelNasal      [136:138] = "ö"
closeVowel           [136:138] = "ö"
syllable             [138:141] = "nga"
consonant            [138:140] = "ng"
vowelNasal           [140:141] = "a"
closeVowelNasal      [140:141] = "a"
closeVowel           [140:141] = "a"
syllable             [141:143] = "na"
consonant            [141:142] = "n"
vowelNasal           [142:143] = "a"
closeVowelNasal      [142:143] = "a"
closeVowel           [142:143] = "a"
sentenceUnmatched    [143:144] = " "
word                 [144:151] = "bâgara"
syllables            [144:151] = "bâgara"
syllable             [144:147] = "bâ"
consonant            [144:145] = "b"
vowelNasal           [145:147] = "â"
closeVowelNasal      [145:147] = "â"
closeVowel           [145:147] = "â"
syllable             [147:149] = "ga"
consonant            [147:148] = "g"
vowelNasal           [148:149] = "a"
closeVowelNasal      [148:149] = "a"
closeVowel           [148:149] = "a"
syllable             [149:151] = "ra"
consonant            [149:150] = "r"
vowelNasal           [150:151] = "a"
closeVowelNasal      [150:151] = "a"
closeVowel           [150:151] = "a"
sentenceUnmatched    [151:153] = ".\n"
sentence             [153:208] = "Tîtene mbï dutï nzönî, mbï pûsu lo kêtê…»\n\n"
word                 [153:160] = "Tîtene"
syllables            [153:160] = "Tîtene"
syllable             [153:156] = "Tî"
consonant            [153:154] = "T"
vowelNasal           [154:156] = "î"
closeVowelNasal      [154:156] = "î"
closeVowel           [154:156] = "î"
syllable             [156:158] = "te"
consonant            [156:157] = "t"
vowelNasal           [157:158] = "e"
closeVowelNasal      [157:158] = "e"
closeVowel           [157:158] = "e"
syllable             [158:160] = "ne"
consonant            [158:159] = "n"
vowelNasal           [159:160] = "e"
closeVowelNasal      [159:160] = "e"
closeVowel           [159:160] = "e"
sentenceUnmatched    [160:161] = " "
word                 [161:165] = "mbï"
syllables            [161:165] = "mbï"
syllable             [161:165] = "mbï"
consonant            [161:163] = "mb"
vowelNasal           [163:165] = "ï"
closeVowelNasal      [163:165] = "ï"
closeVowel           [163:165] = "ï"
sentenceUnmatched    [165:166] = " "
word                 [166:171] = "dutï"
syllables            [166:171] = "dutï"
syllable             [166:168] = "du"
consonant            [166:167] = "d"
vowelNasal           [167:168] = "u"
closeVowelNasal      [167:168] = "u"
closeVowel           [167:168] = "u"
syllable             [168:171] = "tï"
consonant            [168:169] = "t"
vowelNasal           [169:171] = "ï"
closeVowelNasal      [169:171] = "ï"
closeVowel           [169:171] = "ï"
sentenceUnmatched    [171:172] = " "
word                 [172:179] = "nzönî"
syllables            [172:179] = "nzönî"
syllable             [172:176] = "nzö"
consonant            [172:174] = "nz"
vowelNasal           [174:176] = "ö"
closeVowelNasal      [174:176] = "ö"
closeVowel           [174:176] = "ö"
syllable             [176:179] = "nî"
consonant            [176:177] = "n"
vowelNasal           [177:179] = "î"
closeVowelNasal      [177:179] = "î"
closeVowel           [177:179] = "î"
sentenceUnmatched    [179:181] = ", "
word                 [181:185] = "mbï"
syllables            [181:185] = "mbï"
syllable             [181:185] = "mbï"
consonant            [181:183] = "mb"
vowelNasal           [183:185] = "ï"
closeVowelNasal      [183:185] = "ï"
closeVowel           [183:185] = "ï"
sentenceUnmatched    [185:186] = " "
word                 [186:191] = "pûsu"
syllables            [186:191] = "pûsu"
syllable             [186:189] = "pû"
consonant            [186:187] = "p"
vowelNasal           [187:189] = "û"
closeVowelNasal      [187:189] = "û"
closeVowel           [187:189] = "û"
syllable             [189:191] = "su"
consonant            [189:190] = "s"
vowelNasal           [190:191] = "u"
closeVowelNasal      [190:191] = "u"
closeVowel           [190:191] = "u"
sentenceUnmatched    [191:192] = " "
word                 [192:194] = "lo"
syllables            [192:194] = "lo"
syllable             [192:194] = "lo"
consonant            [192:193] = "l"
vowelNasal           [193:194] = "o"
closeVowelNasal      [193:194] = "o"
closeVowel           [193:194] = "o"
sentenceUnmatched    [194:195] = " "
word                 [195:201] = "kêtê"
syllables            [195:201] = "kêtê"
syllable             [195:198] = "kê"
consonant            [195:196] = "k"
vowelNasal           [196:198] = "ê"
closeVowelNasal      [196:198] = "ê"
closeVowel           [196:198] = "ê"
syllable             [198:201] = "tê"
consonant            [198:199] = "t"
vowelNasal           [199:201] = "ê"
closeVowelNasal      [199:201] = "ê"
closeVowel           [199:201] = "ê"
sentenceUnmatched    [201:208] = "…»\n\n"
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}

func TestParseLongSangoText(t *testing.T) {
	s := `
Na yâ tî takisï na Bangî

Benoît Kongbo

«… na pöpö tî âzo otâ sô na pêko tî takisï nî sô mbênî wâlï akono gï töngana bâgara.
Tîtene mbï dutï nzönî, mbï pûsu lo kêtê…»

Mbï löndö tî sïgî na pekô tî büngbïngö-terê na wayïndä tî Alliance française, na lo sô ayeke wayïndä tî Espace langue et formation, na kötä wasüngö-mbëtï Pierre Sammy sô mbï yeke gue na ndo tî lo ândö mîngi, na Gabriel, mbênî wasüngö-mbëtï ngâ na âmbâ tî mbï maseka-kôlï sô ayê kua tî süngö-mbëtï.

Mbï yeke tambûla tî gue na balabâla tî Dipandäa, ndo sô kâ mbï yeke mû takisï daä tî gue na yângâda.
Mbï yeke gbû li tî mbï na ndö tî kua sô amû na ë.
Lo hûnda na ë tî sû mbênî mbëtï sô ë lîngbi tî mû na lo na yâ tî yenga ûse.

Lo tene fadë na ë:

— Kodë nî, ayeke tî sâra mbëtï tî âla nî töngana mbênî bûku sô yângâ tî ködörö nî ayeke dutï sängö na farânzi.
Ndâ tî pialö sô, ayeke tî mû lêgë na âmbênî na pöpö tî âla sô ayeke sû mbëtï na gbemingo tî fa terê tî âla polêlê
na âzo tî Bêafrîka.

Tî taâ tënë nî, pialö sô abâa gïgî na pekô tî gbüngö-li tî vassilis, wasüngö-mbëtï tî farânzi sô taâ marä tî lo ayeke gerêke sô asû mbênî bûku na ndö tî sängö sô ïrï nî ayeke Les mots étrangers.
Na ûse gängö tî lo na sêse tî Bêafrîka, lo sâra na ë tënë tî pialö tî süngö bûku lêgëôko na sängö ngâ na farânzi.
Ë kûê ë yê daä na nî ngbanga tî sô ayeke nzönî kodë tî zîa nzara tî dïköngö-mbëtï na yâ tî âzo tî Bêafriîka, sô ahînga ngâ na âla sô ahînga nî pëpe.
Mawa nî ayeke sô anînga kêtêz mbï hön tënë tî kua tî süngö mbênî bûku na Limoges na sêse tî Farânzi sï mbï wara âsango na ndö tî pialö sô pëpe ngbii asï na lâkûi tî büngbïngö-terê sô mbï löndö tî tene tënë nî sô.

Mbï yeke gi bë tî mbï ngbanga tî sô mbï yeke sû mbênî bûku sô mbï yê tî tokua na da lëkërëngö nî na sêse tî Gameroun.
Mbï zîa bê tî mbï daä mîngi sï mbeto ayeke gbû mbï tî wara lêgë pëpe tî sû mbëtï sô Alliance ahûnda.
Mbï yê tî hûnzi bûku sô hîo sï mbï yê tî sâra mbênî yê ndê pëpe.
Mbï yê tîtene ahûnda mbï hündängö sï mbï sû mbëtï pëpe.
Mbï yeke sâra nî na ngoi sô mbï yeke na âsango tî vunga.
Ânzönî tënë sô mbï yeke na nî tîtene na ndö tî ködörö tî mbï ayeke na yâ tî kôzo bûku tî mbï awe.
Tënë wa ânde sï mbï zîa na ndö nî?

Kötä pîpî tî mbênî takisï sï azî mbï na yâ tî gbüngö-li sô.
Ayeke gï na ngoi sô sï mbï hînga mbï sï na yângâ tî kötä lêgë awe.

Wagbötöngö-takisï nî atene na kötä gô:

— Ben-Zvi-Miskine!
Mo gue?

Mbï sâra mabôko tî mbï na lo tîtene lo lütï sï mbï mû lörö na mbâgë tî takisï nî.

Mbï tene na mâsa nî kôzo kûê tîtene mbï wara parâba pëpe:

— Mâsa, mbï yeke zûu na dambëtï tî Ben-Zvi.
Mbï yeke na pâta balëûse na okü.

Lo dë köngö mo bâa mo tene kâ fadë mbï fa na lo wüngö tî nginza nî pëpe:

— Me mo kö daä ma!

Mbï kö na yâ tî takisï nî na lörö.

— Awe, mâsa.
Ë gue.

Mbï wara lêgë tî dutï nzönî pëpe ngbanga tî sô na pöpö tî âzo otâ sô na pêko tî takisï nî sô, mbênî wâlï akono gï töngana bâgara.
Tîtene mbï dutï nzönî, mbï pûsu lo kêtê me wâlï sô asûku ngbângbâ, lo pete mbï na sïönî bê na terê tî yângâ tî takisï ni, sï lo tene:

— Mâsa, takisï tî mo nî ayeke kêtê töngasô ngbanga tî nye?
Bâa yê, gï na lëkërëngö têre tî mbï tîtene maseka-kôlï sô gerê tî lô akpa gï tî kôlô sô adutï nzönî, yâ tî terê tî mbï kûê ayeke so.

Wagbötöngö-takisï nî atene na ngîâ na yângâ tî lo:

— Mo hînga fadë, mamâ, mo fûta ngêrë tî ndo ûse tîtene mo bâa pâsi pëpe!

Wâlï sô akpa gï doli sô atï na wâ agbï, lo tene:

— Âla mä tënë tî kpändä sô!
Lo hînga mbï äpe, mbï hînga lo äpe.
Lo yeke tene tënë na mbï mo bâa mo tene mbï yeke mbênî kötä konô sô asï yä tî takisï nî kûê.
Mä mbï sï, môlengê.
Mbï lîngbi tî dutï töngana mamâ tî mo.
Gï na särängö terê tî âla töngasô, âla âmaseka, sï âla ke kûi hîo sô.

Ndo adë ngii sï pakara sô na gatï tî mamâ sô azî yângâ:

— Mbï yeke hûnda terê tî mbï, mbï tene, âla wara sïönî mayëre tî zîa âzo osïö na pekö sô na ndo wa?
Takisï tî mo nî ayeke kêtê ahön kopo tî âsusu.

Wagbötöngö-kutukutu nî azî yângâ pëpe tënë tî bîâ sô lo yeke hë na gbegô.

Kôlï sô na terê tî wâlï sô na kotï tî mbï akîri atene:

— Na yâ tî âmbênî ködörö, na tâpandë tî Farânzi sô mbï ke gue daä töngana lêgë tî yäkä tî mbï, âtakisï ayeke mû gï zo ôko.
Na ndönî, âla ayeke zîa lo na ndo sô lo yeke gue daä.

Gï sô lo tene töngasô sô âkperekpere nî alöndö ngangü ahön ndö tî îâ sô ayeke toto ngâ na kötä köngö tî takisï ni.

Masia sô adutï na pöpö tî Wagbötöngö-takisï na mbênî maseka sô amo li tî lo ndâli tî kötä kä sô ayeke daä atene:

— Na mbâgë, tënë tî mo nî ayeke na lêgë nî, me na mbâgë ngâ, alingbi mo hîinga ë dutï nzönî na yâ tî takisï ahön bus.

Môlengê-kôlï sô na kä na li tî lo sô atene na môlengê-wâlï nî:

— Tâtënë laä mo tene sô, îtä.
Mbï sô mo bâa mbï sô, mbï yeke mû bus lâôko pëpe.
Mbï yeke zo sô ayê dütïngö kpô mîngi.
Andâa, na yâ tî bus, zo kûê ayeke tene tënë, ayeke dë köngö, ayeke hë ngîâ na yâ tî mê tî âzo, ngâ ayeke sösö na mbänä tî bê tî lo.

Kôlï tî hümä sô adutï na mbâgë kâ mo bâa mo tene lo löndö na ködörö wandê sô atene:

—Tî mbï, mbï kpë bus ngbanga tî sô mbï yê fünngö tî terê tî âzo pëpe.
Mbï lîngbi tî kânga bê tî mbï na ndö tî fünngö gbe tî mabôko tî âmötöröbôi pëpe, wala ngâ na gbikï tî âla sô ayeke buba bongö tî mbï pëpe.

Ka mbï, âla yeke hûnda terê tî âla, mbï tene tî mbï nye?
Mbï zî yângâ ôko pëpe.
Sô fadë mbï kö na yâ tî takisï sô, sösö ahon ndö tî mbï mîngi.
Mbï yê tî sösö na ngangü tî terê tî mbï kûê sï pupu aya mbï.
Na Bangî, mo sösö, ayeke hë mo ngbii asï na lâ tî kûâ tî mo.
Ka mbï sösö nî sï sasa asïgîgî agbû bongö tî mbï, fadë mbï yeke sâra töngana nye?
Alîngbi mbï dutï kpô tî pete sösö sô ayê gï tî sungba na ngangü sô.

Wagbötöngö-takisï nî azî gô tî fonôno nî ngangü.
Töngana mbênî yê tî ndimä, sösö sô mbï kânga fadë gbä sô asïgî gï yeke bâa mo tene mozoko nî ayê däa tîtene lo sïgî.

Wâlï sô na terê tî mbï sï akpa gï kötä bäbolo sô atene na kötä gô:

— Hêe, mâsa, mo pêe tî kîri na gô tî fonôno-nî kêtê äpe?
Ë yeke na yâ tî nganda äpe, ë yeke na yâ tî takisï laâ!

Âdëköngö alöndö mbâgë na mbâgë:

— Nye laâ afün töngasô, mâsa?

— Mo yö fünngö nyama na yâ tî takisï nî wala nye?

— Kpändä wa laâ asösö sô?

— Mozoko sô ayeke toto ngangü mîngi sô laâ amû lêgë na zo sô tî sösö sô.

— Sösö nî afün gï töngana kûâ tî mbo, dê!

— Alîngbi lo gue abâa yâ tî purû tî lo sï.

Kötä wâlï sô na terê tî mbï sï amû ndo kûê sô atene:

— Mo bâa sô, ayeke gï mâsa laâ asösö sô.

Me wagbötöngö-takisï nî akîri tënë na wâlï nî atene:

— Tî mbï ayeke gï mo, pûpûlenge waâ!
Mbï hînga ka sôsô tî âwâlï sô ngbondä tî âla akono mîngi sô kûê.

Lê tî wâlï sô abe na ngonzo sï lo hûnda tîtene mâsa nî azîa lo na sêse fadëfadësô.
Ë yeke na barâsi balëokü na sängbï-lêgë tî Ben-Zvi-Sica.
Takisï nî alütï na dëmängö terê tî lo.
Wâlï nî asïgîgiî na kugbë tî sâki ôko lo mû na mâsa nî.
voazîne tî mbï nî, sô ayeke voazîne tî mbï mbênî pëpe tëtî sô lo zûu na sêse awe sô, ake na kürü gô:

— Mo dîko tanga tî nginza sô mo kîri na mbï sô laâ?

— Mbï kîri na mo ngbangbo na balë mbârâmbârâ.
Ka tënë nî ayeke na ndo wa, mamâ tî mbï?

— Mbï tene na mo fadë mbï tene mbï yeke gï na balëûse na okü, äpe?

— Mbï yeke mû ka balëûse na okü töngana mo yeke gï na balëûse na okü.

Mâsa nî azî gerê, takisï tî lo nî ahön yongôro na terê tî wâlï tî gïngö tënë sô sï lo tene nï hînga li tî nï na âzo tî Bêafrîka pëpe.
Ngêrë tî takisï ayeke balëotâ.
Ayeke na nzönî bê laâ sï lo yeke yô âzo na balëûse na okü, sô ayeke ngêrë tî köngö na yâ tî bus sô.

Töngana kutukutu nî asï na dambëtï tî Ben-Zvi awe, mbï zûu na sêse ngâ.
Na mbâgë tî lêgë na dawäa tî mbênî dapïkängö-singa, Sylvianne ayeke kü mbï töngana tî sô ë mä terê daä.
Mbï gue na terê tî lo sï lêgëôko, ë mû lêgë sô ague na Centre Baobab tî kîri na yângâda.

Lo hûnda mbï tî hînga:

— Büngbïngö-terê nî atambûla na lêgë nî?

— Iin.
Ahûnda na ë tî bûngbi tî sâra mbênî bûku na farânzi ngâ na sängö.

— Mo hînga yê sô mo yê tî sîra mbëtï na ndö nî sô awe?

— Mbï yê tî vunga sango tî yê sô asï fadë na yâ tî takisï sô.

— Nye laâ asï nye?

—Âzo tî yâ tî takisï nî adîko tënë töna âkûkulû sï mbï sösö tîtene âla dutï kpo.

Sylvianne atï na kötä ngîâ sï nzara tî mbênî sösö agbû mbï.
Me mbï kânga nî.

— Ahînga atene mo laâ mo sösö sô?

— Mbï sösö nî na ngoi sô mâsa nî azî gô tî fonôno nî ngangü.

— Tî taâ tënë nî, atumba fünngö nî äpe?

— Âhînga fadë atene mbï laâ mbï sösö sô, ka kamënë ayeke gbû mbï.

— Sösö ayeke yê tî kamënë pëpe!
Sösö afa atene zo ayeke na sêni.
Atâa töngana âsûru yâ tî zo awe sô, ayeke kü tîtene zo nî asösö awe, sï lo lîngbi tî te yê.
Âmunzû ayeke dema pëpe töngana mbâ tî âla asösö.

Na dabüzë tî Amadi, mbï vo gbâ tî mânga sô ïrï nî Menthol sô sï mbï zîa wâ na terê tî mbênî ôko.
Pialö tî süngö bûku na farânzi-sängö sô akîri na li tî mbï.
Tî särängö mbëtï na farânzi sï mbï kîri mbï lekere nî na sängö sô aso li tî mbï mîngi.
Töngana nye ânde sï mbï sâra mbëtï na yângâ tî ködörö tî mbï sô mbï hînga tî tene nzönî pëpe sô?
Âfa nî fadë na mbî na dambëtï, mbï lîngbi tî wara kpälë tî särängö mbëtï na sängö pëpe.
Me mawa nî ayeke sô!
Gï na fängö na mbï tî tene ngâ na tî sâra mbëtï nzönî na yângâ tî ködörö tî âzo sô amû na ë Dipandäa, âzo tî tënëngö kötä farânzi tî ködörö tî ë asâra sï ë hînga tënëngö sängö nzönî pëpe.
Âmbupa tî farânzi asï sïngö na yâ tî sängö sô ë yeke tene.
Mbï yeke tö ndâ tî sâra mbëtï na marä sängö sô, ngbanga tî sô nzönî ködë tî manda yângâ tî ködörö, ayeke gï tî tene nî.

Mo yê tî bâa, sîngamê tî mbï tî mabôko atoto.

— Alô?

Mbï ma atene:

— Töngana nye, petï?

— Zo wa laâ?

— Mo girisa gô tî mbï awe?
Mbï Yvon laâ äpe.
Mo yeke na ndo wa?

— Mbï yeke gue na yângâda laâ.

— Mo pêe tî gä hîo na ndo sô ë ike wara ka terê daä sô?

— Mbï pêe, me na yâ tî ngbonga ôko ngbanga tî sô mbï yê tî sû mbênî mbëtï na farânzi ngâ na sängö.

— Alîngbi tî kü kêkerêke äpe?

— Mbï yê tî tö ndâ nî gï fafadësô, töngasô äpe sô, mbï lîngbi tî sû nî ânde äpe.

— Mbï lîngbi tî kü gï töngana mo yê daä tî sâra tënë tî mbï na yâ tî mbëtï tî mo nî sï.

Mbï tene na lo na ngîâ na yângâ tî mbï.
`
	parseMatcherForGroupName := map[string]ParseMatcher{
		"sentence": Literal,
		"whole":    Sentence,
	}
	actual := Parse(parseMatcherForGroupName, s, ParseTree{Begin: 0, End: len(s), Name: "whole"})
	var actualBuilder strings.Builder = actual.Annotate(s, false, map[string]bool{
		"sentence": true,
	})
	actualStr := actualBuilder.String()
	expectStr := `
sentence             [  0: 31] = "\nNa yâ tî takisï na Bangî\n\n"
sentence             [ 31: 47] = "Benoît Kongbo\n\n"
sentence             [ 47:153] = "«… na pöpö tî âzo otâ sô na pêko tî takisï nî sô mbênî wâlï akono gï töngana bâgara.\n"
sentence             [153:208] = "Tîtene mbï dutï nzönî, mbï pûsu lo kêtê…»\n\n"
sentence             [208:559] = "Mbï löndö tî sïgî na pekô tî büngbïngö-terê na wayïndä tî Alliance française, na lo sô ayeke wayïndä tî Espace langue et formation, na kötä wasüngö-mbëtï Pierre Sammy sô mbï yeke gue na ndo tî lo ândö mîngi, na Gabriel, mbênî wasüngö-mbëtï ngâ na âmbâ tî mbï maseka-kôlï sô ayê kua tî süngö-mbëtï.\n\n"
sentence             [559:675] = "Mbï yeke tambûla tî gue na balabâla tî Dipandäa, ndo sô kâ mbï yeke mû takisï daä tî gue na yângâda.\n"
sentence             [675:734] = "Mbï yeke gbû li tî mbï na ndö tî kua sô amû na ë.\n"
sentence             [734:827] = "Lo hûnda na ë tî sû mbênî mbëtï sô ë lîngbi tî mû na lo na yâ tî yenga ûse.\n\n"
sentence             [827:849] = "Lo tene fadë na ë:\n\n"
sentence             [849:987] = "— Kodë nî, ayeke tî sâra mbëtï tî âla nî töngana mbênî bûku sô yângâ tî ködörö nî ayeke dutï sängö na farânzi.\n"
sentence             [987:1151] = "Ndâ tî pialö sô, ayeke tî mû lêgë na âmbênî na pöpö tî âla sô ayeke sû mbëtï na gbemingo tî fa terê tî âla polêlê\nna âzo tî Bêafrîka.\n\n"
sentence             [1151:1386] = "Tî taâ tënë nî, pialö sô abâa gïgî na pekô tî gbüngö-li tî vassilis, wasüngö-mbëtï tî farânzi sô taâ marä tî lo ayeke gerêke sô asû mbênî bûku na ndö tî sängö sô ïrï nî ayeke Les mots étrangers.\n"
sentence             [1386:1525] = "Na ûse gängö tî lo na sêse tî Bêafrîka, lo sâra na ë tënë tî pialö tî süngö bûku lêgëôko na sängö ngâ na farânzi.\n"
sentence             [1525:1708] = "Ë kûê ë yê daä na nî ngbanga tî sô ayeke nzönî kodë tî zîa nzara tî dïköngö-mbëtï na yâ tî âzo tî Bêafriîka, sô ahînga ngâ na âla sô ahînga nî pëpe.\n"
sentence             [1708:1966] = "Mawa nî ayeke sô anînga kêtêz mbï hön tënë tî kua tî süngö mbênî bûku na Limoges na sêse tî Farânzi sï mbï wara âsango na ndö tî pialö sô pëpe ngbii asï na lâkûi tî büngbïngö-terê sô mbï löndö tî tene tënë nî sô.\n\n"
sentence             [1966:2106] = "Mbï yeke gi bë tî mbï ngbanga tî sô mbï yeke sû mbênî bûku sô mbï yê tî tokua na da lëkërëngö nî na sêse tî Gameroun.\n"
sentence             [2106:2227] = "Mbï zîa bê tî mbï daä mîngi sï mbeto ayeke gbû mbï tî wara lêgë pëpe tî sû mbëtï sô Alliance ahûnda.\n"
sentence             [2227:2309] = "Mbï yê tî hûnzi bûku sô hîo sï mbï yê tî sâra mbênî yê ndê pëpe.\n"
sentence             [2309:2379] = "Mbï yê tîtene ahûnda mbï hündängö sï mbï sû mbëtï pëpe.\n"
sentence             [2379:2443] = "Mbï yeke sâra nî na ngoi sô mbï yeke na âsango tî vunga.\n"
sentence             [2443:2563] = "Ânzönî tënë sô mbï yeke na nî tîtene na ndö tî ködörö tî mbï ayeke na yâ tî kôzo bûku tî mbï awe.\n"
sentence             [2563:2607] = "Tënë wa ânde sï mbï zîa na ndö nî?\n\n"
sentence             [2607:2683] = "Kötä pîpî tî mbênî takisï sï azî mbï na yâ tî gbüngö-li sô.\n"
sentence             [2683:2765] = "Ayeke gï na ngoi sô sï mbï hînga mbï sï na yângâ tî kötä lêgë awe.\n\n"
sentence             [2765:2813] = "Wagbötöngö-takisï nî atene na kötä gô:\n\n"
sentence             [2813:2834] = "— Ben-Zvi-Miskine!\n"
sentence             [2834:2843] = "Mo gue?\n\n"
sentence             [2843:2944] = "Mbï sâra mabôko tî mbï na lo tîtene lo lütï sï mbï mû lörö na mbâgë tî takisï nî.\n\n"
sentence             [2944:3013] = "Mbï tene na mâsa nî kôzo kûê tîtene mbï wara parâba pëpe:\n\n"
sentence             [3013:3065] = "— Mâsa, mbï yeke zûu na dambëtï tî Ben-Zvi.\n"
sentence             [3065:3104] = "Mbï yeke na pâta balëûse na okü.\n\n"
sentence             [3104:3190] = "Lo dë köngö mo bâa mo tene kâ fadë mbï fa na lo wüngö tî nginza nî pëpe:\n\n"
sentence             [3190:3214] = "— Me mo kö daä ma!\n\n"
sentence             [3214:3258] = "Mbï kö na yâ tî takisï nî na lörö.\n\n"
sentence             [3258:3274] = "— Awe, mâsa.\n"
sentence             [3274:3283] = "Ë gue.\n\n"
sentence             [3283:3441] = "Mbï wara lêgë tî dutï nzönî pëpe ngbanga tî sô na pöpö tî âzo otâ sô na pêko tî takisï nî sô, mbênî wâlï akono gï töngana bâgara.\n"
sentence             [3441:3603] = "Tîtene mbï dutï nzönî, mbï pûsu lo kêtê me wâlï sô asûku ngbângbâ, lo pete mbï na sïönî bê na terê tî yângâ tî takisï ni, sï lo tene:\n\n"
sentence             [3603:3673] = "— Mâsa, takisï tî mo nî ayeke kêtê töngasô ngbanga tî nye?\n"
sentence             [3673:3839] = "Bâa yê, gï na lëkërëngö têre tî mbï tîtene maseka-kôlï sô gerê tî lô akpa gï tî kôlô sô adutï nzönî, yâ tî terê tî mbï kûê ayeke so.\n\n"
sentence             [3839:3901] = "Wagbötöngö-takisï nî atene na ngîâ na yângâ tî lo:\n\n"
sentence             [3901:3989] = "— Mo hînga fadë, mamâ, mo fûta ngêrë tî ndo ûse tîtene mo bâa pâsi pëpe!\n\n"
sentence             [3989:4047] = "Wâlï sô akpa gï doli sô atï na wâ agbï, lo tene:\n\n"
sentence             [4047:4085] = "— Âla mä tënë tî kpändä sô!\n"
sentence             [4085:4127] = "Lo hînga mbï äpe, mbï hînga lo äpe.\n"
sentence             [4127:4238] = "Lo yeke tene tënë na mbï mo bâa mo tene mbï yeke mbênî kötä konô sô asï yä tî takisï nî kûê.\n"
sentence             [4238:4263] = "Mä mbï sï, môlengê.\n"
sentence             [4263:4309] = "Mbï lîngbi tî dutï töngana mamâ tî mo.\n"
sentence             [4309:4396] = "Gï na särängö terê tî âla töngasô, âla âmaseka, sï âla ke kûi hîo sô.\n\n"
sentence             [4396:4463] = "Ndo adë ngii sï pakara sô na gatï tî mamâ sô azî yângâ:\n\n"
sentence             [4463:4583] = "— Mbï yeke hûnda terê tî mbï, mbï tene, âla wara sïönî mayëre tî zîa âzo osïö na pekö sô na ndo wa?\n"
sentence             [4583:4639] = "Takisï tî mo nî ayeke kêtê ahön kopo tî âsusu.\n\n"
sentence             [4639:4730] = "Wagbötöngö-kutukutu nî azî yângâ pëpe tënë tî bîâ sô lo yeke hë na gbegô.\n\n"
sentence             [4730:4798] = "Kôlï sô na terê tî wâlï sô na kotï tî mbï akîri atene:\n\n"
sentence             [4798:4951] = "— Na yâ tî âmbênî ködörö, na tâpandë tî Farânzi sô mbï ke gue daä töngana lêgë tî yäkä tî mbï, âtakisï ayeke mû gï zo ôko.\n"
sentence             [4951:5012] = "Na ndönî, âla ayeke zîa lo na ndo sô lo yeke gue daä.\n\n"
sentence             [5012:5151] = "Gï sô lo tene töngasô sô âkperekpere nî alöndö ngangü ahön ndö tî îâ sô ayeke toto ngâ na kötä köngö tî takisï ni.\n\n"
sentence             [5151:5285] = "Masia sô adutï na pöpö tî Wagbötöngö-takisï na mbênî maseka sô amo li tî lo ndâli tî kötä kä sô ayeke daä atene:\n\n"
sentence             [5285:5427] = "— Na mbâgë, tënë tî mo nî ayeke na lêgë nî, me na mbâgë ngâ, alingbi mo hîinga ë dutï nzönî na yâ tî takisï ahön bus.\n\n"
sentence             [5427:5504] = "Môlengê-kôlï sô na kä na li tî lo sô atene na môlengê-wâlï nî:\n\n"
sentence             [5504:5543] = "— Tâtënë laä mo tene sô, îtä.\n"
sentence             [5543:5603] = "Mbï sô mo bâa mbï sô, mbï yeke mû bus lâôko pëpe.\n"
sentence             [5603:5649] = "Mbï yeke zo sô ayê dütïngö kpô mîngi.\n"
sentence             [5649:5808] = "Andâa, na yâ tî bus, zo kûê ayeke tene tënë, ayeke dë köngö, ayeke hë ngîâ na yâ tî mê tî âzo, ngâ ayeke sösö na mbänä tî bê tî lo.\n\n"
sentence             [5808:5911] = "Kôlï tî hümä sô adutï na mbâgë kâ mo bâa mo tene lo löndö na ködörö wandê sô atene:\n\n"
sentence             [5911:5998] = "—Tî mbï, mbï kpë bus ngbanga tî sô mbï yê fünngö tî terê tî âzo pëpe.\n"
sentence             [5998:6167] = "Mbï lîngbi tî kânga bê tî mbï na ndö tî fünngö gbe tî mabôko tî âmötöröbôi pëpe, wala ngâ na gbikï tî âla sô ayeke buba bongö tî mbï pëpe.\n\n"
sentence             [6167:6233] = "Ka mbï, âla yeke hûnda terê tî âla, mbï tene tî mbï nye?\n"
sentence             [6233:6262] = "Mbï zî yângâ ôko pëpe.\n"
sentence             [6262:6339] = "Sô fadë mbï kö na yâ tî takisï sô, sösö ahon ndö tî mbï mîngi.\n"
sentence             [6339:6414] = "Mbï yê tî sösö na ngangü tî terê tî mbï kûê sï pupu aya mbï.\n"
sentence             [6414:6485] = "Na Bangî, mo sösö, ayeke hë mo ngbii asï na lâ tî kûâ tî mo.\n"
sentence             [6485:6583] = "Ka mbï sösö nî sï sasa asïgîgî agbû bongö tî mbï, fadë mbï yeke sâra töngana nye?\n"
sentence             [6583:6665] = "Alîngbi mbï dutï kpô tî pete sösö sô ayê gï tî sungba na ngangü sô.\n\n"
sentence             [6665:6725] = "Wagbötöngö-takisï nî azî gô tî fonôno nî ngangü.\n"
sentence             [6725:6867] = "Töngana mbênî yê tî ndimä, sösö sô mbï kânga fadë gbä sô asïgî gï yeke bâa mo tene mozoko nî ayê däa tîtene lo sïgî.\n\n"
sentence             [6867:6950] = "Wâlï sô na terê tî mbï sï akpa gï kötä bäbolo sô atene na kötä gô:\n\n"
sentence             [6950:7021] = "— Hêe, mâsa, mo pêe tî kîri na gô tî fonôno-nî kêtê äpe?\n"
sentence             [7021:7087] = "Ë yeke na yâ tî nganda äpe, ë yeke na yâ tî takisï laâ!\n\n"
sentence             [7087:7130] = "Âdëköngö alöndö mbâgë na mbâgë:\n\n"
sentence             [7130:7168] = "— Nye laâ afün töngasô, mâsa?\n\n"
sentence             [7168:7228] = "— Mo yö fünngö nyama na yâ tî takisï nî wala nye?\n\n"
sentence             [7228:7263] = "— Kpändä wa laâ asösö sô?\n\n"
sentence             [7263:7352] = "— Mozoko sô ayeke toto ngangü mîngi sô laâ amû lêgë na zo sô tî sösö sô.\n\n"
sentence             [7352:7407] = "— Sösö nî afün gï töngana kûâ tî mbo, dê!\n\n"
sentence             [7407:7460] = "— Alîngbi lo gue abâa yâ tî purû tî lo sï.\n\n"
sentence             [7460:7527] = "Kötä wâlï sô na terê tî mbï sï amû ndo kûê sô atene:\n\n"
sentence             [7527:7579] = "— Mo bâa sô, ayeke gï mâsa laâ asösö sô.\n\n"
sentence             [7579:7644] = "Me wagbötöngö-takisï nî akîri tënë na wâlï nî atene:\n\n"
sentence             [7644:7689] = "— Tî mbï ayeke gï mo, pûpûlenge waâ!\n"
sentence             [7689:7771] = "Mbï hînga ka sôsô tî âwâlï sô ngbondä tî âla akono mîngi sô kûê.\n\n"
sentence             [7771:7869] = "Lê tî wâlï sô abe na ngonzo sï lo hûnda tîtene mâsa nî azîa lo na sêse fadëfadësô.\n"
sentence             [7869:7935] = "Ë yeke na barâsi balëokü na sängbï-lêgë tî Ben-Zvi-Sica.\n"
sentence             [7935:7983] = "Takisï nî alütï na dëmängö terê tî lo.\n"
sentence             [7983:8052] = "Wâlï nî asïgîgiî na kugbë tî sâki ôko lo mû na mâsa nî.\n"
sentence             [8052:8174] = "voazîne tî mbï nî, sô ayeke voazîne tî mbï mbênî pëpe tëtî sô lo zûu na sêse awe sô, ake na kürü gô:\n\n"
sentence             [8174:8236] = "— Mo dîko tanga tî nginza sô mo kîri na mbï sô laâ?\n\n"
sentence             [8236:8291] = "— Mbï kîri na mo ngbangbo na balë mbârâmbârâ.\n"
sentence             [8291:8339] = "Ka tënë nî ayeke na ndo wa, mamâ tî mbï?\n\n"
sentence             [8339:8418] = "— Mbï tene na mo fadë mbï tene mbï yeke gï na balëûse na okü, äpe?\n\n"
sentence             [8418:8501] = "— Mbï yeke mû ka balëûse na okü töngana mo yeke gï na balëûse na okü.\n\n"
sentence             [8501:8664] = "Mâsa nî azî gerê, takisï tî lo nî ahön yongôro na terê tî wâlï tî gïngö tënë sô sï lo tene nï hînga li tî nï na âzo tî Bêafrîka pëpe.\n"
sentence             [8664:8701] = "Ngêrë tî takisï ayeke balëotâ.\n"
sentence             [8701:8821] = "Ayeke na nzönî bê laâ sï lo yeke yô âzo na balëûse na okü, sô ayeke ngêrë tî köngö na yâ tî bus sô.\n\n"
sentence             [8821:8903] = "Töngana kutukutu nî asï na dambëtï tî Ben-Zvi awe, mbï zûu na sêse ngâ.\n"
sentence             [8903:9028] = "Na mbâgë tî lêgë na dawäa tî mbênî dapïkängö-singa, Sylvianne ayeke kü mbï töngana tî sô ë mä terê daä.\n"
sentence             [9028:9134] = "Mbï gue na terê tî lo sï lêgëôko, ë mû lêgë sô ague na Centre Baobab tî kîri na yângâda.\n\n"
sentence             [9134:9162] = "Lo hûnda mbï tî hînga:\n\n"
sentence             [9162:9215] = "— Büngbïngö-terê nî atambûla na lêgë nî?\n\n"
sentence             [9215:9224] = "— Iin.\n"
sentence             [9224:9304] = "Ahûnda na ë tî bûngbi tî sâra mbênî bûku na farânzi ngâ na sängö.\n\n"
sentence             [9304:9373] = "— Mo hînga yê sô mo yê tî sîra mbëtï na ndö nî sô awe?\n\n"
sentence             [9373:9450] = "— Mbï yê tî vunga sango tî yê sô asï fadë na yâ tî takisï sô.\n\n"
sentence             [9450:9474] = "— Nye laâ asï nye?\n\n"
sentence             [9474:9578] = "—Âzo tî yâ tî takisï nî adîko tënë töna âkûkulû sï mbï sösö tîtene âla dutï kpo.\n\n"
sentence             [9578:9651] = "Sylvianne atï na kötä ngîâ sï nzara tî mbênî sösö agbû mbï.\n"
sentence             [9651:9672] = "Me mbï kânga nî.\n\n"
sentence             [9672:9714] = "— Ahînga atene mo laâ mo sösö sô?\n\n"
sentence             [9714:9791] = "— Mbï sösö nî na ngoi sô mâsa nî azî gô tî fonôno nî ngangü.\n\n"
sentence             [9791:9843] = "— Tî taâ tënë nî, atumba fünngö nî äpe?\n\n"
sentence             [9843:9925] = "— Âhînga fadë atene mbï laâ mbï sösö sô, ka kamënë ayeke gbû mbï.\n\n"
sentence             [9925:9966] = "— Sösö ayeke yê tî kamënë pëpe!\n"
sentence             [9966:10002] = "Sösö afa atene zo ayeke na sêni.\n"
sentence             [10002:10110] = "Atâa töngana âsûru yâ tî zo awe sô, ayeke kü tîtene zo nî asösö awe, sï lo lîngbi tî te yê.\n"
sentence             [10110:10169] = "Âmunzû ayeke dema pëpe töngana mbâ tî âla asösö.\n\n"
sentence             [10169:10287] = "Na dabüzë tî Amadi, mbï vo gbâ tî mânga sô ïrï nî Menthol sô sï mbï zîa wâ na terê tî mbênî ôko.\n"
sentence             [10287:10359] = "Pialö tî süngö bûku na farânzi-sängö sô akîri na li tî mbï.\n"
sentence             [10359:10464] = "Tî särängö mbëtï na farânzi sï mbï kîri mbï lekere nî na sängö sô aso li tî mbï mîngi.\n"
sentence             [10464:10584] = "Töngana nye ânde sï mbï sâra mbëtï na yângâ tî ködörö tî mbï sô mbï hînga tî tene nzönî pëpe sô?\n"
sentence             [10584:10692] = "Âfa nî fadë na mbî na dambëtï, mbï lîngbi tî wara kpälë tî särängö mbëtï na sängö pëpe.\n"
sentence             [10692:10715] = "Me mawa nî ayeke sô!\n"
sentence             [10715:10953] = "Gï na fängö na mbï tî tene ngâ na tî sâra mbëtï nzönî na yângâ tî ködörö tî âzo sô amû na ë Dipandäa, âzo tî tënëngö kötä farânzi tî ködörö tî ë asâra sï ë hînga tënëngö sängö nzönî pëpe.\n"
sentence             [10953:11024] = "Âmbupa tî farânzi asï sïngö na yâ tî sängö sô ë yeke tene.\n"
sentence             [11024:11172] = "Mbï yeke tö ndâ tî sâra mbëtï na marä sängö sô, ngbanga tî sô nzönî ködë tî manda yângâ tî ködörö, ayeke gï tî tene nî.\n\n"
sentence             [11172:11228] = "Mo yê tî bâa, sîngamê tî mbï tî mabôko atoto.\n\n"
sentence             [11228:11239] = "— Alô?\n\n"
sentence             [11239:11255] = "Mbï ma atene:\n\n"
sentence             [11255:11281] = "— Töngana nye, petï?\n\n"
sentence             [11281:11298] = "— Zo wa laâ?\n\n"
sentence             [11298:11330] = "— Mo girisa gô tî mbï awe?\n"
sentence             [11330:11351] = "Mbï Yvon laâ äpe.\n"
sentence             [11351:11371] = "Mo yeke na ndo wa?\n\n"
sentence             [11371:11409] = "— Mbï yeke gue na yângâda laâ.\n\n"
sentence             [11409:11477] = "— Mo pêe tî gä hîo na ndo sô ë ike wara ka terê daä sô?\n\n"
sentence             [11477:11598] = "— Mbï pêe, me na yâ tî ngbonga ôko ngbanga tî sô mbï yê tî sû mbênî mbëtï na farânzi ngâ na sängö.\n\n"
sentence             [11598:11637] = "— Alîngbi tî kü kêkerêke äpe?\n\n"
sentence             [11637:11741] = "— Mbï yê tî tö ndâ nî gï fafadësô, töngasô äpe sô, mbï lîngbi tî sû nî ânde äpe.\n\n"
sentence             [11741:11853] = "— Mbï lîngbi tî kü gï töngana mo yê daä tî sâra tënë tî mbï na yâ tî mbëtï tî mo nî sï.\n\n"
sentence             [11853:11900] = "Mbï tene na lo na ngîâ na yângâ tî mbï.\n"
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}
