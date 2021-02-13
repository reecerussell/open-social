import { combineReducers } from "redux";
import feedReducer from "./feedReducer";
import postReducer from "./postReducer";
import profileReducer from "./profileReducer";

const rootReducer = (state, action) =>
    combineReducers({
        feed: feedReducer,
        post: postReducer,
        profile: profileReducer,
    })(state, action);

export default rootReducer;
