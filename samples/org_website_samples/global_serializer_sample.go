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

package org_website_samples

import (
	. "github.com/hazelcast/hazelcast-go-client"
	. "github.com/hazelcast/hazelcast-go-client/config"
	. "github.com/hazelcast/hazelcast-go-client/serialization"
)

type GlobalSerializer struct {
}

func (*GlobalSerializer) Id() int32 {
	return 20
}

func (*GlobalSerializer) Read(input DataInput) (obj interface{}, err error) {
	// return MyFavoriteSerializer.deserialize(in)
	return
}

func (*GlobalSerializer) Write(output DataOutput, object interface{}) (err error) {
	// output.write(MyFavoriteSerializer.serialize(object))
	return
}

func globalSerializerSampleRun() {
	clientConfig := NewClientConfig()
	clientConfig.SerializationConfig().SetGlobalSerializer(&GlobalSerializer{})
	// Start the Hazelcast Client and connect to an already running Hazelcast Cluster on 127.0.0.1
	hz, _ := NewHazelcastClientWithConfig(clientConfig)

	//GlobalSerializer will serialize/deserialize all non-builtin types

	// Shutdown this hazelcast client
	hz.Shutdown()
}
