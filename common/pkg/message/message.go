package message

type ResponseErrors map[string]interface{}

type ResponseData struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg"`
	Data     interface{} `json:"data"`
}


type Pagination struct {
	Page    uint64 `form:"page"`
	PerPage uint32 `form:"perPage"`
}
