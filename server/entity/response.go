package entity

type PostResponse struct {
	IsValid bool `json:"is_valid"`
}

type GetResponse struct {
	Valid   int64   `json:"count_valid"`
	Invalid int64   `json:"count_invalid"`
	Ratio   float32 `json:"ratio"`
}
