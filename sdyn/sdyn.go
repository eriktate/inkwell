package sdyn

import (
	"github.com/eriktate/inkwell"
	"github.com/eriktate/inkwell/dynamo"
	"github.com/eriktate/inkwell/s3"
)

// Client acts as a container for both an s3 Client and a Dynamo client.
type Client struct {
	s   *s3.Client
	dyn *dynamo.Client

	blogService BlogService
}

// NewClient returns a new sdyn Client struct that will use the given s3 and
// dynamo clients.
func NewClient(s *s3.Client, dyn *dynamo.Client) Client {
	return Client{
		s:           s,
		dyn:         dyn,
		blogService: NewBlogService(s.BlogService(), dyn.BlogService()),
	}
}

// BlogService implements the inkwell.Client interface and returns an sdyn
// implementation of an inkwell.BlogService.
func (c Client) BlogService() inkwell.BlogService {
	return c.blogService
}
