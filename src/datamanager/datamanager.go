// Package datamanager for db persistence
package datamanager

// PostMessage post message for reddit post
type PostMessage struct {
	URL  string
	Text string
}

// SavePostMessage save post
//TODO: Add logging
func SavePostMessage(postMessage PostMessage) error {
	var id int
	err := db.QueryRow(`INSERT INTO posts (url, text) VALUES ($1, $2) RETURNING id`, postMessage.URL, postMessage.Text).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

// AllPosts get all the saved posts
//TODO: Add logging
func AllPosts() ([]*PostMessage, error) {
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*PostMessage, 0)
	for rows.Next() {
		post := new(PostMessage)
		err := rows.Scan(&post.URL, &post.Text)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
