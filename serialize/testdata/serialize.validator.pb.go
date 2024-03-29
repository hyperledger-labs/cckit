// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: testdata/serialize.proto

package testdata

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Payment) Validate() error {
	if this.Deadline != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Deadline); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Deadline", err)
		}
	}
	return nil
}
