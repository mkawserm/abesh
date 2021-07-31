package utility

func GetValue(data map[string]string, key string, defaultValue string) string {
	value, ok := data[key]
	if ok {
		return value
	}

	return defaultValue
}

func Merge(parent map[string]string, current map[string]string) map[string]string {
	h := make(map[string]string)

	if parent != nil {
		for k, v := range parent {
			h[k] = v
		}
	}

	if current != nil {
		for k, v := range current {
			h[k] = v
		}
	}

	return h
}
