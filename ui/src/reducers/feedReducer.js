import { feedTypes, postTypes } from "../actions";
import initialState from "../store/initialState";

const markPostAsLiked = (feedItems, id) => {
    const newItems = [];

    for (let i = 0; i < feedItems.length; i++) {
        const feedItem = feedItems[i];
        console.log(id, feedItem);
        if (feedItem.id === id) {
            feedItem.hasUserLiked = true;
            feedItem.likes += 1;
        }

        newItems.push(feedItem);
    }

    return newItems;
};

const feedReducer = (state = initialState.feed, action) => {
    switch (action.type) {
        case feedTypes.LOAD_FEED:
            return {
                ...state,
                loading: false,
            };
        case feedTypes.LOAD_FEED_SUCCESS:
            return {
                ...state,
                items: action.feedItems,
                loading: false,
                error: null,
            };
        case feedTypes.LOAD_FEED_ERROR:
            return {
                ...state,
                error: action.error,
                loading: false,
            };
        case postTypes.LIKE_POST:
            return {
                ...state,
                loading: true,
            };
        case postTypes.LIKE_POST_SUCCESS:
            return {
                ...state,
                items: markPostAsLiked(state.items, action.id),
                loading: false,
                error: null,
            };
        case postTypes.LIKE_POST_ERROR:
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
