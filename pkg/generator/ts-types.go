package generator

type tsSchemas map[string]tsSchema

type tsProperties map[string]tsSchema

type tsSchema struct {
	Ref     string
	Type    string
	Format  string
	Minimum int
	Maximum int

	Items                *tsSchema
	Properties           tsProperties
	AdditionalProperties interface{}

	Nullable bool

	Description string
	Example     interface{}
}
