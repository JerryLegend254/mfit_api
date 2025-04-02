package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/JerryLegend254/mfit_api/internal/store"
)

func TestBodyPartStore(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	// seed the db for testing
	_, err := db.Exec("INSERT INTO body_part (name, image_url) VALUES ($1, $2);", "Test 1", "Image Url 1")
	if err != nil {
		t.Fatalf("error seeding bodypart table: %v", err)

	}

	s := store.NewStorage(db)

	app := newTestApplication(t, s)
	tests := []struct {
		name string
		args []interface{}
		want interface{}
		f    interface{}
	}{
		{"create bodypart", []interface{}{context.Background(), &store.BodyPart{Name: "Test Name", ImageUrl: "Test Image Url"}}, &store.BodyPart{ID: 2, Name: "Test Name", ImageUrl: "Test Image Url"}, app.store.BodyParts.Create},
		{"get on bodypart", []interface{}{context.Background(), int64(1)}, &store.BodyPart{ID: 1, Name: "Test 1", ImageUrl: "Image Url 1"}, app.store.BodyParts.GetByID},
		{"get all bodypart", []interface{}{context.Background()}, []store.BodyPart{{ID: 1, Name: "Test 1", ImageUrl: "Image Url 1"}, {ID: 2, Name: "Test Name", ImageUrl: "Test Image Url"}}, app.store.BodyParts.GetAll},
		{"update bodypart", []interface{}{context.Background(), &store.BodyPart{ID: 2, Name: "Updated Test Name", ImageUrl: "Test Image Url"}}, nil, app.store.BodyParts.Update},
		{"delete update bodypart", []interface{}{context.Background(), int64(1)}, nil, app.store.BodyParts.Delete},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			results := invokeFunction(tt.f, tt.args)

			for _, result := range results {
				if result.CanInterface() && result.Interface() != nil {

					if !reflect.DeepEqual(result.Interface(), tt.want) {
						t.Errorf("wanted %v but got %v", tt.want, result.Interface())

					}
				}
			}

		})

	}

}

func invokeFunction(fn interface{}, args []interface{}) []reflect.Value {
	fnValue := reflect.ValueOf(fn)
	fnArgs := make([]reflect.Value, len(args))

	for i, arg := range args {
		fnArgs[i] = reflect.ValueOf(arg)
	}

	return fnValue.Call(fnArgs)
}
