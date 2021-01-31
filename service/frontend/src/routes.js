import Feed from "./components/feed";
import { Create } from "./components/post";
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
];

export default routes;
