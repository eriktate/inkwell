package s3

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/eriktate/inkwell"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Client implements a BlogReadWriter in s3.
type Client struct {
	svc    s3iface.S3API
	bucket string
}

// NewClient returns a new s3 Client implementation.
func NewClient(logger *log.Logger, svc s3iface.S3API) Client {
	bucket := os.Getenv("INKWELL_BLOGS_BUCKET")
	if svc != nil {
		logger.Println("Building S3 client with existing S3API...")
		return Client{svc: svc, bucket: bucket}
	}

	logger.Println("Building S3 client with new AWS session...")
	return Client{
		svc:    s3.New(session.New(), aws.NewConfig().WithRegion("us-east-1")),
		bucket: bucket,
	}
}

// Get fetches a blog from S3 using a composite key of authorID and key.
func (c *Client) Get(authorID, key string) (inkwell.Blog, error) {
	var blog inkwell.Blog

	// Get the actual blog data.
	goi := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(jsonPath(authorID, key)),
	}

	goo, err := c.svc.GetObject(goi)
	if err != nil {
		return blog, errors.Wrap(err, "Failed to fetch blog from S3")
	}

	defer goo.Body.Close()

	// Get the ACL for the blog to determine published status.
	aclInput := &s3.GetObjectAclInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(jsonPath(authorID, key)),
	}

	aclOutput, err := c.svc.GetObjectAcl(aclInput)
	if err != nil {
		return blog, errors.Wrap(err, "Failed to fetch ACL from S3")
	}

	// Read object contents.
	data, err := ioutil.ReadAll(goo.Body)
	if err != nil {
		return blog, errors.Wrap(err, "Failed to read s3 object body")
	}

	// need to extract title data in a safe way. NOTE: S3 returns metadata with
	// the keys capitalized for some reason...
	var title string
	if t, ok := goo.Metadata["Title"]; ok {
		title = *t
	}

	blog.AuthorID = authorID
	blog.Key = key
	blog.Content = data
	blog.Title = title
	blog.Published = isPublic(aclOutput.Grants)
	blog.UpdatedAt = *goo.LastModified

	return blog, nil
}

// GetBlogsByAuthor TODO: Need to figure out what this should do in S3.
func (c *Client) GetBlogsByAuthor(authorID, key string) (inkwell.Blog, error) {
	return inkwell.Blog{}, nil
}

// Write uploads a Blog to S3 at the path authorID/key.json
func (c *Client) Write(blog inkwell.Blog) error {
	// blogs are private in s3 by default unless it's marked as published.
	acl := s3.ObjectCannedACLPrivate
	if blog.Published {
		acl = s3.ObjectCannedACLPublicRead
	}

	meta := map[string]*string{"title": aws.String(blog.Title)}

	poi := &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(jsonPath(blog.AuthorID, blog.Key)),
		Body:          bytes.NewReader(blog.Content),
		ContentLength: aws.Int64(int64(len(blog.Content))),
		ContentType:   aws.String("application/json"),
		ACL:           aws.String(acl),
		Metadata:      meta,
	}

	_, err := c.svc.PutObject(poi)
	if err != nil {
		return errors.Wrap(err, "Failed to put blog in S3")
	}

	return nil
}

// Publish sets the ACL for the blog in S3 to PublicRead.
func (c *Client) Publish(authorID, key string) error {
	poi := &s3.PutObjectAclInput{
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
		Bucket: aws.String(c.bucket),
		Key:    aws.String(jsonPath(authorID, key)),
	}

	if _, err := c.svc.PutObjectAcl(poi); err != nil {
		return errors.Wrap(err, "Failed to update ACL for blog post")
	}

	return nil
}

// Redact sets the ACL for the blog in S3 to Private.
func (c *Client) Redact(authorID, key string) error {
	poi := &s3.PutObjectAclInput{
		ACL:    aws.String(s3.ObjectCannedACLPrivate),
		Bucket: aws.String(c.bucket),
		Key:    aws.String(jsonPath(authorID, key)),
	}

	if _, err := c.svc.PutObjectAcl(poi); err != nil {
		return errors.Wrap(err, "Failed to update ACL for blog post")
	}

	return nil
}

// Delete removes a blog completely from S3.
func (c *Client) Delete(authorID, key string) error {
	doi := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(jsonPath(authorID, key)),
	}

	if _, err := c.svc.DeleteObject(doi); err != nil {
		return errors.Wrap(err, "Failed to remove blog from S3")
	}

	return nil
}

func jsonPath(authorID, key string) string {
	return fmt.Sprintf("%s/%s.json", authorID, key)
}

// isPublic cycles through a slice of grants and looks for one that gives the
// READ permission to a "Group" grantee.
// TODO: This function is a little more brittle than I would like.
// Should build a more robust version at some point.
func isPublic(grants []*s3.Grant) bool {
	for _, grant := range grants {
		if *grant.Grantee.Type == "Group" && *grant.Permission == "READ" {
			return true
		}
	}

	return false
}
