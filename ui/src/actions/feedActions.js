import * as types from "./feedActionTypes";

export const loadFeed = () => ({
  type: types.LOAD_FEED,
});

export const loadFeedSuccess = feedItems => ({
  type: types.LOAD_FEED_SUCCESS,
  feedItems: feedItems,
});

export const loadFeedError = error => ({
  type: types.LOAD_FEED_ERROR,
  error: error,
});
