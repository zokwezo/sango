package parse

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set global flags for all tests in this package
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.SetOutput(io.Discard)

	// Run the test suite
	os.Exit(m.Run())
}

func TestParseCanonical(t *testing.T) {
	s := ` § ~ha_HO:-Do:ni^. 日本語は難しい!`
	parseMatchers := []ParseMatcher{SentenceFinal, CanonicalWord, CanonicalSyllable}
	actual := Parse(parseMatchers, s, ParseTree{Begin: 0, End: len(s)})
	actualStr := actual.Annotate(s, true, nil)
	expectStr := `
whole                [  0: 42] =  § ~ha_HO:-Do:ni^. 日本語は難しい!
sentence             [  0: 20] =  § ~ha_HO:-Do:ni^. 
unmatched            [  0:  4] =  § 
word                 [  4: 11] =    ~ha_HO:
syllable             [  4:  8] =    ~ha_
shift                [  4:  5] =    ~
consonant            [  5:  6] =     h
vowelNasal           [  6:  7] =      a
pitch                [  7:  8] =       _
syllable             [  8: 11] =        HO:
consonant            [  8:  9] =        H
vowelNasal           [  9: 10] =         O
pitch                [ 10: 11] =          :
unmatched            [ 11: 12] =           -
word                 [ 12: 18] =            Do:ni^
syllable             [ 12: 15] =            Do:
consonant            [ 12: 13] =            D
vowelNasal           [ 13: 14] =             o
pitch                [ 14: 15] =              :
syllable             [ 15: 18] =               ni^
consonant            [ 15: 16] =               n
vowelNasal           [ 16: 17] =                i
pitch                [ 17: 18] =                 ^
unmatched            [ 18: 20] =                  . 
sentence             [ 20: 42] =                    日本語は難しい!
unmatched            [ 20: 42] =                    日本語は難しい!
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}

func TestParseSango(t *testing.T) {
	s := `Ahöñ-ndönî AHÖÑ-NDÖNÎ, ândɛ bâa-mo-tɛnɛ? BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ!`
	parseMatchers := []ParseMatcher{SentenceFinal, SangoWord, SangoSyllable}
	actual := Parse(parseMatchers, s, ParseTree{Begin: 0, End: len(s)})
	actualStr := actual.Annotate(s, true, nil)
	expectStr := `
whole                [  0:112] = Ahöñ-ndönî AHÖÑ-NDÖNÎ, ândɛ bâa-mo-tɛnɛ? BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ!
sentence             [  0: 54] = Ahöñ-ndönî AHÖÑ-NDÖNÎ, ândɛ bâa-mo-tɛnɛ? 
word                 [  0:  6] = Ahöñ
closeVowel           [  0:  1] = A
closeVowelNasal      [  0:  1] = A
syllable             [  0:  1] = A
vowelNasal           [  0:  1] = A
syllable             [  1:  6] =  höñ
consonant            [  1:  2] =  h
closeVowelNasal      [  2:  6] =   öñ
vowelNasal           [  2:  6] =   öñ
closeVowel           [  2:  4] =   ö
nasal                [  4:  6] =    ñ
unmatched            [  6:  7] =     -
word                 [  7: 14] =      ndönî
syllable             [  7: 12] =      ndön
consonant            [  7:  9] =      nd
closeVowelNasal      [  9: 12] =        ön
vowelNasal           [  9: 12] =        ön
closeVowel           [  9: 11] =        ö
nasal                [ 11: 12] =         n
closeVowel           [ 12: 14] =          î
closeVowelNasal      [ 12: 14] =          î
syllable             [ 12: 14] =          î
vowelNasal           [ 12: 14] =          î
unmatched            [ 14: 15] =            
word                 [ 15: 21] =            AHÖÑ
closeVowel           [ 15: 16] =            A
closeVowelNasal      [ 15: 16] =            A
syllable             [ 15: 16] =            A
vowelNasal           [ 15: 16] =            A
syllable             [ 16: 21] =             HÖÑ
consonant            [ 16: 17] =             H
closeVowelNasal      [ 17: 21] =              ÖÑ
vowelNasal           [ 17: 21] =              ÖÑ
closeVowel           [ 17: 19] =              Ö
nasal                [ 19: 21] =               Ñ
unmatched            [ 21: 22] =                -
word                 [ 22: 29] =                 NDÖNÎ
syllable             [ 22: 27] =                 NDÖN
consonant            [ 22: 24] =                 ND
closeVowelNasal      [ 24: 27] =                   ÖN
vowelNasal           [ 24: 27] =                   ÖN
closeVowel           [ 24: 26] =                   Ö
nasal                [ 26: 27] =                    N
closeVowel           [ 27: 29] =                     Î
closeVowelNasal      [ 27: 29] =                     Î
syllable             [ 27: 29] =                     Î
vowelNasal           [ 27: 29] =                     Î
unmatched            [ 29: 31] =                      , 
word                 [ 31: 37] =                        ândɛ
closeVowelNasal      [ 31: 34] =                        ân
syllable             [ 31: 34] =                        ân
vowelNasal           [ 31: 34] =                        ân
closeVowel           [ 31: 33] =                        â
nasal                [ 33: 34] =                         n
syllable             [ 34: 37] =                          dɛ
consonant            [ 34: 35] =                          d
openVowel            [ 35: 37] =                           ɛ
vowelNasal           [ 35: 37] =                           ɛ
unmatched            [ 37: 38] =                             
word                 [ 38: 42] =                             bâa
syllable             [ 38: 41] =                             bâ
consonant            [ 38: 39] =                             b
closeVowel           [ 39: 41] =                              â
closeVowelNasal      [ 39: 41] =                              â
vowelNasal           [ 39: 41] =                              â
closeVowel           [ 41: 42] =                               a
closeVowelNasal      [ 41: 42] =                               a
syllable             [ 41: 42] =                               a
vowelNasal           [ 41: 42] =                               a
unmatched            [ 42: 43] =                                -
syllable             [ 43: 45] =                                 mo
word                 [ 43: 45] =                                 mo
consonant            [ 43: 44] =                                 m
closeVowel           [ 44: 45] =                                  o
closeVowelNasal      [ 44: 45] =                                  o
vowelNasal           [ 44: 45] =                                  o
unmatched            [ 45: 46] =                                   -
word                 [ 46: 52] =                                    tɛnɛ
syllable             [ 46: 49] =                                    tɛ
consonant            [ 46: 47] =                                    t
openVowel            [ 47: 49] =                                     ɛ
vowelNasal           [ 47: 49] =                                     ɛ
syllable             [ 49: 52] =                                      nɛ
consonant            [ 49: 50] =                                      n
openVowel            [ 50: 52] =                                       ɛ
vowelNasal           [ 50: 52] =                                       ɛ
unmatched            [ 52: 54] =                                        ? 
sentence             [ 54:112] =                                          BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ!
syllable             [ 54: 59] =                                          BƐ̂
word                 [ 54: 59] =                                          BƐ̂
consonant            [ 54: 55] =                                          B
openVowel            [ 55: 59] =                                           Ɛ̂
vowelNasal           [ 55: 59] =                                           Ɛ̂
unmatched            [ 59: 60] =                                             -
word                 [ 60: 77] =                                              kɔ̈mbïtɛbɛ̂
syllable             [ 60: 65] =                                              kɔ̈
consonant            [ 60: 61] =                                              k
openVowel            [ 61: 65] =                                               ɔ̈
vowelNasal           [ 61: 65] =                                               ɔ̈
syllable             [ 65: 69] =                                                 mbï
consonant            [ 65: 67] =                                                 mb
closeVowel           [ 67: 69] =                                                   ï
closeVowelNasal      [ 67: 69] =                                                   ï
vowelNasal           [ 67: 69] =                                                   ï
syllable             [ 69: 72] =                                                    tɛ
consonant            [ 69: 70] =                                                    t
openVowel            [ 70: 72] =                                                     ɛ
vowelNasal           [ 70: 72] =                                                     ɛ
syllable             [ 72: 77] =                                                      bɛ̂
consonant            [ 72: 73] =                                                      b
openVowel            [ 73: 77] =                                                       ɛ̂
vowelNasal           [ 73: 77] =                                                       ɛ̂
unmatched            [ 77: 78] =                                                         -
word                 [ 78: 90] =                                                          kɔ̈mbïtɛ
syllable             [ 78: 83] =                                                          kɔ̈
consonant            [ 78: 79] =                                                          k
openVowel            [ 79: 83] =                                                           ɔ̈
vowelNasal           [ 79: 83] =                                                           ɔ̈
syllable             [ 83: 87] =                                                             mbï
consonant            [ 83: 85] =                                                             mb
closeVowel           [ 85: 87] =                                                               ï
closeVowelNasal      [ 85: 87] =                                                               ï
vowelNasal           [ 85: 87] =                                                               ï
syllable             [ 87: 90] =                                                                tɛ
consonant            [ 87: 88] =                                                                t
openVowel            [ 88: 90] =                                                                 ɛ
vowelNasal           [ 88: 90] =                                                                 ɛ
unmatched            [ 90: 93] =                                                                  難
syllable             [ 93: 98] =                                                                    BƐ̂
word                 [ 93: 98] =                                                                    BƐ̂
consonant            [ 93: 94] =                                                                    B
openVowel            [ 94: 98] =                                                                     Ɛ̂
vowelNasal           [ 94: 98] =                                                                     Ɛ̂
unmatched            [ 98: 99] =                                                                       -
word                 [ 99:111] =                                                                        KƆ̈MBÏTƐ
syllable             [ 99:104] =                                                                        KƆ̈
consonant            [ 99:100] =                                                                        K
openVowel            [100:104] =                                                                         Ɔ̈
vowelNasal           [100:104] =                                                                         Ɔ̈
syllable             [104:108] =                                                                           MBÏ
consonant            [104:106] =                                                                           MB
closeVowel           [106:108] =                                                                             Ï
closeVowelNasal      [106:108] =                                                                             Ï
vowelNasal           [106:108] =                                                                             Ï
syllable             [108:111] =                                                                              TƐ
consonant            [108:109] =                                                                              T
openVowel            [109:111] =                                                                               Ɛ
vowelNasal           [109:111] =                                                                               Ɛ
unmatched            [111:112] =                                                                                !
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}

func TestParseRealSangoText(t *testing.T) {
	s := `
Na yâ tî taksï na Bangî

Benoît Kongbo

«… na pöpö tî âzo otâ sô na pêko tî taksï nî sô mbênî wâlï akono gï töngana bâgara.
Tîtene mbï dutï nzönî, mbï pûsu lo kêtê…»

Mbï löndö tî sïgî na pekô tî büngbïngö-terê na wayïndä tî Alliance française, na lo sô ayeke wayïndä tî Espace langue et formation, na kötä wasüngö-mbëtï Pierre Sammy sô mbï yeke gue na ndo tî lo ândö mîngi, na Gabriel, mbênî wasüngö-mbëtï ngâ na âmbâ tî mbï maseka-kôlï sô ayê kua tî süngö-mbëtï.

Mbï yeke tambûla tî gue na balabâla tî Dipandäa, ndo sô kâ mbï yeke mû taksï daä tî gue na yângâda.
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

Kötä pîpî tî mbênî taksï sï azî mbï na yâ tî gbüngö-li sô.
Ayeke gï na ngoi sô sï mbï hînga mbï sï na yângâ tî kötä lêgë awe.

Wagbötöngö-taksï nî atene na kötä gô:

— Ben-Zvi-Miskine!
Mo gue?

Mbï sâra mabôko tî mbï na lo tîtene lo lütï sï mbï mû lörö na mbâgë tî taksï nî.

Mbï tene na mâsa nî kôzo kûê tîtene mbï wara parâba pëpe:

— Mâsa, mbï yeke zûu na dambëtï tî Ben-Zvi.
Mbï yeke na pâta balëûse na okü.

Lo dë köngö mo bâa mo tene kâ fadë mbï fa na lo wüngö tî nginza nî pëpe:

— Me mo kö daä ma!

Mbï kö na yâ tî taksï nî na lörö.

— Awe, mâsa.
Ë gue.

Mbï wara lêgë tî dutï nzönî pëpe ngbanga tî sô na pöpö tî âzo otâ sô na pêko tî taksï nî sô, mbênî wâlï akono gï töngana bâgara.
Tîtene mbï dutï nzönî, mbï pûsu lo kêtê me wâlï sô asûku ngbângbâ, lo pete mbï na sïönî bê na terê tî yângâ tî taksï ni, sï lo tene:

— Mâsa, taksï tî mo nî ayeke kêtê töngasô ngbanga tî nye?
Bâa yê, gï na lëkërëngö têre tî mbï tîtene maseka-kôlï sô gerê tî lô akpa gï tî kôlô sô adutï nzönî, yâ tî terê tî mbï kûê ayeke so.

Wagbötöngö-taksï nî atene na ngîâ na yângâ tî lo:

— Mo hînga fadë, mamâ, mo fûta ngêrë tî ndo ûse tîtene mo bâa pâsi pëpe!

Wâlï sô akpa gï doli sô atï na wâ agbï, lo tene:

— Âla mä tënë tî kpändä sô!
Lo hînga mbï äpe, mbï hînga lo äpe.
Lo yeke tene tënë na mbï mo bâa mo tene mbï yeke mbênî kötä konô sô asï yä tî taksï nî kûê.
Mä mbï sï, môlengê.
Mbï lîngbi tî dutï töngana mamâ tî mo.
Gï na särängö terê tî âla töngasô, âla âmaseka, sï âla ke kûi hîo sô.

Ndo adë ngii sï pakara sô na gatï tî mamâ sô azî yângâ:

— Mbï yeke hûnda terê tî mbï, mbï tene, âla wara sïönî mayëre tî zîa âzo osïö na pekö sô na ndo wa?
Taksï tî mo nî ayeke kêtê ahön kopo tî âsusu.

Wagbötöngö-kutukutu nî azî yângâ pëpe tënë tî bîâ sô lo yeke hë na gbegô.

Kôlï sô na terê tî wâlï sô na kotï tî mbï akîri atene:

— Na yâ tî âmbênî ködörö, na tâpandë tî Farânzi sô mbï ke gue daä töngana lêgë tî yäkä tî mbï, âtaksï ayeke mû gï zo ôko.
Na ndönî, âla ayeke zîa lo na ndo sô lo yeke gue daä.

Gï sô lo tene töngasô sô âkperekpere nî alöndö ngangü ahön ndö tî îâ sô ayeke toto ngâ na kötä köngö tî taksï ni.

Masia sô adutï na pöpö tî Wagbötöngö-taksï na mbênî maseka sô amo li tî lo ndâli tî kötä kä sô ayeke daä atene:

— Na mbâgë, tënë tî mo nî ayeke na lêgë nî, me na mbâgë ngâ, alingbi mo hîinga ë dutï nzönî na yâ tî taksï ahön bus.

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
Sô fadë mbï kö na yâ tî taksï sô, sösö ahon ndö tî mbï mîngi.
Mbï yê tî sösö na ngangü tî terê tî mbï kûê sï pupu aya mbï.
Na Bangî, mo sösö, ayeke hë mo ngbii asï na lâ tî kûâ tî mo.
Ka mbï sösö nî sï sasa asïgîgî agbû bongö tî mbï, fadë mbï yeke sâra töngana nye?
Alîngbi mbï dutï kpô tî pete sösö sô ayê gï tî sungba na ngangü sô.

Wagbötöngö-taksï nî azî gô tî fonôno nî ngangü.
Töngana mbênî yê tî ndimä, sösö sô mbï kânga fadë gbä sô asïgî gï yeke bâa mo tene mozoko nî ayê däa tîtene lo sïgî.

Wâlï sô na terê tî mbï sï akpa gï kötä bäbolo sô atene na kötä gô:

— Hêe, mâsa, mo pêe tî kîri na gô tî fonôno-nî kêtê äpe?
Ë yeke na yâ tî nganda äpe, ë yeke na yâ tî taksï laâ!

Âdëköngö alöndö mbâgë na mbâgë:

— Nye laâ afün töngasô, mâsa?

— Mo yö fünngö nyama na yâ tî taksï nî wala nye?

— Kpändä wa laâ asösö sô?

— Mozoko sô ayeke toto ngangü mîngi sô laâ amû lêgë na zo sô tî sösö sô.

— Sösö nî afün gï töngana kûâ tî mbo, dê!

— Alîngbi lo gue abâa yâ tî purû tî lo sï.

Kötä wâlï sô na terê tî mbï sï amû ndo kûê sô atene:

— Mo bâa sô, ayeke gï mâsa laâ asösö sô.

Me wagbötöngö-taksï nî akîri tënë na wâlï nî atene:

— Tî mbï ayeke gï mo, pûpûlenge waâ!
Mbï hînga ka sôsô tî âwâlï sô ngbondä tî âla akono mîngi sô kûê.

Lê tî wâlï sô abe na ngonzo sï lo hûnda tîtene mâsa nî azîa lo na sêse fadëfadësô.
Ë yeke na barâsi balëokü na sängbï-lêgë tî Ben-Zvi-Sica.
Taksï nî alütï na dëmängö terê tî lo.
Wâlï nî asïgîgiî na kugbë tî sâki ôko lo mû na mâsa nî.
voazîne tî mbï nî, sô ayeke voazîne tî mbï mbênî pëpe tëtî sô lo zûu na sêse awe sô, ake na kürü gô:

— Mo dîko tanga tî nginza sô mo kîri na mbï sô laâ?

— Mbï kîri na mo ngbangbo na balë mbârâmbârâ.
Ka tënë nî ayeke na ndo wa, mamâ tî mbï?

— Mbï tene na mo fadë mbï tene mbï yeke gï na balëûse na okü, äpe?

— Mbï yeke mû ka balëûse na okü töngana mo yeke gï na balëûse na okü.

Mâsa nî azî gerê, taksï tî lo nî ahön yongôro na terê tî wâlï tî gïngö tënë sô sï lo tene nï hînga li tî nï na âzo tî Bêafrîka pëpe.
Ngêrë tî taksï ayeke balëotâ.
Ayeke na nzönî bê laâ sï lo yeke yô âzo na balëûse na okü, sô ayeke ngêrë tî köngö na yâ tî bus sô.

Töngana kutukutu nî asï na dambëtï tî Ben-Zvi awe, mbï zûu na sêse ngâ.
Na mbâgë tî lêgë na dawäa tî mbênî dapïkängö-singa, Sylvianne ayeke kü mbï töngana tî sô ë mä terê daä.
Mbï gue na terê tî lo sï lêgëôko, ë mû lêgë sô ague na Centre Baobab tî kîri na yângâda.

Lo hûnda mbï tî hînga:

— Büngbïngö-terê nî atambûla na lêgë nî?

— Iin.
Ahûnda na ë tî bûngbi tî sâra mbênî bûku na farânzi ngâ na sängö.

— Mo hînga yê sô mo yê tî sîra mbëtï na ndö nî sô awe?

— Mbï yê tî vunga sango tî yê sô asï fadë na yâ tî taksï sô.

— Nye laâ asï nye?

—Âzo tî yâ tî taksï nî adîko tënë töna âkûkulû sï mbï sösö tîtene âla dutï kpo.

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
	parseMatchers := []ParseMatcher{SentenceFinal, SangoWord, SangoSyllable}
	actual := Parse(parseMatchers, s, ParseTree{Begin: 0, End: len(s)})
	actualStr := actual.Annotate(s, false, map[string]bool{
		"closeVowel":      true,
		"closeVowelNasal": true,
		"consonant":       true,
		"nasal":           true,
		"openVowel":       true,
		"syllable":        true,
		"unmatched":       true,
		"vowelNasal":      true,
		"whole":           true,
		"word":            true,
	})
	expectStr := `
sentence             [  1: 30] = Na yâ tî taksï na Bangî


sentence             [ 30: 46] = Benoît Kongbo


sentence             [ 46:151] = «… na pöpö tî âzo otâ sô na pêko tî taksï nî sô mbênî wâlï akono gï töngana bâgara.

sentence             [151:206] = Tîtene mbï dutï nzönî, mbï pûsu lo kêtê…»


sentence             [206:557] = Mbï löndö tî sïgî na pekô tî büngbïngö-terê na wayïndä tî Alliance française, na lo sô ayeke wayïndä tî Espace langue et formation, na kötä wasüngö-mbëtï Pierre Sammy sô mbï yeke gue na ndo tî lo ândö mîngi, na Gabriel, mbênî wasüngö-mbëtï ngâ na âmbâ tî mbï maseka-kôlï sô ayê kua tî süngö-mbëtï.


sentence             [557:672] = Mbï yeke tambûla tî gue na balabâla tî Dipandäa, ndo sô kâ mbï yeke mû taksï daä tî gue na yângâda.

sentence             [672:731] = Mbï yeke gbû li tî mbï na ndö tî kua sô amû na ë.

sentence             [731:824] = Lo hûnda na ë tî sû mbênî mbëtï sô ë lîngbi tî mû na lo na yâ tî yenga ûse.


sentence             [824:846] = Lo tene fadë na ë:


sentence             [846:984] = — Kodë nî, ayeke tî sâra mbëtï tî âla nî töngana mbênî bûku sô yângâ tî ködörö nî ayeke dutï sängö na farânzi.

sentence             [1123:1148] = na âzo tî Bêafrîka.


sentence             [1148:1383] = Tî taâ tënë nî, pialö sô abâa gïgî na pekô tî gbüngö-li tî vassilis, wasüngö-mbëtï tî farânzi sô taâ marä tî lo ayeke gerêke sô asû mbênî bûku na ndö tî sängö sô ïrï nî ayeke Les mots étrangers.

sentence             [1383:1522] = Na ûse gängö tî lo na sêse tî Bêafrîka, lo sâra na ë tënë tî pialö tî süngö bûku lêgëôko na sängö ngâ na farânzi.

sentence             [1522:1705] = Ë kûê ë yê daä na nî ngbanga tî sô ayeke nzönî kodë tî zîa nzara tî dïköngö-mbëtï na yâ tî âzo tî Bêafriîka, sô ahînga ngâ na âla sô ahînga nî pëpe.

sentence             [1705:1963] = Mawa nî ayeke sô anînga kêtêz mbï hön tënë tî kua tî süngö mbênî bûku na Limoges na sêse tî Farânzi sï mbï wara âsango na ndö tî pialö sô pëpe ngbii asï na lâkûi tî büngbïngö-terê sô mbï löndö tî tene tënë nî sô.


sentence             [1963:2103] = Mbï yeke gi bë tî mbï ngbanga tî sô mbï yeke sû mbênî bûku sô mbï yê tî tokua na da lëkërëngö nî na sêse tî Gameroun.

sentence             [2103:2224] = Mbï zîa bê tî mbï daä mîngi sï mbeto ayeke gbû mbï tî wara lêgë pëpe tî sû mbëtï sô Alliance ahûnda.

sentence             [2224:2306] = Mbï yê tî hûnzi bûku sô hîo sï mbï yê tî sâra mbênî yê ndê pëpe.

sentence             [2306:2376] = Mbï yê tîtene ahûnda mbï hündängö sï mbï sû mbëtï pëpe.

sentence             [2376:2440] = Mbï yeke sâra nî na ngoi sô mbï yeke na âsango tî vunga.

sentence             [2440:2560] = Ânzönî tënë sô mbï yeke na nî tîtene na ndö tî ködörö tî mbï ayeke na yâ tî kôzo bûku tî mbï awe.

sentence             [2560:2604] = Tënë wa ânde sï mbï zîa na ndö nî?


sentence             [2604:2679] = Kötä pîpî tî mbênî taksï sï azî mbï na yâ tî gbüngö-li sô.

sentence             [2679:2761] = Ayeke gï na ngoi sô sï mbï hînga mbï sï na yângâ tî kötä lêgë awe.


sentence             [2761:2808] = Wagbötöngö-taksï nî atene na kötä gô:


sentence             [2808:2829] = — Ben-Zvi-Miskine!

sentence             [2829:2838] = Mo gue?


sentence             [2838:2938] = Mbï sâra mabôko tî mbï na lo tîtene lo lütï sï mbï mû lörö na mbâgë tî taksï nî.


sentence             [2938:3007] = Mbï tene na mâsa nî kôzo kûê tîtene mbï wara parâba pëpe:


sentence             [3007:3059] = — Mâsa, mbï yeke zûu na dambëtï tî Ben-Zvi.

sentence             [3059:3098] = Mbï yeke na pâta balëûse na okü.


sentence             [3098:3184] = Lo dë köngö mo bâa mo tene kâ fadë mbï fa na lo wüngö tî nginza nî pëpe:


sentence             [3184:3208] = — Me mo kö daä ma!


sentence             [3208:3251] = Mbï kö na yâ tî taksï nî na lörö.


sentence             [3251:3267] = — Awe, mâsa.

sentence             [3267:3276] = Ë gue.


sentence             [3276:3433] = Mbï wara lêgë tî dutï nzönî pëpe ngbanga tî sô na pöpö tî âzo otâ sô na pêko tî taksï nî sô, mbênî wâlï akono gï töngana bâgara.

sentence             [3433:3594] = Tîtene mbï dutï nzönî, mbï pûsu lo kêtê me wâlï sô asûku ngbângbâ, lo pete mbï na sïönî bê na terê tî yângâ tî taksï ni, sï lo tene:


sentence             [3594:3663] = — Mâsa, taksï tî mo nî ayeke kêtê töngasô ngbanga tî nye?

sentence             [3663:3829] = Bâa yê, gï na lëkërëngö têre tî mbï tîtene maseka-kôlï sô gerê tî lô akpa gï tî kôlô sô adutï nzönî, yâ tî terê tî mbï kûê ayeke so.


sentence             [3829:3890] = Wagbötöngö-taksï nî atene na ngîâ na yângâ tî lo:


sentence             [3890:3978] = — Mo hînga fadë, mamâ, mo fûta ngêrë tî ndo ûse tîtene mo bâa pâsi pëpe!


sentence             [3978:4036] = Wâlï sô akpa gï doli sô atï na wâ agbï, lo tene:


sentence             [4036:4074] = — Âla mä tënë tî kpändä sô!

sentence             [4074:4116] = Lo hînga mbï äpe, mbï hînga lo äpe.

sentence             [4116:4226] = Lo yeke tene tënë na mbï mo bâa mo tene mbï yeke mbênî kötä konô sô asï yä tî taksï nî kûê.

sentence             [4226:4251] = Mä mbï sï, môlengê.

sentence             [4251:4297] = Mbï lîngbi tî dutï töngana mamâ tî mo.

sentence             [4297:4384] = Gï na särängö terê tî âla töngasô, âla âmaseka, sï âla ke kûi hîo sô.


sentence             [4384:4451] = Ndo adë ngii sï pakara sô na gatï tî mamâ sô azî yângâ:


sentence             [4451:4571] = — Mbï yeke hûnda terê tî mbï, mbï tene, âla wara sïönî mayëre tî zîa âzo osïö na pekö sô na ndo wa?

sentence             [4571:4626] = Taksï tî mo nî ayeke kêtê ahön kopo tî âsusu.


sentence             [4626:4717] = Wagbötöngö-kutukutu nî azî yângâ pëpe tënë tî bîâ sô lo yeke hë na gbegô.


sentence             [4717:4785] = Kôlï sô na terê tî wâlï sô na kotï tî mbï akîri atene:


sentence             [4785:4937] = — Na yâ tî âmbênî ködörö, na tâpandë tî Farânzi sô mbï ke gue daä töngana lêgë tî yäkä tî mbï, âtaksï ayeke mû gï zo ôko.

sentence             [4937:4998] = Na ndönî, âla ayeke zîa lo na ndo sô lo yeke gue daä.


sentence             [4998:5136] = Gï sô lo tene töngasô sô âkperekpere nî alöndö ngangü ahön ndö tî îâ sô ayeke toto ngâ na kötä köngö tî taksï ni.


sentence             [5136:5269] = Masia sô adutï na pöpö tî Wagbötöngö-taksï na mbênî maseka sô amo li tî lo ndâli tî kötä kä sô ayeke daä atene:


sentence             [5269:5410] = — Na mbâgë, tënë tî mo nî ayeke na lêgë nî, me na mbâgë ngâ, alingbi mo hîinga ë dutï nzönî na yâ tî taksï ahön bus.


sentence             [5410:5487] = Môlengê-kôlï sô na kä na li tî lo sô atene na môlengê-wâlï nî:


sentence             [5487:5526] = — Tâtënë laä mo tene sô, îtä.

sentence             [5526:5586] = Mbï sô mo bâa mbï sô, mbï yeke mû bus lâôko pëpe.

sentence             [5586:5632] = Mbï yeke zo sô ayê dütïngö kpô mîngi.

sentence             [5632:5791] = Andâa, na yâ tî bus, zo kûê ayeke tene tënë, ayeke dë köngö, ayeke hë ngîâ na yâ tî mê tî âzo, ngâ ayeke sösö na mbänä tî bê tî lo.


sentence             [5791:5894] = Kôlï tî hümä sô adutï na mbâgë kâ mo bâa mo tene lo löndö na ködörö wandê sô atene:


sentence             [5894:5981] = —Tî mbï, mbï kpë bus ngbanga tî sô mbï yê fünngö tî terê tî âzo pëpe.

sentence             [5981:6150] = Mbï lîngbi tî kânga bê tî mbï na ndö tî fünngö gbe tî mabôko tî âmötöröbôi pëpe, wala ngâ na gbikï tî âla sô ayeke buba bongö tî mbï pëpe.


sentence             [6150:6216] = Ka mbï, âla yeke hûnda terê tî âla, mbï tene tî mbï nye?

sentence             [6216:6245] = Mbï zî yângâ ôko pëpe.

sentence             [6245:6321] = Sô fadë mbï kö na yâ tî taksï sô, sösö ahon ndö tî mbï mîngi.

sentence             [6321:6396] = Mbï yê tî sösö na ngangü tî terê tî mbï kûê sï pupu aya mbï.

sentence             [6396:6467] = Na Bangî, mo sösö, ayeke hë mo ngbii asï na lâ tî kûâ tî mo.

sentence             [6467:6565] = Ka mbï sösö nî sï sasa asïgîgî agbû bongö tî mbï, fadë mbï yeke sâra töngana nye?

sentence             [6565:6647] = Alîngbi mbï dutï kpô tî pete sösö sô ayê gï tî sungba na ngangü sô.


sentence             [6647:6706] = Wagbötöngö-taksï nî azî gô tî fonôno nî ngangü.

sentence             [6706:6848] = Töngana mbênî yê tî ndimä, sösö sô mbï kânga fadë gbä sô asïgî gï yeke bâa mo tene mozoko nî ayê däa tîtene lo sïgî.


sentence             [6848:6931] = Wâlï sô na terê tî mbï sï akpa gï kötä bäbolo sô atene na kötä gô:


sentence             [6931:7002] = — Hêe, mâsa, mo pêe tî kîri na gô tî fonôno-nî kêtê äpe?

sentence             [7002:7067] = Ë yeke na yâ tî nganda äpe, ë yeke na yâ tî taksï laâ!


sentence             [7067:7110] = Âdëköngö alöndö mbâgë na mbâgë:


sentence             [7110:7148] = — Nye laâ afün töngasô, mâsa?


sentence             [7148:7207] = — Mo yö fünngö nyama na yâ tî taksï nî wala nye?


sentence             [7207:7242] = — Kpändä wa laâ asösö sô?


sentence             [7242:7331] = — Mozoko sô ayeke toto ngangü mîngi sô laâ amû lêgë na zo sô tî sösö sô.


sentence             [7331:7386] = — Sösö nî afün gï töngana kûâ tî mbo, dê!


sentence             [7386:7439] = — Alîngbi lo gue abâa yâ tî purû tî lo sï.


sentence             [7439:7506] = Kötä wâlï sô na terê tî mbï sï amû ndo kûê sô atene:


sentence             [7506:7558] = — Mo bâa sô, ayeke gï mâsa laâ asösö sô.


sentence             [7558:7622] = Me wagbötöngö-taksï nî akîri tënë na wâlï nî atene:


sentence             [7622:7667] = — Tî mbï ayeke gï mo, pûpûlenge waâ!

sentence             [7667:7749] = Mbï hînga ka sôsô tî âwâlï sô ngbondä tî âla akono mîngi sô kûê.


sentence             [7749:7847] = Lê tî wâlï sô abe na ngonzo sï lo hûnda tîtene mâsa nî azîa lo na sêse fadëfadësô.

sentence             [7847:7913] = Ë yeke na barâsi balëokü na sängbï-lêgë tî Ben-Zvi-Sica.

sentence             [7913:7960] = Taksï nî alütï na dëmängö terê tî lo.

sentence             [7960:8029] = Wâlï nî asïgîgiî na kugbë tî sâki ôko lo mû na mâsa nî.

sentence             [8029:8151] = voazîne tî mbï nî, sô ayeke voazîne tî mbï mbênî pëpe tëtî sô lo zûu na sêse awe sô, ake na kürü gô:


sentence             [8151:8213] = — Mo dîko tanga tî nginza sô mo kîri na mbï sô laâ?


sentence             [8213:8268] = — Mbï kîri na mo ngbangbo na balë mbârâmbârâ.

sentence             [8268:8316] = Ka tënë nî ayeke na ndo wa, mamâ tî mbï?


sentence             [8316:8395] = — Mbï tene na mo fadë mbï tene mbï yeke gï na balëûse na okü, äpe?


sentence             [8395:8478] = — Mbï yeke mû ka balëûse na okü töngana mo yeke gï na balëûse na okü.


sentence             [8478:8640] = Mâsa nî azî gerê, taksï tî lo nî ahön yongôro na terê tî wâlï tî gïngö tënë sô sï lo tene nï hînga li tî nï na âzo tî Bêafrîka pëpe.

sentence             [8640:8676] = Ngêrë tî taksï ayeke balëotâ.

sentence             [8676:8796] = Ayeke na nzönî bê laâ sï lo yeke yô âzo na balëûse na okü, sô ayeke ngêrë tî köngö na yâ tî bus sô.


sentence             [8796:8878] = Töngana kutukutu nî asï na dambëtï tî Ben-Zvi awe, mbï zûu na sêse ngâ.

sentence             [8878:9003] = Na mbâgë tî lêgë na dawäa tî mbênî dapïkängö-singa, Sylvianne ayeke kü mbï töngana tî sô ë mä terê daä.

sentence             [9003:9109] = Mbï gue na terê tî lo sï lêgëôko, ë mû lêgë sô ague na Centre Baobab tî kîri na yângâda.


sentence             [9109:9137] = Lo hûnda mbï tî hînga:


sentence             [9137:9190] = — Büngbïngö-terê nî atambûla na lêgë nî?


sentence             [9190:9199] = — Iin.

sentence             [9199:9279] = Ahûnda na ë tî bûngbi tî sâra mbênî bûku na farânzi ngâ na sängö.


sentence             [9279:9348] = — Mo hînga yê sô mo yê tî sîra mbëtï na ndö nî sô awe?


sentence             [9348:9424] = — Mbï yê tî vunga sango tî yê sô asï fadë na yâ tî taksï sô.


sentence             [9424:9448] = — Nye laâ asï nye?


sentence             [9448:9551] = —Âzo tî yâ tî taksï nî adîko tënë töna âkûkulû sï mbï sösö tîtene âla dutï kpo.


sentence             [9551:9624] = Sylvianne atï na kötä ngîâ sï nzara tî mbênî sösö agbû mbï.

sentence             [9624:9645] = Me mbï kânga nî.


sentence             [9645:9687] = — Ahînga atene mo laâ mo sösö sô?


sentence             [9687:9764] = — Mbï sösö nî na ngoi sô mâsa nî azî gô tî fonôno nî ngangü.


sentence             [9764:9816] = — Tî taâ tënë nî, atumba fünngö nî äpe?


sentence             [9816:9898] = — Âhînga fadë atene mbï laâ mbï sösö sô, ka kamënë ayeke gbû mbï.


sentence             [9898:9939] = — Sösö ayeke yê tî kamënë pëpe!

sentence             [9939:9975] = Sösö afa atene zo ayeke na sêni.

sentence             [9975:10083] = Atâa töngana âsûru yâ tî zo awe sô, ayeke kü tîtene zo nî asösö awe, sï lo lîngbi tî te yê.

sentence             [10083:10142] = Âmunzû ayeke dema pëpe töngana mbâ tî âla asösö.


sentence             [10142:10260] = Na dabüzë tî Amadi, mbï vo gbâ tî mânga sô ïrï nî Menthol sô sï mbï zîa wâ na terê tî mbênî ôko.

sentence             [10260:10332] = Pialö tî süngö bûku na farânzi-sängö sô akîri na li tî mbï.

sentence             [10332:10437] = Tî särängö mbëtï na farânzi sï mbï kîri mbï lekere nî na sängö sô aso li tî mbï mîngi.

sentence             [10437:10557] = Töngana nye ânde sï mbï sâra mbëtï na yângâ tî ködörö tî mbï sô mbï hînga tî tene nzönî pëpe sô?

sentence             [10557:10665] = Âfa nî fadë na mbî na dambëtï, mbï lîngbi tî wara kpälë tî särängö mbëtï na sängö pëpe.

sentence             [10665:10688] = Me mawa nî ayeke sô!

sentence             [10688:10926] = Gï na fängö na mbï tî tene ngâ na tî sâra mbëtï nzönî na yângâ tî ködörö tî âzo sô amû na ë Dipandäa, âzo tî tënëngö kötä farânzi tî ködörö tî ë asâra sï ë hînga tënëngö sängö nzönî pëpe.

sentence             [10926:10997] = Âmbupa tî farânzi asï sïngö na yâ tî sängö sô ë yeke tene.

sentence             [10997:11145] = Mbï yeke tö ndâ tî sâra mbëtï na marä sängö sô, ngbanga tî sô nzönî ködë tî manda yângâ tî ködörö, ayeke gï tî tene nî.


sentence             [11145:11201] = Mo yê tî bâa, sîngamê tî mbï tî mabôko atoto.


sentence             [11201:11212] = — Alô?


sentence             [11212:11228] = Mbï ma atene:


sentence             [11228:11254] = — Töngana nye, petï?


sentence             [11254:11271] = — Zo wa laâ?


sentence             [11271:11303] = — Mo girisa gô tî mbï awe?

sentence             [11303:11324] = Mbï Yvon laâ äpe.

sentence             [11324:11344] = Mo yeke na ndo wa?


sentence             [11344:11382] = — Mbï yeke gue na yângâda laâ.


sentence             [11382:11450] = — Mo pêe tî gä hîo na ndo sô ë ike wara ka terê daä sô?


sentence             [11450:11571] = — Mbï pêe, me na yâ tî ngbonga ôko ngbanga tî sô mbï yê tî sû mbênî mbëtï na farânzi ngâ na sängö.


sentence             [11571:11610] = — Alîngbi tî kü kêkerêke äpe?


sentence             [11610:11714] = — Mbï yê tî tö ndâ nî gï fafadësô, töngasô äpe sô, mbï lîngbi tî sû nî ânde äpe.


sentence             [11714:11826] = — Mbï lîngbi tî kü gï töngana mo yê daä tî sâra tënë tî mbï na yâ tî mbëtï tî mo nî sï.


sentence             [11826:11873] = Mbï tene na lo na ngîâ na yângâ tî mbï.

sentence             [11873:11874] = 	
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}
