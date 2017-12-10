package mock

import "github.com/eriktate/inkwell"

// MockAuthorReadWriter mocks an AuthorReadWriter.
type MockAuthorReadWriter struct {
	GetFn     func(authorID string) (inkwell.Author, error)
	GetCalled int

	WriteFn     func(author inkwell.Author) error
	WriteCalled int

	PassThru bool
}

// Get implements the AuthorReader interface.
func (m *MockAuthorReadWriter) Get(authorID string) (inkwell.Author, error) {
	m.GetCalled++
	if m.PassThru {
		return m.GetFn(authorID)
	}

	return inkwell.Author{AuthorID: authorID}, nil
}

// Write implements the AuthorWriter interface.
func (m *MockAuthorReadWriter) Write(author inkwell.Author) error {
	m.WriteCalled++
	if m.PassThru {
		return m.WriteFn(author)
	}

	return nil
}
