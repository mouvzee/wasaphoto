import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import ProfileHeader from './components/ProfileHeader.vue'
import ProfileStats from './components/ProfileStats.vue'
import PostCard from './components/PostCard.vue';
import PostModal from './components/PostModal.vue';
import UploadModal from './components/UploadModal.vue';
import FeedPost from './components/FeedPost.vue';

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("ProfileHeader", ProfileHeader);
app.component("ProfileStats", ProfileStats);
app.component("PostCard", PostCard);
app.component("PostModal", PostModal);
app.component("UploadModal", UploadModal);
app.component("FeedPost", FeedPost);

app.use(router)
app.mount('#app')
