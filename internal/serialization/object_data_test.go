package serialization

import (
	"testing"
	"bytes"
	"reflect"
	"github.com/hazelcast/go-client/config"
)

func TestObjectDataOutput_EnsureAvailable(t *testing.T) {
	o := NewObjectDataOutput(2, nil, false)
	o.EnsureAvailable(5)
	buf := o.buffer
	expectedBuf := []byte{0, 0, 0, 0, 0}
	if bytes.Compare(buf, expectedBuf) != 0 {
		t.Errorf("EnsureAvailable() makes ", buf, " expected ", expectedBuf)
	}

}

func TestObjectDataOutput_WriteInt32(t *testing.T) {
	o := NewObjectDataOutput(4, nil, false)
	o.WriteInt32(1)
	o.WriteInt32(2)
	o.WriteInt32(3)

	if o.buffer[0] != 1 || o.buffer[4] != 2 || o.buffer[8] != 3 {
		t.Errorf("WriteInt32() writes to wrong position!")
	}
}

func TestObjectDataInput_AssertAvailable(t *testing.T) {
	o := NewObjectDataInput([]byte{0, 1, 2, 3}, 3, &SerializationService{}, true)
	ret := o.AssertAvailable(2)
	if ret == nil {
		t.Errorf("AssertAvailable() should return error '%s' but it returns nil!", ret)
	}
}

func TestObjectDataInput_ReadInt32(t *testing.T) {
	o := NewObjectDataInput([]byte{0, 0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0}, 4, &SerializationService{}, false)
	expectedRet := 4
	ret, _ := o.ReadInt32()

	if ret != int32(expectedRet) {
		t.Errorf("ReadInt32() returns '%s' expected %s", ret, expectedRet)
	}
}

func TestObjectDataInput_ReadInt32WithPosition(t *testing.T) {
	o := NewObjectDataInput([]byte{0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0}, 4, &SerializationService{}, false)
	expectedRet := 2
	ret, _ := o.ReadInt32WithPosition(8)

	if ret != int32(expectedRet) {
		t.Errorf("ReadInt32WithPosition() returns '%s' expected %s", ret, expectedRet)
	}
}

func TestObjectDataInput_ReadFloat64(t *testing.T) {
	o := NewObjectDataOutput(24, nil, false)
	o.WriteFloat64(1.234)
	o.WriteFloat64(2.544)
	o.WriteFloat64(3.432)
	i := NewObjectDataInput(o.buffer, 16, nil, false)
	var expectedRet float64 = 3.432
	var ret float64
	ret, _ = i.ReadFloat64()
	if ret != expectedRet {
		t.Errorf("ReadFloat64() returns '%s' expected %s", ret, expectedRet)
	}
}

func TestObjectDataInput_ReadFloat64WithPosition(t *testing.T) {
	o := NewObjectDataOutput(24, &SerializationService{}, false)
	o.WriteFloat64(1.234)
	o.WriteFloat64(2.544)
	o.WriteFloat64(3.432)
	i := NewObjectDataInput(o.buffer, 16, &SerializationService{}, false)
	var expectedRet float64 = 2.544
	var ret float64
	ret, _ = i.ReadFloat64WithPosition(8)
	if ret != expectedRet {
		t.Errorf("ReadFloat64WithPosition() returns '%s' expected %s", ret, expectedRet)
	}
}

func TestObjectDataInput_ReadBool(t *testing.T) {
	o := NewObjectDataOutput(9, &SerializationService{}, false)
	o.WriteFloat64(1.234)
	o.WriteBool(true)
	i := NewObjectDataInput(o.buffer, 8, &SerializationService{}, false)
	var expectedRet bool = true
	var ret bool
	ret, _ = i.ReadBool()
	if ret != expectedRet {
		t.Errorf("ReadBool() returns '%s' expected %s", ret, expectedRet)
	}
}

func TestObjectDataInput_ReadBoolWithPosition(t *testing.T) {
	o := NewObjectDataOutput(9, &SerializationService{}, false)
	o.WriteFloat64(1.234)
	o.WriteBool(true)
	i := NewObjectDataInput(o.buffer, 7, &SerializationService{}, false)
	var expectedRet bool = true
	var ret bool
	ret, _ = i.ReadBoolWithPosition(8)
	if ret != expectedRet {
		t.Errorf("ReadBoolWithPosition() returns '%s' expected %s", ret, expectedRet)
	}
}

func TestObjectDataInput_ReadObject(t *testing.T) {
	conf:=config.NewSerializationConfig()
	service:=NewSerializationService(conf)
	o := NewObjectDataOutput(500,service , false)
	var a float64 = 6.739
	var b byte = 125
	var c int32 = 13
	var d bool = true
	var e string = "Hello こんにちは"
	var f []int16 = []int16{3, 4, 5, -50, -123, -34, 22, 0}
	var g []int32 = []int32{3, 2, 1, 7, 23, 56, 42, 51, 66, 76, 53, 123}
	var h []int64 = []int64{123, 25, 83, 8, -23, -47, 51, 0}
	var j []float32 = []float32{12.4, 25.5, 1.24, 3.44, 12.57, 0}
	var k []float64 = []float64{12.45675333444, 25.55677, 1.243232, 3.444666, 12.572424, 0}
	o.WriteObject(a)
	o.WriteObject(b)
	o.WriteObject(c)
	o.WriteObject(d)
	o.WriteObject(e)
	o.WriteObject(f)
	o.WriteObject(g)
	o.WriteObject(h)
	o.WriteObject(j)
	o.WriteObject(k)
	i := NewObjectDataInput(o.buffer, 0, service, false)

	if a != i.ReadObject() || b != i.ReadObject() || c != i.ReadObject() || d != i.ReadObject() ||
		e != i.ReadObject() || !reflect.DeepEqual(f, i.ReadObject()) || !reflect.DeepEqual(g, i.ReadObject()) ||
		!reflect.DeepEqual(h, i.ReadObject()) || !reflect.DeepEqual(j, i.ReadObject()) || !reflect.DeepEqual(k, i.ReadObject()) {
		t.Errorf("There is a problem in WriteObject() or ReadObject()!")
	}

}

func TestObjectDataInput_ReadByte(t *testing.T) {
	o := NewObjectDataOutput(9, nil, false)
	var a byte = 120
	var b byte = 176
	o.WriteByte(a)
	o.WriteByte(b)
	i := NewObjectDataInput(o.buffer, 1, nil, false)
	var expectedRet byte = b
	var ret byte
	ret, _ = i.ReadByte()
	if ret != expectedRet {
		t.Errorf("ReadByte() returns '%s' expected %s", ret, expectedRet)
	}
}

func TestObjectDataInput_ReadByteArray(t *testing.T) {
	var array []byte = []byte{3, 4, 5, 25, 123, 34, 52, 0}
	o := NewObjectDataOutput(0, nil, false)
	o.WriteByteArray(array)
	i := NewObjectDataInput(o.buffer, 0,nil, false)

	if !reflect.DeepEqual(array, i.ReadByteArray()) {
		t.Errorf("There is a problem in WriteByteArray() or ReadByteArray()!")
	}
}

func TestObjectDataInput_ReadBoolArray(t *testing.T) {
	var array []bool = []bool{true, false,true, true, false, false, false, true}
	o := NewObjectDataOutput(0, nil, false)
	o.WriteBoolArray(array)
	i := NewObjectDataInput(o.buffer, 0, nil, false)

	if !reflect.DeepEqual(array, i.ReadBoolArray()) {
		t.Errorf("There is a problem in WriteBoolArray() or ReadBoolArray()!")
	}
}

func TestObjectDataInput_ReadUTFArray(t *testing.T) {
	var array []string = []string{"aAüÜiİıIöÖşŞçÇ","akdha","üğpoıuişlk","üğpreÜaişfçxaaöc"}
	o := NewObjectDataOutput(0, nil, false)
	o.WriteUTFArray(array)
	i := NewObjectDataInput(o.buffer, 0,nil, false)

	if !reflect.DeepEqual(array, i.ReadUTFArray()) {
		t.Errorf("There is a problem in WriteUTFArray() or ReadUTFArray()!")
	}
}

func TestObjectDataInput_ReadInt16Array(t *testing.T) {
	var array []int16 = []int16{3, 4, 5, -50, -123, -34, 22, 0}
	o := NewObjectDataOutput(50, nil, false)
	o.WriteInt16Array(array)
	i := NewObjectDataInput(o.buffer, 0, nil, false)

	if !reflect.DeepEqual(array, i.ReadInt16Array()) {
		t.Errorf("There is a problem in WriteInt16Array() or ReadInt16Array()!")
	}
}

func TestObjectDataInput_ReadInt32Array(t *testing.T) {
	var array []int32 = []int32{321, 122, 14, 0, -123, -34, 67, 0}
	o := NewObjectDataOutput(50, nil, false)
	o.WriteInt32Array(array)
	i := NewObjectDataInput(o.buffer, 0, nil, false)

	if !reflect.DeepEqual(array, i.ReadInt32Array()) {
		t.Errorf("There is a problem in WriteInt32Array() or ReadInt32Array()!")
	}
}

func TestObjectDataInput_ReadInt64Array(t *testing.T) {
	var array []int64 = []int64{123, 25, 83, 8, -23, -47, 51, 0}
	o := NewObjectDataOutput(50, nil, false)
	o.WriteInt64Array(array)
	i := NewObjectDataInput(o.buffer, 0, nil, false)

	if !reflect.DeepEqual(array, i.ReadInt64Array()) {
		t.Errorf("There is a problem in WriteInt64Array() or ReadInt64Array()!")
	}
}

func TestObjectDataInput_ReadFloat32Array(t *testing.T) {
	var array []float32 = []float32{12.4, 25.5, 1.24, 3.44, 12.57, 0}
	o := NewObjectDataOutput(50, nil, false)
	o.WriteFloat32Array(array)
	i := NewObjectDataInput(o.buffer, 0, nil, false)

	if !reflect.DeepEqual(array, i.ReadFloat32Array()) {
		t.Errorf("There is a problem in WriteFloat32Array() or ReadFloat32Array()!")
	}
}

func TestObjectDataInput_ReadFloat64Array(t *testing.T) {
	var array []float64 = []float64{12.45675333444, 25.55677, 1.243232, 3.444666, 12.572424, 0}
	o := NewObjectDataOutput(50, nil, false)
	o.WriteFloat64Array(array)
	i := NewObjectDataInput(o.buffer, 0, nil, false)

	if !reflect.DeepEqual(array, i.ReadFloat64Array()) {
		t.Errorf("There is a problem in WriteFloat64Array() or ReadFloat64Array()!")
	}
}