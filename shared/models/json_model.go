package models

type Meta struct {
	Status    string `json:"status"`
	Code      int    `json:"code"`
	Message   any    `json:"message"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type SingleResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

type PagedResponse struct {
	Meta   Meta          `json:"meta"`
	Data   []interface{} `json:"data",omitempty`
	Paging Paging        `json:"paging,omitempty"`
}
