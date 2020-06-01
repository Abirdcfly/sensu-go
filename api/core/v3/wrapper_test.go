package v3

import (
	"encoding/json"
	fmt "fmt"
	"testing"

	proto "github.com/golang/protobuf/proto"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
)

func TestUnmarshalWrapper(t *testing.T) {
	v := &EntityConfig{
		Metadata: &corev2.ObjectMeta{
			Name: "foo",
		},
	}
	a, err := proto.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	w := &Wrapper{
		TypeMeta: &corev2.TypeMeta{
			APIVersion: "core/v3",
			Type:       "entityconfig",
		},
		Metadata: &corev2.ObjectMeta{
			Name: "foo",
		},
		Value: a,
	}
	b, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", string(b))
}