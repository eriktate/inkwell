package mock

import "github.com/eriktate/inkwell"

type MockBlogService struct {
	GetCalled   int
	GetPassthru func(blogID string) (inkwell.Blog, error)

	GetByAuthorCalled   int
	GetByAuthorPassthru func(authorID string) ([]inkwell.Blog, error)

	WriteCalled   int
	WritePassthru func(blog inkwell.Blog) error

	UpdateCalled   int
	UpdatePassthru func(blog inkwell.Blog) error

	DeleteCalled   int
	DeletePassthru func(blogID string) error

	passthru bool
}

// NewMockBlogService returns a mocked out version of the inkwell BlogService.
func NewMockBlogService(passthru bool) *MockBlogService {
	return &MockBlogService{passthru: passthru}
}

func (m *MockBlogService) Get(blogID string) (inkwell.Blog, error) {
	m.GetCalled += 1
	if m.passThru {
		return m.GetPassthru(blogID)
	}

	return &Blog{}, nil
}

func (m *MockBlogService) GetByAuthor(authorID string) ([]inkwell.Blog, error) {
	m.GetByAuthorCalled += 1
	if m.passThru {
		return m.GetByAuthorPassthru(authorID)
	}

	return []inkwell.Blog{}, nil
}

func (m *MockBlogService) Write(blog inkwell.Blog) error {
	m.WriteCalled += 1
	if m.passThru {
		return m.WritePassthru(blog)
	}

	return nil
}

func (m *MockBlogService) Update(blog inkwell.Blog) error {
	m.UpdateCalled += 1
	if m.passThru {
		return m.UpdatePassthru(blog)
	}

	return nil
}

func (m *MockBlogService) Delete(blogID string) error {
	m.DeleteCalled += 1
	if m.passThru {
		return m.DeletePassthru(blog)
	}

	return nil
}
