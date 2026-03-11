import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";

import MainView from "../views/MainView.vue";
import Cookies from "js-cookie";

import { useAuthStore } from "@/store/auth";

// 同步登入狀態
const syncSessionFromCookie = () => {
  const hasCookie = Cookies.get("session_id") !== undefined;
  useAuthStore().setSession(hasCookie);
};

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "main",
    component: MainView,
    children: [
      {
        path: "",
        name: "home",
        component: () => import("@/pages/MainPage.vue"),
        beforeEnter: (_to, _from, next) => {
          syncSessionFromCookie();
          if (Cookies.get("session_id") === undefined) {
            useAuthStore().openLoginModal();
          }
          next();
        },
      },
      {
        path: "/folder",
        name: "folder",
        component: () => import("@/pages/FolderPage.vue"),
      },
      {
        path: "/markdown-writer",
        name: "markdown-writer",
        component: () => import("@/pages/MdWriter.vue"),
      },
      {
        path: "/error",
        name: "error",
        component: () => import("@/pages/ErrorPage.vue"),
      },
      // match all route
      {
        path: ":pathMatch(.*)*",
        name: "notFound",
        component: () => import("@/pages/ErrorPage.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach((_to, _from, next) => {
  syncSessionFromCookie();
  next();
});

export default router;
