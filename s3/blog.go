package s3

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/eriktate/inkwell"
)

type BlogService struct {
	inkwell.BlogService
	svc    s3iface.S3API
	bucket string
}

func NewBlogService(svc s3iface.S3API, bucket string) BlogService {
	return BlogService{
		svc:    svc,
		bucket: bucket,
	}
}

// Get returns a blog with content from s3.
func (s BlogService) Get(authorID, blogID string) (inkwell.Blog, error) {
	var blog inkwell.Blog

	gbo, err := s.svc.GetObject(s.getBlogInput(fmt.Sprintf("%s/%s", authorID, blogID)))
	if err != nil {
		return blog, err
	}

	defer gbo.Body.Close()

	data, err := ioutil.ReadAll(gbo.Body)
	if err != nil {
		return blog, err
	}

	blog.AuthorID = authorID
	blog.ID = blogID
	blog.Content = string(data)

	return blog, nil
}

func (s BlogService) Write(blog inkwell.Blog) error {
	_, err := s.svc.PutObject(s.putBlogInput(blog))
	if err != nil {
		return err
	}

	// TODO: Maybe do something with the result of PutObject?

	return nil
}

func (s BlogService) Revise(authorID, blogID, content string) error {
	path := fmt.Sprintf("%s/%s", authorID, blogID)
	_, err := s.svc.PutObject(s.reviseBlogInput(path, content))
	return err
}

func (s BlogService) Publish(authorID, blogID string) error {
	path := fmt.Sprintf("%s/%s", authorID, blogID)
	_, err := s.svc.PutObjectAcl(s.publishBlogInput(path))
	return err
}

func (s BlogService) Redact(authorID, blogID string) error {
	path := fmt.Sprintf("%s/%s", authorID, blogID)
	_, err := s.svc.PutObjectAcl(s.redactBlogInput(path))
	return err
}

func (s BlogService) Delete(authorID, blogID string) error {
	path := fmt.Sprintf("%s/%s", authorID, blogID)
	_, err := s.svc.DeleteObject(s.deleteBlogInput(path))
	return err
}

func (s BlogService) getBlogInput(path string) *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s.md", path)),
	}
}

func (s BlogService) putBlogInput(blog inkwell.Blog) *s3.PutObjectInput {
	path := fmt.Sprintf("%s/%s/", blog.AuthorID, blog.ID)
	content := []byte(blog.Content)
	body := bytes.NewReader(content)
	acl := s3.ObjectCannedACLPrivate

	if blog.Published {
		acl = s3.ObjectCannedACLPublicRead
	}

	if blog.Published {
		// TODO: The ACL on this item should be different depending on whether the blog is published.

	}

	return &s3.PutObjectInput{
		Bucket:          aws.String(s.bucket),
		Key:             aws.String(fmt.Sprintf("%sindex.html", path)),
		Body:            body,
		ContentLength:   aws.Int64(int64(len(content))),
		ContentEncoding: aws.String("utf-8"),
		ContentType:     aws.String("text/html"),
		ACL:             aws.String(acl),
	}
}

func (s BlogService) reviseBlogInput(path, newContent string) *s3.PutObjectInput {
	content := []byte(newContent)
	body := bytes.NewReader(content)

	return &s3.PutObjectInput{
		Bucket:          aws.String(s.bucket),
		Key:             aws.String(fmt.Sprintf("%s.html", path)),
		Body:            body,
		ContentLength:   aws.Int64(int64(len(content))),
		ContentEncoding: aws.String("utf-8"),
		ContentType:     aws.String("text/html"),
	}
}

func (s BlogService) publishBlogInput(path string) *s3.PutObjectAclInput {
	return &s3.PutObjectAclInput{
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s.html", path)),
	}
}

func (s BlogService) redactBlogInput(path string) *s3.PutObjectAclInput {
	return &s3.PutObjectAclInput{
		ACL:    aws.String(s3.ObjectCannedACLPrivate),
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s.html", path)),
	}
}

func (s BlogService) deleteBlogInput(path string) *s3.DeleteObjectInput {
	return &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s.html", path)),
	}
}
