package examples

import (
	"testing"

	"github.com/knightso/base/errors"
	"google.golang.org/appengine/aetest"

	"google.golang.org/appengine/datastore"
)

//このファイルは出力されたファイルに対するテストです

func TestGenFile(t *testing.T) {

	//遅延をオフ
	opt := &aetest.Options{StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	r, err := inst.NewRequest("GET", "/gophers", nil)
	if err != nil {
		t.Fatalf("Failed to create req1: %v", err)
	}

	d := WorkType{}

	err = d.GenerateKey(r)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	err = d.Put(r)
	if err != nil {
		t.Fatalf("Failed to put: %v", err)
	}

	keyString := d.GetKey().StringID()
	cKey := d.CreateKey(r, keyString)

	if keyString != cKey.StringID() {
		t.Fatalf("Failed to Key %s != %s", keyString, cKey.StringID())
	}

	newD := WorkType{}
	err = newD.SelectById(r, keyString)
	if err != nil {
		t.Fatalf("Failed to SelectById: %v", err)
	}

	selectKey := newD.GetKey().StringID()
	if keyString != selectKey {
		t.Fatalf("Failed to Key %s != %s", keyString, selectKey)
	}

	slice, err := newD.Select(r, 1)
	if err != nil {
		t.Fatalf("Failed to Select: %v", err)
	}

	if slice == nil {
		t.Fatalf("Failed to Select: return nil")
	}

	if len(slice) != 1 {
		t.Fatalf("Failed to Select: length not 1[%d]", len(slice))
	}

	err = newD.DeleteLogical(r)
	if err != nil {
		t.Fatalf("Failed to Delete Logical: %v", err)
	}

	slice, err = newD.Select(r, 1)
	if err != nil {
		t.Fatalf("Failed to Select(logical delete): %v", err)
	}

	if slice == nil {
		t.Fatalf("Failed to Select(logical delete): return nil")
	}

	if len(slice) != 0 {
		t.Fatalf("Failed to Select(logical delete): length not 0[%d]", len(slice))
	}

	deleteD := WorkType{}
	err = deleteD.SelectById(r, keyString)
	if err != nil {
		t.Fatalf("Failed to SelectById(delete logical): %v", err)
	}

	deleteKey := deleteD.GetKey().StringID()
	if deleteKey != selectKey {
		t.Fatalf("Failed to Key %s != %s", deleteKey, selectKey)
	}

	err = deleteD.Delete(r)
	if err != nil {
		t.Fatalf("Failed to Delete: %v", err)
	}

	dd := WorkType{}
	err = dd.SelectById(r, keyString)
	if err == nil {
		t.Fatalf("Failed to SelectById(delete): %v", err)
	}

	if errors.Root(err) != datastore.ErrNoSuchEntity {
		t.Fatalf("Failed to SelectById(delete):Not NoSutiEntity [%v]", err)
	}

}
