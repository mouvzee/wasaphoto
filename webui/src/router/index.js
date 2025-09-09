import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import HomeView from '../views/HomeView.vue'
import SearchView from '../views/SearchView.vue'
import { validateToken, clearAuthData } from '@/services/axios.js'

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/', 
            redirect: () => {
                const token = localStorage.getItem('token');
                const user = localStorage.getItem('user');
                
                if (token && user && token !== '0' && token !== 'null' && token !== 'undefined') {
                    try {
                        const userData = JSON.parse(user);
                        if (userData.UserID && userData.UserID > 0) {
                            return `/profiles/${userData.UserID}/feed`;
                        }
                    } catch (e) {
                        clearAuthData();
                    }
                }
                return '/login';
            }
        },
        {path: '/login', component: LoginView},
        {
            path: '/profiles/:userID', 
            component: ProfileView,
            meta: { requiresAuth: true }
        },
        {
            path: '/profiles/:userID/feed', 
            component: HomeView,
            meta: { requiresAuth: true }
        },
        {
            path: '/search', 
            component: SearchView,
            meta: { requiresAuth: true }
        },
    ]
})

router.beforeEach(async (to, from, next) => {
    if (to.meta.requiresAuth) {
        try {
            const isValid = await validateToken();
            
            if (!isValid) {
                clearAuthData();
                next('/login');
                return;
            }
            next();
            
        } catch (error) {
            console.error('Router guard error:', error);
            clearAuthData();
            next('/login');
            return;
        }
    } else {
        next();
    }
});

export default router
