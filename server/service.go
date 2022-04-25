package server

import (
	"fmt"
	"log"
	"reflect"
	"runtime"

	"github.com/muque7/my-rpc/pkg"
	"github.com/muque7/my-rpc/utils"
)

type service struct {
	name       string                 // server name
	refVal     reflect.Value          // server reflect value
	refType    reflect.Type           // server reflect type
	methodType map[string]*methodType // server method
}

func newService(server interface{}, serverName string) (*service, error) {
	ser := &service{
		refVal:  reflect.ValueOf(server),
		refType: reflect.TypeOf(server),
	}

	sName := reflect.Indirect(ser.refVal).Type().Name()
	if !utils.IsPublic(sName) {
		return nil, pkg.ErrNonPublic
	}

	ser.name = sName
	if serverName != "" {
		ser.name = serverName
	}

	methods, err := constructionMethods(ser.refType)
	if err != nil {
		return nil, err
	}
	ser.methodType = methods

	for _, v := range methods {
		log.Println("Registry Service: ", ser.name, "   method: ", v.method.Name)
	}

	return ser, nil
}

// call 方法调用
func (s *service) call(mType *methodType, inArgs []reflect.Value) (result []reflect.Value, err error) {
	// recover 捕获堆栈消息
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("[painc service internal error]: %v, method: %s, argv: %+v, stack: %s",
				r, mType.method.Name, inArgs, buf)
			log.Println(err)
		}
	}()

	fn := mType.method.Func
	return fn.Call(append([]reflect.Value{s.refVal}, inArgs...)), nil
}
