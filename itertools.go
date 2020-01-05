package itertools

import (
	"errors"
	"reflect"
)

// ErrIterStop ...
var ErrIterStop = errors.New("There are no more items in iterator")

type any reflect.Value

// Iter ...
type Iter struct {
	Next iterFunc
}

type iterFunc func() (reflect.Value, error)

// type filterFunc func(i interface{}) bool

// Filter values
func (iter Iter) Filter(filterFunc interface{}) Iter {
	filter := reflect.ValueOf(filterFunc).Call
	return Iter{
		Next: func() (item reflect.Value, err error) {
			for item, err = iter.Next(); err == nil; item, err = iter.Next() {
				args := []reflect.Value{item}
				if filter(args)[0].Bool() {
					break
				}
			}
			return item, err
		},
	}
}

// Each execute given func with iterator's value as argument
func (iter Iter) Each(eachFunc interface{}) {
	each := reflect.ValueOf(eachFunc).Call
	for item, err := iter.Next(); err == nil; item, err = iter.Next() {
		args := []reflect.Value{item}
		each(args)
	}
}

type istrs []string

// NewIter ...
func NewIter(slice interface{}) (iter Iter) {
	sliceVal := reflect.ValueOf(slice)
	mssindex := -1
	msslen := sliceVal.Len()
	return Iter{
		Next: func() (item reflect.Value, err error) {
			mssindex++
			if mssindex >= msslen {
				return item, ErrIterStop
			}
			return sliceVal.Index(mssindex), nil
		},
	}
}
