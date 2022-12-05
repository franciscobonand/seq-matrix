package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/franciscobonand/seq-matrix/db"
	"github.com/franciscobonand/seq-matrix/server/entity"
)

type mockGetResponse struct {
	total  int64
	valids int64
	err    error
}

type mockSetResponse struct {
	err error
}

type mockDB struct {
	getResponse mockGetResponse
	setResponse mockSetResponse
}

func (db mockDB) Get() (int64, int64, error) {
	return db.getResponse.total, db.getResponse.valids, db.getResponse.err
}

func (db mockDB) Set(seq []string, valid bool) error {
	return db.setResponse.err
}

func TestHandler_ReceiveSequence(t *testing.T) {
	ctx := context.Background()
	lg := log.Default()

	type fields struct {
		ctx  context.Context
		db   db.Database
		lg   *log.Logger
		body string
	}
	type response struct {
		resp    entity.PostResponse
		status  int
		wantErr bool
	}
	tests := []struct {
		name     string
		fields   fields
		response response
	}{
		{
			name: "Validate sequence",
			fields: fields{
				ctx:  ctx,
				lg:   lg,
				db:   mockDB{},
				body: `{"letters": ["DUHBHB", "DUBUHD", "UBUUHU", "BHBDHH", "DDDDUB", "UDBDUH"]}`,
			},
			response: response{
				resp: entity.PostResponse{
					IsValid: true,
				},
				status:  200,
				wantErr: false,
			},
		},
		{
			name: "Invalid sequence structure",
			fields: fields{
				ctx:  ctx,
				lg:   lg,
				body: `{"letter": ["DUHBHB", "DUBUHD", "UBUUHU", "BHBDHH", "DDDDUB", "UDBDUH"]}`,
			},
			response: response{
				status:  400,
				wantErr: true,
			},
		},
		{
			name: "Sequence validator error",
			fields: fields{
				ctx:  ctx,
				lg:   lg,
				db:   mockDB{},
				body: `{"letters": ["ABCDEF", "DUBUHD", "UBUUHU", "BHBDHH", "DDDDUB", "UDBDUH"]}`,
			},
			response: response{
				status:  400,
				wantErr: true,
			},
		},
		{
			name: "Error sending data to DB",
			fields: fields{
				ctx: ctx,
				lg:  lg,
				db: mockDB{
					setResponse: mockSetResponse{
						err: fmt.Errorf("invalid input"),
					},
				},
				body: `{"letters": ["DUHBHB", "DUBUHD", "UBUUHU", "BHBDHH", "DDDDUB", "UDBDUH"]}`,
			},
			response: response{
				status:  500,
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.fields.body)
			h := New(tt.fields.ctx, tt.fields.db, tt.fields.lg)

			srv := httptest.NewServer(h.ReceiveSequence())
			client := srv.Client()
			defer srv.Close()

			resp, err := client.Post(srv.URL, "application/json", bodyReader)
			if err != nil {
				t.Errorf("Error sending request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if tt.response.status != resp.StatusCode {
				t.Errorf("Expected status: %d, got: %d\n", tt.response.status, resp.StatusCode)
				return
			}

			if !tt.response.wantErr {
				respVal := &entity.PostResponse{}
				err = json.NewDecoder(resp.Body).Decode(respVal)
				if err != nil {
					t.Errorf("Unable to decode response: %s\n", err.Error())
					return
				}
				if respVal.IsValid != tt.response.resp.IsValid {
					t.Errorf("Response error, got: %v, expected: %v\n", respVal.IsValid, tt.response.resp.IsValid)
					return
				}
			}
		})
	}
}

func TestHandler_GetStats(t *testing.T) {
	ctx := context.Background()
	lg := log.Default()

	type fields struct {
		ctx context.Context
		db  db.Database
		lg  *log.Logger
	}
	type response struct {
		resp    entity.GetResponse
		status  int
		wantErr bool
	}
	tests := []struct {
		name     string
		fields   fields
		response response
	}{
		{
			name: "Get stats with values in DB",
			fields: fields{
				ctx: ctx,
				lg:  lg,
				db: mockDB{
					getResponse: mockGetResponse{
						total:  50,
						valids: 10,
						err:    nil,
					},
				},
			},
			response: response{
				resp: entity.GetResponse{
					Valid:   10,
					Invalid: 40,
					Ratio:   0.2,
				},
				status:  200,
				wantErr: false,
			},
		},
		{
			name: "Get stats without values in DB",
			fields: fields{
				ctx: ctx,
				lg:  lg,
				db: mockDB{
					getResponse: mockGetResponse{
						total:  0,
						valids: 0,
						err:    nil,
					},
				},
			},
			response: response{
				resp: entity.GetResponse{
					Valid:   0,
					Invalid: 0,
					Ratio:   0.0,
				},
				status:  200,
				wantErr: false,
			},
		},
		{
			name: "Error getting data from DB",
			fields: fields{
				ctx: ctx,
				lg:  lg,
				db: mockDB{
					getResponse: mockGetResponse{
						total:  0,
						valids: 0,
						err:    fmt.Errorf("connection error"),
					},
				},
			},
			response: response{
				resp:    entity.GetResponse{},
				status:  500,
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.fields.ctx, tt.fields.db, tt.fields.lg)

			srv := httptest.NewServer(h.GetStats())
			client := srv.Client()
			defer srv.Close()

			resp, err := client.Get(srv.URL)
			if err != nil {
				t.Errorf("Error sending request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if tt.response.status != resp.StatusCode {
				t.Errorf("Expected status: %d, got: %d\n", tt.response.status, resp.StatusCode)
				return
			}

			if !tt.response.wantErr {
				respVal := &entity.GetResponse{}
				err = json.NewDecoder(resp.Body).Decode(respVal)
				if err != nil {
					t.Errorf("Unable to decode response: %s\n", err.Error())
					return
				}
				if respVal.Valid != tt.response.resp.Valid ||
					respVal.Invalid != tt.response.resp.Invalid ||
					respVal.Ratio != tt.response.resp.Ratio {
					t.Errorf("Response error, got: %v, expected: %v\n", respVal, tt.response.resp)
					return
				}
			}
		})
	}
}

func TestHandlers_WrongMethodCalls(t *testing.T) {
	ctx := context.Background()
	lg := log.Default()

	t.Run("ReceiveSequence with GET", func(t *testing.T) {
		h := New(ctx, mockDB{}, lg)

		srv := httptest.NewServer(h.ReceiveSequence())
		client := srv.Client()
		defer srv.Close()

		resp, err := client.Get(srv.URL)
		if err != nil {
			t.Errorf("Error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status: %d, got: %d\n", http.StatusMethodNotAllowed, resp.StatusCode)
			return
		}
	})

	t.Run("GetStats with POST", func(t *testing.T) {
		h := New(ctx, mockDB{}, lg)

		srv := httptest.NewServer(h.GetStats())
		client := srv.Client()
		defer srv.Close()

		resp, err := client.Post(srv.URL, "application/json", strings.NewReader(`{}`))
		if err != nil {
			t.Errorf("Error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status: %d, got: %d\n", http.StatusMethodNotAllowed, resp.StatusCode)
			return
		}
	})
}
