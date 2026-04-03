package utils

import "strings"

func GetLocalizedValue(values map[string]string, lang string) string {
	if name, ok := values[lang]; ok && name != "" {
		return name
	}

	if defaultName, ok := values[""]; ok && defaultName != "" {
		return defaultName
	}

	if strings.HasPrefix(lang, "en") {
		for _, fallback := range []string{"en_US", "en_GB", "en", "C"} {
			if name, ok := values[fallback]; ok && name != "" {
				return name
			}
		}
	}

	for _, name := range values {
		if name != "" {
			return name
		}
	}

	return ""
}

func GetAllLocales(m map[string]string, prefix string) (map[string]string, bool) {
	result := make(map[string]string)
	hasDefault := false
	hasAnyTranslation := false

	for key, value := range m {
		if key == prefix {
			result[""] = value
			hasDefault = true
			hasAnyTranslation = true
			continue
		}

		if strings.HasPrefix(key, prefix+"[") && strings.HasSuffix(key, "]") {
			lang := key[len(prefix)+1 : len(key)-1]

			if lang == "" && hasDefault {
				continue
			}

			result[lang] = value
			hasAnyTranslation = true

			if lang == "" {
				hasDefault = true
			}
		}
	}

	if !hasAnyTranslation {
		return nil, false
	}

	if !hasDefault {
		fallbackOrder := []string{"en", "en_US", "en_GB", "C", ""}

		for _, lang := range fallbackOrder {
			if value, ok := result[lang]; ok && value != "" {
				result[""] = value
				hasDefault = true
				break
			}
		}

		if !hasDefault && len(result) > 0 {
			for _, value := range result {
				if value != "" {
					result[""] = value
					break
				}
			}
		}
	}

	return result, true
}
