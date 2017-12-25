package lib

import (
	"fmt"
	"reflect"
)

// Dict :
type Dict map[string]interface{}

func (d Dict) String() string {
	s := ""
	for k, v := range d {
		s = s + fmt.Sprintf("\"%s\": %v", k, v)
	}
	return "{" + s + "}"
}

// Clear :None.  Remove all items from D
func (d *Dict) Clear() {
	for k := range *d {
		delete(*d, k)
	}
}

// Copy :D.copy() -> a shallow copy of D
func (d Dict) Copy() Dict {
	return Dict{}
}

// FromKeys :Returns a new dict with keys from iterable and values equal to value.
func (d Dict) FromKeys() Dict {
	return Dict{}
}

// Get :D.get(k) -> D[k] if k in D, else nil
func (d Dict) Get(k string) interface{} {
	return nil
}

// Items :D.items() -> a set-like object providing a view on D's items
func (d Dict) Items() {

}

// Keys :D.keys() -> a set-like object providing a view on D's keys
func (d Dict) Keys() {

}

// Pop :D.pop(k[,d]) -> v, remove specified key and return the corresponding value.
// If key is not found, d is returned if given, otherwise KeyError is raised
func (d *Dict) Pop() {

}

// PopItem :D.popitem() -> (k, v), remove and return some (key, value) pair as a
// 2-tuple; but raise KeyError if D is empty.
func (d *Dict) PopItem() {

}

// SetDefault :D.setdefault(k[,d]) -> D.get(k,d), also set D[k]=d if k not in D
func (d *Dict) SetDefault() {

}

// Update :D.update([E, ]**F) -> None.  Update D from dict/iterable E and F.
// If E is present and has a .keys() method, then does:  for k in E: D[k] = E[k]
// If E is present and lacks a .keys() method, then does:  for k, v in E: D[k] = v
// In either case, this is followed by: for k in F:  D[k] = F[k]
func (d *Dict) Update() {

}

// Values :D.values() -> an object providing a view on D's values
func (d Dict) Values() {

}

// Struct :
type Struct struct{}

// JSON :
func (s *Struct) JSON() string {
	return Stoj(s)
}

// Map :
func (s *Struct) Map() map[string]interface{} {
	elem := reflect.ValueOf(s).Elem()
	etype := elem.Type()
	rmap := map[string]interface{}{}
	for i := 0; i < etype.NumField(); i++ {
		rmap[etype.Field(i).Name] = elem.Field(i).Interface()
	}
	return rmap
}
