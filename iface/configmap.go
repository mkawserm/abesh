package iface

import (
	"strconv"
	"time"
)

func (v ConfigMap) String(key string, defaultValue string) string {
	if o, ok := v[key]; ok {
		return o
	}
	return defaultValue
}

func (v ConfigMap) Bytes(key string, defaultValue []byte) []byte {
	if o, ok := v[key]; ok {
		return []byte(o)
	}
	return defaultValue
}

func (v ConfigMap) Float32(key string, defaultValue float32) float32 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseFloat(o, 32)
		if err != nil {
			return float32(i)
		}

		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Float64(key string, defaultValue float64) float64 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseFloat(o, 64)
		if err != nil {
			return i
		}
		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Int(key string, defaultValue int) int {
	if o, ok := v[key]; ok {
		i, err := strconv.Atoi(o)
		if err != nil {
			return i
		}

		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Uint(key string, defaultValue uint) uint {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 32)
		if err != nil {
			return uint(i)
		}

		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Int16(key string, defaultValue int16) int16 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseInt(o, 10, 16)
		if err != nil {
			return int16(i)
		}

		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Uint16(key string, defaultValue uint16) uint16 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 16)
		if err != nil {
			return uint16(i)
		}

		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Int32(key string, defaultValue int32) int32 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseInt(o, 10, 32)
		if err != nil {
			return int32(i)
		}

		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Uint32(key string, defaultValue uint32) uint32 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 32)
		if err != nil {
			return uint32(i)
		}

		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Int64(key string, defaultValue int64) int64 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseInt(o, 10, 64)
		if err != nil {
			return i
		}
		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Uint64(key string, defaultValue uint64) uint64 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 64)
		if err != nil {
			return i
		}
		return defaultValue
	}
	return defaultValue
}

func (v ConfigMap) Duration(key string, defaultValue time.Duration) time.Duration {
	if o, ok := v[key]; ok {
		val, err := time.ParseDuration(o)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}
