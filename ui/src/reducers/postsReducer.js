import * as types from "../actions/postActionTypes";
import initialState from "../store/initialState";

const updatePostLikedStatus = (posts, postId, liked) => {
    if (!posts) {
        return [];
    }

    const newPosts = [];

    for (let i = 0; i < posts.length; i++) {
        const post = posts[i];

        if (post.id === postId) {
            if (liked) {
                post.hasLiked = true;
                post.likes += 1;
            } else {
                post.hasLiked = false;
                post.likes -= 1;
            }
        }

        newPosts.push(post);
    }

    return newPosts;
};

const mergePosts = (posts, newPost) => {
    if (!posts) {
        return [];
    }

    return posts.filter(p => p.id !== newPost.Id).concat(newPost);
};

const postReducer = (state = initialState.posts, action) => {
    switch (action.type) {
        case types.CREATE_POST:
            return {
                ...state,
                loading: false,
            };
        case types.CREATE_POST_SUCCESS:
            return {
                ...state,
                loading: false,
                error: null,
            };
        case types.CREATE_POST_ERROR:
            return {
                ...state,
                error: action.error,
                loading: false,
            };
        case types.LIKE_POST:
            return {
                ...state,
                loading: true,
            };
        case types.LIKE_POST_SUCCESS:
            return {
                ...state,
                posts: updatePostLikedStatus(state.posts, action.id, true),
                loading: false,
                error: null,
            };
        case types.LIKE_POST_ERROR:
            return {
                ...state,
                error: action.error,
                loading: false,
            };
        case types.UNLIKE_POST:
            return {
                ...state,
                loading: true,
            };
        case types.UNLIKE_POST_SUCCESS:
            return {
                ...state,
                posts: updatePostLikedStatus(state.posts, action.id, false),
                loading: false,
                error: null,
            };
        case types.UNLIKE_POST_ERROR:
            return {
                ...state,
                error: action.error,
                loading: false,
            };
        case types.LOAD_POST:
            return {
                ...state,
                loading: true,
            };
        case types.LOAD_POST_SUCCESS:
            return {
                ...state,
                posts: mergePosts(state.posts, action.post),
                loading: false,
                error: null,
            };
        case types.LOAD_POST_ERROR:
            return {
                ...state,
                error: action.error,
                loading: false,
            };
        default:
            return state;
    }
};

export default postReducer;
