package goapi

import (
	"fmt"

	"github.com/machinebox/graphql"
)

// Cursor ..
type Cursor struct {
	First  int
	Last   int
	After  string
	Before string
}

// Strings ..
func (c *Cursor) Strings() (string, string) {
	var variableDeclarations string
	var variables string
	var cursorPresent bool

	if c != nil {
		if c.After != "" {
			cursorPresent = true
			if len(variableDeclarations) != 0 {
				variableDeclarations = fmt.Sprintf("%s,", variableDeclarations)
			}
			if len(variables) != 0 {
				variables = fmt.Sprintf("%s,", variables)
			}
			variables = fmt.Sprintf("%safter:$after", variables)
			variableDeclarations = fmt.Sprintf("%s $after: String", variableDeclarations)
		}
		if c.Before != "" {
			cursorPresent = true
			if len(variableDeclarations) != 0 {
				variableDeclarations = fmt.Sprintf("%s,", variableDeclarations)
			}
			if len(variables) != 0 {
				variables = fmt.Sprintf("%s,", variables)
			}
			variables = fmt.Sprintf("%sbefore:$before", variables)
			variableDeclarations = fmt.Sprintf("%s $before: String", variableDeclarations)
		}
		if c.First != 0 {
			cursorPresent = true
			if len(variableDeclarations) != 0 {
				variableDeclarations = fmt.Sprintf("%s,", variableDeclarations)
			}
			if len(variables) != 0 {
				variables = fmt.Sprintf("%s,", variables)
			}
			variables = fmt.Sprintf("%sfirst:$first", variables)
			variableDeclarations = fmt.Sprintf("%s$first: Int", variableDeclarations)
		}
		if c.Last != 0 {
			cursorPresent = true
			if len(variableDeclarations) != 0 {
				variableDeclarations = fmt.Sprintf("%s,", variableDeclarations)
			}
			if len(variables) != 0 {
				variables = fmt.Sprintf("%s,", variables)
			}
			variables = fmt.Sprintf("%slast:$last", variables)
			variableDeclarations = fmt.Sprintf("%s $last: Int", variableDeclarations)
		}

		if cursorPresent {
			variables = fmt.Sprintf("(%s)", variables)
			variableDeclarations = fmt.Sprintf("(%s)", variableDeclarations)
		}
	}

	return variableDeclarations, variables
}

// AddToRequest ..
func (c *Cursor) AddToRequest(req *graphql.Request) {
	if c.After != "" {
		req.Var("after", c.After)
	}
	if c.Before != "" {
		req.Var("before", c.Before)
	}
	if c.First != 0 {
		req.Var("first", c.First)
	}
	if c.Last != 0 {
		req.Var("last", c.Last)
	}
}
