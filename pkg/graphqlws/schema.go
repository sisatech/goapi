package graphqlws

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/graphql-go/graphql"
)

// Update contains information and methods used to help the server decide which
// subscriptions should be updated. A good implementation of Update should
// front-load computation because it should be considered reusable, and it is
// important to minimize the computational complexity of the Overlaps function.
type Update interface {

	// When martialled into JSON, an Update should contain enough
	// information for a Schema to completely reconstruct it and also
	// validate that it is meant for the same Schema. This is necessary to
	// prevent bugs caused by mismatched GraphQL schemas that could occur
	// when Updates are published to a cluster.
	json.Marshaler

	// Overlaps returns true if the two Update objects share at least one
	// field (or leaf) of the GraphQL schema.
	Overlaps(update Update) bool
}

// Schema contains information and methods used to analyze a GraphQL schema and
// generate Update objects.
type Schema interface {

	// Fields returns a list of alphabetically sorted strings representing
	// each leaf of a GraphQL schema's subscriptions.
	Fields() []string

	// NewUpdate returns an Update object created from a list of fields. All
	// fields provided should exists in the list of strings returned by the
	// Fields function on this interface.
	NewUpdate(fields ...string) Update

	// UpdateFromObject is a convenience function that will check if the
	// provided object pointer can be found anywhere in the underlying
	// schema, and use the results of this search to create a new Update
	// from all of the identified fields. The function must be smart enough
	// to find the object if it appears multiple times in the same schema.
	UpdateFromObject(o *graphql.Object) Update

	// UpdateFromJSON must be able to recreate an Update object from a JSON
	// payload (as created by the Update object's MarshalJSON function). It
	// must also be able to validate that the payload perfectly matches the
	// schema. This is necessary to prevent bugs caused by mismatched
	// GraphQL schemas that could occur when Updates are published to a
	// cluster.
	UpdateFromJSON(data []byte) (Update, error)
}

type update struct {
	bitmap []byte
	schema Schema
}

type updateJSON struct {
	Schema []*updateTupleJSON `json:"schema"`
}

type updateTupleJSON struct {
	Field    string `json:"field"`
	Relevant bool   `json:"relevant"`
}

func (u *update) MarshalJSON() ([]byte, error) {

	var fields []*updateTupleJSON
	for i, field := range u.schema.Fields() {
		var relevant bool
		bit := u.bitmap[i/8]
		relevant = (bit & (1 << uint(i%8))) > 0
		fields = append(fields, &updateTupleJSON{
			Field:    field,
			Relevant: relevant,
		})
	}

	return json.Marshal(&updateJSON{
		Schema: fields,
	})

}

func (u *update) Overlaps(x Update) bool {

	y, ok := x.(*update)
	if !ok {
		panic(errors.New("cannot mix and match two different implementations of the Update interface"))
	}

	for i := 0; i < len(u.bitmap); i++ {
		if u.bitmap[i]&y.bitmap[i] > 0 {
			return true
		}
	}

	return false

}

type schema struct {
	fields []string
	schema graphql.Schema
}

func (s *schema) Fields() []string {
	return s.fields
}

func (s *schema) UpdateFromJSON(data []byte) (Update, error) {

	x := new(updateJSON)
	err := json.Unmarshal(data, x)
	if err != nil {
		return nil, err
	}

	if len(s.Fields()) != len(x.Schema) {
		return nil, errors.New("update doesn't match schema")
	}

	l := (len(x.Schema) + 7) / 8

	u := new(update)
	u.schema = s
	u.bitmap = make([]byte, l)

	for i, tuple := range x.Schema {
		if s.fields[i] != tuple.Field {
			return nil, errors.New("update doesn't match schema")
		}
		if tuple.Relevant {
			u.bitmap[i/8] |= 1 << uint(i&8)
		}
	}

	return u, nil

}

// NewSchema analyzes the provided 'graphql.Schema' and returns an
// implementation of the Schema interface.
func NewSchema(gql graphql.Schema) Schema {
	s := new(schema)
	s.schema = gql

	subs := gql.SubscriptionType()
	if subs == nil {
		return s
	}

	var fields []string
	var recurseForFields func(parent string, o *graphql.Object)
	recurseForFields = func(parent string, o *graphql.Object) {
		for k, v := range o.Fields() {
			key := parent + k
			if child, ok := v.Type.(*graphql.Object); ok {
				key += "."
				recurseForFields(key, child)
				continue
			} else if list, ok := v.Type.(*graphql.List); ok {
				if child, ok := list.OfType.(*graphql.Object); ok {
					key += "."
					recurseForFields(key, child)
					continue
				}
			}
			fields = append(fields, key)
		}
	}

	recurseForFields("", subs)
	sort.Strings(fields)
	s.fields = fields

	return s
}

func (s *schema) NewUpdate(fields ...string) Update {
	u := new(update)
	u.schema = s
	l := (len(s.fields) + 7) / 8
	u.bitmap = make([]byte, l)
	for _, f := range fields {
		index := sort.SearchStrings(s.fields, f)
		if index >= len(s.fields) || s.fields[index] != f {
			panic(fmt.Errorf("failed to create update: field '%s' does not exist in the schema", f))
		}
		u.bitmap[index/8] |= 1 << uint(index%8)
	}
	return u
}

func (s *schema) UpdateFromObject(o *graphql.Object) Update {

	schema := s.schema
	x := schema.SubscriptionType()

	var fields []string
	var activate bool

	var flatten func(prefix string, obj *graphql.Object)
	flatten = func(prefix string, obj *graphql.Object) {

		alreadyActive := activate
		if obj == o {
			activate = true
		}

		for k, f := range obj.Fields() {
			subprefix := prefix + k
			if y, ok := f.Type.(*graphql.Object); ok {
				flatten(subprefix+".", y)
				continue
			} else if z, ok := f.Type.(*graphql.List); ok {
				if _, ok = z.OfType.(*graphql.Object); ok {
					flatten(subprefix+".", y)
					continue
				}
			}
			if activate {
				fields = append(fields, subprefix)
			}
		}

		activate = alreadyActive

	}

	flatten("", x)
	return s.NewUpdate(fields...)
}
