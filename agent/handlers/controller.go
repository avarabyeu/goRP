package handlers

import (
	"fmt"
	"net/http"

	"github.com/avarabyeu/goRP/agent/store"
	"github.com/avarabyeu/goRP/gorp"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//NewMux creates new Mux/Controller
func NewMux(basePath string, client *gorp.Client, kvStore store.KVStore) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Heartbeat("/health"))
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(middleware.StripSlashes)

	// creates launch
	mux.Post(fmt.Sprintf("%s/{project}/async/launch", basePath), startLaunchHandler(kvStore, client))

	// finishes launch
	mux.Put(fmt.Sprintf("%s/{project}/async/launch/{launchID}/finish", basePath), finishLaunchHandler(kvStore, client))

	// creates root test item
	mux.Post(fmt.Sprintf("%s/{project}/async/item", basePath), startRootItemHandler(client, kvStore))

	// creates child test item
	mux.Post(fmt.Sprintf("%s/{project}/async/item/{parentID}", basePath), startTestItemHandler(kvStore, client))

	// finishes test item
	mux.Put(fmt.Sprintf("%s/{project}/async/item/{itemID}", basePath), finishItemHandler(kvStore, client))

	// creates log
	mux.Post(fmt.Sprintf("%s/{project}/async/log", basePath), func(w http.ResponseWriter, rq *http.Request) {
	})

	return mux
}
