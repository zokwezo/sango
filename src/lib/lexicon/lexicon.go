// Sango dictionary, stored in row-major and column-major order as Go global variables.
//
// Copyright 2024 Daniel D. Weston
// Use of this source code and data is governed by http://www.apache.org/licenses/LICENSE-2.0
// a copy of which can be found in the LICENSE file.

package lexicon

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// Sango lexicon in row-major and column-major order, outer joined with compatible affixes.
func LexiconRows() DictRows                     { return lexiconRowsAndCols.rows }
func LexiconCols() DictCols                     { return lexiconRowsAndCols.cols }
func HeightlessFromLemma() map[string]string    { return lexiconRowsAndCols.heightlessFromLemma }
func TonelessFromHeightless() map[string]string { return lexiconRowsAndCols.tonelessFromHeightless }
func RowsMatchingToneless() DictRowsMap         { return lexiconRowsAndCols.rowsMatchingToneless }
func RowsMatchingHeightless() DictRowsMap       { return lexiconRowsAndCols.rowsMatchingHeightless }
func RowsMatchingLemma() DictRowsMap            { return lexiconRowsAndCols.rowsMatchingLemma }

type DictRowsMap = map[string]DictRows

type DictRows = []DictRow

type DictRow struct {
	Toneless           string // = Heightless, but with pitch removed
	Heightless         string // = Lemma, but with whitespace, case, and height removed
	Lemma              string
	UDPos              string
	UDFeature          string
	Category           string
	Frequency          int
	EnglishTranslation string
	EnglishDefinition  string
}

type DictCols struct {
	Bytes              []byte
	Runes              []rune
	Toneless           [][]byte
	Heightless         [][]byte
	LemmaRunes         [][]rune
	LemmaUTF8          [][]byte
	UDPos              [][]byte
	UDFeature          [][]byte
	Category           [][]byte
	Frequency          []int
	EnglishTranslation [][]byte
	EnglishDefinition  [][]byte
}

type DictRowRegexp struct {
	TonelessRE           *regexp.Regexp
	HeightlessRE         *regexp.Regexp
	LemmaRE              *regexp.Regexp
	UDPosRE              *regexp.Regexp
	UDFeatureRE          *regexp.Regexp
	CategoryRE           *regexp.Regexp
	EnglishTranslationRE *regexp.Regexp
	EnglishDefinitionRE  *regexp.Regexp
	FrequencyMin         int
	FrequencyMax         int
}

func Lookup(dictRows DictRows, dictRowRegexp DictRowRegexp) (DictRows, error) {
	return lookupMatchingRows(dictRows, dictRowRegexp)
}

//////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func lookupMatchingRows(in DictRows, f DictRowRegexp) (out DictRows, err error) {
	for _, r := range in {
		if r.Toneless != "" &&
			f.FrequencyMin <= r.Frequency &&
			f.FrequencyMax >= r.Frequency &&
			f.TonelessRE.MatchString(r.Toneless) &&
			f.HeightlessRE.MatchString(r.Heightless) &&
			f.LemmaRE.MatchString(r.Lemma) &&
			f.UDPosRE.MatchString(r.UDPos) &&
			f.UDFeatureRE.MatchString(r.UDFeature) &&
			f.CategoryRE.MatchString(r.Category) &&
			f.EnglishTranslationRE.MatchString(r.EnglishTranslation) &&
			f.EnglishDefinitionRE.MatchString(r.EnglishDefinition) {
			out = append(out, r)
		}
	}
	return out, nil
}

type dictRowsAndCols struct {
	rows                   DictRows
	cols                   DictCols
	heightlessFromLemma    map[string]string
	tonelessFromHeightless map[string]string
	rowsMatchingToneless   DictRowsMap
	rowsMatchingHeightless DictRowsMap
	rowsMatchingLemma      DictRowsMap
}

var lexiconRowsAndCols = func() dictRowsAndCols {
	var rows = DictRows{
		{"", "", "", "DO NOT REMOVE THIS ROW", "Copyright=DanielDWeston2024", "HTTP://WWW.APACHE.ORG/LICENSES/LICENSE-2.0", 0, "https://github.com/zokwezo/sango/blob/main/src/lib/lexicon/lexicon.csv", ""},
		{"ababaa", "ababâa", "ababâa", "NOUN", "", "FOOD", 5, "soybean", "soybean"},
		{"ade", "âde", "âdɛ", "ADV", "Mood=Irr", "WHEN", 1, "not-yet", "[lit: if-only|remain]: not yet"},
		{"adu", "âdu", "âdu", "SCONJ", "Mood=Irr", "HOW", 1, "if", "[lit: if-only|exist]: if only"},
		{"ae", "âe", "âe", "INTERJ", "", "SENSE", 3, "oh", "Ouch! Oh!"},
		{"afirika", "afirîka", "afirîka", "NOUN", "", "COUNTRY", 1, "Africa", "Africa"},
		{"ahon", "ahön", "ahön", "ADP", "VerbForm=Fin", "HOW", 2, "more", "more than"},
		{"ahonkue", "ahön kûê", "ahön kûɛ̂", "ADV", "VerbForm=Fin", "HOW", 2, "most", "[lit: more than|all]: most of all"},
		{"ahonndoni", "ahön ndö nî", "ahön ndö nî", "ADV", "VerbForm=Fin", "HOW", 2, "too-much", "too much"},
		{"ai", "âi", "âi", "INTERJ", "", "ALT SP FOR", 9, "oh", "âe"},
		{"akotara", "âkötarä", "âkötarä", "NOUN", "Num=Plur", "FAMILY", 4, "genealogy", "[plural] genealogy, family tree, patrilineage"},
		{"ala", "âla", "âla", "PRON", "Num=Plur|Person=3|PronType=Prs", "WHO", 1, "they-or-you", "they, them; you [plural and/or polite]"},
		{"alameti", "alamëti", "alamëti", "NOUN", "", "ALT SP FOR", 9, "matchstick", "alimëti"},
		{"alamveni", "âla mvenî", "âla mvɛnî", "PRON", "Num=Plur|Person=3|PronType=Prs", "WHO", 1, "themselves", "[lit: they, [plural and/or polite] you|self]: themselves"},
		{"ale", "alë", "alë", "NOUN", "VerbForm=Fin", "WHO", 3, "ancestry", "[lit: it bears fruit]: ancestry"},
		{"alezo", "alë-zo", "alë-zo", "NOUN", "VerbForm=Fin", "WHO", 3, "ancestor", "[lit: it bears fruit|person]: ancestor(s)"},
		{"alimeti", "alimëti", "alimëti", "NOUN", "", "OBJ", 7, "matchstick", "[Fr: allumette]: match (for lighting a fire)"},
		{"ambeso", "âmbeso", "âmbeso", "NOUN", "Num=Plur", "WHO", 2, "ancestors", "[lit: [plural]|formerly]: ancestors"},
		{"ambii", "ambïi", "ambïi", "NOUN", "", "PLANT", 6, "skin-whitener", "skin whitener (for women)"},
		{"amerika", "amerîka", "amerîka", "NOUN", "", "COUNTRY", 1, "America", "America"},
		{"anarumbatee", "ânarûmbatêe", "ânarûmbatêe", "NOUN", "", "CIVIL", 4, "patty-cake", "patty cake game"},
		{"andaa", "andâa", "andâa", "SCONJ", "Mood=Irr", "HOW", 1, "however", "[lit: if-only|end]: however"},
		{"ande", "ânde", "ândɛ", "ADV", "Mood=Irr", "WHEN", 1, "later", "[lit: if-only|next-day]: later"},
		{"ando", "ândö", "ândö", "ADV", "Mood=Irr", "WHEN", 1, "recently", "[lit: if-only|at-the-place]: recently"},
		{"ange", "ânge", "ânge", "INTERJ", "", "SENSE", 1, "watch-out", "[lit: hön=pass|ge=here]: watch out"},
		{"angelee", "angelêe", "angɛlɛ̈ɛ", "NOUN", "", "COUNTRY", 1, "England", "England"},
		{"angoro", "angöro", "angɔ̈rɔ", "ADV", "", "HOW", 7, "again", "[Fr: encore]: once more, again"},
		{"ani", "ânî", "ânî", "PRON", "Animacy=Inan|Case=Acc|Num=Plur|Person=3|PronType=Det", "WHICH", 1, "they", "[neuter]: they, them"},
		{"ani", "ânï", "ânï", "PRON", "Num=Plur|Person=3|PronType=Rel", "WHICH", 2, "they", "[indirect style]: they"},
		{"aparandee", "aparandëe", "aparandëe", "NOUN", "", "WHO", 7, "apprentice", "[Fr: apprenti]: apprentice"},
		{"ape", "äpe", "äpɛ", "PART", "Polarity=Neg", "HOW", 1, "not", "not"},
		{"ara", "ara", "ara", "VERB", "", "MOVE", 4, "crawl", "crawl"},
		{"arabu", "arâbu", "arâbu", "NOUN", "", "COUNTRY", 1, "Arab", "Arab"},
		{"arara", "arara", "arara", "NOUN", "", "OBJ", 4, "umbrella", "umbrella"},
		{"asa", "asa", "asa", "VERB", "Subcat=Tran", "ACT", 4, "scratch", "scratch"},
		{"asina", "asî na", "asî na", "ADP", "", "WHEN", 1, "until", "[lit: arrive|at]: until"},
		{"asinaawe", "asî na ... awe", "asî na ... awe", "ADP", "", "WHEN", 1, "by", "[lit: arrive|at|...|by]: by"},
		{"asingana", "asî ngâ na", "asî ngâ na", "ADP", "", "WHEN", 1, "with-respect-to", "[lit: arrive|also|at]: with respect to"},
		{"ata", "âta", "âta", "NOUN", "", "FAMILY", 4, "maternal-grandrelative", "maternal grandchild; maternal grandparent"},
		{"ataa", "atâa", "atâa", "SCONJ", "Mood=Irr", "HOW", 1, "even-if", "[lit: if-only|true] even (if)"},
		{"ataaso", "atâa sô", "atâa sô", "SCONJ", "", "HOW", 1, "although", "although"},
		{"au", "aû", "aû", "NOUN", "", "ALT WORD FOR", 9, "maternal-relative", "kôya"},
		{"awane", "awâne", "awâne", "NOUN", "", "ANIM", 6, "snail", "snail"},
		{"awe", "awe", "awɛ", "PART", "", "WHEN", 1, "already", "already"},
		{"ayi", "âyi", "âyi", "INTERJ", "", "ALT SP FOR", 9, "oh", "âe"},
		{"ba", "ba", "ba", "VERB", "Subcat=Tran", "ACT", 3, "bend", "fold, bend, twist, roll, undulate, zigzag"},
		{"ba", "bâ", "bâ", "NOUN", "", "WHERE", 1, "silver", "machete; silver"},
		{"ba", "bä", "bä", "NOUN", "", "INTERACT", 3, "oath", "oath"},
		{"ba", "bä", "bä", "NOUN", "", "WHERE", 3, "foundation", "foundation"},
		{"baa", "bâa", "bâa", "VERB", "", "SENSE", 1, "see", "see, look; sense, perceive; understand; meet; experience"},
		{"baamotene", "bâa mo tene", "bâa mo tɛnɛ", "SCONJ", "", "SENSE", 1, "as-if", "[lit: see|you|say]: as if"},
		{"baanga", "bâanga", "bâanga", "VERB", "Aspect=Iter", "SENSE", 3, "examine", "examine"},
		{"baaya", "bâa yä", "bâa yä", "VERB", "Subcat=Intr", "", 5, "menstruate", "menstruate"},
		{"baba", "baba", "baba", "NOUN", "", "SENSE", 4, "vanity", "haughty pride, vanity"},
		{"baba", "babâ", "babâ", "NOUN", "Gender=Masc", "FAMILY", 2, "father", "father"},
		{"baba", "bâbâ", "bâbâ", "NOUN", "", "CIVIL", 5, "initiation-camp", "[lit: place-for|oath]: initiation camp; club; social group"},
		{"babala", "babâla", "babâla", "NOUN", "", "ALT SP FOR", 9, "boulevard", "balabâla"},
		{"babango", "babango", "babango", "NOUN", "", "BODY", 5, "index-finger", "index finger"},
		{"babolo", "bäbolo", "bäbɔlɔ", "NOUN", "", "FOOD", 5, "yam", "sweet potato, yam"},
		{"baboro", "bäboro", "bäbɔrɔ", "NOUN", "", "ALT SP FOR", 9, "yam", "bäbɔlɔ"},
		{"bada", "badâ", "badâ", "NOUN", "", "ANIM", 5, "squirrel", "squirrel"},
		{"bada", "bâda", "bâda", "NOUN", "", "HOUSE", 5, "sanctuary", "[lit: place-for|shelter]: sanctuary"},
		{"badabuku", "bâdabûku", "bâdabûku", "NOUN", "", "HOUSE", 5, "library", "[lit: place-for|shelter|book]: library"},
		{"badahalezo", "bâdahalëzo", "bâdahalëzo", "NOUN", "", "HOUSE", 5, "legislature", "[lit: place-for|shelter|it bears fruit|person]: legislature, National Assembly"},
		{"bagara", "bâgara", "bâgara", "NOUN", "", "ANIM", 2, "cow", "cow"},
		{"bagbara", "bagbara", "bagbara", "NOUN", "", "HOUSE", 3, "bridge", "bridge"},
		{"bage", "bâge", "bâge", "NOUN", "", "OBJ", 7, "finger-ring", "(finger) ring"},
		{"bahule", "bâhülë", "bâhülë", "NOUN", "", "HOW", 2, "stadium", "[lit: place-where|sport]: playing field, stadium"},
		{"bakale", "bakalê", "bakalê", "NOUN", "", "MYTH", 6, "genie", "omniscient genie"},
		{"bakari", "bakarî", "bakarî", "NOUN", "", "OBJ", 4, "dictionary", "dictionary, encyclopedia"},
		{"bake", "bâke", "bâke", "NOUN", "", "CIVIL", 7, "rattan-stool", "balambo"},
		{"bake", "bâke", "bâke", "NOUN", "", "OBJ", 7, "ferry", "ferry, barge"},
		{"bakongo", "bäkongö", "bäkɔngɔ̈", "NOUN", "", "ANIM", 3, "turtle", "turtle"},
		{"bakoya", "bäkoyä", "bäkoyä", "NOUN", "", "ANIM", 3, "baboon", "baboon"},
		{"bakpa", "bäkpä", "bäkpä", "NOUN", "", "FOOD", 5, "nut-butter", "paste of ground nuts"},
		{"bakuru", "bäkürü", "bäkürü", "NOUN", "", "ANIM", 6, "kite", "yellow-billed kite"},
		{"bakutu", "bakûtu", "bakûtu", "NOUN", "", "ACT", 4, "water-drumming", "water drumming"},
		{"bala", "bala", "bala", "NOUN", "", "INTERACT", 1, "greetings", "[lit: see|repeatedly]: greeting"},
		{"bala", "bala", "bala", "VERB", "Subcat=Tran", "INTERACT", 1, "greet", "[lit: see|repeatedly]: greet"},
		{"balabala", "balabâla", "balabâla", "NOUN", "", "CIVIL", 4, "boulevard", "boulevard, avenue, highway"},
		{"balaka", "balaka", "balaka", "NOUN", "", "OBJ", 3, "Balaka", "machete"},
		{"balama", "balama", "balama", "INTERJ", "Mood=Opt", "INTERACT", 1, "Hello", "[enthusiastic]: Hello!"},
		{"balambo", "balambo", "balambo", "NOUN", "", "HOUSE", 4, "stool", "[lit: greet|dog]: rattan stool"},
		{"balangeti", "balangëti", "balangɛ̈ti", "NOUN", "", "OBJ", 7, "blanket", "[En: blanket]: blanket"},
		{"balao", "balaô", "balaɔ̂", "INTERJ", "Polte=Form", "INTERACT", 1, "Hello", "[polite]: Hello!"},
		{"balapaa", "balapâa", "balapâa", "NOUN", "", "TREE", 5, "breadfruit-tree", "[Fr: balle|pain]: Breadfruit tree (South Pacific)"},
		{"balawa", "bâlâwâ", "bâlâwâ", "NOUN", "", "TREE", 5, "shea-tree", "Shea tree (Karite, buttertree)"},
		{"bale", "bale", "bale", "NOUN", "", "NATURE", 2, "river", "large river, Oubangui river"},
		{"bale", "balë", "balë", "ADJ", "NumType=Ord", "NUM", 2, "ten", "ten"},
		{"bale", "balë", "balë", "NUM", "NumType=Card", "NUM", 2, "ten", "ten"},
		{"balee", "balêe", "balêe", "VERB", "", "ACT", 3, "sweep", "[Fr: balet]: sweep"},
		{"balee", "balëe", "balëe", "NOUN", "", "OBJ", 3, "broom", "[Fr: balet]: broom"},
		{"bamara", "bämarä", "bämarä", "NOUN", "", "ANIM", 3, "lion", "lion"},
		{"bambi", "bambî", "bambî", "NOUN", "", "FAMILY", 3, "baby", "[<6mo]: newborn"},
		{"bambinga", "bambinga", "bambinga", "NOUN", "", "WHO", 4, "Pygmy", "Pygmy"},
		{"bambu", "bambü", "bambü", "NOUN", "", "TREE", 5, "bamboo", "[Fr: bambou]: bamboo"},
		{"baminga", "baminga", "baminga", "NOUN", "", "ALT SP FOR", 9, "Pygmy", "bambinga"},
		{"bandeko", "bandeko", "bandeko", "NOUN", "", "ALT SP FOR", 9, "adultery", "ndeko"},
		{"bandembo", "bândembö", "bândembö", "NOUN", "", "HOW", 3, "soccer-field", "[lit: place-where|ball]: soccer field"},
		{"bando", "bandö", "bandö", "NOUN", "", "FOOD", 6, "tripe-fat", "tripe fat"},
		{"banga", "banga", "banga", "NOUN", "", "WHERE", 2, "north", "north; right riverbank"},
		{"banga", "bânga", "bânga", "NOUN", "", "BODY", 4, "chin", "chin"},
		{"banga", "bängâ", "bängâ", "NOUN", "", "TREE", 3, "rubber-tree", "rubber tree"},
		{"bangbi", "bângbi", "bângbi", "NOUN", "", "SENSE", 3, "supervision", "supervision"},
		{"bangbi", "bângbi", "bângbi", "VERB", "Aspect=Imp|Subcat=Tran", "SENSE", 3, "supervise", "supervise"},
		{"bangi", "bangi", "bangi", "NOUN", "", "TREE", 5, "iroko-tree", "iroko"},
		{"bangi", "bangî", "bangî", "NOUN", "", "WHERE", 1, "Bangui", "Bangui"},
		{"bangi", "bângi", "bângi", "NOUN", "", "PLANT", 5, "hemp", "hemp"},
		{"bangu", "bângû", "bângû", "NOUN", "", "NATURE", 3, "trough", "[lit: place|water]: watering hole, trough"},
		{"bao", "bäö", "bäö", "NOUN", "", "ANIM", 6, "python", "python"},
		{"bara", "bara", "bara", "NOUN", "Aspect=Hab", "ALT SP FOR", 9, "greetings", "bala"},
		{"bara", "bara", "bara", "VERB", "Aspect=Hab|Subcat=Tran", "ALT SP FOR", 9, "greet", "bala"},
		{"barama", "barama", "barama", "INTERJ", "Mood=Opt", "ALT SP FOR", 9, "Hello", "balama"},
		{"baramii", "baramïi", "baramïi", "NOUN", "", "OBJ", 7, "crow-bar", "[Fr: barre à mine]: crow bar, pry bar, digging stick"},
		{"barao", "baraô", "baraɔ̂", "INTERJ", "Polite=Form", "ALT SP FOR", 9, "Hello", "balao"},
		{"basenzi", "basënzi", "basɛ̈nzi", "NOUN", "", "CIVIL", 4, "traditional", "traditional, African style, savage"},
		{"bata", "bata", "bata", "VERB", "Aspect=Hab|Subcat=Tran", "ACT", 2, "guard", "guard"},
		{"batoo", "batöo", "batöo", "NOUN", "", "OBJ", 7, "boat", "[Fr: bateau]: boat"},
		{"bawere", "bâwërë", "bâwërë", "NOUN", "", "HOW", 2, "stadium", "[lit: place-where|sport]: playing field, stadium"},
		{"baya", "bâyâ", "bâyâ", "VERB", "", "CIVIL", 4, "repay", "repay, pay off, absolve oneself of"},
		{"bazingere", "bazïngêre", "bazïngêre", "NOUN", "", "WHO", 4, "apostle", "raider, apostle"},
		{"be", "be", "bɛ", "VERB", "Subcat=Intr", "STATE", 3, "ripen", "ripen, ripen, turn red or brown"},
		{"be", "bê", "bɛ̂", "NOUN", "", "BODY", 1, "heart", "heart, liver, mind"},
		{"be", "bê", "bɛ̂", "NOUN", "", "WHERE", 1, "center", "center, middle"},
		{"be", "bë", "bë", "VERB", "Subcat=Tran", "SENSE", 3, "embarrass", "weigh down, bother, embarrass"},
		{"beafrika", "bê-afrîka", "bɛ̂-afrîka", "NOUN", "", "COUNTRY", 1, "Central-Africa", "Central Africa"},
		{"bebee", "bebëe", "bebëe", "NOUN", "", "FAMILY", 7, "baby", "[Fr: bébé]: baby"},
		{"bebi", "bê-bï", "bɛ̂-bï", "NOUN", "", "WHEN", 1, "midnight", "[lit: middle|night]: midnight"},
		{"bekani", "bekâni", "bekâni", "NOUN", "", "OBJ", 7, "bicycle", "[Fr: bécane]: bicycle"},
		{"bekodoro", "bê-ködörö", "bɛ̂-kɔ̈dɔ̈rɔ̈", "NOUN", "", "CIVIL", 2, "urban", "in the city, urban"},
		{"bekombite", "bê-kömbïte", "bɛ̂-kɔ̈mbïtɛ", "NOUN", "", "WHEN", 1, "noon", "[lit: middle|noon]: noon"},
		{"bekpa", "bëkpä", "bëkpä", "NOUN", "", "NATURE", 3, "thunder", "thunder"},
		{"bela", "bê-lâ", "bɛ̂-lâ", "NOUN", "", "WHEN", 1, "midday", "[lit: middle|day]: midday"},
		{"belaawu", "bêlâawü", "bɛ̂lâawü", "NOUN", "", "WHEN", 5, "May", "May"},
		{"bele", "bele", "bele", "VERB", "Aspect=Hab|Subcat=Intr", "ALT SP FOR", 9, "squat", "bere"},
		{"bele", "bele", "bɛlɛ", "VERB", "Aspect=Hab|Subcat=Tran", "ALT SP FOR", 9, "deny", "bɛrɛ"},
		{"bele", "bêle", "bɛ̂lɛ", "NOUN", "Aspect=Hab", "ALT SP FOR", 9, "envy", "bɛ̂rɛ"},
		{"belebele", "belebele", "belebele", "VERB", "|Aspect=Hab|Mood=Emp", "ALT SP FOR", 9, "soaked", "berebere"},
		{"belu", "belü", "belü", "NOUN", "", "ANIM", 5, "porcupine", "porcupine"},
		{"bema", "bema", "bema", "VERB", "Subcat=Intr", "INTERACT", 4, "complain", "moan, whine, complain"},
		{"benda", "benda", "bɛnda", "NOUN", "", "ACT", 6, "win", "win, victory"},
		{"bendambo", "bê-ndâmbo", "bɛ̂-ndâmbo", "NOUN", "", "NUM", 2, "quarter", "[lit: middle|half]: quarter"},
		{"bengba", "bengbä", "bengbä", "ADJ", "", "COLOR", 3, "red", "red"},
		{"bengbabengba", "bengbä-bengbä", "bengbä-bengbä", "ADJ", "Mood=Emp", "COLOR", 3, "bright-red", "bright red"},
		{"bengbakete", "bengbä-kêtê", "bengbä-kɛ̂tɛ̂", "ADJ", "", "COLOR", 3, "pink", "pink"},
		{"benge", "bëngë", "bɛ̈ngɛ̈", "NOUN", "", "PLANT", 6, "poison", "poison, strychnine"},
		{"bengo", "bëngö", "bɛ̈ngɔ̈", "ADJ", "VerbForm=Vnoun", "STATE", 3, "ripe", "reddish, ripe"},
		{"benyama", "bê-nyämä", "bɛ̂-nyämä", "NOUN", "", "CIVIL", 2, "rural", "countryside, rural"},
		{"bere", "bere", "bere", "VERB", "Aspect=Hab|Subcat=Intr", "INTERACT", 6, "squat", "crouch down, hunker down, squat, hide"},
		{"bere", "bere", "bɛrɛ", "VERB", "Aspect=Hab|Subcat=Tran", "INTERACT", 4, "deny", "deny"},
		{"bere", "berë", "berë", "ADV", "", "HOW", 6, "maybe", "maybe"},
		{"bere", "bêre", "bɛ̂rɛ", "NOUN", "Aspect=Hab", "SENSE", 6, "envy", "jealousy, envy"},
		{"berebere", "berebere", "berebere", "ADJ", "Aspect=Hab|Mood=Emp", "HOW", 6, "soaked", "soaked"},
		{"beredele", "beredële", "bɛrɛdɛ̈lɛ", "NOUN", "", "WHO", 7, "prostitute", "prostitute"},
		{"beta", "bêtâ", "bɛ̂tâ", "ADJ", "", "HOW", 6, "true", "[lit:heart|true]: true"},
		{"beta", "bëtä", "bëtä", "NOUN", "", "ANIM", 4, "waterbuck", "waterbuck, large antilope"},
		{"bezongo", "bê-zöngö", "bɛ̂-zöngö", "NOUN", "", "WHEN", 6, "September", "September"},
		{"bi", "bi", "bi", "VERB", "Subcat=Tran", "ACT", 1, "throw", "throw (away)"},
		{"bi", "bî", "bî", "NOUN", "", "ACT", 1, "throw", "throw"},
		{"bi", "bï", "bï", "NOUN", "", "WHEN", 2, "night", "night"},
		{"bi", "bï", "bï", "VERB", "Subcat=Tran", "ACT", 6, "tame", "tame, domesticate"},
		{"bia", "bîâ", "bîâ", "NOUN", "", "INTERACT", 3, "song", "song"},
		{"biaku", "bîakü", "bîakü", "ADV", "", "WHEN", 1, "immediately", "immediately"},
		{"bianga", "bîângâ", "bîângâ", "NOUN", "", "ANIM", 6, "frog", "[? from bîâ=song + yângâ=mouth]: frog"},
		{"biani", "bîanî", "bîanî", "ADV", "", "HOW", 1, "certainly", "certainly"},
		{"bibe", "bi bê", "bi bɛ̂", "VERB", "Subcat=Intr", "FEEL", 1, "consider", "[lit: throw|heart]: wonder (if), consider (whether), reflect"},
		{"bibe", "bi-bê", "bi-bɛ̂", "NOUN", "", "FEEL", 1, "consideration", "[lit: throw|heart]: thought (for), consideration, reflection"},
		{"bibila", "bibila", "bibila", "NOUN", "", "HOW", 4, "filth", "filth, dirtyness"},
		{"bibila", "bibila", "bibila", "VERB", "Subcat=Tran", "HOW", 4, "dirty", "dirty, soil"},
		{"biele", "bîêle", "bîɛ̂lɛ", "NOUN", "", "DRINK", 5, "beer", "beer"},
		{"bikua", "bïkua", "bïkua", "NOUN", "", "WHEN", 2, "weekday", "weekday"},
		{"bikuaoko", "bïkua-ôko", "bïkua-ɔ̂kɔ", "NOUN", "", "WHEN", 2, "Monday", "Monday"},
		{"bikuaoku", "bïkua-okü", "bïkua-ɔkü", "NOUN", "", "WHEN", 2, "Friday", "Friday"},
		{"bikuaota", "bïkua-otâ", "bïkua-otâ", "NOUN", "", "WHEN", 2, "Wednesday", "Wednesday"},
		{"bikuause", "bïkua-ûse", "bïkua-ûse", "NOUN", "", "WHEN", 2, "Tuesday", "Tuesday"},
		{"bikuausio", "bïkua-usïö", "bïkua-usïö", "NOUN", "", "WHEN", 2, "Thursday", "Thursday"},
		{"bilarizi", "bilarïzi", "bilarïzi", "NOUN", "", "SICK", 7, "schistosomiasis", "[Fr: bilharzie]: schistosomiasis"},
		{"bilibili", "bîlîbili", "bîlîbili", "NOUN", "", "DRINK", 4, "millet-beer", "millet beer"},
		{"binabe", "bi na bê", "bi na bɛ̂", "VERB", "Subcat=Tran", "FEEL", 1, "consider", "[lit: throw|in|heart]: wonder about, consider, reflect on"},
		{"bindi", "bindi", "bindi", "NOUN", "", "STATE", 5, "magic", "magic, sorcery"},
		{"bindi", "bindî", "bindî", "NOUN", "", "ANIM", 5, "locust", "locust"},
		{"binga", "bînga", "bînga", "VERB", "Aspect=Iter|Subcat=Tran", "ACT", 3, "disperse", "disperse"},
		{"bingbi", "bingbi", "bingbi", "NOUN", "", "INTERACT", 3, "discussion", "discussion"},
		{"bingbi", "bingbi", "bingbi", "VERB", "Aspect=Imp|Subcat=Tran", "INTERACT", 3, "discuss", "discuss"},
		{"bingbitere", "bingbi terê", "bingbi tɛrɛ̂", "VERB", "Subcat=Intr", "ACT", 3, "have-a-discussion", "[lit: discuss|oneself]: have a discussion"},
		{"bio", "biö", "biö", "NOUN", "", "BODY", 3, "bone", "bone"},
		{"bipatara", "bi patärä", "bi patärä", "VERB", "Subcat=Intr", "GAME", 4, "throw-the-dice", "throw the dice"},
		{"bira", "birâ", "birâ", "NOUN", "", "CIVIL", 4, "battle", "combat, battle"},
		{"biri", "bîrï", "bîrï", "ADV", "", "WHEN", 2, "yesterday", "yesterday"},
		{"biribiri", "bîrîbiri", "bîrîbiri", "NOUN", "", "ALT SP FOR", 9, "millet-beer", "bîlîbili"},
		{"biriki", "birîki", "birîki", "NOUN", "", "OBJ", 7, "brick", "[Fr: brique]: brick"},
		{"biritani", "biritâni", "biritâni", "NOUN", "", "COUNTRY", 1, "Britain", "Britain"},
		{"bisee", "bisêe", "bisêe", "VERB", "Subcat=Tran", "INTERACT", 7, "invite", "[Fr: inviter]: invite"},
		{"biyee", "biyëe", "biyëe", "NOUN", "", "INTERACT", 7, "ticket", "[Fr: billet]: ticket"},
		{"bo", "bô", "bô", "VERB", "", "ACT", 3, "gather", "gather"},
		{"bo", "bö", "bö", "VERB", "Subcat=Tran", "ACT", 4, "lapidate", "lapidate, throw stones at and hit"},
		{"bobo", "bobo", "bobo", "NOUN", "", "ANIM", 4, "termite", "winged worker termite"},
		{"boi", "bôi", "bôi", "NOUN", "", "WHO", 4, "houseboy", "houseboy, domestic servant"},
		{"boingu", "bôingû", "bôingû", "NOUN", "", "ANIM", 4, "ringworm", "moth, ringworm"},
		{"bole", "bole", "bɔlɛ", "NOUN", "", "FAMILY", 3, "newborn", "[<6mo]: newborn"},
		{"bolingo", "bolingo", "bolingo", "NOUN", "", "SENSE", 5, "love", "love"},
		{"boma", "boma", "boma", "NOUN", "", "COOK", 5, "cook-stew", "cooking with the pot sitting directly in the embers"},
		{"bondo", "bôndo", "bôndo", "VERB", "Subcat=Tran", "ALT WORD FOR", 9, "assemble", "bûngbi"},
		{"bondo", "böndö", "böndö", "NOUN", "", "FOOD", 5, "sorghum", "sorghum"},
		{"bongo", "bongô", "bɔngɔ̂", "NOUN", "", "ANIM", 3, "hyena", "hyena"},
		{"bongo", "bongö", "bɔngɔ̈", "NOUN", "", "OBJ", 2, "clothes", "clothes"},
		{"boon", "bôon", "bôon", "NOUN", "", "CIVIL", 7, "debt", "[Fr:bon]: debt"},
		{"boro", "bôrö", "bôrö", "NOUN", "", "OBJ", 6, "head-cushion", "cushion for carrying things on one's head"},
		{"boro", "börö", "börö", "NOUN", "", "SICK", 4, "goiter", "goiter"},
		{"boso", "bôso", "bɔ̂sɔ", "VERB", "Aspect=Hab", "ACT", 3, "pile-up", "pile up"},
		{"bosongbi", "bôsongbi", "bɔ̂sɔngbi", "NOUN", "", "ACT", 3, "partition", "partition"},
		{"bosongbi", "bôsongbi", "bɔ̂sɔngbi", "VERB", "Aspect=Imp|Subcat=Tran", "ACT", 3, "partition", "partition"},
		{"bosongbitere", "bôsongbi terê", "bɔ̂sɔngbi tɛrɛ̂", "VERB", "Subcat=Intr", "ACT", 3, "separate-into-groups", "[lit: partition|oneself]: separate into groups"},
		{"bozo", "bozö", "bozö", "NOUN", "", "OBJ", 2, "bag-pocket-purse", "bag, pocket, purse"},
		{"bua", "buä", "buä", "NOUN", "", "GOD ", 3, "priest", "priest, Father"},
		{"buakete", "buä-kêtê", "buä-kɛ̂tɛ̂", "NOUN", "", "GOD ", 3, "parish-priest", "[lit: priest|small]: parish priest"},
		{"buakota", "buä-kötä", "buä-kötä", "NOUN", "", "GOD ", 3, "parish-pastor", "[lit: priest|big]: parish pastor"},
		{"buamanabe", "buä-mä-na-bê", "buä-mä-na-bɛ̂", "NOUN", "", "GOD ", 3, "pastor", "[lit: priest|Protestant]: pastor"},
		{"buamokonzi", "buä-mokönzi", "buä-mokönzi", "NOUN", "", "GOD ", 3, "archbishop", "[lit: priest|chief]: archbishop, cardinal"},
		{"buasu", "buä-sû", "buä-sû", "NOUN", "", "GOD ", 3, "scribe", "[lit: priest|write]: scribe"},
		{"buate", "buäte", "buäte", "NOUN", "", "OBJ", 7, "can", "can"},
		{"buatokua", "buä-tokua", "buä-tokua", "NOUN", "", "GOD ", 3, "nuncio", "[lit: priest|sent]: nuncio"},
		{"buba", "buba", "buba", "VERB", "", "HOW", 2, "ruin", "ruin"},
		{"buba", "bübä", "bübä", "NOUN", "", "HOW", 2, "stupidity", "stupidity"},
		{"buba", "bübä", "bübä", "NOUN", "", "WHO", 2, "idiot", "idiot"},
		{"bubu", "bubu", "bubu", "NOUN", "", "OBJ", 5, "male-blouse", "formal loose blouse worn by men"},
		{"bubu", "bûbu", "bûbu", "NOUN", "", "ANIM", 4, "ape", "ape"},
		{"buburu", "bûburû", "bûburû", "NOUN", "", "WHO", 4, "deaf-mute", "deaf mute"},
		{"bubuta", "bubûtä", "bubûtä", "ADJ", "", "SENSE", 4, "reticent", "reticent, reserved"},
		{"buku", "bûku", "bûku", "NOUN", "", "OBJ", 3, "book", "book"},
		{"bulee", "bulêe", "bulɛ̂ɛ", "NOUN", "", "FOOD", 2, "banana", "banana"},
		{"bungbi", "bûngbi", "bûngbi", "NOUN", "", "CIVIL", 3, "assembly", "assembly, meeting, corps"},
		{"bungbi", "bûngbi", "bûngbi", "VERB", "Aspect=Imp|Subcat=Tran", "CIVIL", 3, "assemble", "assemble"},
		{"bungbitere", "bûngbi terê", "bûngbi tɛrɛ̂", "VERB", "Subcat=Intr", "ACT", 3, "assemble", "[lit: assemble|oneself]: assemble"},
		{"buru", "burü", "burü", "NOUN", "", "NATURE", 3, "dry-season", "dry season, aridity, drought"},
		{"buruma", "buruma", "buruma", "NOUN", "", "SICK", 4, "leprosy", "leprosy"},
		{"busu", "bûsu", "bûsu", "NOUN", "", "INTERACT", 6, "hypocrisy", "hypocrisy"},
		{"butani", "butâni", "butâni", "NOUN", "", "OBJ", 7, "bottle", "[Fr: bouteille]: bottle"},
		{"butu", "butu", "butu", "NOUN", "", "HOW", 2, "dust", "dust"},
		{"butuma", "butuma", "butuma", "ADV", "", "HOW", 5, "no-matter-how-(much)", "no matter how (much)"},
		{"butuma", "butuma", "butuma", "NOUN", "", "HOW", 5, "disorder", "mess, disorder"},
		{"butuma", "butuma", "butuma", "VERB", "Subcat=Intr", "ACT", 5, "become-corrupted", "become excited, animated; be corrupt, become corrupted, fall from grace; bloom, blossom"},
		{"buze", "büzë", "büzë", "NOUN", "", "CIVIL", 4, "trade", "barter, trade, commerce"},
		{"buzi", "buzî", "buzî", "NOUN", "", "OBJ", 4, "candle", "[Fr: bougie]: candle"},
		{"da", "da", "da", "NOUN", "", "HOUSE", 2, "house", "house, shelter"},
		{"da", "da", "da", "VERB", "Subcat=Tran", "ACT", 3, "put", "[figurative]: put, place"},
		{"da", "dä", "dä", "VERB", "Subcat=Intr", "NATURE", 5, "mildew", "mildew"},
		{"da", "dä", "dä", "VERB", "Subcat=Intr", "STATE", 2, "become", "become"},
		{"daa", "daä", "daä", "PART", "", "WHEN", 1, "then", "then"},
		{"daa", "daä", "daä", "PART", "", "WHERE", 1, "there", "there"},
		{"dabe", "da bê", "da bɛ̂", "VERB", "Subcat=Intr", "FEEL", 1, "recollect", "[lit: place|heart]: think (of), recollect"},
		{"dabe", "da-bê", "da-bɛ̂", "NOUN", "", "FEEL", 1, "recollection", "[lit: place|heart]: thought (of), memory, recollection"},
		{"dakosara", "da-kosâra", "da-kosâra", "NOUN", "", "WHERE", 2, "office", "[lit: house|job]: office"},
		{"dalama", "dâlâmâ", "dâlâmâ", "NOUN", "", "SICK", 4, "epidemic", "epidemic"},
		{"dale", "dâlë", "dâlë", "NOUN", "", "ANIM", 3, "toad", "toad"},
		{"damakongo", "damâköngö", "damâköngö", "NOUN", "", "ANIM", 6, "scorpion", "scorpion"},
		{"damango", "damango", "damangɔ", "NOUN", "", "ANIM", 6, "turtle", "turtle"},
		{"damazani", "damazäni", "damazäni", "NOUN", "", "OBJ", 7, "jug", "[Fr:dame-jeanne]: demijohn, carboy, 20-liter wicker-covered glass jug"},
		{"damba", "dambâ", "dambâ", "NOUN", "", "BODY", 4, "tail", "tail"},
		{"dami", "dâmi", "dâmi", "NOUN", "", "INTERACT", 4, "proverb", "proverb"},
		{"damvene", "damvene", "damvɛnɛ", "NOUN", "", "ANIM", 4, "spider", "spider"},
		{"danabe", "da na bê", "da na bɛ̂", "VERB", "Subcat=Tran", "FEEL", 1, "recall", "[lit: place|in|heart]: think, recall"},
		{"danda", "dândâ", "dândâ", "NOUN", "", "SICK", 4, "headache", "severe headache"},
		{"danga", "dânga", "dânga", "NOUN", "", "HOUSE", 4, "hut", "hut"},
		{"dangalinga", "dangâlingâ", "dangâlingâ", "NOUN", "", "ANIM", 6, "praying-mantis", "praying mantis"},
		{"dangara", "dangara", "dangara", "NOUN", "", "HOUSE", 4, "hangar", "hangar"},
		{"dangbo", "dangbö", "dangbö", "NOUN", "", "ANIM", 6, "chamelion", "chamelion"},
		{"dangere", "dangërë", "dangërë", "NOUN", "", "DRINK", 5, "bamboo-palm-wine", "bamboo palm wine"},
		{"dangi", "dangi", "dangi", "NOUN", "", "NATURE", 6, "termite-mound", "termite mound"},
		{"dara", "dara", "dara", "VERB", "Aspect=Hab|Subcat=Tran", "ACT", 4, "shake", "[lit: place|repeatedly]: shake, caress, smooth out"},
		{"daraa", "daräa", "daräa", "NOUN", "", "HOUSE", 7, "bedsheet", "[Fr:drap]: bedsheet"},
		{"daturu", "da-turu", "da-turu", "NOUN", "", "WHO", 4, "blacksmith-shop", "[lit: house|forge]: blacksmith shop"},
		{"daveke", "daveke", "davɛkɛ", "NOUN", "", "SICK", 5, "syphilis", "syphilis"},
		{"dawaa", "dawäa", "dawäa", "NOUN", "", "WHERE", 7, "in-front-of", "[lit: underneath|front]: in front of"},
		{"dazo", "dazo", "dazo", "NOUN", "", "FOOD", 5, "potato", "potato"},
		{"de", "de", "de", "VERB", "Subcat=Intr", "BODY", 3, "vomit", "vomit"},
		{"de", "de", "dɛ", "VERB", "", "STATE", 2, "remain", "remain"},
		{"de", "dê", "dê", "NOUN", "", "HOW", 3, "coldness", "coldness, shade"},
		{"de", "dë", "dë", "VERB", "Subcat=Intr", "HOW", 3, "be-cold", "be cold"},
		{"de", "dë", "dɛ̈", "VERB", "Subcat=Tran", "ACT", 2, "cut-or-grow", "cut, slice; grow, cultivate"},
		{"de", "dë", "dɛ̈", "VERB", "Subcat=Tran", "INTERACT", 3, "emit", "emit"},
		{"deba", "dë-bä", "dɛ̈-bä", "VERB", "", "INTERACT", 3, "swear-an-oath", "[lit: emit|oath]: swear, swear in, administer an oath"},
		{"deba", "dëbä", "dɛ̈bä", "NOUN", "", "INTERACT", 3, "blessing-or-curse", "[lit: emit|oath]: blessing, curse"},
		{"debango", "dëbängö", "dɛ̈bängɔ̈", "VERB", "VerbForm=Vnoun", "INTERACT", 3, "oath", "oath"},
		{"debanzoni", "dë-bä-nzönî", "dɛ̈-bä-nzɔ̈nî", "VERB", "Subcat=Intr", "INTERACT", 3, "bless", "[lit: emit|oath|good]: administer a blessing"},
		{"debanzoninandoti", "dë-bä-nzönî na ndö tî", "dɛ̈-bä-nzɔ̈nî na ndö tî", "VERB", "Subcat=Tran", "INTERACT", 3, "bless", "[lit: emit|oath|good|on|place]: bless"},
		{"debasioni", "dë-bä-sïönî", "dɛ̈-bä-sïɔ̈nî", "VERB", "Subcat=Intr", "INTERACT", 3, "curse", "[lit: emit|oath|good]: administer a curse"},
		{"debasioninandoti", "dë-bä-sïönî na ndö tî", "dɛ̈-bä-sïɔ̈nî na ndö tî", "VERB", "Subcat=Tran", "INTERACT", 3, "curse", "[lit: emit|oath|good|on|place]: curse"},
		{"debuze", "dë-büzë", "dɛ̈-büzë", "VERB", "Subcat=Intr", "INTERACT", 2, "do-commerce", "[lit: cultivate|commerce]: engage in commerce"},
		{"defa", "dêfa", "dêfa", "VERB", "Subcat=Tran", "CIVIL", 3, "borrow-or-lend", "borrow, lend"},
		{"dekite", "dë-kîte", "dɛ̈-kîtɛ", "VERB", "Subcat=Intr", "SENSE", 3, "doubt", "[lit: emit|doubt]: doubt"},
		{"dekongo", "dë-köngö", "dɛ̈-kɔ̈ngɔ̈", "VERB", "Subcat=Intr|VerbForm=Vnoun", "INTERACT", 2, "cry-out", "[lit: emit|cry]: cry out"},
		{"deku", "deku", "dɛku", "NOUN", "", "ANIM", 3, "mouse-or-rat", "mouse, rat"},
		{"dema", "dema", "dema", "VERB", "Subcat=Tran", "SENSE", 3, "pity", "pity"},
		{"demangotere", "dëmängö-terê", "dëmängɔ̈-tɛrɛ̂", "VERB", "Subcat=Intr", "SENSE", 3, "lamentations", "lamentations"},
		{"dematere", "dema-terê", "dema-tɛrɛ̂", "VERB", "Subcat=Intr", "SENSE", 3, "lament", "lament"},
		{"dengbe", "dengbe", "dɛngbɛ", "NOUN", "", "ANIM", 4, "small-antilope", "Maxwell's duiker (small antilope)"},
		{"denge", "dênge", "dênge", "VERB", "Subcat=Tran", "ACT", 4, "bend-down", "bend down"},
		{"dengi", "dêngi", "dêngi", "NOUN", "", "INTERACT", 4, "curse", "curse"},
		{"dengo", "dëngö", "dëngɔ̈", "VERB", "VerbForm=Vnoun", "HOW", 3, "cold", "cold"},
		{"dengo", "dëngö", "dɛ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "INTERACT", 4, "pronunciation", "pronunciation"},
		{"denzoba", "dë-nzö-bä", "dɛ̈-nzɔ̈-bä", "VERB", "Subcat=Intr", "INTERACT", 3, "bless", "[lit: emit|good|oath]: administer a blessing"},
		{"denzobanandoti", "dë-nzö-bä na ndö tî", "dɛ̈-nzɔ̈-bä na ndö tî", "VERB", "Subcat=Tran", "INTERACT", 3, "bless", "[lit: emit|good|oath|on|place]: bless"},
		{"dere", "derë", "dɛrɛ̈", "NOUN", "", "HOUSE", 3, "fence-or-dam", "wall, fence, dam"},
		{"desioba", "dë-sïö-bä", "dɛ̈-sïɔ̈-bä", "VERB", "Subcat=Intr", "INTERACT", 3, "curse", "[lit: emit|bad|oath]: administer a curse"},
		{"desiobanandoti", "dë-sïö-bä na ndö tî", "dɛ̈-sïɔ̈-bä na ndö tî", "VERB", "Subcat=Tran", "INTERACT", 3, "curse", "[lit: emit|bad|oath|on|place]: curse"},
		{"deyaka", "dë-yäkä", "dɛ̈-yäkä", "VERB", "Subcat=Intr", "ACT", 2, "grow-crops", "grow crops, cultivate a field"},
		{"di", "di", "di", "VERB", "Subcat=Intr", "STATE", 3, "stick", "stick, adhere"},
		{"di", "dî", "dî", "NOUN", "", "INTERACT", 3, "pronunciation", "pronunciation"},
		{"di", "dï", "dï", "VERB", "", "INTERACT", 3, "pronounce", "pronounce"},
		{"didi", "didi", "didi", "NOUN", "", "BODY", 4, "animal-horn", "animal horn"},
		{"didiri", "dïdïrï", "dïdïrï", "NOUN", "", "ANIM", 5, "mud-wasp", "mud wasp"},
		{"diiriti", "dï ïrï tî", "dï ïrï tî", "VERB", "Subcat=Tran", "INTERACT", 3, "denounce", "denounce"},
		{"diki", "dïkï", "dïkï", "NOUN", "", "HOUSE", 4, "hearth", "hearth"},
		{"dikinzi", "dïkïnzï", "dïkïnzï", "NOUN", "", "ANIM", 5, "shrimp", "shrimp"},
		{"diko", "dîko", "dîkɔ", "VERB", "Subcat=Tran", "SENSE", 2, "read-or-count", "read, count"},
		{"diko", "dîkô", "dîkɔ̂", "NOUN", "", "SENSE", 2, "reading-or-counting", "reading, counting"},
		{"do", "do", "do", "NOUN", "", "WHERE", 2, "west-or-downstream", "west; downstream"},
		{"do", "dô", "dɔ̂", "VERB", "Subcat=Intr", "BODY", 3, "tremble", "tremble"},
		{"do", "dö", "dɔ̈", "NOUN", "", "OBJ", 4, "hatchet", "hatchet"},
		{"dodo", "dodo", "dɔdɔ", "VERB", "Mood=Emp", "BODY", 3, "dance", "dance"},
		{"dodo", "dödö", "dɔ̈dɔ̈", "NOUN", "Mood=Emp", "BODY", 3, "dance", "dance"},
		{"dodoro", "dödörö", "dödörö", "NOUN", "", "ANIM", 6, "partridge", "partridge"},
		{"dokpa", "dokpa", "dokpa", "ADJ", "", "HOW", 4, "unripe", "unripe"},
		{"dokpala", "dokpâlâ", "dɔkpâlâ", "NOUN", "", "ANIM", 6, "kite", "kite"},
		{"doli", "doli", "doli", "NOUN", "", "ANIM", 3, "elephant", "elephant"},
		{"dolo", "dolö", "dolö", "NOUN", "", "DRINK", 5, "corn-beer", "corn beer"},
		{"donali", "dö na li", "dɔ̈ na li", "VERB", "Subcat=Tran", "OBJ", 4, "seduce", "charm, seduce"},
		{"dondo", "dondö", "dondö", "NOUN", "", "BODY", 5, "vagina", "vagina"},
		{"dondo", "dôndô", "dôndô", "NOUN", "", "FOOD", 5, "corn-paste-bar", "corn paste bar"},
		{"dongba", "döngbä", "döngbä", "NOUN", "", "FISH", 6, "catfish", "catfish"},
		{"dongo", "dongo", "dɔngɔ", "VERB", "Subcat=Tran", "ACT", 3, "classify", "classify"},
		{"dongododo", "dongö dödö", "dɔngɔ̈ dɔ̈dɔ̈", "VERB", "Subcat=Intr|VerbForm=Vnoun", "BODY", 3, "dancing", "dancing"},
		{"dongongbi", "dongöngbi", "dɔngɔ̈ngbi", "NOUN", "", "ACT", 3, "arrangement", "arrangement"},
		{"dongongbi", "dongöngbi", "dɔngɔ̈ngbi", "VERB", "Aspect=Imp|Subcat=Tran", "ACT", 3, "arrange", "arrange"},
		{"dongongbitere", "dongöngbi terê", "dɔngɔ̈ngbi tɛrɛ̂", "VERB", "Subcat=Intr", "ACT", 3, "get-in-line", "[lit: arrange|oneself]: arrange oneself, get in order, get in line"},
		{"doro", "dörö", "dɔ̈rɔ̈", "VERB", "Mood=Emp", "BODY", 4, "shiver", "shiver"},
		{"doroko", "doroko", "dɔrɔkɔ", "VERB", "Subcat=Tran", "ACT", 4, "eviscerate", "eviscerate"},
		{"du", "du", "du", "VERB", "Subcat=Intr", "STATE", 2, "sit", "sit"},
		{"du", "dû", "dû", "NOUN", "", "NATURE", 3, "hole", "hole"},
		{"du", "dü", "dü", "VERB", "", "BODY", 3, "give-birth-or-be-born", "give birth, be born"},
		{"dudu", "dudu", "dudu", "NOUN", "", "CIVIL", 6, "poverty", "poverty"},
		{"duma", "duma", "duma", "NOUN", "", "DRINK", 5, "honey-beer", "honey beer"},
		{"dungo", "düngö", "düngɔ̈", "VERB", "VerbForm=Vnoun", "BODY", 4, "birth", "birth, existence"},
		{"dungu", "dû-ngû", "dû-ngû", "NOUN", "", "NATURE", 3, "well", "[lit: hole|water]: well, cistern"},
		{"dunia", "dûnîa", "dûnîa", "NOUN", "", "CIVIL", 3, "world", "world, universe, history"},
		{"dunyene", "dû-nyenë", "dû-nyɛnɛ̈", "NOUN", "", "BODY", 6, "anus", "[lit: hole|buttocks]: anus"},
		{"duru", "dûru", "dûru", "VERB", "Subcat=Tran", "COOK", 4, "boil", "boil"},
		{"duti", "dutï", "dutï", "VERB", "Subcat=Intr", "STATE", 1, "sit-or-live", "sit, live"},
		{"dutinzoni", "dutï nzönî", "dutï nzɔ̈nî", "INTERJ", "Mood=Opt", "STATE", 1, "Goodbye", "[lit: sit|well]: Goodbye! (said to those staying)"},
		{"e", "e", "ɛ", "VERB", "Subcat=Intr", "HOW", 6, "be-sharp", "be sharp"},
		{"e", "ë", "ë", "PRON", "Num=Plur|Person=1|PronType=Prs", "WHO", 1, "we", "we, us"},
		{"ekalitise", "ëkälïtïse", "ëkälïtïse", "NOUN", "", "TREE", 6, "eucalyptus", "eucalyptus tree"},
		{"emveni", "ë mvenî", "ë mvɛnî", "PRON", "Num=Plur|Person=1|PronType=Prs", "WHO", 1, "ourselves", "[lit: we|self]: ourselves"},
		{"epatite", "ëpätîte", "ëpätîte", "NOUN", "", "SICK", 5, "hepatitis", "hepatitis"},
		{"ere", "ere", "ere", "VERB", "Subcat=Tran", "ACT", 4, "pluck", "pluck"},
		{"ere", "ere", "ɛrɛ", "VERB", "Subcat=Intr", "ALT WORD FOR", 6, "disappear", "yîkɔ"},
		{"erege", "êrêge", "ɛ̂rɛ̂gɛ", "NOUN", "", "OBJ", 6, "liquor", "liquor, distilled alcohol"},
		{"fa", "fa", "fa", "VERB", "Subcat=Tran", "INTERACT", 1, "show", "show"},
		{"fa", "fâ", "fâ", "NOUN", "", "ACT", 3, "bet-or-fraction", "bet, fraction"},
		{"fa", "fä", "fä", "NOUN", "", "PLANT", 5, "wildflower", "wildflower"},
		{"fa", "fä", "fä", "VERB", "Subcat=Tran", "ACT", 3, "wager", "bet, wager"},
		{"faa", "fâa", "fâa", "VERB", "Subcat=Tran", "ACT", 2, "cross-cut-strike-break-kill", "cross, traverse; cut, strike, break; wound, kill"},
		{"fade", "fadë", "fadë", "ADV", "", "WHEN", 1, "right-now-or-will", "[postverbal]: right now; [preverbal]: will, shall"},
		{"fadeso", "fadësô", "fadësô", "ADV", "", "WHEN", 1, "now", "now"},
		{"fafadeso", "fafadësô", "fafadësô", "ADV", "", "WHEN", 1, "immediately", "immediately"},
		{"fala", "fâla", "fâla", "VERB", "Aspect=Hab|Subcat=Tran", "ALT SP FOR", 9, "chop", "fâra"},
		{"falambio", "fâla mbïö", "fâla mbïɔ̈", "VERB", "Subcat=Intr", "ACT", 4, "apply-Western-makeup", "put on Western makeup"},
		{"falazua", "fâla zûâ", "fâla zûâ", "VERB", "Subcat=Intr", "ACT", 4, "apply-traditional-makeup", "put on traditional makeup"},
		{"fangbi", "fângbi", "fângbi", "NOUN", "", "ACT", 3, "section", "section"},
		{"fangbi", "fângbi", "fângbi", "VERB", "Aspect=Imp|Subcat=Tran", "ACT", 3, "cut-up", "cut up"},
		{"fani", "fâ-nî", "fâ-nî", "NOUN", "", "WHEN", 1, "times", "times"},
		{"fara", "fâra", "fâra", "VERB", "Aspect=Hab|Subcat=Tran", "ACT", 4, "chop", "[lit: cut|repeatedly]: chop"},
		{"faranzi", "farânzi", "farânzi", "NOUN", "", "COUNTRY", 1, "France", "France"},
		{"farini", "farïni", "farïni", "NOUN", "", "FOOD", 4, "wheat", "wheat"},
		{"fen", "fên", "fên", "ADV", "", "SENSE", 3, "feel-bad", "feel bad"},
		{"ferere", "fêrêrê", "fɛ̂rɛ̂rɛ̂", "NOUN", "", "ACT", 6, "whistle", "whistle"},
		{"fi", "fi", "fi", "VERB", "Subcat=Tran", "ACT", 3, "stir", "stir; whistle"},
		{"fimbo", "fîmbo", "fîmbo", "NOUN", "", "OBJ", 4, "whip", "whip"},
		{"fingi", "fingi", "fingi", "NOUN", "", "ANIM", 6, "monkey", "colobus monkey"},
		{"fini", "finî", "finî", "ADJ", "", "WHEN", 3, "new", "new, fresh"},
		{"fini", "finî", "finî", "NOUN", "", "CIVIL", 3, "life", "life, living"},
		{"finon", "finön", "finön", "NOUN", "", "CIVIL", 4, "pain", "pain, poverty"},
		{"fo", "fo", "fo", "NOUN", "", "FAMILY", 2, "colleague", "colleague, coworker"},
		{"fo", "fö", "fö", "NOUN", "", "WHEN", 4, "interval", "interval"},
		{"fondo", "fondo", "fɔndɔ", "NOUN", "", "FOOD", 2, "plantain", "plantain"},
		{"fondo", "föndo", "föndo", "NOUN", "", "WHEN", 5, "June", "June"},
		{"fono", "fono", "fɔnɔ", "VERB", "", "MOVE", 2, "wander", "wander"},
		{"fono", "fönö", "fɔ̈nɔ̈", "NOUN", "", "MOVE", 2, "stroll", "stroll"},
		{"fu", "fû", "fû", "VERB", "Subcat=Tran", "BODY", 3, "grab-handful", "grab handful"},
		{"fu", "fü", "fü", "VERB", "Subcat=Tran", "ACT", 2, "tailor", "tailor"},
		{"fufu", "fufû", "fufû", "NOUN", "", "BODY", 4, "lungs", "lungs"},
		{"fufu", "fufû", "fufû", "NOUN", "", "FOOD", 3, "manioc-ball", "manioc ball"},
		{"fufulafu", "fufulafu", "fufulafu", "NOUN", "", "SICK", 4, "rabies-or-convulsions", "rabies, convulsions"},
		{"fuku", "fûku", "fûku", "NOUN", "", "FOOD", 2, "flour", "flour"},
		{"fulundingi", "fulundïngi", "fulundïngi", "NOUN", "", "WHEN", 5, "February", "February"},
		{"fun", "fûn", "fûn", "NOUN", "", "SENSE", 3, "smell", "smell, odor"},
		{"fun", "fün", "fün", "VERB", "Subcat=Intr", "SENSE", 3, "smell", "smell"},
		{"funga", "fungâ", "fungâ", "NOUN", "", "ACT", 4, "embroidery", "embroidery"},
		{"funga", "fûnga", "fûnga", "VERB", "Aspect=Iter|Subcat=Tran", "ACT", 4, "embroider", "embroider"},
		{"fungula", "fungûla", "fungûla", "NOUN", "", "OBJ", 4, "lock-or-key", "lock, key"},
		{"fungula", "fungûla", "fungûla", "VERB", "Subcat=Tran", "OBJ", 4, "unlock", "unlock"},
		{"funngo", "fünngö", "fünngɔ̈", "VERB", "VerbForm=Vnoun", "SENSE", 4, "odor", "odor"},
		{"furu", "fûru", "fûru", "NOUN", "", "BODY", 4, "foam", "foam, scum, drool"},
		{"furu", "fûru", "fûru", "VERB", "Aspect=Hab|Subcat=Tran", "ACT", 4, "knead-or-brew", "knead, brew"},
		{"futa", "fûta", "fûta", "NOUN", "", "CIVIL", 2, "payment", "payment, salary"},
		{"futa", "fûta", "fûta", "VERB", "Subcat=Tran", "CIVIL", 2, "pay", "pay, repay"},
		{"fuu", "füu", "füu", "NOUN", "", "SENSE", 4, "crazy", "crazy"},
		{"ga", "gä", "gä", "VERB", "Subcat=Intr", "MOVE", 1, "come", "come"},
		{"ga", "gä", "gä", "VERB", "Subcat=Tran", "STATE", 1, "become", "become"},
		{"gagi", "gagi", "gagi", "NOUN", "", "FISH", 6, "rayfinned-fish", "ray finned fish"},
		{"gana", "gana", "gana", "VERB", "Subcat=Intr", "ACT", 3, "wrap", "wrap"},
		{"gana", "gana", "gana", "VERB", "Subcat=Tran", "ACT", 4, "lapidate", "throw stones at"},
		{"gana", "gä na", "gä na", "VERB", "Subcat=Tran", "MOVE", 1, "come-with", "come with, bring"},
		{"ganda", "gânda", "gânda", "NOUN", "", "ANIM", 5, "mud-wasp", "mud wasp"},
		{"ganga", "gängä", "gängä", "NOUN", "", "FISH", 6, "caiman-fish", "caiman fish"},
		{"gangara", "gângârâ", "gângârâ", "NOUN", "", "HOUSE", 4, "roof-beams", "roof beams"},
		{"gangba", "gangba", "gangba", "ADJ", "", "HOW", 4, "ripe", "ripe"},
		{"gangbi", "gângbi", "gângbi", "NOUN", "", "MOVE", 3, "convergence", "coming at the same time, convergence"},
		{"gangbi", "gângbi", "gângbi", "VERB", "Aspect=Imp|Subcat=Intr", "MOVE", 3, "converge", "come at the same time, converge"},
		{"ganza", "ganzâ", "ganzâ", "NOUN", "", "BODY", 4, "circumcision", "circumcision, excision"},
		{"gao", "gao", "gao", "NOUN", "", "HOW", 3, "beauty", "beauty"},
		{"gapa", "gapa", "gapa", "NOUN", "", "INTERACT", 3, "menace", "menace"},
		{"gara", "gara", "gara", "VERB", "Aspect=Hab", "MOVE", 3, "come-often", "[lit: come|repeatedly]: come often; [lit: come|repeatedly]: flow; strip leaves off"},
		{"gara", "garâ", "garâ", "NOUN", "", "CIVIL", 2, "market", "market"},
		{"gasa", "gasa", "gasa", "NOUN", "", "FISH", 6, "catfish", "upside down catfish"},
		{"gasa", "gasa", "gasa", "VERB", "", "STATE", 3, "endure", "put up with"},
		{"gati", "gatï", "gatï", "ADV", "", "WHERE", 4, "left", "left hand, left side"},
		{"gba", "gba", "gba", "ADJ", "", "NUM", 2, "much-or-many", "much, many"},
		{"gba", "gba", "gba", "VERB", "Subcat=Tran", "CIVIL", 4, "have-sex-with", "have sex with"},
		{"gba", "gbâ", "gbâ", "NOUN", "", "NUM", 2, "bunch", "bunch, group, bundle"},
		{"gba", "gbâ", "gbâ", "VERB", "Subcat=Tran", "STATE", 3, "squeeze-or-crowd", "squeeze, crowd"},
		{"gba", "gbä", "gbä", "ADV", "", "HOW", 1, "in-vain", "in vain"},
		{"gbadola", "gbadöla", "gbadöla", "NOUN", "", "ANIM", 5, "locust", "locust"},
		{"gbafu", "gbafu", "gbafu", "NOUN", "", "NATURE", 4, "flood", "tidal flow; flood"},
		{"gbaga", "gbagä", "gbagä", "NOUN", "", "ANIM", 5, "mongoose", "mongoose"},
		{"gbagba", "gbägbä", "gbägbä", "ADJ", "Mood=Emp", "HOW", 4, "spoiled", "spoiled"},
		{"gbagba", "gbägbä", "gbägbä", "NOUN", "", "ANIM", 5, "red-ant", "red ant"},
		{"gbagba", "gbägbä", "gbägbä", "NOUN", "", "HOUSE", 3, "fence", "fence,  enclosure,  yard,  property"},
		{"gbagbara", "gbagbara", "gbagbara", "VERB", "Subcat=Tran", "ACT", 4, "shake", "[lit: squeeze|squeeze|repeatedly]: shake"},
		{"gbagbara", "gbâgbârâ", "gbâgbârâ", "NOUN", "", "OBJ", 4, "iron-brush", "iron brush"},
		{"gbaka", "gbâka", "gbâka", "ADV", "", "HOW", 4, "giant", "giant"},
		{"gbakaragba", "gbâkarâgba", "gbâkarâgba", "NOUN", "", "WHEN", 6, "February", "February"},
		{"gbako", "gbakô", "gbakô", "NOUN", "", "NATURE", 3, "forest", "forest"},
		{"gbakuru", "gbâkûrû", "gbâkûrû", "NOUN", "", "OBJ", 6, "tool", "tool"},
		{"gbalaka", "gbalâka", "gbalâka", "NOUN", "", "GOD", 4, "altar", "altar"},
		{"gbambingo", "gba-mbîngo", "gba-mbîngo", "NOUN", "", "INTERACT", 4, "secret", "[lit:gbx=underneath|darkness]: secret"},
		{"gbanambana", "gba-na-mbänä", "gba-na-mbänä", "NOUN", "", "WHO", 4, "prostitute", "[lit:have sex with|in|wickedness]: prostitute"},
		{"gbanda", "gbânda", "gbânda", "NOUN", "", "OBJ", 3, "net", "net"},
		{"gbanda", "gbândä", "gbândä", "ADV", "", "WHEN", 1, "later", "later"},
		{"gbandasango", "gbânda-sango", "gbânda-sango", "NOUN", "", "COMPUTER", 6, "internet", "[lit: net|news]: internet"},
		{"gbandatitere", "gbânda tî tere", "gbânda tî tɛrɛ", "NOUN", "", "ANIM", 4, "spider-web", "[lit: net|of|spider]: spider web"},
		{"gbanga", "gbängä", "gbängä", "NOUN", "", "TREE", 5, "nutmeg-tree", "nutmeg tree"},
		{"gbanza", "gbanza", "gbanza", "NOUN", "", "FOOD", 6, "corn", "corn"},
		{"gbanzi", "gbânzi", "gbânzi", "VERB", "Subcat=Tran", "ACT", 4, "prevent", "prevent"},
		{"gbanzia", "gbanzia", "gbanzia", "NOUN", "", "FOOD", 6, "corn", "corn"},
		{"gbara", "gbara", "gbara", "NOUN", "", "ACT", 4, "grill", "grill, frying pan"},
		{"gbara", "gbara", "gbara", "VERB", "Subcat=Tran", "ACT", 4, "spread-out", "[lit: squeeze|repeatedly]: spread out"},
		{"gbaragaza", "gbaragaza", "gbaragaza", "NOUN", "", "OBJ", 4, "broom", "broom; throwing knife"},
		{"gbaraka", "gbârâka", "gbârâka", "NOUN", "", "OBJ", 4, "drying-rack", "drying rack"},
		{"gbari", "gbari", "gbari", "NOUN", "", "FOOD", 6, "bean", "bean"},
		{"gbata", "gbätä", "gbätä", "NOUN", "", "CIVIL", 4, "post-or-assignment", "post, station"},
		{"gbaza", "gbâzâ", "gbâzâ", "NOUN", "", "OBJ", 3, "wheel", "wheel, circle"},
		{"gbazabanga", "gbâzâ-bängâ", "gbâzâ-bängâ", "NOUN", "", "OBJ", 3, "bicycle", "[lit: wheel|rubber]: bicycle"},
		{"gbazagbo", "gbäzägbö", "gbäzägbö", "NOUN", "", "ANIM", 6, "antilope", "antilope"},
		{"gbe", "gbe", "gbɛ", "NOUN", "", "WHERE", 2, "underneath", "underneath"},
		{"gbe", "gbë", "gbë", "VERB", "Subcat=Tran", "ACT", 3, "embarrass", "attach; weigh down, embarrass"},
		{"gbee", "gbëe", "gbɛ̈ɛ", "VERB", "Subcat=Intr", "HOW", 3, "grow-old", "grow old"},
		{"gbefa", "gbefâ", "gbefâ", "NOUN", "", "HOUSE", 4, "veranda", "veranda"},
		{"gbegbere", "gbegbëre", "gbɛgbɛ̈rɛ", "NOUN", "", "ALT WORD FOR", 6, "mange", "särä"},
		{"gbele", "gbe-lê", "gbɛ-lɛ̂", "NOUN", "", "WHERE", 2, "in-front-of", "[lit: underneath|front]: in front of"},
		{"gbelewele", "gbêlêwele", "gbêlêwele", "NOUN", "", "MYTH", 6, "wildcat", "wildcat"},
		{"gbene", "gbene", "gbɛnɛ", "VERB", "Subcat=Tran", "ACT", 6, "measure", "cut out; measure, define"},
		{"gbene", "gbënë", "gbɛ̈nɛ̈", "NOUN", "", "ACT", 6, "definiteness", "definiteness"},
		{"gbenga", "gbênga", "gbênga", "VERB", "Aspect=Iter|Subcat=Tran", "ACT", 4, "tie-up", "tie up"},
		{"gbengbi", "gbêngbi", "gbêngbi", "VERB", "Aspect=Imp|Subcat=Tran", "CIVIL", 3, "fight-over", "fight over"},
		{"gbengbitere", "gbêngbi terê", "gbêngbi tɛrɛ̂", "VERB", "Subcat=Intr", "ACT", 3, "change", "change, transform"},
		{"gbengbitere", "gbêngbi terê", "gbêngbi tɛrɛ̂", "VERB", "Subcat=Intr", "CIVIL", 3, "fight-each-other", "fight amongst oneselves"},
		{"gbenyongbia", "gbe-nyön-gbïä", "gbɛ-nyön-gbïä", "NOUN", "", "WHO", 2, "government-minister", "[lit: underneath|mouth|king]: (government) minister"},
		{"gbenzi", "gbenzï", "gbenzï", "NOUN", "", "PLANT", 5, "mistletoe", "mistletoe"},
		{"gbere", "gberê", "gbɛrɛ̂", "NOUN", "", "WHERE", 2, "in-front-of", "[alt: gbe-le]: in front of"},
		{"gbi", "gbï", "gbï", "VERB", "Subcat=Intr", "HOW", 3, "be-burnt-or-feverish", "burn; be feverish, have malaria"},
		{"gbi", "gbï", "gbï", "VERB", "Subcat=Tran", "HOW", 3, "burn", "burn, heat"},
		{"gbia", "gbïä", "gbïä", "NOUN", "", "GOD", 2, "king-or-Lord", "king, Lord"},
		{"gbiangbi", "gbîangbi", "gbîangbi", "VERB", "Aspect=Imp|Subcat=Tran", "ACT", 3, "change", "change, transform"},
		{"gbigbi", "gbigbi", "gbigbi", "NOUN", "", "FISH", 6, "electric-eel", "electric eel"},
		{"gbiki", "gbikï", "gbikï", "NOUN", "", "BODY", 4, "sweat", "sweat"},
		{"gbingo", "gbïngö", "gbïngɔ̈", "VERB", "VerbForm=Vnoun", "HOW", 3, "heat", "heat"},
		{"gbogbo", "gbogbo", "gbogbo", "NOUN", "", "HOUSE", 3, "bed", "bed"},
		{"gbogbolinda", "gbôgbôlinda", "gbôgbôlinda", "NOUN", "", "SICK", 6, "rabies", "rabies"},
		{"gbokoro", "gbôkôrô", "gbôkôrô", "NOUN", "", "FOOD", 6, "peas", "peas"},
		{"gbongu", "gbo-ngû", "gbɔ-ngû", "VERB", "Subcat=Intr", "ACT", 3, "bathe-or-swim", "[lit: take|water]: bathe, swim"},
		{"gboro", "gboro", "gbɔrɔ", "NOUN", "", "FOOD", 6, "okra", "okra"},
		{"gbote", "gbo-te", "gbɔ-tɛ", "VERB", "Subcat=Intr", "ACT", 4, "swim", "swim"},
		{"gboto", "gbôto", "gbɔ̂tɔ", "VERB", "Subcat=Tran", "ACT", 3, "pull", "pull"},
		{"gbu", "gbû", "gbû", "VERB", "Subcat=Intr", "ACT", 3, "grab", "grab; grab, seize"},
		{"gbugbu", "gbûgbû", "gbûgbû", "ADV", "", "HOW", 6, "somersault", "somersault"},
		{"gbugburu", "gbugburu", "gbugburu", "NOUN", "", "FISH", 6, "electric-eel", "electric eel"},
		{"gbugburu", "gbugburu", "gbugburu", "VERB", "Subcat=Tran", "CIVIL", 4, "fight-over", "fight over"},
		{"ge", "ge", "ge", "ADV", "", "WHERE", 1, "here", "here"},
		{"gekoro", "gêkôrô", "gêkôrô", "NOUN", "", "ANIM", 6, "python", "python"},
		{"gene", "gene", "gɛnɛ", "NOUN", "", "CIVIL", 3, "visitor", "visitor; visit"},
		{"genia", "genia", "genia", "ADJ", "", "HOW", 4, "meticulous", "meticulous"},
		{"genyengo", "gënyëngö", "gɛ̈nyɛ̈ngɔ̈", "NOUN", "", "OBJ", 4, "throwing-knife", "[lit: (wondering) what's here?] throwing knife"},
		{"gere", "gerê", "gɛrɛ̂", "NOUN", "", "BODY", 2, "leg-or-foot", "leg, foot; support, foundation, underpinnings"},
		{"gerere", "gerere", "gɛrɛrɛ", "ADJ", "", "HOW", 3, "useless", "ordinary, useless, empty"},
		{"gerere", "gerere", "gɛrɛrɛ", "ADV", "", "HOW", 3, "in-vain", "in vain"},
		{"gerewungo", "gerê-wüngö", "gɛrɛ̂-wüngɔ̈", "NOUN", "", "NUM", 2, "digit", "[lit: leg|number]: digit"},
		{"gete", "gëtë", "gɛ̈tɛ̈", "NOUN", "", "FISH", 6, "rayfinned-fish", "ray finned fish"},
		{"gi", "gi", "gi", "VERB", "Subcat=Tran", "ACT", 1, "look-for", "look for, seek, search for, hunt for; annoy"},
		{"gi", "gï", "gï", "ADV", "", "HOW", 1, "only", "only"},
		{"gia", "gîâ", "gîâ", "VERB", "Subcat=Tran", "CIVIL", 4, "repay", "repay, avenge"},
		{"gibe", "gi bê", "gi bɛ̂", "VERB", "Subcat=Intr", "FEEL", 1, "meditate", "[lit: seek|heart]: think (about), meditate"},
		{"gibe", "gi-bê", "gi-bɛ̂", "NOUN", "", "FEEL", 1, "meditation", "[lit: seek|heart]: thought (about), meditation"},
		{"gidi", "gidi", "gidi", "NOUN", "", "GAME", 4, "dice-game", "dice game"},
		{"gigi", "gïgî", "gïgî", "NOUN", "", "WHERE", 2, "outside", "outside"},
		{"gilisa", "gilisa", "gilisa", "VERB", "", "SENSE", 2, "lose", "lose, forget"},
		{"ginabe", "gi na bê", "gi na bɛ̂", "VERB", "Subcat=Tran", "FEEL", 1, "think-about", "[lit: seek|in|heart]: consider carefully, think about, meditate on"},
		{"gindi", "gindî", "gindî", "NOUN", "", "OBJ", 4, "hunting-bow", "bow(for arrows)"},
		{"ginon", "ginon", "ginon", "NOUN", "", "CIVIL", 5, "bravery", "bravery, courage, ardor in battle"},
		{"gio", "gio", "giɔ", "VERB", "", "ACT", 3, "pull", "pull"},
		{"giriri", "giriri", "giriri", "ADV", "", "WHEN", 1, "long-ago", "long ago"},
		{"go", "gô", "gɔ̂", "NOUN", "", "BODY", 2, "throat-or-voice", "throat, voice"},
		{"gobi", "gobi", "gobi", "NOUN", "", "WHO", 5, "half-breed", "half breed"},
		{"gobo", "gobo", "gɔbɔ", "NOUN", "", "BODY", 3, "fist", "fist"},
		{"godobe", "godobe", "gɔdɔbx", "NOUN", "", "WHO", 4, "street-youth", "juvenile delinquent, street youth"},
		{"gogo", "gögö", "gögö", "NOUN", "", "FISH", 6, "catfish", "upside down catfish"},
		{"gogoro", "gogoro", "gɔgɔrɔ", "NOUN", "", "HOUSE", 6, "grange", "grange, grain storage room"},
		{"gogua", "gögüä", "gögüä", "NOUN", "", "ANIM", 6, "buffalo", "buffalo"},
		{"goigoi", "goigôî", "gɔigɔ̂î", "NOUN", "", "HOW", 4, "laziness", "laziness"},
		{"gon", "gön", "gön", "VERB", "", "SENSE", 4, "resonate", "resonate"},
		{"gonda", "gônda", "gônda", "NOUN", "", "INTERACT", 3, "praise", "praise"},
		{"gonda", "gônda", "gônda", "VERB", "", "INTERACT", 3, "praise", "praise"},
		{"goro", "gôro", "gɔ̂rɔ", "NOUN", "", "CIVIL", 4, "bribe", "[lit: kola nut]: bribe"},
		{"goro", "gôro", "gɔ̂rɔ", "NOUN", "", "FOOD", 4, "kola-nut", "kola nut"},
		{"goro", "gôro", "gɔ̂rɔ", "NOUN", "", "HOW", 4, "bitterness", "[lit: kola nut flavor]: bitterness"},
		{"gosa", "gôsâ", "gôsâ", "NOUN", "", "FOOD", 6, "eggplant", "eggplant"},
		{"goyongo", "gôyongö", "gôyongö", "NOUN", "", "PLANT", 6, "henna", "henna"},
		{"gozo", "gozo", "gɔzɔ", "NOUN", "", "FOOD", 2, "manioc-root", "manioc root"},
		{"gua", "gua", "gua", "NOUN", "", "SENSE", 4, "child-labor-pains", "child labor pains"},
		{"gua", "gûâ", "gûâ", "VERB", "", "ACT", 3, "hang", "hang"},
		{"guagua", "guagua", "guagua", "ADV", "", "SENSE", 4, "conceitedly", "conceitedly"},
		{"guagua", "güägüä", "güägüä", "NOUN", "", "ANIM", 6, "buffalo", "buffalo"},
		{"gue", "gue", "gue", "VERB", "Subcat=Intr", "MOVE", 1, "go", "go"},
		{"guena", "gue na", "gue na", "VERB", "Subcat=Tran", "MOVE", 1, "leave-with", "leave with"},
		{"guenzoni", "gue nzönî", "gue nzɔ̈nî", "INTERJ", "Mood=Opt", "INTERACT", 1, "Goodbye", "[lit: go|well]: Goodbye! (said to those leaving)"},
		{"gugu", "gügü", "gügü", "NOUN", "", "FOOD", 4, "mushroom", "mushroom"},
		{"guguma", "gügümä", "gügümä", "NOUN", "", "BODY", 4, "stutter", "stutter, stammer"},
		{"guguru", "gugûrû", "gugûrû", "NOUN", "", "FISH", 6, "sardine", "sardine"},
		{"gui", "gûî", "gûî", "NOUN", "", "FOOD", 4, "yam", "yam"},
		{"gumbaya", "gümbâyä", "gümbâyä", "ADJ", "NumType=Ord", "NUM", 2, "nine", "nine"},
		{"gumbaya", "gümbâyä", "gümbâyä", "NUM", "NumType=Card", "NUM", 2, "nine", "nine"},
		{"gunda", "gündâ", "gündâ", "NOUN", "", "BODY", 3, "base", "foot, base, root"},
		{"guru", "gûrû", "gûrû", "NOUN", "", "TREE", 5, "jackfruit-tree", "jackfruit tree"},
		{"guru", "gürü", "gürü", "NOUN", "", "NATURE", 3, "smoke", "smoke"},
		{"ha", "hä", "hä", "VERB", "Subcat=Tran", "ACT", 2, "open-wide", "open completely; braid, weave"},
		{"haa", "hâa", "hâa", "VERB", "Subcat=Tran", "ACT", 4, "try-on", "try on, measure, compare, evaluate"},
		{"haka", "hâka", "hâka", "VERB", "Subcat=Tran", "ACT", 3, "compare", "try on, measure, compare, evaluate"},
		{"hakango", "häkängö", "häkängɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 3, "comparison", "comparison"},
		{"hako", "häko", "häko", "NOUN", "", "CIVIL", 4, "temp-work", "temporary work"},
		{"hale", "halë", "halë", "NOUN", "", "ALT SP FOR", 3, "ancestry", "alë"},
		{"han", "hân", "hân", "VERB", "Subcat=Tran", "SICK", 3, "heal", "treat, cure, heal"},
		{"handa", "hânda", "hânda", "NOUN", "", "INTERACT", 3, "deception", "trick, deception"},
		{"handa", "hânda", "hânda", "VERB", "Subcat=Tran", "INTERACT", 3, "deceive", "trick, deceive"},
		{"hariya", "hâriya", "hâriya", "NOUN", "", "FOOD", 5, "millet", "fonio millet"},
		{"he", "he", "he", "VERB", "Subcat=Tran", "ALT SP FOR", 9, "laugh-at", "hë"},
		{"he", "hë", "hë", "VERB", "Subcat=Tran", "INTERACT", 2, "laugh-at", "laugh at, mock, ridicule"},
		{"hene", "hene", "hɛnɛ", "VERB", "", "INTERACT", 4, "exaggerate", "exaggerate"},
		{"hene", "hënë", "hɛ̈nɛ̈", "NOUN", "", "INTERACT", 4, "exaggeration", "exaggeration"},
		{"hengia", "he ngîâ", "he ngîâ", "VERB", "Subcat=Intr", "INTERACT", 2, "laugh", "laugh, amuse oneself"},
		{"hinga", "hînga", "hînga", "VERB", "Aspect=Iter", "SENSE", 2, "know", "know"},
		{"hingango", "hïngängö", "hïngängɔ̈", "VERB", "VerbForm=Vnoun", "SENSE", 2, "knowledge", "knowledge"},
		{"hini", "hîni", "hîni", "VERB", "Subcat=Tran", "ACT", 3, "paint", "tint, paint"},
		{"hini", "hîni", "hîni", "VERB", "Subcat=Tran", "GOD", 4, "annoint", "annoint"},
		{"hio", "hîo", "hîo", "ADV", "", "HOW", 1, "fast", "fast"},
		{"hiohio", "hîo-hîo", "hîo-hîo", "ADV", "", "HOW", 1, "very-fast", "very fast"},
		{"homba", "hömba", "hömba", "NOUN", "", "ALT SP FOR", 9, "relative", "wömba"},
		{"hon", "hôn", "hôn", "NOUN", "", "BODY", 2, "nose", "nose"},
		{"hon", "hön", "hön", "VERB", "", "ACT", 2, "pass", "pass, surpass, exceed"},
		{"honde", "hônde", "hɔ̂ndɛ", "VERB", "Subcat=Tran", "ACT", 3, "hide", "hide, conceal"},
		{"hondengo", "höndëngö", "hɔ̈ndɛ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 3, "secret", "secret"},
		{"hondesioye", "hônde-sïö-yê", "hɔ̂ndɛ-sïɔ̈-yê", "VERB", "Subcat=Intr", "GOD", 3, "forgive", "forgive"},
		{"hongere", "hôn-gerê", "hôn-gɛrɛ̂", "NOUN", "", "ANIM", 4, "caterpillar", "[lit: nose|foot]: caterpillar"},
		{"honndoti", "hön ndö tî", "hön ndö tî", "VERB", "Subcat=Tran", "HOW", 2, "surpass", "[lit: pass|place|of]: surpass, triumph over, vanquish, dominate"},
		{"honti", "hôn-tï", "hôn-tï", "NOUN", "", "BODY", 5, "wrist", "[lit: nose|arm]: wrist"},
		{"hu", "hû", "hû", "VERB", "Subcat=Tran", "SENSE", 3, "see", "see"},
		{"hule", "hûle", "hûlɛ", "VERB", "Subcat=Intr", "HOW", 2, "dry-out", "dry, dry out, dry up"},
		{"hule", "hülë", "hülë", "NOUN", "", "HOW", 2, "game", "[Sg: dry, because sports are played in the dry season]: game, sport"},
		{"hulengo", "hülëngö", "hülɛ̈ngɔ̈", "ADJ", "VerbForm=Vnoun", "HOW", 2, "dry", "dry, dried, scrawny"},
		{"hulengo", "hülëngö", "hülɛ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "HOW", 2, "dryness", "dryness"},
		{"hunda", "hûnda", "hûnda", "NOUN", "", "INTERACT", 1, "question", "question"},
		{"hunda", "hûnda", "hûnda", "VERB", "Subcat=Tran", "INTERACT", 1, "ask", "ask"},
		{"hunu", "hunu", "hunu", "VERB", "Subcat=Tran", "SENSE", 4, "sniff", "sniff"},
		{"hunzi", "hûnzi", "hûnzi", "VERB", "Subcat=Intr", "STATE", 1, "be-used-up", "end, be used up"},
		{"huru", "huru", "huru", "VERB", "Subcat=Intr", "MOVE", 4, "fly", "fly"},
		{"huru", "hürü", "hürü", "NOUN", "", "MOVE", 4, "flight", "flight"},
		{"i", "ï", "ï", "PRON", "Num=Plur|Person=2|PronType=Prs", "WHO", 1, "you", "[Catholic,formal]: you [plural]"},
		{"imveni", "ï mvenî", "ï mvɛnî", "PRON", "Num=Plur|Person=2|PronType=Prs", "WHO", 1, "yourselves", "[Catholic,formal]: you [plural]"},
		{"in", "in", "in", "INTERJ", "Polarity=Pos", "INTERACT", 1, "yes", "yes"},
		{"ingo", "îngö", "îngɔ̈", "NOUN", "", "FOOD", 3, "salt", "salt"},
		{"inin", "in-in", "in-in", "INTERJ", "Polarity=Neg", "INTERACT", 1, "no", "no"},
		{"ino", "înö", "înɔ̈", "NOUN", "", "BODY", 3, "urine", "urine"},
		{"iri", "îri", "îri", "VERB", "Subcat=Tran", "INTERACT", 1, "call", "call, name"},
		{"iri", "ïrï", "ïrï", "NOUN", "", "INTERACT", 1, "name", "name"},
		{"ita", "îtä", "îtä", "NOUN", "", "FAMILY", 2, "sibling", "brother, sister, fellow tribesman"},
		{"itabua", "îtä-buä", "îtä-buä", "NOUN", "", "GOD ", 3, "friar", "[lit: brother|priest]: friar"},
		{"ka", "ka", "ka", "CCONJ", "", "HOW", 1, "and", "and (between clauses)"},
		{"ka", "kâ", "kâ", "ADV", "", "WHERE", 1, "over-there", "over there"},
		{"ka", "kâ", "kâ", "CCONJ", "", "HOW", 1, "then", "then, in that case"},
		{"ka", "kä", "kä", "NOUN", "", "SICK", 4, "wound", "wound, sore, ulcer"},
		{"ka", "kä", "kä", "VERB", "Subcat=Tran", "CIVIL", 2, "sell", "sell, barter; betray"},
		{"kabi", "kâbî", "kâbî", "NOUN", "", "OBJ", 4, "head-pillow", "head pillow"},
		{"kabinee", "kabinêe", "kabinêe", "NOUN", "", "HOUSE", 3, "latrine", "[Fr: cabinet]: latrine"},
		{"kada", "kadâ", "kadâ", "NOUN", "", "ANIM", 6, "gecko", "gecko"},
		{"kafe", "kâfe", "kâfe", "NOUN", "", "PLANT", 5, "coffee", "coffee (plant, bean, grounds)"},
		{"kaga", "kaga", "kaga", "NOUN", "", "NATURE", 4, "hill", "hill, mountain"},
		{"kai", "kâi", "kâi", "VERB", "Subcat=Intr", "SICK", 3, "heal", "heal, get well, be cured, grow calm"},
		{"kai", "kâî", "kâî", "NOUN", "", "ACT", 4, "paddle", "paddle"},
		{"kaka", "kakâ", "kakâ", "NOUN", "", "FAMILY", 4, "maternal-grandfather", "[term of address]: maternal grandparent, old person"},
		{"kakara", "kakara", "kakara", "ADV", "", "CIVIL", 4, "evenly-matched", "evenly matched"},
		{"kakauka", "kakauka", "kakauka", "NOUN", "", "WHEN", 5, "December", "December"},
		{"kakere", "käkërë", "käkërë", "NOUN", "", "PLANT", 5, "seed-pod-plant", "seed pod plant"},
		{"kako", "kakö", "kakö", "NOUN", "", "PLANT", 4, "shell-or-husk", "(peanut)shell, (corn)husk, (seed)pod"},
		{"kakoro", "käkorö", "käkorö", "NOUN", "", "ANIM", 6, "anteater", "anteater"},
		{"kala", "kalâ", "kalâ", "NOUN", "", "ANIM", 5, "snail", "snail"},
		{"kalambo", "kalambo", "kalambo", "NOUN", "", "NATURE", 6, "lake", "lake"},
		{"kamata", "kamâta", "kamâta", "NOUN", "", "CIVIL", 4, "arrest", "seizure, arrest"},
		{"kamata", "kamâta", "kamâta", "VERB", "", "CIVIL", 4, "arrest", "seize, confiscate, arrest"},
		{"kamba", "kamba", "kamba", "NOUN", "", "OBJ", 3, "machete", "machete"},
		{"kamba", "kâmba", "kâmba", "NOUN", "", "OBJ", 3, "rope-belt-wire", "vine, fiber, cord, rope, belt, wire; [lit: fiber (of my soul)]: my love, the one I love"},
		{"kambiri", "kambîri", "kambîri", "ADJ", "", "COLOR", 3, "yellow", "[lit: cooking oil]: yellow"},
		{"kambiri", "kambîri", "kambîri", "NOUN", "", "FOOD", 3, "cooking-oil", "cooking oil"},
		{"kambisa", "kambisa", "kambisa", "VERB", "Subcat=Tran", "INTERACT", 4, "explain", "explain, prove, solve"},
		{"kambisa", "kambisä", "kambisä", "NOUN", "", "INTERACT", 4, "explanation", "explanation, proof, solution"},
		{"kambusu", "kämbûsu", "kämbûsu", "NOUN", "", "OBJ", 4, "loincloth", "loincloth"},
		{"kamene", "kamënë", "kamɛ̈nɛ̈", "NOUN", "", "CIVIL", 3, "shame", "shame"},
		{"kanana", "kanâna", "kanâna", "NOUN", "", "ANIM", 5, "duck", "duck"},
		{"kanda", "kandä", "kandä", "NOUN", "", "FOOD", 3, "manioc-termite-bar", "manioc termite bar"},
		{"kandaa", "kandâa", "kandâa", "SCONJ", "Mood=Irr", "HOW", 1, "however", "[lit: and|if-only|end] however"},
		{"kanga", "kânga", "kânga", "NOUN", "", "ACT", 3, "prison", "cover, enclosure, prison"},
		{"kanga", "kânga", "kânga", "VERB", "Aspect=Iter", "ACT", 3, "imprison", "cover, enclose, imprison; tie, attach"},
		{"kanga", "kângâ", "kângâ", "NOUN", "", "ANIM", 6, "antilope", "antilope"},
		{"kanga", "kângâ", "kângâ", "NOUN", "", "OBJ", 4, "loincloth-or-pickaxe", "loincloth; pick(axe)"},
		{"kangama", "kangamä", "kangamä", "NOUN", "", "OBJ", 4, "robe", "robe"},
		{"kangba", "kangba", "kangba", "ADJ", "", "FAMILY", 3, "adult", "[>21yrs or married or pregnant]: adult"},
		{"kangba", "kangba", "kangba", "NOUN", "", "FAMILY", 3, "adult", "[>21yrs or married or pregnant]: adult"},
		{"kangba", "kângbâ", "kângbâ", "NOUN", "", "HOUSE", 4, "metal-roof", "metal roof"},
		{"kangba", "kängbä", "kängbä", "NOUN", "", "ANIM", 5, "crab", "crab"},
		{"kangbi", "kângbi", "kângbi", "NOUN", "", "CIVIL", 3, "division", "division, separation, share, portion"},
		{"kangbi", "kângbi", "kângbi", "VERB", "Aspect=Imp|Subcat=Tran", "CIVIL", 3, "divide", "divide, separate, share"},
		{"kangbitere", "kângbi terê", "kângbi tɛrɛ̂", "VERB", "Subcat=Intr", "CIVIL", 3, "separate", "[lit: divide|oneself]: separate, break apart"},
		{"kangi", "kangi", "kangi", "NOUN", "", "ANIM", 3, "termite", "wingless soldier termite"},
		{"kango", "kängö", "kängɔ̈", "VERB", "VerbForm=Vnoun", "CIVIL", 3, "sale", "sale"},
		{"kangoya", "kangoya", "kangoya", "NOUN", "", "DRINK", 5, "oil-palm-wine", "oil palm wine"},
		{"kanya", "kânyâ", "kânyâ", "NOUN", "", "OBJ", 2, "fork", "fork"},
		{"kanza", "kanza", "kanza", "NOUN", "", "WHAT", 4, "raw-materials", "raw materials"},
		{"kanzago", "kanzagö", "kanzagɔ̈", "NOUN", "", "OBJ", 5, "woman's-blouse", "woman's blouse"},
		{"kapi", "kapï", "kapï", "NOUN", "", "ANIM", 6, "mongoose", "mongoose"},
		{"kapitani", "kapitäni", "kapitäni", "NOUN", "", "FISH", 6, "Nile-perch", "Nile perch, [French]: capitaine"},
		{"kara", "kara", "kara", "VERB", "Subcat=Tran", "SENSE", 4, "weigh-down", "weigh down, embarrass, overtax"},
		{"kara", "kâra", "kâra", "VERB", "", "ACT", 4, "demolish", "break, demolish, uproot"},
		{"kara", "kârâ", "kârâ", "ADJ", "", "HOW", 3, "tightly-closed", "tightly closed"},
		{"karagba", "karagba", "karagba", "NOUN", "", "ANIM", 5, "grasshopper", "grasshopper"},
		{"karagoro", "kârâgorö", "kârâgɔrɔ̈", "NOUN", "", "ANIM", 6, "pidgeon", "green pidgeon"},
		{"karako", "kârâkö", "kârâkö", "NOUN", "", "FOOD", 3, "peanut", "peanut"},
		{"karangba", "karangbâ", "karangbâ", "NOUN", "", "OBJ", 5, "xylophone", "xylophone"},
		{"kasa", "kâsa", "kâsa", "NOUN", "", "FOOD", 3, "side-dish", "side dish; prey"},
		{"kasakasa", "kasakasa", "kasakasa", "NOUN", "", "WHEN", 6, "April", "April"},
		{"kasi", "kasï", "kasï", "SCONJ", "", "HOW", 1, "but", "but"},
		{"kate", "kate", "katɛ", "NOUN", "", "BODY", 3, "ribs", "ribs, chest, trunk"},
		{"katikati", "katikâti", "katikâti", "NOUN", "", "CIVIL", 4, "border", "limit, border"},
		{"katisima", "kätîsima", "kätîsima", "NOUN", "", "GOD", 5, "catechism", "catechism"},
		{"kawa", "kâwa", "kâwa", "NOUN", "", "DRINK", 2, "coffee", "coffee"},
		{"kawoya", "kawoya", "kawoya", "NOUN", "", "PLANT", 4, "pumpkin", "pumpkin"},
		{"kaye", "kayë", "kayë", "NOUN", "", "PLANT", 4, "peel", "(banana) peel, (fish) scales, (snail) shell"},
		{"kayee", "kayëe", "kayëe", "NOUN", "", "OBJ", 4, "notebook", "[Fr: cahier]: notebook"},
		{"ke", "ke", "kɛ", "VERB", "", "SENSE", 1, "reject", "refuse, deny, reject, divorce"},
		{"keke", "këkë", "kɛ̈kɛ̈", "NOUN", "", "PLANT", 2, "tree-or-wood", "tree, trunk, branch, wood, stick"},
		{"kekere", "kekere", "kɛkɛrɛ", "NOUN", "", "ANIM", 5, "ant", "ant"},
		{"kekereke", "kêkerêke", "kêkerêke", "NOUN", "", "WHEN", 2, "tomorrow", "tomorrow"},
		{"kele", "kêlê", "kêlê", "VERB", "Subcat=Tran", "BODY", 4, "blind", "blind"},
		{"kelele", "kêlêlê", "kɛ̂lɛ̂lɛ̂", "NOUN", "", "OBJ", 5, "lock-or-key", "[Fr: clé]: lock, key"},
		{"kema", "kêma", "kêma", "NOUN", "", "ANIM", 6, "monkey", "monkey"},
		{"kembe", "kembe", "kɛmbɛ", "NOUN", "", "WHERE", 5, "Kembe", "Kembe"},
		{"kenda", "kênda", "kênda", "NOUN", "", "BODY", 4, "corpse", "corpse, carcass"},
		{"kene", "kêne", "kɛ̂nɛ", "VERB", "Subcat=Tran", "ACT", 4, "roll-up", "roll up"},
		{"kene", "kënë", "kɛ̈nɛ̈", "NOUN", "", "DRINK", 5, "manioc-wine", "manioc wine"},
		{"kenge", "këngë", "kɛ̈ngɛ̈", "NOUN", "", "BODY", 5, "penis", "penis"},
		{"kengere", "kengêre", "kengêre", "NOUN", "", "SENSE", 4, "supposition", "supposition, speculation"},
		{"kengo", "këngö", "kɛ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 1, "rejection", "rejection, refusal"},
		{"kepaka", "kepaka", "kepaka", "NOUN", "Gender=Masc", "WHO", 6, "Mister", "Mister, honest man"},
		{"kepakara", "kepakara", "kepakara", "NOUN", "Gender=Masc", "WHO", 6, "Mister", "Mister, honest man"},
		{"kere", "kêrë", "kɛ̂rɛ̈", "NOUN", "", "SICK", 4, "heartburn", "heartburn"},
		{"kerebende", "kerebende", "kerebende", "ADJ", "", "OBJ", 4, "round", "round"},
		{"kerebende", "kerebende", "kerebende", "NOUN", "", "OBJ", 4, "circle", "circle"},
		{"kerekpa", "kerekpa", "kerekpa", "NOUN", "", "HOUSE", 4, "rattan-bed", "traditional bed made of rattan"},
		{"kerekpa", "kerekpä", "kerekpä", "NOUN", "", "CIVIL", 4, "mutual-savings-plan", "mutual savings plan"},
		{"kete", "kêtê", "kɛ̂tɛ̂", "ADJ", "", "NUM", 2, "small", "little, few, small, short"},
		{"kete", "kêtê", "kɛ̂tɛ̂", "ADV", "", "NUM", 2, "little-bit", "a little bit"},
		{"kete", "kêtê", "kɛ̂tɛ̂", "NOUN", "", "NUM", 2, "younger-person", "person younger than you"},
		{"ketebaba", "kêtê-babâ", "kɛ̂tɛ̂-babâ", "NOUN", "Gender=Masc", "FAMILY", 4, "paternal-uncle", "paternal uncle (father's younger brother)"},
		{"keteita", "kêtê-îtä", "kɛ̂tɛ̂-îtä", "NOUN", "", "FAMILY", 4, "younger-sibling", "[lit: little|sibling]: younger brother, younger sister"},
		{"ketemama", "kêtê-mamâ", "kɛ̂tɛ̂-mamâ", "NOUN", "Gender=Fem", "FAMILY", 4, "maternal-aunt", "maternal aunt (mother's younger sister)"},
		{"ki", "kî", "kî", "NOUN", "", "ANIM", 4, "spine", "(porcupine or plant) spine, quill"},
		{"ki", "kî", "kî", "VERB", "Subcat=Tran", "ACT", 3, "build", "build, construct"},
		{"kiki", "kîki", "kîki", "VERB", "Subcat=Tran", "ACT", 4, "tickle", "tickle"},
		{"kinda", "kinda", "kinda", "VERB", "", "CIVIL", 4, "defeat", "defeat, knock down"},
		{"kindanda", "kindânda", "kindânda", "NOUN", "", "MUSIC", 5, "accordion", "accordion"},
		{"kindango", "kïndängö", "kïndängɔ̈", "VERB", "VerbForm=Vnoun", "CIVIL", 4, "defeat", "defeat, knock-down"},
		{"kinde", "kindë", "kindë", "NOUN", "", "OBJ", 4, "club", "club, truncheon"},
		{"kindere", "kîndêrê", "kîndɛ̂rɛ̂", "VERB", "Subcat=Intr", "STATE", 4, "be-submerged", "be submerged"},
		{"kinini", "kinîni", "kinîni", "NOUN", "", "SICK", 4, "quinine", "quinine"},
		{"kio", "kîo", "kîɔ", "VERB", "", "ACT", 4, "shave", "shave, scrape, grate"},
		{"kiri", "kîri", "kîri", "VERB", "Subcat=Intr", "MOVE", 1, "return", "return, respond, lower (price)"},
		{"kirikiri", "kîrîkiri", "kîrîkiri", "ADV", "", "CIVIL", 2, "disorderly", "disorderly, on the wrong track"},
		{"kirikiri", "kîrîkiri", "kîrîkiri", "NOUN", "", "CIVIL", 2, "disorderliness", "disorderliness, carelessness"},
		{"kiringo", "kïrïngö", "kïrïngɔ̈", "VERB", "VerbForm=Vnoun", "MOVE", 1, "return", "return"},
		{"kiro", "kîrô", "kîrɔ̂", "NOUN", "", "ACT", 4, "clay-pan", "frying pan made of baked clay"},
		{"kisoro", "kisoro", "kisɔrɔ", "NOUN", "", "GAME", 5, "board-game", "board game moving stones around an egg carton"},
		{"kite", "kîte", "kîtɛ", "NOUN", "", "SENSE", 3, "doubt", "doubt"},
		{"kiti", "kîti", "kîti", "NOUN", "", "HOUSE", 6, "chair", "chair"},
		{"kizi", "kîzi", "kîzi", "NOUN", "", "NATURE", 5, "pearl", "pearl"},
		{"ko", "kô", "kô", "VERB", "Subcat=Intr", "ACT", 3, "go-up-or-down", "go up or down, embark, disembark"},
		{"ko", "kô", "kɔ̂", "VERB", "Subcat=Tran", "ACT", 3, "gather", "gather, pick (fruit)"},
		{"ko", "kö", "kö", "NOUN", "", "STATE", 2, "utmost", "utmost; root, growth"},
		{"ko", "kö", "kö", "VERB", "Subcat=Intr", "STATE", 3, "grow", "grow, bear fruit"},
		{"kobe", "kôbe", "kɔ̂bɛ", "NOUN", "", "FOOD", 2, "food", "food"},
		{"kobela", "kobêla", "kobêla", "NOUN", "", "SICK", 3, "illness", "illness"},
		{"kobelatiwa", "kobêla tî wâ", "kobêla tî wâ", "NOUN", "", "SICK", 3, "fever", "fever, malaria"},
		{"kode", "kodë", "kɔdɛ̈", "NOUN", "", "FEEL", 5, "-istics", "cleverness, ingenuity, field of study"},
		{"kodekua", "kodëkua", "kɔdɛ̈kua", "NOUN", "", "FEEL", 5, "technique", "[lit: cleverness|work] technique"},
		{"kodoro", "ködörö", "kɔ̈dɔ̈rɔ̈", "NOUN", "", "CIVIL", 2, "village", "village, neighborhood"},
		{"kodorosese", "ködörö-sêse", "kɔ̈dɔ̈rɔ̈-sêse", "NOUN", "", "CIVIL", 2, "republic", "republic, state"},
		{"kogara", "kögarä", "kɔ̈garä", "NOUN", "Gender=Masc", "FAMILY", 4, "in-laws", "father-in-law, parents-in-law"},
		{"koka", "kôkâ", "kôkâ", "NOUN", "", "GAME", 5, "dice-game", "dice game"},
		{"koko", "koko", "kɔkɔ", "NOUN", "", "FOOD", 3, "sliced-greens", "plant leaves finely cut steamed or fried"},
		{"koko", "kôko", "kôko", "NOUN", "", "FISH", 6, "catfish", "upside down catfish"},
		{"koko", "kôkô", "kɔ̂kɔ̂", "NOUN", "", "ANIM", 5, "lizard", "lizard"},
		{"kokombe", "kokombe", "kɔkɔmbɛ", "NOUN", "", "SICK", 4, "yaws", "yaws"},
		{"kokora", "kokora", "kɔkɔra", "NOUN", "", "OBJ", 4, "arrow", "arrow, dart"},
		{"koli", "koli", "koli", "NOUN", "", "OBJ", 4, "pillow", "pillow, cushion"},
		{"koli", "kôlï", "kɔ̂lï", "NOUN", "", "WHERE", 2, "right-side", "right side"},
		{"koli", "kôlï", "kɔ̂lï", "NOUN", "Gender=Masc", "FAMILY", 2, "husband", "husband"},
		{"koli", "kôlï", "kɔ̂lï", "NOUN", "Gender=Masc", "WHO", 1, "man", "man, male"},
		{"kolikoli", "kôlï-kôlï", "kɔ̂lï-kɔ̂lï", "NOUN", "Gender=Masc", "WHO", 6, "gay-man", "[lit: man|man]: 'butch' gay man"},
		{"kolingo", "kolîngo", "kolîngo", "NOUN", "", "ANIM", 6, "chamelion", "chamelion"},
		{"koliti", "kô-li-tï", "kɔ̂-li-tï", "NOUN", "", "BODY", 4, "middle-finger", "[lit: male|finger]: middle finger"},
		{"koliwali", "kôlï-wâlï", "kɔ̂lï-wâlï", "NOUN", "Gender=Masc", "WHO", 6, "gay-man", "[lit: man|woman]: 'fem' gay man"},
		{"kolo", "kôlo", "kôlo", "NOUN", "", "ANIM", 3, "giraffe", "giraffe"},
		{"kolofia", "kolôfîa", "kolôfîa", "NOUN", "", "ANIM", 6, "shrew", "shrew"},
		{"kolokoto", "kôlökôtö", "kɔ̂lɔ̈kɔ̂tɔ̈", "NOUN", "", "ANIM", 6, "turtledove", "turtledove"},
		{"kolongo", "kolongo", "kɔlɔngɔ", "NOUN", "", "TREE", 5, "palm-tree", "sugar palm, fan palm"},
		{"kolongo", "kolôngo", "kolôngo", "NOUN", "", "OBJ", 4, "wooden-bowl", "wooden bowl"},
		{"kombe", "kömbë", "kömbë", "ADJ", "", "COLOR", 6, "yellow", "[lit: yellowfruit]: yellow"},
		{"kombe", "kömbë", "kömbë", "NOUN", "", "TREE", 6, "yellowfruit", "yellowfruit"},
		{"kombuka", "kombûka", "kombûka", "NOUN", "", "ACT", 4, "uprising", "revolt, uprising"},
		{"kombuka", "kombûka", "kombûka", "VERB", "Subcat=Intr", "ACT", 4, "revolt", "revolt"},
		{"kome", "kome", "kɔmɛ", "NOUN", "", "ANIM", 6, "lizard", "Nile monitor lizard"},
		{"kondo", "kôndo", "kɔ̂ndɔ", "NOUN", "", "ANIM", 2, "chicken", "chicken"},
		{"konga", "konga", "konga", "VERB", "Aspect=Iter|Subcat=Tran", "ACT", 3, "select", "select, filter"},
		{"konga", "kongä", "kongä", "NOUN", "", "ACT", 3, "selection", "selection, filter"},
		{"kongba", "kongba", "kongba", "NOUN", "", "ANIM", 6, "toad", "toad"},
		{"kongba", "kongba", "kongba", "NOUN", "", "HOW", 4, "eccentric", "eccentric"},
		{"kongba", "köngbä", "köngbä", "NOUN", "", "OBJ", 6, "bellows", "blacksmith bellows"},
		{"kongo", "kongö", "kongö", "NOUN", "", "ANIM", 4, "parrot", "parrot"},
		{"kongo", "kongö", "kongö", "NOUN", "", "PLANT", 4, "flower", "flower"},
		{"kongo", "kongö", "kɔngɔ̈", "NOUN", "", "NATURE", 4, "rainbow", "rainbow"},
		{"kongo", "kôngô", "kôngô", "NOUN", "", "OBJ", 4, "hoe", "hoe"},
		{"kongo", "köngö", "köngö", "NOUN", "", "ACT", 4, "dam-fishing", "dam fishing"},
		{"kongo", "köngö", "kɔ̈ngɔ̈", "NOUN", "VerbForm=Vnoun", "INTERACT", 2, "exclamation", "cry, exclamation"},
		{"kono", "kono", "kɔnɔ", "VERB", "Subcat=Intr", "HOW", 2, "get-big", "grow up, get big, large, fat"},
		{"kono", "konô", "kɔnɔ̂", "NOUN", "", "ANIM", 3, "hippopotamus", "hippopotamus"},
		{"konongo", "könöngö", "kɔ̈nɔ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "HOW", 2, "large-size", "large size, grandeur"},
		{"konza", "konza", "konza", "NOUN", "", "OBJ", 4, "mat", "mat"},
		{"konzongoro", "konzöngörö", "kɔnzɔ̈ngɔ̈rɔ̈", "NOUN", "", "ANIM", 6, "lizard", "blue orange lizard"},
		{"kopo", "kopo", "kopo", "NOUN", "", "HOUSE", 3, "metal-roof-or-can", "metal roof, metal can"},
		{"kopo", "köpö", "kɔ̈pɔ̈", "NOUN", "", "HOUSE", 3, "cup", "cup, goblet"},
		{"koro", "koro", "kɔrɔ", "VERB", "", "SICK", 4, "cough", "cough, catch cold"},
		{"koro", "kôro", "kôro", "VERB", "Subcat=Tran", "ACT", 3, "pierce", "pierce, dig; de-louse"},
		{"koro", "körö", "kɔ̈rɔ̈", "NOUN", "", "SICK", 4, "cough", "cold, cough"},
		{"korobo", "korobö", "korobö", "NOUN", "", "BODY", 5, "testicles", "testicles"},
		{"korokongbo", "korôkongbô", "korôkongbô", "NOUN", "", "WHO", 5, "leprechaun", "leprechaun"},
		{"koromenge", "körömëngë", "körömëngë", "NOUN", "", "SICK", 6, "hernia", "hernia"},
		{"kororo", "korôro", "korôro", "NOUN", "", "ANIM", 5, "donkey", "donkey"},
		{"kosala", "kosâla", "kosâla", "NOUN", "", "ALT SP FOR", 9, "job", "kusâra"},
		{"kosara", "kosâra", "kosâra", "NOUN", "", "ALT SP FOR", 9, "job", "kusâra"},
		{"koso", "koso", "kɔsɔ", "NOUN", "", "ANIM", 2, "pork", "pork"},
		{"koso", "kôso", "kɔ̂sɔ", "VERB", "Subcat=Tran", "ACT", 4, "pull-towards-oneself", "pull towards oneself"},
		{"koso", "kôsö", "kɔ̂sɔ̈", "NOUN", "", "FOOD", 4, "squash", "squash, cucumber, melon"},
		{"kosotingonda", "koso tî ngonda", "kɔsɔ tî ngonda", "NOUN", "", "ANIM", 2, "boar", "boar"},
		{"kota", "kötä", "kötä", "ADJ", "", "NUM", 1, "big", "big, large, tall"},
		{"kota", "kötä", "kötä", "NOUN", "", "NUM", 1, "older-person", "person older than you"},
		{"kotababa", "kötä-babâ", "kötä-babâ", "NOUN", "Gender=Masc", "FAMILY", 4, "paternal-uncle", "paternal uncle (father's older brother)"},
		{"kotabe", "kötä-bë", "kötä-bɛ̈", "NOUN", "", "FEEL", 1, "envy", "[lit: big|heart]: envy, jealousy"},
		{"kotabua", "kötä-buä", "kötä-buä", "NOUN", "", "GOD ", 3, "bishop", "[lit: big|priest]: bishop"},
		{"kotaita", "kötä-îtä", "kötä-îtä", "NOUN", "", "FAMILY", 4, "older-sibling", "[lit: big|sibling]: older brother, older sister"},
		{"kotamama", "kötä-mamâ", "kötä-mamâ", "NOUN", "Gender=Fem", "FAMILY", 4, "maternal-aunt", "maternal aunt (mother's older sister)"},
		{"kotangu", "kötä ngû", "kötä ngû", "NOUN", "", "NATURE", 1, "river", "river"},
		{"kotara", "kötarä", "kötarä", "NOUN", "Gender=Masc", "FAMILY", 4, "paternal-grandfather", "paternal grandfather"},
		{"kotazo", "kötä-zo", "kötä-zo", "NOUN", "", "WHO", 1, "VIP-or-master", "[lit: big|person]: important person, master"},
		{"koti", "kötï", "kɔ̈tï", "ADV", "", "WHERE", 4, "right-hand", "right hand, right side"},
		{"koto", "koto", "koto", "NOUN", "", "NATURE", 3, "mound", "mound"},
		{"koto", "koto", "kɔtɔ", "NOUN", "", "HOUSE", 4, "house-foundation", "foundation of a house"},
		{"koto", "koto", "kɔtɔ", "VERB", "Subcat=Tran", "ACT", 4, "scratch", "claw, scratch"},
		{"koto", "kôto", "kôto", "NOUN", "", "BODY", 5, "Adam's-apple", "Adam's apple"},
		{"kotoon", "kotöon", "kɔtɔ̈ɔn", "NOUN", "", "PLANT", 5, "cotton", "cotton"},
		{"koya", "kôya", "kôya", "NOUN", "", "FAMILY", 4, "maternal-relative", "maternal uncle/niece/nephew"},
		{"kozo", "kôzo", "kɔ̂zɔ", "ADJ", "", "WHEN", 1, "first", "first"},
		{"kozo", "kôzo", "kɔ̂zɔ", "NOUN", "", "FAMILY", 4, "oldest", "[lit: first]: firstborn child, oldest"},
		{"kozoni", "kôzonî", "kɔ̂zɔnî", "ADV", "", "WHEN", 1, "first-of-all", "beforehand, first of all"},
		{"kozoti", "kôzo tî", "kɔ̂zɔ tî", "ADP", "", "WHEN", 1, "before", "before"},
		{"kpa", "kpa", "kpa", "VERB", "Subcat=Tran", "ACT", 2, "resemble", "resemble; scratch"},
		{"kpa", "kpâ", "kpâ", "NOUN", "", "BODY", 3, "hair", "hair"},
		{"kpaa", "kpâa", "kpâa", "ADV", "", "WHEN", 1, "only-just-now", "only just now"},
		{"kpaa", "kpâa", "kpâa", "NOUN", "", "ANIM", 6, "dwarf-monkey", "dwarf monkey"},
		{"kpaka", "kpaka", "kpaka", "VERB", "", "ACT", 3, "shave", "shave, scrape, grate"},
		{"kpakata", "kpaka-ta", "kpaka-ta", "NOUN", "", "OBJ", 3, "saucepan", "[lit: grate|pan]: saucepan, pot, kettle"},
		{"kpakpa", "kpäkpä", "kpäkpä", "NOUN", "", "OBJ", 4, "soap", "soap"},
		{"kpalakongo", "kpâlâköngö", "kpâlâköngö", "NOUN", "", "ANIM", 4, "scorpion", "scorpion"},
		{"kpale", "kpälë", "kpälë", "NOUN", "", "INTERACT", 4, "declaration", "declaration"},
		{"kpangaba", "kpängäbä", "kpängäbä", "NOUN", "", "PLANT", 5, "root", "root, manioc root"},
		{"kpangba", "kpangba", "kpangba", "NOUN", "", "OBJ", 3, "machete", "machete"},
		{"kpangbara", "kpangbara", "kpangbara", "ADJ", "", "HOW", 3, "wide", "wide"},
		{"kpangbara", "kpangbara", "kpangbara", "NOUN", "", "OBJ", 3, "machete", "[? from kpangba=machete + ra=iterative(ly)]: machete"},
		{"kpangbara", "kpângbârâ", "kpângbârâ", "ADJ", "", "HOW", 3, "flat", "flat"},
		{"kpangbara", "kpängbärä", "kpängbärä", "NOUN", "", "ANIM", 6, "bat", "[? from kpângi=wing + badâ=squirrel] bat"},
		{"kpangi", "kpângi", "kpângi", "NOUN", "", "BODY", 4, "wing", "wing"},
		{"kpata", "kpätä", "kpätä", "NOUN", "", "DRINK", 4, "corn-beer", "boiled corn beer; boiled termite drink"},
		{"kpe", "kpê", "kpɛ̂", "NOUN", "", "INTERACT", 4, "honor", "respect, honor"},
		{"kpe", "kpê", "kpɛ̂", "NOUN", "", "MOVE", 3, "fleeing", "moving about, traffic; run, flight, escape, avoidance, desertion"},
		{"kpe", "kpë", "kpë", "NOUN", "", "FOOD", 3, "nut-butter", "nut butter"},
		{"kpe", "kpë", "kpɛ̈", "VERB", "", "INTERACT", 4, "honor", "respect, honor"},
		{"kpe", "kpë", "kpɛ̈", "VERB", "", "MOVE", 3, "flee", "move about, circulate; run, flee, escape, avoid, desert"},
		{"kpee", "kpêe", "kpêe", "VERB", "Subcat=Intr", "FOOD", 4, "ferment", "ferment, be acidic"},
		{"kpeke", "kpêkê", "kpêkê", "NOUN", "", "COMPUTER", 5, "mouse-click", "(mouse) click"},
		{"kpekeuse", "kpêkê-ûse", "kpêkê-ûse", "NOUN", "", "COMPUTER", 5, "mouse-double-click", "(mouse) double-click"},
		{"kpeli", "kpë-li", "kpë-li", "NOUN", "", "FOOD", 5, "brain", "[lit: butter|head]: brain"},
		{"kpembeto", "kpë mbeto", "kpɛ̈ mbɛtɔ", "NOUN", "", "FEEL", 3, "fear", "fear, revere"},
		{"kpenda", "kpenda", "kpenda", "VERB", "Subcat=Tran", "ACT", 4, "compress", "compress"},
		{"kpengba", "kpëngba", "kpëngba", "VERB", "Subcat=Intr", "HOW", 3, "be-hard-or-strong", "be hard, solid, strong, serious, important"},
		{"kpengba", "kpëngbä", "kpëngbä", "ADJ", "", "HOW", 3, "hard-or-strong", "hard, solid, strong, serious, important"},
		{"kpengbango", "kpëngbängö", "kpëngbängɔ̈", "VERB", "VerbForm=Vnoun", "HOW", 3, "hardness-or-strength", "hardness, solidity, strength, seriousness, importance"},
		{"kpengbere", "kpëngbërë", "kpɛ̈ngbɛ̈rɛ̈", "NOUN", "", "NATURE", 4, "savanna", "savanna"},
		{"kpere", "kpere", "kpere", "NOUN", "", "ANIM", 6, "antilope", "antilope"},
		{"kperekpere", "kperekpere", "kpɛrɛkpɛrɛ", "ADV", "", "HOW", 4, "garrulously", "garrulously"},
		{"kpete", "kpete", "kpɛtɛ", "NOUN", "", "FISH", 6, "elephantfish", "elephantfish"},
		{"kpikara", "kpîkara", "kpîkara", "NOUN", "", "ANIM", 6, "anteater", "anteater"},
		{"kpo", "kpo", "kpɔ", "VERB", "Subcat=Tran", "ACT", 3, "pierce", "pierce; plant; thatch"},
		{"kpo", "kpô", "kpɔ̂", "ADJ", "", "HOW", 3, "silent", "calm, silent"},
		{"kpoka", "kpöka", "kpöka", "NOUN", "", "OBJ", 4, "hoe", "hoe"},
		{"kpokpo", "kpôkpô", "kpôkpô", "NOUN", "", "OBJ", 4, "pipe", "pipe"},
		{"kporo", "kporo", "kpɔrɔ", "VERB", "", "ACT", 3, "boil", "boil"},
		{"kpoto", "kpoto", "kpɔtɔ", "NOUN", "", "OBJ", 3, "hat-or-hairstyle", "hat, hairstyle"},
		{"kpu", "kpu", "kpu", "NOUN", "", "OBJ", 3, "mortar", "[onomatopeia]: mortar"},
		{"kpu", "kpû", "kpû", "NOUN", "", "OBJ", 6, "connection", "connection"},
		{"kpukangbi", "kpû-kângbi", "kpû-kângbi", "NOUN", "", "OBJ", 6, "dash", "[lit: connection|divide]: dash (punctuation)"},
		{"kpukpu", "kpûkpû", "kpûkpû", "NOUN", "", "OBJ", 3, "motorcycle", "[onomatopeia]: motorcycle"},
		{"kpuku", "kpûkû", "kpûkû", "NOUN", "", "INTERACT", 4, "riddle", "riddle"},
		{"kpunakpu", "kpunakpu", "kpunakpu", "ADJ", "", "WHEN", 5, "eternal", "eternal"},
		{"kpunakpu", "kpunakpu", "kpunakpu", "ADV", "", "WHEN", 5, "forever", "forever, eternally, indefinitely"},
		{"kputa", "kpütä", "kpütä", "NOUN", "", "FISH", 6, "snakehead-fish", "snakehead fish"},
		{"kputengbi", "kpû-têngbi", "kpû-tɛ̂ngbi", "NOUN", "", "OBJ", 6, "hyphen", "[lit: connection|join]: hyphen (punctuation)"},
		{"ku", "ku", "ku", "VERB", "Subcat=Intr", "BODY", 3, "spit", "spit"},
		{"ku", "kü", "kü", "VERB", "", "STATE", 1, "wait-for", "wait, wait for"},
		{"kua", "kua", "kua", "NOUN", "", "CIVIL", 2, "work", "work, job, duty"},
		{"kua", "kûâ", "kûâ", "NOUN", "", "STATE", 2, "death", "death"},
		{"kua", "küä", "küä", "NOUN", "", "BODY", 3, "hair", "hair, fur, pelt, feathers, down"},
		{"kuale", "kualë", "kualë", "NOUN", "", "ANIM", 6, "partridge", "partridge"},
		{"kue", "kûê", "kûɛ̂", "ADV", "", "NUM", 1, "completely", "completely"},
		{"kugbe", "kugbë", "kugbë", "NOUN", "", "PLANT", 3, "leaf", "leaf, leafy vegetable, sheet (of paper)"},
		{"kuii", "kûîi", "kûîi", "NOUN", "", "STATE", 2, "dying", "dying"},
		{"kuii", "kûîi", "kûîi", "VERB", "", "STATE", 2, "die", "die"},
		{"kuku", "kûku", "kûku", "NOUN", "", "HOUSE", 3, "cooking", "cooking"},
		{"kuku", "kûku", "kûku", "VERB", "Subcat=Intr", "MOVE", 4, "kneel", "kneel"},
		{"kukuru", "kûkurû", "kûkurû", "NOUN", "", "BODY", 4, "wig", "wig"},
		{"kukuru", "kûkürû", "kûkürû", "NOUN", "", "WHEN", 5, "August", "August"},
		{"kukuru", "kükürü", "kükürü", "NOUN", "", "FOOD", 5, "cucumber", "cucumber"},
		{"kulu", "kulü", "kulü", "NOUN", "", "FOOD", 5, "baby-formula", "baby formula"},
		{"kuma", "kûma", "kûma", "NOUN", "", "ANIM", 6, "python", "python"},
		{"kunda", "kunda", "kunda", "VERB", "Subcat=Tran", "MOVE", 4, "pull-up", "pull up what was slowly falling down, hitch up"},
		{"kunda", "kundâ", "kundâ", "NOUN", "", "ANIM", 5, "turtle", "turtle"},
		{"kundi", "kundi", "kundi", "NOUN", "", "OBJ", 5, "harp", "harp"},
		{"kungba", "kûngbâ", "kûngbâ", "NOUN", "", "OBJ", 2, "baggage", "baggage"},
		{"kungbi", "kûngbi", "kûngbi", "NOUN", "", "ACT", 3, "breakage", "breakage, debris"},
		{"kungbi", "kûngbi", "kûngbi", "VERB", "Aspect=Imp", "ACT", 3, "break", "break, shatter"},
		{"kungu", "kûngü", "kûngü", "NOUN", "", "PLANT", 6, "flower", "flower"},
		{"kupu", "kupu", "kupu", "NOUN", "", "TREE", 5, "Kapok-tree", "Kapok tree"},
		{"kura", "kürä", "kürä", "NOUN", "", "SENSE", 4, "spite", "spite, grudge"},
		{"kuru", "kürü", "kürü", "ADJ", "", "HOW", 3, "dry", "dry, brusque"},
		{"kurukuru", "kûrûkürü", "kûrûkürü", "NOUN", "", "FOOD", 5, "peanut-brittle", "peanut brittle"},
		{"kurungu", "kürüngü", "kürüngü", "NOUN", "", "ANIM", 6, "bluebird", "blue Turaco bird"},
		{"kusala", "kusâla", "kusâla", "NOUN", "", "ALT SP FOR", 9, "job", "kusâra"},
		{"kusara", "kusâra", "kusâra", "NOUN", "", "ACT", 2, "job", "job, work, profession, power"},
		{"kutu", "kûtu", "kûtu", "ADJ", "NumType=Ord", "NUM", 2, "million", "million"},
		{"kutu", "kûtu", "kûtu", "NOUN", "", "NUM", 3, "knot", "knot, hump, bump, boil"},
		{"kutu", "kûtu", "kûtu", "NUM", "NumType=Card", "NUM", 2, "million", "million"},
		{"kutugere", "kûtu-gerë", "kûtu-gɛrɛ̈", "NOUN", "", "BODY", 3, "ankle", "[lit: knot-foot] ankle"},
		{"kutukutu", "kutukutu", "kutukutu", "NOUN", "", "OBJ", 3, "car-or-truck", "[onomatopeia]: car, truck"},
		{"kuzu", "kuzü", "kuzü", "NOUN", "", "BODY", 4, "death", "death, cadaver"},
		{"la", "lâ", "lâ", "NOUN", "", "NATURE", 2, "sun", "sun"},
		{"la", "lâ", "lâ", "NOUN", "", "WHEN", 2, "day-or-daytime", "[lit: sun]: day[when], daytime"},
		{"laa", "laâ", "laâ", "PART", "", "INTERACT", 2, "behold", "behold!"},
		{"labada", "lâbâdâ", "lâbâdâ", "NOUN", "", "SICK", 6, "yaws", "yaws (illness of hands or feet)"},
		{"lagbada", "lägbädä", "lägbädä", "NOUN", "", "OBJ", 5, "tambourine", "[European-made]: tambourine"},
		{"lai", "lâi", "lâi", "NOUN", "", "FOOD", 3, "garlic", "[Fr: l'ail]: garlic"},
		{"lakere", "lakërë", "lakërë", "NOUN", "", "NATURE", 4, "drying-rock", "large rock or cleanly-swept dirt for drying clothes or manioc"},
		{"lakpangba", "lakpängbä", "lakpängbä", "ADJ", "", "BODY", 4, "bald", "bald"},
		{"lakue", "lâkûê", "lâkûɛ̂", "ADV", "", "WHEN", 1, "always", "always"},
		{"lakuelakue", "lâkûê-lâkûê", "lâkûɛ̂-lâkûɛ̂", "ADV", "", "WHEN", 1, "continually", "continually, constantly"},
		{"lakui", "lâ-kûî", "lâ-kûî", "NOUN", "", "WHEN", 2, "sunset", "sunset, dusk, evening"},
		{"lamba", "lâmbâ", "lâmbâ", "ADJ", "", "HOW", 4, "threadbare", "threadbare"},
		{"lando", "lando", "lando", "NOUN", "", "NATURE", 4, "marshland", "marshland, stadium"},
		{"langa", "langä", "langä", "NOUN", "", "PLANT", 5, "taro", "taro"},
		{"lango", "längö", "längɔ̈", "NOUN", "", "STATE", 1, "sleep", "sleep"},
		{"lango", "längö", "längɔ̈", "NOUN", "", "WHEN", 1, "day", "day [how long]"},
		{"lango", "längö", "längɔ̈", "VERB", "Subcat=Intr", "STATE", 1, "lie-down-or-stay", "lie down, stay, sleep"},
		{"lani", "lâ-nî", "lâ-nî", "ADV", "", "WHEN", 2, "on-that-day", "then, on that day"},
		{"laniso", "lâ-nî-sô", "lâ-nî-sô", "ADV", "", "WHEN", 2, "on-the-day-that", "on that day, on the day that"},
		{"lapara", "lapärä", "lapärä", "NOUN", "", "OBJ", 3, "airplane", "airplane"},
		{"laposo", "lâ-pôso", "lâ-pɔ̂sɔ", "NOUN", "", "WHEN", 2, "Saturday", "[lit: day|[Fr: portion]: ration] Saturday"},
		{"laso", "lâ-sô", "lâ-sô", "ADV", "", "WHEN", 2, "today", "today"},
		{"lavu", "lavu", "lavu", "NOUN", "", "ANIM", 6, "bee", "bee"},
		{"lawa", "lâ-wa", "lâ-wa", "ADV", "", "WHEN", 2, "when", "[lit: day|which]: when"},
		{"lawu", "lawü", "lawü", "NOUN", "", "OBJ", 4, "cuttingboard", "cutting/crushing board"},
		{"layenga", "lâ-yenga", "lâ-yenga", "NOUN", "", "WHEN", 2, "Sunday", "[lit: day|feast]: Sunday"},
		{"le", "lê", "lê", "NOUN", "", "STATE", 3, "seeds", "sprout, fruit, grain, seeds"},
		{"le", "lê", "lɛ̂", "NOUN", "", "BODY", 2, "eye-face-surface", "eye; face, surface; front, before one's eyes; blade"},
		{"le", "lë", "lë", "VERB", "Subcat=Intr", "STATE", 3, "bear-fruit", "sprout, bear fruit"},
		{"le", "lë", "lɛ̈", "NOUN", "", "ALT WORD FOR", 6, "mange", "särä"},
		{"lege", "lêgë", "lêgë", "NOUN", "", "HOW", 1, "road-times-way", "num times; way, manner, means, how to; way, path, road"},
		{"legeoko", "lêgë-ôko", "lêgë-ɔ̂kɔ", "ADV", "", "HOW", 1, "together", "once; the same way, similarly, identical; together, at the same time"},
		{"leke", "leke", "lɛkɛ", "VERB", "Subcat=Tran", "ACT", 1, "fix", "fix, repair, put in order, resolve; prepare"},
		{"lekere", "lekere", "lɛkɛrɛ", "VERB", "Subcat=Tran", "ACT", 1, "work-on", "work on, edit, produce"},
		{"lekpa", "lekpa", "lekpa", "NOUN", "", "ANIM", 6, "antilope", "bushbuck antilope"},
		{"lele", "lele", "lele", "NOUN", "", "ANIM", 6, "porcupine", "brushtailed porcupine"},
		{"lele", "lele", "lɛlɛ", "NOUN", "", "NATURE", 3, "pond", "pond, lake"},
		{"lele", "lêlê", "lêlê", "NOUN", "", "PLANT", 5, "beans", "beans"},
		{"lele", "lëlë", "lëlë", "NOUN", "", "ANIM", 4, "donkey", "donkey"},
		{"lele", "lëlë", "lɛ̈lɛ̈", "NOUN", "", "ALT WORD FOR", 6, "mange", "särä"},
		{"lembe", "lembe", "lembe", "NOUN", "", "FISH", 6, "catfish", "schilbid catfish"},
		{"lenda", "lëndâ", "lëndâ", "NOUN", "", "BODY", 5, "clitoris", "clitoris"},
		{"lende", "lendë", "lendë", "NOUN", "", "NATURE", 6, "lake", "lake"},
		{"lengbetoro", "lêngbêtôrô", "lêngbêtɔ̂rɔ̂", "NOUN", "", "FOOD", 6, "soybean", "[lit: seed|October]: soybean"},
		{"lenge", "lenge", "lenge", "NOUN", "", "NATURE", 5, "pearl", "[lit: seed|water?]: pearl"},
		{"lengua", "lêngua", "lɛ̂ngua", "NOUN", "", "WHEN", 5, "July", "July"},
		{"letibekpa", "lê-tî-bëkpä", "lɛ̂-tî-bëkpä", "NOUN", "", "NATURE", 3, "lightning", "[lit: eye|of|thunder]: lightning bolt"},
		{"letimbeti", "lê-tî-mbëtï", "lɛ̂-tî-mbɛ̈tï", "NOUN", "", "OBJ", 5, "page", "[lit: face|of|book] page"},
		{"letindo", "lê-tî-ndö", "lɛ̂-tî-ndö", "NOUN", "", "OBJ", 5, "site", "[lit: face|of|surface] site"},
		{"leyaka", "lê-yäkä", "lê-yäkä", "NOUN", "", "STATE", 4, "harvest", "harvest"},
		{"li", "li", "li", "NOUN", "", "BODY", 1, "head", "head"},
		{"li", "li", "li", "NOUN", "", "STATE", 3, "mildew", "mildew"},
		{"li", "li", "li", "NOUN", "", "WHEN", 1, "beginning", "beginning"},
		{"li", "li", "li", "NOUN", "", "WHERE", 1, "top", "top, front of a line, summit, point"},
		{"li", "li", "li", "VERB", "", "NUM", 2, "number", "number, count, quantity"},
		{"li", "lï", "lï", "VERB", "Subcat=Intr", "MOVE", 1, "enter", "enter, break into"},
		{"li", "lï", "lï", "VERB", "Subcat=Intr", "STATE", 3, "be-deep", "be deep"},
		{"lia", "lîâ", "lîâ", "NOUN", "", "OBJ", 4, "fishing-net", "handheld fishing net"},
		{"lifilo", "lîfïlo", "lîfïlo", "NOUN", "", "GOD", 4, "hell", "[Fr: l'enfer]: hell"},
		{"likisi", "likisi", "likisi", "NOUN", "", "INTERACT", 4, "fraud", "fraud, swindle"},
		{"likongo", "likongô", "likɔngɔ̂", "NOUN", "", "OBJ", 5, "javelin", "lance, javelin"},
		{"likundu", "likundû", "likundû", "NOUN", "", "STATE", 4, "sorcery", "magic, sorcery, evil spirit"},
		{"likune", "li-kûne", "li-kûnɛ", "NOUN", "", "STATE", 5, "automatic", "automatic"},
		{"linda", "linda", "linda", "NOUN", "", "STATE", 3, "entrance", "entrance, submergence"},
		{"linda", "linda", "linda", "VERB", "Subcat=Intr", "STATE", 3, "enter", "enter, be submerged"},
		{"lindo", "li-ndö", "li-ndö", "NOUN", "", "OBJ", 5, "address", "[lit: head|of|place] address"},
		{"lindotisinga", "li-ndö-tî-sînga", "li-ndö-tî-sînga", "NOUN", "", "OBJ", 5, "email-address", "[lit: head|of|surface|of|wire] email address"},
		{"lindotitokua", "li-ndö-tî-tokua", "li-ndö-tî-tokua", "NOUN", "", "OBJ", 5, "mailing-address", "[lit: head|of|surface|of|message] mailing address"},
		{"linga", "lïngä", "lïngä", "NOUN", "", "OBJ", 5, "tambourine", "[African-made]: wooden tambourine"},
		{"lingbi", "lîngbi", "lîngbi", "VERB", "Aspect=Imp|Mood=Nec|Subcat=Intr", "HOW", 1, "must-or-may", "must, may, can, should"},
		{"lingbi", "lîngbi", "lîngbi", "VERB", "Aspect=Imp|Mood=Pot|Subcat=Intr", "HOW", 3, "according", "suffice, be equal, according"},
		{"lingo", "lïngö", "lïngɔ̈", "VERB", "VerbForm=Vnoun", "MOVE", 1, "entrance", "entrance, entering, breaking in"},
		{"lingu", "li-ngû", "li-ngû", "NOUN", "", "NATURE", 3, "water-source", "[lit: head|water]: source of drinking water"},
		{"lio", "lîo", "lîɔ", "ADJ", "", "HOW", 5, "dwarf", "dwarf"},
		{"lisoro", "lisoro", "lisoro", "NOUN", "", "INTERACT", 2, "conversation", "conversation, chat"},
		{"litene", "li-tënë", "li-tɛ̈nɛ̈", "NOUN", "", "INTERACT", 5, "chapter", "[lit: head|speech]: chapter"},
		{"liti", "li-tï", "li-tï", "NOUN", "", "BODY", 3, "finger", "finger"},
		{"lititurungu", "li tî tûrûngu", "li tî tûrûngu", "NOUN", "", "BODY", 5, "umbilical-cord", "[lit: head|of|navel]: umbilical cord"},
		{"lo", "lo", "lo", "PRON", "Num=Sing|Person=3|PronType=Prs", "WHO", 1, "he-she-it", "he, she, it"},
		{"lo", "lö", "lö", "NOUN", "", "INTERACT", 4, "phrase", "phrase, pronouncement"},
		{"lobia", "lö-bîâ", "lö-bîâ", "NOUN", "", "INTERACT", 4, "poem", "[lit: phrase|song]: ode, poem, chant"},
		{"logbia", "lö-gbïä", "lö-gbïä", "NOUN", "", "GOD", 4, "gospel", "[lit: phrase|lord]: gospel"},
		{"lokpoto", "lokpoto", "lɔkpɔtɔ", "NOUN", "", "DRINK", 4, "malt", "malt, fermented dregs"},
		{"lokutu", "lö-kûtu", "lö-kûtu", "NOUN", "", "INTERACT", 4, "problem", "[lit: phrase|knot]: problem"},
		{"lolo", "lolo", "lolo", "NOUN", "", "FISH", 6, "tilapia", "carp, tilapia"},
		{"lombo", "lömbö", "lɔ̈mbɔ̈", "NOUN", "", "ANIM", 6, "frog", "frog"},
		{"lomveni", "lo mvenî", "lo mvɛnî", "PRON", "Num=Sing|Person=3|PronType=Prs", "WHO", 1, "himself-herself-itself", "[lit: (s)he,it|self]: himself, herself, itself"},
		{"londa", "lö-ndâ", "lö-ndâ", "NOUN", "", "INTERACT", 4, "formula", "[lit: phrase|end]: formula"},
		{"londo", "löndö", "löndö", "VERB", "Subcat=Intr", "MOVE", 1, "stand-up-or-leave", "depart, leave, wander off; start, begin; start standing up, rise"},
		{"londona", "löndö na", "löndö na", "VERB", "Subcat=Tran", "STATE", 1, "come-from", "come from, be from"},
		{"longo", "longo", "lɔngɔ", "NOUN", "", "ANIM", 6, "cobra", "cobra"},
		{"loro", "lörö", "lɔ̈rɔ̈", "NOUN", "", "MOVE", 3, "run", "run, race"},
		{"loso", "lôso", "lɔ̂sɔ", "NOUN", "", "FOOD", 3, "rice", "rice"},
		{"lu", "lü", "lü", "VERB", "Subcat=Tran", "ACT", 3, "bury", "bury"},
		{"lungula", "lungûla", "lungûla", "VERB", "Subcat=Tran", "ACT", 3, "remove", "remove, take away, take off, open"},
		{"lupa", "lûpa", "lûpa", "NOUN", "", "OBJ", 5, "ladle", "ladle"},
		{"luti", "lütï", "lütï", "VERB", "Subcat=Intr", "STATE", 1, "be-standing", "be standing, stop moving"},
		{"luu", "lûu", "lûu", "VERB", "Subcat=Intr", "STATE", 3, "be-worn-out", "be worn out"},
		{"luu", "lûu", "lûu", "VERB", "Subcat=Tran", "STATE", 3, "miss-or-fail", "miss, fail"},
		{"ma", "ma", "ma", "PART", "Mood=Opt", "INTERACT", 1, "[emphasis]", "[insistance]"},
		{"ma", "mâ", "mâ", "NOUN", "", "ALT SP FOR", 9, "ear", "mɛ̂"},
		{"ma", "mä", "mä", "VERB", "", "SENSE", 1, "hear-listen-smell-understand", "hear, listen, smell, understand"},
		{"mabaya", "mabaya", "mabaya", "NOUN", "", "OBJ", 6, "board", "board to sit or cut on"},
		{"mabe", "ma bê", "ma bɛ̂", "VERB", "Subcat=Intr", "FEEL", 1, "believe", "[lit: hear|heart]: believe, have faith"},
		{"mabe", "ma-bê", "ma-bɛ̂", "NOUN", "Subcat=Intr", "FEEL", 1, "belief", "[lit: hear|heart]: faith"},
		{"maboko", "mabôko", "mabɔ̂kɔ", "NOUN", "", "BODY", 2, "hand-arm-shoulder", "hand, arm, shoulder"},
		{"mafuta", "mafüta", "mafüta", "NOUN", "", "FOOD", 2, "oil", "oil, fat, grease"},
		{"mafutatimburu", "mafüta tî mbûrü", "mafüta tî mbûrü", "NOUN", "", "FOOD", 5, "palm-oil", "palm oil"},
		{"mafutatingutimetibagara", "mafüta tî ngû tî me tî bâgara", "mafüta tî ngû tî mɛ tî bâgara", "NOUN", "", "NATURE", 5, "butter", "[lit: oil|of|water|of|teat|of|cow]: butter"},
		{"mafutatiwotoro", "mafüta tî wôtoro", "mafüta tî wɔ̂tɔrɔ", "NOUN", "", "ANIM", 3, "honey", "[lit: oil|of|bee]: honey"},
		{"magbonga", "magböngä", "magbɔ̈ngä", "NOUN", "", "PLANT", 6, "banana-stalk", "banana stalk"},
		{"magia", "magia", "magia", "NOUN", "", "OBJ", 6, "throwing-knife", "throwing knife"},
		{"makako", "makâko", "makâko", "NOUN", "", "ANIM", 2, "monkey", "monkey"},
		{"makala", "makala", "makala", "NOUN", "", "FOOD", 5, "donut", "donut"},
		{"makango", "makângo", "makângo", "NOUN", "", "FAMILY", 5, "concubine", "mistress, lover, concubine"},
		{"makela", "makelâ", "makelâ", "NOUN", "", "STATE", 3, "good-luck", "good luck"},
		{"makobe", "makobe", "makobe", "NOUN", "", "WHEN", 6, "December", "December"},
		{"makongo", "makongö", "makongö", "NOUN", "", "ANIM", 4, "caterpillar", "caterpillar"},
		{"makoroo", "makoröo", "makoröo", "NOUN", "", "WHO", 5, "traitor-hypocrite", "traitor, hypocrite"},
		{"malangi", "malangi", "malangi", "NOUN", "", "OBJ", 6, "bottle", "bottle, flask, vase"},
		{"malinga", "malînga", "malînga", "NOUN", "", "ACT", 6, "dance", "modern style dance"},
		{"mama", "mamâ", "mamâ", "NOUN", "Gender=Fem", "FAMILY", 2, "mother", "mother"},
		{"mamatimapa", "mamâ tî mâpa", "mamâ tî mâpa", "NOUN", "", "FOOD", 5, "yeast", "bread yeast"},
		{"mamiwata", "mamîwätä", "mamîwätä", "NOUN", "", "GOD", 5, "water-nymph", "[En: mommy water]: albino water nymph blamed for drowning accidents"},
		{"manabe", "mä na bê", "mä na bɛ̂", "VERB", "Subcat=Tran", "FEEL", 1, "believe-in", "believe in, have faith in"},
		{"manabe", "mä-na-bê", "mä-na-bɛ̂", "NOUN", "", "GOD", 1, "Protestant", "[lit: hear|in|heart]: Protestant, Evangelical"},
		{"manda", "manda", "manda", "VERB", "Subcat=Tran", "ACT", 1, "learn", "learn, study, imitate"},
		{"manda", "mä ndâ", "mä ndâ", "VERB", "Subcat=Intr", "SENSE", 1, "understand", "understand how and why"},
		{"mandako", "mandako", "mandako", "NOUN", "", "INTERACT", 6, "canoe-race", "canoe race"},
		{"mando", "mä-ndo", "mä-ndo", "INTERJ", "Mood=Jus", "INTERACT", 4, "Listen-up", "Listen up!"},
		{"manga", "mânga", "mânga", "NOUN", "", "OBJ", 3, "cigarette-or-tobacco", "cigarette, tobacco"},
		{"mangbere", "mangbêrê", "mangbɛ̂rɛ̂", "NOUN", "", "FOOD", 3, "manioc-bar", "manioc bar"},
		{"mangbi", "mângbi", "mângbi", "NOUN", "", "SENSE", 1, "agreement", "agreement"},
		{"mangbi", "mângbi", "mângbi", "VERB", "Aspect=Imp|Subcat=Intr", "SENSE", 1, "agree", "agree"},
		{"mangboko", "mangbökö", "mangbökö", "NOUN", "", "OBJ", 3, "ship", "ship"},
		{"mango", "mângo", "mângo", "NOUN", "", "FOOD", 3, "mango", "mango"},
		{"manzeke", "manzêke", "manzêke", "NOUN", "", "FOOD", 6, "black-pepper", "black pepper"},
		{"manzinzi", "manzinzi", "manzinzi", "NOUN", "", "FOOD", 6, "black-pepper", "black pepper"},
		{"mapa", "mâpa", "mâpa", "NOUN", "", "FOOD", 3, "bread", "bread"},
		{"mapia", "mapîâ", "mapîâ", "NOUN", "", "OBJ", 6, "pagne", "pagne (wrap worn by woman)"},
		{"mapo", "mapô", "mapô", "NOUN", "", "OBJ", 4, "scissors", "scissors"},
		{"mara", "marä", "marä", "NOUN", "", "NUM", 1, "tribe-or-type", "tribe, race, type, kind, sort, variety"},
		{"mara", "märä", "märä", "ADJ", "", "SICK", 5, "barren", "sterile, barren"},
		{"masango", "masango", "masango", "NOUN", "", "OBJ", 6, "newsletter", "newsletter, bulletin"},
		{"masaragba", "mâsarâgba", "mâsarâgba", "NOUN", "", "ANIM", 4, "rhinoceros", "rhinoceros"},
		{"maseka", "maseka", "maseka", "NOUN", "", "FAMILY", 3, "youth", "[12-16yrs]: youth"},
		{"masekatikoli", "maseka tî kôlï", "maseka tî kɔ̂lï", "NOUN", "Gender=Masc", "FAMILY", 3, "young-man", "[12-16yrs]: young man"},
		{"masekatiwali", "maseka tî wâlï", "maseka tî wâlï", "NOUN", "Gender=Fem", "FAMILY", 3, "young-woman", "[12-16yrs]: young woman"},
		{"masini", "masïni", "masïni", "NOUN", "", "OBJ", 5, "machine", "machine"},
		{"masua", "masua", "masua", "NOUN", "", "OBJ", 3, "boat", "boat"},
		{"matabisi", "matabïsi", "matabïsi", "NOUN", "", "INTERACT", 3, "advantage-or-gift", "advantage, tip, gratuity, gift"},
		{"matanga", "matânga", "matânga", "NOUN", "", "CIVIL", 2, "feast-or-holiday", "feast, banquet, ceremony, national holiday"},
		{"mawa", "mawa", "mawa", "NOUN", "", "CIVIL", 2, "misery", "misery, suffering, unhappiness; pity, compassion"},
		{"mawoya", "mawôya", "mawôya", "NOUN", "", "SICK", 6, "epidemic", "epidemic"},
		{"mayanga", "mä-yângâ", "mä-yângâ", "NOUN", "", "INTERACT", 3, "obey", "obey"},
		{"mayere", "mayëre", "mayɛ̈rɛ", "NOUN", "", "HOW", 3, "means-or-ability", "manner, customs, habits; skill, means, ability"},
		{"mba", "mbâ", "mbâ", "NOUN", "", "WHO", 2, "compatriot", "comrade, compatriot, fellow citizen"},
		{"mbadi", "mbadi", "mbadi", "NOUN", "", "GOD", 4, "divination", "divination, fortune telling"},
		{"mbage", "mbâgë", "mbâgë", "NOUN", "", "WHERE", 1, "side-or-direction", "side, direction"},
		{"mbagetikoli", "mbâgë tî kôlï", "mbâgë tî kɔ̂lï", "NOUN", "", "WHERE", 1, "right-side", "[lit: side|of|man]: right side"},
		{"mbagetiwali", "mbâgë tî wâlï", "mbâgë tî wâlï", "NOUN", "", "WHERE", 1, "left-side", "[lit: side|of|woman]: left side"},
		{"mbai", "mbai", "mbai", "NOUN", "", "INTERACT", 6, "proverb", "proverb"},
		{"mbakele", "mbâkêlê", "mbâkêlê", "NOUN", "", "FOOD", 6, "yellow-squash", "yellow squash"},
		{"mbakoro", "mbäkôro", "mbäkɔ̂rɔ", "ADJ", "", "HOW", 4, "old", "old, aged"},
		{"mbakoro", "mbäkôro", "mbäkɔ̂rɔ", "NOUN", "", "WHO", 4, "old-person", "old person"},
		{"mbala", "mbala", "mbala", "NOUN", "", "ANIM", 6, "elephant", "elephant"},
		{"mbamba", "mbamba", "mbamba", "NOUN", "", "ANIM", 5, "oyster", "oyster, mussel"},
		{"mbamba", "mbamba", "mbamba", "NOUN", "", "NATURE", 4, "chalk", "chalk, whitewash"},
		{"mbana", "mbänä", "mbänä", "NOUN", "", "FEEL", 4, "wickedness", "malice, wickedness"},
		{"mbangu", "mbängü", "mbängü", "NOUN", "", "WHEN", 5, "March", "March"},
		{"mbanu", "mbanu", "mbanu", "NOUN", "", "OBJ", 6, "cross-bow", "cross bow"},
		{"mbarambara", "mbârâmbârâ", "mbârâmbârâ", "ADJ", "NumType=Ord", "NUM", 2, "seven", "seven"},
		{"mbarambara", "mbârâmbârâ", "mbârâmbârâ", "NUM", "NumType=Card", "NUM", 2, "seven", "seven"},
		{"mbarata", "mbârâtâ", "mbârâtâ", "NOUN", "", "ANIM", 3, "horse", "horse"},
		{"mbarawara", "mbârâwârâ", "mbârâwârâ", "NOUN", "", "ANIM", 6, "iguana", "iguana"},
		{"mbasa", "mbâsa", "mbâsa", "NOUN", "", "FOOD", 5, "metallic-lead", "[metal] lead"},
		{"mbasala", "mbâsala", "mbâsala", "NOUN", "", "FOOD", 6, "green-onion", "green onion"},
		{"mbasambara", "mbâsâmbârâ", "mbâsâmbârâ", "ADJ", "NumType=Ord", "ALT SP FOR", 9, "seven", "mbârâmbârâ"},
		{"mbasambara", "mbâsâmbârâ", "mbâsâmbârâ", "NUM", "NumType=Card", "ALT SP FOR", 9, "seven", "mbârâmbârâ"},
		{"mbata", "mbata", "mbata", "NOUN", "", "ANIM", 6, "lion", "lion"},
		{"mbata", "mbätä", "mbätä", "NOUN", "", "HOUSE", 3, "chair", "bench, chair, stool"},
		{"mbea", "mbêâ", "mbɛ̂â", "NOUN", "", "WHERE", 5, "opposite-bank", "opposite bank"},
		{"mbenge", "mbëngë", "mbɛ̈ngɛ̈", "NOUN", "", "ALT SP FOR", 9, "poison", "bɛ̈ngɛ̈"},
		{"mbenge", "mbëngë", "mbɛ̈ngɛ̈", "NOUN", "", "ANIM", 6, "bush-pig", "bush pig"},
		{"mbeni", "mbênî", "mbɛ̂nî", "ADJ", "", "NUM", 1, "some", "another, some...other..."},
		{"mbeni", "mbênî", "mbɛ̂nî", "ADV", "", "NUM", 1, "again", "again, still, once more"},
		{"mbeni", "mbênî", "mbɛ̂nî", "NOUN", "", "NUM", 1, "another-one", "another one, some...others..."},
		{"mbeniape", "mbênî äpe", "mbɛ̂nî äpɛ", "ADV", "", "NUM", 1, "no-more", "no more, no longer"},
		{"mbenibiri", "mbênî-bîrï", "mbɛ̂nî-bîrï", "ADV", "", "WHEN", 2, "day-before-yesterday", "day before yesterday"},
		{"mbenikekereke", "mbênî-kêkerêke", "mbɛ̂nî-kêkerêke", "ADV", "", "WHEN", 2, "day-after-tomorrow", "day after tomorrow"},
		{"mbenila", "mbênî-lâ", "mbɛ̂nî-lâ", "ADV", "", "WHEN", 2, "someday", "someday, another day"},
		{"mbenipepe", "mbênî pëpe", "mbɛ̂nî pɛ̈pɛ", "ADV", "", "NUM", 1, "no-more", "no more, no longer"},
		{"mbere", "mbere", "mbere", "NOUN", "", "ANIM", 6, "blood-pact-or-trapdoor-spider", "alliance, blood pact; trapdoor spider"},
		{"mbereke", "mbêrêkê", "mbɛ̂rɛ̂kɛ̂", "NOUN", "", "FOOD", 6, "watermelon", "watermelon"},
		{"mbeso", "mbeso", "mbeso", "ADV", "", "WHEN", 1, "formerly", "formerly, lately, in times past"},
		{"mbeti", "mbëtï", "mbɛ̈tï", "NOUN", "", "OBJ", 2, "writing", "paper, letter, writing, document, receipt"},
		{"mbetikua", "mbëtï-kua", "mbɛ̈tï-kua", "NOUN", "", "OBJ", 2, "work-permit", "work permit"},
		{"mbetilege", "mbëtï-lêgë", "mbɛ̈tï-lêgë", "NOUN", "", "OBJ", 2, "travel-documents", "travel documents"},
		{"mbetisango", "mbëtï-sango", "mbɛ̈tï-sango", "NOUN", "", "OBJ", 2, "newspaper-or-magazine", "newspaper, magazine"},
		{"mbetitinzapa", "mbëtï tî nzapä", "mbɛ̈tï tî nzapä", "NOUN", "", "GOD", 2, "bible-or-catechism", "bible, catechism"},
		{"mbetitokua", "mbëtï-tokua", "mbɛ̈tï-tokua", "NOUN", "", "OBJ", 2, "message", "message, missive"},
		{"mbeto", "mbeto", "mbɛtɔ", "NOUN", "", "FEEL", 3, "fear", "fear"},
		{"mbi", "mbï", "mbï", "PRON", "Num=Sing|Person=1|PronType=Prs", "WHO", 1, "I-me", "I, me"},
		{"mbimveni", "mbï mvenî", "mbï mvɛnî", "PRON", "Num=Sing|Person=1|PronType=Prs", "WHO", 1, "myself", "[lit: I,me|self]: myself"},
		{"mbinda", "mbîndä", "mbîndä", "NOUN", "", "NATURE", 3, "cloud-fog-mist", "cloud, fog, mist"},
		{"mbingo", "mbîngo", "mbîngo", "NOUN", "", "NATURE", 3, "darkness-ignorance", "darkness, ignorance"},
		{"mbio", "mbïö", "mbïɔ̈", "NOUN", "", "TREE", 6, "padauk-tree", "African padauk tree"},
		{"mbirimbiri", "mbîrîmbîrî", "mbîrîmbîrî", "ADJ", "", "HOW", 1, "straight-or-honest", "straight, just, loyal, honest, moral"},
		{"mbirimbiri", "mbîrîmbîrî", "mbîrîmbîrî", "ADV", "", "HOW", 1, "honestly", "honestly, correctly, perfectly"},
		{"mbo", "mbo", "mbo", "NOUN", "", "ANIM", 2, "dog", "dog"},
		{"mbo", "mbô", "mbô", "VERB", "Subcat=Tran", "ACT", 3, "wipe-erase-clean", "wipe, erase, clean"},
		{"mbo", "mbö", "mbö", "NOUN", "", "BODY", 3, "breath", "breath"},
		{"mboko", "mbôko", "mbôko", "VERB", "", "BODY", 4, "be-bald-or-bruised-or-shed-skin", "be bald; become bruised, become scratched; shed"},
		{"mbokoli", "mbôko-li", "mbôko-li", "ADJ", "", "BODY", 4, "bald", "bald"},
		{"mbokoro", "mbökôro", "mbɔ̈kɔ̂rɔ", "ADJ", "", "HOW", 4, "old", "old, aged"},
		{"mbombo", "mbômbô", "mbômbô", "ADJ", "", "BODY", 4, "bald", "bald"},
		{"mbomboli", "mbômbô-li", "mbômbô-li", "NOUN", "", "BODY", 4, "baldness", "baldness"},
		{"mbongo", "mbongo", "mbongo", "NOUN", "", "WHERE", 2, "south", "south; left riverbank"},
		{"mbongo", "mbôngo", "mbɔ̂ngɔ", "NOUN", "", "CIVIL", 6, "money", "money"},
		{"mbongo", "mböngö", "mböngɔ̈", "VERB", "VerbForm=Vnoun", "WHERE", 3, "joking", "[lit: breathing]: joking"},
		{"mboro", "mborô", "mborô", "NOUN", "", "BODY", 6, "acne", "pimple, acne"},
		{"mbororo", "mbôrôrô", "mbôrôrô", "NOUN", "", "WHO", 5, "Fulani", "Fulani, Muslim cow herders"},
		{"mboto", "mbotö", "mbɔtɔ̈", "NOUN", "", "FISH", 6, "baby-catfish", "baby catfish"},
		{"mbuki", "mbûki", "mbûki", "NOUN", "", "CIVIL", 3, "alliance", "pact, alliance"},
		{"mbuma", "mbuma", "mbuma", "NOUN", "", "FOOD", 6, "almond-or-kernel", "almond, kernel"},
		{"mburu", "mburu", "mburu", "NOUN", "", "NATURE", 2, "powder", "powder"},
		{"mburu", "mbûrü", "mbûrü", "NOUN", "", "TREE", 5, "oil-palm", "oil palm"},
		{"mburutiwa", "mburu tî wâ", "mburu tî wâ", "NOUN", "", "COLOR", 3, "gray", "[lit: cinders]: gray"},
		{"mburutiwa", "mburu tî wâ", "mburu tî wâ", "NOUN", "", "NATURE", 2, "cinders", "cinders"},
		{"mbutu", "mbütü", "mbütü", "NOUN", "", "NATURE", 3, "sand", "sand"},
		{"me", "me", "mɛ", "NOUN", "", "BODY", 2, "breast", "breast"},
		{"me", "me", "mɛ", "VERB", "Subcat=Intr", "MOVE", 2, "climb-ascend", "climb, ascend"},
		{"me", "mê", "mɛ̂", "NOUN", "", "BODY", 2, "ear", "ear"},
		{"me", "mê", "mɛ̂", "VERB", "Subcat=Tran", "ACT", 4, "knead", "knead"},
		{"mea", "mêa", "mêa", "NOUN", "", "WHO", 6, "twin", "twin"},
		{"mee", "meë", "mɛɛ̈", "CCONJ", "", "HOW", 1, "but", "but"},
		{"meka", "meka", "mɛka", "NOUN", "", "CIVIL", 4, "limit-border", "limit, border"},
		{"meka", "meka", "mɛka", "VERB", "Subcat=Tran", "INTERACT", 4, "measure-compete", "size up, measure, compete"},
		{"mene", "mene", "mɛnɛ", "VERB", "Subcat=Tran", "ACT", 3, "swallow", "swallow"},
		{"mene", "mênë", "mɛ̂nɛ̈", "NOUN", "", "BODY", 3, "blood", "blood, one related by blood"},
		{"menga", "mëngä", "mɛ̈ngä", "NOUN", "", "BODY", 3, "tongue", "tongue"},
		{"mengo", "mëngö", "mɛ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "MOVE", 2, "climbing", "climbing"},
		{"mesa", "mêsa", "mɛ̂sa", "NOUN", "", "GOD", 5, "mass", "mass"},
		{"meti", "metï", "mɛtï", "VERB", "Subcat=Intr", "MOVE", 2, "climb", "climb vertically"},
		{"mi", "mî", "mî", "NOUN", "", "BODY", 3, "flesh-thickness", "flesh; thickness"},
		{"mimi", "mîmi", "mîmi", "ADJ", "", "CIVIL", 5, "brave", "brave, courageous"},
		{"mingi", "mîngi", "mîngi", "ADV", "", "NUM", 1, "very", "very, too[ much]"},
		{"mingo", "mîngo", "mîngɔ", "VERB", "Subcat=Tran", "ACT", 3, "extinguish", "extinguish, close"},
		{"meambe", "meambe", "meambe", "ADJ", "NumType=Ord", "NUM", 2, "eight", "eight"},
		{"meambe", "meambe", "meambe", "NUM", "NumType=Card", "NUM", 2, "eight", "eight"},
		{"miombe", "miombe", "miɔmbɛ", "ADJ", "NumType=Ord", "ALT SP FOR", 9, "eight", "meambe"},
		{"miombe", "miombe", "miɔmbɛ", "NUM", "NumType=Card", "ALT SP FOR", 9, "eight", "meambe"},
		{"misuiya", "mî-suïya", "mî-suïya", "NOUN", "", "FOOD", 3, "skewered-meat", "[lit: flesh|skewer]: skewered meat,kabob"},
		{"mitere", "mî-terê", "mî-tɛrɛ̂", "NOUN", "", "BODY", 3, "muscle", "muscle"},
		{"mo", "mo", "mɔ", "PRON", "Num=Sing|Person=2|PronType=Prs", "WHO", 1, "you", "[singular]: you"},
		{"mobaamotene", "mo-bâa-mo-tene", "mɔ-bâa-mɔ-tɛnɛ", "SCONJ", "", "HOW", 1, "as-if", "[lit: (if) you|see (it)|(then) you (would) say]: as if"},
		{"modogere", "modögerê", "mɔdɔ̈gɛrɛ̂", "NOUN", "", "FOOD", 6, "soybean", "soybean"},
		{"mokiri", "mokiri", "mɔkiri", "NOUN", "", "CIVIL", 6, "world", "world, universe, history"},
		{"mokondo", "mokondö", "mɔkondö", "ADJ", "", "GOD", 4, "holy", "holy, saintly, honest, just, pure"},
		{"mokondo", "mokondö", "mɔkondö", "NOUN", "", "GOD", 4, "sainthood", "sainthood"},
		{"mokonzi", "mokönzi", "mokönzi", "NOUN", "", "CIVIL", 3, "village-chief", "village chief"},
		{"molenge", "môlengê", "môlɛngɛ̂", "NOUN", "", "FAMILY", 3, "child", "[5-11yrs]: child; (man's) grandchild"},
		{"molengetikoli", "môlengê tî kôlï", "môlɛngɛ̂ tî kɔ̂lï", "NOUN", "Gender=Masc", "FAMILY", 3, "boy", "[5-11yrs]: boy"},
		{"molengetindeko", "môlengê tî ndeko", "môlɛngɛ̂ tî ndeko", "NOUN", "", "FAMILY", 5, "bastard", "[lit: child|of|adultery]: bastard"},
		{"molengetiwali", "môlengê tî wâlï", "môlɛngɛ̂ tî wâlï", "NOUN", "Gender=Fem", "FAMILY", 3, "girl", "[5-11yrs]: girl"},
		{"molongo", "molongö", "molongö", "NOUN", "", "WHERE", 2, "row-or-column", "row, column, line, alignment"},
		{"momveni", "mo mvenî", "mɔ mvɛnî", "PRON", "Num=Sing|Person=2|PronType=Prs", "WHO", 1, "yourself", "[lit: you|self]: yourself [singular]"},
		{"mondelepako", "mondelepâko", "mondelepâko", "NOUN", "", "FOOD", 6, "sweet-manioc", "sweet manioc"},
		{"monganga", "mongânga", "mongânga", "NOUN", "", "WHO", 6, "shaman", "witchdoctor, shaman, sorceror"},
		{"mongoli", "möngö-li", "möngö-li", "NOUN", "", "BODY", 6, "brain", "brain"},
		{"mopi", "mopï", "mɔpï", "NOUN", "", "STATE", 3, "bad-luck", "bad luck"},
		{"mosongoli", "mosongôli", "mosongôli", "NOUN", "", "HOW", 6, "pink", "pink"},
		{"mosoro", "mosoro", "mɔsɔrɔ", "NOUN", "", "CIVIL", 4, "wealth", "[lit: you|choose]: wealth"},
		{"mosuma", "mosümä", "mosümä", "NOUN", "", "FEEL", 3, "dream", "[lit: you|dream]: dream"},
		{"moyetibaa", "mo-yê-tî-bâa", "mɔ-yê-tî-bâa", "SCONJ", "", "HOW", 1, "as-it-were", "[lit: you|want|to|see][Fr: si vous voulez]: as it were"},
		{"mozingo", "mozïngö", "mɔzïngɔ̈", "NOUN", "", "CIVIL", 6, "poverty", "[lit: you|rising up]: poverty"},
		{"mu", "mû", "mû", "VERB", "Subcat=Tran", "INTERACT", 1, "give-or-take", "give [+na=to]; take [+na=as]"},
		{"mua", "müä", "müä", "NOUN", "", "CIVIL", 3, "mourning", "common purpose; mourning, fasting, penitance; castration"},
		{"muen", "muen", "muen", "NOUN", "", "NATURE", 4, "tidal-ebb", "tidal ebb"},
		{"mukoli", "mû kôlï", "mû kɔ̂lï", "VERB", "Subcat=Intr", "CIVIL", 2, "marry-a-man", "marry a man, take a husband"},
		{"mulege", "mû lêgë", "mû lêgë", "VERB", "Subcat=Intr", "MOVE", 1, "set-off", "[lit: take|road]: start a journey, set off"},
		{"mulegena", "mû lêgë na", "mû lêgë na", "VERB", "Subcat=Tran", "INTERACT", 1, "authorize", "[lit: give|road|to]: authorize (someone) [+ti](do something)"},
		{"mulu", "mûlu", "mûlu", "NOUN", "", "HOUSE", 5, "brick-baking-pan", "[Fr: moule]: pan to make bricks"},
		{"mumabokona", "mû mabôko na", "mû mabɔ̂kɔ na", "VERB", "Subcat=Tran", "INTERACT", 2, "help", "[lit: give|hand|to]: help, lend a hand to"},
		{"munambi", "mû na mbï", "mû na mbï", "VERB", "Subcat=Tran", "INTERACT", 1, "give-me", "[lit: give|to|me]: I beg you for"},
		{"mungbi", "mûngbi", "mûngbi", "NOUN", "", "CIVIL", 2, "marriage", "marriage"},
		{"mungbi", "mûngbi", "mûngbi", "VERB", "Aspect=Imp|Subcat=Intr", "CIVIL", 2, "get-married", "get married"},
		{"mungianabeti", "mû ngîâ na bê tî", "mû ngîâ na bɛ̂ tî", "VERB", "Subcat=Tran", "FEEL", 1, "please", "[lit: give|joy|to|heart|of]: please"},
		{"munzu", "munzû", "munzû", "NOUN", "", "WHO", 1, "white-foreigner", "white person, foreigner, European"},
		{"munzunzapa", "munzû-nzapä", "munzû-nzapä", "NOUN", "", "GOD", 1, "missionary", "missionary"},
		{"munzuvuko", "munzû-vukö", "munzû-vukɔ̈", "NOUN", "", "WHO", 1, "functionary", "African acting like a foreigner, functionary"},
		{"muru", "mûrû", "mûrû", "NOUN", "", "ANIM", 6, "panther", "panther, leopard"},
		{"muwa", "mû wâ", "mû wâ", "VERB", "Subcat=Intr", "STATE", 3, "be-hot", "catch fire, be hot, have a fever"},
		{"muwali", "mû wâlï", "mû wâlï", "VERB", "Subcat=Intr", "CIVIL", 2, "marry-a-woman", "marry a woman, take a wife"},
		{"muwango", "mû wängö", "mû wängɔ̈", "VERB", "Subcat=Intr", "INTERACT", 3, "give-advice", "[lit: give|advising]: give advice"},
		{"muwangona", "mû wängö na", "mû wängɔ̈ na", "VERB", "Subcat=Tran", "INTERACT", 3, "advise", "[lit: give|advising|to]: advise"},
		{"muyangati", "mû yângâ tî", "mû yângâ tî", "VERB", "Subcat=Tran", "STATE", 1, "overwhelm-or-discourage", "[lit: take|mouth|of]: overwhelm, discourage"},
		{"mvele", "mvele", "mvɛlɛ", "NOUN", "", "NATURE", 5, "copper", "copper"},
		{"mvene", "mvene", "mvɛnɛ", "NOUN", "", "INTERACT", 1, "untruth", "lie, untruth"},
		{"mveni", "mvenî", "mvɛnî", "NOUN", "", "WHO", 2, "VIP", "proprietor, owner; self"},
		{"mvitasioon", "mvitasiöon", "mvitasiöon", "NOUN", "", "INTERACT", 7, "invitation", "[Fr: invitation]: invitation"},
		{"mvuka", "mvüka", "mvüka", "NOUN", "", "WHEN", 5, "September", "September"},
		{"na", "na", "na", "ADP", "", "HOW", 1, "at-in-with", "at, in, with"},
		{"na", "na", "na", "CCONJ", "", "HOW", 1, "and", "and (between nouns)"},
		{"nabanduru", "nabändurü", "nabändurü", "NOUN", "", "WHEN", 5, "November", "November"},
		{"nambagewa", "na mbâgë wa", "na mbâgë wa", "ADV", "", "WHICH", 1, "whither", "[lit: at|direction|which]: in which direction"},
		{"nandowa", "na ndo wa", "na ndo wa", "ADV", "", "WHICH", 1, "where", "[lit: at|place|which]: where"},
		{"nda", "ndâ", "ndâ", "NOUN", "", "WHERE", 1, "end-rear", "end, behind, rear, base, essence"},
		{"ndagere", "ndâ-gerê", "ndâ-gɛrɛ̂", "NOUN", "", "BODY", 2, "heel", "[lit: base|leg]: heel"},
		{"ndali", "ndâ-li", "ndâ-li", "NOUN", "", "WHY", 1, "intent-purpose-cause", "[lit: end|beginning]: intent, purpose, cause; nape of the neck"},
		{"ndaliti", "ndâ-li tî", "ndâ-li tî", "ADP", "", "HOW", 1, "in-order-for", "[lit: purpose|of]: in order for"},
		{"ndalitinye", "ndâ-li tî nye", "ndâ-li tî nyɛ", "ADV", "", "HOW", 1, "why", "[lit: purpose|of|what]: to what end, why"},
		{"ndalitiso", "ndâ-li tî sô", "ndâ-li tî sô", "CCONJ", "", "HOW", 1, "because", "[lit: reason|of|that]: in order to"},
		{"ndalo", "ndâ-lö", "ndâ-lö", "NOUN", "", "WHY", 1, "solution", "[lit: base|phrase]: solution"},
		{"ndambo", "ndâmbo", "ndâmbo", "NOUN", "", "NUM", 2, "half", "half, part, portion, leftovers"},
		{"ndangba", "ndângbâ", "ndângbâ", "NOUN", "", "FAMILY", 4, "youngest", "[lit: last]: lastborn child, younger"},
		{"ndangbaliti", "ndângbâ-li-tï", "ndângbâ-li-tï", "NOUN", "", "BODY", 4, "pinky-finger", "[lit: last|finger]: pinky"},
		{"ndao", "ndao", "ndao", "NOUN", "", "ACT", 6, "forge", "forge"},
		{"ndapere", "ndäpêrê", "ndäpêrê", "NOUN", "", "WHEN", 2, "morning", "morning"},
		{"ndaperere", "ndäpêrêrê", "ndäpêrêrê", "NOUN", "", "WHEN", 2, "morning", "morning"},
		{"ndara", "ndarä", "ndarä", "NOUN", "", "STATE", 2, "intelligence", "intelligence, wisdom"},
		{"ndara", "ndârâ", "ndârâ", "NOUN", "", "STATE", 4, "tension", "tension"},
		{"ndaramba", "ndaramba", "ndaramba", "NOUN", "", "ANIM", 4, "rabbit", "rabbit, hare"},
		{"ndatu", "ndatu", "ndatu", "NOUN", "", "WHEN", 4, "dawn", "dawn"},
		{"ndawo", "ndawo", "ndawo", "NOUN", "", "ACT", 6, "accident", "accident; forge"},
		{"nde", "ndê", "ndê", "ADV", "", "HOW", 1, "differently", "differently"},
		{"nde", "ndê", "ndê", "DET", "PronType=Ind", "HOW", 1, "different", "different"},
		{"ndeke", "ndeke", "ndɛkɛ", "NOUN", "", "ANIM", 2, "bird", "bird"},
		{"ndeko", "ndeko", "ndeko", "NOUN", "", "FAMILY", 2, "friend-or-friendship", "friend, friendship, adultery"},
		{"ndekozande", "ndeko-zändë", "ndeko-zändë", "NOUN", "", "FAMILY", 6, "homosexuality", "[lit: friendship|(of)|Zande (tribe)]: homosexuality"},
		{"ndembe", "ndembë", "ndembë", "NOUN", "", "WHEN", 2, "season", "(cultural) season, epoch, festivities"},
		{"ndembo", "ndembö", "ndembö", "NOUN", "", "OBJ", 3, "rubber-ball-soccer", "rubber, ball, soccer"},
		{"ndende", "ndê-ndê", "ndê-ndê", "ADV", "", "HOW", 1, "each-in-their-own-way", "each in their own way"},
		{"ndende", "ndê-ndê", "ndê-ndê", "DET", "PronType=Ind", "HOW", 1, "various", "various"},
		{"ndendia", "ndêndïä", "ndêndïä", "NOUN", "", "INTERACT", 6, "hypocrisy-treachery-fraud", "[lit: other|law]: hypocrisy, treachery, fraud"},
		{"ndiba", "ndiba", "ndiba", "NOUN", "", "SICK", 6, "yaws", "yaws (illness of hands or feet)"},
		{"ndika", "ndikâ", "ndikâ", "NOUN", "", "CIVIL", 4, "law-commandment-decree", "law, commandment, decree"},
		{"ndiri", "ndïrï", "ndïrï", "NOUN", "", "STATE", 2, "stupidity", "stupidity"},
		{"ndo", "ndo", "ndo", "NOUN", "", "WHERE", 1, "place", "place; weather, climate, atmosphere"},
		{"ndo", "ndô", "ndô", "NOUN", "", "NATURE", 6, "pottery", "pottery clay, porcelain"},
		{"ndo", "ndö", "ndö", "NOUN", "", "WHERE", 1, "on-above-surface", "on, above, surface"},
		{"ndoi", "ndoî", "ndoî", "NOUN", "", "FAMILY", 4, "homonym", "homonym; person with same name, friend"},
		{"ndoko", "ndökö", "ndɔ̈kɔ̈", "NOUN", "", "PLANT", 4, "flower", "flower"},
		{"ndokoro", "ndokôrö", "ndɔkɔ̂rɔ̈", "NOUN", "", "PLANT", 5, "sponge", "sponge"},
		{"ndole", "ndö-lê", "ndö-lɛ̂", "NOUN", "", "BODY", 6, "forehead-eyebrow", "[lit: above|eye]: forehead, eyebrow"},
		{"ndombe", "ndombe", "ndɔmbɛ", "NOUN", "", "WHO", 4, "shopkeeper", "shopkeeper"},
		{"ndombo", "ndombö", "ndɔmbɔ̈", "NOUN", "", "PLANT", 4, "rattan", "rattan"},
		{"ndongo", "ndôngô", "ndôngô", "NOUN", "", "FOOD", 4, "pepper", "cayenne pepper, bell pepper"},
		{"ndoni", "ndö nî", "ndö nî", "NOUN", "", "WHERE", 1, "above-all", "addition, above all"},
		{"ndoo", "ndôo", "ndôo", "VERB", "Subcat=Intr", "HOW", 2, "be-short", "be short"},
		{"ndoti", "ndö-tî", "ndö-tî", "NOUN", "", "BODY", 6, "shoulder", "[lit: above|hand]: shoulder"},
		{"ndowa", "ndo-wâ", "ndo-wâ", "NOUN", "", "WHERE", 3, "heat", "heat"},
		{"ndoye", "ndö-yê", "ndö-yê", "NOUN", "", "GOD", 3, "God's-love-for-man", "[lit: above|want]: love of God for man"},
		{"ndoye", "ndö-yê", "ndö-yê", "VERB", "Subcat=Tran", "GOD", 3, "God's-love-for-man", "[lit: above|want]: love of God for man"},
		{"ndoyengo", "ndö-yëngö", "ndö-yëngɔ̈", "NOUN", "", "GOD", 3, "man's-love-For-God", "[lit: above|wanting]: love of God by man"},
		{"ndoyengo", "ndö-yëngö", "ndö-yëngɔ̈", "VERB", "Subcat=Tran", "GOD", 3, "man's-love-For-God", "[lit: above|wanting]: love of God by man"},
		{"ndu", "ndû", "ndû", "VERB", "Subcat=Tran", "ACT", 3, "touch", "touch"},
		{"ndu", "ndü", "ndü", "NOUN", "", "CIVIL", 3, "widowed", "widowed, in mourning"},
		{"ndumba", "ndûmba", "ndûmba", "NOUN", "", "WHO", 4, "self-supporting-unmarried-woman", "self-supporting unmarried woman"},
		{"nduru", "ndurü", "ndurü", "ADJ", "", "HOW", 2, "short-near", "short, near"},
		{"ndurukpa", "ndûrûkpâ", "ndûrûkpâ", "ADJ", "", "WHO", 5, "dwarf", "dwarf"},
		{"ndutu", "ndütü", "ndütü", "NOUN", "", "OBJ", 5, "basket-or-jug", "grain basket, wine jug (5-10 liters)"},
		{"nduzu", "ndüzü", "ndüzü", "NOUN", "", "NATURE", 1, "sky", "sky, upwards"},
		{"ne", "ne", "nɛ", "VERB", "Subcat=Intr", "STATE", 2, "weigh", "weigh, be heavy"},
		{"neka", "neka", "nɛka", "VERB", "Subcat=Tran", "ACT", 2, "grind", "grind, crush, mill"},
		{"nengo", "nëngö", "nɛ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "STATE", 2, "weight", "weight, thickness, heaviness"},
		{"nga", "ngâ", "ngâ", "ADV", "", "NUM", 1, "also", "also"},
		{"nga", "ngä", "ngä", "NOUN", "", "OBJ", 1, "bottle", "bottle"},
		{"ngaakoo", "ngâakôo", "ngâakôo", "NOUN", "", "PLANT", 5, "sugarcane", "sugarcane"},
		{"ngago", "ngâgö", "ngâgö", "NOUN", "", "PLANT", 4, "eggplant", "eggplant"},
		{"ngambi", "ngambi", "ngambi", "NOUN", "", "FAMILY", 4, "younger-sibling", "younger sibling"},
		{"ngan", "ngän", "ngän", "NOUN", "", "PLANT", 6, "sugarcane", "sugarcane"},
		{"nganga", "nganga", "nganga", "NOUN", "", "CIVIL", 4, "charm", "charm, talisman, magic; medicine"},
		{"nganga", "ngängä", "ngängä", "NOUN", "", "PLANT", 4, "gourd", "gourd"},
		{"ngango", "ngango", "ngango", "NOUN", "", "PLANT", 6, "palm-nut-fibers", "palm nut fibers"},
		{"ngangu", "ngangü", "ngangü", "ADJ", "", "HOW", 1, "powerful", "powerful, firm; difficult"},
		{"ngangu", "ngangü", "ngangü", "ADV", "", "HOW", 1, "powerfully", "powerfully, firmly"},
		{"ngangu", "ngangü", "ngangü", "NOUN", "", "HOW", 1, "power", "hardship, difficulty, force, power"},
		{"ngao", "ngâo", "ngâo", "NOUN", "", "NATURE", 3, "smoke", "smoke"},
		{"ngao", "ngâo", "ngâo", "NOUN", "", "OBJ", 6, "tobacco-or-cigarette", "tobacco, cigarette"},
		{"ngapo", "ngâpô", "ngâpô", "NOUN", "", "OBJ", 6, "hoe", "hoe"},
		{"ngasa", "ngäsa", "ngäsa", "NOUN", "", "ANIM", 2, "goat", "goat"},
		{"ngasi", "ngâsî", "ngâsî", "NOUN", "", "BODY", 4, "sneeze", "sneeze"},
		{"ngba", "ngbâ", "ngbâ", "VERB", "Subcat=Intr", "STATE", 1, "stay-unchanged", "remain (the same), stay (unchanged)"},
		{"ngba", "ngbä", "ngbä", "NOUN", "", "ANIM", 5, "buffalo", "buffalo"},
		{"ngbaa", "ngbâa", "ngbâa", "VERB", "Subcat=Intr", "CIVIL", 3, "slave", "slave"},
		{"ngbadara", "ngbadârâ", "ngbadârâ", "NOUN", "", "CIVIL", 4, "theater", "theater, play"},
		{"ngbako", "ngbâko", "ngbâko", "NOUN", "", "OBJ", 5, "liquor", "liquor distilled from corn and manioc"},
		{"ngbakongo", "ngbâ-kongö", "ngbâ-kongö", "NOUN", "", "FAMILY", 5, "metis", "[lit: remain|Congo]: half European, metis"},
		{"ngbalo", "ngbâlo", "ngbâlo", "NOUN", "", "SICK", 4, "headache", "headache, migraine"},
		{"ngbanga", "ngbanga", "ngbanga", "NOUN", "", "HOW", 1, "reason-court-justice", "reason, process, judgment, court, justice"},
		{"ngbangati", "ngbanga tî", "ngbanga tî", "ADP", "", "HOW", 1, "because-of", "[lit: reason|of]: because of"},
		{"ngbangatinye", "ngbanga tî nye", "ngbanga tî nyɛ", "ADV", "", "HOW", 1, "why", "[lit: reason|of|what]: why"},
		{"ngbangatiso", "ngbanga tî sô", "ngbanga tî sô", "ADV", "", "HOW", 1, "because-of-this", "[lit: reason|of|this]: because of this"},
		{"ngbangatiso", "ngbanga tî sô", "ngbanga tî sô", "CCONJ", "", "HOW", 1, "because", "[lit: reason|of|that]: because"},
		{"ngbangba", "ngbângbä", "ngbângbä", "NOUN", "", "ANIM", 6, "sea-toad", "sea toad"},
		{"ngbangba", "ngbängbâ", "ngbängbâ", "NOUN", "", "BODY", 3, "jaw-cheek", "jaw, cheek"},
		{"ngbangbo", "ngbangbo", "ngbangbo", "ADJ", "NumType=Ord", "NUM", 2, "hundred", "hundred"},
		{"ngbangbo", "ngbangbo", "ngbangbo", "NUM", "NumType=Card", "NUM", 2, "hundred", "hundred"},
		{"ngbangbotukia", "ngbangbo-tukîa", "ngbangbo-tukîa", "NOUN", "", "NUM", 5, "hectare", "[lit: hecto|are]: hectare=~2.5 acres"},
		{"ngbangerengu", "ngbângêrêngû", "ngbângêrêngû", "NOUN", "", "SICK", 6, "convulsions", "convulsions"},
		{"ngbanzoni", "ngbâ nzönî", "ngbâ nzɔ̈nî", "INTERJ", "Mood=Opt", "INTERACT", 1, "Goodbye", "[lit: remain|well]: Goodbye! (said to those staying)"},
		{"ngbene", "ngbêne", "ngbɛ̂nɛ", "ADJ", "", "HOW", 4, "old", "old, aged"},
		{"ngbengbe", "ngbëngbë", "ngbëngbë", "NOUN", "", "OBJ", 4, "fish-net", "fish net"},
		{"ngbenge", "ngbengë", "ngbengë", "NOUN", "", "OBJ", 5, "guitar", "guitar"},
		{"ngberena", "ngbêrênâ", "ngbêrênâ", "NOUN", "", "OBJ", 4, "bracelet-anklet", "bracelet, anklet"},
		{"ngberere", "ngberere", "ngbɛrɛrɛ", "NOUN", "", "WHEN", 5, "October", "October"},
		{"ngbii", "ngbii", "ngbii", "ADV", "", "WHEN", 1, "long-time", "long time"},
		{"ngbiiasina", "ngbii asî na", "ngbii asî na", "ADP", "", "WHEN", 1, "not-until", "[lit: long time|arrive|at]: not until"},
		{"ngbiisi", "ngbii sï", "ngbii sï", "SCONJ", "", "WHEN", 1, "not-until", "[lit: long time|arrive]: not until"},
		{"ngbo", "ngbö", "ngbɔ̈", "NOUN", "", "ANIM", 2, "snake", "snake"},
		{"ngbo", "ngbö", "ngbɔ̈", "NOUN", "", "FAMILY", 4, "twins", "twins"},
		{"ngboko", "ngbökö", "ngbɔ̈kɔ̈", "NOUN", "", "PLANT", 6, "sugarcane", "sugarcane"},
		{"ngbonda", "ngbondä", "ngbondä", "NOUN", "", "BODY", 3, "bottom-rear-base", "butt, bottom, rear, base"},
		{"ngbondo", "ngbondô", "ngbondô", "ADJ", "", "CIVIL", 6, "valuable", "precious, valuable"},
		{"ngbonga", "ngbonga", "ngbonga", "NOUN", "", "OBJ", 5, "tambourine", "[African-made]: wooden tambourine"},
		{"ngbonga", "ngbonga", "ngbonga", "NOUN", "", "TREE", 4, "trunk", "trunk"},
		{"ngbongboro", "ngbongbôro", "ngbongbôro", "ADJ", "", "NUM", 3, "huge-massive", "huge, massive, chunk of"},
		{"ngbongboro", "ngbongbörö", "ngbɔngbɔ̈rɔ̈", "NOUN", "", "OBJ", 3, "musket", "musket"},
		{"ngbongboto", "ngbongboto", "ngbongboto", "NOUN", "", "FOOD", 6, "palm-oil-dregs", "palm oil dregs"},
		{"ngbongo", "ngbôngô", "ngbɔ̂ngɔ̂", "NOUN", "", "BODY", 5, "hock", "[animal leg] hock"},
		{"ngboto", "ngbôto", "ngbôto", "NOUN", "", "WHERE", 3, "detour", "detour"},
		{"ngbuku", "ngbûku", "ngbûku", "NOUN", "", "INTERACT", 4, "riddle", "riddle"},
		{"ngbundangbu", "ngbundangbu", "ngbundangbu", "ADJ", "NumType=Ord", "NUM", 2, "billion", "billion"},
		{"ngbundangbu", "ngbundangbu", "ngbundangbu", "NUM", "NumType=Card", "NUM", 2, "billion", "billion"},
		{"ngbungbu", "ngbungbu", "ngbungbu", "NOUN", "", "WHEN", 6, "August", "August"},
		{"ngbuta", "ngbuta", "ngbuta", "NOUN", "", "OBJ", 4, "harpoon", "harpoon"},
		{"nge", "nge", "ngɛ", "NOUN", "", "PLANT", 6, "sugarcane", "sugarcane"},
		{"nge", "nge", "ngɛ", "VERB", "Subcat=Intr", "STATE", 2, "be-thin", "shrink, be thin, scrawny"},
		{"ngende", "ngendë", "ngɛndɛ̈", "NOUN", "", "OBJ", 2, "chair", "chair"},
		{"ngenge", "ngenge", "ngenge", "NOUN", "", "FISH", 6, "little-fish", "small fry, little fish"},
		{"ngengo", "ngëngö", "ngɛ̈ngɔ̈", "VERB", "Subcat=Intr|VerbForm=Vnoun", "STATE", 3, "thinness", "thinness, scrawniness"},
		{"ngere", "ngêrë", "ngêrë", "NOUN", "", "CIVIL", 2, "price", "business, commerce, shopping; price, cost, value"},
		{"ngia", "ngîâ", "ngîâ", "NOUN", "", "FEEL", 1, "pleasure", "joy, pleasure"},
		{"ngiba", "ngiba", "ngiba", "NOUN", "", "SICK", 6, "leprosy", "leprosy"},
		{"nginza", "nginza", "nginza", "NOUN", "", "CIVIL", 2, "money", "money"},
		{"ngira", "ngira", "ngira", "NOUN", "", "CIVIL", 5, "taboo-prohibition", "taboo, prohibition"},
		{"ngiriba", "ngiriba", "ngiriba", "NOUN", "", "SICK", 6, "leprosy", "leprosy"},
		{"ngiriki", "ngïrïkï", "ngïrïkï", "NOUN", "", "TREE", 6, "kola-tree", "kola tree"},
		{"ngo", "ngo", "ngo", "NOUN", "", "OBJ", 3, "handle-crook", "handle, crook"},
		{"ngo", "ngo", "ngɔ", "NOUN", "", "BODY", 2, "fetus-pregnancy", "fetus, pregnancy"},
		{"ngo", "ngo", "ngɔ", "NOUN", "", "OBJ", 5, "tambourine", "[African-made]: wooden tambourine"},
		{"ngo", "ngô", "ngô", "VERB", "Subcat=Tran", "ACT", 3, "bend-twist-roll", "bend, twist, roll"},
		{"ngo", "ngö", "ngɔ̈", "NOUN", "", "BODY", 3, "canoe", "canoe"},
		{"ngo", "ngö", "ngɔ̈", "NOUN", "", "OBJ", 3, "canoe", "canoe"},
		{"ngoi", "ngoi", "ngoi", "NOUN", "", "WHEN", 2, "season-time", "season, time, moment, epoch"},
		{"ngoitiburu", "ngoi tî burü", "ngoi tî burü", "NOUN", "", "WHEN", 2, "dry-season", "dry season"},
		{"ngoitingu", "ngoi tî ngû", "ngoi tî ngû", "NOUN", "", "WHEN", 2, "rainy-season", "rainy season"},
		{"ngolo", "ngôlö", "ngɔ̂lɔ̈", "NOUN", "", "OBJ", 4, "fish-net", "wicker fish net"},
		{"ngolo", "ngölo", "ngɔ̈lɔ", "NOUN", "", "INTERACT", 4, "provocation-sarcasm", "provocation, sarcasm"},
		{"ngombe", "ngombe", "ngombe", "NOUN", "", "OBJ", 3, "gun-tube", "gun, tube"},
		{"ngonda", "ngonda", "ngonda", "NOUN", "", "NATURE", 2, "wilderness", "wilderness"},
		{"ngonga", "ngonga", "ngonga", "NOUN", "", "WHEN", 2, "hour", "hour"},
		{"ngongbi", "ngôngbi", "ngôngbi", "NOUN", "", "ACT", 3, "fold", "fold, pleat, wrinkle"},
		{"ngongbi", "ngôngbi", "ngôngbi", "VERB", "Aspect=Imp", "ACT", 3, "fold-in-half", "fold in half, roll"},
		{"ngonza", "ngonzâ", "ngonzâ", "NOUN", "", "ANIM", 6, "snail", "snail"},
		{"ngonzo", "ngonzo", "ngɔnzɔ", "NOUN", "", "BODY", 5, "bile-gall-anger", "bile, gall, anger"},
		{"ngoro", "ngoro", "ngoro", "VERB", "Subcat=Tran", "WHERE", 5, "surround", "surround, encircle; welcome a guest or be welcomed as one"},
		{"ngoro", "ngoro", "ngɔrɔ", "NOUN", "", "FISH", 6, "catfish", "catfish"},
		{"ngorongbi", "ngôrôngbi", "ngôrôngbi", "VERB", "Aspect=Imp|Subcat=Tran", "WHERE", 5, "surround-a-group", "surround a group, encircle a group"},
		{"ngorongbo", "ngöröngbö", "ngöröngbö", "VERB", "Subcat=Tran", "WHERE", 5, "gutter", "gutter, trench"},
		{"ngoropangi", "ngoropangi", "ngoropangi", "NOUN", "", "INTERACT", 6, "echo", "echo"},
		{"ngoti", "ngo-tï", "ngo-tï", "NOUN", "", "BODY", 6, "elbow", "elbow"},
		{"ngu", "ngû", "ngû", "NOUN", "", "NATURE", 1, "water-or-year", "water, liquid, humidity; year"},
		{"nguba", "nguba", "nguba", "NOUN", "", "NATURE", 6, "bark", "bark, papyrus, fabric, ancient manuscript"},
		{"ngube", "ngubë", "ngubë", "NOUN", "", "WHEN", 5, "April", "April"},
		{"ngubu", "ngubü", "ngubü", "NOUN", "", "ANIM", 6, "hippopotamus", "hippopotamus"},
		{"ngui", "ngui", "ngui", "NOUN", "", "ANIM", 6, "monkey", "colobus monkey"},
		{"nguingo", "ngû-îngö", "ngû-îngɔ̈", "NOUN", "", "FOOD", 3, "ocean", "[lit: water|salt]: ocean"},
		{"ngulavu", "ngû-lavu", "ngû-lavu", "NOUN", "", "FOOD", 6, "honey", "[lit: water|bee]: honey"},
		{"ngule", "ngû-lê", "ngû-lɛ̂", "NOUN", "", "BODY", 2, "tears", "[lit: water|eye]: tears"},
		{"ngumba", "ngûmbâ", "ngûmbâ", "NOUN", "", "SICK", 4, "enlarged-liver", "enlarged liver"},
		{"ngumu", "ngumu", "ngumu", "NOUN", "", "ALT WORD FOR", 6, "mange", "särä"},
		{"ngunde", "ngundë", "ngundë", "NOUN", "", "ANIM", 4, "crocodile", "crocodile"},
		{"ngungu", "ngungu", "ngungu", "NOUN", "", "ANIM", 3, "mosquito", "mosquito"},
		{"ngungunza", "ngû-ngunzä", "ngû-ngunzä", "ADJ", "", "COLOR", 3, "green", "[lit: water (of)|manioc leaves]: green"},
		{"ngunza", "ngunzä", "ngunzä", "NOUN", "", "FOOD", 2, "manioc-leaves", "manioc leaves"},
		{"ngunzapa", "ngû-nzapä", "ngû-nzapä", "NOUN", "", "NATURE", 1, "rain", "[lit: water|God]: rain"},
		{"ngunzapa", "ngûnzapä", "ngûnzapä", "NOUN", "", "NATURE", 1, "rain", "[lit: water|God]: rain"},
		{"nguru", "ngûru", "ngûru", "NOUN", "", "ANIM", 2, "pig", "pig"},
		{"ngusu", "ngusü", "ngusü", "NOUN", "", "ANIM", 5, "larva", "jigger, sandflea, larva"},
		{"ngutikoli", "ngû tî kôlï", "ngû tî kɔ̂lï", "NOUN", "", "BODY", 6, "semen", "[lit: water|of|man]: sperm, semen"},
		{"ngutile", "ngû tî lê", "ngû tî lɛ̂", "NOUN", "", "BODY", 2, "tears", "[lit: water|of|eye]: tears"},
		{"ngutimbeti", "ngû tî mbëtï", "ngû tî mbɛ̈tï", "NOUN", "", "OBJ", 2, "ink", "[lit: water|of|paper]: ink"},
		{"ngutime", "ngû tî me", "ngû tî mɛ", "NOUN", "", "DRINK", 3, "breast-milk", "[lit: water|of|teat]: breast milk"},
		{"ngutimetibagara", "ngû tî me tî bâgara", "ngû tî mɛ tî bâgara", "NOUN", "", "DRINK", 3, "milk", "[lit: water|of|teat|of|cow]: milk"},
		{"ngutinyon", "ngû tî nyön", "ngû tî nyön", "NOUN", "", "NATURE", 3, "potable-water", "[lit: water|to|drink]: potable water"},
		{"ngutinzapa", "ngû tî nzapä", "ngû tî nzapä", "NOUN", "", "ALT SP FOR", 9, "rain", "ngûnzapä"},
		{"ngutitere", "ngû tî terê", "ngû tî tɛrɛ̂", "NOUN", "", "BODY", 2, "sweat", "[lit: water|of|body]: sweat"},
		{"ngutivuru", "ngû tî vurü", "ngû tî vurü", "NOUN", "", "BODY", 2, "pus", "[lit: water|of|whiteness]: pus"},
		{"ngutiyanga", "ngû tî yängâ", "ngû tî yängâ", "NOUN", "", "BODY", 2, "saliva", "[lit: water|of|mouth]: saliva"},
		{"nguyenga", "ngû-yenga", "ngû-yenga", "NOUN", "", "WHEN", 2, "anniversary", "[lit: year|feast]: anniversary"},
		{"ni", "nî", "nî", "DET", "PronType=Art", "WHICH", 1, "the", "the"},
		{"ni", "nî", "nî", "PRON", "Animacy=Inan|Case=Acc|Num=Sing|Person=3|PronType=Det", "WHICH", 1, "it", "it"},
		{"ni", "nï", "nï", "PRON", "Num=Sing|Person=3|PronType=Rel", "WHICH", 2, "he-she", "[indirect style]: he, she"},
		{"nigisi", "nîgisi", "nîgisi", "NOUN", "", "NUM", 6, "zero", "zero"},
		{"nika", "nika", "nika", "VERB", "Subcat=Tran", "ACT", 2, "grind", "grind, crush, mill"},
		{"nikpa", "nïkpä", "nïkpä", "NOUN", "", "ANIM", 5, "leech", "leech"},
		{"ninga", "nînga", "nînga", "VERB", "Aspect=Iter|Subcat=Intr", "STATE", 2, "last-linger", "last, delay, stay"},
		{"no", "nô", "nô", "NOUN", "", "MOVE", 6, "gait-walk-step", "gait, walk, step"},
		{"no", "nö", "nö", "VERB", "Subcat=Intr", "MOVE", 6, "go-walk-step", "go, walk, step"},
		{"noko", "nökö", "nɔ̈kɔ̈", "NOUN", "", "ALT WORD FOR", 9, "maternal-relative", "kôya"},
		{"nyama", "nyama", "nyama", "NOUN", "", "ANIM", 1, "animal-or-meat", "animal, meat, beast, [after Fr: bete]: idiot"},
		{"nyau", "nyâu", "nyâu", "NOUN", "", "ANIM", 2, "cat-or-hypocrite", "cat, [fig]: hypocrite"},
		{"nye", "nye", "nyɛ", "NOUN", "", "WHICH", 1, "what", "what"},
		{"nyene", "nyenë", "nyɛnɛ̈", "NOUN", "", "BODY", 6, "buttocks", "butt, buttocks"},
		{"nyenye", "nyenye", "nyɛnyɛ", "NOUN", "", "WHEN", 5, "January", "January"},
		{"nyenyeke", "nyenyekê", "nyɛnyɛkɛ̂", "ADV", "", "BODY", 4, "graceful", "graceful, elegant"},
		{"nyi", "nyï", "nyï", "NOUN", "", "FAMILY", 3, "child", "[6mo-4yrs]: child"},
		{"nyikoli", "nyï-kôlï", "nyï-kɔ̂lï", "NOUN", "Gender=Masc", "FAMILY", 3, "boy", "[6mo-4yrs]: boy"},
		{"nyiliti", "nyï-li-tï", "nyï-li-tï", "NOUN", "", "BODY", 4, "ring-finger", "[lit: baby|finger]: ring finger"},
		{"nyindu", "nyï-ndü", "nyï-ndü", "NOUN", "", "FAMILY", 4, "orphan", "[lit: baby|widowed]: orphan"},
		{"nyingambi", "nyï-ngambi", "nyï-ngambi", "NOUN", "", "FAMILY", 4, "baby-younger-sibling", "[lit: baby|younger]: baby brother, baby sister"},
		{"nyiwali", "nyï-wâlï", "nyï-wâlï", "NOUN", "Gender=Fem", "FAMILY", 3, "girl", "[6mo-4yrs]: girl"},
		{"nyiwanda", "nyï-wanda", "nyï-wanda", "NOUN", "", "FAMILY", 6, "bastard", "[lit: baby|adultery]: bastard"},
		{"nyon", "nyön", "nyön", "VERB", "Subcat=Tran", "BODY", 2, "drink-or-inhale", "drink, inhale"},
		{"nyonmanga", "nyön mânga", "nyön mânga", "VERB", "Subcat=Intr", "BODY", 2, "smoke", "[lit: drink|tobacco]: smoke"},
		{"nyonmene", "nyön mênë", "nyön mɛ̂nɛ̈", "VERB", "Subcat=Intr", "BODY", 2, "make-a-blood-pact", "[lit: drink|blood]: make a blood pact"},
		{"nza", "nza", "nza", "NOUN", "", "OBJ", 5, "antenna", "antenna"},
		{"nza", "nzä", "nzä", "NOUN", "", "PLANT", 5, "palm-fruit", "fan palm fruit"},
		{"nzabi", "nzabï", "nzabï", "NOUN", "", "FISH", 6, "Nile-perch", "Nile perch, [French]: capitaine"},
		{"nzai", "nzaï", "nzaï", "NOUN", "", "FISH", 5, "Nile-perch", "Nile perch, [French]: capitaine"},
		{"nzangi", "nzângi", "nzângi", "NOUN", "", "OBJ", 5, "backpack", "backpack"},
		{"nzanza", "nzanza", "nzanza", "NOUN", "", "PLANT", 6, "reed", "reed, rush"},
		{"nzanze", "nzanzë", "nzanzë", "NOUN", "", "PLANT", 6, "twig", "twig"},
		{"nzapa", "nzapä", "nzapä", "INTERJ", "", "GOD", 1, "Lord", "[lit: God]: Lord (willing)! Oh my!, I swear!"},
		{"nzapa", "nzapä", "nzapä", "NOUN", "", "GOD", 1, "God", "God"},
		{"nzapababa", "nzapä-babâ", "nzapä-babâ", "NOUN", "", "GOD", 1, "God-the-Father", "God the Father, (Our) Father"},
		{"nzara", "nzara", "nzara", "NOUN", "", "STATE", 1, "hunger-famine-appetite-desire", "hunger, famine, appetite, desire, yearning"},
		{"nzaratingu", "nzara tî ngû", "nzara tî ngû", "NOUN", "", "STATE", 1, "thirst", "thirst"},
		{"nzayu", "nzayü", "nzayü", "NOUN", "", "FISH", 6, "Nile-perch", "Nile perch, [French]: capitaine"},
		{"nze", "nze", "nzɛ", "NOUN", "", "BODY", 2, "menstruation", "[lit: moon]: menstruation"},
		{"nze", "nze", "nzɛ", "NOUN", "", "NATURE", 2, "moon", "moon"},
		{"nze", "nze", "nzɛ", "NOUN", "", "WHEN", 2, "month", "[lit: moon]: month"},
		{"nzebaleoko", "nze balë-ôko", "nzɛ balë-ɔ̂kɔ", "NOUN", "", "WHEN", 2, "October", "[lit: moon|ten]: October"},
		{"nzebaleokonaoko", "nze balë-ôko na ôko", "nzɛ balë-ɔ̂kɔ na ɔ̂kɔ", "NOUN", "", "WHEN", 2, "November", "[lit: moon|eleven]: November"},
		{"nzebaleokonause", "nze balë-ôko na ûse", "nzɛ balë-ɔ̂kɔ na ûse", "NOUN", "", "WHEN", 2, "December", "[lit: moon|twelve]: December"},
		{"nzeen", "nzêen", "nzêen", "VERB", "Subcat=Intr", "FEEL", 3, "be-discouraged", "be discouraged, depressed, tired"},
		{"nzeennapekoti", "nzêen na pekö tî", "nzêen na pekö tî", "VERB", "Subcat=Tran", "FEEL", 3, "be-exasperated-by", "[lit: be discouraged|to|back|of]: be exasperated by, be at wits end with"},
		{"nzege", "nzëgë", "nzëgë", "NOUN", "", "ANIM", 5, "frog", "frog"},
		{"nzeli", "nzeli", "nzɛli", "NOUN", "", "ACT", 3, "razor", "razor"},
		{"nzembarambara", "nze mbârâmbârâ", "nzɛ mbârâmbârâ", "NOUN", "", "WHEN", 2, "July", "[lit: moon|seven]: July"},
		{"nzemeambe", "nze meambe", "nzɛ meambe", "NOUN", "", "WHEN", 2, "August", "[lit: moon|eight]: August"},
		{"nzene", "nzêne", "nzɛ̂nɛ", "ADJ", "", "NUM", 4, "small", "small, little"},
		{"nzene", "nzënë", "nzɛ̈nɛ̈", "NOUN", "", "BODY", 4, "claw-fingernail-toenail", "claw, fingernail, toenail"},
		{"nzengumbaya", "nze ngümbâyä", "nzɛ ngümbâyä", "NOUN", "", "WHEN", 2, "September", "[lit: moon|nine]: September"},
		{"nzenze", "nzenze", "nzɛnzɛ", "NOUN", "", "OBJ", 3, "machete", "machete"},
		{"nzeoko", "nze ôko", "nzɛ ɔ̂kɔ", "NOUN", "", "WHEN", 2, "January", "[lit: moon|one]: January"},
		{"nzeoku", "nze okü", "nzɛ ɔkü", "NOUN", "", "WHEN", 2, "May", "[lit: moon|five]: May"},
		{"nzeomene", "nze omenë", "nzɛ omɛnë", "NOUN", "", "WHEN", 2, "June", "[lit: moon|six]: June"},
		{"nzeota", "nze otâ", "nzɛ otâ", "NOUN", "", "WHEN", 2, "March", "[lit: moon|three]: March"},
		{"nzepere", "nzêpêrê", "nzêpêrê", "NOUN", "", "OBJ", 6, "arrow", "arrow"},
		{"nzere", "nzere", "nzɛrɛ", "VERB", "Subcat=Intr", "FEEL", 1, "be-delicious-or-satisfying", "be delicious, tasty, agreeable, pleasant, satisfying"},
		{"nzere", "nzerë", "nzerë", "NOUN", "", "COLOR", 3, "color", "color, tint"},
		{"nzerenabeti", "nzere na bê tî", "nzɛrɛ na bɛ̂ tî", "VERB", "Subcat=Tran", "FEEL", 1, "please", "[lit: be tasty|to|heart|of]: please"},
		{"nzerengo", "nzërëngö", "nzɛ̈rɛ̈ngɔ̈", "ADJ", "", "FEEL", 1, "delicious-or-satisfying", "delicious, tasty, agreeable, pleasant, satisfying"},
		{"nzeretinduzu", "nzerë tî ndüzü", "nzerë tî ndüzü", "NOUN", "", "COLOR", 3, "blue", "blue"},
		{"nzeretiye", "nzerë-tî-yê", "nzerë-tî-yê", "NOUN", "", "COLOR", 3, "image-or-picture", "image, picture, drawing, icon"},
		{"nzeuse", "nze ûse", "nzɛ ûse", "NOUN", "", "WHEN", 2, "February", "[lit: moon|two]: February"},
		{"nzeusio", "nze usïö", "nzɛ usïö", "NOUN", "", "WHEN", 2, "April", "[lit: moon|four]: April"},
		{"nzi", "nzï", "nzï", "VERB", "", "ACT", 2, "steal-burglerize-rob", "steal, burglerize, rob"},
		{"nzinangonga", "nzîna-ngonga", "nzîna-ngonga", "NOUN", "", "WHEN", 2, "minute", "[lit: deci|hour]: minute"},
		{"nzingo", "nzïngö", "nzïngɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 2, "theft", "theft"},
		{"nzo", "nzö", "nzɔ̈", "NOUN", "", "FOOD", 3, "corn", "corn"},
		{"nzoba", "nzö-bä", "nzɔ̈-bä", "NOUN", "", "INTERACT", 1, "benediction", "benediction"},
		{"nzobe", "nzö-bê", "nzɔ̈-bê", "NOUN", "", "HOW", 1, "kindness", "kindness"},
		{"nzobia", "nzö-bîâ", "nzɔ̈-bîâ", "NOUN", "", "GOD", 1, "psalm", "psalm"},
		{"nzodeba", "nzö-dëbä", "nzɔ̈-dëbä", "NOUN", "", "INTERACT", 3, "blessing", "blessing"},
		{"nzombo", "nzombö", "nzɔmbɔ̈", "NOUN", "", "FISH", 6, "electric-catfish", "electric catfish"},
		{"nzongoro", "nzöngörö", "nzɔ̈ngɔ̈rɔ̈", "NOUN", "", "ANIM", 6, "blue-orange-lizard", "blue orange lizard"},
		{"nzoni", "nzönî", "nzɔ̈nî", "ADJ", "", "HOW", 2, "good", "good"},
		{"nzoni", "nzönî", "nzɔ̈nî", "ADV", "", "HOW", 1, "well", "well"},
		{"nzoni", "nzönî", "nzɔ̈nî", "NOUN", "", "HOW", 3, "goodness", "goodness"},
		{"nzonisango", "nzönî-sango", "nzɔ̈nî-sango", "NOUN", "", "GOD", 4, "gospel", "[lit: good|news]: gospel, good tidings"},
		{"nzoroko", "nzorôko", "nzɔrɔ̂kɔ", "ADJ", "", "COLOR", 6, "yellow", "[lit: bodypaint]: yellow"},
		{"nzoroko", "nzorôko", "nzɔrɔ̂kɔ", "NOUN", "", "COMPUTER", 6, "site-website", "[neologism]: (web)site"},
		{"nzoroko", "nzorôko", "nzɔrɔ̂kɔ", "NOUN", "", "TREE", 6, "body-paint-tattoo-scarification", "body paint, tattoo, scarification"},
		{"nzorokokombe", "nzorôko kömbë", "nzɔrɔ̂kɔ kömbë", "ADJ", "", "COLOR", 6, "yellow", "[lit: bodypaint|yellowfruit]: yellow"},
		{"nzosango", "nzö-sango", "nzɔ̈-sango", "NOUN", "", "GOD", 4, "gospel", "[lit: good|news]: gospel, good tidings"},
		{"o", "o", "o", "PART", "Polite=Form", "INTERACT", 1, "[politeness]", "[politeness]"},
		{"o", "o", "o", "VERB", "Subcat=Tran", "ACT", 6, "kill", "kill"},
		{"oke", "ôke", "ɔ̂kɛ", "ADV", "", "NUM", 2, "how-many", "how many"},
		{"oko", "ôko", "ɔ̂kɔ", "ADJ", "NumType=Ord", "NUM", 2, "one", "one"},
		{"oko", "ôko", "ɔ̂kɔ", "NUM", "NumType=Card", "NUM", 2, "one", "one"},
		{"okoape", "ôko äpe", "ɔ̂kɔ äpɛ", "ADJ", "", "NUM", 2, "no", "no"},
		{"okoape", "ôko äpe", "ɔ̂kɔ äpɛ", "ADV", "", "NUM", 2, "never-no", "never, no"},
		{"okooko", "ôko-ôko", "ɔ̂kɔ-ɔ̂kɔ", "ADJ", "", "NUM", 2, "each-one", "each one"},
		{"okopepe", "ôko pëpe", "ɔ̂kɔ pɛ̈pɛ", "ADJ", "", "NUM", 2, "no", "no"},
		{"okopepe", "ôko pëpe", "ɔ̂kɔ pɛ̈pɛ", "ADV", "", "NUM", 2, "never-no", "never, no"},
		{"oku", "okü", "ɔkü", "ADJ", "NumType=Ord", "NUM", 2, "five", "five"},
		{"oku", "okü", "ɔkü", "NUM", "NumType=Card", "NUM", 2, "five", "five"},
		{"omene", "omenë", "omɛnë", "ADJ", "NumType=Ord", "NUM", 2, "six", "six"},
		{"omene", "omenë", "omɛnë", "NUM", "NumType=Card", "NUM", 2, "six", "six"},
		{"ota", "otâ", "otâ", "ADJ", "NumType=Ord", "NUM", 2, "three", "three"},
		{"ota", "otâ", "otâ", "NUM", "NumType=Card", "NUM", 2, "three", "three"},
		{"oto", "ôtö", "ɔ̂tɔ̈", "NOUN", "", "NATURE", 3, "hill", "hill, mountain"},
		{"pa", "pâ", "pâ", "NOUN", "", "INTERACT", 3, "slander", "slander, false accusation"},
		{"pa", "pä", "pä", "VERB", "Subcat=Tran", "INTERACT", 3, "slander", "slander, falsely accuse"},
		{"pafungula", "pâ-fungûla", "pâ-fungûla", "NOUN", "", "OBJ", 5, "password", "[lit: false-unlock] password"},
		{"pairiti", "pä ïrï tî", "pä ïrï tî", "VERB", "Subcat=Tran", "INTERACT", 3, "slander-the-good-name-of", "slander the good name of"},
		{"pakapaka", "pakapâka", "pakapâka", "NOUN", "", "NATURE", 5, "turbulence", "[onomatopeia]: wake, turbulence, propellor"},
		{"pakara", "pakara", "pakara", "NOUN", "", "WHO", 5, "Mister-honest-man", "Mister, honest man"},
		{"palata", "palâta", "palâta", "NOUN", "", "OBJ", 5, "medal", "medal"},
		{"pambo", "pambo", "pambo", "NOUN", "", "ANIM", 5, "earthworm", "earthworm"},
		{"pamboti", "pâmbo-tï", "pâmbo-tï", "NOUN", "", "BODY", 5, "shoulder", "shoulder"},
		{"pande", "pandë", "pandë", "NOUN", "", "HOW", 2, "model-example", "model, example, type, norm, answer key, pattern"},
		{"panga", "pängä", "pängä", "NOUN", "", "SICK", 6, "asthma", "asthma"},
		{"papa", "papa", "papa", "NOUN", "", "INTERACT", 2, "argument", "argument, dispute"},
		{"papa", "papa", "papa", "NOUN", "", "OBJ", 2, "spoon", "spoon, spoonful"},
		{"papa", "pâpa", "pâpa", "VERB", "Subcat=Intr", "INTERACT", 2, "argue", "argue"},
		{"papa", "pâpâ", "pâpâ", "NOUN", "", "OBJ", 2, "sandal", "sandal"},
		{"papaye", "papayë", "papayë", "NOUN", "", "FOOD", 2, "papaya", "papaya"},
		{"para", "pärä", "pärä", "NOUN", "", "FOOD", 2, "egg", "egg"},
		{"para", "pärä", "pärä", "NOUN", "", "NUM", 2, "zero", "zero"},
		{"paragere", "pärä-gerë", "pärä-gɛrɛ̈", "NOUN", "", "BODY", 4, "heel", "heel"},
		{"parati", "pärä-tï", "pärä-tï", "NOUN", "", "BODY", 4, "fist", "fist"},
		{"pasa", "päsä", "päsä", "NOUN", "", "STATE", 6, "good-luck", "good luck"},
		{"pasaporo", "päsäpôro", "päsäpɔ̂rɔ", "NOUN", "", "CIVIL", 4, "passport", "passport"},
		{"pasee", "pasëe", "pasëe", "VERB", "Subcat=Tran", "ACT", 4, "iron", "[Fr: repasser]: iron (clothes)"},
		{"pasi", "pâsi", "pâsi", "NOUN", "", "CIVIL", 3, "misery", "misery"},
		{"pata", "patâ", "patâ", "ADJ", "", "INTERACT", 4, "whispering", "whispering"},
		{"pata", "pâta", "pâta", "NOUN", "", "CIVIL", 1, "penny", "penny, cent [lowest valued coin=5 francs CFA in CAR]"},
		{"patara", "patärä", "patärä", "NOUN", "", "GAME", 4, "dice-or-negotiation", "dice; negotiation"},
		{"pe", "pe", "pe", "VERB", "Subcat=Tran", "ACT", 2, "intertwine-braid", "entwine, intertwine, braid"},
		{"pe", "pë", "pɛ̈", "VERB", "Subcat=Intr", "INTERACT", 2, "fan", "[hand movement]: wave, fan, winnow"},
		{"peke", "pekë", "pɛkɛ̈", "NOUN", "", "DRINK", 5, "palm-wine", "Raffia palm wine"},
		{"peko", "pekô", "pekô", "NOUN", "", "ALT SP FOR", 9, "back", "pekö"},
		{"peko", "pekö", "pekö", "NOUN", "", "BODY", 1, "back", "back"},
		{"peko", "pekö", "pekö", "NOUN", "", "WHEN", 1, "moment", "moment, duration of time"},
		{"peko", "pekö", "pekö", "NOUN", "", "WHERE", 1, "behind-after-following", "behind, after, following"},
		{"pembe", "pëmbë", "pɛ̈mbɛ̈", "NOUN", "", "BODY", 2, "tooth", "tooth"},
		{"penda", "pendä", "pendä", "NOUN", "", "WHEN", 2, "influence-result-consequence", "trace, influence, result, consequence"},
		{"pendere", "pendere", "pɛndɛrɛ", "ADJ", "", "HOW", 1, "beautiful", "beautiful, pretty"},
		{"penderekoli", "pendere kôlï", "pɛndɛrɛ kɔ̂lï", "NOUN", "Gender=Masc", "FAMILY", 3, "adolescent-man", "[17-21yrs and unmarried]: adolescent man"},
		{"penderewali", "pendere wâlï", "pɛndɛrɛ wâlï", "NOUN", "Gender=Fem", "FAMILY", 3, "adolescent-woman", "[17-21yrs and unmarried]: adolescent woman"},
		{"pengo", "pëngö", "pɛ̈ngɔ̈", "VERB", "Subcat=Intr", "INTERACT", 2, "wink-or-flap", "wink, blink, smack (lips), bat (eyes), flap (arms or wings)"},
		{"penze", "penze", "pɛnzɛ", "NOUN", "", "NUM", 6, "segment", "segment, lengthwise piece (e.g. of rope)"},
		{"pepe", "pëpe", "pɛ̈pɛ", "PART", "Polarity=Neg", "HOW", 1, "not", "not"},
		{"pere", "pêrë", "pêrë", "NOUN", "", "NATURE", 1, "straw-grass-brush", "straw, grass, brush"},
		{"pete", "pete", "pɛtɛ", "VERB", "Subcat=Tran", "ACT", 2, "mash", "mash, puree, press"},
		{"pete", "pête", "pɛ̂tɛ", "NOUN", "", "WHEN", 6, "March", "March"},
		{"pete", "pëtë", "pɛ̈tɛ̈", "NOUN", "", "ACT", 2, "pressure-printing", "pressure, imprint, printing"},
		{"pida", "pîda", "pîda", "VERB", "Subcat=Intr", "STATE", 3, "stick", "stick, adhere"},
		{"pika", "pîka", "pîka", "VERB", "Subcat=Tran", "ACT", 2, "beat-strike-fall-play", "beat, strike, pound, shoot, nail, [rain]: fall, play (music, sports)"},
		{"pikahon", "pîka-hôn", "pîka-hôn", "VERB", "Subcat=Intr", "ACT", 2, "sneeze", "[lit: strike|nose]: sneeze"},
		{"pikakpekena", "pîka kpêkê na", "pîka kpêkê na", "VERB", "Subcat=Tran", "COMPUTER", 5, "mouse-click", "(mouse) click"},
		{"pikakpekeusena", "pîka kpêkê-ûse na", "pîka kpêkê-ûse na", "VERB", "Subcat=Tran", "COMPUTER", 5, "mouse-double-click", "(mouse) double-click"},
		{"pikamaboko", "pîka mabôko", "pîka mabɔ̂kɔ", "VERB", "Subcat=Intr", "ACT", 2, "applaud", "[lit: strike|hand]: applaud"},
		{"pikambeti", "pîka mbëtï", "pîka mbɛ̈tï", "VERB", "Subcat=Intr", "ACT", 2, "type", "[lit: strike|writing]: type"},
		{"pikandembo", "pîka ndembö", "pîka ndembö", "NOUN", "", "ACT", 2, "play-soccer", "[lit: strike|ball]: play soccer"},
		{"pikangasi", "pîka-ngâsî", "pîka-ngâsî", "VERB", "Subcat=Intr", "ACT", 2, "sneeze", "[lit: strike|sneeze]: sneeze"},
		{"pikangombe", "pîka ngombe", "pîka ngombe", "NOUN", "", "ACT", 2, "shoot-a-gun", "[lit: strike|gun]: shoot a gun"},
		{"pikangu", "pîka ngû", "pîka ngû", "NOUN", "", "ACT", 2, "swim", "[lit: strike|water]: swim"},
		{"pikapatara", "pîka patärä", "pîka patärä", "NOUN", "", "ACT", 4, "play-dice", "[lit: strike|dice]: play dice"},
		{"pikapatara", "pîka patärä", "pîka patärä", "VERB", "Subcat=Intr", "INTERACT", 4, "negotiate", "negotiate"},
		{"pikawen", "pîka wên", "pîka wên", "NOUN", "", "ACT", 4, "forge", "[lit: strike|metal]: forge"},
		{"pilipili", "pilipîli", "pilipîli", "NOUN", "", "FOOD", 4, "hot-sauce", "hot sauce"},
		{"pindiri", "pïndïrï", "pïndïrï", "NOUN", "", "OBJ", 4, "charcoal", "charred wood, black"},
		{"pindiritiwa", "pïndïrï tî wâ", "pïndïrï tî wâ", "NOUN", "", "OBJ", 4, "charred-wood", "[lit: charred wood|of|fire]: charred wood used for controlled burning"},
		{"pipi", "pîpï", "pîpï", "NOUN", "", "ANIM", 6, "army-ant", "army ant"},
		{"piri", "pîri", "pîri", "NOUN", "", "OBJ", 5, "mourning-clothes", "mourning clothes"},
		{"pito", "pito", "pito", "NOUN", "", "BODY", 6, "foreskin", "foreskin"},
		{"polele", "polêlê", "polêlê", "ADV", "", "INTERACT", 5, "frankly", "frankly, openly, not mincing words"},
		{"polisi", "polîsi", "polîsi", "NOUN", "", "WHO", 5, "police", "police"},
		{"pome", "pömë", "pɔ̈mɛ̈", "NOUN", "", "FOOD", 6, "apple", "apple"},
		{"pomesitere", "pôme-sitëre", "pɔ̂me-sitɛ̈rɛ", "NOUN", "", "FOOD", 5, "golden-apple", "[Fr: pomme Cythère]: golden apple"},
		{"pometere", "pömëtêre", "pɔ̈mɛ̈tɛ̂rɛ", "NOUN", "", "FOOD", 5, "potato", "potato"},
		{"pongi", "pongi", "pɔngi", "VERB", "Subcat=Intr", "STATE", 3, "relax", "rest, relax, be calm, idle, take a break"},
		{"pongi", "pöngï", "pɔ̈ngï", "NOUN", "", "STATE", 3, "relaxation", "rest, relaxation, calm, idleness"},
		{"pono", "ponö", "pɔnɔ̈", "NOUN", "", "CIVIL", 6, "pain-poverty", "pain, poverty"},
		{"popo", "popô", "popô", "NOUN", "", "BODY", 6, "uncircumcised", "uncircumcised"},
		{"popo", "popö", "popö", "NOUN", "", "TREE", 6, "tattoo-scarification", "tattoo, scarification"},
		{"popo", "pöpö", "pöpö", "NOUN", "", "WHERE", 1, "among", "between, among, inter-"},
		{"poporo", "pôpôrô", "pôpôrô", "NOUN", "", "SICK", 4, "rosaceae", "rosaceae (dermatitis)"},
		{"poro", "pörö", "pɔ̈rɔ̈", "NOUN", "", "BODY", 2, "skin", "skin"},
		{"porole", "pörö-lë", "pɔ̈rɔ̈-lɛ̈", "NOUN", "", "BODY", 3, "pupil", "[lit: skin|eye]: pupil"},
		{"porotigere", "pörö tî gerê", "pɔ̈rɔ̈ tî gɛrɛ̂", "NOUN", "", "OBJ", 2, "shoe", "[lit: skin|of|feet]: shoe"},
		{"porotikeke", "pörö tî këkë", "pɔ̈rɔ̈ tî kɛ̈kɛ̈", "NOUN", "", "OBJ", 2, "bark", "[lit: skin|of|tree]: bark"},
		{"porotimaboko", "pörö tî mabôko", "pɔ̈rɔ̈ tî mabɔ̂kɔ", "NOUN", "", "OBJ", 2, "glove", "[lit: skin|of|hand]: glove"},
		{"poroyanga", "pörö-yângâ", "pɔ̈rɔ̈-yângâ", "NOUN", "", "BODY", 3, "lips", "[lit: skin|mouth]: lips"},
		{"poto", "poto", "pɔtɔ", "VERB", "", "INTERACT", 4, "meddle", "meddle, mix up, screw up"},
		{"potopoto", "potopôto", "pɔtɔpɔ̂tɔ", "NOUN", "", "NATURE", 2, "mud-or-mortar", "mud, mortar, paste"},
		{"pulusu", "pulûsu", "pulûsu", "NOUN", "", "WHO", 5, "police", "police"},
		{"pupu", "pupu", "pupu", "NOUN", "", "HOW", 2, "crumb", "crumb, bit, morsel"},
		{"pupu", "pupu", "pupu", "NOUN", "", "NATURE", 2, "wind", "wind"},
		{"pupulenge", "pûpûlenge", "pûpûlɛngɛ", "NOUN", "", "ANIM", 3, "butterfly", "[lit: wind|pearl?]: butterfly"},
		{"pupulenge", "pûpûlenge", "pûpûlɛngɛ", "NOUN", "", "WHO", 3, "prostitute", "[lit: butterfly]: prostitute"},
		{"pupusese", "pupu-sêse", "pupu-sêse", "NOUN", "", "HOW", 2, "dust", "[lit: crumb|earth]: dust"},
		{"puru", "purû", "purû", "NOUN", "", "BODY", 2, "feces-spoor", "feces, spoor"},
		{"purutingu", "purû tî ngû", "purû tî ngû", "NOUN", "", "BODY", 2, "diarrhea", "[lit: feces|of|water]: diarrhea"},
		{"pusu", "pûsu", "pûsu", "VERB", "Subcat=Tran", "ACT", 3, "push", "push"},
		{"pusupusu", "pûsu-pûsu", "pûsu-pûsu", "NOUN", "", "OBJ", 3, "pushcart", "pushcart"},
		{"pusupusu", "pûsu-pûsu", "pûsu-pûsu", "NOUN", "", "WHO", 3, "pushcart-porter", "pushcart porter"},
		{"saa", "sâa", "sâa", "NOUN", "", "OBJ", 2, "measuring-device", "measuring device"},
		{"saa", "sâa", "sâa", "VERB", "Subcat=Tran", "ACT", 2, "pour-deliver-provoke", "disperse, pour, pour out, pour into; escort, deliver; induce, provoke, initiate, cause, unleash"},
		{"saakiloo", "sâa-kilöo", "sâa-kilöo", "NOUN", "", "OBJ", 2, "weight-scale", "[lit: measure|kilogram]: weight scale"},
		{"saangonga", "sâa-ngonga", "sâa-ngonga", "NOUN", "", "OBJ", 2, "clock-watch", "[lit: measure|hour]: clock, watch"},
		{"saapenda", "sâa pendä", "sâa pendä", "VERB", "Subcat=Intr", "ACT", 2, "influence", "(have an) influence"},
		{"saapete", "sâa-pëtë", "sâa-pɛ̈tɛ̈", "NOUN", "", "OBJ", 2, "barometer", "[lit: measure|pressure]: voltmeter, barometer"},
		{"saato", "sâa to", "sâa to", "VERB", "Subcat=Intr", "ACT", 2, "wage-war", "wage war"},
		{"saawa", "sâa-wâ", "sâa-wâ", "NOUN", "", "OBJ", 2, "thermometer", "[lit: measure|heat]: thermometer"},
		{"saba", "saba", "saba", "NOUN", "", "OBJ", 4, "tongs", "tongs"},
		{"sagba", "sägbä", "sägbä", "NOUN", "", "INTERACT", 6, "rumor", "rumor"},
		{"sai", "sâi", "sâi", "NOUN", "", "FOOD", 5, "tea", "brewing yeast, tea"},
		{"saki", "sâki", "sâki", "ADJ", "NumType=Ord", "NUM", 2, "thousand", "thousand"},
		{"saki", "sâki", "sâki", "NUM", "NumType=Card", "NUM", 2, "thousand", "thousand"},
		{"sakpa", "sakpä", "sakpä", "NOUN", "", "OBJ", 2, "basket", "basket"},
		{"sala", "sâla", "sâla", "VERB", "Subcat=Intr", "ACT", 1, "happen-occur", "happen, occur"},
		{"sala", "sâla", "sâla", "VERB", "Subcat=Tran", "ACT", 1, "do-make", "do, make, spend (time)"},
		{"salaada", "saläada", "saläada", "NOUN", "", "OBJ", 4, "lettuce", "[Fr: salade]: lettuce"},
		{"salana", "sâla na", "sâla na", "VERB", "Subcat=Tran", "ACT", 1, "serve", "serve"},
		{"salanabeoko", "sâla na bê ôko", "sâla na bɛ̂ ɔ̂kɔ", "VERB", "Subcat=Intr", "ACT", 1, "harmonize", "[lit: do|with|heart|one]: act as one, act together, act in harmony"},
		{"salanabeuse", "sâla na bê ûse", "sâla na bɛ̂ ûse", "VERB", "Subcat=Intr", "ACT", 1, "disharmonize", "[lit: do|with|heart|two]: act against each other, act in opposition"},
		{"salayanga", "sâla yângâ", "sâla yângâ", "VERB", "Subcat=Tran", "ACT", 1, "promise", "[lit: do|mouth]: promise"},
		{"samba", "sambâ", "sambâ", "NOUN", "", "FAMILY", 5, "second-wife", "co-spouse, second wife"},
		{"samba", "sâmba", "sâmba", "NOUN", "", "DRINK", 3, "alcohol", "alcohol, any alcoholic beverage"},
		{"sambatibengba", "sâmba tî bengbä", "sâmba tî bengbä", "NOUN", "", "DRINK", 3, "red-wine", "[lit: alcohol|of|red]: red wine"},
		{"sambativuru", "sâmba tî vurü", "sâmba tî vurü", "NOUN", "", "DRINK", 3, "white-wine", "[lit: alcohol|of|white]: white wine"},
		{"sambela", "sambêla", "sambêla", "NOUN", "", "GOD", 4, "prayer", "prayer"},
		{"sambela", "sambêla", "sambêla", "VERB", "", "GOD", 4, "pray", "pray"},
		{"sandaga", "sândâga", "sândâga", "NOUN", "", "GOD", 4, "ritual-feast", "sacrifice, ritual feast"},
		{"sandugu", "sandûgu", "sandûgu", "NOUN", "", "OBJ", 4, "chest-case-trunk-coffin", "chest, case, trunk, coffin"},
		{"sangbi", "sangbi", "sangbi", "NOUN", "", "MOVE", 2, "cross", "cross, throw across, diverge, fork"},
		{"sangbi", "sangbi", "sangbi", "VERB", "Aspect=Imp", "MOVE", 2, "cross", "cross, throw across, diverge, fork"},
		{"sangbilege", "sangbi-lêgë", "sangbi-lêgë", "NOUN", "", "MOVE", 2, "intersection", "[lit: cross|road]: intersection, fork in the road"},
		{"sangbiwa", "sangbi-wâ", "sangbi-wâ", "NOUN", "", "MOVE", 2, "crossfire", "[lit: cross|fire]: crossfire"},
		{"sangi", "sängï", "sängï", "NOUN", "", "NUM", 4, "bunch-of-bananas", "bunch of bananas"},
		{"sangibulee", "sängï bulêe", "sängï bulɛ̂ɛ", "NOUN", "", "NUM", 4, "bunch-of-sweet-bananas", "bunch of sweet bananas"},
		{"sangifondo", "sängï fondo", "sängï fɔndɔ", "NOUN", "", "NUM", 4, "bunch-of-plantains", "bunch of plantains"},
		{"sango", "sango", "sango", "NOUN", "", "INTERACT", 1, "news", "news, tidings"},
		{"sango", "sängö", "sängɔ̈", "NOUN", "", "COUNTRY", 1, "Sango", "Sango"},
		{"santini", "sântînî", "sântînî", "NOUN", "", "ALT WORD FOR", 9, "watchman", "sânzîrî"},
		{"sanze", "sanze", "sanze", "NOUN", "", "PLAY", 4, "thumb-harp", "thumb harp"},
		{"sanziri", "sânzîrî", "sânzîrî", "NOUN", "", "WHO", 4, "watchman", "[En: sentry]: watchman, house guard"},
		{"sara", "särä", "särä", "NOUN", "", "SICK", 5, "mange", "scabies, mange (dermatitis)"},
		{"sarawisi", "sarawîsi", "sarawîsi", "NOUN", "", "CIVIL", 4, "civilized", "civilized, European style"},
		{"sasa", "sasa", "sasa", "VERB", "Subcat=Intr", "BODY", 3, "have-diarrhea", "defecate, have diarrhea"},
		{"sasa", "sasa", "sasa", "VERB", "Subcat=Tran", "BODY", 3, "cause-to-have-diarrhea", "cause to have diarrhea"},
		{"se", "sê", "sê", "VERB", "Subcat=Intr", "STATE", 2, "be-bitter", "be bitter"},
		{"se", "sê", "sɛ̂", "NOUN", "Prefix=Yes", "HOW", 1, "state", "state, manner, proper place, -ence"},
		{"seko", "seko", "seko", "NOUN", "", "ANIM", 5, "chimpanzee", "chimpanzee"},
		{"seleka", "selêka", "selêka", "NOUN", "", "CIVIL", 4, "alliance", "alliance"},
		{"sembe", "sembë", "sɛmbɛ̈", "NOUN", "", "OBJ", 2, "plate", "plate, plateful, disk"},
		{"senda", "sêndâ", "sɛ̂ndâ", "NOUN", "", "HOW", 1, "science", "science, -ology"},
		{"sende", "sëndë", "sɛ̈ndɛ̈", "NOUN", "", "WHERE", 6, "tomb", "tomb, cemetery"},
		{"sene", "sënë", "sɛ̈nɛ̈", "NOUN", "", "SICK", 4, "intestinal-worms", "intestinal worms"},
		{"senge", "sêngê", "sɛ̂ngɛ̂", "ADJ", "", "HOW", 1, "simple-free-okay", "simple, free, empty, just as it is, okay, nude, ordinary, unimportant, without difficulty or opposition"},
		{"senge", "sêngê", "sɛ̂ngɛ̂", "ADV", "", "HOW", 1, "simply-freely-okay", "simply, freely, without difficulty or opposition"},
		{"sepe", "sepë", "sɛpɛ̈", "NOUN", "", "WHEN", 6, "January", "January"},
		{"sepela", "sepela", "sɛpɛla", "VERB", "Subcat=Intr", "CIVIL", 3, "be-on-good-behavior", "be on good behavior"},
		{"sepela", "sepela", "sɛpɛla", "VERB", "Subcat=Tran", "CIVIL", 3, "be-polite-to", "be polite to"},
		{"sepelangozo", "sëpëlängö-zo", "sɛ̈pɛ̈längɔ̈-zo", "NOUN", "", "CIVIL", 3, "politeness-good-manners", "[lit: be polite to|people]: politeness, good manners"},
		{"sere", "serê", "serê", "NOUN", "", "ANIM", 6, "sardine", "sardine"},
		{"sese", "sêse", "sêse", "NOUN", "", "NATURE", 2, "ground-floor-earth-dirt", "ground, floor, earth, dirt"},
		{"sesee", "sesêe", "sesêe", "ADJ", "", "STATE", 2, "bitter", "bitter"},
		{"seta", "sêtâ", "sêtâ", "NOUN", "", "BODY", 4, "bowels-guts-intestines", "bowels, guts, intestines"},
		{"sete", "sëtë", "sɛ̈tɛ̈", "NOUN", "", "OBJ", 4, "ring-nail", "ring, nail"},
		{"setika", "së tî kä", "sɛ̈ tî kä", "NOUN", "", "SICK", 4, "scar", "[lit: state|of|wound]: scar"},
		{"sewa", "sêwâ", "sêwâ", "NOUN", "", "FAMILY", 2, "family", "family"},
		{"si", "sî", "sî", "VERB", "Subcat=Intr", "MOVE", 1, "arrive-finish-have-stood-up", "arrive; finish, end; finish standing up, have risen"},
		{"si", "sï", "sï", "ADV", "", "WHEN", 1, "beforehand-first", "beforehand, first"},
		{"si", "sï", "sï", "SCONJ", "", "WHEN", 1, "before-then-until", "before, then, until"},
		{"si", "sï", "sï", "VERB", "Subcat=Intr", "STATE", 1, "be-full", "be full"},
		{"si", "sï", "sï", "VERB", "Subcat=Tran", "STATE", 1, "fill", "fill"},
		{"sigi", "sïgî", "sïgî", "VERB", "Subcat=Intr", "MOVE", 1, "go-out", "[abbr: sîgïgî] go out, exit"},
		{"sigigi", "sîgïgî", "sîgïgî", "VERB", "Subcat=Intr", "MOVE", 1, "go-out", "[lit: arrive|outside] go out, exit"},
		{"simba", "simba", "simba", "VERB", "Subcat=Intr", "MOVE", 5, "travel", "[ar: 'Simbad'?] travel, voyage, navigate"},
		{"simba", "simbä", "simbä", "NOUN", "", "MOVE", 5, "travel", "[ar: 'Simbad'?] travel, voyage, navigation"},
		{"simisi", "simîsi", "simîsi", "NOUN", "", "SICK", 4, "gonorrhea", "gonorrhea"},
		{"sina", "sî na", "sî na", "VERB", "Subcat=Tran", "STATE", 1, "attain", "end up at, attain"},
		{"sindi", "sindi", "sindi", "NOUN", "", "FOOD", 3, "sesame", "sesame"},
		{"singa", "sînga", "sînga", "NOUN", "", "OBJ", 3, "wire", "wire"},
		{"singa", "sîngâ", "sîngâ", "NOUN", "", "SICK", 4, "fungal-infection", "fungal infection, dry cracking skin"},
		{"singi", "singi", "singi", "NOUN", "", "PLANT", 6, "ginger-plant", "ginger plant"},
		{"singila", "singîla", "singîla", "INTERJ", "", "INTERACT", 1, "Thanks", "Thanks!"},
		{"singila", "singîla", "singîla", "NOUN", "", "INTERACT", 1, "thank", "thank"},
		{"singila", "singîla", "singîla", "VERB", "Subcat=Tran", "INTERACT", 1, "thank", "thank"},
		{"singo", "sïngö", "sïngɔ̈", "VERB", "VerbForm=Vnoun", "MOVE", 1, "arrival", "arrival"},
		{"singola", "sïngö-lâ", "sïngɔ̈-lâ", "NOUN", "", "WHEN", 2, "sunrise", "sunrise"},
		{"sioba", "sïö-bä", "sïɔ̈-bä", "NOUN", "", "INTERACT", 1, "malediction", "malediction"},
		{"siodeba", "sïö-dëbä", "sïɔ̈-dëbä", "NOUN", "", "INTERACT", 3, "curse", "curse"},
		{"siokpale", "sïö-kpälë", "sïɔ̈-kpälë", "NOUN", "", "GOD", 6, "sin", "sin"},
		{"siokpari", "sïö-kpärï", "sïɔ̈-kpärï", "NOUN", "", "GOD", 6, "sin", "sin"},
		{"sioni", "sïönî", "sïɔ̈nî", "ADJ", "", "HOW", 2, "bad", "bad"},
		{"sioni", "sïönî", "sïɔ̈nî", "ADV", "", "HOW", 1, "badly", "badly"},
		{"sioni", "sïönî", "sïɔ̈nî", "NOUN", "", "HOW", 3, "badness", "badness, evil"},
		{"sionipere", "sïönî pêrë", "sïɔ̈nî pêrë", "NOUN", "", "PLANT", 1, "underbrush", "[lit: bad|grass]: weed, underbrush"},
		{"sioye", "sïö-yê", "sïɔ̈-yê", "NOUN", "", "GOD", 3, "sin", "sin"},
		{"siozo", "sïö-zo", "sïɔ̈-zo", "NOUN", "", "GOD", 3, "wicked-person", "wicked-person"},
		{"siri", "siri", "siri", "NOUN", "", "ANIM", 4, "flea-lice", "flea, lice"},
		{"siriri", "sîrîrî", "sîrîrî", "ADV", "", "STATE", 3, "peaceful", "peaceful, calm, tranquil"},
		{"siriri", "sîrîrî", "sîrîrî", "NOUN", "", "STATE", 3, "peace", "peace, calm, tranquillity"},
		{"sisa", "sisa", "sisa", "NOUN", "", "BODY", 4, "tendon-nerve-vein", "tendon, nerve, vein"},
		{"sisi", "sisi", "sisi", "NOUN", "", "OBJ", 6, "needle", "needle"},
		{"so", "so", "so", "VERB", "Subcat=Tran", "FEEL", 2, "afflict-torment-make-suffer", "afflict, torment, make suffer"},
		{"so", "sô", "sô", "DET", "PronType=Art|PronType=Rel", "WHICH", 1, "this-these", "this, these"},
		{"so", "sô", "sô", "SCONJ", "", "INTERACT", 1, "that", "that"},
		{"so", "sô", "sô", "VERB", "Subcat=Intr", "ACT", 3, "be-saved", "be-saved"},
		{"so", "sô", "sô", "VERB", "Subcat=Tran", "ACT", 3, "save", "save"},
		{"so", "sô", "sɔ̂", "VERB", "Subcat=Tran", "ACT", 2, "strike", "strike"},
		{"sobenda", "sô benda", "sɔ̂ bɛnda", "VERB", "Subcat=Intr", "ACT", 6, "win", "win, be victorious"},
		{"sombee", "sombêe", "sombêe", "VERB", "Subcat=Tran", "ACT", 6, "accumulate", "amass, accumulate, rack up"},
		{"sombere", "sombere", "sombere", "NOUN", "", "OBJ", 4, "barb", "barb"},
		{"somoye", "sô mo yê", "sô mo yê", "ADV", "", "INTERACT", 1, "-ever", "[this|you|want]: -ever"},
		{"somvenisi", "sô mvenî sï", "sô mvɛnî sï", "SCONJ", "", "INTERACT", 1, "it's-just-what", "[this|self|then]: it's just what"},
		{"son", "sôn", "sôn", "NOUN", "", "OBJ", 5, "rat-trap", "rat trap"},
		{"son", "sôn", "sôn", "VERB", "Subcat=Tran", "OBJ", 5, "light", "light (a source of illumination)"},
		{"songo", "songö", "sɔngɔ̈", "VERB", "VerbForm=Vnoun", "FEEL", 2, "pain", "pain"},
		{"songo", "söngö", "söngö", "NOUN", "", "FAMILY", 2, "family-member", "family member"},
		{"songo", "söngö", "söngö", "NOUN", "", "FEEL", 2, "filial-love", "familial affection, filial love"},
		{"songobe", "songö-bê", "sɔngɔ̈-bɛ̂", "NOUN", "", "FEEL", 2, "mental-anguish", "mental anguish"},
		{"songosongo", "songosongo", "songosongo", "NOUN", "", "PLANT", 6, "elephant-grass", "elephant grass"},
		{"soro", "soro", "soro", "VERB", "Subcat=Tran", "ACT", 5, "choose", "choose"},
		{"soro", "sorö", "sorö", "NOUN", "", "ACT", 5, "choice", "choice"},
		{"soro", "sörö", "sɔ̈rɔ̈", "NOUN", "", "BODY", 4, "rapids-foam-scum", "drool, river rapids, foam, scum"},
		{"soronga", "soronga", "soronga", "VERB", "Aspect=Iter|Subcat=Tran", "ACT", 4, "differentiate", "differentiate"},
		{"soso", "soso", "sɔsɔ", "NOUN", "", "ACT", 2, "pound-in-a-mortar", "pound in a mortar"},
		{"soso", "sösö", "sösö", "NOUN", "", "BODY", 5, "fart", "fart"},
		{"su", "su", "su", "VERB", "Subcat=Intr", "NATURE", 2, "be-brilliantly-lit-up", "be brilliantly lit up"},
		{"su", "su", "su", "VERB", "Subcat=Tran", "ACT", 2, "suck-or-lick", "suck, lick"},
		{"su", "sû", "sû", "VERB", "Subcat=Tran", "INTERACT", 2, "draw-design-trace", "draw, design, trace"},
		{"sua", "sua", "sua", "VERB", "Subcat=Intr", "NATURE", 3, "flow", "flow"},
		{"sua", "sua", "sua", "VERB", "Subcat=Tran", "BODY", 3, "comb", "comb"},
		{"sua", "süä", "süä", "NOUN", "", "OBJ", 3, "needle", "needle"},
		{"suali", "suali", "suali", "NOUN", "Subcat=Tran", "BODY", 3, "comb-or-brush", "[lit: comb|hair]: comb, brush"},
		{"suali", "süäli", "süäli", "VERB", "Subcat=Tran", "BODY", 3, "hairpin", "[lit: needle|hair]: hair pin"},
		{"sui", "sûî", "sûî", "NUM", "NumType=Frac|Prefix=Yes", "NUM", 2, "deca-", "deca-"},
		{"suiya", "suïya", "suïya", "NOUN", "", "OBJ", 3, "skewer", "skewer"},
		{"sukani", "sukâni", "sukâni", "NOUN", "", "FOOD", 3, "sugar", "sugar"},
		{"suku", "sûku", "sûku", "VERB", "Subcat=Intr", "ACT", 3, "be-inflated", "be inflated"},
		{"suku", "sûku", "sûku", "VERB", "Subcat=Tran", "ACT", 3, "inflate", "inflate"},
		{"sukula", "sukûla", "sukûla", "VERB", "Subcat=Tran", "ACT", 2, "washclean", "wash,clean"},
		{"sukulabe", "sukûla bê", "sukûla bɛ̂", "VERB", "Subcat=Intr", "GOD", 2, "forgive-one's-sins", "forgive one's sins"},
		{"sukulangongu", "sükülängö ngû", "sükülängɔ̈ ngû", "NOUN", "", "ACT", 2, "bath", "bath"},
		{"sukulangu", "sukûla ngû", "sukûla ngû", "VERB", "Subcat=Intr", "ACT", 2, "bathe", "bathe"},
		{"sukulu", "sukûlu", "sukûlu", "NOUN", "", "ANIM", 4, "owl", "owl"},
		{"sulee", "sulëe", "sulëe", "VERB", "Subcat=Intr", "STATE", 4, "be-drunk", "[Fr: soul]: be drunk"},
		{"suma", "suma", "suma", "VERB", "Subcat=Tran", "FEEL", 3, "dream", "dream"},
		{"suma", "sümä", "sümä", "NOUN", "", "FEEL", 3, "dream", "dream"},
		{"sumasuma", "suma süma", "suma süma", "VERB", "Subcat=Intr", "FEEL", 3, "dream", "dream"},
		{"sumbeti", "sû-mbëtï", "sû-mbɛ̈tï", "VERB", "Subcat=Intr", "INTERACT", 2, "write", "[lit: trace|writing]: write"},
		{"sumbu", "sumbu", "sumbu", "NOUN", "", "ANIM", 5, "gorilla", "gorilla"},
		{"sungba", "sungba", "sungba", "VERB", "Subcat=Intr", "NATURE", 3, "explode-thunder-blossom", "explode, thunder, blossom"},
		{"sungba", "sungbä", "sungbä", "NOUN", "", "NATURE", 3, "debris-from-explosion", "debris from explosion"},
		{"sungbango", "süngbängö", "süngbängɔ̈", "VERB", "VerbForm=Vnoun", "NATURE", 3, "explosion", "explosion"},
		{"sungombeti", "süngö-mbëtï", "süngɔ̈-mbɛ̈tï", "VERB", "VerbForm=Vnoun", "INTERACT", 6, "writing", "[lit: tracing|writing]: writing"},
		{"supu", "sûpu", "sûpu", "NOUN", "", "FOOD", 4, "soup", "[Fr: soupe]: soup"},
		{"sura", "sura", "sura", "VERB", "Subcat=Tran", "INTERACT", 5, "cut-divide-partition", "cut, divide, partition"},
		{"sura", "surä", "surä", "NOUN", "", "INTERACT", 5, "section", "section"},
		{"suru", "sûru", "sûru", "VERB", "Subcat=Intr", "ACT", 2, "tear-rip-split", "tear, rip, split"},
		{"suru", "sûru", "sûru", "VERB", "Subcat=Tran", "ACT", 2, "tear-rip-split", "tear, rip, split, pluck, cut, sharpen"},
		{"susu", "susu", "susu", "NOUN", "", "ANIM", 2, "fish", "fish"},
		{"ta", "ta", "ta", "NOUN", "", "PLANT", 3, "gourd-or-kettle", "gourd, pumpkin; pot, kettle"},
		{"taa", "taâ", "taâ", "ADJ", "", "INTERACT", 1, "true-real-authentic", "true, real, authentic"},
		{"taapande", "taä-pandë", "taä-pandë", "NOUN", "", "HOW", 2, "paradigm", "[lit: true-example]: paradigm, archetype, best example"},
		{"taasango", "taä-sängö", "taä-sängɔ̈", "NOUN", "", "COUNTRY", 1, "ethnic-Sango", "[lit: true|Sango]: ethnic Sango, Ngbandi"},
		{"taatene", "taä-tene", "taä-tɛnɛ", "NOUN", "", "INTERACT", 1, "truth", "[lit: true|story]: truth"},
		{"taba", "taba", "taba", "NOUN", "", "ANIM", 2, "sheep", "sheep"},
		{"tagba", "tâgba", "tâgba", "NOUN", "", "ANIM", 6, "antilope", "medium-sized antilope"},
		{"taka", "tâkâ", "tâkâ", "ADJ", "Prefix=Yes", "STATE", 6, "original", "original"},
		{"takasa", "ta-kâsa", "ta-kâsa", "NOUN", "", "PLANT", 3, "saucepan", "[lit: pot|sauce]: saucepan"},
		{"taliti", "tâ-li-tï", "tâ-li-tï", "NOUN", "", "BODY", 4, "thumb", "[lit: mother|finger]: thumb"},
		{"tambula", "tambûla", "tambûla", "NOUN", "", "MOVE", 2, "walk", "walk, promenade, travel"},
		{"tambula", "tambûla", "tambûla", "VERB", "Subcat=Intr", "MOVE", 2, "walk", "walk, walk around, promenade; work, function"},
		{"tanga", "tanga", "tanga", "NOUN", "", "NUM", 1, "leftovers", "rest, remaining portion, leftovers"},
		{"tangbi", "tangbi", "tangbi", "NOUN", "", "INTERACT", 4, "connection", "connection"},
		{"tangbi", "tangbi", "tangbi", "VERB", "Aspect=Imp", "INTERACT", 4, "connect", "connect"},
		{"tangbo", "tâ-ngbö", "tâ-ngbɔ̈", "NOUN", "", "FAMILY", 4, "mother-of-twins", "mother of twins"},
		{"tangbo", "tä-ngbö", "tä-ngbɔ̈", "NOUN", "Gender=Fem", "FAMILY", 4, "mother-of-twins", "mother of twins"},
		{"tange", "tangë", "tangë", "NOUN", "", "HOUSE", 4, "traditional-bed", "traditional bed made of wood"},
		{"tango", "tângo", "tângo", "NOUN", "", "WHEN", 2, "time-era-epoch", "time, era, epoch"},
		{"tangu", "ta-ngû", "ta-ngû", "NOUN", "", "OBJ", 3, "water-barrel", "[lit: pot|water]: water barrel"},
		{"tapare", "täpärë", "täpärë", "NOUN", "", "INTERACT", 4, "argument", "argument, dispute, quarrel"},
		{"tara", "tara", "tara", "VERB", "Subcat=Tran", "FEEL", 1, "try-taste", "try, try on, try out, taste, attempt; tempt, seduce"},
		{"tara", "tarä", "tarä", "NOUN", "", "FAMILY", 4, "paternal-grandrelative", "paternal grandchild; paternal grandmother"},
		{"tasese", "ta-sêse", "ta-sêse", "NOUN", "", "OBJ", 3, "crock", "[lit: pot|earth]: crock"},
		{"tatalita", "tatalîta", "tatalîta", "NOUN", "", "SENSE", 5, "bugle", "[onomatopoeia]: bugle"},
		{"tatara", "tatärä", "tatärä", "NOUN", "", "OBJ", 3, "glass-mirror", "glass, mirror"},
		{"tatarale", "tatärä-lê", "tatärä-lɛ̂", "NOUN", "", "OBJ", 3, "glasses", "[lit: glass|eye]: glasses"},
		{"tatarando", "tatara ndo", "tatara ndo", "NOUN", "", "OBJ", 6, "feel-one's-way", "[lit: touch|place]: feel-one's-way"},
		{"taza", "taza", "taza", "NOUN", "", "PLANT", 5, "reed", "reed"},
		{"te", "te", "tɛ", "VERB", "", "ACT", 1, "eat-bite-gnaw", "eat, bite, gnaw"},
		{"tekiri", "te-kîri", "tɛ-kîri", "NOUN", "", "WHEN", 6, "November", "November"},
		{"teloti", "te-lötï", "tɛ-lötï", "NOUN", "", "WHEN", 6, "October", "October"},
		{"tembe", "tembe", "tembe", "NOUN", "", "INTERACT", 4, "concurrence-or-rivalry", "concurrence, emulation, rivalry"},
		{"tende", "tende", "tende", "NOUN", "", "PLANT", 4, "cotton", "cotton"},
		{"tene", "tene", "tɛnɛ", "VERB", "Subcat=Tran", "INTERACT", 1, "say-tell", "say, tell"},
		{"tene", "tênë", "tɛ̂nɛ̈", "NOUN", "", "NATURE", 3, "rock-stone", "rock, stone, gravel, pebble"},
		{"tene", "tënë", "tɛ̈nɛ̈", "NOUN", "", "INTERACT", 1, "speech-issue-problem-argument", "problem, quarrel; speech, talk, words, tale"},
		{"tenemvene", "tene mvene", "tɛnɛ mvɛnɛ", "VERB", "Subcat=Intr", "INTERACT", 1, "tell-a-lie", "tell a lie"},
		{"tenengotene", "tënëngö tënë", "tɛ̈nɛ̈ngɔ̈ tɛ̈nɛ̈", "VERB", "Subcat=Intr", "INTERACT", 1, "speaking-talking", "speaking, talking"},
		{"tenetaatene", "tene taä-tënë", "tɛnɛ taä-tɛ̈nɛ̈", "VERB", "Subcat=Intr", "INTERACT", 1, "tell-the-truth", "tell the truth"},
		{"tenetene", "tene tënë", "tɛnɛ tɛ̈nɛ̈", "VERB", "Subcat=Intr", "INTERACT", 1, "speak-talk", "speak, talk"},
		{"teneti", "tënë tî", "tɛ̈nɛ̈ tî", "ADP", "", "HOW", 1, "because-of", "[lit: account|of]: because of"},
		{"tenetinye", "tënë tî nye", "tɛ̈nɛ̈ tî nyɛ", "ADV", "", "HOW", 1, "why", "[lit: account|of|what]: why"},
		{"tenetinzapa", "tënë tî nzapä", "tɛ̈nɛ̈ tî nzapä", "NOUN", "", "GOD", 1, "gospel", "[lit: word|of|God]: gospel, gospel truth"},
		{"tenetiso", "tënë tî sô", "tɛ̈nɛ̈ tî sô", "ADV", "", "HOW", 1, "because-of-this", "[lit: account|of|this]: because of this"},
		{"tenetiso", "tënë tî sô", "tɛ̈nɛ̈ tî sô", "CCONJ", "", "HOW", 1, "because", "[lit: account|of|that]: because"},
		{"tenga", "tênga", "tɛ̂nga", "VERB", "Aspect=Iter|Subcat=Tran", "CIVIL", 3, "bring-together", "organize a meeting, bring together"},
		{"tengawa", "tênga-wâ", "tɛ̂nga-wâ", "NOUN", "", "OBJ", 3, "match-lighter", "[lit: bring together|fire]: match, lighter"},
		{"tengbi", "têngbi", "tɛ̂ngbi", "NOUN", "", "CIVIL", 3, "meeting-interview", "junction, meeting, interview"},
		{"tengbi", "têngbi", "tɛ̂ngbi", "VERB", "Aspect=Imp|Subcat=Tran", "CIVIL", 3, "join-meet", "join, rejoin, meet, meet up"},
		{"tere", "tere", "tɛrɛ", "NOUN", "", "ANIM", 4, "spider", "spider"},
		{"tere", "tere", "tɛrɛ", "NOUN", "", "MYTH", 4, "Spider", "Spider, trickster-hero in many fables"},
		{"tere", "terê", "tɛrɛ̂", "NOUN", "", "BODY", 1, "body", "body, tree trunk, (house) walls"},
		{"tere", "terê", "tɛrɛ̂", "NOUN", "", "WHERE", 1, "next-to", "surroundings, next to"},
		{"tere", "terê", "tɛrɛ̂", "NOUN", "Reflex=Yes", "WHO", 1, "each-other", "each other, oneself"},
		{"ti", "tî", "tî", "ADP", "", "STATE", 1, "of-or-to", "[+noun => adj]: of, from, pertaining to; [+verb => infinitive]: to"},
		{"ti", "tï", "tï", "NOUN", "", "BODY", 4, "hand-arm", "hand, arm"},
		{"ti", "tï", "tï", "VERB", "Subcat=Intr", "MOVE", 1, "fall", "fall, (sun) set"},
		{"tia", "tîa", "tîa", "VERB", "Subcat=Tran", "NUM", 3, "be-insufficient-for", "be lacking for, missing in, insufficient for"},
		{"tiaa", "tîâa", "tîâa", "VERB", "Subcat=Tran", "ALT SP FOR", 9, "be-insufficient-for", "tîa"},
		{"tiki", "tikî", "tikî", "NOUN", "", "PLANT", 6, "cotton", "cotton"},
		{"tikisa", "tikîsa", "tikîsa", "VERB", "Subcat=Tran", "INTERACT", 5, "betray", "betray"},
		{"tiko", "tîko", "tîkɔ", "VERB", "Subcat=Intr", "SICK", 6, "cough", "cough"},
		{"tiko", "tîkö", "tîkɔ̈", "NOUN", "", "SICK", 6, "cough", "cold, cough"},
		{"tindani", "tî-ndâ-nî", "tî-ndâ-nî", "ADV", "", "WHEN", 1, "lastly", "lastly, finally"},
		{"tipoi", "tipôi", "tipôi", "NOUN", "", "OBJ", 6, "sedan-chair", "sedan chair"},
		{"tiri", "tiri", "tiri", "VERB", "Subcat=Intr", "INTERACT", 4, "battle", "battle"},
		{"tirika", "tirika", "tirika", "VERB", "Subcat=Intr", "INTERACT", 4, "struggle", "struggle"},
		{"tiringbi", "tîrîngbi", "tîrîngbi", "VERB", "Aspect=Imp|Subcat=Intr", "INTERACT", 4, "do-hand-to-hand-combat", "do hand-to-hand combat"},
		{"tisa", "tisa", "tisa", "VERB", "", "INTERACT", 5, "invite-advise", "invite, notify, advise"},
		{"tisa", "tisä", "tisä", "NOUN", "", "INTERACT", 5, "invitation-advice", "invitation, notification, advice"},
		{"titene", "tî-tene", "tî-tɛnɛ", "CCONJ", "", "INTERACT", 1, "i.e.", "[lit: to|say]: which is to say..."},
		{"to", "to", "to", "VERB", "Subcat=Tran", "ACT", 2, "assign-to-post", "send, assign (to a post)"},
		{"to", "tö", "tö", "VERB", "Subcat=Tran", "ACT", 2, "convey-or-transport", "convey, transport (e.g. by cart); draw (e.g. water from a well, grain from storage); hit a target"},
		{"to", "tö", "tɔ̈", "NOUN", "", "WHERE", 2, "east-or-upstream", "east; upstream"},
		{"tobua", "tö-buä", "tö-buä", "NOUN", "", "GOD ", 3, "pope", "[lit: father|priest]: pope"},
		{"tokua", "to-kua", "to-kua", "VERB", "Subcat=Intr", "INTERACT", 2, "send-a-message", "[lit: send|work]: send a message, report"},
		{"tokua", "tokua", "tokua", "NOUN", "", "INTERACT", 2, "report", "[lit: send|work]: message, report, mail"},
		{"toli", "toli", "toli", "NOUN", "", "MYTH", 6, "fable", "parable, fable, legend"},
		{"toli", "tolï", "tolï", "NOUN", "", "INTERACT", 4, "advice", "advice, counsel"},
		{"toliti", "tö-li-tï", "tö-li-tï", "NOUN", "", "BODY", 4, "index-finger", "[lit: father|finger]: index finger"},
		{"tolo", "tôlo", "tôlo", "NOUN", "", "HOUSE", 6, "metal-roof", "[Fr: tôle]: metal roof"},
		{"tomati", "tomâti", "tomâti", "NOUN", "", "FOOD", 4, "tomato", "[Fr: tomate]: tomato"},
		{"tomboka", "tombôka", "tombôka", "VERB", "Subcat=Intr", "FEEL", 5, "go-crazy", "be troubled, go crazy"},
		{"tonda", "töndâ", "töndâ", "VERB", "Subcat=Intr", "STATE", 2, "start", "start, begin"},
		{"tondo", "tondo", "tɔndɔ", "VERB", "Subcat=Intr", "ACT", 3, "check", "control, verify, check"},
		{"tondo", "töndö", "töndö", "NOUN", "", "FOOD", 5, "ginger-root", "ginger root"},
		{"tondongo", "töndöngö", "tɔ̈ndɔ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 3, "check", "control, verification, check"},
		{"tonga", "tonga", "tonga", "NOUN", "", "OBJ", 6, "needle", "needle"},
		{"tongana", "töngana", "töngana", "ADP", "", "HOW", 1, "as", "[lit: as*|to] as, like"},
		{"tongana", "töngana", "töngana", "SCONJ", "", "HOW", 1, "if", "if, while"},
		{"tongananye", "töngana nye", "töngana nyɛ", "ADV", "", "HOW", 1, "how", "how"},
		{"tonganatiso", "töngana tî sô", "töngana tî sô", "CCONJ", "", "HOW", 1, "for-example", "for example"},
		{"tongaso", "töngasô", "töngasô", "ADV", "", "HOW", 1, "thus", "thus, like this, like that"},
		{"tongaso", "töngasô", "töngasô", "CCONJ", "", "HOW", 1, "thus", "[lit: as*|so] in this/that case, like this/that"},
		{"tongbi", "tôngbi", "tôngbi", "VERB", "Aspect=Imp|Subcat=Tran", "ACT", 2, "exchange", "exchange, permute, alternate"},
		{"tongbi", "tôngbï", "tôngbï", "NOUN", "", "ACT", 2, "exchange", "exchange, permutation"},
		{"tongbo", "tô-ngbö", "tô-ngbɔ̈", "NOUN", "Gender=Masc", "FAMILY", 4, "father-of-twins", "father of twins"},
		{"tongo", "tongo", "tongo", "NOUN", "", "NATURE", 4, "metallic-lead", "[metal] lead"},
		{"tongo", "töngö", "töngɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 2, "drawing-water", "drawing (e.g. water from a well, grain from storage)"},
		{"tongo", "töngö", "tɔ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 2, "cooking-or-boiling", "cooking, boiling"},
		{"tongolo", "tongolo", "tongolo", "NOUN", "", "NATURE", 3, "star", "star"},
		{"tongoro", "tongoro", "tongoro", "NOUN", "", "ALT SP FOR", 9, "star", "tongolo"},
		{"tono", "tono", "tɔnɔ", "VERB", "Subcat=Intr", "ACT", 3, "drain-or-drip", "drain, drip"},
		{"too", "tôo", "tɔ̂ɔ", "VERB", "Subcat=Tran", "ACT", 2, "cook-or-boil", "cook, boil"},
		{"toro", "törö", "tɔ̈rɔ̈", "NOUN", "", "MYTH", 4, "spirit", "spirit, soul of the deceased, ghost, phantom"},
		{"tororo", "torôrô", "tɔrɔ̂rɔ̂", "ADV", "", "HOW", 4, "repeatedly", "repeatedly"},
		{"toto", "toto", "toto", "NOUN", "", "FEEL", 2, "sound-noise", "sound, noise, crying"},
		{"toto", "toto", "toto", "VERB", "Subcat=Intr", "FEEL", 2, "make-a-sound", "make a sound, make a noise, cry"},
		{"toto", "toto", "tɔtɔ", "NOUN", "", "ANIM", 6, "mongoose", "mongoose"},
		{"tukia", "tukîa", "tukîa", "NOUN", "", "NUM", 5, "are", "are=100 sq meters=~1000 sq ft"},
		{"tukia", "tukîa", "tukîa", "NOUN", "", "PLANT", 6, "row-of-cotton", "row of cotton"},
		{"tuku", "tûku", "tûku", "NOUN", "", "OBJ", 2, "barrel", "[Fr: touque]: barrel"},
		{"tuku", "tûku", "tûku", "VERB", "Subcat=Intr", "ACT", 2, "flip-over-capsize", "flip over, capsize"},
		{"tuku", "tûku", "tûku", "VERB", "Subcat=Tran", "ACT", 2, "throw-out", "throw out, pour, disperse, dissipate"},
		{"tukumolenge", "tûku môlengê", "tûku môlɛngɛ̂", "VERB", "Subcat=Intr", "ACT", 6, "abort-a-pregnancy", "[lit: throw out|child]: abort a pregnancy"},
		{"tumba", "tumba", "tumba", "NOUN", "", "ACT", 2, "fight", "fight, combat, battle, war, offensive"},
		{"tumba", "tumba", "tumba", "VERB", "Subcat=Tran", "ACT", 2, "hunt-or-chase-away", "hunt, pursue, chase away, dismiss"},
		{"tungu", "tungu", "tungu", "NOUN", "", "NATURE", 5, "iron-tin-aluminum", "iron, cast iron, lead, tin, aluminum"},
		{"turu", "turu", "turu", "NOUN", "", "ACT", 4, "forge", "forge"},
		{"turugu", "turûgu", "turûgu", "NOUN", "", "WHO", 4, "soldier", "[Fr: turc]: soldier"},
		{"turungu", "tûrûngu", "tûrûngu", "NOUN", "", "BODY", 5, "navel", "navel"},
		{"tutu", "tutü", "tutü", "VERB", "Subcat=Intr", "COLOR", 3, "be-blue", "be blue"},
		{"uru", "uru", "uru", "VERB", "Subcat=Intr", "MOVE", 2, "jump", "jump, leap, fly away"},
		{"uru", "uru", "uru", "VERB", "Subcat=Tran", "MOVE", 2, "jump-over", "jump over, mount"},
		{"uru", "ûru", "ûru", "VERB", "Subcat=Intr", "MOVE", 2, "blow", "blow, play (wind instrument)"},
		{"urulu", "ûrûlû", "ûrûlû", "NOUN", "", "INTERACT", 4, "quarrel", "argument, dispute, quarrel"},
		{"use", "ûse", "ûse", "ADJ", "NumType=Ord", "NUM", 2, "two", "two"},
		{"use", "ûse", "ûse", "NUM", "NumType=Card", "NUM", 2, "two", "two"},
		{"usio", "usïö", "usïö", "ADJ", "NumType=Ord", "NUM", 2, "four", "four"},
		{"usio", "usïö", "usïö", "NUM", "NumType=Card", "NUM", 2, "four", "four"},
		{"va", "va", "va", "NOUN", "", "WHO", 2, "servant-or-disciple", "servant, disciple"},
		{"vaka", "vaka", "vaka", "NOUN", "", "CIVIL", 4, "district", "district of a town, housing tract"},
		{"vara", "vârä", "vârä", "NOUN", "", "OBJ", 5, "shield", "shield"},
		{"vatanda", "vâtândâ", "vâtândâ", "ADV", "", "ACT", 5, "gulp-down", "(drink down) in one gulp"},
		{"veke", "vekë", "vɛkɛ̈", "NOUN", "", "FOOD", 4, "okra", "okra"},
		{"vii", "vîi", "vîi", "ADV", "", "INTERACT", 5, "frankly", "frankly, openly, not mincing words"},
		{"vo", "vo", "vo", "VERB", "", "ACT", 4, "eat-with-one's-fingers", "eat with one's fingers"},
		{"vo", "vo", "vɔ", "VERB", "Subcat=Tran", "ACT", 2, "buy", "buy"},
		{"vongba", "vongbâ", "vongbâ", "NOUN", "", "ANIM", 4, "warthog", "warthog"},
		{"vongere", "vo ngêrë", "vɔ ngêrë", "VERB", "Subcat=Intr", "ACT", 2, "shop", "[lit: buy|commerce]: shop"},
		{"vongo", "vöngö", "vɔ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 2, "shopping", "shopping"},
		{"voro", "voro", "vɔrɔ", "VERB", "Subcat=Tran", "INTERACT", 2, "worship", "beg, implore, pray to, adore, worship"},
		{"vorongo", "vöröngö", "vɔ̈rɔ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "INTERACT", 2, "worship", "begging, praying, adoration, worship"},
		{"vorotere", "voro terê", "vɔrɔ tɛrɛ̂", "VERB", "Subcat=Intr", "INTERACT", 2, "beg-forgiveness", "[lit: beg|for oneself]: ask for pardon, beg forgiveness"},
		{"voto", "vo to ", "vɔ to ", "VERB", "Subcat=Intr", "ACT", 2, "lodge-a-complaint", "[lit: buy|tears]: lodge a complaint"},
		{"vovongo", "vo vöngö ", "vɔ vɔ̈ngɔ̈ ", "VERB", "Subcat=Intr", "ACT", 2, "go-shopping", "[lit: buy|buying]: go shopping"},
		{"vovoro", "vovoro", "vɔvɔrɔ", "NOUN", "", "BODY", 4, "spine", "spine, lower back"},
		{"vu", "vü", "vü", "VERB", "Subcat=Intr", "ACT", 5, "be-propagated-as-darkness", "become dirty, become dark; become famous, be spread (news), get around"},
		{"vu", "vü", "vü", "VERB", "Subcat=Tran", "ACT", 4, "propagate-darkness", "make dirty, darken; make famous, spread (news), propagate"},
		{"vuko", "vukö", "vukɔ̈", "ADJ", "", "COLOR", 3, "black", "black"},
		{"vuko", "vûkö", "vûkɔ̈", "VERB", "Subcat=Intr", "COLOR", 3, "be-black", "be black"},
		{"vukokete", "vukö-kêtê", "vukɔ̈-kɛ̂tɛ̂", "ADJ", "", "COLOR", 3, "dark-gray", "dark gray"},
		{"vukokete", "vûkö kêtê", "vûkɔ̈ kɛ̂tɛ̂", "VERB", "Subcat=Intr", "COLOR", 3, "be-dark-gray", "be dark gray"},
		{"vukole", "vûko-lë", "vûkɔ-lɛ̈", "NOUN", "", "BODY", 3, "iris", "[lit: black|eye]: iris"},
		{"vukombunzu", "vukö mbunzû", "vukɔ̈ mbunzû", "NOUN", "", "WHO", 1, "black-foreigner", "black foreigner, African American"},
		{"vukomingi", "vûkö mîngi", "vûkɔ̈ mîngi", "VERB", "Subcat=Intr", "COLOR", 3, "be-jet-black", "be jet black"},
		{"vukovuko", "vukö-vukö", "vukɔ̈-vukɔ̈", "ADJ", "", "COLOR", 3, "jet-black", "jet black"},
		{"vuma", "vümä", "vümä", "NOUN", "", "ANIM", 3, "fly", "fly"},
		{"vundu", "vundû", "vundû", "NOUN", "", "FEEL", 2, "pity-resentment", "pity, chagrin, bitterness, resentment"},
		{"vunga", "vünga", "vünga", "VERB", "Aspect=Iter|Subcat=Tran", "ACT", 2, "announce-publish", "announce, publish"},
		{"vuru", "vuru", "vuru", "VERB", "Subcat=Intr", "COLOR", 3, "be-white", "be white"},
		{"vuru", "vurü", "vurü", "ADJ", "", "COLOR", 3, "white", "white"},
		{"vuru", "vurü", "vurü", "NOUN", "", "HOW", 3, "dry", "dry"},
		{"vurukete", "vuru kêtê", "vuru kɛ̂tɛ̂", "VERB", "Subcat=Intr", "COLOR", 3, "be-light-gray", "be light gray"},
		{"vurukete", "vurü-kêtê", "vurü-kɛ̂tɛ̂", "ADJ", "", "COLOR", 3, "light-gray", "light gray"},
		{"vurumingi", "vuru mîngi", "vuru mîngi", "VERB", "Subcat=Intr", "COLOR", 3, "be-bright-white", "be bright white"},
		{"vuruvuru", "vurü-vurü", "vurü-vurü", "ADJ", "", "COLOR", 3, "bright-white", "bright white"},
		{"wa", "wa", "wa", "DET", "PronType=Int|Suffix=Yes", "WHICH", 1, "which", "[following noun, esp. la, zo, ndo]: which"},
		{"wa", "wa", "wa", "NOUN", "Prefix=Yes", "WHO", 1, "owner-or-agent", "proprietor, agent, one in charge, master"},
		{"wa", "wa", "wa", "VERB", "Subcat=Tran", "INTERACT", 3, "advise-or-blame", "advise, counsel; blame, reprimand"},
		{"wa", "wâ", "wâ", "NOUN", "", "NATURE", 2, "fire", "fire, flame, heat, light"},
		{"waawa", "waäwa", "waäwa", "ADJ", "", "STATE", 5, "vague", "vague, indefinite, imprecise"},
		{"waawa", "waäwa", "waäwa", "ADV", "", "STATE", 5, "pell-mell", "pell-mell"},
		{"wabe", "wa-bê", "wa-bɛ̂", "VERB", "Subcat=Intr", "INTERACT", 3, "have-regrets", "have regrets"},
		{"wabindi", "wa-bindi", "wa-bindi", "NOUN", "", "STATE", 5, "sorcerer", "[lit: one who|magic]: magician, sorcerer, mage"},
		{"wafangokua", "wa-fängö-kua", "wa-fängɔ̈-kua", "NOUN", "", "WHO", 3, "trainer", "[lit: one who|teaching|work]: trainer"},
		{"wafangombeti", "wa-fängö-mbëtï", "wa-fängɔ̈-mbɛ̈tï", "NOUN", "", "WHO", 3, "primary-school-teacher", "[lit: one who|teaching|writing]: primary school teacher"},
		{"wafangoye", "wa-fängö-yê", "wa-fängɔ̈-yê", "NOUN", "", "WHO", 3, "teacher", "[lit: one who|teaching|thing]: teacher"},
		{"wahanda", "wa-hânda", "wa-hânda", "NOUN", "", "GOD", 4, "trickster-or-devil", "trickster, devil, demon, Satan (Protestant)"},
		{"wala", "wala", "wala", "CCONJ", "", "WHICH", 1, "or", "or, or else (between nouns or clauses)"},
		{"wali", "wâlï", "wâlï", "NOUN", "", "FAMILY", 2, "woman", "woman, wife"},
		{"wali", "wâlï", "wâlï", "NOUN", "", "WHERE", 2, "left-side", "left side"},
		{"wali", "wâlï", "wâlï", "NOUN", "Gender=Fem", "WHO", 1, "woman-female", "woman, female"},
		{"walikoli", "wâlï-kôlï", "wâlï-kɔ̂lï", "NOUN", "Gender=Fem", "WHO", 6, "lesbian", "[lit: woman|man]: 'butch' lesbian"},
		{"waliwali", "wâlï-wâlï", "wâlï-wâlï", "NOUN", "Gender=Fem", "WHO", 6, "lesbian", "[lit: woman|woman]: 'lipstick' lesbian"},
		{"wamabe", "wa-mä-bê", "wa-mä-bɛ̂", "NOUN", "", "GOD", 1, "believer", "believer, one who has faith"},
		{"wamandangokua", "wa-mändängö-kua", "wa-mändängɔ̈-kua", "NOUN", "", "WHO", 3, "apprentice", "[lit: one who|learning|work]: apprentice"},
		{"wamandangombeti", "wa-mändängö-mbëtï", "wa-mändängɔ̈-mbɛ̈tï", "NOUN", "", "WHO", 3, "pupil", "[lit: one who|learning|writing]: pupil"},
		{"wamandangoye", "wa-mändängö-yê", "wa-mändängɔ̈-yê", "NOUN", "", "WHO", 3, "student", "[lit: one who|learning|thing]: student"},
		{"wande", "wa-ndê", "wa-ndê", "NOUN", "", "WHO", 3, "foreigner", "[lit: one who|different]: foreigner"},
		{"wango", "wängö", "wängɔ̈", "VERB", "VerbForm=Vnoun", "INTERACT", 3, "advice", "advice, counsel"},
		{"wangobe", "wängö-bê", "wängɔ̈-bɛ̂", "VERB", "VerbForm=Vnoun", "INTERACT", 3, "regret", "regret"},
		{"wanzi", "wa-nzï", "wa-nzï", "NOUN", "", "ACT", 2, "thief", "[lit: one who|steal]: thief"},
		{"wapolisi", "wa-polîsi", "wa-polîsi", "NOUN", "", "WHO", 5, "policeman", "policeman"},
		{"wapulusu", "wa-pulûsu", "wa-pulûsu", "NOUN", "", "WHO", 5, "policeman", "policeman"},
		{"wara", "wara", "wara", "VERB", "Subcat=Tran", "STATE", 1, "obtain", "obtain, receive, find, earn, win, discover"},
		{"wasenda", "wa-sêndâ", "wa-sɛ̂ndâ", "NOUN", "", "WHO", 1, "scientist", "scientist"},
		{"wasungombeti", "wa-süngö-mbëtï", "wa-süngɔ̈-mbɛ̈tï", "VERB", "VerbForm=Vnoun", "INTERACT", 6, "writer", "[lit: one who|tracing|writing]: writer, author"},
		{"wataka", "wätäkä", "wätäkä", "NOUN", "", "INTERACT", 5, "untruth", "lie, untruth"},
		{"watokua", "wa-tokua", "wa-tokua", "NOUN", "", "ACT", 2, "ambassador-envoy", "[lit: one who|send|work]: ambassador, envoy"},
		{"waturu", "wa-turu", "wa-turu", "NOUN", "", "WHO", 4, "blacksmith", "[lit: one who|forge]: blacksmith"},
		{"wayanganzapa", "wa-yângâ-nzapä", "wa-yângâ-nzapä", "NOUN", "", "WHO", 3, "prophet", "[lit: one who|mouth|God]: prophet"},
		{"waziba", "wa-zibä", "wa-zibä", "NOUN", "", "SICK", 4, "blind-person", "blind person"},
		{"wen", "wên", "wên", "NOUN", "", "NATURE", 2, "iron-metal", "iron, metal"},
		{"were", "wêre", "wɛ̂rɛ", "VERB", "Subcat=Intr", "HOW", 2, "dry-out", "dry, dry out, dry up"},
		{"were", "wërë", "wërë", "NOUN", "", "HOW", 2, "game-or-sport", "[Sg: dry, because sports are played in the dry season]: game, sport"},
		{"werengo", "wërëngö", "wɛ̈rɛ̈ngɔ̈", "ADJ", "", "HOW", 2, "dry", "dry, dried, scrawny"},
		{"werengo", "wërëngö", "wɛ̈rɛ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "HOW", 2, "dryness", "dryness"},
		{"wo", "wo", "wɔ", "VERB", "Subcat=Intr", "ACT", 2, "breathe-or-rest", "breathe, rest"},
		{"wo", "wö", "wö", "VERB", "Subcat=Intr", "STATE", 2, "deflate-or-shrink", "diminish, go down, deflate, shrink"},
		{"wobe", "wo bê", "wɔ bɛ̂", "VERB", "Subcat=Tran", "ACT", 2, "calm-down", "calm down"},
		{"woga", "woga", "woga", "NOUN", "", "ANIM", 6, "antilope", " red-flanked duiker (small antilope)"},
		{"wogara", "wögarâ", "wögarâ", "NOUN", "Gender=Fem", "FAMILY", 4, "mother-in-law", "mother-in-law"},
		{"woko", "wôko", "wôko", "VERB", "Subcat=Intr", "STATE", 2, "become-weak-or-soft", "become weak, be fragile, soft or tender"},
		{"wokongo", "wököngö", "wököngɔ̈", "VERB", "VerbForm=Vnoun", "STATE", 2, "weakness-or-softness", "weakness, fragility, softness, tenderness"},
		{"womba", "wömba", "wömba", "NOUN", "", "FAMILY", 4, "paternal-relative", "woman's brother's child"},
		{"womba", "wömba", "wömba", "NOUN", "Gender=Fem", "FAMILY", 4, "paternal-aunt", "paternal aunt"},
		{"woro", "wörö", "wɔ̈rɔ̈", "NOUN", "", "CIVIL", 6, "poverty", "poverty"},
		{"wotere", "wo terê", "wɔ tɛrɛ̂", "VERB", "Subcat=Intr", "ACT", 2, "rest-relax", "rest, relax"},
		{"wotoro", "wôtoro", "wɔ̂tɔrɔ", "NOUN", "", "ANIM", 3, "bee", "bee"},
		{"woza", "woza", "woza", "VERB", "Subcat=Tran", "ACT", 4, "deceive", "erase footprints, fool, deceive, dupe"},
		{"wu", "wü", "wü", "VERB", "Subcat=Intr", "STATE", 2, "disperse", "diffuse, disperse, expand, be spread"},
		{"wuku", "wûku", "wûku", "VERB", "Subcat=Tran", "ACT", 3, "put-into", "put into (e.g. a sack)"},
		{"wungo", "wüngö", "wüngɔ̈", "NOUN", "", "NUM", 2, "numeral", "[lit: expanding]: numeral"},
		{"wungo", "wüngö", "wüngɔ̈", "VERB", "VerbForm=Vnoun", "STATE", 2, "diffusion", "diffusion, expansion, spread, fame, reputation"},
		{"wuruwuru", "wûrûwûrû", "wûrûwûrû", "NOUN", "", "STATE", 5, "tumult", "tumult, confusion, disorder"},
		{"wusuwusu", "wûsûwusu", "wûsûwusu", "ADV", "", "STATE", 5, "hairy", "hairy, hirsute, disheveled"},
		{"wuyawuya", "wûyâwuya", "wûyâwuya", "NOUN", "", "STATE", 5, "troubling", "troubling, sowing confusion"},
		{"ya", "ya", "ya", "NOUN", "Gender=Fem|Prefix=Yes", "WHO", 5, "Mrs.", "[followed by husband's name]: Mrs."},
		{"ya", "yâ", "yâ", "NOUN", "", "BODY", 1, "stomach", "stomach"},
		{"ya", "yâ", "yâ", "NOUN", "", "WHERE", 1, "inside", "inside, interior"},
		{"yaka", "yäkä", "yäkä", "NOUN", "", "WHERE", 1, "farm", "plantation, garden, farm, (cultivated) fields"},
		{"yakepaka", "yakepaka", "yakepaka", "NOUN", "Gender=Masc", "WHO", 6, "Mr.", "Mister"},
		{"yakere", "yakêrê", "yakɛ̂rɛ̂", "NUM", "NumType=Frac|Prefix=Yes", "NUM", 2, "milli-", "milli-"},
		{"yakerengonga", "yakêrê-ngonga", "yakɛ̂rɛ̂-ngonga", "NOUN", "", "WHEN", 2, "instant", "[lit: milli|hour]: instant, brief moment"},
		{"yanda", "yändä", "yändä", "NOUN", "", "STATE", 4, "sorcery", "enchantment, sorcery"},
		{"yanga", "yângâ", "yângâ", "NOUN", "", "BODY", 3, "mouth-or-language", "mouth, language"},
		{"yangada", "yângâ-da", "yângâ-da", "NOUN", "", "HOUSE", 3, "door", "[lit: mouth-house]: door"},
		{"yangba", "yangba", "yangba", "NOUN", "", "SICK", 6, "smallpox-or-chickenpox", "smallpox, chickenpox"},
		{"yango", "yangö", "yangö", "NOUN", "", "OBJ", 4, "hook", "hook, fish-hook"},
		{"yapakara", "yapakara", "yapakara", "NOUN", "Gender=Fem|Prefix=Yes", "WHO", 5, "Mrs.", "[followed by husband's name]: Mrs."},
		{"yapu", "yapü", "yapü", "ADJ", "", "STATE", 2, "light", "light (in weight)"},
		{"yapu", "yapü", "yapü", "NOUN", "", "STATE", 2, "lightness", "lightness (in weight)"},
		{"yapu", "yâpu", "yâpu", "VERB", "Subcat=Intr", "STATE", 2, "be-light-or-lighten", "be light (in weight), lighten (the weight)"},
		{"yaya", "yaya", "yaya", "NOUN", "Gender=Fem", "FAMILY", 5, "older-sibling", "[lit: wife|wife]: older sibling"},
		{"yayu", "yäyû", "yäyû", "NOUN", "", "GOD", 4, "heaven", "heaven"},
		{"yazo", "yazo", "yazo", "NOUN", "Gender=Fem", "WHO", 5, "Madam", "[lit: wife of|person]: Madam"},
		{"ye", "ye", "yɛ", "NOUN", "", "OBJ", 4, "fishing-net", "handheld fishing net"},
		{"ye", "yê", "yê", "NOUN", "", "FEEL", 1, "thing", "thing, object"},
		{"ye", "yê", "yê", "VERB", "", "FEEL", 1, "want", "want, wish, be about to"},
		{"yedaa", "yê-daä", "yê-daä", "VERB", "Subcat=Intr", "FEEL", 1, "be-willing", "be willing"},
		{"yeke", "yeke", "yɛkɛ", "ADV", "", "HOW", 1, "slow", "slow"},
		{"yeke", "yeke", "yɛkɛ", "VERB", "", "STATE", 1, "be", "be; [progressive modal verb] be"},
		{"yekena", "yeke na", "yɛkɛ na", "VERB", "Subcat=Tran", "STATE", 1, "have", "[lit: be|with]: have"},
		{"yeketi", "yeke tî", "yɛkɛ tî", "VERB", "Subcat=Tran", "STATE", 1, "belong-to", "[lit: be|of]: belong to"},
		{"yekeyeke", "yeke-yeke", "yɛkɛ-yɛkɛ", "ADV", "", "HOW", 1, "very-slow", "very slow"},
		{"yekpa", "yêkpâ", "yêkpâ", "NOUN", "", "NATURE", 4, "lightning", "sheet lightning"},
		{"yeme", "yême", "yême", "NOUN", "", "WHEN", 6, "June", "June"},
		{"yenga", "yenga", "yenga", "NOUN", "", "WHEN", 2, "week", "week, holiday, feast"},
		{"yengbi", "yêngbi", "yêngbi", "VERB", "Aspect=Imp|Subcat=Tran", "INTERACT", 2, "align", "harmonize, align"},
		{"yenge", "yênge", "yênge", "VERB", "Subcat=Tran", "ACT", 3, "shake-or-irrigate", "shake, stir, bump; water, moisten, irrigate"},
		{"yengengosese", "yëngëngö-sêse", "yëngëngɔ̈-sêse", "VERB", "VerbForm=Vnoun", "NATURE", 3, "earthquake", "earthquake"},
		{"yengere", "yengere", "yɛngɛrɛ", "VERB", "Subcat=Tran", "ACT", 3, "sift", "sift"},
		{"yengere", "yëngërë", "yɛ̈ngɛ̈rɛ̈", "NOUN", "", "ACT", 3, "sieve", "sieve"},
		{"yengodaa", "yëngö-daä", "yëngɔ̈-daä", "VERB", "VerbForm=Vnoun", "FEEL", 1, "willingness", "willingness, authorization, permission, permit, agreement"},
		{"yengondo", "yëngö-ndö", "yëngɔ̈-ndö", "NOUN", "", "GOD", 3, "God's-love-for-man", "[lit: above|want]: love of God for man"},
		{"yengondo", "yëngö-ndö", "yëngɔ̈-ndö", "VERB", "Subcat=Tran", "GOD", 3, "God's-love-for-man", "[lit: above|want]: love of God for man"},
		{"yere", "yërë", "yɛ̈rɛ̈", "NOUN", "", "CIVIL", 6, "poverty", "poverty"},
		{"yiko", "yîko", "yîkɔ", "VERB", "Subcat=Intr", "SENSE", 4, "disappear", "move away, go out of sight, disappear"},
		{"yingo", "yingö", "yingɔ̈", "NOUN", "", "GOD", 2, "spirit", "spirit"},
		{"yingo", "yingö", "yingɔ̈", "NOUN", "", "NATURE", 2, "shadow", "shadow"},
		{"yingogbia", "yingö-gbïä", "yingɔ̈-gbïä", "NOUN", "", "GOD", 2, "Holy-Spirit", "[lit: spirit|lord]: Holy Spirit"},
		{"yingova", "yingö-va", "yingɔ̈-va", "NOUN", "", "GOD", 2, "angel", "[lit: spirit|servant]: angel"},
		{"yingovuru", "yingö-vurü", "yingɔ̈-vurü", "NOUN", "", "GOD", 6, "Holy-Spirit", "[lit: spirit|white]: Holy Spirit"},
		{"yo", "yo", "yɔ", "VERB", "Subcat=Intr", "WHERE", 1, "be-far", "be long, far, tall"},
		{"yo", "yô", "yɔ̂", "VERB", "Subcat=Tran", "ACT", 2, "carry-off", "carry (in arms, on head), make off with"},
		{"yombo", "yömbö", "yɔ̈mbɔ̈", "NOUN", "", "OBJ", 4, "perfume", "perfume"},
		{"yongba", "yongba", "yongba", "NOUN", "", "BODY", 6, "stutter", "stutter, stammer"},
		{"yongoro", "yongôro", "yɔngɔ̂rɔ", "ADJ", "", "WHERE", 1, "far", "long, far, tall"},
		{"yoro", "yorö", "yɔrɔ̈", "NOUN", "", "SICK", 3, "medicine", "medicine, antidote, charm, talisman"},
		{"yoro", "yôro", "yɔ̂rɔ", "VERB", "Subcat=Tran", "ACT", 3, "fry", "fry, grill, roast, broil; introduce, insert, stuff"},
		{"yorongo", "yöröngö", "yɔ̈rɔ̈ngɔ̈", "ADJ", "", "ACT", 3, "roasted", "roasted"},
		{"yorongo", "yöröngö", "yɔ̈rɔ̈ngɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 3, "roasting", "roasting"},
		{"yu", "yü", "yü", "VERB", "Subcat=Tran", "ACT", 2, "wear", "wear (clothes, shoes)"},
		{"yuru", "yuru", "yuru", "VERB", "Subcat=Intr", "MOVE", 3, "flow", "flow, have diarrhea"},
		{"yuru", "yûru", "yûru", "VERB", "Subcat=Intr", "ACT", 3, "shove", "shove, repel; stretch out"},
		{"za", "zä", "zä", "VERB", "Subcat=Intr", "NATURE", 2, "be-lit-up", "be brilliantly lit up; be mean; be sharp; be sonorous"},
		{"za", "zä", "zä", "VERB", "Subcat=Tran", "NATURE", 2, "illuminate", "illuminate, light (lamp, fire)"},
		{"zabolo", "zâbolo", "zâbolo", "NOUN", "", "GOD", 4, "devil", "[Fr: diable]: devil, demon, Satan (Catholic)"},
		{"zakari", "zâkâri", "zâkâri", "VERB", "Subcat=Intr", "INTERACT", 4, "be-tangled", "be tangled"},
		{"zakarima", "zakarima", "zakarima", "VERB", "Subcat=Intr", "INTERACT", 6, "be-tangled", "be tangled"},
		{"zakazaka", "zâkâzaka", "zâkâzaka", "ADV", "", "SENSE", 5, "prickly", "prickly"},
		{"zalamaa", "zalamäa", "zalamäa", "NOUN", "", "COUNTRY", 1, "Germany", "Germany"},
		{"zambala", "zambala", "zambala", "NOUN", "", "ANIM", 5, "camel", "camel"},
		{"zango", "zängö", "zängɔ̈", "VERB", "VerbForm=Vnoun", "NATURE", 2, "brilliance", "brilliance, sharpness"},
		{"zaza", "zaza", "zaza", "NOUN", "", "OBJ", 5, "whip", "whip made of a bundle of twigs"},
		{"ze", "ze", "zɛ", "NOUN", "", "ANIM", 5, "panther", "panther, leopard"},
		{"zee", "zêe", "zêe", "VERB", "Subcat=Tran", "INTERACT", 3, "alert", "alert, warn"},
		{"zegbe", "zegbe", "zɛgbɛ", "NUM", "NumType=Frac|Prefix=Yes", "NUM", 2, "centi-", "centi-"},
		{"zegbengonga", "zegbe-ngonga", "zɛgbɛ-ngonga", "NOUN", "", "WHEN", 2, "second", "[lit: centi|hour]: second"},
		{"zege", "zegë", "zegë", "NOUN", "", "GAME", 4, "dice", "dice, divination with cards or dice"},
		{"zembe", "zembe", "zɛmbɛ", "NOUN", "", "OBJ", 2, "knife", "knife"},
		{"zeme", "zeme", "zɛmɛ", "NOUN", "", "OBJ", 6, "knife", "knife"},
		{"zen", "zën", "zën", "NOUN", "", "SICK", 2, "infirmity", "infirmity, disability"},
		{"zen", "zën", "zën", "VERB", "Subcat=Tran", "INTERACT", 2, "borrow-or-lend", "borrow; lend"},
		{"zengondo", "zëngö ndo", "zëngɔ̈ ndo", "VERB", "VerbForm=Vnoun", "INTERACT", 3, "alert", "alert, warning"},
		{"zi", "zî", "zî", "VERB", "Subcat=Tran", "ACT", 1, "open-release-remove", "liberate, release; open; remove, detach"},
		{"zi", "zï", "zï", "VERB", "Subcat=Tran", "ACT", 1, "dig-up", "dig, dig up"},
		{"zia", "zîâ", "zîâ", "VERB", "Subcat=Tran", "ACT", 1, "put-or-leave", "put, place, apply, add; imply; install, appoint; leave, abandon, cease; let, permit, authorize"},
		{"ziabe", "zîâ bê", "zîâ bɛ̂", "VERB", "Subcat=Intr", "FEEL", 1, "apply-oneself", "[lit: put|heart]: be faithful, commit (to), apply oneself"},
		{"ziabe", "zîâ-bê", "zîâ-bɛ̂", "NOUN", "", "FEEL", 1, "applying-oneself", "[lit: put|heart]: devotion, applying oneself"},
		{"ziakoli", "zîâ kôlï", "zîâ kɔ̂lï", "VERB", "Subcat=Intr", "ACT", 1, "divorce-one's-husband", "[lit: leave|man]: divorce one's husband"},
		{"zialegena", "zîâ lêgë na", "zîâ lêgë na", "VERB", "Subcat=Tran", "MOVE", 1, "authorize", "[lit: put|road|to]: authorize (someone) [+ti](do something)"},
		{"zianabe", "zîâ na bê", "zîâ na bɛ̂", "VERB", "Subcat=Tran", "FEEL", 1, "memorize", "[lit: put|in|heart]: memorize"},
		{"zianando", "zîâ na ndö", "zîâ na ndö", "VERB", "Subcat=Tran", "ACT", 1, "cover", "[lit: put|at|on top]: cover"},
		{"ziangianabeti", "zîâ ngîâ na bê tî", "zîâ ngîâ na bɛ̂ tî", "VERB", "Subcat=Tran", "FEEL", 1, "fill-with-joy", "[lit: put|joy|in|heart|of]: fill with joy"},
		{"ziango", "zïängö", "zïängɔ̈", "VERB", "VerbForm=Vnoun", "ACT", 1, "abandonment", "leaving, abandonment, cessation"},
		{"ziangu", "zîâ ngû", "zîâ ngû", "VERB", "Subcat=Intr", "ACT", 1, "irrigate", "[lit: apply|water]: water, irrigate"},
		{"ziawali", "zîâ wâlï", "zîâ wâlï", "VERB", "Subcat=Intr", "ACT", 1, "divorce-one's-wife", "[lit: leave|wife]: divorce one's wife"},
		{"ziayanga", "zîâ yângâ", "zîâ yângâ", "VERB", "Subcat=Tran", "FEEL", 1, "meddle", "[lit: put|mouth]: meddle, interfere"},
		{"ziba", "zibä", "zibä", "NOUN", "", "SICK", 4, "blindness", "blindness"},
		{"zibongo", "zî bongö", "zî bɔngɔ̈", "VERB", "Subcat=Intr", "ACT", 2, "undress", "[lit: remove|clothes]: undress"},
		{"zibongo", "zî-bongö", "zî-bɔngɔ̈", "VERB", "Subcat=Intr", "ACT", 6, "hornet", "[lit: remove|clothes]: yellowjacket, hornet"},
		{"zidoro", "zîdoro", "zîdɔrɔ", "NOUN", "", "FOOD", 4, "lemon", "lemon"},
		{"zigida", "zigidâ", "zigidâ", "NOUN", "", "OBJ", 6, "plastic-pearls", "plastic pearls"},
		{"zingo", "zîngo", "zîngɔ", "VERB", "Subcat=Intr", "ACT", 1, "get-up", "wake up, get up, rise up"},
		{"zingona", "zîngo na", "zîngɔ na", "VERB", "Subcat=Tran", "ACT", 1, "reprimand", "[lit: rise up|at]: reprimand"},
		{"zo", "zo", "zo", "NOUN", "", "WHO", 1, "person", "person, human being"},
		{"zo", "zö", "zɔ̈", "VERB", "Subcat=Tran", "ACT", 2, "burn", "burn, ignite, grill, roast"},
		{"zokuezo", "zo kûê zo", "zo kûɛ̂ zo", "INTERJ", "", "CIVIL", 1, "Everyone-is-someone", "[motto of CAR, cf. Ubuntu]: Everyone is someone"},
		{"zonga", "zonga", "zonga", "NOUN", "", "INTERACT", 2, "insult", "insult, offense, curse"},
		{"zonga", "zonga", "zonga", "VERB", "Aspect=Iter|Subcat=Tran", "INTERACT", 2, "insult", "insult, give offense, curse"},
		{"zongo", "zöngö", "zöngö", "NOUN", "", "NATURE", 6, "rainy-season", "rainy season"},
		{"zongo", "zöngö", "zɔ̈ngɔ̈", "ADJ", "", "ACT", 2, "grilled", "grilled"},
		{"zotinzi", "zo tî nzï", "zo tî nzï", "NOUN", "", "ACT", 2, "thief", "[lit: person|of|steal]: thief"},
		{"zowa", "zo wa", "zo wa", "NOUN", "", "WHICH", 1, "who", "[lit: person|which]: who"},
		{"zowa", "zö wâ", "zɔ̈ wâ", "VERB", "Subcat=Intr", "ACT", 2, "light-a-fire", "light a fire"},
		{"zua", "zûâ", "zûâ", "NOUN", "", "NATURE", 3, "island", "island"},
		{"zungo", "züngö", "züngɔ̈", "VERB", "VerbForm=Vnoun", "MOVE", 2, "descendant", "descendant"},
		{"zuu", "zûu", "zûu", "VERB", "Subcat=Intr", "MOVE", 2, "descend", "descend"},
		{"zuu", "zûu", "zûu", "VERB", "Subcat=Tran", "MOVE", 2, "lower-bow", "lower, bow"},
	}

	// Add rows derived from existing entries by the affixing of prefixes or suffixes.
	// These don't have to be completely productive and may include incorrect lexemes.
	// The goal is to increase recall, not precision.
	//
	// For each verb not ending in "nga" or "ngbi", add it.
	derivedRows := DictRows{}
	for _, row := range rows {
		if row.Lemma != "" && row.UDPos == "VERB" &&
			!strings.ContainsAny(row.Lemma, "- ") && !strings.Contains(row.UDFeature, "VerbForm=") {
			if !strings.HasSuffix(row.Lemma, "nga") && !strings.Contains(row.UDFeature, "Aspect=Iter") {
				r := row
				r.Toneless += "nga"
				r.Heightless += "nga"
				r.Lemma += "nga"
				if r.UDFeature != "" {
					r.UDFeature += "|"
				}
				r.UDFeature += "Aspect=Iter"
				r.EnglishTranslation += "-repeatedly"
				if strings.HasPrefix(row.Category, "ALT") {
					r.EnglishDefinition += "nga"
				} else {
					r.EnglishDefinition += "-repeatedly"
				}
				derivedRows = append(derivedRows, r)
			}
			if !strings.HasSuffix(row.Lemma, "ngbi") && !strings.Contains(row.UDFeature, "Aspect=Imp") {
				r := row
				r.Toneless += "ngbi"
				r.Heightless += "ngbi"
				r.Lemma += "ngbi"
				if r.UDFeature != "" {
					r.UDFeature += "|"
				}
				r.UDFeature += "Aspect=Imp"
				r.EnglishTranslation += "-together"
				if strings.HasPrefix(row.Category, "ALT") {
					r.EnglishDefinition += "ngbi"
				} else {
					r.EnglishDefinition += "-together"
				}
				derivedRows = append(derivedRows, r)
			}
		}
	}
	rows = append(rows, derivedRows...)

	// For each verb not ending in "ngɔ̈, add it and change all other vowels to mid pitch.
	derivedRows = DictRows{}
	for _, row := range rows {
		if row.Lemma != "" && row.UDPos == "VERB" && !strings.ContainsAny(row.Lemma, "- ") &&
			!strings.Contains(row.UDFeature, "VerbForm=") && !strings.HasSuffix(row.Lemma, "ngɔ̈") {
			r := row
			r.Toneless += "ngo"
			r.Heightless += "ngö"
			r.Lemma = norm.NFD.String(r.Lemma)
			r.Lemma = strings.ReplaceAll(r.Lemma, "a", "a\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "e", "e\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "ɛ", "ɛ\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "ə", "ə\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "i", "i\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "o", "o\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "ɔ", "ɔ\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "ø", "ø\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "u", "u\u0308")
			r.Lemma = strings.ReplaceAll(r.Lemma, "\u0302", "")
			r.Lemma = strings.ReplaceAll(r.Lemma, "\u0308\u0308", "\u0308")
			r.Lemma = norm.NFC.String(r.Lemma)
			r.Lemma += "ngɔ̈"
			if r.UDFeature != "" {
				r.UDFeature += "|"
			}
			// Gerunds can function as both NOUN and VERB.
			// Leave it as a VERB but note that it is a noun in the feature list.
			r.UDFeature += "VerbForm=Vnoun"
			r.EnglishTranslation += "ing"
			if strings.HasPrefix(row.Category, "ALT") {
				r.EnglishDefinition = norm.NFD.String(r.EnglishDefinition)
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "a", "a\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "e", "e\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "ɛ", "ɛ\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "ə", "ə\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "i", "i\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "o", "o\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "ɔ", "ɔ\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "ø", "ø\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "u", "u\u0308")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "\u0302", "")
				r.EnglishDefinition = strings.ReplaceAll(r.EnglishDefinition, "\u0308\u0308", "\u0308")
				r.EnglishDefinition = norm.NFC.String(r.EnglishDefinition)
				r.EnglishDefinition += "ngɔ̈"
			} else {
				r.EnglishDefinition += "ing"
			}
			derivedRows = append(derivedRows, r)
		}
	}

	// For each verb, prefix a "wa-".
	derivedRows = DictRows{}
	for _, row := range rows {
		if row.Lemma != "" && row.UDPos == "VERB" {
			if !strings.ContainsAny(row.Lemma, "- ") && !strings.HasPrefix(row.Lemma, "wa") {
				r := row
				r.Toneless = "wa" + r.Toneless
				r.Heightless = "wa" + r.Heightless
				r.Lemma = "wa-" + r.Lemma
				r.UDPos = "NOUN"
				r.Category = "WHO"
				r.EnglishTranslation = "one-who-" + r.EnglishTranslation
				if strings.HasPrefix(row.Category, "ALT") {
					r.EnglishDefinition = "wa" + r.EnglishDefinition
				} else {
					r.EnglishDefinition = "one-who-" + r.EnglishDefinition
				}
				derivedRows = append(derivedRows, r)
			}
		}
	}
	rows = append(rows, derivedRows...)

	// Pluralize each noun. Finitize each verb.
	derivedRows = DictRows{}
	for _, row := range rows {
		if row.Lemma != "" && !strings.ContainsAny(row.Lemma, "- ") && !strings.Contains(row.UDFeature, "VerbForm=") {
			if (row.UDPos == "NOUN" || row.UDPos == "ADJ") && !strings.HasPrefix(row.Lemma, "â") &&
				!strings.Contains(row.UDFeature, "Number=") {
				r := row
				r.Toneless = "a" + r.Toneless
				r.Heightless = "â" + r.Heightless
				r.Lemma = "â" + r.Lemma
				if r.UDFeature != "" {
					r.UDFeature += "|"
				}
				r.UDFeature += "Number=Plur"
				r.EnglishTranslation += "s"
				if strings.HasPrefix(row.Category, "ALT") {
					r.EnglishDefinition = "â" + r.EnglishDefinition
				} else {
					r.EnglishDefinition += "s"
				}
				derivedRows = append(derivedRows, r)
			}
			if row.UDPos == "VERB" && !strings.HasPrefix(row.Lemma, "a") &&
				!strings.Contains(row.UDFeature, "Person=") {
				r := row
				r.Toneless = "a" + r.Toneless
				r.Heightless = "a" + r.Heightless
				r.Lemma = "a" + r.Lemma
				if r.UDFeature != "" {
					r.UDFeature += "|"
				}
				r.UDFeature += "Person=3|VerbForm=Fin"
				if hyphenIndex := strings.Index(r.EnglishTranslation, "-"); hyphenIndex >= 0 {
					r.EnglishTranslation = r.EnglishTranslation[:hyphenIndex] + "s" + r.EnglishTranslation[hyphenIndex:]
				} else {
					r.EnglishTranslation += "s"
				}
				if strings.HasPrefix(row.Category, "ALT") {
					r.EnglishDefinition = "a" + r.EnglishDefinition
				} else {
					if hyphenIndex := strings.Index(r.EnglishDefinition, "-"); hyphenIndex >= 0 {
						r.EnglishDefinition = r.EnglishDefinition[:hyphenIndex] + "s" + r.EnglishDefinition[hyphenIndex:]
					} else {
						r.EnglishDefinition += "s"
					}
				}
				derivedRows = append(derivedRows, r)
			}
		}
	}
	rows = append(rows, derivedRows...)
	derivedRows = DictRows{}
	if len(derivedRows) != 0 {
		panic("Derived lexicon entries were left orphaned")
	}

	// Unique-sort the UDFeature set.
	for k, _ := range rows {
		features := strings.Split(rows[k].UDFeature, "|")
		slices.Sort(features)
		rows[k].UDFeature = strings.Join(slices.Compact(features), "|")
	}
	// Sort lexicon entries.
	rowLess := func(lhs, rhs DictRow) int {
		if c := strings.Compare(lhs.Toneless, rhs.Toneless); c != 0 {
			return c
		}
		if c := compareUnicode(lhs.Heightless, rhs.Heightless); c != 0 {
			return c
		}
		if c := compareUnicode(lhs.Lemma, rhs.Lemma); c != 0 {
			return c
		}
		if c := strings.Compare(lhs.UDPos, rhs.UDPos); c != 0 {
			return c
		}
		if c := strings.Compare(lhs.UDFeature, rhs.UDFeature); c != 0 {
			return c
		}
		if c := strings.Compare(lhs.Category, rhs.Category); c != 0 {
			return c
		}
		return rhs.Frequency - lhs.Frequency
	}
	rowEquiv := func(lhs, rhs DictRow) bool {
		return rowLess(lhs, rhs) == 0
	}
	// By sorting stably, static entries take precedence over derived ones.
	slices.SortStableFunc(rows, rowLess)
	rows = slices.CompactFunc(rows, rowEquiv)

	// Preallocate memory to save space.
	numRows := len(rows)
	var cols DictCols
	cols.Toneless = make([][]byte, numRows)
	cols.Heightless = make([][]byte, numRows)
	cols.LemmaRunes = make([][]rune, numRows)
	cols.LemmaUTF8 = make([][]byte, numRows)
	cols.UDPos = make([][]byte, numRows)
	cols.UDFeature = make([][]byte, numRows)
	cols.Category = make([][]byte, numRows)
	cols.Frequency = make([]int, numRows)
	cols.EnglishTranslation = make([][]byte, numRows)
	cols.EnglishDefinition = make([][]byte, numRows)

	// Count the total number of bytes and runes.
	numRunes := 0
	numBytes := 0
	for _, r := range rows {
		lenToneless := len(r.Toneless)
		lenHeightless := len(r.Heightless)
		lenLemma := utf8.RuneCountInString(r.Lemma)
		lenLemmaUTF8 := len(r.Lemma)
		lenUDPos := len(r.UDPos)
		lenUDFeature := len(r.UDFeature)
		lenCategory := len(r.Category)
		lenEnglishTranslation := len(r.EnglishTranslation)
		lenEnglishDefinition := len(r.EnglishDefinition)

		numBytes += lenToneless
		numBytes += lenHeightless
		numRunes += lenLemma
		numBytes += lenLemmaUTF8
		numBytes += lenUDPos
		numBytes += lenUDFeature
		numBytes += lenCategory
		numBytes += lenEnglishTranslation
		numBytes += lenEnglishDefinition
	}

	// Copy over the data into backing arrays for maximum locality.
	cols.Bytes = make([]byte, numBytes)
	cols.Runes = make([]rune, numRunes)
	endRune := 0
	endByte := 0
	for k, r := range rows {
		startByte := endByte
		endByte += copy(cols.Bytes[startByte:], []byte(r.Toneless))
		cols.Toneless[k] = cols.Bytes[startByte:endByte]

		startByte = endByte
		endByte += copy(cols.Bytes[startByte:], []byte(r.Heightless))
		cols.Heightless[k] = cols.Bytes[startByte:endByte]

		startRune := endRune
		endRune += copy(cols.Runes[startRune:], []rune(r.Lemma))
		cols.LemmaRunes[k] = cols.Runes[startRune:endRune]

		startByte = endByte
		endByte += copy(cols.Bytes[startByte:], []byte(r.Lemma))
		cols.LemmaUTF8[k] = cols.Bytes[startByte:endByte]

		startByte = endByte
		endByte += copy(cols.Bytes[startByte:], []byte(r.UDPos))
		cols.UDPos[k] = cols.Bytes[startByte:endByte]

		startByte = endByte
		endByte += copy(cols.Bytes[startByte:], []byte(r.UDFeature))
		cols.UDFeature[k] = cols.Bytes[startByte:endByte]

		startByte = endByte
		endByte += copy(cols.Bytes[startByte:], []byte(r.Category))
		cols.Category[k] = cols.Bytes[startByte:endByte]

		cols.Frequency[k] = r.Frequency

		startByte = endByte
		endByte += copy(cols.Bytes[startByte:], []byte(r.EnglishTranslation))
		cols.EnglishTranslation[k] = cols.Bytes[startByte:endByte]

		startByte = endByte
		endByte += copy(cols.Bytes[startByte:], []byte(r.EnglishDefinition))
		cols.EnglishDefinition[k] = cols.Bytes[startByte:endByte]
	}
	if endByte != numBytes {
		panic("Bad endByte")
	}
	if len(cols.Bytes) != numBytes {
		panic("Bad len(cols.Bytes)")
	}
	if cap(cols.Bytes) != numBytes {
		panic("Bad cap(cols.Bytes)")
	}
	if endRune != numRunes {
		panic("Bad endRune")
	}
	if len(cols.Runes) != numRunes {
		panic("Bad len(cols.Runes)")
	}
	if cap(cols.Runes) != numRunes {
		panic("Bad cap(cols.Runes)")
	}
	if len(cols.Frequency) != numRows {
		panic("Bad len(cols.Frequency)")
	}
	if cap(cols.Frequency) != numRows {
		panic("Bad cap(cols.Frequency)")
	}

	heightlessFromLemma := make(map[string]string)
	tonelessFromHeightless := make(map[string]string)
	for _, row := range rows {
		heightlessFromLemma[row.Lemma] = row.Heightless
		tonelessFromHeightless[row.Heightless] = row.Toneless
	}

	rowsMatchingLemma := make(DictRowsMap)
	rowsMatchingHeightless := make(DictRowsMap)
	rowsMatchingToneless := make(DictRowsMap)

	boundsMatchingLemma := make(map[string][2]int)
	boundsMatchingHeightless := make(map[string][2]int)
	boundsMatchingToneless := make(map[string][2]int)
	for k, row := range rows {
		newBounds := [2]int{k, k + 1}

		if bounds, found := boundsMatchingLemma[row.Lemma]; found {
			newBounds[0] = min(bounds[0], newBounds[0])
			newBounds[1] = max(bounds[1], newBounds[1])
		}
		boundsMatchingLemma[row.Lemma] = newBounds

		if bounds, found := boundsMatchingHeightless[row.Lemma]; found {
			newBounds[0] = min(bounds[0], newBounds[0])
			newBounds[1] = max(bounds[1], newBounds[1])
		}
		boundsMatchingHeightless[row.Lemma] = newBounds

		if bounds, found := boundsMatchingHeightless[row.Heightless]; found {
			newBounds[0] = min(bounds[0], newBounds[0])
			newBounds[1] = max(bounds[1], newBounds[1])
		}
		boundsMatchingHeightless[row.Heightless] = newBounds

		if bounds, found := boundsMatchingToneless[row.Lemma]; found {
			newBounds[0] = min(bounds[0], newBounds[0])
			newBounds[1] = max(bounds[1], newBounds[1])
		}
		boundsMatchingToneless[row.Lemma] = newBounds
		if bounds, found := boundsMatchingToneless[row.Heightless]; found {
			newBounds[0] = min(bounds[0], newBounds[0])
			newBounds[1] = max(bounds[1], newBounds[1])
		}
		boundsMatchingToneless[row.Heightless] = newBounds
		if bounds, found := boundsMatchingToneless[row.Toneless]; found {
			newBounds[0] = min(bounds[0], newBounds[0])
			newBounds[1] = max(bounds[1], newBounds[1])
		}
		boundsMatchingToneless[row.Toneless] = newBounds
	}

	for token, bound := range boundsMatchingToneless {
		rowsMatchingToneless[token] = rows[bound[0]:bound[1]]
	}
	for token, bound := range boundsMatchingHeightless {
		rowsMatchingHeightless[token] = rows[bound[0]:bound[1]]
	}
	for token, bound := range boundsMatchingLemma {
		rowsMatchingLemma[token] = rows[bound[0]:bound[1]]
	}

	return dictRowsAndCols{
		rows:                   rows,
		cols:                   cols,
		heightlessFromLemma:    heightlessFromLemma,
		tonelessFromHeightless: tonelessFromHeightless,
		rowsMatchingToneless:   rowsMatchingToneless,
		rowsMatchingHeightless: rowsMatchingHeightless,
		rowsMatchingLemma:      rowsMatchingLemma,
	}
}()

func compareUnicode(lhs, rhs string) int {
	return strings.Compare(toComparable(lhs), toComparable(rhs))
}

func toComparable(s string) string {
	s = strings.ReplaceAll(s, "A", "A01")
	s = strings.ReplaceAll(s, "Ə", "E01")
	s = strings.ReplaceAll(s, "Ɛ", "E11")
	s = strings.ReplaceAll(s, "E", "E21")
	s = strings.ReplaceAll(s, "I", "I01")
	s = strings.ReplaceAll(s, "Ø", "O01")
	s = strings.ReplaceAll(s, "Ɔ", "O11")
	s = strings.ReplaceAll(s, "O", "O21")
	s = strings.ReplaceAll(s, "U", "U01")
	s = strings.ReplaceAll(s, "a", "a01")
	s = strings.ReplaceAll(s, "ə", "e01")
	s = strings.ReplaceAll(s, "ɛ", "e11")
	s = strings.ReplaceAll(s, "e", "e21")
	s = strings.ReplaceAll(s, "i", "i01")
	s = strings.ReplaceAll(s, "ø", "o01")
	s = strings.ReplaceAll(s, "ɔ", "o11")
	s = strings.ReplaceAll(s, "o", "o21")
	s = strings.ReplaceAll(s, "u", "u01")
	s = strings.ReplaceAll(s, "1\u0323", "0")
	s = strings.ReplaceAll(s, "1\u0308", "2")
	s = strings.ReplaceAll(s, "1\u0302", "3")
	if strings.ContainsRune(s, 0x0302) || strings.ContainsRune(s, 0x0308) || strings.ContainsRune(s, 0x0323) {
		panic("Bad orphaned combining mark")
	}
	return s
}
