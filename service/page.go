package service

type Page struct {
	Page   uint64 `json:"page" form:"page"`
	Number uint64 `json:"number" form:"number"`
}

type PagingContent struct {
	Page  *Page       `json:"page"`
	Datas interface{} `json:"datas"`
}
