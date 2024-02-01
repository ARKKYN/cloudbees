// post_dao.go
package dao

import (
	e "cloudbees/errors"
	m "cloudbees/models"
	"sync"
)

type PostDAO struct {
	posts map[uint64]*m.Post
	mu    sync.Mutex
}

var instance *PostDAO
var once sync.Once

func NewPostDAO() *PostDAO {
	once.Do(func() {
		instance = &PostDAO{
			posts: make(map[uint64]*m.Post),
			mu:    sync.Mutex{},
		}
	})
	return instance
}

func (dao *PostDAO) Create(post *m.Post) error {

	dao.mu.Lock()
	defer dao.mu.Unlock()
	print(dao.posts[post.PostId])
	dao.posts[post.PostId] = post
	return nil
}

func (dao *PostDAO) Read(id uint64) (*m.Post, error) {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	post, exists := dao.posts[id]
	if !exists {
		return nil, e.EnitityNotFoundError
	}
	return post, nil
}

func (dao *PostDAO) Update(post *m.Post) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	_, exists := dao.posts[post.PostId]
	if !exists {
		return e.EnitityNotFoundError
	}
	dao.posts[post.PostId] = post
	return nil
}

func (dao *PostDAO) Delete(id uint64) error {
	dao.mu.Lock()
	defer dao.mu.Unlock()
	_, exists := dao.posts[id]
	if !exists {
		return e.EnitityNotFoundError
	}
	delete(dao.posts, id)
	return nil
}
