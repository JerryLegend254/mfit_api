package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JerryLegend254/mfit_api/internal/store"
	"github.com/JerryLegend254/mfit_api/internal/store/mocks"
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
		expectedBody       interface{}
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
				fmt.Sprintf(`{"error":%q}`, ErrBadRequest.Error()),
			},
		},
		{"should return 400 - missing name",
			[]byte(`{"image_url": "Test Image Url"}`),
			response{
				http.StatusBadRequest,
				fmt.Sprintf(`{"error":%q}`, ErrBadRequest.Error()),
			},
		},
		{"should return 400 - missing image_url",
			[]byte(`{"name": "Test Name"}`),
			response{
				http.StatusBadRequest,
				fmt.Sprintf(`{"error":%q}`, ErrBadRequest.Error()),
			},
		},
		{"should return 400 - invalid json",
			[]byte(`{"name": "Test Name", "image_url": "Test Image Url"`),
			response{
				http.StatusBadRequest,
				fmt.Sprintf(`{"error":%q}`, ErrBadRequest.Error()),
			},
		},
		{"should return 201",
			[]byte(`{"name": "Test Name", "image_url": "Test Image Url"}`),
			response{
				http.StatusCreated,
				`{"data":{"id": 0, "name": "Test Name", "image_url": "Test Image Url"}}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create the response
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/bodyparts", bytes.NewReader(tt.payload))
			res := httptest.NewRecorder()

			mux.ServeHTTP(res, req)

			// compare the status codes
			assertStatusCode(t, res.Code, int(tt.response.expectedStatusCode))

			fmt.Println("body", res.Body, tt.name)
			fmt.Println("expectedBody", tt.response.expectedBody)
			// compare the response bodies
			assertResponse(t, res.Body, tt.response.expectedBody)

		})
	}

}
