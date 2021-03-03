import { Api, Auth } from "../utils";
import { authActions } from "../actions";

export const submitRegistration = data => dispatch => {
  dispatch(authActions.register());

  return Api.post("users/register", data)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      const { accessToken, expires } = res.data.token;
      Auth.setAccessToken(accessToken, expires);

      dispatch(authActions.registerSuccess());
    })
    .catch(err => dispatch(authActions.registerError(err.toString())));
};
