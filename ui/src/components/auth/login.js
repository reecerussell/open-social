import React, { useState } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { Redirect, Link } from "react-router-dom";
import { authApi } from "../../api";
import { Textbox } from "../form";

const defaultState = {
  username: "",
  password: "",
};

const Login = ({ error, loading, success, submitLogin }) => {
  const [formData, setFormData] = useState(defaultState);

  const handleUpdate = e => {
    const { name, value } = e.target;
    const newState = {
      ...formData,
    };
    newState[name] = value;
    setFormData(newState);
  };

  const handleSubmit = e => {
    e.preventDefault();

    if (loading) {
      return;
    }

    submitLogin({
      username: formData.username,
      password: formData.password,
    });
  };

  if (success) {
    return <Redirect to="/" />;
  }

  return (
    <div className="container sections">
      <div className="row">
        <div className="col-md-8 col-lg-6 col-xl-5">
          <div className="section p-4">
            <h1 className="header-1">Open Social</h1>
            <p className="text-muted small pt-1 info">
              The brand new, simple social media platform, using leading
              open-source technologies!
            </p>
            <hr />
            <h2 className="display-3">Log In</h2>
            <p className="text-center">
              Enter your account details to login to your account.
            </p>
            {error && <p className="text-center text-danger">{error}</p>}
            <form onSubmit={handleSubmit}>
              <Textbox
                label="Username"
                name="username"
                placeholder="Enter a username"
                max="255"
                onChange={handleUpdate}
                value={formData.username}
                autoComplete="username"
                required
              />
              <Textbox
                label="Password"
                name="password"
                placeholder="Enter a password"
                onChange={handleUpdate}
                value={formData.password}
                type="password"
                autoComplete="current-password"
                required
              />

              <div className="form-group">
                <button
                  type="submit"
                  className="btn btn-primary btn-block float-right"
                >
                  {loading ? "LOGGING IN" : "LOG IN"}
                </button>
              </div>
            </form>
            <p className="text-center mb-0">
              <Link to="/register">Don't have an account?</Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

Login.propTypes = {
  loading: PropTypes.bool,
  error: PropTypes.string,
  success: PropTypes.bool,
  submitLogin: PropTypes.func.isRequired,
};

Login.defaultProps = {
  loading: false,
  error: null,
  success: false,
};

const mapStateToProps = state => ({
  error: state.auth.error,
  loading: state.auth.loading,
  success: state.auth.success,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      submitLogin: authApi.submitLogin,
    },
    dispatch
  );

export default connect(mapStateToProps, mapDispatchToProps)(Login);
