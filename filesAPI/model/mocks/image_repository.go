package mocks

import (
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// MockImageRepository is a mock type for model.ImageRepository
type MockImageRepository struct {
	mock.Mock
}

// Save is mock of representations of ImageRepository Save
func (m *MockImageRepository) Save(imageType string, imageFile multipart.File) (*model.ImageInfo, error) {
	// args that will be passed to "Return" in the tests, when function
	// is called with a uid. Hence the name "ret"
	ret := m.Called(imageType, imageFile)

	// first value passed to "Return"
	var r0 *model.ImageInfo
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*model.ImageInfo)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}