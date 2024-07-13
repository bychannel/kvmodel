// File convert_test.go tests the conversion
// to and from go data structures of a variety of types.

package kvmodel

import (
	"reflect"
	"testing"
	"time"
)

func TestConvertPrimatives(t *testing.T) {
	testingSetUp()
	defer testingTearDown()
	testConvertType(t, indexedPrimativesModels, createIndexedPrimativesModel())
}

func TestConvertPointers(t *testing.T) {
	testingSetUp()
	defer testingTearDown()
	testConvertType(t, indexedPointersModels, createIndexedPointersModel())
}

func TestTimeDuration(t *testing.T) {
	testingSetUp()
	defer testingTearDown()

	type durationModel struct {
		Duration time.Duration
		RandomID
	}
	durationModels, err := testPool.NewCollection(&durationModel{})
	if err != nil {
		t.Errorf("Unexpected error in testPool.NewCollection: %s", err.Error())
	}
	model := &durationModel{
		Duration: 43 * time.Second,
	}
	testConvertType(t, durationModels, model)
}

func TestGobFallback(t *testing.T) {
	testingSetUp()
	defer testingTearDown()

	type gobModel struct {
		Complex     complex128
		IntSlice    []int
		StringSlice []string
		IntArray    [3]int
		StringArray [3]string
		StringMap   map[string]string
		IntMap      map[int]int
		RandomID
	}
	options := DefaultCollectionOptions.WithFallbackMarshalerUnmarshaler(GobMarshalerUnmarshaler)
	gobModels, err := testPool.NewCollectionWithOptions(&gobModel{}, options)
	if err != nil {
		t.Errorf("Unexpected error in testPool.NewCollection: %s", err.Error())
	}
	model := &gobModel{
		Complex:     randomComplex(),
		IntSlice:    []int{randomInt(), randomInt(), randomInt()},
		StringSlice: []string{randomString(), randomString(), randomString()},
		IntArray:    [3]int{randomInt(), randomInt(), randomInt()},
		StringArray: [3]string{randomString(), randomString(), randomString()},
		StringMap:   map[string]string{randomString(): randomString(), randomString(): randomString()},
		IntMap:      map[int]int{randomInt(): randomInt(), randomInt(): randomInt()},
	}
	testConvertType(t, gobModels, model)
}

func TestJSONFallback(t *testing.T) {
	testingSetUp()
	defer testingTearDown()

	type jsonModel struct {
		IntSlice       []int
		StringSlice    []string
		IntArray       [3]int
		StringArray    [3]string
		StringMap      map[string]string
		EmptyInterface interface{}
		RandomID
	}
	options := DefaultCollectionOptions.WithFallbackMarshalerUnmarshaler(JSONMarshalerUnmarshaler)
	jsonModels, err := testPool.NewCollectionWithOptions(&jsonModel{}, options)
	if err != nil {
		t.Errorf("Unexpected error in testPool.NewCollection: %s", err.Error())
	}
	model := &jsonModel{
		IntSlice:       []int{randomInt(), randomInt(), randomInt()},
		StringSlice:    []string{randomString(), randomString(), randomString()},
		IntArray:       [3]int{randomInt(), randomInt(), randomInt()},
		StringArray:    [3]string{randomString(), randomString(), randomString()},
		StringMap:      map[string]string{randomString(): randomString(), randomString(): randomString()},
		EmptyInterface: map[string]interface{}{"key": []interface{}{"This satisfies empty interface"}},
	}
	testConvertType(t, jsonModels, model)
}

type Embeddable struct {
	Int    int
	String string
	Bool   bool
}

func TestConvertEmbeddedStruct(t *testing.T) {
	testingSetUp()
	defer testingTearDown()

	type embeddedStructModel struct {
		Embeddable
		RandomID
	}
	embededStructModels, err := testPool.NewCollection(&embeddedStructModel{})
	if err != nil {
		t.Errorf("Unexpected error in testPool.NewCollection: %s", err.Error())
	}
	model := &embeddedStructModel{
		Embeddable: Embeddable{
			Int:    randomInt(),
			String: randomString(),
			Bool:   randomBool(),
		},
	}
	testConvertType(t, embededStructModels, model)
}

func TestEmbeddedPointerToStruct(t *testing.T) {
	testingSetUp()
	defer testingTearDown()

	type embeddedPointerToStructModel struct {
		*Embeddable
		RandomID
	}
	embededPointerToStructModels, err := testPool.NewCollection(&embeddedPointerToStructModel{})
	if err != nil {
		t.Errorf("Unexpected error in testPool.NewCollection: %s", err.Error())
	}
	model := &embeddedPointerToStructModel{
		Embeddable: &Embeddable{
			Int:    randomInt(),
			String: randomString(),
			Bool:   randomBool(),
		},
	}
	testConvertType(t, embededPointerToStructModels, model)
}

// testConvertType is a general test that uses reflection. It saves model to the databse then finds it. If
// the found copy does not exactly match the original, it reports an error via t.Error or t.Errorf
func testConvertType(t *testing.T, collection *Collection, model Model) {
	// Make sure we can save the model without errors
	if err := collection.Save(model); err != nil {
		t.Errorf("Unexpected error in Save: %s", err.Error())
	}
	// Find the model from the database and scan it into a new copy
	modelCopy, ok := reflect.New(collection.spec.typ.Elem()).Interface().(Model)
	if !ok {
		t.Fatalf("Unexpected error: Could not convert type %s to Model", collection.spec.typ.String())
	}
	if err := collection.Find(model.ModelID(), modelCopy); err != nil {
		t.Errorf("Unexpected error in Find: %s", err.Error())
	}
	// Make sure the copy equals the original
	if !reflect.DeepEqual(model, modelCopy) {
		t.Errorf("Model of type %T was not saved/retrieved correctly.\nExpected: %+v\nGot:      %+v", model, model, modelCopy)
	}
	// Make sure we can save a model with all nil fields. This should
	// not cause an error.
	emptyModel, ok := reflect.New(collection.spec.typ.Elem()).Interface().(Model)
	if !ok {
		t.Fatalf("Unexpected error: Could not convert type %s to Model", collection.spec.typ.String())
	}
	if err := collection.Save(emptyModel); err != nil {
		t.Errorf("Unexpected error saving an empty model: %s", err.Error())
	}
	emptyModelCopy, _ := reflect.New(collection.spec.typ.Elem()).Interface().(Model)
	if err := collection.Find(emptyModel.ModelID(), emptyModelCopy); err != nil {
		t.Errorf("Unexpected error in Find: %s", err.Error())
	}
	// Make sure the copy equals the original
	if !reflect.DeepEqual(emptyModel, emptyModelCopy) {
		t.Errorf("Model of type %T was not saved/retrieved correctly.\nExpected: %+v\nGot:      %+v", emptyModel, emptyModel, emptyModelCopy)
	}
}
