import * as types from "../actions/postActionTypes";
import initialState from "../store/initialState";

const postReducer = (state = initialState.post, action) => {
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
                hasLiked: true,
                likes: state.likes + 1,
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
                hasLiked: false,
                likes: state.likes - 1,
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
                loading: false,
                error: null,
                ...action.post,
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
