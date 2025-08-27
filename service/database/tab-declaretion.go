package database

var user_table = `CREATE TABLE IF NOT EXISTS User
				(
					userID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					username TEXT NOT NULL UNIQUE CHECK (LENGTH(username) >= 3 AND LENGTH(username) <= 13)
				);`

var post_table = `CREATE TABLE IF NOT EXISTS Post
				(
					PhotoID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					userID INTEGER NOT NULL,
					image BLOB NOT NULL,
					caption TEXT NOT NULL,
					created_at TEXT NOT NULL,
					CONSTRAINT fk_post
						FOREIGN KEY(userID) REFERENCES user(userID) 
						ON DELETE CASCADE
				);`

// like, comment, follow, ban
var like_table = `CREATE TABLE IF NOT EXISTS Like
				(
					userID INTEGER NOT NULL,
					PhotoID INTEGER NOT NULL,
					PRIMARY KEY(userID, PhotoID),
					CONSTRAINT fk_like
						FOREIGN KEY(userID) REFERENCES user(userID) ON DELETE CASCADE,
						FOREIGN KEY(PhotoID) REFERENCES post(PhotoID) ON DELETE CASCADE
				);`

var comment_table = `CREATE TABLE IF NOT EXISTS Comment
				(
					commentID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					userID INTEGER NOT NULL,
					textComment TEXT NOT NULL,
					PhotoID INTEGER NOT NULL,
					created_at TEXT NOT NULL,
					CONSTRAINT fk_comment
						FOREIGN KEY(userID) REFERENCES user(userID) ON DELETE CASCADE,
						FOREIGN KEY(PhotoID) REFERENCES post(PhotoID) ON DELETE CASCADE
				);`

var follow_table = `CREATE TABLE IF NOT EXISTS Follow
				(
					followerID INTEGER NOT NULL,
					followedID INTEGER NOT NULL,
					PRIMARY KEY(followerID, followedID),
					CONSTRAINT fk_follow
						FOREIGN KEY(followerID) REFERENCES user(userID) ON DELETE CASCADE,
						FOREIGN KEY(followedID) REFERENCES user(userID) ON DELETE CASCADE
				);`

var ban_table = `CREATE TABLE IF NOT EXISTS Ban
				(
					bannerID INTEGER NOT NULL,
					bannedID INTEGER NOT NULL,
					PRIMARY KEY(bannerID, bannedID),
					CONSTRAINT fk_ban
						FOREIGN KEY(bannerID) REFERENCES user(userID) ON DELETE CASCADE,
						FOREIGN KEY(bannedID) REFERENCES user(userID) ON DELETE CASCADE
				);`
