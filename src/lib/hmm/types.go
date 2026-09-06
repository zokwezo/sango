package hmm

type (
	Syllable int8
	Tag      int8
)

const NumSyllables int = 13

func (syllable Syllable) String() string {
	return [NumSyllables]string{
		"UNKNOWN",
		"alice",
		"bob",
		"charlie",
		"cupertino",
		"jobs",
		"london",
		"paris",
		"saw",
		"steve",
		"to",
		"visited",
		"went",
	}[syllable]
}

var syllableFromNameMap_ = func() map[string]Syllable {
	m := make(map[string]Syllable, NumSyllables)
	for i := range NumSyllables {
		s := Syllable(i)
		m[s.String()] = s
	}
	return m
}()

func S(syllableName string) Syllable {
	return syllableFromNameMap_[syllableName]
}

// TAGS
const (
	UnknownPitch Tag = iota
	LowPitch
	MedPitch
	HighPitch
	NumTags
)

func (tag Tag) String() string {
	return [NumTags]string{
		"?",
		"_",
		":",
		"^",
	}[tag]
}

var tagFromNameMap_ = func() map[string]Tag {
	m := make(map[string]Tag, NumTags)
	for i := range NumTags {
		s := Tag(i)
		m[s.String()] = s
	}
	return m
}()

func T(tagName string) Tag {
	return tagFromNameMap_[tagName]
}

type HMM struct {
	Transition   [NumTags][NumTags]float64      // [prevTag][nextTag] = log P(nextTag  | prevTag)
	Emission     [NumTags][NumSyllables]float64 // [tag][syllable]    = log P(syllable | tag)
	StartTags    [NumTags]float64               // [tag]              = log P(tag)
	TagCounts    [NumTags]int                   // [tag] = # occurances of tag in training corpus
	NumSentences int                            // < 0 during Accumulate, > 0 after Generate
}

// TODO: Replace with Sango strings containing diacritics.
type SangoSyllable struct {
	Token     Syllable
	TokenName string
	TagIndex  Tag
	TagName   string
}

// Metrics holds the classification performance statistics
type Metrics struct {
	Precision float64 // fraction of predicted that are correct
	Recall    float64 // fraction of correct that are predicted
	F1Score   float64 // harmonic mean of Precision and Recall
	TP        int     // # true  positives
	FP        int     // # false positives
	FN        int     // # false negatives
}
