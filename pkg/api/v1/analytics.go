package v1

type AnalyticsResponse struct {
	Views int `json:"views"`
	OS map[string]int64 `json:"os"`
	Browser map[string]int64 `json:"browser"`
}
