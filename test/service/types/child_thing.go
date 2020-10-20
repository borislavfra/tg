package types

type ChildThing struct {
	Anything interface{} `json:"anything" bson:"anything"`
}
