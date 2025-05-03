package kset

func Select[Key, Value any](selector func(Value) Key, values ...Value) []Key {
	keys := make([]Key, 0, len(values))

	for i := range values {
		keys = append(keys, selector(values[i]))
	}

	return keys
}
