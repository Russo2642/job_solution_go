package utils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	regNonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`)
	regMultipleDashes  = regexp.MustCompile(`-+`)
)

func GenerateSlug(name string) string {
	slug := strings.ToLower(name)

	slug = Transliterate(slug)

	slug = regNonAlphaNumeric.ReplaceAllString(slug, "-")

	slug = regMultipleDashes.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}

func GenerateUniqueSlug(name string, id int) string {
	baseSlug := GenerateSlug(name)
	return baseSlug + "-" + strconv.Itoa(id)
}

var transliterationMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh",
	'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
	'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
	'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu",
	'я': "ya",
	'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo", 'Ж': "Zh",
	'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M", 'Н': "N", 'О': "O",
	'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U", 'Ф': "F", 'Х': "Kh", 'Ц': "Ts",
	'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch", 'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu",
	'Я': "Ya",
}

func Transliterate(text string) string {
	var result strings.Builder
	result.Grow(len(text))

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, _ := transform.String(t, text)

	for _, r := range normalized {
		if transliterated, ok := transliterationMap[r]; ok {
			result.WriteString(transliterated)
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '_' {
			result.WriteRune(r)
		} else {
			result.WriteRune('-')
		}
	}

	return result.String()
}
