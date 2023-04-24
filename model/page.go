package model

type PageInput struct {
	PageNum  int `json:"page_num"`  //页数
	PageSize int `json:"page_size"` //每页条数
}
