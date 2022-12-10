package translation

// Repository defines domain translation repository methods
type Repository interface {
	Create(translation Translation) error
	Update(translation Translation) error
	Get(id, authorId string) (*Translation, error)
	Delete(id, authorId string) error
}
