<template>
  <div class="home-view">

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
      offset: 0
    }
  },
  computed: {
    currentUserId() {
      try {
        const user = JSON.parse(localStorage.getItem('user'));
        return user ? user.UserID : null;
      } catch (error) {
        return null;
      }
    }
  },
  async created() {
    await this.loadPosts();
    
    // Listen for new posts
    window.addEventListener('postUploaded', this.onPostUploaded);
  },
  beforeUnmount() {
    window.removeEventListener('postUploaded', this.onPostUploaded);
  },
  methods: {
    async loadPosts(reset = true) {
      if (!this.currentUserId) return;

      this.loading = true;
      this.error = null;

      try {
        const currentOffset = reset ? 0 : this.offset;
        const response = await this.$axios.get(`/profiles/${this.currentUserId}/feed`, {
          params: {
            limit: this.limit,
            offset: currentOffset
          }
        });

        const newPosts = response.data || [];
        
        if (reset) {
          this.posts = newPosts;
          this.offset = newPosts.length;
        } else {
          this.posts.push(...newPosts);
          this.offset += newPosts.length;
        }

        // Check if there are more posts
        this.hasMore = newPosts.length === this.limit;

      } catch (error) {
        console.error('Feed load error:', error);
        if (error.response?.status === 401) {
          this.$router.push('/login');
          return;
        }
        this.error = 'Failed to load feed. Please try again.';
      } finally {
        this.loading = false;
      }
    },

    async loadMorePosts() {
      if (this.loading || !this.hasMore) return;
      await this.loadPosts(false);
    },

    async onPostUploaded() {
      // Reload the feed when a new post is uploaded
      await this.loadPosts(true);
    },

    openPostModal(post) {
      this.selectedPost = { ...post };
    },

    closePostModal() {
      this.selectedPost = null;
    },

    onPostDeleted(deletedPost) {
      // Remove the deleted post from the feed
      this.posts = this.posts.filter(p => p.PhotoID !== deletedPost.PhotoID);
      this.closePostModal();
    },

    onPostUpdated(updatedData) {
      // Update the post in the feed
      const postIndex = this.posts.findIndex(p => p.PhotoID === updatedData.PhotoID);
      if (postIndex !== -1) {
        this.posts[postIndex] = {
          ...this.posts[postIndex],
          Nlike: updatedData.Nlike,
          Liked: updatedData.Liked,
          Ncomment: updatedData.Ncomment
        };
        
        this.$forceUpdate();
      }

      // Update also the modal if open
      if (this.selectedPost && this.selectedPost.PhotoID === updatedData.PhotoID) {
        this.selectedPost = {
          ...this.selectedPost,
          Nlike: updatedData.Nlike,
          Liked: updatedData.Liked,
          Ncomment: updatedData.Ncomment
        };
      }
    }
  }
}
</script>

<style scoped>
.home-view {
  max-width: 700px; /* ✅ AUMENTATO: da 600px a 700px */
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
  max-width: 600px; /* ✅ AUMENTATO: da 470px a 600px */
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
