package dynamo_test

import (
	"testing"

	"github.com/eriktate/inkwell"
	"github.com/eriktate/inkwell/dynamo"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {

}

func Test_NewBlogService(t *testing.T) {
	Convey("Given a dynamodb interface", t, func() {
		db := dynamo.NewMockDynamo()

		Convey("When building a new blog service", func() {
			svc := dynamo.NewBlogService(db, "blogs")
			Convey("We should get back a blog service that meets the inkwell interface", func() {
				var blogInterface inkwell.BlogService

				blogInterface = svc
				So(blogInterface, ShouldNotBeNil)
			})
		})
	})
}

func Test_GetBlog(t *testing.T) {
	Convey("Given a valid blog service", t, func() {
		dynBlog := inkwell.Blog{
			ID:        "test",
			Title:     "Test Blog",
			AuthorID:  "auth1",
			Content:   "A great blog",
			Published: true,
		}

		db := dynamo.NewMockDynamo()
		if err := db.AddItem("blogs", dynBlog); err != nil {
			log.WithError(err).Println("Failed to add item")
		}
		svc := dynamo.NewBlogService(db, "blogs")

		Convey("When attempt to get a blog", func() {
			blog, err := svc.Get("test")
			Convey("A non-zero blog should be returned", func() {
				So(blog.Title, ShouldNotEqual, "")
			})
			Convey("And no error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

/*
func Test_GetBlogsByAuthor(t *testing.T) {
	Convey("Given a valid blog service", t, func() {
		dynBlog := inkwell.Blog{
			ID:        "test",
			Title:     "Test Blog",
			AuthorID:  "auth1",
			Content:   "A great blog",
			Published: true,
		}

		db := dynamo.NewMockDynamo()
		if err := db.AddItem("blogs", dynBlog); err != nil {
			log.WithError(err).Println("Failed to add item")
		}
		svc := dynamo.NewBlogService(db, "blogs")
		Convey("When we attempt to get blogs by author", func() {
			blogs, err := svc.GetByAuthor("auth1")
			Convey("A non-empty slice of blogs should be returned", func() {
				So(blogs, ShouldNotBeEmpty)
			})
			Convey("And no error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
*/

func Test_WriteBlog(t *testing.T) {
	Convey("Given a valid blog service", t, func() {
		dynBlog := inkwell.Blog{
			ID:        "test",
			Title:     "Test Blog",
			AuthorID:  "auth1",
			Content:   "A great blog",
			Published: true,
		}

		db := dynamo.NewMockDynamo()
		svc := dynamo.NewBlogService(db, "blogs")

		Convey("When we attempt to write a blog", func() {
			err := svc.Write(dynBlog)
			Convey("A dynamo put should be called at least once", func() {
				So(db.PutItemCalled, ShouldBeGreaterThan, 0)
			})
			Convey("And no error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
