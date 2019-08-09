// Package datamanager for db persistence
package datamanager

// PostMessage post message for reddit post
type PostMessage struct {
	URL  string
	Text string
}

// SavePostMessage save post
func SavePostMessage(postMessage *PostMessage) error {
	_, err := db.Exec(`INSERT INTO posts (url, text) VALUES ($1, $2) RETURNING id`, postMessage.URL, postMessage.Text)
	if err != nil {
		return err
	}
	return nil
}

// IsPostExist check if post exists via url
func IsPostExist(url string) (bool, error) {
	var exists bool
	row := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM posts WHERE url=$1)`, url)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists != false {
		return true, nil
	}
	return false, nil
}

// AllPosts get all the saved posts
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
