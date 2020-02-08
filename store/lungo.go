package store

import "github.com/256dpi/lungo"

func NewStorageMemory(name string) (lungo.IDatabase, *lungo.Engine) {
	client, engine, err := lungo.Open(nil, lungo.Options{
		Store: lungo.NewMemoryStore(),
	})
	if err != nil {
		panic(err)
	}
	return client.Database(name), engine
}
