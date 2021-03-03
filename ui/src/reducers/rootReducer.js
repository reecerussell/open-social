import { combineReducers } from "redux";
import feedReducer from "./feedReducer";
import postsReducer from "./postsReducer";
import profileReducer from "./profileReducer";
import userReducer from "./userReducer";
import authReducer from "./authReducer";

const rootReducer = (state, action) =>
  combineReducers({
    feed: feedReducer,
    posts: postsReducer,
    profile: profileReducer,
    user: userReducer,
    auth: authReducer,
  })(state, action);

export default rootReducer;
