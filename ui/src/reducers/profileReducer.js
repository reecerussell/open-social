import { postTypes, profileTypes } from "../actions";
import initialState from "../store/initialState";

const updatePostLikedStatus = (feedItems, id, liked) => {
    const newItems = [];

    for (let i = 0; i < feedItems.length; i++) {
        const feedItem = feedItems[i];
        if (feedItem.id === id) {
            if (liked) {
                feedItem.hasUserLiked = true;
                feedItem.likes += 1;
            } else {
                feedItem.hasUserLiked = false;
                feedItem.likes -= 1;
            }
        }

        newItems.push(feedItem);
    }

    return newItems;
};

const profileReducer = (state = initialState.post, action) => {
    switch (action.type) {
        case profileTypes.LOAD_PROFILE:
            return {
                ...state,
                loading: false,
            };
        case profileTypes.LOAD_PROFILE_SUCCESS:
            return {
                ...state,
                data: action.profile,
                loading: false,
                error: null,
            };
        case profileTypes.LOAD_PROFILE_ERROR:
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
                data: {
                    ...state.data,
                    feed: updatePostLikedStatus(
                        state.data.feed,
                        action.id,
                        true
                    ),
                },
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
                data: {
                    ...state.data,
                    feed: updatePostLikedStatus(
                        state.data.feed,
                        action.id,
                        false
                    ),
                },
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

export default profileReducer;
