package util

import (
	"encoding/json"

	"github.com/flowershop/backend/internal/model"
)

// MarshalJSON 序列化任意值为 JSON 字符串。
func MarshalJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// UnmarshalStrings 反序列化 JSON 字符串为字符串切片。
func UnmarshalStrings(s string) []string {
	var out []string
	if s == "" {
		return out
	}
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return out
	}
	return out
}

// UnmarshalJSONTo 反序列化 JSON 字符串到目标对象。
func UnmarshalJSONTo(s string, v interface{}) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}

// UnmarshalLogisticsEvents 反序列化物流轨迹 JSON。
func UnmarshalLogisticsEvents(s string) []model.LogisticsEvent {
	var out []model.LogisticsEvent
	if s == "" {
		return out
	}
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return out
	}
	return out
}
