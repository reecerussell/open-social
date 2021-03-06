import { Api, Auth } from "../utils";
import { authActions } from "../actions";

export const submitRegistration = data => dispatch => {
  dispatch(authActions.register());

  return Api.post("auth/register", data)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      const { token, expires } = res.data.accessToken;
      Auth.setAccessToken(token, expires);

      dispatch(authActions.registerSuccess());
    })
    .catch(err => dispatch(authActions.registerError(err.toString())));
};

export const submitLogin = data => dispatch => {
  dispatch(authActions.login());

  return Api.post("auth/token", data)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      const { token, expires } = res.data;
      Auth.setAccessToken(token, expires);

      dispatch(authActions.loginSuccess());
    })
    .catch(err => dispatch(authActions.loginError(err.toString())));
};
