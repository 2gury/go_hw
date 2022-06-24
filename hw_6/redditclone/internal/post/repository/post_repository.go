package repository

import (
	"context"
	"database/sql"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/post"
	"log"
)

type PostRepository struct {
	dbConn *sql.DB
}

func NewPostRepository(conn *sql.DB) post.PostRepository {
	return &PostRepository{
		dbConn: conn,
	}
}

func (r *PostRepository) InsertPost(post *models.Post) (uint64, error) {
	tx, err := r.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	var lastID int64
	err = tx.QueryRow(`INSERT INTO posts (user_id, category, created, score, text, title, type, upvote_percentage, url, views) 
					   VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		post.UserID, post.Category, post.CreatedAt, post.Score, post.Text, post.Title, post.Type, post.UpvotePercentage, post.URL, post.Views).Scan(&lastID)
	if err != nil {
		if rollBackError := tx.Rollback(); rollBackError != nil {
			log.Fatal(rollBackError.Error())
		}
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(lastID), nil
}

func (r *PostRepository) DeletePostByID(postID uint64) error {
	tx, err := r.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`DELETE FROM posts
		 WHERE id=$1 `, postID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatal(rollbackErr)
		}
		return nil
	}
	_, err = tx.Exec(
		`DELETE FROM comments cmt
		 WHERE cmt.post_id=$1 `, postID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatal(rollbackErr)
		}
		return nil
	}
	_, err = tx.Exec(
		`DELETE FROM votes vts
		 WHERE vts.post_id=$1 `, postID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatal(rollbackErr)
		}
		return nil
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) SelectAuthorPost(postID uint64) (*models.User, error) {
	author := &models.User{}

	err := r.dbConn.QueryRow(
		`SELECT usr.id, usr.username
		FROM posts pst
		INNER JOIN users usr
		ON pst.user_id = usr.id
		WHERE pst.id=$1`, postID).
		Scan(&author.ID, &author.Username)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (r *PostRepository) SelectAuthorComment(commentID uint64) (*models.User, error) {
	author := &models.User{}

	err := r.dbConn.QueryRow(
		`SELECT usr.id, usr.username
		FROM comments cmt
		INNER JOIN users usr
		ON cmt.user_id = usr.id
		WHERE cmt.id=$1`, commentID).
		Scan(&author.ID, &author.Username)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (r *PostRepository) DeleteVoteFromPostByUserID(userID uint64, postID uint64) error {
	tx, err := r.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`DELETE FROM votes vts
		 WHERE vts.user_id=$1 AND vts.post_id=$2`, userID, postID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatal(rollbackErr)
		}
		return nil
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) InsertVote(vote *models.Vote) (uint64, error) {
	tx, err := r.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return 0, err
		}
		var lastID int64
		err = tx.QueryRow(`INSERT INTO votes(user_id, post_id, vote)
						   VALUES($1, $2, $3) RETURNING id`,
			vote.UserID, vote.PostID, vote.Vote).Scan(&lastID)
		if err != nil {
			if rollBackError := tx.Rollback(); rollBackError != nil {
				log.Fatal(rollBackError.Error())
			}
			return 0, err
		}
		if err := tx.Commit(); err != nil {
			return 0, err
		}
	
		return uint64(lastID), nil
}

func (r *PostRepository) SelectAllPosts() ([]*models.Post, error) {
	rows, err := r.dbConn.Query(
		`SELECT usr.id, usr.username, pst.id, pst.user_id, pst.category, pst.created, pst.score, pst.text, pst.title, pst.type, pst.upvote_percentage, pst.url, pst.views 
		FROM posts pst
		INNER JOIN users usr
		ON pst.user_id = usr.id
		ORDER BY pst.score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}
	for rows.Next() {
		post := &models.Post{}
		author := &models.User{}
		err := rows.Scan(&author.ID, &author.Username, &post.ID, &post.UserID, &post.Category, &post.CreatedAt, &post.Score, &post.Text, &post.Title, &post.Type, &post.UpvotePercentage, &post.URL, &post.Views)
		if err != nil {
			return nil, err
		}
		post.Author = author
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) SelectAllComments() ([]*models.Comment, error) {
	rows, err := r.dbConn.Query(
		`SELECT usr.id, usr.username, cmt.post_id, cmt.id, cmt.body, cmt.created 
		 FROM comments cmt
		 INNER JOIN users usr
		 ON cmt.user_id = usr.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{}
		author := &models.User{}
		err := rows.Scan(&author.ID, &author.Username, &comment.PostID, &comment.ID, &comment.Body, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comment.Author = author
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *PostRepository) SelectAllVotes() ([]*models.Vote, error) {
	rows, err := r.dbConn.Query(
		`SELECT vts.post_id, vts.user_id, vts.vote
		 FROM votes vts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votes := []*models.Vote{}
	for rows.Next() {
		vote := &models.Vote{}
		err := rows.Scan(&vote.PostID, &vote.User, &vote.Vote)
		if err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return votes, nil
}

func (r *PostRepository) SelectPostsByCategory(categoryName string) ([]*models.Post, error) {
	rows, err := r.dbConn.Query(
		`SELECT usr.id, usr.username, pst.id, pst.user_id, pst.category, pst.created, pst.score, pst.text, pst.title, pst.type, pst.upvote_percentage, pst.url, pst.views 
		FROM posts pst
		INNER JOIN users usr
		ON pst.user_id = usr.id
		WHERE pst.category=$1
		ORDER BY pst.score DESC`, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}
	for rows.Next() {
		post := &models.Post{}
		author := &models.User{}
		err := rows.Scan(&author.ID, &author.Username, &post.ID, &post.UserID, &post.Category, &post.CreatedAt, &post.Score, &post.Text, &post.Title, &post.Type, &post.UpvotePercentage, &post.URL, &post.Views)
		if err != nil {
			return nil, err
		}
		post.Author = author
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) SelectCommentsByCategory(categoryName string) ([]*models.Comment, error) {
	rows, err := r.dbConn.Query(
		`SELECT usr.id, usr.username, cmt.post_id, cmt.id, cmt.body, cmt.created 
		 FROM comments cmt
		 INNER JOIN users usr
		 ON cmt.user_id = usr.id
		 INNER JOIN posts pst
		 ON cmt.post_id = pst.id
		 WHERE pst.category=$1`, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{}
		author := &models.User{}
		err := rows.Scan(&author.ID, &author.Username, &comment.PostID, &comment.ID, &comment.Body, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comment.Author = author
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *PostRepository) SelectVotesByCategory(categoryName string) ([]*models.Vote, error) {
	rows, err := r.dbConn.Query(
		`SELECT vts.post_id, vts.user_id, vts.vote
		 FROM votes vts
		 INNER JOIN posts pst
		 ON vts.post_id = pst.id
		 WHERE pst.category=$1`, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votes := []*models.Vote{}
	for rows.Next() {
		vote := &models.Vote{}
		err := rows.Scan(&vote.PostID, &vote.User, &vote.Vote)
		if err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return votes, nil
}

func (r *PostRepository) SelectPostsByUsername(userLogin string) ([]*models.Post, error) {
	rows, err := r.dbConn.Query(
		`SELECT usr.id, usr.username, pst.id, pst.user_id, pst.category, pst.created, pst.score, pst.text, pst.title, pst.type, pst.upvote_percentage, pst.url, pst.views 
		FROM posts pst
		INNER JOIN users usr
		ON pst.user_id = usr.id
		WHERE usr.username=$1`, userLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}
	for rows.Next() {
		post := &models.Post{}
		author := &models.User{}
		err := rows.Scan(&author.ID, &author.Username, &post.ID, &post.UserID, &post.Category, &post.CreatedAt, &post.Score, &post.Text, &post.Title, &post.Type, &post.UpvotePercentage, &post.URL, &post.Views)
		if err != nil {
			return nil, err
		}
		post.Author = author
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) SelectCommentsByUsername(userLogin string) ([]*models.Comment, error) {
	rows, err := r.dbConn.Query(
		`SELECT usr.id, usr.username, cmt.post_id, cmt.id, cmt.body, cmt.created 
		 FROM comments cmt
		 INNER JOIN posts pst
		 ON cmt.post_id = pst.id
		 INNER JOIN users usr
		 ON pst.user_id = usr.id
		 WHERE usr.username=$1`, userLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{}
		author := &models.User{}
		err := rows.Scan(&author.ID, &author.Username, &comment.PostID, &comment.ID, &comment.Body, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comment.Author = author
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *PostRepository) SelectVotesByUsername(userLogin string) ([]*models.Vote, error) {
	rows, err := r.dbConn.Query(
		`SELECT vts.post_id, vts.user_id, vts.vote
		 FROM votes vts
		 INNER JOIN posts pst
		 ON vts.post_id = pst.id
		 INNER JOIN users usr
		 ON pst.user_id = usr.id
		 WHERE usr.username=$1`, userLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votes := []*models.Vote{}
	for rows.Next() {
		vote := &models.Vote{}
		err := rows.Scan(&vote.PostID, &vote.User, &vote.Vote)
		if err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return votes, nil
}

func (r *PostRepository) SelectPostByID(postID uint64) (*models.Post, error) {
	post := &models.Post{}
	author := &models.User{}

	err := r.dbConn.QueryRow(
		`SELECT usr.id, usr.username, pst.id, pst.user_id, pst.category, pst.created, pst.score, pst.text, pst.title, pst.type, pst.upvote_percentage, pst.url, pst.views 
		FROM posts pst
		INNER JOIN users usr
		ON pst.user_id = usr.id
        WHERE pst.id=$1`, postID).
		Scan(&author.ID, &author.Username, &post.ID, &post.UserID, &post.Category, &post.CreatedAt, &post.Score, &post.Text, &post.Title, &post.Type, &post.UpvotePercentage, &post.URL, &post.Views)
	if err != nil {
		return nil, err
	}
	post.Author = author

	return post, nil
}

func (r *PostRepository) SelectCommentsByPostID(postID uint64) ([]*models.Comment, error) {
	rows, err := r.dbConn.Query(
		`SELECT usr.id, usr.username, cmt.id, cmt.body, cmt.created 
		 FROM comments cmt
		 INNER JOIN users usr
		 ON cmt.user_id = usr.id
		 WHERE cmt.post_id=$1`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{}
		author := &models.User{}
		err := rows.Scan(&author.ID, &author.Username, &comment.ID, &comment.Body, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comment.Author = author
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *PostRepository) SelectVotesByPostID(postID uint64) ([]*models.Vote, error) {
	rows, err := r.dbConn.Query(
		`SELECT vts.user_id, vts.vote
		FROM votes vts
		WHERE vts.post_id=$1`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votes := []*models.Vote{}
	for rows.Next() {
		vote := &models.Vote{}
		err := rows.Scan(&vote.User, &vote.Vote)
		if err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return votes, nil
}

func (r *PostRepository) InsertComment(comment *models.Comment) (uint64, error) {
	tx, err := r.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	var lastID int64
	err = tx.QueryRow(`INSERT INTO comments(user_id, post_id, body, created) 
                       VALUES($1, $2, $3, $4) RETURNING id`,
		comment.UserID, comment.PostID, comment.Body, comment.CreatedAt).Scan(&lastID)
	if err != nil {
		if rollBackError := tx.Rollback(); rollBackError != nil {
			log.Fatal(rollBackError.Error())
		}
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(lastID), nil
}

func (r *PostRepository) DeleteCommentByID(commentID uint64) (error) {
	tx, err := r.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`DELETE FROM comments cmt
		 WHERE cmt.id=$1 `, commentID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatal(rollbackErr)
		}
		return nil
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
