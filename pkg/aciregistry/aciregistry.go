package aciregistry

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/helloacai/spindle/pkg/log"
	. "github.com/helloacai/spindle/pkg/util"
)

var (
	subgraphURL string
)

func init() {
	//subgraphURL = fmt.Sprintf("https://gateway.thegraph.com/api/%s/subgraphs/id/qDXHMsvxwV5VTkYz14PYUVY96z5CTicrVEiEK6Gwger",
	//	os.Getenv("SUBGRAPH_API_KEY"))
	subgraphURL = "https://api.studio.thegraph.com/query/63407/aciregistry-polygon-amoy/version/latest"

	log.Info().Str("subgraph_url", subgraphURL).Msg("subgraph initialized")
}

func pinataURL(uri string) string {
	return "https://aquamarine-big-mink-835.mypinata.cloud/ipfs/" + uri
}

func Get(ctx context.Context, aciUID []byte) (*Metadata, error) {
	// fetch metadata URI from subgraph
	metadataURI, err := fetchMetadataURI(ctx, aciUID)
	if err != nil {
		return nil, err
	}

	// fetch metadata from Pinata
	return fetchMetadata(ctx, metadataURI)
}

type MetadataQuery struct {
	Query         string            `json:"query"`
	OperationName string            `json:"operationName"`
	Variables     map[string]string `json:"variables"`
}

func fetchMetadataURI(_ context.Context, uid []byte) (string, error) {
	q := MetadataQuery{
		Query:         `{ registereds(where: {uid: "` + Hex(uid) + `"}) { aci_metadataURI } }`,
		OperationName: "Subgraphs",
		Variables:     map[string]string{},
	}
	body, err := json.Marshal(&q)
	if err != nil {
		return "", err
	}
	log.Debug().Str("body", string(body)).Msg("posting to subgraph")
	r, err := http.Post(subgraphURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	var response SubgraphResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Data.Registereds) == 0 {
		return "", errors.Errorf("uid %s not found", Hex(uid))
	}

	return response.Data.Registereds[0].ACIMetadataURI, nil
}

type SubgraphResponse struct {
	Data SubgraphResponseData `json:"data"`
}

type SubgraphResponseData struct {
	Registereds []Registered `json:"registereds"`
}

type Registered struct {
	ACIMetadataURI string `json:"aci_metadataURI"`
}

type Metadata struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Tools       StringSlice  `json:"tools"`
	BaseURL     string       `json:"baseUrl"`
	PostRoute   RequestRoute `json:"requestRoute"`
	PatchRoute  RequestRoute `json:"requestRoute"`
}

var _ zerolog.LogObjectMarshaler = &Metadata{}

func (m *Metadata) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("name", m.Name).
		Str("description", m.Description).
		Array("tools", m.Tools).
		Str("baseURL", m.BaseURL).
		Object("postRoute", &m.PostRoute).
		Object("patchRoute", &m.PatchRoute)
}

type StringSlice []string

var _ zerolog.LogArrayMarshaler = StringSlice{}

func (s StringSlice) MarshalZerologArray(a *zerolog.Array) {
	for _, v := range s {
		a.Str(v)
	}
}

type RequestRoute struct {
	URI         string     `json:"uri"`
	Method      string     `json:"method"`
	QueryParams ParamSlice `json:"queryParams"`
	BodyParams  ParamSlice `json:"bodyParams"`
}

var _ zerolog.LogObjectMarshaler = &RequestRoute{}

func (r *RequestRoute) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("uri", r.URI).
		Str("method", r.Method).
		Array("queryParams", r.QueryParams).
		Array("bodyParams", r.BodyParams)
}

type Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var _ zerolog.LogObjectMarshaler = &Param{}

func (p *Param) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("name", p.Name).
		Str("value", p.Value)
}

type ParamSlice []Param

func (s ParamSlice) MarshalZerologArray(a *zerolog.Array) {
	for _, p := range s {
		a.Object(&p)
	}
}

func fetchMetadata(_ context.Context, metadataURI string) (*Metadata, error) {
	log.Debug().Str("pinata_url", pinataURL(metadataURI)).Msg("getting metadata from pinata")
	r, err := http.Get(pinataURL(metadataURI))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var response Metadata
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
