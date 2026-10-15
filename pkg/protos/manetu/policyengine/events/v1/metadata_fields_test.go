//
//  Copyright © Manetu Inc. All rights reserved.
//

package eventsv1

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

// The phase 2 sequence fields must survive a wire round trip, and older encoders (fields
// absent) must decode as the "no sequence information" zero values.
func TestMetadataSequenceFieldsRoundTrip(t *testing.T) {
	in := &AccessRecord{Metadata: &AccessRecord_Metadata{
		Id: "abc", PdpInstanceId: "ns/pod-1", PdpEpoch: "epoch-1", Seq: 42,
	}}
	b, err := proto.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out AccessRecord
	if err := proto.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	m := out.GetMetadata()
	if m.GetPdpInstanceId() != "ns/pod-1" || m.GetPdpEpoch() != "epoch-1" || m.GetSeq() != 42 || m.GetId() != "abc" {
		t.Fatalf("got %+v", m)
	}

	old, _ := proto.Marshal(&AccessRecord{Metadata: &AccessRecord_Metadata{Id: "old"}})
	var dec AccessRecord
	if err := proto.Unmarshal(old, &dec); err != nil {
		t.Fatal(err)
	}
	if dec.GetMetadata().GetSeq() != 0 || dec.GetMetadata().GetPdpEpoch() != "" {
		t.Fatalf("absent fields must decode as zero values: %+v", dec.GetMetadata())
	}
}
