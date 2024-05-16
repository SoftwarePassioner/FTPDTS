// Copyright 2021 The Starship Troopers Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//webserver with rest api endpoints
package webserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type DataGetResponse struct {
	Response
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"createdAt"`
	TTL       uint        `json:"ttl"`
}

type DataPostResponse struct {
	Response
	UID string `json:"uid"`
}

var errNFound = Response{10, "Not found"}

//webserver options
type Opts struct {
	Port           uint
	Host           string
	MaxRequestBody int64
	DataStorage    DataStorage //data storage
	UIDGenerator   UID
	Logger         *log.Logger //Where log will be written to (default to stdout)
}

type DataStorage interface {
	Get(uid string) (payload interface{}, createdAt time.Time, ttl time.Duration, err error)
	Put(uid string, payload interface{}, ttl *time.Duration) error
}

//UID validator
type UID interface {
	//searching the UID in the string
	Validate(string) (string, error)
	New() string
}

type WebServer struct {
	logger         *log.Logger
	ds             DataStorage
	port           uint
	maxRequestBody int64
	uidGenerator   UID
	server         *http.Server
}

func New(o Opts) *WebServer {
	var mux http.ServeMux

	s := &WebServer{
		o.Logger,
		o.DataStorage,
		o.Port,
		o.MaxRequestBody,
		o.UIDGenerator,
		&http.Server{
			Addr:    fmt.Sprintf("%s:%d", o.Host, o.Port),
			Handler: &mux,
		},
	}
	mux.HandleFunc("/data", s.dataRequest)

	return s
}

func (s *WebServer) Run() error {
	s.logger.Printf("WEB server has been started at %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *WebServer) dataRequest(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	if req.Method == http.MethodPost {
		if res.Header().Get("Content-type") != "application/json" {
			http.Error(res, "wrong content-type (not a json)", http.StatusBadRequest)
			return
		}
		var d interface{}
		err := s.readBodyAsJSON(req, &d)
		if err != nil {
			http.Error(res, "wrong request data", http.StatusBadRequest)
			return
		}

		var ttl *time.Duration
		v, err := strconv.Atoi(req.FormValue("ttl"))
		if err == nil {
			d := time.Second * time.Duration(v)
			ttl = &d
		}

		//store data
		uid := s.uidGenerator.New()
		err = s.ds.Put(uid, d, ttl)
		if err != nil {
			s.logger.Printf("Can't store data into the datastorage: %v", err)
			http.Error(res, "Internal error", http.StatusInternalServerError)
			return
		}
		_, _ = res.Write(s.jsonResponse(DataPostResponse{Response{0, "OK"}, uid}))
		s.logger.Printf("New data has been stored into the storage with uid %s", uid)
		return
	}

	if req.Method == http.MethodGet {
		uid := req.FormValue("uid")
		if uid == "" {
			_, _ = res.Write(s.jsonResponse(errNFound))
			return
		}

		d, c, ttl, err := s.ds.Get(uid)
		if err != nil {
			_, _ = res.Write(s.jsonResponse(errNFound))
			return
		}

		_, _ = res.Write(s.jsonResponse(DataGetResponse{Response{0, "OK"}, d, c, uint(ttl / time.Second)}))
		s.logger.Printf("Data with uid %s has been presented", uid)
		return
	}

	http.Error(res, "Bad request", http.StatusBadRequest)
}

func (s *WebServer) Shutdown() {
	ctx := context.Background()
	s.logger.Printf("Shutting down the web server")
	_ = s.server.Shutdown(ctx)
}

func (s *WebServer) jsonResponse(d interface{}) []byte {
	b, err := json.Marshal(d)
	if err != nil {
		s.logger.Printf("ERROR Can't create a JSON response: %v", err)
		return nil
	}
	return b
}

func (s *WebServer) readBody(req *http.Request) ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0))
	if req.ContentLength > s.maxRequestBody {
		s.logger.Printf("Request body length greater then %d (maxRequestBody)\n", s.maxRequestBody)
		return nil, errors.New("Body is too large")
	}

	_, err := b.ReadFrom(io.LimitReader(req.Body, req.ContentLength))

	if err != nil {
		s.logger.Printf("Can't read request body %v\n", err)
		return nil, errors.New("Body read error")
	}
	_ = req.Body.Close()

	return b.Bytes(), nil
}

func (s *WebServer) readBodyAsJSON(req *http.Request, j *interface{}) error {
	b, err := s.readBody(req)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, j)
	if err != nil {
		s.logger.Printf("Can't parse json from request body %v\n", err)
		return errors.New("JSON error")
	}
	return nil
}
