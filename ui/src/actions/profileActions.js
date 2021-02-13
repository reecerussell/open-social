import * as types from "./profileActionTypes";

export const loadProfile = () => ({
    type: types.LOAD_PROFILE,
});

export const loadProfileSuccess = profile => ({
    type: types.LOAD_PROFILE_SUCCESS,
    profile: profile,
});

export const loadProfileError = error => ({
    type: types.LOAD_PROFILE_ERROR,
    error: error,
});
