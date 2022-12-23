import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  // {
  //   path: '/login',
  //   component: () => import('@/views/login/index'),
  //   hidden: true
  // },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  // {
  //   path: '/',
  //   component: Layout,
  //   redirect: '/dashboard',
  //   children: [{
  //     path: 'dashboard',
  //     name: 'Dashboard',
  //     component: () => import('@/views/dashboard/index'),
  //     meta: { title: 'Dashboard', icon: 'dashboard' }
  //   }]
  // },
  {
    path: '/',
    component: Layout,
    redirect: '/probe/http'
  },
  {
    path: '/probe',
    component: Layout,
    // redirect: 'http',
    meta: { title: 'Probe', icon: 'radar' },
    children: [
      {
        path: 'http',
        name: 'http',
        component: () => import('@/views/probe/index'),
        meta: { title: 'HTTPProbe', icon: 'radar' }
      },
    ]
  },
  {
    path: '/scenarios',
    component: Layout,
    // redirect: "/scenarios",
    meta: { title: 'Scenarios', icon: 'el-icon-s-help' },
    children: [
        {
          path: '',
          name: 'list',
          hidden: true,
          component: () => import('@/views/scenarios/list'),
          meta: { title: 'list'},
        },
        {
          path: 'editor/:namespace/:name',
          hidden: true,
          component: () => import('@/views/scenarios/editor'),
          meta: { title: 'Scenarios Detail'},
        },
        {
          path: 'editor',
          hidden: true,
          component: () => import('@/views/scenarios/editor'),
          meta: { title: 'Create Scenarios'},
        }
    ]
  },
  {
    path: '/contact',
    component: Layout,
    meta: { title: 'Contact', icon: 'el-icon-s-platform' },
    children: [
      {
        path: 'contactgroup',
        name: 'ContactGroup',
        component: () => import('@/views/contact_group/index'),
        meta: { title: 'ContactGroup', icon: 'user-group' }
      }
    ]
  },
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
