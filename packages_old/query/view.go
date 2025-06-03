package query

type View interface {
	Contains(predicate Predicate) bool
	Count(predicate Predicate) int
	Find(predicate Predicate) (Object, bool)
	FindAll(predicate Predicate) []Object
	Get(id string) (Object, bool)
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

func (v *view) Contains(predicate Predicate) bool {
	return Contains(v.objects, predicate)
}

func (v *view) Count(predicate Predicate) int {
	return Count(v.objects, predicate)
}

func (v *view) Find(predicate Predicate) (Object, bool) {
	return Find(v.objects, predicate)
}

func (v *view) FindAll(predicate Predicate) []Object {
	return FindAll(v.objects, predicate)
}

func (v *view) Get(id string) (Object, bool) {
	return Get(v.objects, id)
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
