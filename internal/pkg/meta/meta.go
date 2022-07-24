package meta

import (
	"time"

	"github.com/costa92/go-web/internal/pkg/scheme"
)

type ObjectMetaAccessor interface {
	GetObjectMeta() Object
}

type Object interface {
	GetID() uint64
	SetID(id uint64)
	GetName() string
	SetName(name string)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(Updated time.Time)
}

type ListInterface interface {
	GetTotalCount() int64
	SetTotalCount(count int64)
}

type Type interface {
	GetAPIVersion() string
	SetAPIVersion(version string)
	GetKind() string
	SetKind(kind string)
}

var _ ListInterface = &ListMeta{}

func (meta *ListMeta) GetTotalCount() int64      { return meta.TotalCount }
func (meta *ListMeta) SetTotalCount(count int64) { meta.TotalCount = count }

var _ Type = &TypeMeta{}

func (meta *TypeMeta) GetObjectKind() scheme.ObjectKind { return meta }

func (meta *TypeMeta) SetGroupVersionKind(gvk scheme.GroupVersionKind) {
	meta.APIVersion, meta.Kind = gvk.ToAPIVersionAndKind()
}

func (meta *TypeMeta) GroupVersionKind() scheme.GroupVersionKind {
	return scheme.FromAPIVersionAndKind(meta.APIVersion, meta.Kind)
}

func (meta *TypeMeta) GetAPIVersion() string        { return meta.APIVersion }
func (meta *TypeMeta) SetAPIVersion(version string) { meta.APIVersion = version }
func (meta *TypeMeta) GetKind() string              { return meta.Kind }
func (meta *TypeMeta) SetKind(kind string)          { meta.Kind = kind }

func (meta *ListMeta) GetListMeta() ListInterface { return meta }

func (obj *ObjectMeta) GetObjectMeta() Object { return obj }

var _ Object = &ObjectMeta{}

func (obj *ObjectMeta) GetID() uint64                    { return obj.ID }
func (obj *ObjectMeta) SetID(id uint64)                  { obj.ID = id }
func (obj *ObjectMeta) GetName() string                  { return obj.Name }
func (obj *ObjectMeta) SetName(name string)              { obj.Name = name }
func (obj *ObjectMeta) GetCreatedAt() time.Time          { return obj.CreatedAt }
func (obj *ObjectMeta) SetCreatedAt(createdAt time.Time) { obj.CreatedAt = createdAt }
func (obj *ObjectMeta) GetUpdatedAt() time.Time          { return obj.UpdatedAt }
func (obj *ObjectMeta) SetUpdatedAt(updatedAt time.Time) { obj.UpdatedAt = updatedAt }
