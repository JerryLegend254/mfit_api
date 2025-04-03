package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/JerryLegend254/mfit_api/internal/store"
	"github.com/JerryLegend254/mfit_api/internal/store/mocks"
)

var (
	BodyPartUrl = newCollectionPath("bodyparts")
)

func TestCreateBodyPart(t *testing.T) {
	mockBodyPartStore := new(mocks.MockBodyPartStore)

	store := store.Storage{
		BodyParts: mockBodyPartStore,
	}

	app := newTestApplication(t, store)
	mux := app.mount()

	type response struct {
		expectedStatusCode int
		expectedBody       []byte
	}
	tests := []struct {
		name     string
		payload  []byte
		response response
	}{
		{"should return 400 - empty body",
			[]byte(`{}`),
			response{
				http.StatusBadRequest,
				func() []byte {
					res := fmt.Sprintf(`{"error":%q}`, ErrBadRequest.Error())
					return []byte(res)
				}(),
			},
		},
		{"should return 400 - missing name",
			[]byte(`{"image_url": "Test Image Url"}`),
			response{
				http.StatusBadRequest,
				func() []byte {
					res := fmt.Sprintf(`{"error":%q}`, ErrBadRequest.Error())
					return []byte(res)
				}(),
			},
		},
		{"should return 400 - missing image_url",
			[]byte(`{"name": "Test Name"}`),
			response{
				http.StatusBadRequest,
				func() []byte {
					res := fmt.Sprintf(`{"error":%q}`, ErrBadRequest.Error())
					return []byte(res)
				}(),
			},
		},
		{"should return 400 - invalid json",
			[]byte(`{"name": "Test Name", "image_url": "Test Image Url"`),
			response{
				http.StatusBadRequest,
				func() []byte {
					res := fmt.Sprintf(`{"error":%q}`, ErrBadRequest.Error())
					return []byte(res)
				}(),
			},
		},
		{"should return 201",
			[]byte(`{"name": "Test Name", "image_url": "Test Image Url"}`),
			response{
				http.StatusCreated,
				[]byte(`{"data":{"id": 0, "name": "Test Name", "image_url": "Test Image Url"}}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create the response
			req := newPostBodyPartRequest(tt.payload)

			res := execRequest(mux, req)

			// compare the status codes
			assertStatusCode(t, res.Code, int(tt.response.expectedStatusCode))

			// compare the response bodies
			assertResponse(t, res.Body, tt.response.expectedBody)

			// compare the content type header
			assertContentType(t, res.Header().Get("content-type"), jsonContentType)

		})
	}

}

func TestGetBodyPart(t *testing.T) {
	mockBodyPartStore := new(mocks.MockBodyPartStore)

	store := store.Storage{
		BodyParts: mockBodyPartStore,
	}

	app := newTestApplication(t, store)
	mux := app.mount()
	t.Run("valid route path", func(t *testing.T) {
		req := newGetBodyPartRequest(1)

		res := execRequest(mux, req)

		assertStatusCode(t, res.Code, http.StatusOK)

		assertContentType(t, res.Result().Header.Get("content-type"), jsonContentType)

	})

}

func newPostBodyPartRequest(payload []byte) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, BodyPartUrl, bytes.NewReader(payload))
	return req
}

func newGetBodyPartRequest(id int64) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", BodyPartUrl, id), nil)
	return req
}

func newGetBodyPartsRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, BodyPartUrl, nil)
	return req
}

func newPatchBodyPartRequest(id int64, payload []byte) *http.Request {
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%d", BodyPartUrl, id), bytes.NewReader(payload))
	return req
}

func newDeleteBodyPartRequest(id int64) *http.Request {
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/%d", BodyPartUrl, id), nil)
	return req
}
