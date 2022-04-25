package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"

	"github.com/muque7/my-rpc/codec"
	"github.com/muque7/my-rpc/transport"
)

// Server struct
type Server struct {
	addr       string
	serviceMap map[string]*service
}

// NewServer creates a new server
func NewServer(addr string) *Server {
	return &Server{
		addr:       addr,
		serviceMap: map[string]*service{},
	}
}

func (s *Server) Register(service interface{}, serviceName string) error {
	ser, err := newService(service, serviceName)
	if err != nil {
		return err
	}

	log.Println(ser.name)
	s.serviceMap[ser.name] = ser
	return nil
}

// Run server
func (s *Server) Run() {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("listen on %s err: %v\n", s.addr, err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}

		go func() {
			srvTransport := transport.NewTransport(conn)

			for {
				// read request from client
				req, err := srvTransport.Receive()
				if err != nil {
					if err != io.EOF {
						log.Printf("read err: %v\n", err)
					}
					return
				}
				// get method by name
				log.Print(s.serviceMap)
				service, ok := s.serviceMap[req.ServiceName]
				if !ok { // if method requested does not exist
					e := fmt.Sprintf("service %s does not exist", req.ServiceName)
					log.Println(e)
					if err = srvTransport.Send(codec.Data{ServiceName: req.ServiceName, MethodName: req.MethodName, Err: e}); err != nil {
						log.Printf("transport write err: %v\n", err)
					}
					continue
				}

				method, ok := service.methodType[req.MethodName]
				if !ok { // if method requested does not exist
					e := fmt.Sprintf("func %s does not exist", req.MethodName)
					log.Println(e)
					if err = srvTransport.Send(codec.Data{ServiceName: req.ServiceName, MethodName: req.MethodName, Err: e}); err != nil {
						log.Printf("transport write err: %v\n", err)
					}
					continue
				}
				log.Printf("func %s.%s is called\n", req.ServiceName, req.MethodName)
				// unpackage request arguments

				// myReq := utils.RefNew(method.RequestType)

				// slice2Struct(req.Args, myReq)

				inArgs := make([]reflect.Value, len(req.Args))
				log.Print(req.Args)
				// mType := method.method.Type
				for i := range req.Args {
					// log.Printf("args %s is %s", i, req.Args[i])
					// convert json number to specific type like int
					// if mType.In(i+1).Kind() != reflect.TypeOf(req.Args[i]).Kind() {
					// inArgs[i] = reflect.ValueOf(req.Args[i])
					// } else {
					inArgs[i] = reflect.ValueOf(req.Args[i])
					// }
				}
				log.Print(inArgs)
				// invoke requested method
				out, _ := service.call(method, inArgs)
				// package response arguments (except error)
				outArgs := make([]interface{}, len(out)-1)
				for i := 0; i < len(out)-1; i++ {
					outArgs[i] = out[i].Interface()
				}
				// package error argument
				var e string
				if _, ok := out[len(out)-1].Interface().(error); !ok {
					e = ""
				} else {
					e = out[len(out)-1].Interface().(error).Error()
				}
				// send response to client
				err = srvTransport.Send(codec.Data{ServiceName: req.ServiceName, MethodName: req.MethodName, Args: outArgs, Err: e})
				if err != nil {
					log.Printf("transport write err: %v\n", err)
				}
			}
		}()
	}
}

func slice2Struct(arr []interface{}, u interface{}) error {
	valueOf := reflect.ValueOf(u)
	if valueOf.Kind() != reflect.Ptr {
		return errors.New("must ptr")
	}
	valueOf = valueOf.Elem()
	if valueOf.Kind() != reflect.Struct {
		return errors.New("must struct")
	}
	for i := 0; i < valueOf.NumField(); i++ {
		if i >= len(arr) {
			break
		}
		val := arr[i]
		if val != nil && reflect.ValueOf(val).Kind() == valueOf.Field(i).Kind() {
			valueOf.Field(i).Set(reflect.ValueOf(val))
		}
	}
	return nil
}

func strcutToSlice(in interface{}) []interface{} {
	v := reflect.ValueOf(in)
	ss := make([]interface{}, v.NumField())
	for i := range ss {
		ss[i] = v.Field(i)
	}
	return ss
}
