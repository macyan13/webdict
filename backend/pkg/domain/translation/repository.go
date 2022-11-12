package translation

type Repository interface {
	Save(translation Translation) error
	GetById(id string) *Translation
	Get() []Translation
	Delete(translation Translation) error
}
