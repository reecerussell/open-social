import React, { Suspense, useEffect } from "react";
import PropTypes from "prop-types";
import {
  BrowserRouter as Router,
  Route,
  Redirect,
  Switch,
} from "react-router-dom";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import routes from "./routes";
import * as userApi from "./api/user";
import { isAuthenticated } from "./utils/auth";

import "./scss/styles.scss";

const renderRoute = (route, props, userInfo) => {
  if (route.auth && !isAuthenticated()) {
    return <Redirect to="/login" />;
  }

  const component = <route.component {...props} />;

  if (route.layout) {
    return <route.layout userInfo={userInfo}>{component}</route.layout>;
  }

  return component;
};

const App = ({ fetchUserInfo, userInfo, error, loading }) => {
  useEffect(() => {
    if (isAuthenticated()) {
      fetchUserInfo();
    }
  }, [fetchUserInfo]);

  return (
    <Router>
      <Switch>
        <Suspense fallback={"loading..."}>
          {routes.map((route, key) => (
            <Route
              key={key}
              path={route.path}
              exact={route.exact}
              render={props => renderRoute(route, props, userInfo)}
            />
          ))}
        </Suspense>
        <Redirect to={isAuthenticated() ? "/" : "/login"} />
      </Switch>
    </Router>
  );
};

App.propTypes = {
  userInfo: PropTypes.shape({
    id: PropTypes.string.isRequired,
    username: PropTypes.string.isRequired,
    mediaId: PropTypes.string,
    followerCount: PropTypes.number.isRequired,
  }),
  error: PropTypes.string,
  loading: PropTypes.bool,
};

App.defaultProps = {
  userInfo: {
    id: "",
    username: "",
    mediaId: null,
    followerCount: 0,
  },
  error: null,
  loading: false,
};

const mapStateToProps = state => ({
  userInfo: {
    id: state.user.id,
    username: state.user.username,
    mediaId: state.user.mediaId,
    followerCount: state.user.followerCount,
  },
  error: state.user.error,
  loading: state.user.loading,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      fetchUserInfo: userApi.fetchInfo,
    },
    dispatch
  );

export default connect(mapStateToProps, mapDispatchToProps)(App);
