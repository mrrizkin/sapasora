package hash

import nanoid "github.com/matoous/go-nanoid/v2"

func NanoID(l ...int) string {
	id, _ := nanoid.New(l...)
	return id
}

func GenerateNanoID(chars string, size int) string {
	id, _ := nanoid.Generate(chars, size)
	return id
}
