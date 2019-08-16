package models

import (
	"tokensky_bg_admin/enums"
)

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Code    enums.JsonResultCode `json:"code"`
	Msg     string               `json:"msg"`
	Content interface{}          `json:"content"`
}

type JsonResult2 struct {
	Status  enums.JsonResultCode `json:"status"`
	Message string               `json:"message"`
}

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
}
