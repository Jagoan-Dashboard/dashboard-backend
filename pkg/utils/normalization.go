package utils

import (
	"strings"
	"unicode"
)


func NormalizeString(s string) string {
	if s == "" {
		return s
	}
	
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Replace multiple spaces with single space first
	s = strings.Join(strings.Fields(s), " ")
	
	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")
	
	// Remove special characters except underscore and alphanumeric
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// NormalizeStringKeepSpaces converts string to lowercase but keeps spaces
// Example: "Jalan Nasional" -> "jalan nasional"
func NormalizeStringKeepSpaces(s string) string {
	if s == "" {
		return s
	}
	
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Replace multiple spaces with single space
	s = strings.Join(strings.Fields(s), " ")
	
	return strings.TrimSpace(s)
}

// NormalizeEnum normalizes enum values (uppercase with underscore)
// Example: "Jalan Nasional" -> "JALAN_NASIONAL"
func NormalizeEnum(s string) string {
	if s == "" {
		return s
	}
	
	// Convert to uppercase
	s = strings.ToUpper(s)
	
	// Replace multiple spaces with single space first
	s = strings.Join(strings.Fields(s), " ")
	
	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")
	
	// Remove special characters except underscore and alphanumeric
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// NormalizeCommodityType normalizes commodity type values
func NormalizeCommodityType(s string) string {
	s = NormalizeEnum(s)
	
	// Map common variations
	switch s {
	case "FOOD", "FOOD_CROPS", "PANGAN":
		return "PANGAN"
	case "HORTICULTURE", "HORTIKULTURA":
		return "HORTIKULTURA"
	case "PLANTATION", "PERKEBUNAN":
		return "PERKEBUNAN"
	default:
		return s
	}
}

// NormalizeLocation normalizes location names (village, district)
// Keeps original casing for proper names but normalizes for search
func NormalizeLocation(s string) string {
	if s == "" {
		return s
	}
	
	// For locations, we keep spaces but normalize casing
	// Title case for proper names
	words := strings.Fields(strings.ToLower(s))
	for i, word := range words {
		if len(word) > 0 {
			// Capitalize first letter of each word
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	
	return strings.Join(words, " ")
}

// NormalizeForSearch creates a search-friendly version of string
// Example: "Jalan Raya No. 123" -> "jalan_raya_no_123"
func NormalizeForSearch(s string) string {
	if s == "" {
		return s
	}
	
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Remove punctuation
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)
	
	// Replace multiple spaces with single space
	s = strings.Join(strings.Fields(s), " ")
	
	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")
	
	return s
}

// NormalizeArrayStrings normalizes all strings in an array
func NormalizeArrayStrings(arr []string) []string {
	result := make([]string, len(arr))
	for i, s := range arr {
		result[i] = NormalizeString(s)
	}
	return result
}

// NormalizeMapKeys normalizes all keys in a map
func NormalizeMapKeys(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[NormalizeString(k)] = v
	}
	return result
}

// CompareNormalized compares two strings after normalization
func CompareNormalized(s1, s2 string) bool {
	return NormalizeString(s1) == NormalizeString(s2)
}

// ContainsNormalized checks if normalized s1 contains normalized s2
func ContainsNormalized(s1, s2 string) bool {
	return strings.Contains(NormalizeString(s1), NormalizeString(s2))
}

// NormalizeBeforeSave prepares a string for database storage
// Uses lowercase with underscores for consistency
func NormalizeBeforeSave(s string) string {
	return NormalizeString(s)
}

// NormalizeForDisplay prepares a string for display to users
// Title case with spaces
func NormalizeForDisplay(s string) string {
	if s == "" {
		return s
	}
	
	// Replace underscores with spaces
	s = strings.ReplaceAll(s, "_", " ")
	
	// Title case
	words := strings.Fields(strings.ToLower(s))
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	
	return strings.Join(words, " ")
}