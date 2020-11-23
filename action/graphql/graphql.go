package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Schema : GraphQL schema definition
var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryTypes,
		Mutation: mutationTypes,
	},
)

<<<<<<< HEAD
// GraphqlHandler : GraphQL schema handler
//
// Handler config options
=======
// Handler : GraphQL schema handler
//
// Config options of handler
>>>>>>> f41ff24f626bd8c0587cb05747b5a3edd16976db
//
// Schema : GraphQL schema definition variable name
//
// Pretty : Show sorted json code in GraphiQL
//
// GraphiQL : Show GraphQL GUI request form in web browser
<<<<<<< HEAD
var GraphqlHandler = handler.New(&handler.Config{
=======
var Handler = handler.New(&handler.Config{
>>>>>>> f41ff24f626bd8c0587cb05747b5a3edd16976db
	Schema:   &Schema,
	Pretty:   true,
	GraphiQL: true,
})
