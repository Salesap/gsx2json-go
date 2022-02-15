package cache

type Manager interface {
	Save(json []byte, key string) error

	Load(key string) ([]byte, error)

	List() []string

	Flush()
}
