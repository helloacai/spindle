package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog"

	"github.com/helloacai/spindle/pkg/aciregistry"
	"github.com/helloacai/spindle/pkg/log"
	. "github.com/helloacai/spindle/pkg/util"
)

var agentClient http.Client

type Response struct {
	Messages MessageSlice `json:"messages"`
}

func (r *Response) MarshalZerologObject(e *zerolog.Event) {
	e.Array("messages", r.Messages)
}

type Message struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}

func (m *Message) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("status", string(m.Status)).
		Str("message", m.Message)
}

type MessageSlice []Message

func (s MessageSlice) MarshalZerologArray(a *zerolog.Array) {
	for _, msg := range s {
		a.Object(&msg)
	}
}

type Status string

const (
	Status_Debug    Status = "debug"
	Status_Info     Status = "info"
	Status_Waiting  Status = "waiting"
	Status_Complete Status = "complete"
)

func replaceString(s, requestRef, threadHex string) string {
	s = strings.ReplaceAll(s, "$requestRef", requestRef)
	s = strings.ReplaceAll(s, "$threadUID", threadHex)
	return s
}

func Call(ctx context.Context, metadata *aciregistry.Metadata, requestRef string, threadUID []byte, isNew bool) (*Response, error) {
	agentURL, err := url.Parse(metadata.BaseURL)
	if err != nil {
		return nil, err
	}

	var route aciregistry.RequestRoute
	if isNew {
		route = metadata.PostRoute
	} else {
		route = metadata.PatchRoute
	}

	agentURL = agentURL.JoinPath(route.URI)
	if len(route.QueryParams) > 0 {
		query := url.Values{}
		for _, p := range route.QueryParams {
			query.Add(p.Name, replaceString(p.Value, requestRef, Hex(threadUID)))
		}
		agentURL.RawQuery = query.Encode()
	}

	body := map[string]string{}
	for _, p := range route.BodyParams {
		body[p.Name] = replaceString(p.Value, requestRef, Hex(threadUID))
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("url", agentURL.String()).Str("body", string(bodyBytes)).Msg("querying agent")
	req, err := http.NewRequest(route.Method, agentURL.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	r, err := agentClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var response Response
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	log.Debug().Object("response", &response).Msg("agent response")

	return &response, nil
}
