package pagination

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/yunling101/cayo/pkg/global"
)

// PageResponse 分页传入结构体
type PageResponse struct {
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Order  string     `json:"order"`
	Total  int        `json:"total"`
	Query  *mgo.Query `json:"query"`
}

// HTTPResponse 分页返回结构体
type HTTPResponse struct {
	Results interface{} `json:"results"`
	Total   int         `json:"total"`
}

// Pagination 分页结构体
type Pagination struct {
	Controller  *gin.Context
	Table       string                 `json:"table"`
	SearchField []string               `json:"search_field"`
	FilterField []string               `json:"filter_field"`
	NumberField []string               `json:"number_field"`
	NewQs       map[string]interface{} `json:"new"`
}

// Query 查询方法
func (p Pagination) Query() PageResponse {
	var (
		limit  int
		offset int
		order  string
	)

	requestLimit := p.Controller.Query("limit")
	requestOffset := p.Controller.Query("offset")
	requestAll := p.Controller.Query("all")
	if requestLimit == "" || requestOffset == "" {
		limit = 10
		offset = 0
	} else {
		limit, _ = strconv.Atoi(requestLimit)
		offset, _ = strconv.Atoi(requestOffset)
	}
	requestOrder := p.Controller.Query("order")
	pk := p.Controller.Query("pk")
	result := PageResponse{
		Limit:  limit,
		Offset: offset,
	}
	if requestOrder != "" && pk != "" {
		if requestOrder == "asc" {
			order = pk
		} else if requestOrder == "desc" {
			order = fmt.Sprintf("-%s", pk)
		} else {
			order = "-create_time"
		}
		result.Order = order
	} else {
		result.Order = "-create_time"
	}

	if len(p.NewQs) == 0 {
		p.NewQs = make(map[string]interface{})
	}
	for _, f := range p.FilterField {
		field := p.Controller.Query(f)
		if field != "" {
			if f == "id" || f == "status_code" || p.isContain(p.NumberField, f) {
				p.NewQs[f], _ = strconv.Atoi(field)
			} else if f == "task.id" {
				p.NewQs = bson.M{"task.id": bson.M{"$eq": field}}
			} else if f == "task.name" {
				p.NewQs = bson.M{"task.name": field}
			} else {
				switch field {
				case "true":
					p.NewQs[f] = true
				case "false":
					p.NewQs[f] = false
				default:
					p.NewQs[f] = field
				}
			}
		}
	}
	result.Query = global.DB.C(p.Table).Find(p.NewQs)
	if len(p.NewQs) == 0 {
		var cond []bson.M
		for _, s := range p.SearchField {
			search := p.Controller.Query("search")
			if search != "" {
				if s == "id" {
					ns, _ := strconv.Atoi(search)
					cond = append(cond, bson.M{s: bson.M{"$regex": ns}})
				} else {
					cond = append(cond, bson.M{s: bson.M{"$regex": search}})
				}
			}
		}
		if len(cond) != 0 {
			p.NewQs["$or"] = cond
			result.Query = global.DB.C(p.Table).Find(p.NewQs)
		}
	}
	result.Total, _ = result.Query.Count()

	// 条件查询所有不分页展示
	if requestAll == "true" {
		result.Limit = result.Total
	}
	return result
}

// isContain 是否存在列表
func (p Pagination) isContain(item []string, k string) bool {
	for _, v := range item {
		if v == k {
			return true
		}
	}
	return false
}
