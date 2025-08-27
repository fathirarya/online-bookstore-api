package model

type WebResponse[T any] struct {
	Paging *PageMetadata `json:"paging,omitempty"`
	Data   T             `json:"data"`
	Errors string        `json:"errors,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}
