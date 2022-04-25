package main

import (
	"errors"
	"log"

	"github.com/muque7/my-rpc/example/helloworld/public"
	"github.com/muque7/my-rpc/server"
)

type User struct{}

func (u *User) QueryUser(uid float64) (public.ResponseQueryUser, error) {
	db := make(map[float64]public.User)
	db[0] = public.User{Name: "Jiahonzheng", Age: 70}
	db[1] = public.User{Name: "ChiuSinYing", Age: 75}
	if u, ok := db[uid]; ok {
		return public.ResponseQueryUser{User: u, Msg: "success"}, nil
	}
	return public.ResponseQueryUser{User: public.User{}, Msg: "fail"}, errors.New("uid is not in database")
}

func main() {
	// gob.Register(public.ResponseQueryUser{})

	addr := "0.0.0.0:2333"
	srv := server.NewServer(addr)
	e := srv.Register(&User{}, "myService")
	if e != nil {
		log.Println(e)
		return
	}
	log.Println("service is running")
	go srv.Run()

	for {
	}
}
