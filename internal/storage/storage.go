package storage

type Storage interface {
	Save(u string) (string, error)
	Get(id string) (string, error)
}
