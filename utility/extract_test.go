package utility

import "testing"

func TestGetLanguage(t *testing.T) {
	headers := map[string]string{"accept-language": "en-US,en;q=0.5"}

	lang := GetLanguage(headers)
	if lang != "en" {
		t.Errorf("expected language en but got %s", lang)
	}

	headers = map[string]string{"accept-language": "en"}
	lang = GetLanguage(headers)
	if lang != "en" {
		t.Errorf("expected language en but got %s", lang)
	}

	headers = map[string]string{"accept-language": "en-US"}
	lang = GetLanguage(headers)
	if lang != "en" {
		t.Errorf("expected language en but got %s", lang)
	}

	headers = map[string]string{"accept-language": "bn"}
	lang = GetLanguage(headers)
	if lang != "bn" {
		t.Errorf("expected language en but got %s", lang)
	}
}
