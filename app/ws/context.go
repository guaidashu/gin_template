package ws

import (
	"encoding/json"
	"fmt"
	"gin_template/app/enum"
	"gin_template/app/libs"
	"gin_template/app/libs/serror"
	"sync"
	"time"
)

type Context struct {
	values   map[string]interface{}
	clientId string
	client   *Client
	event    enum.WsEventEnum
	lock     sync.RWMutex
	data     []byte // 本次请求所携带到方法里的数据
}

func NewContext() *Context {
	return &Context{
		values: make(map[string]interface{}),
	}
}

func (c *Context) SetClient(clientId string, client *Client) {
	c.clientId = clientId
	c.client = client
}

func (c *Context) Client() *Client {
	return c.client
}

func (c *Context) ClientId() string {
	return c.clientId
}

func (c *Context) SetEvent(event enum.WsEventEnum) {
	c.event = event
}

func (c *Context) Set(key string, value interface{}) {
	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
	}()

	c.values[key] = value
}

func (c *Context) Get(key string) (interface{}, error) {
	c.lock.RLock()
	defer func() {
		c.lock.RUnlock()
	}()

	if _, ok := c.values[key]; !ok {
		return nil, serror.ErrContextKeyNotExistError
	}

	return c.values[key], nil
}

func (c *Context) SetData(data []byte) {
	c.data = data
}

func (c *Context) OriginData() []byte {
	return c.data
}

func (c *Context) send(data interface{}) {
	if c.client == nil {
		libs.Logger.Error(fmt.Sprintf("客户端 %v 已断开", c.clientId))
		return
	}

	c.client.Send(c.event, data)
}

func (c *Context) Send(data interface{}) {
	c.send(libs.GetSuccessReply(data))
}

func (c *Context) Error(err error) {
	c.send(libs.GetErrorReply(err))
}

func (c *Context) BindJson(obj interface{}) error {
	if c.data == nil {
		return serror.ErrDataIsNil
	}

	return json.Unmarshal(c.data, obj)
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err returns a non-nil error value after Done is closed,
// successive calls to Err return the same error.
// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
func (c *Context) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return 0
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}

	return nil
}
