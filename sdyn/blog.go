package sdyn

import (
	"github.com/eriktate/inkwell"
	"github.com/pkg/errors"
)

// BlogService performs blog related actions in both dynamo and s3
type BlogService struct {
	s3svc  inkwell.BlogService
	dynsvc inkwell.BlogService
}

// NewBlogService returns a new sdyn BlogService given an s3 and dynamo
// BlogService.
func NewBlogService(s3svc inkwell.BlogService, dynsvc inkwell.BlogService) BlogService {
	return BlogService{
		s3svc:  s3svc,
		dynsvc: dynsvc,
	}
}

// Get will retrieve a Blog by reconciling dynamo and s3 data.
func (s BlogService) Get(authorID, blogID string) (inkwell.Blog, error) {
	var blog inkwell.Blog

	content, err := s.s3svc.Get(authorID, blogID)
	if err != nil {
		return blog, errors.Wrap(err, "Failed to retrieve Blog content from s3")
	}

	meta, err := s.dynsvc.Get(authorID, blogID)
	if err != nil {
		return blog, errors.Wrap(err, "Failed to retrieve Blog meta from dynamo")
	}

	blog = meta
	blog.Content = content.Content

	return blog, nil
}

// Write will write relevant partf of the given Blog to dynamo and s3.
func (s BlogService) Write(blog inkwell.Blog) error {
	if err := s.s3svc.Write(blog); err != nil {
		return errors.Wrap(err, "Failed to write blog to s3")
	}

	if err := s.dynsvc.Write(blog); err != nil {
		return errors.Wrap(err, "Failed to write blog to dynamo")
	}

	return nil
}

// Revise will overwrite the content of a blog in s3.
func (s BlogService) Revise(authorID, blogID, content string) error {
	if err := s.s3svc.Revise(authorID, blogID, content); err != nil {
		return errors.Wrap(err, "Failed to revise blog in s3")
	}

	if err := s.dynsvc.Revise(authorID, blogID, content); err != nil {
		return errors.Wrap(err, "Failed to revise blog in dynamo")
	}

	return nil
}

// Publish will mark the given blog as published in both dynamo and s3.
func (s BlogService) Publish(authorID, blogID string) error {
	if err := s.s3svc.Publish(authorID, blogID); err != nil {
		return errors.Wrap(err, "Failed to publish blog in s3")
	}

	if err := s.dynsvc.Publish(authorID, blogID); err != nil {
		return errors.Wrap(err, "Failed to publish blog in dynamo")
	}

	return nil
}

// Redact will mark the given blog as unpublished in both dynamo and s3.
func (s BlogService) Redact(authorID, blogID string) error {
	if err := s.s3svc.Redact(authorID, blogID); err != nil {
		return errors.Wrap(err, "Failed to redact blog in s3")
	}

	if err := s.dynsvc.Redact(authorID, blogID); err != nil {
		return errors.Wrap(err, "Failed to redact blog in dynamo")
	}

	return nil

}
