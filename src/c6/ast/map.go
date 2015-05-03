package ast

type Map struct {
}

func (self *Map) Set() {

}

func (self *Map) Get() {
}

func (self Map) GetValueType() ValueType {
	return MapValue
}

func NewMap() *Map {
	return &Map{}
}
