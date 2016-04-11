package chanrpc

import (
	"reflect"
)

type AgentCaller struct {
	rpcRouter map[reflect.Type]reflect.Value
}

func (self *AgentCaller) Regist(msg interface{}, fn interface{}) bool {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return false
	}
	if self.rpcRouter == nil {
		self.rpcRouter = make(map[reflect.Type]reflect.Value, 10)
	}
	self.rpcRouter[msgType] = reflect.ValueOf(fn)
	return true
}

func (self *AgentCaller) Call(args []interface{}) {
	msgType := reflect.TypeOf(args[0])
	fn, found := self.rpcRouter[msgType]
	if found {
		fn.Call([]reflect.Value{reflect.ValueOf(args[0]), reflect.ValueOf(args[1])})
	}
}
