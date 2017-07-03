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
	svc    s3iface.S3API
	bucket string
}

func NewBlogService(svc s3iface.S3API, bucket string) BlogService {
	return BlogService{
		svc:    svc,
		bucket: bucket,
	}
}

// GetBlog returns a blog with content from s3.
func (s BlogService) Get(authorID, blogID string) (inkwell.Blog, error) {
	var blog inkwell.Blog

	gbo, err := s.svc.GetObject(s.getBlogInput(fmt.Sprintf("%s/%s", authorID, blogID)))
	if err != nil {
		return blog, err
	}

	data, err := ioutil.ReadAll(gbo.Body)
	if err != nil {
		return blog, err
	}

	blog.AuthorID = authorID
	blog.BlogID = blogID
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

func (s BlogService) Publish(authorID, blogID string) error {
	_, err := s.svc.PutObjectAcl(putBlogAcl(fmt.Sprintf("%s/%s", authorID, blogID)))
	return err
}

func (s BlogService) getBlogInput(path string) *s3.GetObjetInput {
	return &s3.GetObjectInput{
		BucketName: aws.String(s.bucket),
		Key:        aws.String(fmt.Sprintf("%s.md", path)),
	}, nil
}

func (s BlogService) putBlogInput(blog inkwell.Blog) *s3.PutObjectInput {
	path := fmt.Sprintf("%s/%s", blog.AuthorID, blog.BlogID)
	content := []byte(blog.Content)
	body := bytes.NewReader(content)

	if blog.Published {
		// TODO: The ACL on this item should be different depending on whether the blog is published.

	}

	return &s3.PutObjectInput{
		BucketName:      aws.String(s.bucket),
		Key:             aws.String(fmt.Sprintf("%s.md", path)),
		Body:            body,
		ContentLength:   int64(len(content)),
		ContentEncoding: aws.String("utf-8"),
	}
}

func (s BlogService) putBlogAcl(path string, published bool) *s3.PutObjectAclInput {
	return &s3.PutObjectAclInput{}
}
