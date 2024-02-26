package helpers

import "database/sql"

func GetTotalLikesForPost(dbConnection *sql.DB, postID string) (int, error) {
	var totalLikes int
	err := dbConnection.QueryRow("SELECT COUNT(*) FROM postLikes WHERE post_id = ?", postID).Scan(&totalLikes)
	if err != nil {
		return 0, err
	}
	return totalLikes, nil
}
