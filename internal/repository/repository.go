package repository

var KeyValueStore map[string]interface{}
var MetadataStore map[string]interface{}
var ExpireKey = "EXPIRE"

func Init() {
	KeyValueStore = make(map[string]interface{})
	MetadataStore = make(map[string]interface{})
}
