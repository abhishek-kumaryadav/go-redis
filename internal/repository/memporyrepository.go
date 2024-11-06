package repository

var MemKeyValueStore map[string]interface{}
var MemMetadataStore map[string]interface{}

func InitMemoryRepository() {
	MemKeyValueStore = make(map[string]interface{})
	MemMetadataStore = make(map[string]interface{})
}
