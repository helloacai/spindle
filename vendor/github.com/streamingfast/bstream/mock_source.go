// Copyright 2019 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bstream

import (
	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	"github.com/streamingfast/shutter"
	"go.uber.org/zap"
)

type MockSource struct {
	*shutter.Shutter

	handler Handler
	blocks  []*pbbstream.Block
}

func NewMockSource(blocks []*pbbstream.Block, handler Handler) *MockSource {
	return &MockSource{
		Shutter: shutter.New(),
		blocks:  blocks,
		handler: handler,
	}
}

func (s *MockSource) Run() {
	for _, block := range s.blocks {
		err := s.handler.ProcessBlock(block, nil)
		if err != nil {
			s.Shutdown(err)
			return
		}
	}
}

func (s *MockSource) SetLogger(logger *zap.Logger) {
}
