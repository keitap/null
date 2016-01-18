package null

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/mholt/binding"
)

type TestStruct struct {
	B Bool
	I Int
	F Float
	S String
}

func (this *TestStruct) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&this.B: "bool",
		&this.I: "int",
		&this.F: "float",
		&this.S: "string",
	}
}

func TestBind(t *testing.T) {
	req := &http.Request{}
	req.PostForm = url.Values{}
	req.PostForm.Set("bool", "true")
	req.PostForm.Set("int", "12345")
	req.PostForm.Set("float", "1.2345")
	req.PostForm.Set("string", "test")

	ts := &TestStruct{}
	errs := binding.Form(req, ts)

	if 0 < errs.Len() {
		t.Error(errs.Error())
	}

	assertBool(t, ts.B, "binding")
	assertInt(t, ts.I, "binding")
	assertFloat(t, ts.F, "binding")
	assertStr(t, ts.S, "binding")
}

func TestBindNull(t *testing.T) {
	req := &http.Request{}
	req.PostForm = url.Values{}

	ts := &TestStruct{}
	errs := binding.Form(req, ts)

	if 0 < errs.Len() {
		t.Error(errs.Error())
	}

	assertNullBool(t, ts.B, "binding")
	assertNullInt(t, ts.I, "binding")
	assertNullFloat(t, ts.F, "binding")
	assertNullStr(t, ts.S, "binding")
}
