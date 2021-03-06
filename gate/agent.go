package gate

import (
	"github.com/abel/leaf/network"
)

type Agent interface {
	WriteMsg(msg interface{})
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	Conn() network.Conn
}
