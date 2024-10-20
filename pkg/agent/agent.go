package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/helloacai/spindle/pkg/aciregistry"
	"github.com/helloacai/spindle/pkg/log"
	. "github.com/helloacai/spindle/pkg/util"
)

var agentClient http.Client

func init() {
	agentClient = http.Client{
		Timeout: 2 * time.Minute,
	}
}

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

func Call(ctx context.Context, metadata *aciregistry.Metadata, requestRef string, parentThreadUID, threadUID []byte, isNew bool) (*Response, error) {
	agentURL, err := url.Parse(metadata.BaseURL)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing url")
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
	body["parentThreadUID"] = Hex(parentThreadUID)
	for _, p := range route.BodyParams {
		body[p.Name] = replaceString(p.Value, requestRef, Hex(threadUID))
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling request body")
	}

	log.Debug().Str("url", agentURL.String()).Str("body", string(bodyBytes)).Str("method", route.Method).Msg("querying agent")
	req, err := http.NewRequest(route.Method, agentURL.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, errors.Wrap(err, "error crafting http request to agent")
	}

	r, err := agentClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error doing http request to agent")
	}
	defer r.Body.Close()
	var response Response
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "error decoding agent request body")
	}
	log.Debug().Object("response", &response).Msg("agent response")

	return &response, nil
}
