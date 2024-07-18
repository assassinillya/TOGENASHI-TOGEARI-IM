package ctype

import (
	"database/sql/driver"
	"encoding/json"
)

type SystemMsg struct {
	Type int8 `json:"type"` // 违规类型 1 涉黄 2 涉恐 3 涉政 4 不正当言论
}

// Scan 入库的数据
func (c *SystemMsg) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}

// Value 入库的数据
func (c *SystemMsg) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
