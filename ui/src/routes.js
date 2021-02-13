import Feed from "./components/feed";
import { Create, Post } from "./components/post";
import { Layout } from "./components/layout";
import Profile from "./components/profile";

const routes = [
    {
        path: "/",
        exact: true,
        component: Feed,
        layout: Layout,
    },
    {
        path: "/post",
        exact: true,
        component: Create,
        layout: Layout,
    },
    {
        path: "/p/:id",
        component: Post,
        layout: Layout,
    },
    {
        path: "/u/:username",
        component: Profile,
        layout: Layout,
    },
];

export default routes;
