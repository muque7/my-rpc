package server

import (
	"log"
	"reflect"

	"github.com/muque7/my-rpc/pkg"
	"github.com/muque7/my-rpc/utils"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type methodType struct {
	method reflect.Method
	// RequestType  reflect.Type
	// ResponseType reflect.Type
}

// constructionMethods Get specific method
func constructionMethods(typ reflect.Type) (map[string]*methodType, error) {
	methods := make(map[string]*methodType)
	log.Println(methods)
	for idx := 0; idx < typ.NumMethod(); idx++ {
		method := typ.Method(idx)
		// mType := method.Type
		mName := method.Name

		log.Println(mName)

		if !utils.IsPublic(mName) {
			return nil, pkg.ErrNonPublic
		}

		// // 默认是3个参数 func(*server.Method, *server.MethodReq, *server.MethodResp) error
		// if mType.NumIn() != 3 {
		// 	continue
		// }

		// // request 参数检查
		// requestType := mType.In(1)
		// if requestType.Kind() != reflect.Ptr {
		// 	continue
		// }

		// if !utils.IsPublicOrBuiltinType(requestType) {
		// 	continue
		// }

		// // response 参数检查
		// responseType := mType.In(2)
		// if responseType.Kind() != reflect.Ptr {
		// 	continue
		// }

		// if !utils.IsPublicOrBuiltinType(responseType) {
		// 	continue
		// }

		// // 校验返回参数
		// if mType.NumOut() != 1 {
		// 	continue
		// }

		// returnType := mType.Out(0)
		// if returnType != typeOfError {
		// 	continue
		// }

		methods[mName] = &methodType{
			method: method,
			// RequestType:  requestType,
			// ResponseType: responseType,
		}
	}

	if len(methods) == 0 {
		return nil, pkg.ErrNoAvailable
	}

	return methods, nil
}
