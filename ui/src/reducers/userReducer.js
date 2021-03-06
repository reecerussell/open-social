import * as types from "../actions/userActionTypes";
import initialState from "../store/initialState";

const userReducer = (state = initialState.user, action) => {
  switch (action.type) {
    case types.LOAD_INFO:
      return {
        ...state,
        loading: true,
      };
    case types.LOAD_INFO_SUCCESS:
      return {
        ...state,
        ...action.info,
        loading: false,
        error: null,
      };
    case types.LOAD_INFO_ERROR:
      return {
        ...state,
        loading: false,
        error: action.error,
      };
    default:
      return state;
  }
};

export default userReducer;
