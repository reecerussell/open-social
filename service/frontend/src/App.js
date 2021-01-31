import React, { Suspense } from "react";
import {
    BrowserRouter as Router,
    Route,
    Redirect,
    Switch,
} from "react-router-dom";

import "./scss/styles.scss";
import routes from "./routes";

const renderRoute = (route, props) => {
    const component = <route.component {...props} />;

    if (route.layout) {
        return <route.layout>{component}</route.layout>;
    }

    return component;
};

const App = () => (
    <Router>
        <Switch>
            <Suspense fallback={"loading..."}>
                {routes.map((route, key) => (
                    <Route
                        key={key}
                        path={route.path}
                        exact={route.exact}
                        render={props => renderRoute(route, props)}
                    />
                ))}
            </Suspense>
            <Redirect to="/" />
        </Switch>
    </Router>
);

export default App;
