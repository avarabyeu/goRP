package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/avarabyeu/goRP/agent/store"
	"github.com/avarabyeu/goRP/gorp"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Save Log Request part with json data
const LogRequestJsonPart = "json_request_part"

// Save Log Request binary part
const LogRequestBinaryPart = "binary_part"

func logHandler(kvStore store.KVStore, client *gorp.Client) http.HandlerFunc {
	return JSONHandler(func(rq *http.Request) (interface{}, error) {

		//
		//itemID := chi.URLParam(rq, "itemID")

		err := rq.ParseMultipartForm(5 * 1024 * 1024)
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.New("Cannot read request body"))
		}

		rqBody := rq.FormValue(LogRequestBinaryPart)
		var logRQ gorp.SaveLogRQ
		err = json.Unmarshal([]byte(rqBody), &logRQ)
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.New("Cannot read request body"))
		}

		realItemID, err := kvStore.FindString("uuids", logRQ.ItemID)
		if err != nil {
			return nil, NewStatusErr(http.StatusBadRequest, errors.New("Unable to find real request UUID"))
		}

		for fName, f := range rq.MultipartForm.File {
			fmt.Println(fName)
			fmt.Println(f)
		}

		go func() {
			logRQ.ItemID = realItemID
			if _, err := client.SaveLog(&logRQ); err != nil {
				log.Error(err)
			}
		}()

		return &gorp.MsgRS{Msg: "Test Item has been finished"}, nil
	})
}

func finishItemHandler(kvStore store.KVStore, client *gorp.Client) http.HandlerFunc {
	return JSONHandler(func(rq *http.Request) (interface{}, error) {
		itemID := chi.URLParam(rq, "itemID")
		realItemID, err := kvStore.FindString("uuids", itemID)
		if err != nil {
			return nil, NewStatusErr(http.StatusBadRequest, errors.New("Unable to find real request UUID"))
		}

		body, err := ioutil.ReadAll(rq.Body)
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.New("Cannot read request body"))
		}

		go func() {
			if _, err := client.FinishTestRaw(realItemID, bytes.NewBuffer(body)); err != nil {
				log.Error(err)
			}
		}()

		return &gorp.MsgRS{Msg: "Test Item has been finished"}, nil
	})
}

func startTestItemHandler(kvStore store.KVStore, client *gorp.Client) http.HandlerFunc {
	return JSONHandler(func(rq *http.Request) (interface{}, error) {
		uid, err := uuid.NewV4()
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.Wrap(errors.New("ok"), "Unable to generate UUID"))
		}

		parentID := chi.URLParam(rq, "parentID")
		realParentID, err := kvStore.FindString("uuids", parentID)
		if err != nil {
			return nil, NewStatusErr(http.StatusBadRequest, errors.New("Unable to find request UUID"))
		}

		body, err := ioutil.ReadAll(rq.Body)
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.New("Cannot read request body"))
		}

		go func() {
			rs, err := client.StartChildTestRaw(realParentID, bytes.NewBuffer(body))
			if err != nil {
				if err := kvStore.Store("pending-item", uid.String(), body); err != nil {
					log.Error(err)
				}
			} else {
				if err := kvStore.Store("uuids", uid.String(), rs.ID); err != nil {
					log.Error(err)
				}
			}
		}()

		return &gorp.EntryCreatedRS{ID: uid.String()}, nil
	})
}

func startRootItemHandler(client *gorp.Client, kvStore store.KVStore) http.HandlerFunc {
	return JSONHandler(func(rq *http.Request) (interface{}, error) {
		uid, err := uuid.NewV4()
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.Wrap(errors.New("ok"), "Unable to generate UUID"))
		}

		body, err := ioutil.ReadAll(rq.Body)
		if err != nil {
			return nil, NewStatusErr(http.StatusInternalServerError, errors.New("Cannot read request body"))
		}

		go func() {
			rs, err := client.StartTestRaw(bytes.NewBuffer(body))
			if err != nil {
				if err := kvStore.Store("pending-item", uid.String(), body); err != nil {
					log.Error(err)
				}
			} else {
				if err := kvStore.Store("uuids", uid.String(), rs.ID); err != nil {
					log.Error(err)
				}
			}
		}()

		return &gorp.EntryCreatedRS{ID: uid.String()}, nil
	})
}

func finishLaunchHandler(kvStore store.KVStore, client *gorp.Client) http.HandlerFunc {
	return JSONHandler(func(rq *http.Request) (interface{}, error) {
		launchID := chi.URLParam(rq, "launchID")
		realLaunchID, err := kvStore.FindString("uuids", launchID)
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
	})
}

func startLaunchHandler(kvStore store.KVStore, client *gorp.Client) http.HandlerFunc {
	return JSONHandler(func(rq *http.Request) (interface{}, error) {
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
			if err != nil {
				if err := kvStore.Store("pending-launch", uid.String(), b); err != nil {
					log.Error(err)
				}
			} else {
				if err := kvStore.Store("uuids", uid.String(), rs.ID); err != nil {
					log.Error(err)
				}
			}

		}(body)

		return &gorp.EntryCreatedRS{ID: uid.String()}, nil
	})
}
