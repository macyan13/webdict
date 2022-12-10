package tag

type Repository interface {
	Create(tag Tag) error
	Update(tag Tag) error
	Get(id, authorId string) (*Tag, error)
	Delete(id, authorId string) error
	AllExist(ids []string, authorId string) (bool, error)
}
