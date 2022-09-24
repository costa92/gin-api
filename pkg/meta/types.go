package meta

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Extend map[string]interface{}

func (ext Extend) String() string {
	data, _ := json.Marshal(ext)
	return string(data)
}

func (ext Extend) Merge(extendShadow string) Extend {
	var extend Extend
	_ = json.Unmarshal([]byte(extendShadow), &extend)
	for k, v := range extend {
		if _, ok := ext[k]; !ok {
			ext[k] = v
		}
	}
	return ext
}

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}

type ObjectMeta struct {
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`

	InstanceID string `json:"instanceID,omitempty" gorm:"unique;column:instanceID;type:varchar(32);not null"`

	Name string `json:"name,omitempty" gorm:"column:name;type:varchar(64);not null" validate:"name"`

	// Extend store the fields that need to be added, but do not want to add a new table column, will not be stored in db.
	Extend Extend `json:"extend,omitempty" gorm:"-" validate:"omitempty"`

	// ExtendShadow is the shadow of Extend. DO NOT modify directly.
	ExtendShadow string `json:"-" gorm:"column:extendShadow" validate:"omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`

	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (obj *ObjectMeta) BeforeCreate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()
	return nil
}

func (obj *ObjectMeta) BeforeUpdate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()
	return nil
}

func (obj *ObjectMeta) AfterFind(tx *gorm.DB) error {
	if obj.ExtendShadow != "" {
		if err := json.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
			return err
		}
	}
	return nil
}

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	TypeMeta `json:",inline"`

	// LabelSelector is used to find matching REST resources.
	LabelSelector string `json:"labelSelector,omitempty" form:"labelSelector"`

	// FieldSelector restricts the list of returned objects by their fields. Defaults to everything.
	FieldSelector string `json:"fieldSelector,omitempty" form:"fieldSelector"`

	// TimeoutSeconds specifies the seconds of ClientIP type session sticky time.
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	// Offset specify the number of records to skip before starting to return the records.
	Offset *int64 `json:"offset,omitempty" form:"offset"`

	// Limit specify the number of records to be retrieved.
	Limit *int64 `json:"limit,omitempty" form:"limit"`
}

// ExportOptions is the query options to the standard REST get call.
// Deprecated. Planned for removal in 1.18.
type ExportOptions struct {
	TypeMeta `json:",inline"`

	// Should this value be exported.  Export strips fields that a user can not specify.
	// Deprecated. Planned for removal in 1.18.
	Export bool `json:"export"`
	// Should the export be exact.  Exact export maintains cluster-specific fields like 'Namespace'.
	// Deprecated. Planned for removal in 1.18.
	Exact bool `json:"exact"`
}

// GetOptions is the standard query options to the standard REST get call.
type GetOptions struct {
	TypeMeta `json:",inline"`
}

// DeleteOptions may be provided when deleting an API object.
type DeleteOptions struct {
	TypeMeta `json:",inline"`

	// +optional
	Unscoped bool `json:"unscoped"`
}

// CreateOptions may be provided when creating an API object.
type CreateOptions struct {
	TypeMeta `json:",inline"`

	// When present, indicates that modifications should not be
	// persisted. An invalid or unrecognized dryRun directive will
	// result in an error response and no further processing of the
	// request. Valid values are:
	// - All: all dry run stages will be processed
	// +optional
	DryRun []string `json:"dryRun,omitempty"`
}

// PatchOptions may be provided when patching an API object.
// PatchOptions is meant to be a superset of UpdateOptions.
type PatchOptions struct {
	TypeMeta `json:",inline"`

	// When present, indicates that modifications should not be
	// persisted. An invalid or unrecognized dryRun directive will
	// result in an error response and no further processing of the
	// request. Valid values are:
	// - All: all dry run stages will be processed
	// +optional
	DryRun []string `json:"dryRun,omitempty"`

	// Force is going to "force" Apply requests. It means user will
	// re-acquire conflicting fields owned by other people. Force
	// flag must be unset for non-apply patch requests.
	// +optional
	Force bool `json:"force,omitempty"`
}

// UpdateOptions may be provided when updating an API object.
// All fields in UpdateOptions should also be present in PatchOptions.
type UpdateOptions struct {
	TypeMeta `json:",inline"`

	// When present, indicates that modifications should not be
	// persisted. An invalid or unrecognized dryRun directive will
	// result in an error response and no further processing of the
	// request. Valid values are:
	// - All: all dry run stages will be processed
	// +optional
	DryRun []string `json:"dryRun,omitempty"`
}

// AuthorizeOptions may be provided when authorize an API object.
type AuthorizeOptions struct {
	TypeMeta `json:",inline"`
}

// TableOptions are used when a Table is requested by the caller.
type TableOptions struct {
	TypeMeta `json:",inline"`

	// NoHeaders is only exposed for internal callers. It is not included in our OpenAPI definitions
	// and may be removed as a field in a future release.
	NoHeaders bool `json:"-"`
}
