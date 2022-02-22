package data

import (
	"testing"
)

type StructTest struct {
	Name string `col:"nama"`
	Age  int
	City string `col:"CITY"`
}

var JsonTest = `
	{
		"nama": "Ziska Zarkasyi",
		"Age": 42,
		"CITY": "Depok"
	}
`

var ObjectTest = StructTest{
	Name: "Ziska Zarkasyi",
	Age:  42,
	City: "Depok",
}

func TestRawBytes(t *testing.T) {
	d, err := FromStruct(ObjectTest)
	if err != nil {
		t.Fail()
	}

	djs, err := FromJson([]byte(JsonTest))
	if err != nil {
		t.Fail()
	}

	t.Logf("from Object: %#v\n", d)
	t.Logf("from Json: %#v\n", djs)
}
