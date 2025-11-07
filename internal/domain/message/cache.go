package message

//go:generate go tool moq -out zmock_cache.go . Cache
type Cache interface {
	Set(v string)
	Exists(v string) bool
}
