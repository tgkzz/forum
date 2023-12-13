package filter

import (
	"database/sql"
	"fmt"
	"forum/internal/model"
)

type FilterRepo struct {
	DB *sql.DB
}

type Filter interface {
	GetUserPostsById(id int) ([]model.Post, error)
	GetPostsByCategory(categories []int) ([]model.Post, error)
}

func NewFilterRepo(db *sql.DB) *FilterRepo {
	return &FilterRepo{
		DB: db,
	}
}

func (f *FilterRepo) GetUserPostsById(id int) ([]model.Post, error) {
	query := `
	SELECT p.Id, p.Name, p.Text, p.CreationTime, p.UserId, u.Username
	FROM Post p
	JOIN Users u ON p.UserId = u.Id
	WHERE p.UserId = $1
	`

	result := []model.Post{}

	rows, err := f.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.Id, &post.Name, &post.Text, &post.CreationTime, &post.UserId, &post.Username); err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (f *FilterRepo) GetPostsByCategory(categories []int) ([]model.Post, error) {
	var inParams string
	for i := range categories {
		if i > 0 {
			inParams += ", "
		}
		inParams += "?"
	}

	query := fmt.Sprintf(`
	SELECT p.Id, p.Name, p.Text, p.CreationTime, p.UserId, u.Username
	FROM Post p
	JOIN Users u ON p.UserId = u.Id
	JOIN PostCategory pc ON p.Id = pc.PostId
	WHERE pc.CategoryId IN (%s)
	`, inParams)

	args := make([]interface{}, len(categories))
	for i, v := range categories {
		args[i] = v
	}

	result := []model.Post{}
	rows, err := f.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.Id, &post.Name, &post.Text, &post.CreationTime, &post.UserId, &post.Username); err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil

}
