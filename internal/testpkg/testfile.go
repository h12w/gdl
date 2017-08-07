package testpkg

type (
	TestStruct struct {
		IntField    int         `id:"1"`
		StringField string      `id:"2"`
		StructField TestStructA `id:"3"`
	}
	TestStructA struct {
		IntField int
	}
)
