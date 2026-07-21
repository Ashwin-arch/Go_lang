package schema

import (
	"errors"

	"github.com/graphql-go/graphql"
	"graphql-api/models"
)

// BuildSchema creates the GraphQL schema for Books and Authors with queries, mutations, and resolvers
func BuildSchema(store *models.Store) (graphql.Schema, error) {
	var authorType *graphql.Object
	var bookType *graphql.Object

	authorType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Author",
		Description: "Represents a book author",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.ID),
					Description: "Unique author ID",
				},
				"name": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the author",
				},
				"bio": &graphql.Field{
					Type:        graphql.String,
					Description: "Short biography of the author",
				},
				"books": &graphql.Field{
					Type:        graphql.NewList(bookType),
					Description: "List of books written by this author",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						author, ok := p.Source.(models.Author)
						if !ok {
							return nil, nil
						}
						return store.GetBooksByAuthorID(author.ID), nil
					},
				},
			}
		}),
	})

	bookType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Book",
		Description: "Represents a book entity in the library",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.ID),
					Description: "Unique book ID",
				},
				"title": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Title of the book",
				},
				"authorId": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.ID),
					Description: "Author ID associated with the book",
				},
				"price": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.Float),
					Description: "Price of the book in USD",
				},
				"author": &graphql.Field{
					Type:        authorType,
					Description: "Author object associated with this book",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						book, ok := p.Source.(models.Book)
						if !ok {
							return nil, nil
						}
						author, exists := store.GetAuthorByID(book.AuthorID)
						if !exists {
							return nil, nil
						}
						return author, nil
					},
				},
			}
		}),
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"books": &graphql.Field{
				Type:        graphql.NewList(bookType),
				Description: "Get all books with optional filters for authorId and maxPrice",
				Args: graphql.FieldConfigArgument{
					"authorId": &graphql.ArgumentConfig{Type: graphql.String},
					"maxPrice": &graphql.ArgumentConfig{Type: graphql.Float},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					authorID, _ := p.Args["authorId"].(string)
					maxPrice, _ := p.Args["maxPrice"].(float64)
					return store.GetAllBooks(authorID, maxPrice), nil
				},
			},
			"book": &graphql.Field{
				Type:        bookType,
				Description: "Get a book by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					book, exists := store.GetBookByID(id)
					if !exists {
						return nil, errors.New("book not found")
					}
					return book, nil
				},
			},
			"authors": &graphql.Field{
				Type:        graphql.NewList(authorType),
				Description: "Get all authors",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.GetAllAuthors(), nil
				},
			},
			"author": &graphql.Field{
				Type:        authorType,
				Description: "Get an author by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, _ := p.Args["id"].(string)
					author, exists := store.GetAuthorByID(id)
					if !exists {
						return nil, errors.New("author not found")
					}
					return author, nil
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createBook": &graphql.Field{
				Type:        bookType,
				Description: "Create a new book entry",
				Args: graphql.FieldConfigArgument{
					"title":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"authorId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"price":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					title := p.Args["title"].(string)
					authorID := p.Args["authorId"].(string)
					price := p.Args["price"].(float64)
					return store.CreateBook(title, authorID, price)
				},
			},
			"updateBook": &graphql.Field{
				Type:        bookType,
				Description: "Update an existing book entry",
				Args: graphql.FieldConfigArgument{
					"id":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"title":    &graphql.ArgumentConfig{Type: graphql.String},
					"authorId": &graphql.ArgumentConfig{Type: graphql.ID},
					"price":    &graphql.ArgumentConfig{Type: graphql.Float},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)
					var titlePtr *string
					if val, ok := p.Args["title"].(string); ok {
						titlePtr = &val
					}
					var authorIDPtr *string
					if val, ok := p.Args["authorId"].(string); ok {
						authorIDPtr = &val
					}
					var pricePtr *float64
					if val, ok := p.Args["price"].(float64); ok {
						pricePtr = &val
					}
					return store.UpdateBook(id, titlePtr, authorIDPtr, pricePtr)
				},
			},
			"deleteBook": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a book by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)
					err := store.DeleteBook(id)
					if err != nil {
						return false, err
					}
					return true, nil
				},
			},
			"createAuthor": &graphql.Field{
				Type:        authorType,
				Description: "Create a new author entry",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"bio":  &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					name := p.Args["name"].(string)
					bio, _ := p.Args["bio"].(string)
					return store.CreateAuthor(name, bio)
				},
			},
			"updateAuthor": &graphql.Field{
				Type:        authorType,
				Description: "Update an existing author entry",
				Args: graphql.FieldConfigArgument{
					"id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"name": &graphql.ArgumentConfig{Type: graphql.String},
					"bio":  &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)
					var namePtr *string
					if val, ok := p.Args["name"].(string); ok {
						namePtr = &val
					}
					var bioPtr *string
					if val, ok := p.Args["bio"].(string); ok {
						bioPtr = &val
					}
					return store.UpdateAuthor(id, namePtr, bioPtr)
				},
			},
			"deleteAuthor": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete an author by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)
					err := store.DeleteAuthor(id)
					if err != nil {
						return false, err
					}
					return true, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}
