// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: cockroach/proto/status.proto

#define INTERNAL_SUPPRESS_PROTOBUF_FIELD_DEPRECATION
#include "cockroach/proto/status.pb.h"

#include <algorithm>

#include <google/protobuf/stubs/common.h>
#include <google/protobuf/stubs/once.h>
#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/wire_format_lite_inl.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/generated_message_reflection.h>
#include <google/protobuf/reflection_ops.h>
#include <google/protobuf/wire_format.h>
// @@protoc_insertion_point(includes)

namespace cockroach {
namespace proto {

namespace {

const ::google::protobuf::Descriptor* StoreStatus_descriptor_ = NULL;
const ::google::protobuf::internal::GeneratedMessageReflection*
  StoreStatus_reflection_ = NULL;

}  // namespace


void protobuf_AssignDesc_cockroach_2fproto_2fstatus_2eproto() {
  protobuf_AddDesc_cockroach_2fproto_2fstatus_2eproto();
  const ::google::protobuf::FileDescriptor* file =
    ::google::protobuf::DescriptorPool::generated_pool()->FindFileByName(
      "cockroach/proto/status.proto");
  GOOGLE_CHECK(file != NULL);
  StoreStatus_descriptor_ = file->message_type(0);
  static const int StoreStatus_offsets_[6] = {
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(StoreStatus, store_id_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(StoreStatus, node_id_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(StoreStatus, range_count_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(StoreStatus, started_at_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(StoreStatus, updated_at_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(StoreStatus, stats_),
  };
  StoreStatus_reflection_ =
    new ::google::protobuf::internal::GeneratedMessageReflection(
      StoreStatus_descriptor_,
      StoreStatus::default_instance_,
      StoreStatus_offsets_,
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(StoreStatus, _has_bits_[0]),
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(StoreStatus, _unknown_fields_),
      -1,
      ::google::protobuf::DescriptorPool::generated_pool(),
      ::google::protobuf::MessageFactory::generated_factory(),
      sizeof(StoreStatus));
}

namespace {

GOOGLE_PROTOBUF_DECLARE_ONCE(protobuf_AssignDescriptors_once_);
inline void protobuf_AssignDescriptorsOnce() {
  ::google::protobuf::GoogleOnceInit(&protobuf_AssignDescriptors_once_,
                 &protobuf_AssignDesc_cockroach_2fproto_2fstatus_2eproto);
}

void protobuf_RegisterTypes(const ::std::string&) {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedMessage(
    StoreStatus_descriptor_, &StoreStatus::default_instance());
}

}  // namespace

void protobuf_ShutdownFile_cockroach_2fproto_2fstatus_2eproto() {
  delete StoreStatus::default_instance_;
  delete StoreStatus_reflection_;
}

void protobuf_AddDesc_cockroach_2fproto_2fstatus_2eproto() {
  static bool already_here = false;
  if (already_here) return;
  already_here = true;
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  ::cockroach::proto::protobuf_AddDesc_cockroach_2fproto_2fdata_2eproto();
  ::gogoproto::protobuf_AddDesc_gogoproto_2fgogo_2eproto();
  ::google::protobuf::DescriptorPool::InternalAddGeneratedFile(
    "\n\034cockroach/proto/status.proto\022\017cockroac"
    "h.proto\032\032cockroach/proto/data.proto\032\024gog"
    "oproto/gogo.proto\"\346\001\n\013StoreStatus\022,\n\010sto"
    "re_id\030\001 \001(\005B\032\310\336\037\000\342\336\037\007StoreID\332\336\037\007StoreID\022"
    ")\n\007node_id\030\002 \001(\005B\030\310\336\037\000\342\336\037\006NodeID\332\336\037\006Node"
    "ID\022\031\n\013range_count\030\003 \001(\005B\004\310\336\037\000\022\030\n\nstarted"
    "_at\030\004 \001(\003B\004\310\336\037\000\022\030\n\nupdated_at\030\005 \001(\003B\004\310\336\037"
    "\000\022/\n\005stats\030\006 \001(\0132\032.cockroach.proto.MVCCS"
    "tatsB\004\310\336\037\000B\023Z\005proto\340\342\036\001\310\342\036\001\320\342\036\001", 351);
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedFile(
    "cockroach/proto/status.proto", &protobuf_RegisterTypes);
  StoreStatus::default_instance_ = new StoreStatus();
  StoreStatus::default_instance_->InitAsDefaultInstance();
  ::google::protobuf::internal::OnShutdown(&protobuf_ShutdownFile_cockroach_2fproto_2fstatus_2eproto);
}

// Force AddDescriptors() to be called at static initialization time.
struct StaticDescriptorInitializer_cockroach_2fproto_2fstatus_2eproto {
  StaticDescriptorInitializer_cockroach_2fproto_2fstatus_2eproto() {
    protobuf_AddDesc_cockroach_2fproto_2fstatus_2eproto();
  }
} static_descriptor_initializer_cockroach_2fproto_2fstatus_2eproto_;

// ===================================================================

#ifndef _MSC_VER
const int StoreStatus::kStoreIdFieldNumber;
const int StoreStatus::kNodeIdFieldNumber;
const int StoreStatus::kRangeCountFieldNumber;
const int StoreStatus::kStartedAtFieldNumber;
const int StoreStatus::kUpdatedAtFieldNumber;
const int StoreStatus::kStatsFieldNumber;
#endif  // !_MSC_VER

StoreStatus::StoreStatus()
  : ::google::protobuf::Message() {
  SharedCtor();
  // @@protoc_insertion_point(constructor:cockroach.proto.StoreStatus)
}

void StoreStatus::InitAsDefaultInstance() {
  stats_ = const_cast< ::cockroach::proto::MVCCStats*>(&::cockroach::proto::MVCCStats::default_instance());
}

StoreStatus::StoreStatus(const StoreStatus& from)
  : ::google::protobuf::Message() {
  SharedCtor();
  MergeFrom(from);
  // @@protoc_insertion_point(copy_constructor:cockroach.proto.StoreStatus)
}

void StoreStatus::SharedCtor() {
  _cached_size_ = 0;
  store_id_ = 0;
  node_id_ = 0;
  range_count_ = 0;
  started_at_ = GOOGLE_LONGLONG(0);
  updated_at_ = GOOGLE_LONGLONG(0);
  stats_ = NULL;
  ::memset(_has_bits_, 0, sizeof(_has_bits_));
}

StoreStatus::~StoreStatus() {
  // @@protoc_insertion_point(destructor:cockroach.proto.StoreStatus)
  SharedDtor();
}

void StoreStatus::SharedDtor() {
  if (this != default_instance_) {
    delete stats_;
  }
}

void StoreStatus::SetCachedSize(int size) const {
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
}
const ::google::protobuf::Descriptor* StoreStatus::descriptor() {
  protobuf_AssignDescriptorsOnce();
  return StoreStatus_descriptor_;
}

const StoreStatus& StoreStatus::default_instance() {
  if (default_instance_ == NULL) protobuf_AddDesc_cockroach_2fproto_2fstatus_2eproto();
  return *default_instance_;
}

StoreStatus* StoreStatus::default_instance_ = NULL;

StoreStatus* StoreStatus::New() const {
  return new StoreStatus;
}

void StoreStatus::Clear() {
#define OFFSET_OF_FIELD_(f) (reinterpret_cast<char*>(      \
  &reinterpret_cast<StoreStatus*>(16)->f) - \
   reinterpret_cast<char*>(16))

#define ZR_(first, last) do {                              \
    size_t f = OFFSET_OF_FIELD_(first);                    \
    size_t n = OFFSET_OF_FIELD_(last) - f + sizeof(last);  \
    ::memset(&first, 0, n);                                \
  } while (0)

  if (_has_bits_[0 / 32] & 63) {
    ZR_(store_id_, updated_at_);
    range_count_ = 0;
    if (has_stats()) {
      if (stats_ != NULL) stats_->::cockroach::proto::MVCCStats::Clear();
    }
  }

#undef OFFSET_OF_FIELD_
#undef ZR_

  ::memset(_has_bits_, 0, sizeof(_has_bits_));
  mutable_unknown_fields()->Clear();
}

bool StoreStatus::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!(EXPRESSION)) goto failure
  ::google::protobuf::uint32 tag;
  // @@protoc_insertion_point(parse_start:cockroach.proto.StoreStatus)
  for (;;) {
    ::std::pair< ::google::protobuf::uint32, bool> p = input->ReadTagWithCutoff(127);
    tag = p.first;
    if (!p.second) goto handle_unusual;
    switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // optional int32 store_id = 1;
      case 1: {
        if (tag == 8) {
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::int32, ::google::protobuf::internal::WireFormatLite::TYPE_INT32>(
                 input, &store_id_)));
          set_has_store_id();
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(16)) goto parse_node_id;
        break;
      }

      // optional int32 node_id = 2;
      case 2: {
        if (tag == 16) {
         parse_node_id:
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::int32, ::google::protobuf::internal::WireFormatLite::TYPE_INT32>(
                 input, &node_id_)));
          set_has_node_id();
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(24)) goto parse_range_count;
        break;
      }

      // optional int32 range_count = 3;
      case 3: {
        if (tag == 24) {
         parse_range_count:
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::int32, ::google::protobuf::internal::WireFormatLite::TYPE_INT32>(
                 input, &range_count_)));
          set_has_range_count();
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(32)) goto parse_started_at;
        break;
      }

      // optional int64 started_at = 4;
      case 4: {
        if (tag == 32) {
         parse_started_at:
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::int64, ::google::protobuf::internal::WireFormatLite::TYPE_INT64>(
                 input, &started_at_)));
          set_has_started_at();
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(40)) goto parse_updated_at;
        break;
      }

      // optional int64 updated_at = 5;
      case 5: {
        if (tag == 40) {
         parse_updated_at:
          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::int64, ::google::protobuf::internal::WireFormatLite::TYPE_INT64>(
                 input, &updated_at_)));
          set_has_updated_at();
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(50)) goto parse_stats;
        break;
      }

      // optional .cockroach.proto.MVCCStats stats = 6;
      case 6: {
        if (tag == 50) {
         parse_stats:
          DO_(::google::protobuf::internal::WireFormatLite::ReadMessageNoVirtual(
               input, mutable_stats()));
        } else {
          goto handle_unusual;
        }
        if (input->ExpectAtEnd()) goto success;
        break;
      }

      default: {
      handle_unusual:
        if (tag == 0 ||
            ::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_END_GROUP) {
          goto success;
        }
        DO_(::google::protobuf::internal::WireFormat::SkipField(
              input, tag, mutable_unknown_fields()));
        break;
      }
    }
  }
success:
  // @@protoc_insertion_point(parse_success:cockroach.proto.StoreStatus)
  return true;
failure:
  // @@protoc_insertion_point(parse_failure:cockroach.proto.StoreStatus)
  return false;
#undef DO_
}

void StoreStatus::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // @@protoc_insertion_point(serialize_start:cockroach.proto.StoreStatus)
  // optional int32 store_id = 1;
  if (has_store_id()) {
    ::google::protobuf::internal::WireFormatLite::WriteInt32(1, this->store_id(), output);
  }

  // optional int32 node_id = 2;
  if (has_node_id()) {
    ::google::protobuf::internal::WireFormatLite::WriteInt32(2, this->node_id(), output);
  }

  // optional int32 range_count = 3;
  if (has_range_count()) {
    ::google::protobuf::internal::WireFormatLite::WriteInt32(3, this->range_count(), output);
  }

  // optional int64 started_at = 4;
  if (has_started_at()) {
    ::google::protobuf::internal::WireFormatLite::WriteInt64(4, this->started_at(), output);
  }

  // optional int64 updated_at = 5;
  if (has_updated_at()) {
    ::google::protobuf::internal::WireFormatLite::WriteInt64(5, this->updated_at(), output);
  }

  // optional .cockroach.proto.MVCCStats stats = 6;
  if (has_stats()) {
    ::google::protobuf::internal::WireFormatLite::WriteMessageMaybeToArray(
      6, this->stats(), output);
  }

  if (!unknown_fields().empty()) {
    ::google::protobuf::internal::WireFormat::SerializeUnknownFields(
        unknown_fields(), output);
  }
  // @@protoc_insertion_point(serialize_end:cockroach.proto.StoreStatus)
}

::google::protobuf::uint8* StoreStatus::SerializeWithCachedSizesToArray(
    ::google::protobuf::uint8* target) const {
  // @@protoc_insertion_point(serialize_to_array_start:cockroach.proto.StoreStatus)
  // optional int32 store_id = 1;
  if (has_store_id()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteInt32ToArray(1, this->store_id(), target);
  }

  // optional int32 node_id = 2;
  if (has_node_id()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteInt32ToArray(2, this->node_id(), target);
  }

  // optional int32 range_count = 3;
  if (has_range_count()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteInt32ToArray(3, this->range_count(), target);
  }

  // optional int64 started_at = 4;
  if (has_started_at()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteInt64ToArray(4, this->started_at(), target);
  }

  // optional int64 updated_at = 5;
  if (has_updated_at()) {
    target = ::google::protobuf::internal::WireFormatLite::WriteInt64ToArray(5, this->updated_at(), target);
  }

  // optional .cockroach.proto.MVCCStats stats = 6;
  if (has_stats()) {
    target = ::google::protobuf::internal::WireFormatLite::
      WriteMessageNoVirtualToArray(
        6, this->stats(), target);
  }

  if (!unknown_fields().empty()) {
    target = ::google::protobuf::internal::WireFormat::SerializeUnknownFieldsToArray(
        unknown_fields(), target);
  }
  // @@protoc_insertion_point(serialize_to_array_end:cockroach.proto.StoreStatus)
  return target;
}

int StoreStatus::ByteSize() const {
  int total_size = 0;

  if (_has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    // optional int32 store_id = 1;
    if (has_store_id()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::Int32Size(
          this->store_id());
    }

    // optional int32 node_id = 2;
    if (has_node_id()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::Int32Size(
          this->node_id());
    }

    // optional int32 range_count = 3;
    if (has_range_count()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::Int32Size(
          this->range_count());
    }

    // optional int64 started_at = 4;
    if (has_started_at()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::Int64Size(
          this->started_at());
    }

    // optional int64 updated_at = 5;
    if (has_updated_at()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::Int64Size(
          this->updated_at());
    }

    // optional .cockroach.proto.MVCCStats stats = 6;
    if (has_stats()) {
      total_size += 1 +
        ::google::protobuf::internal::WireFormatLite::MessageSizeNoVirtual(
          this->stats());
    }

  }
  if (!unknown_fields().empty()) {
    total_size +=
      ::google::protobuf::internal::WireFormat::ComputeUnknownFieldsSize(
        unknown_fields());
  }
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = total_size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
  return total_size;
}

void StoreStatus::MergeFrom(const ::google::protobuf::Message& from) {
  GOOGLE_CHECK_NE(&from, this);
  const StoreStatus* source =
    ::google::protobuf::internal::dynamic_cast_if_available<const StoreStatus*>(
      &from);
  if (source == NULL) {
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
    MergeFrom(*source);
  }
}

void StoreStatus::MergeFrom(const StoreStatus& from) {
  GOOGLE_CHECK_NE(&from, this);
  if (from._has_bits_[0 / 32] & (0xffu << (0 % 32))) {
    if (from.has_store_id()) {
      set_store_id(from.store_id());
    }
    if (from.has_node_id()) {
      set_node_id(from.node_id());
    }
    if (from.has_range_count()) {
      set_range_count(from.range_count());
    }
    if (from.has_started_at()) {
      set_started_at(from.started_at());
    }
    if (from.has_updated_at()) {
      set_updated_at(from.updated_at());
    }
    if (from.has_stats()) {
      mutable_stats()->::cockroach::proto::MVCCStats::MergeFrom(from.stats());
    }
  }
  mutable_unknown_fields()->MergeFrom(from.unknown_fields());
}

void StoreStatus::CopyFrom(const ::google::protobuf::Message& from) {
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void StoreStatus::CopyFrom(const StoreStatus& from) {
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool StoreStatus::IsInitialized() const {

  return true;
}

void StoreStatus::Swap(StoreStatus* other) {
  if (other != this) {
    std::swap(store_id_, other->store_id_);
    std::swap(node_id_, other->node_id_);
    std::swap(range_count_, other->range_count_);
    std::swap(started_at_, other->started_at_);
    std::swap(updated_at_, other->updated_at_);
    std::swap(stats_, other->stats_);
    std::swap(_has_bits_[0], other->_has_bits_[0]);
    _unknown_fields_.Swap(&other->_unknown_fields_);
    std::swap(_cached_size_, other->_cached_size_);
  }
}

::google::protobuf::Metadata StoreStatus::GetMetadata() const {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::Metadata metadata;
  metadata.descriptor = StoreStatus_descriptor_;
  metadata.reflection = StoreStatus_reflection_;
  return metadata;
}


// @@protoc_insertion_point(namespace_scope)

}  // namespace proto
}  // namespace cockroach

// @@protoc_insertion_point(global_scope)