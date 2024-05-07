package storages

type MemStorage struct {
	*AbstractStorage
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		AbstractStorage: NewAbstractStorage(),
	}
}

func (ms *MemStorage) WriteMetrics() error {
	return nil
}

func (ms *MemStorage) ReadMetrics() error {
	return nil
}
