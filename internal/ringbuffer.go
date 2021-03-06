// Copyright (c) 2008-2018, Hazelcast, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client/core"
	. "github.com/hazelcast/hazelcast-go-client/internal/protocol"
	. "github.com/hazelcast/hazelcast-go-client/internal/serialization"
)

type RingbufferProxy struct {
	*partitionSpecificProxy
	capacity int64
}

func newRingbufferProxy(client *HazelcastClient, serviceName *string, name *string) (*RingbufferProxy, error) {
	parSpecProxy, err := newPartitionSpecificProxy(client, serviceName, name)
	if err != nil {
		return nil, err
	}
	return &RingbufferProxy{parSpecProxy, -1}, nil
}

func (rp *RingbufferProxy) Capacity() (capacity int64, err error) {
	if rp.capacity == -1 {
		request := RingbufferCapacityEncodeRequest(rp.name)
		responseMessage, err := rp.invoke(request)
		capacity, err := rp.decodeToInt64AndError(responseMessage, err, RingbufferCapacityDecodeResponse)
		if err != nil {
			return 0, nil
		}
		rp.capacity = capacity
	}
	return rp.capacity, nil
}
func (rp *RingbufferProxy) Size() (size int64, err error) {
	request := RingbufferSizeEncodeRequest(rp.name)
	responseMessage, err := rp.invoke(request)
	return rp.decodeToInt64AndError(responseMessage, err, RingbufferSizeDecodeResponse)
}

func (rp *RingbufferProxy) TailSequence() (tailSequence int64, err error) {
	request := RingbufferTailSequenceEncodeRequest(rp.name)
	responseMessage, err := rp.invoke(request)
	return rp.decodeToInt64AndError(responseMessage, err, RingbufferTailSequenceDecodeResponse)
}

func (rp *RingbufferProxy) HeadSequence() (headSequence int64, err error) {
	request := RingbufferHeadSequenceEncodeRequest(rp.name)
	responseMessage, err := rp.invoke(request)
	return rp.decodeToInt64AndError(responseMessage, err, RingbufferHeadSequenceDecodeResponse)
}

func (rp *RingbufferProxy) RemainingCapacity() (remainingCapacity int64, err error) {
	request := RingbufferRemainingCapacityEncodeRequest(rp.name)
	responseMessage, err := rp.invoke(request)
	return rp.decodeToInt64AndError(responseMessage, err, RingbufferRemainingCapacityDecodeResponse)
}

func (rp *RingbufferProxy) Add(item interface{}, overflowPolicy core.OverflowPolicy) (sequence int64, err error) {
	itemData, err := rp.validateAndSerialize(item)
	if err != nil {
		return
	}
	request := RingbufferAddEncodeRequest(rp.name, int32(overflowPolicy.Policy()), itemData)
	responseMessage, err := rp.invoke(request)
	return rp.decodeToInt64AndError(responseMessage, err, RingbufferAddDecodeResponse)
}

func (rp *RingbufferProxy) AddAll(items []interface{}, overflowPolicy core.OverflowPolicy) (lastSequence int64, err error) {
	itemsData, err := rp.validateAndSerializeSlice(items)
	if err != nil {
		return
	}
	request := RingbufferAddAllEncodeRequest(rp.name, itemsData, int32(overflowPolicy.Policy()))
	responseMessage, err := rp.invoke(request)
	return rp.decodeToInt64AndError(responseMessage, err, RingbufferAddAllDecodeResponse)
}

func (rp *RingbufferProxy) ReadOne(sequence int64) (item interface{}, err error) {
	if err = rp.validateSequenceNotNegative(sequence, "sequence"); err != nil {
		return
	}
	request := RingbufferReadOneEncodeRequest(rp.name, sequence)
	responseMessage, err := rp.invoke(request)
	return rp.decodeToObjectAndError(responseMessage, err, RingbufferReadOneDecodeResponse)
}

func (rp *RingbufferProxy) ReadMany(startSequence int64, minCount int32, maxCount int32, filter interface{}) (readResultSet core.ReadResultSet, err error) {
	filterData, err := rp.toData(filter)
	if err != nil {
		return
	}
	if err = rp.validateSequenceNotNegative(startSequence, "start sequence"); err != nil {
		return
	}
	if err = rp.checkCounts(minCount, maxCount); err != nil {
		return
	}
	request := RingbufferReadManyEncodeRequest(rp.name, startSequence, minCount, maxCount, filterData)
	responseMessage, err := rp.invoke(request)
	if err != nil {
		return
	}
	readCount, itemsData, itemSeqs := RingbufferReadManyDecodeResponse(responseMessage)()
	return NewLazyReadResultSet(readCount, itemsData, itemSeqs, rp.client.SerializationService), nil
}

func (rp *RingbufferProxy) validateSequenceNotNegative(value int64, argName string) (err error) {
	if value < 0 {
		err = core.NewHazelcastIllegalArgumentError(fmt.Sprintf("%v %v can't be smaller than 0", argName, value), nil)
	}
	return
}

func (rp *RingbufferProxy) checkCounts(minCount int32, maxCount int32) (err error) {
	if minCount < 0 {
		return core.NewHazelcastIllegalArgumentError(fmt.Sprintf("min count %v can't be smaller than 0", minCount), nil)
	}
	if minCount > maxCount {
		return core.NewHazelcastIllegalArgumentError(fmt.Sprintf("min count %v can't be larger than max count %v", minCount, maxCount), nil)
	}
	return
}

type LazyReadResultSet struct {
	readCount int32
	// This slice includes both data and de-serialized objects.
	lazyItems            []interface{}
	itemSequences        []int64
	serializationService *SerializationService
}

const sequenceUnavailable int64 = -1

func NewLazyReadResultSet(readCount int32, itemsData []*Data, itemSeqs []int64, ss *SerializationService) (rs *LazyReadResultSet) {
	rs = &LazyReadResultSet{readCount: readCount, itemSequences: itemSeqs, serializationService: ss}
	lazyItems := make([]interface{}, len(itemsData))
	for i, itemData := range itemsData {
		lazyItems[i] = itemData
	}
	rs.lazyItems = lazyItems
	return
}

func (rs *LazyReadResultSet) ReadCount() int32 {
	return rs.readCount
}

func (rs *LazyReadResultSet) Get(index int32) (result interface{}, err error) {
	if err = rs.rangeCheck(index); err != nil {
		return
	}
	if itemData, ok := rs.lazyItems[index].(*Data); ok {
		item, err := rs.serializationService.ToObject(itemData)
		if err != nil {
			return nil, err
		}
		rs.lazyItems[index] = item
	}
	return rs.lazyItems[index], nil
}

func (rs *LazyReadResultSet) Sequence(index int32) (sequence int64, err error) {
	if rs.itemSequences == nil {
		return sequenceUnavailable, nil
	}
	if err = rs.rangeCheck(index); err != nil {
		return
	}
	return rs.itemSequences[index], nil
}

func (rs *LazyReadResultSet) Size() int32 {
	return int32(len(rs.lazyItems))
}

func (rs *LazyReadResultSet) rangeCheck(index int32) (err error) {
	size := len(rs.lazyItems)
	if index < 0 || index >= int32(size) {
		err = core.NewHazelcastIllegalArgumentError(fmt.Sprintf("index=", index, "size=", size), nil)
	}
	return
}
