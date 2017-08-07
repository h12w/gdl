package testpkg

type (
	TestStruct struct {
		IntField    int
		StringField string
		StructField TestStructA
	}
	TestStructA struct {
		IntField int
	}
)
