package node

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yunling101/cayo/pkg/cayoserver/controller"
	"github.com/yunling101/cayo/pkg/model/node"
	"github.com/yunling101/cayo/pkg/types"
)

// NodeController 控制器
type NodeController struct {
	controller.BaseController
	Bind struct {
		Node  node.EndPoint
		State struct {
			ID    int
			State bool
		}
	}
}

// OperatorList 运营商
func (w NodeController) OperatorList(c *gin.Context) {
	w.RenderSuccess(c, node.GetOperatorMap())
}

// List 列表
func (w NodeController) List(c *gin.Context) {
	query, err := node.EndPoint{}.QueryList(c)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, query)
}

// Probe 获取
func (w NodeController) Probe(c *gin.Context) {
	defaultMenu, selectedMenu, err := node.EndPoint{}.Probe()
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, map[string]interface{}{
		"default":  defaultMenu,
		"selected": selectedMenu,
	})
}

// Add 添加
func (w NodeController) Add(c *gin.Context) {
	if err := c.ShouldBindJSON(&w.Bind.Node); err != nil {
		w.RenderFail(c, "运营商不能为空!")
		return
	}

	if v := w.Params(c, []types.Param{
		{Data: w.Bind.Node.Province, Desc: "省份"},
		{Data: w.Bind.Node.City, Desc: "城市"},
	}); !v {
		return
	}
	err := node.EndPoint{}.Add(w.Bind.Node)
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// SelectedStatus 选中状态
func (w NodeController) SelectedStatus(c *gin.Context) {
	if v := w.BindJSON(c, &w.Bind.State); !v {
		return
	}
	err := node.EndPoint{}.UpdateOne(map[string]interface{}{"id": w.Bind.State.ID}, map[string]interface{}{
		"selected": w.Bind.State.State,
	})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}

// Delete 删除
func (w NodeController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		w.RenderFail(c, "ID错误!")
		return
	}
	err = node.EndPoint{}.Delete(map[string]interface{}{"id": id})
	if v := w.RenderIfFail(c, err); !v {
		return
	}
	w.RenderSuccess(c, nil)
}
