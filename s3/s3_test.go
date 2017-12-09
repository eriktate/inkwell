package s3_test

import (
	"testing"

	"github.com/eriktate/inkwell"
	"github.com/eriktate/inkwell/s3"
	log "github.com/sirupsen/logrus"
)

func Test_Integration_Get(t *testing.T) {
	// SETUP
	log.Println("Running integration test for Get")
	client := s3.NewClient(log.New(), nil)
	authorID := "int-test-user"
	key := "int-test-blog"
	content := "This is a test"
	title := "Int Test Blog"

	testBlog := inkwell.Blog{
		AuthorID: authorID,
		Key:      key,
		Title:    title,
		Content:  []byte(content),
	}

	if err := client.Write(testBlog); err != nil {
		t.Fatalf("Failed to write int test blog")
	}

	// RUN
	log.Println("About to run test")
	blog, err := client.Get(authorID, key)

	// CLEANUP
	// if err := client.Delete(authorID, key); err != nil {
	// 	t.Fatalf("Failed to cleanup integration test blog: %s", err)
	// }

	// ASSERT
	if err != nil {
		t.Fatalf("Failed to fetch blog from s3: %s", err)
	}

	if string(blog.Content) != content {
		t.Fatal("Fetched content does not match what was saved")
	}

	if blog.Title != title {
		t.Fatal("Fetched title does not match what was saved")
	}
}

func Test_Integration_Publish(t *testing.T) {
	// SETUP
	log.Println("Running integration test for Get")
	client := s3.NewClient(log.New(), nil)
	authorID := "int-test-user"
	key := "int-test-blog"
	content := "This is a test"
	title := "Int Test Blog"

	testBlog := inkwell.Blog{
		AuthorID: authorID,
		Key:      key,
		Title:    title,
		Content:  []byte(content),
	}

	if err := client.Write(testBlog); err != nil {
		t.Fail()
	}

	// RUN
	err := client.Publish(authorID, key)
	newBlog, gErr := client.Get(authorID, key)
	if gErr != nil {
		t.Fatalf("Failed to fetch published test blog: %s", err)
	}

	// CLEANUP
	if err := client.Delete(authorID, key); err != nil {
		t.Fatalf("Failed to cleanup integration test blog: %s", err)
	}

	// ASSERT
	if err != nil {
		t.Fatalf("Failed to fetch blog from s3: %s", err)
	}

	if !newBlog.Published {
		t.Fatal("Updated blog was not set to published")
	}
}

func Test_Integration_Redact(t *testing.T) {
	// SETUP
	log.Println("Running integration test for Get")
	client := s3.NewClient(log.New(), nil)
	authorID := "int-test-user"
	key := "int-test-blog"
	content := "This is a test"
	title := "Int Test Blog"

	testBlog := inkwell.Blog{
		AuthorID:  authorID,
		Key:       key,
		Title:     title,
		Published: true,
		Content:   []byte(content),
	}

	if err := client.Write(testBlog); err != nil {
		t.Fail()
	}

	// RUN
	err := client.Redact(authorID, key)
	newBlog, gErr := client.Get(authorID, key)
	if gErr != nil {
		t.Fatalf("Failed to fetch published test blog: %s", err)
	}

	// CLEANUP
	if err := client.Delete(authorID, key); err != nil {
		t.Fatalf("Failed to cleanup integration test blog: %s", err)
	}

	// ASSERT
	if err != nil {
		t.Fatalf("Failed to fetch blog from s3: %s", err)
	}

	if newBlog.Published {
		t.Fatal("Updated blog was not set to redacted")
	}
}
