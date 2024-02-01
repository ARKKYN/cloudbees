package dao

type DAO interface {
	Create(entity interface{}) error
	Read(id int32) (interface{}, error)
	Update(id int32, updatedEntity interface{}) error
	Delete(id int32) error
}
