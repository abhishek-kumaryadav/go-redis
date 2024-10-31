package repository

var KeyValueStore map[string]interface{}
var MetadataStore map[string]interface{}
var EXPIRE_KEY = "EXPIRE"

func InitRepositories() {
	KeyValueStore = make(map[string]interface{})
	MetadataStore = make(map[string]interface{})
}
