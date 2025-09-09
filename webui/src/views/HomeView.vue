<template>
  <div class="home-view">
    <!-- Header -->
    <div class="home-header">
      <h1 class="page-title">Feed</h1>
    </div>

    <!-- Loading spinner -->
    <LoadingSpinner :loading="loading && posts.length === 0">
      <!-- Error message -->
      <ErrorMsg v-if="error" :msg="error" />
      
      <!-- Feed -->
      <div v-else class="feed-container">
        <!-- Posts -->
        <div class="feed-posts" v-if="posts.length > 0">
          <FeedPost
            v-for="post in posts"
            :key="post.PhotoID"
            :post="post"
            @openModal="openPostModal"
            @postUpdated="onPostUpdated"
            @showLikes="showLikesList"
          />
          
          <!-- Load more button -->
          <div class="load-more-container" v-if="hasMore">
            <button 
              class="btn btn-outline-primary load-more-btn"
              @click="loadMorePosts"
              :disabled="loading"
            >
              {{ loading ? 'Loading...' : 'Load More' }}
            </button>
          </div>
        </div>

        <!-- Empty state -->
        <div v-else class="empty-feed">
          <div class="empty-content">
            <i class="fas fa-camera fa-3x text-muted mb-3"></i>
            <h3>Welcome to WASAPhoto!</h3>
            <p class="text-muted">
              Follow some users to see their posts here, or 
              <RouterLink to="/search" class="text-primary">discover new people</RouterLink>.
            </p>
          </div>
        </div>
      </div>
    </LoadingSpinner>

    <!-- Post Modal -->
    <PostModal
      v-if="selectedPost"
      :postData="selectedPost"
      @close="closePostModal"
      @deleted="onPostDeleted"
      @postUpdated="onPostUpdated"
    />

    <div v-if="showLikesModal" class="modal-overlay" @click="closeLikesModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h5>Likes</h5>
          <button class="btn-close" @click="closeLikesModal">&times;</button>
        </div>
        <div class="modal-body">
          <LoadingSpinner :loading="loadingLikes">
            <div v-if="likes.length > 0" class="users-list">
              <div v-for="like in likes" :key="like.UserID" class="user-item">
                <div class="user-info">
                  <div class="user-avatar" :style="getUserAvatarStyle(like.Username)">
                    <span class="avatar-initials">{{ getUserInitials(like.Username) }}</span>
                  </div>
                  <RouterLink 
                    :to="`/profiles/${like.UserID}`" 
                    class="username"
                    @click="closeLikesModal"
                  >
                    {{ like.Username }}
                  </RouterLink>
                </div>
              </div>
            </div>
            <div v-else class="text-center text-muted">
              No likes yet
            </div>
          </LoadingSpinner>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomeView',
  data() {
    return {
      posts: [],
      loading: false,
      error: null,
      selectedPost: null,
      hasMore: true,
      limit: 10,
      offset: 0,
      showLikesModal: false,
      likes: [],
      loadingLikes: false,
      currentLikesPhoto: null
    }
  },
  
  async created() {
    await this.loadPosts();
    
    window.addEventListener('postUploaded', this.handlePostUploaded);
  },

  beforeUnmount() {
    window.removeEventListener('postUploaded', this.handlePostUploaded);
  },

  methods: {
    async loadPosts() {
      this.loading = true;
      this.error = null;

      try {
        const user = JSON.parse(localStorage.getItem('user'));
        if (!user?.UserID) {
          this.error = 'User not found';
          return;
        }

        const response = await this.$axios.get(`/profiles/${user.UserID}/feed`, {
          params: {
            limit: this.limit,
            offset: this.offset
          }
        });

        const newPosts = response.data || [];
        
        if (this.offset === 0) {
          this.posts = newPosts;
        } else {
          this.posts = [...this.posts, ...newPosts];
        }

        this.hasMore = newPosts.length === this.limit;

      } catch (error) {
        console.error('Feed load error:', error);
        if (error.response?.status === 404 || error.response?.status === 403) {
          this.error = 'Unable to load feed. Please check your following list.';
        } else {
          this.error = 'Failed to load feed. Please try again.';
        }
      } finally {
        this.loading = false;
      }
    },

    async loadMorePosts() {
      if (this.loading || !this.hasMore) return;
      
      this.offset += this.limit;
      await this.loadPosts();
    },

    async handlePostUploaded() {
      this.offset = 0;
      this.hasMore = true;
      await this.loadPosts();
    },

    openPostModal(post) {
      this.selectedPost = { ...post };
    },

    closePostModal() {
      this.selectedPost = null;
    },

    showLikesList(photoData) {
      this.currentLikesPhoto = photoData;
      this.showLikesModal = true;
      this.loadLikes();
    },

    closeLikesModal() {
      this.showLikesModal = false;
      this.likes = [];
      this.currentLikesPhoto = null;
    },

    async loadLikes() {
      if (!this.currentLikesPhoto) return;
      
      this.loadingLikes = true;
      try {
        const response = await this.$axios.get(
          `/profiles/${this.currentLikesPhoto.UserID}/posts/${this.currentLikesPhoto.PhotoID}/likes`
        );
        
        this.likes = response.data || [];
        
      } catch (error) {
        console.error('Error loading likes:', error);
        this.likes = [];
      } finally {
        this.loadingLikes = false;
      }
    },

    getUserInitials(username) {
      if (!username) return '?';
      return username.charAt(0).toUpperCase();
    },

    getUserAvatarStyle(username) {
      const colors = [
        '#e74c3c', '#3498db', '#2ecc71', '#f39c12', 
        '#9b59b6', '#1abc9c', '#34495e', '#e67e22'
      ];
      
      const usernameStr = username || 'default';
      let hash = 0;
      for (let i = 0; i < usernameStr.length; i++) {
        hash = usernameStr.charCodeAt(i) + ((hash << 5) - hash);
      }
      
      const color = colors[Math.abs(hash) % colors.length];
      
      return {
        backgroundColor: color,
        color: 'white'
      };
    },

    onPostUpdated(updatedData) {
      const postIndex = this.posts.findIndex(p => p.PhotoID === updatedData.PhotoID);
      if (postIndex !== -1) {
        this.posts[postIndex] = {
          ...this.posts[postIndex],
          Nlike: updatedData.Nlike,
          Liked: updatedData.Liked,
          Ncomment: updatedData.Ncomment
        };
      }

      if (this.selectedPost && this.selectedPost.PhotoID === updatedData.PhotoID) {
        this.selectedPost = {
          ...this.selectedPost,
          Nlike: updatedData.Nlike,
          Liked: updatedData.Liked,
          Ncomment: updatedData.Ncomment
        };
      }
    },

    onPostDeleted(deletedPost) {
      this.posts = this.posts.filter(p => p.PhotoID !== deletedPost.PhotoID);
      this.closePostModal();
    }
  }
}
</script>

<style scoped>
.home-view {
  max-width: 700px;
  margin: 0 auto;
  padding: 20px;
}

.home-header {
  margin-bottom: 30px;
  text-align: center;
}

.page-title {
  font-size: 28px;
  font-weight: 300;
  color: #262626;
  margin: 0;
}

.feed-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.feed-posts {
  width: 100%;
  max-width: 600px;
}

.load-more-container {
  text-align: center;
  margin: 20px 0;
}

.load-more-btn {
  padding: 12px 24px;
  border-radius: 8px;
  font-weight: 500;
}

.empty-feed {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  text-align: center;
}

.empty-content h3 {
  color: #262626;
  margin-bottom: 12px;
}

.empty-content p {
  max-width: 300px;
  line-height: 1.5;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 400px;
  max-height: 80vh;
  overflow: hidden;
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e1e8ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h5 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
  max-height: 60vh;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #8e8e8e;
  line-height: 1;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.user-item:last-child {
  border-bottom: none;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  font-weight: 600;
  font-size: 18px;
}

.avatar-initials {
  line-height: 1;
}

.username {
  color: #262626;
  font-weight: 500;
  text-decoration: none;
}

.username:hover {
  color: #262626;
  text-decoration: none;
}

@media (max-width: 768px) {
  .home-view {
    padding: 10px 0;
    max-width: none;
  }
  
  .home-header {
    padding: 0 20px;
    margin-bottom: 20px;
  }
  
  .page-title {
    font-size: 24px;
  }
  
  .feed-posts {
    max-width: none;
  }
}
</style>
