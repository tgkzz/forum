package post

import (
	"database/sql"
	"forum/internal/model"
)

type PostRepo struct {
	DB *sql.DB
}

type Post interface {
	CreatePost(model.Post) error
	GetAllPost() ([]model.Post, error)
	GetPostById(id int) (model.Post, error)
	CreateComment(comment model.Comment) error
	AddGradeToPost(grade model.Grade) error
	AddGradeToComment(grade model.Grade) error
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{
		DB: db,
	}
}

func (p *PostRepo) AddGradeToPost(grade model.Grade) error {
	var existingGrade int
	var existingGradeId int

	err := p.DB.QueryRow(`
        SELECT Id, GradeValue FROM Grade 
        WHERE UserId = $1 AND PostId = $2`,
		grade.UserId,
		grade.PostId,
	).Scan(&existingGradeId, &existingGrade)

	if err == sql.ErrNoRows {
		_, err := p.DB.Exec(`
            INSERT INTO Grade (UserId, PostId, GradeValue) 
            VALUES ($1, $2, $3)`,
			grade.UserId,
			grade.PostId,
			grade.GradeValue,
		)
		return err
	} else if err != nil {
		return err
	} else {
		if existingGrade == grade.GradeValue {
			_, err := p.DB.Exec("DELETE FROM Grade WHERE Id = $1", existingGradeId)
			return err
		} else {
			_, err := p.DB.Exec(`
                UPDATE Grade 
                SET GradeValue = $1 
                WHERE Id = $2`,
				grade.GradeValue,
				existingGradeId,
			)
			return err
		}
	}
}

func (p *PostRepo) AddGradeToComment(grade model.Grade) error {
	var existingGrade int
	var existingGradeId int

	err := p.DB.QueryRow(`
        SELECT Id, GradeValue FROM Grade 
        WHERE UserId = $1 AND CommentId = $2`,
		grade.UserId,
		grade.CommentId,
	).Scan(&existingGradeId, &existingGrade)

	if err == sql.ErrNoRows {
		_, err := p.DB.Exec(`
            INSERT INTO Grade (UserId, CommentId, GradeValue) 
            VALUES ($1, $2, $3)`,
			grade.UserId,
			grade.CommentId,
			grade.GradeValue,
		)
		return err
	} else if err != nil {
		return err
	} else {
		if existingGrade == grade.GradeValue {
			_, err := p.DB.Exec("DELETE FROM Grade WHERE Id = $1", existingGradeId)
			return err
		} else {
			_, err := p.DB.Exec(`
                UPDATE Grade 
                SET GradeValue = $1 
                WHERE Id = $2`,
				grade.GradeValue,
				existingGradeId,
			)
			return err
		}
	}
}

func (p *PostRepo) CreatePost(post model.Post) error {
	query := "INSERT INTO Post (Name, Text, CreationTime, UserId, ImagePath) VALUES ($1, $2, $3, $4, $5) RETURNING Id"
	var postId int
	err := p.DB.QueryRow(query, post.Name, post.Text, post.CreationTime, post.UserId, post.PhotoPath).Scan(&postId)
	if err != nil {
		return err
	}

	for _, categoryId := range post.CategoryId {
		query = "INSERT INTO PostCategory (PostId, CategoryId) VALUES ($1, $2)"
		_, err = p.DB.Exec(query, postId, categoryId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostRepo) GetAllPost() ([]model.Post, error) {
	// query := "SELECT Id, Name, Text, UserId, CategoryId FROM Post"
	// query := "SELECT Id, Name, Text, UserId FROM Post"
	query := `
    SELECT p.Id, p.Name, p.Text, p.UserId, u.Username, p.CreationTime, p.ImagePath 
    FROM Post p 
    JOIN Users u ON p.UserId = u.Id
    `

	result := []model.Post{}

	rows, err := p.DB.Query(query)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.Id, &post.Name, &post.Text, &post.UserId, &post.Username, &post.CreationTime, &post.PhotoPath); err != nil {
			return result, err
		}

		result = append(result, post)
	}

	return result, rows.Err()
}

func (p *PostRepo) GetPostById(id int) (model.Post, error) {
	var result model.Post

	postQuery := `
    SELECT p.Id, p.Name, p.Text, p.UserId, p.CreationTime, p.ImagePath
    FROM Post p
    WHERE p.Id = $1
    `
	row := p.DB.QueryRow(postQuery, id)
	if err := row.Scan(&result.Id, &result.Name, &result.Text, &result.UserId, &result.CreationTime, &result.PhotoPath); err != nil {
		return result, model.ErrInvalidId
	}

	categoriesQuery := `
    SELECT CategoryId
    FROM PostCategory
    WHERE PostId = $1
    `
	rows, err := p.DB.Query(categoriesQuery, id)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryId int
		if err := rows.Scan(&categoryId); err != nil {
			return result, err
		}
		result.CategoryId = append(result.CategoryId, categoryId)
	}

	if err = rows.Err(); err != nil {
		return result, err
	}

	gradeQuery := `
    SELECT 
        COALESCE(SUM(CASE WHEN GradeValue = 1 THEN 1 ELSE 0 END), 0) as Likes,
        COALESCE(SUM(CASE WHEN GradeValue = -1 THEN 1 ELSE 0 END), 0) as Dislikes
    FROM Grade WHERE PostId = $1
    `
	err = p.DB.QueryRow(gradeQuery, id).Scan(&result.Likes, &result.Dislikes)
	if err != nil {
		return result, err
	}

	commentsQuery := `
	SELECT 
		c.Id, c.Text, c.PostId, c.UserId, u.Username,
		COALESCE(SUM(CASE WHEN g.GradeValue = 1 THEN 1 ELSE 0 END), 0) as Likes,
		COALESCE(SUM(CASE WHEN g.GradeValue = -1 THEN 1 ELSE 0 END), 0) as Dislikes
	FROM Comment c
	JOIN Users u ON c.UserId = u.Id
	LEFT JOIN Grade g ON c.Id = g.CommentId
	WHERE c.PostId = $1
	GROUP BY c.Id, u.Username, c.Text, c.PostId, c.UserId
	`
	rows, err = p.DB.Query(commentsQuery, id)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.Id, &comment.Text, &comment.PostId, &comment.UserId, &comment.Username, &comment.Likes, &comment.Dislikes); err != nil {
			return result, err
		}
		result.Comment = append(result.Comment, comment)
	}

	if err = rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func (p *PostRepo) CreateComment(comment model.Comment) error {
	query := "INSERT INTO Comment (UserId, PostId, Text) VALUES ($1, $2, $3)"

	_, err := p.DB.Exec(query, comment.UserId, comment.PostId, comment.Text)

	return err
}
