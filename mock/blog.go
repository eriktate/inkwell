package mock

import (
	"errors"

	"github.com/eriktate/inkwell"
)

type MockBlogReadWriter struct {
	GetFn     func(authorID, key string) (inkwell.Blog, error)
	GetCalled int

	GetByAuthorFn     func(authorID string) ([]inkwell.Blog, error)
	GetByAuthorCalled int

	WriteFn     func(blog inkwell.Blog) error
	WriteCalled int

	SetKeyFn     func(authorID, key, newKey string) error
	SetKeyCalled int

	SetContentFn     func(authorID, key string, content []byte) error
	SetContentCalled int

	SetTitleFn     func(authorID, key, title string) error
	SetTitleCalled int

	PublishFn     func(authorID, key string) error
	PublishCalled int

	RedactFn     func(authorID, key string) error
	RedactCalled int

	DeleteFn     func(authorID, key string) error
	DeleteCalled int

	PassThru bool
	Fail     bool
}

func (m *MockBlogReadWriter) Get(authorID, key string) (inkwell.Blog, error) {
	m.GetCalled++
	if m.PassThru {
		return m.GetFn(authorID, key)
	}

	if m.Fail {
		return inkwell.Blog{}, errors.New("Mock failure")
	}

	return inkwell.Blog{AuthorID: authorID, Key: key}, nil
}

func (m *MockBlogReadWriter) GetByAuthor(authorID string) ([]inkwell.Blog, error) {
	m.GetByAuthorCalled++
	if m.PassThru {
		return m.GetByAuthor(authorID)
	}

	if m.Fail {
		return nil, errors.New("Mock failure")
	}

	return []inkwell.Blog{}, nil
}

func (m *MockBlogReadWriter) Write(blog inkwell.Blog) error {
	m.WriteCalled++
	if m.PassThru {
		return m.WriteFn(blog)
	}

	if m.Fail {
		return errors.New("Mock failure")
	}

	return nil
}

func (m *MockBlogReadWriter) SetKey(authorID, key, newKey string) error {
	m.SetKeyCalled++
	if m.PassThru {
		return m.SetKeyFn(authorID, key, newKey)
	}

	if m.Fail {
		return errors.New("Mock failure")
	}

	return nil
}

func (m *MockBlogReadWriter) SetContent(authorID, key string, content []byte) error {
	m.SetContentCalled++
	if m.PassThru {
		return m.SetContentFn(authorID, key, content)
	}

	if m.Fail {
		return errors.New("Mock failure")
	}

	return nil
}

func (m *MockBlogReadWriter) SetTitle(authorID, key, title string) error {
	m.SetTitleCalled++
	if m.PassThru {
		return m.SetTitleFn(authorID, key, title)
	}

	if m.Fail {
		return errors.New("Mock failure")
	}

	return nil
}

func (m *MockBlogReadWriter) Publish(authorID, key string) error {
	m.PublishCalled++
	if m.PassThru {
		return m.PublishFn(authorID, key)
	}

	if m.Fail {
		return errors.New("Mock failure")
	}

	return nil
}

func (m *MockBlogReadWriter) Redact(authorID, key string) error {
	m.RedactCalled++
	if m.PassThru {
		return m.RedactFn(authorID, key)
	}

	if m.Fail {
		return errors.New("Mock failure")
	}

	return nil
}

func (m *MockBlogReadWriter) Delete(authorID, key string) error {
	m.DeleteCalled++
	if m.PassThru {
		return m.DeleteFn(authorID, key)
	}

	if m.Fail {
		return errors.New("Mock failure")
	}

	return nil
}
