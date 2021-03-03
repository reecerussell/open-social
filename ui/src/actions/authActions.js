import * as types from "./authActionTypes";

export const register = () => ({
  type: types.REGISTER,
});

export const registerSuccess = () => ({
  type: types.REGISTER_SUCCESS,
});

export const registerError = error => ({
  type: types.REGISTER_ERROR,
  error: error,
});
