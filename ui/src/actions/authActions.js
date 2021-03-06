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

export const login = () => ({
  type: types.LOGIN,
});

export const loginSuccess = () => ({
  type: types.LOGIN_SUCCESS,
});

export const loginError = error => ({
  type: types.LOGIN_ERROR,
  error: error,
});
