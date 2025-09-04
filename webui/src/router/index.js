import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import HomeView from '../views/HomeView.vue'
import SearchView from '../views/SearchView.vue'

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
                        // Se i dati sono corrotti, pulisci e vai al login
                        localStorage.clear();
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

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    
    let isAuthenticated = false;
    let currentUserId = null;
    
    if (token && user && token !== '0' && token !== 'null' && token !== 'undefined') {
        try {
            const userData = JSON.parse(user);
            if (userData && userData.UserID && userData.UserID > 0) {
                isAuthenticated = true;
                currentUserId = userData.UserID;
            }
        } catch (e) {
            // Dati corrotti, pulisci e vai al login
            localStorage.clear();
            isAuthenticated = false;
        }
    }

    // Se non autenticato e sta andando su una pagina protetta
    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login');
        return;
    }

    // Se autenticato e va al login, reindirizza alla home
    if (to.path === '/login' && isAuthenticated) {
        next(`/profiles/${currentUserId}/feed`);
        return;
    }

    // Controlla se sta andando su un profilo/feed che non è il suo
    // ma potrebbe essere rimasto in cache dalla sessione precedente
    if (isAuthenticated && (to.path.includes('/profiles/') || to.path.includes('/feed'))) {
        const routeUserId = parseInt(to.params.userID);
        
        // Se l'utente nella route non corrisponde a quello autenticato
        // E sta andando sulla "sua" home/profilo (probabilmente cache vecchia)
        if (!isNaN(routeUserId) && routeUserId !== currentUserId && 
            (to.path.endsWith('/feed') || to.path === `/profiles/${routeUserId}`)) {
            
            // Solo se sembra che stia andando sulla sua home ma con ID sbagliato
            if (to.path.endsWith('/feed')) {
                next(`/profiles/${currentUserId}/feed`);
                return;
            }
        }
    }

    next();
});

export default router
