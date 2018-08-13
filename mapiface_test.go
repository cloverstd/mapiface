package mapiface_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/cloverstd/mapiface"
	"github.com/stretchr/testify/assert"
)

func TestTest(t *testing.T) {
	var mapiface = map[string]interface{}{}
	var ifaceType = reflect.TypeOf(mapiface).Elem()
	log.Println(ifaceType.Kind())
}

func TestConvert(t *testing.T) {
	s := "a"
	m := map[int]*string{
		1: &s,
	}
	v, err := mapiface.Convert(m)
	log.Println(err)
	log.Println(reflect.TypeOf(v))
	log.Println(reflect.TypeOf(v.(map[int]interface{})[1]))
}

func TestSlice(t *testing.T) {
	assert := assert.New(t)
	v := []string{"1"}
	s, err := mapiface.Convert(v)
	assert.Nil(err)
	assert.IsType([]interface{}{}, s)
	assert.Equal(v[0], s.([]interface{})[0])
}

func TestSlicePtr(t *testing.T) {
	assert := assert.New(t)
	s := "str"
	var a interface{} = []*string{&s}
	b, err := mapiface.Convert(a)
	assert.Nil(err)
	assert.IsType([]interface{}{}, b)
	assert.Equal(s, *b.([]interface{})[0].(*string))
}

func TestError(t *testing.T) {
	assert := assert.New(t)
	var a interface{}
	b, err := mapiface.Convert(a)
	assert.NotNil(err)
	assert.Nil(b)

	a = map[string]interface{}{
		"nil": nil,
	}
	b, err = mapiface.Convert(a)
	assert.Nil(err)
	log.Println(b)
}

func TestStruct(t *testing.T) {
	assert := assert.New(t)
	v := struct {
		UserID   int
		Name     string      `json:"name"`
		Value    interface{} `json:"value"`
		Point    float64     `json:"point"`
		Empty    interface{} `json:"empty,omitempty"`
		NotEmpty interface{} `json:"not_empty"`
		private  int
	}{
		UserID:  1,
		Name:    "name",
		Value:   map[string]interface{}{"a": 1},
		Point:   10,
		private: 1,
	}

	b, err := mapiface.Convert(v)
	assert.Nil(err)
	assert.IsType(map[string]interface{}{}, b)
	vv := b.(map[string]interface{})
	assert.Equal(v.UserID, vv["UserID"])
	_, exist := vv["user_id"]
	assert.False(exist)
	assert.Equal(v.Name, vv["name"])
	assert.Equal(v.Value, vv["value"])
	assert.Equal(v.Point, vv["point"])
	_, exist = vv["empty"]
	assert.False(exist)
	_, exist = vv["private"]
	assert.False(exist)
	_, exist = vv["not_empty"]
	assert.True(exist)
	assert.Nil(vv["not_empty"])
}
