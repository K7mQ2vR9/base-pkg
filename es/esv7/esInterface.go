package esv7

type IndexTable interface {
	IndexName() string
}

// Schema 提供 Elasticsearch 7 文档 ID，FindOne 和 Search 成功后将原生 _id 写入目标或结果元素。
type Schema interface {
	GetId() string
	SetId(id string)
}
