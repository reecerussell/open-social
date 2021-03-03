import Feed from "./components/feed";
import { Create, Post } from "./components/post";
import { Layout } from "./components/layout";
import Profile from "./components/profile";
import { Register } from "./components/auth";

const routes = [
  {
    path: "/",
    exact: true,
    component: Feed,
    layout: Layout,
    auth: true,
  },
  {
    path: "/post",
    exact: true,
    component: Create,
    layout: Layout,
    auth: true,
  },
  {
    path: "/p/:id",
    component: Post,
    layout: Layout,
    auth: true,
  },
  {
    path: "/u/:username",
    component: Profile,
    layout: Layout,
    auth: true,
  },
  {
    path: "/register",
    component: Register,
  },
];

export default routes;
