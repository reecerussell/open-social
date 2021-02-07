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

export const likePost = id => ({
    type: types.LIKE_POST,
    id: id,
});

export const likePostSuccess = id => ({
    type: types.LIKE_POST_SUCCESS,
    id: id,
});

export const likePostError = error => ({
    type: types.LIKE_POST_ERROR,
    error: error,
});
