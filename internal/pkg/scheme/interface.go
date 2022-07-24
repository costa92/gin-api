package scheme

type ObjectKind interface {
	SetGroupVersionKind(kind GroupVersionKind)
	GroupVersionKind() GroupVersionKind
}

var EmptyObjectKind = emptyObjectKind{}

type emptyObjectKind struct{}

func (emptyObjectKind) SetGroupVersionKind(gvk GroupVersionKind) {}

// GroupVersionKind implements the ObjectKind interface.
func (emptyObjectKind) GroupVersionKind() GroupVersionKind { return GroupVersionKind{} }
