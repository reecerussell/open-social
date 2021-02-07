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

export const unlikePost = id => ({
    type: types.UNLIKE_POST,
    id: id,
});

export const unlikePostSuccess = id => ({
    type: types.UNLIKE_POST_SUCCESS,
    id: id,
});

export const unlikePostError = error => ({
    type: types.UNLIKE_POST_ERROR,
    error: error,
});

export const loadPost = () => ({
    type: types.LOAD_POST,
});

export const loadPostSuccess = post => ({
    type: types.LOAD_POST_SUCCESS,
    post: post,
});

export const loadPostError = error => ({
    type: types.LOAD_POST_ERROR,
    error: error,
});
