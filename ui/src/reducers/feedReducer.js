import { feedTypes, postTypes } from "../actions";
import initialState from "../store/initialState";

const updateLikedStatus = (feedItems, id, liked) => {
    const newItems = [];

    for (let i = 0; i < feedItems.length; i++) {
        const feedItem = feedItems[i];
        console.log(id, feedItem);
        if (feedItem.id === id) {
            if (liked) {
                feedItem.hasLiked = true;
                feedItem.likes += 1;
            } else {
                feedItem.hasLiked = false;
                feedItem.likes -= 1;
            }
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
                items: updateLikedStatus(state.items, action.id, true),
                loading: false,
                error: null,
            };
        case postTypes.LIKE_POST_ERROR:
            return {
                ...state,
                error: action.error,
                loading: false,
            };
        case postTypes.UNLIKE_POST:
            return {
                ...state,
                loading: true,
            };
        case postTypes.UNLIKE_POST_SUCCESS:
            return {
                ...state,
                items: updateLikedStatus(state.items, action.id, false),
                loading: false,
                error: null,
            };
        case postTypes.UNLIKE_POST_ERROR:
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
