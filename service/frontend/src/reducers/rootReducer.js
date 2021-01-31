import { combineReducers } from "redux";
import feedReducer from "./feedReducer";
import postReducer from "./postReducer";

const rootReducer = (state, action) =>
    combineReducers({
        feed: feedReducer,
        post: postReducer,
    })(state, action);

export default rootReducer;
