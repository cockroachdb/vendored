// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: errorspb/errors.proto

package errorspb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// EncodedError is the wire-encodable representation
// of an error (or error cause chain).
type EncodedError struct {
	// Types that are valid to be assigned to Error:
	//	*EncodedError_Leaf
	//	*EncodedError_Wrapper
	Error isEncodedError_Error `protobuf_oneof:"error"`
}

func (m *EncodedError) Reset()         { *m = EncodedError{} }
func (m *EncodedError) String() string { return proto.CompactTextString(m) }
func (*EncodedError) ProtoMessage()    {}
func (*EncodedError) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_b785792bfacbca0b, []int{0}
}
func (m *EncodedError) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EncodedError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *EncodedError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncodedError.Merge(dst, src)
}
func (m *EncodedError) XXX_Size() int {
	return m.Size()
}
func (m *EncodedError) XXX_DiscardUnknown() {
	xxx_messageInfo_EncodedError.DiscardUnknown(m)
}

var xxx_messageInfo_EncodedError proto.InternalMessageInfo

type isEncodedError_Error interface {
	isEncodedError_Error()
	MarshalTo([]byte) (int, error)
	Size() int
}

type EncodedError_Leaf struct {
	Leaf *EncodedErrorLeaf `protobuf:"bytes,1,opt,name=leaf,proto3,oneof"`
}
type EncodedError_Wrapper struct {
	Wrapper *EncodedWrapper `protobuf:"bytes,2,opt,name=wrapper,proto3,oneof"`
}

func (*EncodedError_Leaf) isEncodedError_Error()    {}
func (*EncodedError_Wrapper) isEncodedError_Error() {}

func (m *EncodedError) GetError() isEncodedError_Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *EncodedError) GetLeaf() *EncodedErrorLeaf {
	if x, ok := m.GetError().(*EncodedError_Leaf); ok {
		return x.Leaf
	}
	return nil
}

func (m *EncodedError) GetWrapper() *EncodedWrapper {
	if x, ok := m.GetError().(*EncodedError_Wrapper); ok {
		return x.Wrapper
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*EncodedError) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _EncodedError_OneofMarshaler, _EncodedError_OneofUnmarshaler, _EncodedError_OneofSizer, []interface{}{
		(*EncodedError_Leaf)(nil),
		(*EncodedError_Wrapper)(nil),
	}
}

func _EncodedError_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*EncodedError)
	// error
	switch x := m.Error.(type) {
	case *EncodedError_Leaf:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Leaf); err != nil {
			return err
		}
	case *EncodedError_Wrapper:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Wrapper); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("EncodedError.Error has unexpected type %T", x)
	}
	return nil
}

func _EncodedError_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*EncodedError)
	switch tag {
	case 1: // error.leaf
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(EncodedErrorLeaf)
		err := b.DecodeMessage(msg)
		m.Error = &EncodedError_Leaf{msg}
		return true, err
	case 2: // error.wrapper
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(EncodedWrapper)
		err := b.DecodeMessage(msg)
		m.Error = &EncodedError_Wrapper{msg}
		return true, err
	default:
		return false, nil
	}
}

func _EncodedError_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*EncodedError)
	// error
	switch x := m.Error.(type) {
	case *EncodedError_Leaf:
		s := proto.Size(x.Leaf)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *EncodedError_Wrapper:
		s := proto.Size(x.Wrapper)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// EncodedErrorLeaf is the wire-encodable representation
// of an error leaf.
type EncodedErrorLeaf struct {
	// The main error message (mandatory), that can be printed to human
	// users and may contain PII. This contains the value of the leaf
	// error's Error(), or using a registered encoder.
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	// The error details.
	Details EncodedErrorDetails `protobuf:"bytes,2,opt,name=details,proto3" json:"details"`
}

func (m *EncodedErrorLeaf) Reset()         { *m = EncodedErrorLeaf{} }
func (m *EncodedErrorLeaf) String() string { return proto.CompactTextString(m) }
func (*EncodedErrorLeaf) ProtoMessage()    {}
func (*EncodedErrorLeaf) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_b785792bfacbca0b, []int{1}
}
func (m *EncodedErrorLeaf) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EncodedErrorLeaf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *EncodedErrorLeaf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncodedErrorLeaf.Merge(dst, src)
}
func (m *EncodedErrorLeaf) XXX_Size() int {
	return m.Size()
}
func (m *EncodedErrorLeaf) XXX_DiscardUnknown() {
	xxx_messageInfo_EncodedErrorLeaf.DiscardUnknown(m)
}

var xxx_messageInfo_EncodedErrorLeaf proto.InternalMessageInfo

type EncodedErrorDetails struct {
	// The original fully qualified error type name (mandatory).
	// This is primarily used to print out error details
	// in error reports and Format().
	//
	// It is additionally used to populate the error mark
	// below when the family name is not known/set.
	// See the `markers` error package and the
	// RFC on error handling for details.
	OriginalTypeName string `protobuf:"bytes,1,opt,name=original_type_name,json=originalTypeName,proto3" json:"original_type_name,omitempty"`
	// The error mark. This is used to determine error equivalence and
	// identifying a decode function.
	// See the `markers` error package and the
	// RFC on error handling for details.
	ErrorTypeMark ErrorTypeMark `protobuf:"bytes,2,opt,name=error_type_mark,json=errorTypeMark,proto3" json:"error_type_mark"`
	// The reportable payload (optional), which is as descriptive as
	// possible but may not contain PII.
	//
	// This is extracted automatically using a registered encoder, if
	// any, or the SafeDetailer interface.
	ReportablePayload []string `protobuf:"bytes,3,rep,name=reportable_payload,json=reportablePayload,proto3" json:"reportable_payload,omitempty"`
	// An arbitrary payload that (presumably) encodes the
	// native error object. This is also optional.
	//
	// This is extracted automatically using a registered encoder, if
	// any.
	FullDetails *types.Any `protobuf:"bytes,4,opt,name=full_details,json=fullDetails,proto3" json:"full_details,omitempty"`
}

func (m *EncodedErrorDetails) Reset()         { *m = EncodedErrorDetails{} }
func (m *EncodedErrorDetails) String() string { return proto.CompactTextString(m) }
func (*EncodedErrorDetails) ProtoMessage()    {}
func (*EncodedErrorDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_b785792bfacbca0b, []int{2}
}
func (m *EncodedErrorDetails) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EncodedErrorDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *EncodedErrorDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncodedErrorDetails.Merge(dst, src)
}
func (m *EncodedErrorDetails) XXX_Size() int {
	return m.Size()
}
func (m *EncodedErrorDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_EncodedErrorDetails.DiscardUnknown(m)
}

var xxx_messageInfo_EncodedErrorDetails proto.InternalMessageInfo

// EncodedWrapper is the wire-encodable representation
// of an error wrapper.
type EncodedWrapper struct {
	// The cause error. Mandatory.
	Cause EncodedError `protobuf:"bytes,1,opt,name=cause,proto3" json:"cause"`
	// The wrapper message prefix (which may be empty). This
	// isbprinted before the cause's own message when
	// constructing a full message. This may contain PII.
	//
	// This is extracted automatically:
	//
	// - for wrappers that have a registered encoder,
	// - otherwise, when the wrapper's Error() has its cause's Error() as suffix.
	MessagePrefix string `protobuf:"bytes,2,opt,name=message_prefix,json=messagePrefix,proto3" json:"message_prefix,omitempty"`
	// The error details.
	Details EncodedErrorDetails `protobuf:"bytes,3,opt,name=details,proto3" json:"details"`
}

func (m *EncodedWrapper) Reset()         { *m = EncodedWrapper{} }
func (m *EncodedWrapper) String() string { return proto.CompactTextString(m) }
func (*EncodedWrapper) ProtoMessage()    {}
func (*EncodedWrapper) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_b785792bfacbca0b, []int{3}
}
func (m *EncodedWrapper) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EncodedWrapper) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *EncodedWrapper) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncodedWrapper.Merge(dst, src)
}
func (m *EncodedWrapper) XXX_Size() int {
	return m.Size()
}
func (m *EncodedWrapper) XXX_DiscardUnknown() {
	xxx_messageInfo_EncodedWrapper.DiscardUnknown(m)
}

var xxx_messageInfo_EncodedWrapper proto.InternalMessageInfo

// ErrorTypeMark identifies an error type for the purpose of determining
// error equivalences and looking up decoder functions.
type ErrorTypeMark struct {
	// The family name identifies the error type.
	// This is equal to original_type_name above in the common case, but
	// can be overridden when e.g. the package that defines the type
	// changes path.
	// This is the field also used for looking up a decode function.
	FamilyName string `protobuf:"bytes,1,opt,name=family_name,json=familyName,proto3" json:"family_name,omitempty"`
	// This marker string is used in combination with
	// the family name for the purpose of determining error equivalence.
	// This can be used to separate error instances that have the same type
	// into separate equivalence classes.
	// See the `markers` error package and the
	// RFC on error handling for details.
	Extension string `protobuf:"bytes,2,opt,name=extension,proto3" json:"extension,omitempty"`
}

func (m *ErrorTypeMark) Reset()         { *m = ErrorTypeMark{} }
func (m *ErrorTypeMark) String() string { return proto.CompactTextString(m) }
func (*ErrorTypeMark) ProtoMessage()    {}
func (*ErrorTypeMark) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_b785792bfacbca0b, []int{4}
}
func (m *ErrorTypeMark) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ErrorTypeMark) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *ErrorTypeMark) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorTypeMark.Merge(dst, src)
}
func (m *ErrorTypeMark) XXX_Size() int {
	return m.Size()
}
func (m *ErrorTypeMark) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorTypeMark.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorTypeMark proto.InternalMessageInfo

func init() {
	proto.RegisterType((*EncodedError)(nil), "cockroach.errorspb.EncodedError")
	proto.RegisterType((*EncodedErrorLeaf)(nil), "cockroach.errorspb.EncodedErrorLeaf")
	proto.RegisterType((*EncodedErrorDetails)(nil), "cockroach.errorspb.EncodedErrorDetails")
	proto.RegisterType((*EncodedWrapper)(nil), "cockroach.errorspb.EncodedWrapper")
	proto.RegisterType((*ErrorTypeMark)(nil), "cockroach.errorspb.ErrorTypeMark")
}
func (m *EncodedError) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EncodedError) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Error != nil {
		nn1, err := m.Error.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	return i, nil
}

func (m *EncodedError_Leaf) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Leaf != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintErrors(dAtA, i, uint64(m.Leaf.Size()))
		n2, err := m.Leaf.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}
func (m *EncodedError_Wrapper) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Wrapper != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintErrors(dAtA, i, uint64(m.Wrapper.Size()))
		n3, err := m.Wrapper.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}
func (m *EncodedErrorLeaf) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EncodedErrorLeaf) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.Message)))
		i += copy(dAtA[i:], m.Message)
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintErrors(dAtA, i, uint64(m.Details.Size()))
	n4, err := m.Details.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	return i, nil
}

func (m *EncodedErrorDetails) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EncodedErrorDetails) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.OriginalTypeName) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.OriginalTypeName)))
		i += copy(dAtA[i:], m.OriginalTypeName)
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintErrors(dAtA, i, uint64(m.ErrorTypeMark.Size()))
	n5, err := m.ErrorTypeMark.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	if len(m.ReportablePayload) > 0 {
		for _, s := range m.ReportablePayload {
			dAtA[i] = 0x1a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if m.FullDetails != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintErrors(dAtA, i, uint64(m.FullDetails.Size()))
		n6, err := m.FullDetails.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	return i, nil
}

func (m *EncodedWrapper) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EncodedWrapper) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintErrors(dAtA, i, uint64(m.Cause.Size()))
	n7, err := m.Cause.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n7
	if len(m.MessagePrefix) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.MessagePrefix)))
		i += copy(dAtA[i:], m.MessagePrefix)
	}
	dAtA[i] = 0x1a
	i++
	i = encodeVarintErrors(dAtA, i, uint64(m.Details.Size()))
	n8, err := m.Details.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n8
	return i, nil
}

func (m *ErrorTypeMark) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ErrorTypeMark) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.FamilyName) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.FamilyName)))
		i += copy(dAtA[i:], m.FamilyName)
	}
	if len(m.Extension) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintErrors(dAtA, i, uint64(len(m.Extension)))
		i += copy(dAtA[i:], m.Extension)
	}
	return i, nil
}

func encodeVarintErrors(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *EncodedError) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Error != nil {
		n += m.Error.Size()
	}
	return n
}

func (m *EncodedError_Leaf) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Leaf != nil {
		l = m.Leaf.Size()
		n += 1 + l + sovErrors(uint64(l))
	}
	return n
}
func (m *EncodedError_Wrapper) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Wrapper != nil {
		l = m.Wrapper.Size()
		n += 1 + l + sovErrors(uint64(l))
	}
	return n
}
func (m *EncodedErrorLeaf) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	l = m.Details.Size()
	n += 1 + l + sovErrors(uint64(l))
	return n
}

func (m *EncodedErrorDetails) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OriginalTypeName)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	l = m.ErrorTypeMark.Size()
	n += 1 + l + sovErrors(uint64(l))
	if len(m.ReportablePayload) > 0 {
		for _, s := range m.ReportablePayload {
			l = len(s)
			n += 1 + l + sovErrors(uint64(l))
		}
	}
	if m.FullDetails != nil {
		l = m.FullDetails.Size()
		n += 1 + l + sovErrors(uint64(l))
	}
	return n
}

func (m *EncodedWrapper) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Cause.Size()
	n += 1 + l + sovErrors(uint64(l))
	l = len(m.MessagePrefix)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	l = m.Details.Size()
	n += 1 + l + sovErrors(uint64(l))
	return n
}

func (m *ErrorTypeMark) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.FamilyName)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	l = len(m.Extension)
	if l > 0 {
		n += 1 + l + sovErrors(uint64(l))
	}
	return n
}

func sovErrors(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozErrors(x uint64) (n int) {
	return sovErrors(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EncodedError) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: EncodedError: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EncodedError: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leaf", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &EncodedErrorLeaf{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Error = &EncodedError_Leaf{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Wrapper", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &EncodedWrapper{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Error = &EncodedError_Wrapper{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipErrors(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthErrors
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *EncodedErrorLeaf) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: EncodedErrorLeaf: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EncodedErrorLeaf: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Details", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Details.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipErrors(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthErrors
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *EncodedErrorDetails) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: EncodedErrorDetails: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EncodedErrorDetails: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginalTypeName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OriginalTypeName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorTypeMark", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ErrorTypeMark.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReportablePayload", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReportablePayload = append(m.ReportablePayload, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FullDetails", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.FullDetails == nil {
				m.FullDetails = &types.Any{}
			}
			if err := m.FullDetails.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipErrors(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthErrors
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *EncodedWrapper) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: EncodedWrapper: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EncodedWrapper: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cause", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Cause.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessagePrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessagePrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Details", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Details.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipErrors(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthErrors
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ErrorTypeMark) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ErrorTypeMark: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ErrorTypeMark: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FamilyName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FamilyName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Extension", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthErrors
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Extension = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipErrors(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthErrors
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipErrors(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowErrors
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowErrors
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthErrors
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowErrors
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipErrors(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthErrors = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowErrors   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("errorspb/errors.proto", fileDescriptor_errors_b785792bfacbca0b) }

var fileDescriptor_errors_b785792bfacbca0b = []byte{
	// 482 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0x4d, 0x6b, 0x13, 0x51,
	0x14, 0x9d, 0x31, 0xa9, 0x31, 0x37, 0x4d, 0xad, 0xcf, 0x0a, 0xb1, 0xc8, 0x34, 0x06, 0xc5, 0x22,
	0x3a, 0x01, 0x5d, 0x08, 0x22, 0x82, 0xc1, 0x62, 0x17, 0x5a, 0xcb, 0x20, 0x08, 0x6e, 0x86, 0x97,
	0xc9, 0x9d, 0x71, 0xc8, 0x9b, 0x79, 0x8f, 0x37, 0x13, 0xec, 0xfc, 0x0b, 0xc1, 0xff, 0xe3, 0x3a,
	0xcb, 0x2e, 0xbb, 0x12, 0x4d, 0xfe, 0x86, 0x0b, 0x99, 0xf7, 0x41, 0x53, 0x2d, 0x76, 0xd1, 0xdd,
	0x9d, 0x73, 0xcf, 0x99, 0x7b, 0xde, 0x3d, 0x17, 0x6e, 0xa1, 0x94, 0x5c, 0x16, 0x62, 0x3c, 0xd4,
	0x85, 0x2f, 0x24, 0x2f, 0x39, 0x21, 0x11, 0x8f, 0xa6, 0x92, 0xd3, 0xe8, 0xb3, 0x6f, 0x09, 0xdb,
	0xb7, 0x13, 0xce, 0x13, 0x86, 0x43, 0xc5, 0x18, 0xcf, 0xe2, 0x21, 0xcd, 0x2b, 0x4d, 0xdf, 0xde,
	0x4a, 0x78, 0xc2, 0x55, 0x39, 0xac, 0x2b, 0x8d, 0x0e, 0xbe, 0xb9, 0xb0, 0xbe, 0x97, 0x47, 0x7c,
	0x82, 0x93, 0xbd, 0xfa, 0x27, 0xe4, 0x39, 0x34, 0x19, 0xd2, 0xb8, 0xe7, 0xf6, 0xdd, 0xdd, 0xce,
	0x93, 0x7b, 0xfe, 0xbf, 0x43, 0xfc, 0x55, 0xfe, 0x5b, 0xa4, 0xf1, 0xbe, 0x13, 0x28, 0x0d, 0x79,
	0x09, 0xad, 0x2f, 0x92, 0x0a, 0x81, 0xb2, 0x77, 0x45, 0xc9, 0x07, 0xff, 0x91, 0x7f, 0xd4, 0xcc,
	0x7d, 0x27, 0xb0, 0xa2, 0x51, 0x0b, 0xd6, 0x14, 0x6b, 0x30, 0x83, 0xcd, 0xbf, 0x87, 0x90, 0x1e,
	0xb4, 0x32, 0x2c, 0x0a, 0x9a, 0xa0, 0xf2, 0xd6, 0x0e, 0xec, 0x27, 0x79, 0x03, 0xad, 0x09, 0x96,
	0x34, 0x65, 0x85, 0x19, 0xfb, 0xe0, 0x22, 0xd7, 0xaf, 0x35, 0x7d, 0xd4, 0x9c, 0xff, 0xd8, 0x71,
	0x02, 0xab, 0x1e, 0xfc, 0x76, 0xe1, 0xe6, 0x39, 0x34, 0xf2, 0x08, 0x08, 0x97, 0x69, 0x92, 0xe6,
	0x94, 0x85, 0x65, 0x25, 0x30, 0xcc, 0x69, 0x66, 0x5d, 0x6c, 0xda, 0xce, 0x87, 0x4a, 0xe0, 0x01,
	0xcd, 0x90, 0xbc, 0x87, 0xeb, 0x6a, 0xa8, 0xa6, 0x66, 0x54, 0x4e, 0x8d, 0xad, 0xbb, 0xe7, 0xda,
	0xaa, 0x8b, 0x5a, 0xfb, 0x8e, 0xca, 0xa9, 0x31, 0xd4, 0xc5, 0x55, 0x90, 0x3c, 0x06, 0x22, 0x51,
	0x70, 0x59, 0xd2, 0x31, 0xc3, 0x50, 0xd0, 0x8a, 0x71, 0x3a, 0xe9, 0x35, 0xfa, 0x8d, 0xdd, 0x76,
	0x70, 0xe3, 0xb4, 0x73, 0xa8, 0x1b, 0xe4, 0x19, 0xac, 0xc7, 0x33, 0xc6, 0x42, 0xbb, 0x93, 0xa6,
	0x1a, 0xbe, 0xe5, 0xeb, 0xd3, 0xf0, 0xed, 0x69, 0xf8, 0xaf, 0xf2, 0x2a, 0xe8, 0xd4, 0x4c, 0xf3,
	0xcc, 0xc1, 0x77, 0x17, 0x36, 0xce, 0x86, 0x43, 0x5e, 0xc0, 0x5a, 0x44, 0x67, 0x05, 0x9a, 0x73,
	0xe8, 0x5f, 0xb4, 0x58, 0xf3, 0x00, 0x2d, 0x22, 0xf7, 0x61, 0xc3, 0x64, 0x14, 0x0a, 0x89, 0x71,
	0x7a, 0xa4, 0x16, 0xd1, 0x0e, 0xba, 0x06, 0x3d, 0x54, 0xe0, 0x6a, 0x7e, 0x8d, 0x4b, 0xe5, 0x77,
	0x00, 0xdd, 0x33, 0xeb, 0x24, 0x3b, 0xd0, 0x89, 0x69, 0x96, 0xb2, 0x6a, 0x35, 0x31, 0xd0, 0x90,
	0xca, 0xea, 0x0e, 0xb4, 0xf1, 0xa8, 0xc4, 0xbc, 0x48, 0x79, 0x6e, 0xcc, 0x9d, 0x02, 0xa3, 0x87,
	0xf3, 0x5f, 0x9e, 0x33, 0x5f, 0x78, 0xee, 0xf1, 0xc2, 0x73, 0x4f, 0x16, 0x9e, 0xfb, 0x73, 0xe1,
	0xb9, 0x5f, 0x97, 0x9e, 0x73, 0xbc, 0xf4, 0x9c, 0x93, 0xa5, 0xe7, 0x7c, 0xba, 0x66, 0xed, 0x8d,
	0xaf, 0xaa, 0xbd, 0x3e, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0x8e, 0x05, 0xe2, 0xd1, 0xad, 0x03,
	0x00, 0x00,
}
