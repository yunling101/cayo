package client

import (
	"context"
	"fmt"

	"github.com/yunling101/cayo/pkg/global"
	"github.com/yunling101/cayo/pkg/model/client"
	"github.com/yunling101/cayo/pkg/propb"
)

// handlerFunc
type handlerFunc func(ctx context.Context, r *propb.Request) (*propb.Reply, error)

// RpcController 控制器
type RpcController struct{}

// verify
func (c *RpcController) verify(r *propb.Request) error {
	if r.Key == global.VerifyKey {
		return nil
	}
	return fmt.Errorf("%s", "Key verification failed")
}

// RenderSuccess 成功返回
func (c *RpcController) RenderSuccess(result map[string]string) (*propb.Reply, error) {
	return &propb.Reply{Code: 1, Result: result}, nil
}

// RenderFail 失败返回
func (c *RpcController) RenderFail(err error) (*propb.Reply, error) {
	return &propb.Reply{Code: 0, Msg: err.Error()}, nil
}

// decorator
func (c *RpcController) decorator(h handlerFunc) handlerFunc {
	return func(ctx context.Context, r *propb.Request) (*propb.Reply, error) {
		if err := c.verify(r); err != nil {
			return c.RenderFail(err)
		}
		return h(ctx, r)
	}
}

// Heartbeat 心跳监测
func (c *RpcController) Heartbeat(ctx context.Context, r *propb.Request) (*propb.Reply, error) {
	handle := func(ctx context.Context, r *propb.Request) (*propb.Reply, error) {
		h := client.New(r)
		if err := h.Heartbeat(); err != nil {
			return c.RenderFail(err)
		}
		return c.RenderSuccess(nil)
	}
	return c.decorator(handle)(ctx, r)
}

// ObtainTask 任务获取
func (c *RpcController) ObtainTask(ctx context.Context, r *propb.Request) (*propb.Reply, error) {
	handle := func(ctx context.Context, r *propb.Request) (*propb.Reply, error) {
		h := client.New(r)
		if b, err := h.ObtainTask(); err != nil {
			return c.RenderFail(err)
		} else {
			return c.RenderSuccess(b)
		}
	}
	return c.decorator(handle)(ctx, r)
}

// ReceiptTask 任务回执
func (c *RpcController) ReceiptTask(ctx context.Context, r *propb.Request) (*propb.Reply, error) {
	handle := func(ctx context.Context, r *propb.Request) (*propb.Reply, error) {
		go func() {
			h := client.New(r)
			h.ReceiptMetric()
		}()
		return c.RenderSuccess(nil)
	}
	return c.decorator(handle)(ctx, r)
}
