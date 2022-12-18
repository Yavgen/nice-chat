package main

import (
	"chat/internal"
	"chat/internal/action"
	"chat/internal/auth"
	"chat/internal/client"
	"chat/internal/domain/store"
	"chat/internal/handler"
	"chat/internal/pipe"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	registeredClientsCh := client.NewRegisteredClientChannel()
	unregisteredClientsCh := client.NewUnregisteredClientsChannel()
	broadcastCh := client.NewBroadcastChannel()

	registeredUsersStore := store.NewRegisteredUsersStore()
	loginUsersStore := store.NewLoginUsersStore()
	clientsStore := client.NewClientsStore()
	roomsStore := store.NewRoomsStore()

	actionFactory := action.NewFactory(broadcastCh, loginUsersStore, registeredUsersStore, roomsStore, clientsStore)

	authorizer := auth.NewAuthorizer(registeredUsersStore, loginUsersStore, clientsStore, roomsStore)

	//TODO вынести в папку к client и в отдельный неймспейс
	readPipe := pipe.NewReadPipe(actionFactory, authorizer)
	writePipe := pipe.NewWritePipe(authorizer, loginUsersStore)

	indexHandler := handler.NewIndexHandler()
	loginHandler := handler.NewLoginHandler(authorizer)
	wsHandler := handler.NewWSHandler(
		loginUsersStore,
		registeredClientsCh,
		authorizer,
		readPipe,
		writePipe,
		broadcastCh,
	)

	chatKernel := internal.NewKernel(registeredClientsCh, unregisteredClientsCh, broadcastCh, clientsStore)

	go chatKernel.Run()

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler.Handle)
	router.HandleFunc("/login", loginHandler.Handle)
	router.HandleFunc("/ws", wsHandler.Handle)
	http.Handle("/", router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
