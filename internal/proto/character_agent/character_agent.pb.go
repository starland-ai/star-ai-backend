// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.24.3
// source: character_agent.proto

package character_agent

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ChatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role    string `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{0}
}

func (x *ChatMessage) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *ChatMessage) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type ImageMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Thumbnail string `protobuf:"bytes,2,opt,name=thumbnail,proto3" json:"thumbnail,omitempty"`
	Enable3D  bool   `protobuf:"varint,3,opt,name=enable3D,proto3" json:"enable3D,omitempty"`
	Id        string `protobuf:"bytes,4,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ImageMeta) Reset() {
	*x = ImageMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImageMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImageMeta) ProtoMessage() {}

func (x *ImageMeta) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImageMeta.ProtoReflect.Descriptor instead.
func (*ImageMeta) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{1}
}

func (x *ImageMeta) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ImageMeta) GetThumbnail() string {
	if x != nil {
		return x.Thumbnail
	}
	return ""
}

func (x *ImageMeta) GetEnable3D() bool {
	if x != nil {
		return x.Enable3D
	}
	return false
}

func (x *ImageMeta) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CharacterSetting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Gender       string            `protobuf:"bytes,2,opt,name=gender,proto3" json:"gender,omitempty"`
	Description  string            `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Tags         map[string]string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Introduction string            `protobuf:"bytes,5,opt,name=introduction,proto3" json:"introduction,omitempty"`
}

func (x *CharacterSetting) Reset() {
	*x = CharacterSetting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CharacterSetting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CharacterSetting) ProtoMessage() {}

func (x *CharacterSetting) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CharacterSetting.ProtoReflect.Descriptor instead.
func (*CharacterSetting) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{2}
}

func (x *CharacterSetting) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CharacterSetting) GetGender() string {
	if x != nil {
		return x.Gender
	}
	return ""
}

func (x *CharacterSetting) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CharacterSetting) GetTags() map[string]string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *CharacterSetting) GetIntroduction() string {
	if x != nil {
		return x.Introduction
	}
	return ""
}

type ChatCompletionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string       `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Message   *ChatMessage `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ChatCompletionRequest) Reset() {
	*x = ChatCompletionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatCompletionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatCompletionRequest) ProtoMessage() {}

func (x *ChatCompletionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatCompletionRequest.ProtoReflect.Descriptor instead.
func (*ChatCompletionRequest) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{3}
}

func (x *ChatCompletionRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *ChatCompletionRequest) GetMessage() *ChatMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

type ChatCompletionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code        uint32         `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"` // resp code, success: 0, failed: <> 0
	ErrMsg      string         `protobuf:"bytes,2,opt,name=err_msg,json=errMsg,proto3" json:"err_msg,omitempty"`
	Messages    []*ChatMessage `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages,omitempty"`
	ImageMetas  []*ImageMeta   `protobuf:"bytes,4,rep,name=image_metas,json=imageMetas,proto3" json:"image_metas,omitempty"`
	NeedConfirm bool           `protobuf:"varint,5,opt,name=need_confirm,json=needConfirm,proto3" json:"need_confirm,omitempty"`
	ConfirmType string         `protobuf:"bytes,6,opt,name=confirm_type,json=confirmType,proto3" json:"confirm_type,omitempty"`
}

func (x *ChatCompletionResponse) Reset() {
	*x = ChatCompletionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatCompletionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatCompletionResponse) ProtoMessage() {}

func (x *ChatCompletionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatCompletionResponse.ProtoReflect.Descriptor instead.
func (*ChatCompletionResponse) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{4}
}

func (x *ChatCompletionResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ChatCompletionResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *ChatCompletionResponse) GetMessages() []*ChatMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *ChatCompletionResponse) GetImageMetas() []*ImageMeta {
	if x != nil {
		return x.ImageMetas
	}
	return nil
}

func (x *ChatCompletionResponse) GetNeedConfirm() bool {
	if x != nil {
		return x.NeedConfirm
	}
	return false
}

func (x *ChatCompletionResponse) GetConfirmType() string {
	if x != nil {
		return x.ConfirmType
	}
	return ""
}

type ConfirmCharacterSettingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *ConfirmCharacterSettingRequest) Reset() {
	*x = ConfirmCharacterSettingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfirmCharacterSettingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfirmCharacterSettingRequest) ProtoMessage() {}

func (x *ConfirmCharacterSettingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfirmCharacterSettingRequest.ProtoReflect.Descriptor instead.
func (*ConfirmCharacterSettingRequest) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{5}
}

func (x *ConfirmCharacterSettingRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type ConfirmCharacterSettingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code             uint32            `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"` // resp code, success: 0, failed: <> 0
	ErrMsg           string            `protobuf:"bytes,2,opt,name=err_msg,json=errMsg,proto3" json:"err_msg,omitempty"`
	Messages         []*ChatMessage    `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages,omitempty"`
	NeedConfirm      bool              `protobuf:"varint,4,opt,name=need_confirm,json=needConfirm,proto3" json:"need_confirm,omitempty"`
	ConfirmType      string            `protobuf:"bytes,5,opt,name=confirm_type,json=confirmType,proto3" json:"confirm_type,omitempty"`
	CharacterSetting *CharacterSetting `protobuf:"bytes,6,opt,name=character_setting,json=characterSetting,proto3" json:"character_setting,omitempty"`
}

func (x *ConfirmCharacterSettingResponse) Reset() {
	*x = ConfirmCharacterSettingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfirmCharacterSettingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfirmCharacterSettingResponse) ProtoMessage() {}

func (x *ConfirmCharacterSettingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfirmCharacterSettingResponse.ProtoReflect.Descriptor instead.
func (*ConfirmCharacterSettingResponse) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{6}
}

func (x *ConfirmCharacterSettingResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ConfirmCharacterSettingResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *ConfirmCharacterSettingResponse) GetMessages() []*ChatMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *ConfirmCharacterSettingResponse) GetNeedConfirm() bool {
	if x != nil {
		return x.NeedConfirm
	}
	return false
}

func (x *ConfirmCharacterSettingResponse) GetConfirmType() string {
	if x != nil {
		return x.ConfirmType
	}
	return ""
}

func (x *ConfirmCharacterSettingResponse) GetCharacterSetting() *CharacterSetting {
	if x != nil {
		return x.CharacterSetting
	}
	return nil
}

type ChatCompletionStreamResponseChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkType         uint32 `protobuf:"varint,1,opt,name=chunk_type,json=chunkType,proto3" json:"chunk_type,omitempty"` // chat chunk = 1; image chunk = 2; need confirm chunk = 3; setting chunk = 4;
	ChunkSessionIndex uint32 `protobuf:"varint,2,opt,name=chunk_session_index,json=chunkSessionIndex,proto3" json:"chunk_session_index,omitempty"`
	// { "name": "xxxx", "gender": "xxxx", "description": "xxxx" }
	// { "name": "xxxx", "gender": "xxxx", "description": "xxxx123455" }
	ChatChunk        *ChatMessage `protobuf:"bytes,3,opt,name=chat_chunk,json=chatChunk,proto3" json:"chat_chunk,omitempty"` // I am
	ImageChunk       []*ImageMeta `protobuf:"bytes,4,rep,name=image_chunk,json=imageChunk,proto3" json:"image_chunk,omitempty"`
	NeedConfirmChunk bool         `protobuf:"varint,5,opt,name=need_confirm_chunk,json=needConfirmChunk,proto3" json:"need_confirm_chunk,omitempty"`
	SettingChunk     string       `protobuf:"bytes,6,opt,name=setting_chunk,json=settingChunk,proto3" json:"setting_chunk,omitempty"`
}

func (x *ChatCompletionStreamResponseChunk) Reset() {
	*x = ChatCompletionStreamResponseChunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatCompletionStreamResponseChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatCompletionStreamResponseChunk) ProtoMessage() {}

func (x *ChatCompletionStreamResponseChunk) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatCompletionStreamResponseChunk.ProtoReflect.Descriptor instead.
func (*ChatCompletionStreamResponseChunk) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{7}
}

func (x *ChatCompletionStreamResponseChunk) GetChunkType() uint32 {
	if x != nil {
		return x.ChunkType
	}
	return 0
}

func (x *ChatCompletionStreamResponseChunk) GetChunkSessionIndex() uint32 {
	if x != nil {
		return x.ChunkSessionIndex
	}
	return 0
}

func (x *ChatCompletionStreamResponseChunk) GetChatChunk() *ChatMessage {
	if x != nil {
		return x.ChatChunk
	}
	return nil
}

func (x *ChatCompletionStreamResponseChunk) GetImageChunk() []*ImageMeta {
	if x != nil {
		return x.ImageChunk
	}
	return nil
}

func (x *ChatCompletionStreamResponseChunk) GetNeedConfirmChunk() bool {
	if x != nil {
		return x.NeedConfirmChunk
	}
	return false
}

func (x *ChatCompletionStreamResponseChunk) GetSettingChunk() string {
	if x != nil {
		return x.SettingChunk
	}
	return ""
}

type ChatCompletionStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code   uint32                             `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"` // resp code, success: 0, failed: <> 0
	ErrMsg string                             `protobuf:"bytes,2,opt,name=err_msg,json=errMsg,proto3" json:"err_msg,omitempty"`
	Chunk  *ChatCompletionStreamResponseChunk `protobuf:"bytes,3,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (x *ChatCompletionStreamResponse) Reset() {
	*x = ChatCompletionStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_character_agent_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatCompletionStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatCompletionStreamResponse) ProtoMessage() {}

func (x *ChatCompletionStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_character_agent_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatCompletionStreamResponse.ProtoReflect.Descriptor instead.
func (*ChatCompletionStreamResponse) Descriptor() ([]byte, []int) {
	return file_character_agent_proto_rawDescGZIP(), []int{8}
}

func (x *ChatCompletionStreamResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ChatCompletionStreamResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *ChatCompletionStreamResponse) GetChunk() *ChatCompletionStreamResponseChunk {
	if x != nil {
		return x.Chunk
	}
	return nil
}

var File_character_agent_proto protoreflect.FileDescriptor

var file_character_agent_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x22, 0x3b,
	0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x69, 0x0a, 0x09, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e,
	0x61, 0x62, 0x6c, 0x65, 0x33, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x65, 0x6e,
	0x61, 0x62, 0x6c, 0x65, 0x33, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xf4, 0x01, 0x0a, 0x10, 0x43, 0x68, 0x61, 0x72, 0x61,
	0x63, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x35, 0x0a, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x2e, 0x54, 0x61, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73,
	0x12, 0x22, 0x0a, 0x0c, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x37, 0x0a, 0x09, 0x54, 0x61, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x64, 0x0a,
	0x15, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43,
	0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0xee, 0x01, 0x0a, 0x16, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x2e, 0x0a, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x31, 0x0a, 0x0b, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x4d, 0x65,
	0x74, 0x61, 0x52, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x73, 0x12, 0x21,
	0x0a, 0x0c, 0x6e, 0x65, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x6e, 0x65, 0x65, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72,
	0x6d, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x22, 0x3f, 0x0a, 0x1e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x43,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x8a, 0x02, 0x0a, 0x1f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72,
	0x6d, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x17, 0x0a,
	0x07, 0x65, 0x72, 0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x2e, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x65, 0x65, 0x64, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x6e, 0x65,
	0x65, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x72, 0x6d, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x44, 0x0a, 0x11,
	0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e,
	0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x22, 0xab, 0x02, 0x0a, 0x21, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x53, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x31, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x74, 0x5f,
	0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x09, 0x63, 0x68, 0x61, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x31, 0x0a, 0x0b, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x5f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x74,
	0x61, 0x52, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x2c, 0x0a,
	0x12, 0x6e, 0x65, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x5f, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x6e, 0x65, 0x65, 0x64, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x23, 0x0a, 0x0d, 0x73,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x22, 0x8b, 0x01, 0x0a, 0x1c, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x5f, 0x6d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x3e,
	0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x32, 0xa5,
	0x02, 0x0a, 0x05, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x50, 0x0a, 0x0f, 0x43, 0x68, 0x61, 0x74,
	0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1c, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5e, 0x0a, 0x15, 0x43, 0x68,
	0x61, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x12, 0x1c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74,
	0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x23, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x43, 0x6f,
	0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x6a, 0x0a, 0x17, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x72, 0x6d, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x25, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x72, 0x6d, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x43, 0x68, 0x61, 0x72,
	0x61, 0x63, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x18, 0x5a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_character_agent_proto_rawDescOnce sync.Once
	file_character_agent_proto_rawDescData = file_character_agent_proto_rawDesc
)

func file_character_agent_proto_rawDescGZIP() []byte {
	file_character_agent_proto_rawDescOnce.Do(func() {
		file_character_agent_proto_rawDescData = protoimpl.X.CompressGZIP(file_character_agent_proto_rawDescData)
	})
	return file_character_agent_proto_rawDescData
}

var file_character_agent_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_character_agent_proto_goTypes = []interface{}{
	(*ChatMessage)(nil),                       // 0: agent.ChatMessage
	(*ImageMeta)(nil),                         // 1: agent.ImageMeta
	(*CharacterSetting)(nil),                  // 2: agent.CharacterSetting
	(*ChatCompletionRequest)(nil),             // 3: agent.ChatCompletionRequest
	(*ChatCompletionResponse)(nil),            // 4: agent.ChatCompletionResponse
	(*ConfirmCharacterSettingRequest)(nil),    // 5: agent.ConfirmCharacterSettingRequest
	(*ConfirmCharacterSettingResponse)(nil),   // 6: agent.ConfirmCharacterSettingResponse
	(*ChatCompletionStreamResponseChunk)(nil), // 7: agent.ChatCompletionStreamResponseChunk
	(*ChatCompletionStreamResponse)(nil),      // 8: agent.ChatCompletionStreamResponse
	nil,                                       // 9: agent.CharacterSetting.TagsEntry
}
var file_character_agent_proto_depIdxs = []int32{
	9,  // 0: agent.CharacterSetting.tags:type_name -> agent.CharacterSetting.TagsEntry
	0,  // 1: agent.ChatCompletionRequest.message:type_name -> agent.ChatMessage
	0,  // 2: agent.ChatCompletionResponse.messages:type_name -> agent.ChatMessage
	1,  // 3: agent.ChatCompletionResponse.image_metas:type_name -> agent.ImageMeta
	0,  // 4: agent.ConfirmCharacterSettingResponse.messages:type_name -> agent.ChatMessage
	2,  // 5: agent.ConfirmCharacterSettingResponse.character_setting:type_name -> agent.CharacterSetting
	0,  // 6: agent.ChatCompletionStreamResponseChunk.chat_chunk:type_name -> agent.ChatMessage
	1,  // 7: agent.ChatCompletionStreamResponseChunk.image_chunk:type_name -> agent.ImageMeta
	7,  // 8: agent.ChatCompletionStreamResponse.chunk:type_name -> agent.ChatCompletionStreamResponseChunk
	3,  // 9: agent.Agent.ChatCompletions:input_type -> agent.ChatCompletionRequest
	3,  // 10: agent.Agent.ChatCompletionsStream:input_type -> agent.ChatCompletionRequest
	5,  // 11: agent.Agent.ConfirmCharacterSetting:input_type -> agent.ConfirmCharacterSettingRequest
	4,  // 12: agent.Agent.ChatCompletions:output_type -> agent.ChatCompletionResponse
	8,  // 13: agent.Agent.ChatCompletionsStream:output_type -> agent.ChatCompletionStreamResponse
	6,  // 14: agent.Agent.ConfirmCharacterSetting:output_type -> agent.ConfirmCharacterSettingResponse
	12, // [12:15] is the sub-list for method output_type
	9,  // [9:12] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_character_agent_proto_init() }
func file_character_agent_proto_init() {
	if File_character_agent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_character_agent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_character_agent_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImageMeta); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_character_agent_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CharacterSetting); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_character_agent_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatCompletionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_character_agent_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatCompletionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_character_agent_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfirmCharacterSettingRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_character_agent_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfirmCharacterSettingResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_character_agent_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatCompletionStreamResponseChunk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_character_agent_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatCompletionStreamResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_character_agent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_character_agent_proto_goTypes,
		DependencyIndexes: file_character_agent_proto_depIdxs,
		MessageInfos:      file_character_agent_proto_msgTypes,
	}.Build()
	File_character_agent_proto = out.File
	file_character_agent_proto_rawDesc = nil
	file_character_agent_proto_goTypes = nil
	file_character_agent_proto_depIdxs = nil
}
