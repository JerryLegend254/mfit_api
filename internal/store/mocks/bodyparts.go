package mocks

import (
	"context"

	"github.com/JerryLegend254/mfit_api/internal/store"
)

type MockBodyPartStore struct {
}

func (m *MockBodyPartStore) Create(context.Context, *store.BodyPart) error {
	return nil
}

func (m *MockBodyPartStore) GetByID(context.Context, int64) (*store.BodyPart, error) {
	return nil, nil
}

func (m *MockBodyPartStore) GetAll(context.Context) ([]store.BodyPart, error) {
	return nil, nil
}

func (m *MockBodyPartStore) Update(context.Context, *store.BodyPart) error {
	return nil
}

func (m *MockBodyPartStore) Delete(context.Context, int64) error {
	return nil
}
