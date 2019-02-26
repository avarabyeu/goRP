package handlers

import (
	"bytes"
	"fmt"
	"github.com/avarabyeu/goRP/agent/store"
	"github.com/avarabyeu/goRP/gorp"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func NewMux(basePath string, client *gorp.Client, store store.KVStore) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.Heartbeat("/health"))
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(middleware.StripSlashes)

	//creates launch
	mux.Post(fmt.Sprintf("%s/{project}/launch", basePath), JSONHandler(func(rq *http.Request) (interface{}, error) {
		uid, err := uuid.NewV4()
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.Wrap(errors.New("ok"), "Unable to generate UUID"))
		}

		body, err := ioutil.ReadAll(rq.Body)
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.New("Cannot read request body"))
		}

		go func(b []byte) {
			rs, err := client.StartLaunchRaw(bytes.NewBuffer(b))
			if nil != err {
				if err := store.Store("pending-launch", uid.String(), b); err != nil {
					log.Error(err)
				}
			} else {
				if err := store.Store("uuids", uid.String(), rs.ID); err != nil {
					log.Error(err)
				}
			}

		}(body)

		return &gorp.EntryCreatedRS{ID: uid.String()}, nil
	}))

	//finishes launch
	mux.Put(fmt.Sprintf("%s/{project}/launch/{launchID}/finish", basePath), JSONHandler(func(rq *http.Request) (interface{}, error) {
		launchID := chi.URLParam(rq, "launchID")
		realLaunchID, err := store.FindString("uuids", launchID)
		if err != nil {
			return nil, NewStatusErr(http.StatusBadRequest, errors.New("Unable to find request UUID"))
		}

		body, err := ioutil.ReadAll(rq.Body)
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.New("Cannot read request body"))
		}

		go func() {
			if _, err := client.FinishLaunchRaw(realLaunchID, bytes.NewBuffer(body)); err != nil {
				log.Error(err)
			}
		}()

		return &gorp.MsgRS{Msg: "Launch has been finished"}, nil
	}))

	//creates root test item
	mux.Post(fmt.Sprintf("%s/{project}/item", basePath), func(w http.ResponseWriter, rq *http.Request) {
	})

	//creates child test item
	mux.Post(fmt.Sprintf("%s/{project}/item/{parentID}", basePath), func(w http.ResponseWriter, rq *http.Request) {
	})

	//finishes test item
	mux.Put(fmt.Sprintf("%s/{project}/item/{itemID}", basePath), func(w http.ResponseWriter, rq *http.Request) {
	})

	//creates log
	mux.Post(fmt.Sprintf("%s/{project}/log", basePath), func(w http.ResponseWriter, rq *http.Request) {
	})

	return mux
}
