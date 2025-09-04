<template>
  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Search Users</h1>
  </div>
  
  <div class="mb-3">
    <input 
      type="text" 
      class="form-control" 
      placeholder="Search for users..."
      v-model="searchQuery"
      @input="searchUsers"
    >
  </div>
  
  <div v-if="results.length">
    <div class="list-group">
      <div v-for="user in results" :key="user.UserID" class="list-group-item">
        <RouterLink :to="`/profiles/${user.UserID}`" class="text-decoration-none">
          {{ user.Username }}
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SearchView',
  data() {
    return {
      searchQuery: '',
      results: []
    }
  },
  methods: {
    async searchUsers() {
      if (this.searchQuery.length < 2) {
        this.results = []
        return
      }
      
      try {
        const response = await this.$axios.get(`/profiles?search=${this.searchQuery}`)
        this.results = response.data
      } catch (error) {
        console.error('Search error:', error)
      }
    }
  }
}
</script>