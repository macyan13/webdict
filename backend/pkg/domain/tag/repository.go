package tag

type Repository interface {
	Save(tag Tag) error
	GetById(id string) *Tag
	GetByIds(ids []string) []*Tag // May be remove it
	Get() []Tag
	Delete(tag Tag) error
	AllExist(ids []string, AuthorId string) bool
}
