package types

type CoolThing struct {
	Name        string       `json:"name" bson:"name"`
	ThirdName   string       `json:"thirdName,omitempty" bson:"thirdName"`
	UsefulArray []bool       `json:"usefulArray" bson:"usefulArray"`
	ChildThing  ChildThing   `json:"childThing" bson:"childThing"`
	ManyThings  []ChildThing `json:"manyThings" bson:"manyThings"`
}
