package itertools

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrIterStop ...
var ErrIterStop = errors.New("There are no more items in iterator")

// NewErrIter create an interator that produce err from given string
func NewErrIter(s string) Iter {
	return Iter{
		Next: func() (i interface{}, err error) {
			return i, fmt.Errorf(s)
		},
	}
}

// Iter ...
type Iter struct {
	Next iterFunc
}

type iterFunc func() (interface{}, error)

// type filterFunc func(i interface{}) bool

// Filter values
func (iter Iter) Filter(filterFunc interface{}) Iter {
	filter := reflect.ValueOf(filterFunc).Call
	return Iter{
		Next: func() (item interface{}, err error) {
			for item, err = iter.Next(); err == nil; item, err = iter.Next() {
				args := []reflect.Value{reflect.ValueOf(item)}
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
	for {
		item, err := iter.Next()
		switch err {
		case nil:
			args := []reflect.Value{reflect.ValueOf(item)}
			each(args)
		case ErrIterStop:
			return
		default:
			panic(err)
		}
	}
}

// SliceIter ...
func SliceIter(slice interface{}) (iter Iter) {
	sliceVal := reflect.ValueOf(slice)
	mssindex := 0
	msslen := sliceVal.Len()
	return Iter{
		Next: func() (item interface{}, err error) {
			if mssindex >= msslen {
				return item, ErrIterStop
			}
			item = sliceVal.Index(mssindex).Interface()
			mssindex++
			return item, nil
		},
	}
}
