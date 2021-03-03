import { createStore } from "redux";
import rootReducer from "../reducers";

const configureStore = initialState => {
  return createStore(rootReducer, initialState, null);
};

export default configureStore;
