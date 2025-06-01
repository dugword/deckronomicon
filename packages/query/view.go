package query

type View interface {
	Get(id string) (Object, error)
	GetAll() []Object
	Name() string
	Size() int
}

type view struct {
	name    string
	objects []Object
}

func NewView[T Object](name string, objects []T) View {
	v := view{
		name: name,
	}
	for _, obj := range objects {
		v.objects = append(v.objects, obj)
	}
	return &v
}

func (v *view) Get(id string) (Object, error) {
	for _, obj := range v.objects {
		if obj.ID() == id {
			return obj, nil
		}
	}
	return nil, ErrNotFound
}

func (v *view) GetAll() []Object {
	return v.objects
}

func (v *view) Name() string {
	return v.name
}

func (v *view) Size() int {
	return len(v.objects)
}
