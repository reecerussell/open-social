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

export const followProfile = userId => ({
    type: types.FOLLOW_PROFILE,
    userId: userId,
});

export const followProfileSuccess = userId => ({
    type: types.FOLLOW_PROFILE_SUCCESS,
    userId: userId,
});

export const followProfileError = error => ({
    type: types.FOLLOW_PROFILE,
    error: error,
});

export const unfollowProfile = userId => ({
    type: types.UNFOLLOW_PROFILE,
    userId: userId,
});

export const unfollowProfileSuccess = userId => ({
    type: types.UNFOLLOW_PROFILE_SUCCESS,
    userId: userId,
});

export const unfollowProfileError = error => ({
    type: types.UNFOLLOW_PROFILE,
    error: error,
});
