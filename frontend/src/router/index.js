import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home'
import Error from '../views/Error'
import store from '../store'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('@/views/About.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/error',
    name: 'Error',
    component: Error
  },
  {
    path: '/newTag',
    name: 'NewTag',
    component: () => import('@/views/NewTag.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/editTag/:id',
    name: 'EditTag',
    component: () => import('@/views/EditTag.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/tags',
    name: 'Tags',
    component: () => import('@/views/Tags.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/newLang',
    name: 'NewLang',
    component: () => import('@/views/NewLang.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/editLang/:id',
    name: 'EditLang',
    component: () => import('@/views/EditLang.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/langs',
    name: 'Langs',
    component: () => import('@/views/Langs.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/newTranslation',
    name: 'NewTranslation',
    component: () => import('@/views/NewTranslation.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/editTranslation/:id',
    name: 'EditTranslation',
    component: () => import('@/views/EditTranslation.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/newUser',
    name: 'NewUser',
    component: () => import('@/views/NewUser.vue'),
    meta: {
      requiresAuth: true
    }
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // this route requires auth, check if logged in
    // if not, redirect to login page.
    if (!store.getters['auth/isLoggedIn']) {
      next({
        path: '/login'
      })
    } else {
      next();
    }
  } else {
    next(); // make sure to always call next()!
  }
});

export default router
