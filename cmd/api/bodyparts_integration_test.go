package main

import (
	"net/http"
	"testing"

	"github.com/JerryLegend254/mfit_api/internal/store"
)

func TestBodyPartIT(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	// seed the db for testing
	_, err := db.Exec("INSERT INTO body_part (name, image_url) VALUES ($1, $2);", "Test 1", "Image Url 1")
	if err != nil {
		t.Fatalf("error seeding bodypart table: %v", err)

	}

	s := store.NewStorage(db)

	app := newTestApplication(t, s)
	mux := app.mount()

	ts := []struct {
		name string
		req  *http.Request
		want []byte
	}{
		{"create bodypart", newPostBodyPartRequest([]byte(`{"name": "Test Name", "image_url": "Test Image Url"}`)), []byte(`{"data": {"id": 2, "name": "Test Name", "image_url": "Test Image Url"}}`)},
		{"get one bodypart", newGetBodyPartRequest(1), []byte(`{"data": {"id": 1, "name": "Test 1", "image_url": "Image Url 1"}}`)},
		{"get all bodyparts", newGetBodyPartsRequest(), []byte(`{"data": [{"id": 1, "name": "Test 1", "image_url": "Image Url 1"},{"id": 2, "name": "Test Name", "image_url": "Test Image Url"}]}`)},
		{"update bodypart", newPatchBodyPartRequest(2, []byte(`{"name": "Update Title", "image_url": "Updated Image Url"}`)), []byte(`{"data": {"id": 2, "name": "Update Title", "image_url": "Updated Image Url"}}`)},
		{"delete bodypart", newDeleteBodyPartRequest(1), nil},
	}

	for _, tt := range ts {
		t.Run(tt.name, func(t *testing.T) {
			res := execRequest(mux, tt.req)
			assertResponse(t, res.Body, tt.want)
		})
	}

}
