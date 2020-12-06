package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
	log "github.com/sirupsen/logrus"
)

const (
	defaultTimeout = 30 * time.Second
)

type Adapter interface {
	Write(context.Context, *prompb.WriteRequest) error
	Read(context.Context, *prompb.ReadRequest) (*prompb.ReadResponse, error)
}

type Server struct {
	adapter Adapter
	logger  log.FieldLogger
	timeout time.Duration
}

type Option func(s *Server) error

func WithTimeout(t time.Duration) Option {
	return func(s *Server) error {
		s.timeout = t
		return nil
	}
}

func New(adapter Adapter, opts ...Option) (*Server, error) {
	s := &Server{
		adapter: adapter,
		logger:  log.New().WithField("component", "server"),
		timeout: defaultTimeout,
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) Write(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithField("method", "remoteWrite")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.WithField("error", err).Errorf("error while reading body")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	data, err := snappy.Decode(nil, body)
	if err != nil {
		logger.WithField("error", err).Errorf("error while decoding snappy compressed body")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var req prompb.WriteRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		logger.WithField("error", err).Errorf("error while unmarshaling protobuf")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctxx, cancel := context.WithTimeout(r.Context(), s.timeout)
	defer cancel()
	if err := s.adapter.Write(ctxx, &req); err != nil {
		logger.WithField("error", err).Errorf("error while writing data")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write([]byte("ok")); err != nil {
		logger.WithField("error", err).Errorf("error while writing response")
		return
	}
}

func (s *Server) Read(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithField("method", "read")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.WithField("error", err).Errorf("error while reading body")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data, err := snappy.Decode(nil, body)
	if err != nil {
		logger.WithField("error", err).Errorf("error while decoding snappy data")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var req prompb.ReadRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		logger.WithField("error", err).Errorf("error while unmarshaling protobuf")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctxx, cancel := context.WithTimeout(r.Context(), s.timeout)
	defer cancel()
	resp, err := s.adapter.Read(ctxx, &req)
	if err != nil {
		logger.WithField("error", err).Errorf("error while reading data")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	data, err = proto.Marshal(resp)
	if err != nil {
		logger.WithField("error", err).Errorf("error while marshaling protobuf response")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Content-Encoding", "snappy")
	respBody := snappy.Encode(nil, data)
	if _, err := w.Write(respBody); err != nil {
		logger.WithField("error", err).Errorf("error while writing response body")
		return
	}
}
