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
        default:
            return state;
    }
};

export default postReducer;
