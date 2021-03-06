import * as types from "./userActionTypes";

export const loadInfo = () => ({
  type: types.LOAD_INFO,
});

export const loadInfoSuccess = info => ({
  type: types.LOAD_INFO_SUCCESS,
  info: info,
});

export const loadInfoError = error => ({
  type: types.LOAD_INFO_ERROR,
  error: error,
});
