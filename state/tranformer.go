package state

type (

	// KeyTransformer is used before putState operation for convert key
	KeyTransformer func(Key) (Key, error)

	// StringTransformer is used before setEvent operation for convert name
	StringTransformer func(string) (string, error)
)

// KeyAsIs returns string parts of composite key
func KeyAsIs(key Key) (Key, error) {
	return key, nil
}

func NameAsIs(name string) (string, error) {
	return name, nil
}
