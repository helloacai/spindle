package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
}

func (r *Response) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("status", string(r.Status)).
		Str("message", r.Message)
}

type ResponseStatus string

const (
	ResponseStatus_Waiting  = "waiting"
	ResponseStatus_Complete = "complete"
)

func replaceString(s, requestRef, threadHex string) string {
	s = strings.ReplaceAll(s, "$requestRef", requestRef)
	s = strings.ReplaceAll(s, "$threadUID", threadHex)
	return s
}

func Call(ctx context.Context, metadata *aciregistry.Metadata, requestRef string, threadUID []byte) (*Response, error) {
	agentURL, err := url.Parse(metadata.BaseURL)
	if err != nil {
		return nil, err
	}
	agentURL = agentURL.JoinPath(metadata.RequestRoute.URI)
	if len(metadata.RequestRoute.QueryParams) > 0 {
		query := url.Values{}
		for _, p := range metadata.RequestRoute.QueryParams {
			query.Add(p.Name, replaceString(p.Value, requestRef, Hex(threadUID)))
		}
		agentURL.RawQuery = query.Encode()
	}

	body := map[string]string{}
	for _, p := range metadata.RequestRoute.BodyParams {
		body[p.Name] = replaceString(p.Value, requestRef, Hex(threadUID))
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("url", agentURL.String()).Str("body", string(bodyBytes)).Msg("querying agent")
	req, err := http.NewRequest(metadata.RequestRoute.Method, agentURL.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	r, err := agentClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var response Response
	//if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
	//	return nil, err
	//}
	buf := new(strings.Builder)
	if _, err = io.Copy(buf, r.Body); err != nil {
		return nil, err
	}
	log.Debug().Str("response", buf.String()).Msg("agent response")
	response.Status = ResponseStatus_Complete
	response.Message = buf.String()

	return &response, nil
}
