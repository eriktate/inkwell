package inkwell_test

import (
	"log"
	"testing"

	"github.com/eriktate/inkwell"
)

// func Test_Toml(t *testing.T) {
// 	testToml := `
// 		# This is a test post configuration

// 		id = "54321"
// 		title = "Example post"

// 		[[row]]
// 			id = "12345"
// 			[[row.column]]
// 				type = "text"
// 				[[row.column.element]]
// 					id = "someid"
// 					text = "Some text..?"
// 		`

// 	var post inkwell.Post

// 	if _, err := toml.Decode(testToml, &post); err != nil {
// 		t.Fatalf("Failed to decode toml: %s", err)
// 	}
// }

func Test_Render(t *testing.T) {
	testJson := `
	{
		"id": "54321",
		"title": "Example post",
		"rows": [
		{
			"id": "12345",
			"columns": [
			{
				"id": "56789",
				"contents": [
				{
					"type": "paragraph",
					"element": {
						"id": "someid",
						"text": "Some text..?"
					}
				}
				]
			}
			]
		}
		]
	}
	`

	post, err := inkwell.BuildPost([]byte(testJson))
	if err != nil {
		t.Fatalf("Failed to decode json: %s", err)
	}

	log.Println(inkwell.RenderPost(post))
}
