// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcp_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestContent(t *testing.T) {
	tests := []struct {
		in   mcp.Content
		want string // json serialization
	}{
		{
			&mcp.TextContent{Text: "hello"},
			`{"type":"text","text":"hello"}`,
		},
		{
			&mcp.TextContent{
				Text:        "hello",
				Meta:        mcp.Meta{"key": "value"},
				Annotations: &mcp.Annotations{Priority: 1.0},
			},
			`{"type":"text","text":"hello","_meta":{"key":"value"},"annotations":{"priority":1}}`,
		},
		{
			&mcp.ImageContent{
				Data:     []byte("a1b2c3"),
				MIMEType: "image/png",
			},
			`{"type":"image","mimeType":"image/png","data":"YTFiMmMz"}`,
		},
		{
			&mcp.ImageContent{
				Data:        []byte("a1b2c3"),
				MIMEType:    "image/png",
				Meta:        mcp.Meta{"key": "value"},
				Annotations: &mcp.Annotations{Priority: 1.0},
			},
			`{"type":"image","mimeType":"image/png","data":"YTFiMmMz","_meta":{"key":"value"},"annotations":{"priority":1}}`,
		},
		{
			&mcp.AudioContent{
				Data:     []byte("a1b2c3"),
				MIMEType: "audio/wav",
			},
			`{"type":"audio","mimeType":"audio/wav","data":"YTFiMmMz"}`,
		},
		{
			&mcp.AudioContent{
				Data:        []byte("a1b2c3"),
				MIMEType:    "audio/wav",
				Meta:        mcp.Meta{"key": "value"},
				Annotations: &mcp.Annotations{Priority: 1.0},
			},
			`{"type":"audio","mimeType":"audio/wav","data":"YTFiMmMz","_meta":{"key":"value"},"annotations":{"priority":1}}`,
		},
		{
			&mcp.EmbeddedResource{
				Resource: &mcp.ResourceContents{URI: "file://foo", MIMEType: "text", Text: "abc"},
			},
			`{"type":"resource","resource":{"uri":"file://foo","mimeType":"text","text":"abc"}}`,
		},
		{
			&mcp.EmbeddedResource{
				Resource: &mcp.ResourceContents{URI: "file://foo", MIMEType: "image/png", Blob: []byte("a1b2c3")},
			},
			`{"type":"resource","resource":{"uri":"file://foo","mimeType":"image/png","blob":"YTFiMmMz"}}`,
		},
		{
			&mcp.EmbeddedResource{
				Resource:    &mcp.ResourceContents{URI: "file://foo", MIMEType: "text", Text: "abc"},
				Meta:        mcp.Meta{"key": "value"},
				Annotations: &mcp.Annotations{Priority: 1.0},
			},
			`{"type":"resource","resource":{"uri":"file://foo","mimeType":"text","text":"abc"},"_meta":{"key":"value"},"annotations":{"priority":1}}`,
		},
	}

	for _, test := range tests {
		got, err := json.Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(test.want, string(got)); diff != "" {
			t.Errorf("json.Marshal(%v) mismatch (-want +got):\n%s", test.in, diff)
		}
		result := fmt.Sprintf(`{"content":[%s]}`, string(got))
		log.Println(result)
		var out mcp.CallToolResult
		if err := json.Unmarshal([]byte(result), &out); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(test.in, out.Content[0]); diff != "" {
			t.Errorf("json.Unmarshal(%q) mismatch (-want +got):\n%s", string(got), diff)
		}
	}
}

func TestEmbeddedResource(t *testing.T) {
	for _, tt := range []struct {
		rc   *mcp.ResourceContents
		want string // marshaled JSON
	}{
		{
			&mcp.ResourceContents{URI: "u", Text: "t"},
			`{"uri":"u","text":"t"}`,
		},
		{
			&mcp.ResourceContents{URI: "u", MIMEType: "m", Text: "t", Meta: mcp.Meta{"key": "value"}},
			`{"uri":"u","mimeType":"m","text":"t","_meta":{"key":"value"}}`,
		},
		{
			&mcp.ResourceContents{URI: "u"},
			`{"uri":"u"}`,
		},
		{
			&mcp.ResourceContents{URI: "u", Blob: []byte{}},
			`{"uri":"u","blob":""}`,
		},
		{
			&mcp.ResourceContents{URI: "u", Blob: []byte{1}},
			`{"uri":"u","blob":"AQ=="}`,
		},
	} {
		data, err := json.Marshal(tt.rc)
		if err != nil {
			t.Fatal(err)
		}
		if got := string(data); got != tt.want {
			t.Errorf("%#v:\ngot  %s\nwant %s", tt.rc, got, tt.want)
		}
		urc := new(mcp.ResourceContents)
		if err := json.Unmarshal(data, urc); err != nil {
			t.Fatal(err)
		}
		// Since Blob is omitempty, the empty slice changes to nil.
		if diff := cmp.Diff(tt.rc, urc); diff != "" {
			t.Errorf("mismatch (-want, +got):\n%s", diff)
		}
	}
}
