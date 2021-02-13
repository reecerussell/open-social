import { combineReducers } from "redux";
import feedReducer from "./feedReducer";
import postsReducer from "./postsReducer";
import profileReducer from "./profileReducer";

const rootReducer = (state, action) =>
    combineReducers({
        feed: feedReducer,
        posts: postsReducer,
        profile: profileReducer,
    })(state, action);

export default rootReducer;
