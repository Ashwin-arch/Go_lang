package schema_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"graphql-api/models"
	"graphql-api/schema"
)

func executeQuery(query string, gqlSchema graphql.Schema) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:        gqlSchema,
		RequestString: query,
	})
}

func TestGraphQLQueriesAndMutations(t *testing.T) {
	store := models.NewStore()
	gqlSchema, err := schema.BuildSchema(store)
	if err != nil {
		t.Fatalf("Failed to build GraphQL schema: %v", err)
	}

	// 1. Test Fetching Books with Nested Author
	t.Run("Query Books with Nested Author", func(t *testing.T) {
		q := `{
			books {
				id
				title
				price
				author {
					name
				}
			}
		}`
		res := executeQuery(q, gqlSchema)
		if len(res.Errors) > 0 {
			t.Fatalf("Unexpected query errors: %v", res.Errors)
		}
		data, ok := res.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Invalid response data structure")
		}
		books, ok := data["books"].([]interface{})
		if !ok || len(books) != 3 {
			t.Fatalf("Expected 3 books, got %v", books)
		}
	})

	// 2. Test Fetching Single Book by ID
	t.Run("Query Book By ID", func(t *testing.T) {
		q := `{
			book(id: "1") {
				id
				title
				author {
					name
				}
			}
		}`
		res := executeQuery(q, gqlSchema)
		if len(res.Errors) > 0 {
			t.Fatalf("Unexpected query errors: %v", res.Errors)
		}
		data := res.Data.(map[string]interface{})
		book := data["book"].(map[string]interface{})
		if book["title"] != "The Go Programming Language" {
			t.Errorf("Expected title 'The Go Programming Language', got '%v'", book["title"])
		}
	})

	// 3. Test Create Author Mutation
	t.Run("Mutation Create Author", func(t *testing.T) {
		q := `mutation {
			createAuthor(name: "Rob Pike", bio: "Co-creator of Go programming language.") {
				id
				name
				bio
			}
		}`
		res := executeQuery(q, gqlSchema)
		if len(res.Errors) > 0 {
			t.Fatalf("Unexpected mutation errors: %v", res.Errors)
		}
		data := res.Data.(map[string]interface{})
		author := data["createAuthor"].(map[string]interface{})
		if author["name"] != "Rob Pike" {
			t.Errorf("Expected name 'Rob Pike', got '%v'", author["name"])
		}
	})

	// 4. Test Create Book Mutation
	t.Run("Mutation Create Book", func(t *testing.T) {
		q := `mutation {
			createBook(title: "Go in Action", authorId: "1", price: 34.99) {
				id
				title
				price
				author {
					name
				}
			}
		}`
		res := executeQuery(q, gqlSchema)
		if len(res.Errors) > 0 {
			t.Fatalf("Unexpected mutation errors: %v", res.Errors)
		}
		data := res.Data.(map[string]interface{})
		book := data["createBook"].(map[string]interface{})
		if book["title"] != "Go in Action" {
			t.Errorf("Expected title 'Go in Action', got '%v'", book["title"])
		}
	})

	// 5. Test Update Book Mutation
	t.Run("Mutation Update Book", func(t *testing.T) {
		q := `mutation {
			updateBook(id: "1", price: 42.50) {
				id
				price
			}
		}`
		res := executeQuery(q, gqlSchema)
		if len(res.Errors) > 0 {
			t.Fatalf("Unexpected mutation errors: %v", res.Errors)
		}
		data := res.Data.(map[string]interface{})
		book := data["updateBook"].(map[string]interface{})
		if book["price"] != 42.50 {
			t.Errorf("Expected price 42.50, got %v", book["price"])
		}
	})

	// 6. Test Delete Book Mutation
	t.Run("Mutation Delete Book", func(t *testing.T) {
		q := `mutation {
			deleteBook(id: "2")
		}`
		res := executeQuery(q, gqlSchema)
		if len(res.Errors) > 0 {
			t.Fatalf("Unexpected mutation errors: %v", res.Errors)
		}
		data := res.Data.(map[string]interface{})
		success := data["deleteBook"].(bool)
		if !success {
			t.Errorf("Expected deleteBook to return true")
		}
	})
}
