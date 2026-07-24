import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'

const routes = [
  {
    path: '/',
    name: 'dashboard',
    component: DashboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('../views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile/timeline',
    name: 'timeline',
    component: () => import('../views/TimelineView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile/timeline/new',
    name: 'log-event',
    component: () => import('../views/LogEventView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/assets',
    name: 'assets',
    component: () => import('../views/AssetsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/goals',
    name: 'goals',
    component: () => import('../views/GoalsView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/portfolio',
    name: 'portfolio',
    component: () => import('../views/PortfolioView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/review',
    name: 'review',
    component: () => import('../views/ReviewView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/onboarding/interview',
    name: 'onboarding-interview',
    component: () => import('../views/OnboardingView.vue'),
    meta: { requiresAuth: true }
  },

  {
    path: '/ips',
    name: 'ips',
    component: () => import('../views/IPSView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/thesis',
    name: 'thesis',
    component: () => import('../views/ThesisView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/thesis/new',
    name: 'thesis-new',
    component: () => import('../views/ThesisDetailView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/thesis/:ticker',
    name: 'thesis-detail',
    component: () => import('../views/ThesisDetailView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/signup',
    name: 'signup',
    component: () => import('../views/SignupView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/review',
    name: 'review',
    component: () => import('../views/ReviewView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/review/active',
    name: 'review-active',
    component: () => import('../views/ActiveReviewView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/review/:month',
    name: 'review-detail',
    component: () => import('../views/ActiveReviewView.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Simple Route Guard using local token
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
