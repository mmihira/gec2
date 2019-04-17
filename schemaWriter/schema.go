package schemaWriter

type Tag struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type NodeSchema struct {
	Name    string   `json:"name"`
	KeyName string   `json:"keyname"`
	Roles   []string `json:"roles"`
	Ip      string   `json:"ip"`
}

type Schema map[string]NodeSchema
