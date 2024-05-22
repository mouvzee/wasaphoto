package database

var user_table = `CREATE TABLE IF NOT EXISTS User
				(
					userID INTEGER AUTO_INCREMENT,
					username INTEGER NOT NULL,
					PRIMARY KEY(userID)
				);`

var post_table = `CREATE TABLE IF NOT EXISTS Post
				(
					postID INTEGER NOT NULL,
					userID INTEGER NOT NULL,
					caption TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY(postID, userID),
					CONSTRAINT fk_post
						FOREIGN KEY(userID) REFERENCES user(userID) ON DELETE CASCADE
				);`

//like, comment, follow, ban
var like_table = `CREATE TABLE IF NOT EXISTS Like
				(
					userID INTEGER NOT NULL,
					postID INTEGER NOT NULL,
					creatorID INTEGER NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY(userID, postID, creatorID),
					CONSTRAINT fk_like
						FOREIGN KEY(userID) REFERENCES user(userID) ON DELETE CASCADE,
						FOREIGN KEY(postID, creatorID) REFERENCES post(postID,userID)
				);`

var comment_table = `CREATE TABLE IF NOT EXISTS Comment
				(
					userID INTEGER NOT NULL,
					commentID INTEGER NOT NULL,
					textComment TEXT NOT NULL,
					postID INTEGER NOT NULL,
					creatorID INTEGER NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					PRIMARY KEY(commentID, creatorID, postID),
					CONSTRAINT fk_comment
						FOREIGN KEY(userID) REFERENCES user(userID) ON DELETE CASCADE,
						FOREIGN KEY(creatorID, postID) REFERENCES post(userID, postID)
				);`

var follow_table = `CREATE TABLE IF NOT EXISTS Follow
				(
					followerID INTEGER NOT NULL,
					followedID INTEGER NOT NULL,
					PRIMARY KEY(followerID, followedID),
					CONSTRAINT fk_follow
						FOREIGN KEY(followerID) REFERENCES user(userID),
						FOREIGN KEY(followedID) REFERENCES user(userID)
				);`

var ban_table = `CREATE TABLE IF NOT EXISTS Ban
				(
					bannerID INTEGER NOT NULL,
					bannedID INTEGER NOT NULL,
					PRIMARY KEY(bannerID, bannedID),
					CONSTRAINT fk_ban
						FOREIGN KEY(bannerID) REFERENCES user(userID),
						FOREIGN KEY(bannedID) REFERENCES user(userID)
				);`
