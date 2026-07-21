package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
)

// GraphQLHandler manages GraphQL HTTP requests and serves the GraphiQL IDE interface
type GraphQLHandler struct {
	Schema graphql.Schema
}

// NewGraphQLHandler constructs a GraphQLHandler instance
func NewGraphQLHandler(schema graphql.Schema) *GraphQLHandler {
	return &GraphQLHandler{Schema: schema}
}

type requestBody struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Allow CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Serve GraphiQL web UI if GET request has no query string or accepts text/html
	if r.Method == http.MethodGet {
		queryParam := r.URL.Query().Get("query")
		acceptHeader := r.Header.Get("Accept")
		if queryParam == "" || strings.Contains(acceptHeader, "text/html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(graphiqlHTML))
			return
		}
	}

	var reqBody requestBody

	if r.Method == http.MethodPost {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
				http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
				return
			}
		}
	} else if r.Method == http.MethodGet {
		reqBody.Query = r.URL.Query().Get("query")
		reqBody.OperationName = r.URL.Query().Get("operationName")
		varsStr := r.URL.Query().Get("variables")
		if varsStr != "" {
			_ = json.Unmarshal([]byte(varsStr), &reqBody.Variables)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         h.Schema,
		RequestString:  reqBody.Query,
		OperationName:  reqBody.OperationName,
		VariableValues: reqBody.Variables,
		Context:        r.Context(),
	})

	w.Header().Set("Content-Type", "application/json")
	if len(result.Errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(result)
}

// Embedded GraphiQL IDE HTML interface using CDN resources
const graphiqlHTML = `<!DOCTYPE html>
<html>
  <head>
    <title>GraphiQL IDE - Go GraphQL API</title>
    <link rel="stylesheet" href="https://unpkg.com/graphiql@3.0.6/graphiql.min.css" />
    <style>
      body {
        margin: 0;
        height: 100vh;
        font-family: system-ui, -apple-system, sans-serif;
      }
      #graphiql {
        height: 100vh;
      }
    </style>
  </head>
  <body>
    <div id="graphiql">Loading GraphiQL IDE...</div>
    <script src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
    <script src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
    <script src="https://unpkg.com/graphiql@3.0.6/graphiql.min.js"></script>
    <script>
      const fetcher = GraphiQL.createFetcher({ url: '/graphql' });
      ReactDOM.render(
        React.createElement(GraphiQL, {
          fetcher: fetcher,
          defaultQuery: "# Welcome to the Go GraphQL API!\n# Try running this nested query:\n\nquery {\n  books {\n    id\n    title\n    price\n    author {\n      name\n      bio\n    }\n  }\n}\n",
        }),
        document.getElementById('graphiql')
      );
    </script>
  </body>
</html>`
