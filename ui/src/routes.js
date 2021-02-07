import Feed from "./components/feed";
import { Create, Post } from "./components/post";
import { Layout } from "./components/layout";

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
        path: "/post/:id",
        component: Post,
        layout: Layout,
    },
];

export default routes;
