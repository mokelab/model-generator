package golang

import (
	"testing"
)

func TestCommon_0000_toCamelcase(t *testing.T) {
	g := &generator{}
	out := g.toCamelcase("user_name")
	if out != "userName" {
		t.Errorf("Unexpected output : %s", out)
	}
}

func TestCommon_0001_toCamelcase_no(t *testing.T) {
	g := &generator{}
	out := g.toCamelcase("user")
	if out != "user" {
		t.Errorf("Unexpected output : %s", out)
	}
}

func TestCommon_0010_toGolangType(t *testing.T) {
	g := &generator{}
	out := g.toGolangType("varchar")
	if out != "string" {
		t.Errorf("Unexpected output : %s", out)
	}
}

func TestCommon_0011_toGolangType_text(t *testing.T) {
	g := &generator{}
	out := g.toGolangType("text")
	if out != "string" {
		t.Errorf("Unexpected output : %s", out)
	}
}

func TestCommon_0012_toGolangType_int(t *testing.T) {
	g := &generator{}
	out := g.toGolangType("int")
	if out != "int" {
		t.Errorf("Unexpected output : %s", out)
	}
}

func TestCommon_0013_toGolangType_long(t *testing.T) {
	g := &generator{}
	out := g.toGolangType("long")
	if out != "int64" {
		t.Errorf("Unexpected output : %s", out)
	}
}

func TestCommon_0014_toGolangType_float(t *testing.T) {
	g := &generator{}
	out := g.toGolangType("float")
	if out != "float32" {
		t.Errorf("Unexpected output : %s", out)
	}
}

func TestCommon_0015_toGolangType_double(t *testing.T) {
	g := &generator{}
	out := g.toGolangType("double")
	if out != "float64" {
		t.Errorf("Unexpected output : %s", out)
	}
}
