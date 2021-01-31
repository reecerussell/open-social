import * as types from "../actions/feedActionTypes";
import initialState from "../store/initialState";

const feedReducer = (state = initialState.feed, action) => {
    switch (action.type) {
        case types.LOAD_FEED:
            return {
                ...state,
                loading: false,
            };
        case types.LOAD_FEED_SUCCESS:
            return {
                ...state,
                items: action.feedItems,
                loading: false,
                error: null,
            };
        case types.LOAD_FEED_ERROR:
            return {
                ...state,
                error: action.error,
                loading: false,
            };
        default:
            return state;
    }
};

export default feedReducer;
