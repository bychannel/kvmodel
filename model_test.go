// File model_test.go contains code for testing the model.go file.

package kvmodel

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestCompileModelSpec(t *testing.T) {
	type Primitive struct {
		Int    int
		String string
		Bool   bool
	}
	type Unexported struct {
		privateInt    int
		privateString string
		privateBool   bool
	}
	type Pointer struct {
		Int    *int
		String *string
		Bool   *bool
	}
	type Indexed struct {
		Int    int    `zoom:"index"`
		String string `zoom:"index"`
		Bool   bool   `zoom:"index"`
	}
	type Ignored struct {
		Int    int    `redis:"-"`
		String string `redis:"-"`
		Bool   bool   `redis:"-"`
	}
	type CustomName struct {
		Int    int    `redis:"myInt"`
		String string `redis:"myString"`
		Bool   bool   `redis:"myBool"`
	}
	type Inconvertible struct {
		Time time.Time
	}
	type InconvertibleIndexed struct {
		Time time.Time `zoom:"index"`
	}
	type Embedded struct {
		Primitive
	}
	type private struct {
		Int int
	}
	type EmbeddedPrivate struct {
		private
	}
	testCases := []struct {
		model         interface{}
		expectedSpec  *modelSpec
		expectedError error
	}{
		{
			model: &Primitive{},
			expectedSpec: &modelSpec{
				typ:  reflect.TypeOf(&Primitive{}),
				name: "Primitive",
				fieldsByName: map[string]*fieldSpec{
					"Int": &fieldSpec{
						kind:      primativeField,
						name:      "Int",
						redisName: "Int",
						typ:       reflect.TypeOf(Primitive{}.Int),
						indexKind: noIndex,
					},
					"String": &fieldSpec{
						kind:      primativeField,
						name:      "String",
						redisName: "String",
						typ:       reflect.TypeOf(Primitive{}.String),
						indexKind: noIndex,
					},
					"Bool": &fieldSpec{
						kind:      primativeField,
						name:      "Bool",
						redisName: "Bool",
						typ:       reflect.TypeOf(Primitive{}.Bool),
						indexKind: noIndex,
					},
				},
				fields: []*fieldSpec{
					{
						kind:      primativeField,
						name:      "Int",
						redisName: "Int",
						typ:       reflect.TypeOf(Primitive{}.Int),
						indexKind: noIndex,
					},
					{
						kind:      primativeField,
						name:      "String",
						redisName: "String",
						typ:       reflect.TypeOf(Primitive{}.String),
						indexKind: noIndex,
					},
					{
						kind:      primativeField,
						name:      "Bool",
						redisName: "Bool",
						typ:       reflect.TypeOf(Primitive{}.Bool),
						indexKind: noIndex,
					},
				},
			},
		},
		{
			model: &Unexported{},
			expectedSpec: &modelSpec{
				typ:          reflect.TypeOf(&Unexported{}),
				name:         "Unexported",
				fieldsByName: map[string]*fieldSpec{},
			},
		},
		{
			model: &Pointer{},
			expectedSpec: &modelSpec{
				typ:  reflect.TypeOf(&Pointer{}),
				name: "Pointer",
				fieldsByName: map[string]*fieldSpec{
					"Int": &fieldSpec{
						kind:      pointerField,
						name:      "Int",
						redisName: "Int",
						typ:       reflect.TypeOf(Pointer{}.Int),
						indexKind: noIndex,
					},
					"String": &fieldSpec{
						kind:      pointerField,
						name:      "String",
						redisName: "String",
						typ:       reflect.TypeOf(Pointer{}.String),
						indexKind: noIndex,
					},
					"Bool": &fieldSpec{
						kind:      pointerField,
						name:      "Bool",
						redisName: "Bool",
						typ:       reflect.TypeOf(Pointer{}.Bool),
						indexKind: noIndex,
					},
				},
				fields: []*fieldSpec{
					{
						kind:      pointerField,
						name:      "Int",
						redisName: "Int",
						typ:       reflect.TypeOf(Pointer{}.Int),
						indexKind: noIndex,
					},
					{
						kind:      pointerField,
						name:      "String",
						redisName: "String",
						typ:       reflect.TypeOf(Pointer{}.String),
						indexKind: noIndex,
					},
					{
						kind:      pointerField,
						name:      "Bool",
						redisName: "Bool",
						typ:       reflect.TypeOf(Pointer{}.Bool),
						indexKind: noIndex,
					},
				},
			},
		},
		{
			model: &Indexed{},
			expectedSpec: &modelSpec{
				typ:  reflect.TypeOf(&Indexed{}),
				name: "Indexed",
				fieldsByName: map[string]*fieldSpec{
					"Int": &fieldSpec{
						kind:      primativeField,
						name:      "Int",
						redisName: "Int",
						typ:       reflect.TypeOf(Indexed{}.Int),
						indexKind: numericIndex,
					},
					"String": &fieldSpec{
						kind:      primativeField,
						name:      "String",
						redisName: "String",
						typ:       reflect.TypeOf(Indexed{}.String),
						indexKind: stringIndex,
					},
					"Bool": &fieldSpec{
						kind:      primativeField,
						name:      "Bool",
						redisName: "Bool",
						typ:       reflect.TypeOf(Indexed{}.Bool),
						indexKind: booleanIndex,
					},
				},
				fields: []*fieldSpec{
					{
						kind:      primativeField,
						name:      "Int",
						redisName: "Int",
						typ:       reflect.TypeOf(Indexed{}.Int),
						indexKind: numericIndex,
					},
					{
						kind:      primativeField,
						name:      "String",
						redisName: "String",
						typ:       reflect.TypeOf(Indexed{}.String),
						indexKind: stringIndex,
					},
					{
						kind:      primativeField,
						name:      "Bool",
						redisName: "Bool",
						typ:       reflect.TypeOf(Indexed{}.Bool),
						indexKind: booleanIndex,
					},
				},
			},
		},
		{
			model: &Ignored{},
			expectedSpec: &modelSpec{
				typ:          reflect.TypeOf(&Ignored{}),
				name:         "Ignored",
				fieldsByName: map[string]*fieldSpec{},
			},
		},
		{
			model: &CustomName{},
			expectedSpec: &modelSpec{
				typ:  reflect.TypeOf(&CustomName{}),
				name: "CustomName",
				fieldsByName: map[string]*fieldSpec{
					"Int": &fieldSpec{
						kind:      primativeField,
						name:      "Int",
						redisName: "myInt",
						typ:       reflect.TypeOf(CustomName{}.Int),
						indexKind: noIndex,
					},
					"String": &fieldSpec{
						kind:      primativeField,
						name:      "String",
						redisName: "myString",
						typ:       reflect.TypeOf(CustomName{}.String),
						indexKind: noIndex,
					},
					"Bool": &fieldSpec{
						kind:      primativeField,
						name:      "Bool",
						redisName: "myBool",
						typ:       reflect.TypeOf(CustomName{}.Bool),
						indexKind: noIndex,
					},
				},
				fields: []*fieldSpec{
					{
						kind:      primativeField,
						name:      "Int",
						redisName: "myInt",
						typ:       reflect.TypeOf(CustomName{}.Int),
						indexKind: noIndex,
					},
					{
						kind:      primativeField,
						name:      "String",
						redisName: "myString",
						typ:       reflect.TypeOf(CustomName{}.String),
						indexKind: noIndex,
					},
					{
						kind:      primativeField,
						name:      "Bool",
						redisName: "myBool",
						typ:       reflect.TypeOf(CustomName{}.Bool),
						indexKind: noIndex,
					},
				},
			},
		},
		{
			model: &Inconvertible{},
			expectedSpec: &modelSpec{
				typ:  reflect.TypeOf(&Inconvertible{}),
				name: "Inconvertible",
				fieldsByName: map[string]*fieldSpec{
					"Time": &fieldSpec{
						kind:      inconvertibleField,
						name:      "Time",
						redisName: "Time",
						typ:       reflect.TypeOf(Inconvertible{}.Time),
						indexKind: noIndex,
					},
				},
				fields: []*fieldSpec{
					{
						kind:      inconvertibleField,
						name:      "Time",
						redisName: "Time",
						typ:       reflect.TypeOf(Inconvertible{}.Time),
						indexKind: noIndex,
					},
				},
			},
		},
		{
			model:         &InconvertibleIndexed{},
			expectedSpec:  nil,
			expectedError: errors.New("zoom: Requested index on unsupported type time.Time"),
		},
		{
			model: &Embedded{},
			expectedSpec: &modelSpec{
				typ:  reflect.TypeOf(&Embedded{}),
				name: "Embedded",
				fieldsByName: map[string]*fieldSpec{
					"Primitive": {
						kind:      inconvertibleField,
						name:      "Primitive",
						redisName: "Primitive",
						typ:       reflect.TypeOf(Primitive{}),
						indexKind: noIndex,
					},
				},
				fields: []*fieldSpec{
					{
						kind:      inconvertibleField,
						name:      "Primitive",
						redisName: "Primitive",
						typ:       reflect.TypeOf(Primitive{}),
						indexKind: noIndex,
					},
				},
			},
		},
		{
			model: &EmbeddedPrivate{},
			expectedSpec: &modelSpec{
				typ:          reflect.TypeOf(&EmbeddedPrivate{}),
				name:         "EmbeddedPrivate",
				fieldsByName: map[string]*fieldSpec{},
			},
		},
	}
	for _, tc := range testCases {
		gotSpec, err := compileModelSpec(reflect.TypeOf(tc.model))
		if tc.expectedError == nil {
			if err != nil {
				t.Error("Error compiling model spec: ", err.Error())
				continue
			}
			if !reflect.DeepEqual(tc.expectedSpec, gotSpec) {
				t.Errorf(
					"Incorrect model spec.\nExpected: %s\nBut got:  %s\n",
					spew.Sprint(tc.expectedSpec),
					spew.Sprint(gotSpec),
				)
			}
		} else {
			if err == nil {
				t.Errorf(
					"Didn't get an error but expected: %s\n",
					spew.Sprint(tc.expectedError),
				)
			}
			if !reflect.DeepEqual(tc.expectedError, err) {
				t.Errorf(
					"Incorrect error.\nExpected: %s\nBut got:  %s\n",
					spew.Sprint(tc.expectedError),
					spew.Sprint(err),
				)
			}
		}
	}
}
