## HTTP Endpoints

### AuthorHandler
- GET /author/:authorID/blog/:key | Retrieves a single blog.
- GET /author/:authorID/blog | Retrieves all blogs for an author (probably just metadata)
- POST /author | Create a new author (?)

### BlogHandler
- POST /author/:authorID/blog | Write a blog
- POST /author/:authorID/blog/:key/publish | Publish a blog
- POST /author/:authorID/blog/:key/redact | Redact a blog
- POST /author/:authorID/blog/:key/content | Set new content
- POST /author/:authorID/blog/:key/title | Set new title
