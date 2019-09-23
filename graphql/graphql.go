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

// Handler : GraphQL schema handler
//
// Config options of handler
//
// Schema : GraphQL schema definition variable name
//
// Pretty : Show sorted json code in GraphiQL
//
// GraphiQL : Show GraphQL GUI request form in web browser
var Handler = handler.New(&handler.Config{
	Schema:   &Schema,
	Pretty:   true,
	GraphiQL: true,
})
