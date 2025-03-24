package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"slices"
)

type Item struct {
	Title string
	Body  string
}

type API int

var Db []Item

func (a *API) GetDB(_ string, reply *[]Item) error {
	*reply = Db
	return nil
}

func (a *API) GetByName(title string, reply *Item) error {
	for _, val := range Db {
		if val.Title == title {
			*reply = val
		}
	}

	return nil
}

func (a *API) AddItem(item Item, reply *Item) error {
	Db = append(Db, item)
	*reply = item
	return nil
}

func (a *API) EditItem(edit Item, reply *Item) error {
	for idx, val := range Db {
		if val.Title == edit.Title {
			Db[idx] = Item{edit.Title, edit.Body}
			*reply = Db[idx]
			break
		}
	}

	return nil
}

func (a *API) DeleteItem(item Item, reply *Item) error {
	for idx, val := range Db {
		if val.Title == item.Title && val.Body == item.Body {
			Db = slices.Delete(Db, idx, idx+1)
			*reply = item
			break
		}
	}
	return nil
}

func main() {
	var api = new(API)
	err := rpc.Register(api)

	if err != nil {
		log.Fatal("Error registering API", err)
	}

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatal("Listener error: ", err)
	}
	log.Printf("Serving RPC on port: %d", 4040)

	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}
