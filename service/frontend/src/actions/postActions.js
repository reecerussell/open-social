import * as types from "./postActionTypes";

export const createPost = () => ({
    type: types.CREATE_POST,
});

export const createPostSuccess = () => ({
    type: types.CREATE_POST_SUCCESS,
});

export const createPostError = error => ({
    type: types.CREATE_POST_ERROR,
    error: error,
});
