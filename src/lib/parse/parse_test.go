package parse

import (
	"log"
	"os"
	"strings"
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
	actual := Parse(parseMatchers, s, ParseTree{Begin: 0, End: len(s), Name: "whole"})
	var actualBuilder strings.Builder = actual.Annotate(s, true, map[string]bool{
		"closeVowel":      true,
		"closeVowelNasal": true,
		"consonant":       true,
		"nasal":           true,
		"openVowel":       true,
		"pitch":           true,
		"sentence":        true,
		"shift":           true,
		"syllable":        true,
		"unmatched":       true,
		"vowelNasal":      true,
		"whole":           true,
		"word":            true,
	})
	actualStr := actualBuilder.String()
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
	actual := Parse(parseMatchers, s, ParseTree{Begin: 0, End: len(s), Name: "whole"})
	var actualBuilder strings.Builder = actual.Annotate(s, true, map[string]bool{
		"closeVowel":      true,
		"closeVowelNasal": true,
		"consonant":       true,
		"nasal":           true,
		"openVowel":       true,
		"pitch":           true,
		"sentence":        true,
		"shift":           true,
		"syllable":        true,
		"unmatched":       true,
		"vowelNasal":      true,
		"whole":           true,
		"word":            true,
	})
	actualStr := actualBuilder.String()
	expectStr := `
whole                [  0:112] = Ahöñ-ndönî AHÖÑ-NDÖNÎ, ândɛ bâa-mo-tɛnɛ? BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ!
sentence             [  0: 54] = Ahöñ-ndönî AHÖÑ-NDÖNÎ, ândɛ bâa-mo-tɛnɛ? 
word                 [  0:  6] = Ahöñ
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
unmatched            [  6:  7] =     -
word                 [  7: 14] =      ndönî
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
unmatched            [ 14: 15] =            
word                 [ 15: 21] =            AHÖÑ
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
unmatched            [ 21: 22] =                -
word                 [ 22: 29] =                 NDÖNÎ
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
unmatched            [ 29: 31] =                      , 
word                 [ 31: 37] =                        ândɛ
syllable             [ 31: 33] =                        â
vowelNasal           [ 31: 33] =                        â
closeVowelNasal      [ 31: 33] =                        â
closeVowel           [ 31: 33] =                        â
syllable             [ 33: 37] =                         ndɛ
consonant            [ 33: 35] =                         nd
vowelNasal           [ 35: 37] =                           ɛ
openVowel            [ 35: 37] =                           ɛ
unmatched            [ 37: 38] =                             
word                 [ 38: 42] =                             bâa
syllable             [ 38: 41] =                             bâ
consonant            [ 38: 39] =                             b
vowelNasal           [ 39: 41] =                              â
closeVowelNasal      [ 39: 41] =                              â
closeVowel           [ 39: 41] =                              â
syllable             [ 41: 42] =                               a
vowelNasal           [ 41: 42] =                               a
closeVowelNasal      [ 41: 42] =                               a
closeVowel           [ 41: 42] =                               a
unmatched            [ 42: 43] =                                -
word                 [ 43: 45] =                                 mo
syllable             [ 43: 45] =                                 mo
consonant            [ 43: 44] =                                 m
vowelNasal           [ 44: 45] =                                  o
closeVowelNasal      [ 44: 45] =                                  o
closeVowel           [ 44: 45] =                                  o
unmatched            [ 45: 46] =                                   -
word                 [ 46: 52] =                                    tɛnɛ
syllable             [ 46: 49] =                                    tɛ
consonant            [ 46: 47] =                                    t
vowelNasal           [ 47: 49] =                                     ɛ
openVowel            [ 47: 49] =                                     ɛ
syllable             [ 49: 52] =                                      nɛ
consonant            [ 49: 50] =                                      n
vowelNasal           [ 50: 52] =                                       ɛ
openVowel            [ 50: 52] =                                       ɛ
unmatched            [ 52: 54] =                                        ? 
sentence             [ 54:112] =                                          BƐ̂-kɔ̈mbïtɛbɛ̂-kɔ̈mbïtɛ難BƐ̂-KƆ̈MBÏTƐ!
word                 [ 54: 59] =                                          BƐ̂
syllable             [ 54: 59] =                                          BƐ̂
consonant            [ 54: 55] =                                          B
vowelNasal           [ 55: 59] =                                           Ɛ̂
openVowel            [ 55: 59] =                                           Ɛ̂
unmatched            [ 59: 60] =                                             -
word                 [ 60: 77] =                                              kɔ̈mbïtɛbɛ̂
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
unmatched            [ 77: 99] =                                                         -kɔ̈mbïtɛ難BƐ̂-
word                 [ 99:111] =                                                                        KƆ̈MBÏTƐ
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
unmatched            [111:112] =                                                                                !
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}

func TestParseRealSangoText(t *testing.T) {
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
	parseMatchers := []ParseMatcher{SentenceFinal, SangoWord, SangoSyllable}
	actual := Parse(parseMatchers, s, ParseTree{Begin: 0, End: len(s), Name: "whole"})
	var actualBuilder strings.Builder = actual.Annotate(s, false, map[string]bool{
		"sentence": true,
	})
	actualStr := actualBuilder.String()
	expectStr := `
sentence             [  0: 31] = 
Na yâ tî takisï na Bangî


sentence             [ 31: 47] = Benoît Kongbo


sentence             [ 47:153] = «… na pöpö tî âzo otâ sô na pêko tî takisï nî sô mbênî wâlï akono gï töngana bâgara.

sentence             [153:208] = Tîtene mbï dutï nzönî, mbï pûsu lo kêtê…»


sentence             [208:559] = Mbï löndö tî sïgî na pekô tî büngbïngö-terê na wayïndä tî Alliance française, na lo sô ayeke wayïndä tî Espace langue et formation, na kötä wasüngö-mbëtï Pierre Sammy sô mbï yeke gue na ndo tî lo ândö mîngi, na Gabriel, mbênî wasüngö-mbëtï ngâ na âmbâ tî mbï maseka-kôlï sô ayê kua tî süngö-mbëtï.


sentence             [559:675] = Mbï yeke tambûla tî gue na balabâla tî Dipandäa, ndo sô kâ mbï yeke mû takisï daä tî gue na yângâda.

sentence             [675:734] = Mbï yeke gbû li tî mbï na ndö tî kua sô amû na ë.

sentence             [734:827] = Lo hûnda na ë tî sû mbênî mbëtï sô ë lîngbi tî mû na lo na yâ tî yenga ûse.


sentence             [827:849] = Lo tene fadë na ë:


sentence             [849:987] = — Kodë nî, ayeke tî sâra mbëtï tî âla nî töngana mbênî bûku sô yângâ tî ködörö nî ayeke dutï sängö na farânzi.

sentence             [987:1151] = Ndâ tî pialö sô, ayeke tî mû lêgë na âmbênî na pöpö tî âla sô ayeke sû mbëtï na gbemingo tî fa terê tî âla polêlê
na âzo tî Bêafrîka.


sentence             [1151:1386] = Tî taâ tënë nî, pialö sô abâa gïgî na pekô tî gbüngö-li tî vassilis, wasüngö-mbëtï tî farânzi sô taâ marä tî lo ayeke gerêke sô asû mbênî bûku na ndö tî sängö sô ïrï nî ayeke Les mots étrangers.

sentence             [1386:1525] = Na ûse gängö tî lo na sêse tî Bêafrîka, lo sâra na ë tënë tî pialö tî süngö bûku lêgëôko na sängö ngâ na farânzi.

sentence             [1525:1708] = Ë kûê ë yê daä na nî ngbanga tî sô ayeke nzönî kodë tî zîa nzara tî dïköngö-mbëtï na yâ tî âzo tî Bêafriîka, sô ahînga ngâ na âla sô ahînga nî pëpe.

sentence             [1708:1966] = Mawa nî ayeke sô anînga kêtêz mbï hön tënë tî kua tî süngö mbênî bûku na Limoges na sêse tî Farânzi sï mbï wara âsango na ndö tî pialö sô pëpe ngbii asï na lâkûi tî büngbïngö-terê sô mbï löndö tî tene tënë nî sô.


sentence             [1966:2106] = Mbï yeke gi bë tî mbï ngbanga tî sô mbï yeke sû mbênî bûku sô mbï yê tî tokua na da lëkërëngö nî na sêse tî Gameroun.

sentence             [2106:2227] = Mbï zîa bê tî mbï daä mîngi sï mbeto ayeke gbû mbï tî wara lêgë pëpe tî sû mbëtï sô Alliance ahûnda.

sentence             [2227:2309] = Mbï yê tî hûnzi bûku sô hîo sï mbï yê tî sâra mbênî yê ndê pëpe.

sentence             [2309:2379] = Mbï yê tîtene ahûnda mbï hündängö sï mbï sû mbëtï pëpe.

sentence             [2379:2443] = Mbï yeke sâra nî na ngoi sô mbï yeke na âsango tî vunga.

sentence             [2443:2563] = Ânzönî tënë sô mbï yeke na nî tîtene na ndö tî ködörö tî mbï ayeke na yâ tî kôzo bûku tî mbï awe.

sentence             [2563:2607] = Tënë wa ânde sï mbï zîa na ndö nî?


sentence             [2607:2683] = Kötä pîpî tî mbênî takisï sï azî mbï na yâ tî gbüngö-li sô.

sentence             [2683:2765] = Ayeke gï na ngoi sô sï mbï hînga mbï sï na yângâ tî kötä lêgë awe.


sentence             [2765:2813] = Wagbötöngö-takisï nî atene na kötä gô:


sentence             [2813:2834] = — Ben-Zvi-Miskine!

sentence             [2834:2843] = Mo gue?


sentence             [2843:2944] = Mbï sâra mabôko tî mbï na lo tîtene lo lütï sï mbï mû lörö na mbâgë tî takisï nî.


sentence             [2944:3013] = Mbï tene na mâsa nî kôzo kûê tîtene mbï wara parâba pëpe:


sentence             [3013:3065] = — Mâsa, mbï yeke zûu na dambëtï tî Ben-Zvi.

sentence             [3065:3104] = Mbï yeke na pâta balëûse na okü.


sentence             [3104:3190] = Lo dë köngö mo bâa mo tene kâ fadë mbï fa na lo wüngö tî nginza nî pëpe:


sentence             [3190:3214] = — Me mo kö daä ma!


sentence             [3214:3258] = Mbï kö na yâ tî takisï nî na lörö.


sentence             [3258:3274] = — Awe, mâsa.

sentence             [3274:3283] = Ë gue.


sentence             [3283:3441] = Mbï wara lêgë tî dutï nzönî pëpe ngbanga tî sô na pöpö tî âzo otâ sô na pêko tî takisï nî sô, mbênî wâlï akono gï töngana bâgara.

sentence             [3441:3603] = Tîtene mbï dutï nzönî, mbï pûsu lo kêtê me wâlï sô asûku ngbângbâ, lo pete mbï na sïönî bê na terê tî yângâ tî takisï ni, sï lo tene:


sentence             [3603:3673] = — Mâsa, takisï tî mo nî ayeke kêtê töngasô ngbanga tî nye?

sentence             [3673:3839] = Bâa yê, gï na lëkërëngö têre tî mbï tîtene maseka-kôlï sô gerê tî lô akpa gï tî kôlô sô adutï nzönî, yâ tî terê tî mbï kûê ayeke so.


sentence             [3839:3901] = Wagbötöngö-takisï nî atene na ngîâ na yângâ tî lo:


sentence             [3901:3989] = — Mo hînga fadë, mamâ, mo fûta ngêrë tî ndo ûse tîtene mo bâa pâsi pëpe!


sentence             [3989:4047] = Wâlï sô akpa gï doli sô atï na wâ agbï, lo tene:


sentence             [4047:4085] = — Âla mä tënë tî kpändä sô!

sentence             [4085:4127] = Lo hînga mbï äpe, mbï hînga lo äpe.

sentence             [4127:4238] = Lo yeke tene tënë na mbï mo bâa mo tene mbï yeke mbênî kötä konô sô asï yä tî takisï nî kûê.

sentence             [4238:4263] = Mä mbï sï, môlengê.

sentence             [4263:4309] = Mbï lîngbi tî dutï töngana mamâ tî mo.

sentence             [4309:4396] = Gï na särängö terê tî âla töngasô, âla âmaseka, sï âla ke kûi hîo sô.


sentence             [4396:4463] = Ndo adë ngii sï pakara sô na gatï tî mamâ sô azî yângâ:


sentence             [4463:4583] = — Mbï yeke hûnda terê tî mbï, mbï tene, âla wara sïönî mayëre tî zîa âzo osïö na pekö sô na ndo wa?

sentence             [4583:4639] = Takisï tî mo nî ayeke kêtê ahön kopo tî âsusu.


sentence             [4639:4730] = Wagbötöngö-kutukutu nî azî yângâ pëpe tënë tî bîâ sô lo yeke hë na gbegô.


sentence             [4730:4798] = Kôlï sô na terê tî wâlï sô na kotï tî mbï akîri atene:


sentence             [4798:4951] = — Na yâ tî âmbênî ködörö, na tâpandë tî Farânzi sô mbï ke gue daä töngana lêgë tî yäkä tî mbï, âtakisï ayeke mû gï zo ôko.

sentence             [4951:5012] = Na ndönî, âla ayeke zîa lo na ndo sô lo yeke gue daä.


sentence             [5012:5151] = Gï sô lo tene töngasô sô âkperekpere nî alöndö ngangü ahön ndö tî îâ sô ayeke toto ngâ na kötä köngö tî takisï ni.


sentence             [5151:5285] = Masia sô adutï na pöpö tî Wagbötöngö-takisï na mbênî maseka sô amo li tî lo ndâli tî kötä kä sô ayeke daä atene:


sentence             [5285:5427] = — Na mbâgë, tënë tî mo nî ayeke na lêgë nî, me na mbâgë ngâ, alingbi mo hîinga ë dutï nzönî na yâ tî takisï ahön bus.


sentence             [5427:5504] = Môlengê-kôlï sô na kä na li tî lo sô atene na môlengê-wâlï nî:


sentence             [5504:5543] = — Tâtënë laä mo tene sô, îtä.

sentence             [5543:5603] = Mbï sô mo bâa mbï sô, mbï yeke mû bus lâôko pëpe.

sentence             [5603:5649] = Mbï yeke zo sô ayê dütïngö kpô mîngi.

sentence             [5649:5808] = Andâa, na yâ tî bus, zo kûê ayeke tene tënë, ayeke dë köngö, ayeke hë ngîâ na yâ tî mê tî âzo, ngâ ayeke sösö na mbänä tî bê tî lo.


sentence             [5808:5911] = Kôlï tî hümä sô adutï na mbâgë kâ mo bâa mo tene lo löndö na ködörö wandê sô atene:


sentence             [5911:5998] = —Tî mbï, mbï kpë bus ngbanga tî sô mbï yê fünngö tî terê tî âzo pëpe.

sentence             [5998:6167] = Mbï lîngbi tî kânga bê tî mbï na ndö tî fünngö gbe tî mabôko tî âmötöröbôi pëpe, wala ngâ na gbikï tî âla sô ayeke buba bongö tî mbï pëpe.


sentence             [6167:6233] = Ka mbï, âla yeke hûnda terê tî âla, mbï tene tî mbï nye?

sentence             [6233:6262] = Mbï zî yângâ ôko pëpe.

sentence             [6262:6339] = Sô fadë mbï kö na yâ tî takisï sô, sösö ahon ndö tî mbï mîngi.

sentence             [6339:6414] = Mbï yê tî sösö na ngangü tî terê tî mbï kûê sï pupu aya mbï.

sentence             [6414:6485] = Na Bangî, mo sösö, ayeke hë mo ngbii asï na lâ tî kûâ tî mo.

sentence             [6485:6583] = Ka mbï sösö nî sï sasa asïgîgî agbû bongö tî mbï, fadë mbï yeke sâra töngana nye?

sentence             [6583:6665] = Alîngbi mbï dutï kpô tî pete sösö sô ayê gï tî sungba na ngangü sô.


sentence             [6665:6725] = Wagbötöngö-takisï nî azî gô tî fonôno nî ngangü.

sentence             [6725:6867] = Töngana mbênî yê tî ndimä, sösö sô mbï kânga fadë gbä sô asïgî gï yeke bâa mo tene mozoko nî ayê däa tîtene lo sïgî.


sentence             [6867:6950] = Wâlï sô na terê tî mbï sï akpa gï kötä bäbolo sô atene na kötä gô:


sentence             [6950:7021] = — Hêe, mâsa, mo pêe tî kîri na gô tî fonôno-nî kêtê äpe?

sentence             [7021:7087] = Ë yeke na yâ tî nganda äpe, ë yeke na yâ tî takisï laâ!


sentence             [7087:7130] = Âdëköngö alöndö mbâgë na mbâgë:


sentence             [7130:7168] = — Nye laâ afün töngasô, mâsa?


sentence             [7168:7228] = — Mo yö fünngö nyama na yâ tî takisï nî wala nye?


sentence             [7228:7263] = — Kpändä wa laâ asösö sô?


sentence             [7263:7352] = — Mozoko sô ayeke toto ngangü mîngi sô laâ amû lêgë na zo sô tî sösö sô.


sentence             [7352:7407] = — Sösö nî afün gï töngana kûâ tî mbo, dê!


sentence             [7407:7460] = — Alîngbi lo gue abâa yâ tî purû tî lo sï.


sentence             [7460:7527] = Kötä wâlï sô na terê tî mbï sï amû ndo kûê sô atene:


sentence             [7527:7579] = — Mo bâa sô, ayeke gï mâsa laâ asösö sô.


sentence             [7579:7644] = Me wagbötöngö-takisï nî akîri tënë na wâlï nî atene:


sentence             [7644:7689] = — Tî mbï ayeke gï mo, pûpûlenge waâ!

sentence             [7689:7771] = Mbï hînga ka sôsô tî âwâlï sô ngbondä tî âla akono mîngi sô kûê.


sentence             [7771:7869] = Lê tî wâlï sô abe na ngonzo sï lo hûnda tîtene mâsa nî azîa lo na sêse fadëfadësô.

sentence             [7869:7935] = Ë yeke na barâsi balëokü na sängbï-lêgë tî Ben-Zvi-Sica.

sentence             [7935:7983] = Takisï nî alütï na dëmängö terê tî lo.

sentence             [7983:8052] = Wâlï nî asïgîgiî na kugbë tî sâki ôko lo mû na mâsa nî.

sentence             [8052:8174] = voazîne tî mbï nî, sô ayeke voazîne tî mbï mbênî pëpe tëtî sô lo zûu na sêse awe sô, ake na kürü gô:


sentence             [8174:8236] = — Mo dîko tanga tî nginza sô mo kîri na mbï sô laâ?


sentence             [8236:8291] = — Mbï kîri na mo ngbangbo na balë mbârâmbârâ.

sentence             [8291:8339] = Ka tënë nî ayeke na ndo wa, mamâ tî mbï?


sentence             [8339:8418] = — Mbï tene na mo fadë mbï tene mbï yeke gï na balëûse na okü, äpe?


sentence             [8418:8501] = — Mbï yeke mû ka balëûse na okü töngana mo yeke gï na balëûse na okü.


sentence             [8501:8664] = Mâsa nî azî gerê, takisï tî lo nî ahön yongôro na terê tî wâlï tî gïngö tënë sô sï lo tene nï hînga li tî nï na âzo tî Bêafrîka pëpe.

sentence             [8664:8701] = Ngêrë tî takisï ayeke balëotâ.

sentence             [8701:8821] = Ayeke na nzönî bê laâ sï lo yeke yô âzo na balëûse na okü, sô ayeke ngêrë tî köngö na yâ tî bus sô.


sentence             [8821:8903] = Töngana kutukutu nî asï na dambëtï tî Ben-Zvi awe, mbï zûu na sêse ngâ.

sentence             [8903:9028] = Na mbâgë tî lêgë na dawäa tî mbênî dapïkängö-singa, Sylvianne ayeke kü mbï töngana tî sô ë mä terê daä.

sentence             [9028:9134] = Mbï gue na terê tî lo sï lêgëôko, ë mû lêgë sô ague na Centre Baobab tî kîri na yângâda.


sentence             [9134:9162] = Lo hûnda mbï tî hînga:


sentence             [9162:9215] = — Büngbïngö-terê nî atambûla na lêgë nî?


sentence             [9215:9224] = — Iin.

sentence             [9224:9304] = Ahûnda na ë tî bûngbi tî sâra mbênî bûku na farânzi ngâ na sängö.


sentence             [9304:9373] = — Mo hînga yê sô mo yê tî sîra mbëtï na ndö nî sô awe?


sentence             [9373:9450] = — Mbï yê tî vunga sango tî yê sô asï fadë na yâ tî takisï sô.


sentence             [9450:9474] = — Nye laâ asï nye?


sentence             [9474:9578] = —Âzo tî yâ tî takisï nî adîko tënë töna âkûkulû sï mbï sösö tîtene âla dutï kpo.


sentence             [9578:9651] = Sylvianne atï na kötä ngîâ sï nzara tî mbênî sösö agbû mbï.

sentence             [9651:9672] = Me mbï kânga nî.


sentence             [9672:9714] = — Ahînga atene mo laâ mo sösö sô?


sentence             [9714:9791] = — Mbï sösö nî na ngoi sô mâsa nî azî gô tî fonôno nî ngangü.


sentence             [9791:9843] = — Tî taâ tënë nî, atumba fünngö nî äpe?


sentence             [9843:9925] = — Âhînga fadë atene mbï laâ mbï sösö sô, ka kamënë ayeke gbû mbï.


sentence             [9925:9966] = — Sösö ayeke yê tî kamënë pëpe!

sentence             [9966:10002] = Sösö afa atene zo ayeke na sêni.

sentence             [10002:10110] = Atâa töngana âsûru yâ tî zo awe sô, ayeke kü tîtene zo nî asösö awe, sï lo lîngbi tî te yê.

sentence             [10110:10169] = Âmunzû ayeke dema pëpe töngana mbâ tî âla asösö.


sentence             [10169:10287] = Na dabüzë tî Amadi, mbï vo gbâ tî mânga sô ïrï nî Menthol sô sï mbï zîa wâ na terê tî mbênî ôko.

sentence             [10287:10359] = Pialö tî süngö bûku na farânzi-sängö sô akîri na li tî mbï.

sentence             [10359:10464] = Tî särängö mbëtï na farânzi sï mbï kîri mbï lekere nî na sängö sô aso li tî mbï mîngi.

sentence             [10464:10584] = Töngana nye ânde sï mbï sâra mbëtï na yângâ tî ködörö tî mbï sô mbï hînga tî tene nzönî pëpe sô?

sentence             [10584:10692] = Âfa nî fadë na mbî na dambëtï, mbï lîngbi tî wara kpälë tî särängö mbëtï na sängö pëpe.

sentence             [10692:10715] = Me mawa nî ayeke sô!

sentence             [10715:10953] = Gï na fängö na mbï tî tene ngâ na tî sâra mbëtï nzönî na yângâ tî ködörö tî âzo sô amû na ë Dipandäa, âzo tî tënëngö kötä farânzi tî ködörö tî ë asâra sï ë hînga tënëngö sängö nzönî pëpe.

sentence             [10953:11024] = Âmbupa tî farânzi asï sïngö na yâ tî sängö sô ë yeke tene.

sentence             [11024:11172] = Mbï yeke tö ndâ tî sâra mbëtï na marä sängö sô, ngbanga tî sô nzönî ködë tî manda yângâ tî ködörö, ayeke gï tî tene nî.


sentence             [11172:11228] = Mo yê tî bâa, sîngamê tî mbï tî mabôko atoto.


sentence             [11228:11239] = — Alô?


sentence             [11239:11255] = Mbï ma atene:


sentence             [11255:11281] = — Töngana nye, petï?


sentence             [11281:11298] = — Zo wa laâ?


sentence             [11298:11330] = — Mo girisa gô tî mbï awe?

sentence             [11330:11351] = Mbï Yvon laâ äpe.

sentence             [11351:11371] = Mo yeke na ndo wa?


sentence             [11371:11409] = — Mbï yeke gue na yângâda laâ.


sentence             [11409:11477] = — Mo pêe tî gä hîo na ndo sô ë ike wara ka terê daä sô?


sentence             [11477:11598] = — Mbï pêe, me na yâ tî ngbonga ôko ngbanga tî sô mbï yê tî sû mbênî mbëtï na farânzi ngâ na sängö.


sentence             [11598:11637] = — Alîngbi tî kü kêkerêke äpe?


sentence             [11637:11741] = — Mbï yê tî tö ndâ nî gï fafadësô, töngasô äpe sô, mbï lîngbi tî sû nî ânde äpe.


sentence             [11741:11853] = — Mbï lîngbi tî kü gï töngana mo yê daä tî sâra tënë tî mbï na yâ tî mbëtï tî mo nî sï.


sentence             [11853:11900] = Mbï tene na lo na ngîâ na yângâ tî mbï.

sentence             [11900:11901] = 	
`
	if actualStr != expectStr {
		t.Errorf("bad Parse:\nfound %v but expected %v\n", actualStr, expectStr)
	}
}
