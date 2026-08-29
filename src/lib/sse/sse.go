// Sango Syllabic Encoding (SSE)
//
// Encoding for up to 4 unicode runes or a phonemically-valid Sango syllable.
//
// This is used internally because UTF8 (with both precomposed and combining
// diacritics) and the lossy standard Sango orthography (no vowel height)
// make it too cumbersome to use these latter formats.

package sse

import (
	"strings"
)

type SSE uint64

type WriteUtf8Options struct {
	ForSpaceUse, ForHyphenUse                    string
	WithShift, WithHeight, WithNTilde, WithPitch bool
}

var (
	AsToneless   = WriteUtf8Options{ForSpaceUse: "", ForHyphenUse: "", WithShift: false, WithHeight: false, WithNTilde: false, WithPitch: false}
	AsHeightless = WriteUtf8Options{ForSpaceUse: " ", ForHyphenUse: "-", WithShift: true, WithHeight: false, WithNTilde: false, WithPitch: true}
	AsLemma      = WriteUtf8Options{ForSpaceUse: " ", ForHyphenUse: "-", WithShift: true, WithHeight: true, WithNTilde: false, WithPitch: true}
	AsUtf8       = WriteUtf8Options{ForSpaceUse: " ", ForHyphenUse: "-", WithShift: true, WithHeight: true, WithNTilde: true, WithPitch: true}
)

type FromUtf8Options struct {
	TreatClosedVowelAsUnknownHeight  bool
	TreatUnmarkedPitchAsUnknownPitch bool
}

var (
	FromToneless   = FromUtf8Options{TreatClosedVowelAsUnknownHeight: true, TreatUnmarkedPitchAsUnknownPitch: true}
	FromHeightless = FromUtf8Options{TreatClosedVowelAsUnknownHeight: true, TreatUnmarkedPitchAsUnknownPitch: false}
	FromLemma      = FromUtf8Options{TreatClosedVowelAsUnknownHeight: false, TreatUnmarkedPitchAsUnknownPitch: false}
)

func (sse SSE) String() string             { return sse.toString() }
func (sse SSE) Less(rhs SSE) bool          { return sse.less(rhs) }
func CanonicalKey(canonical string) string { return canonicalKey(canonical) }
func CanonicalCompare(lhs, rhs string) int { return canonicalCompare(lhs, rhs) }

type SSEs []SSE

func (sses SSEs) Len() int           { return len(sses) }
func (sses SSEs) Swap(i, j int)      { sses[i], sses[j] = sses[j], sses[i] }
func (sses SSEs) Less(i, j int) bool { return sses[i].Less(sses[j]) }

func FromShortCode(shortCode uint64) SSE { return SSE(padRight(shortCode)) }
func (sse SSE) GetShortCode() uint64     { return unpadRight(uint64(sse)) }

func (sse SSE) WriteUtf8To(s *strings.Builder, options WriteUtf8Options) {
	writeUtf8To(s, uint64(sse), options)
}

func (sse SSE) WriteAsUtf8To(s *strings.Builder) {
	writeUtf8To(s, uint64(sse), AsUtf8)
}

func (sse SSE) WriteAsTonelessTo(s *strings.Builder) {
	writeUtf8To(s, uint64(sse), AsToneless)
}

func (sse SSE) WriteAsHeightlessTo(s *strings.Builder) {
	writeUtf8To(s, uint64(sse), AsHeightless)
}

func (sse SSE) WriteAsLemmaTo(s *strings.Builder) {
	writeUtf8To(s, uint64(sse), AsLemma)
}

func (sse SSE) WriteAsCanonicalTo(s *strings.Builder) {
	writeAsCanonicalTo(s, uint64(sse))
}

func Utf8ToSSEs(s string, options FromUtf8Options) (SSEs, error) {
	utf8ToCodes := func(u string) ([]sseCode, int) {
		return utf8ToCodes(u, options)
	}
	return toSSEs(s, utf8ToCodes)
}

func CanonicalToSSEs(s string) (SSEs, error) {
	return toSSEs(s, canonicalToCodes)
}
