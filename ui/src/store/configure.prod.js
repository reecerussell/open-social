import { createStore, compose, applyMiddleware } from "redux";
import thunk from "redux-thunk";
import rootReducer from "../reducers";

const configureStore = initialState => {
  return createStore(rootReducer, initialState, compose(applyMiddleware(thunk)));
};

export default configureStore;
