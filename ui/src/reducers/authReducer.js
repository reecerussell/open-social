import { authTypes } from "../actions";
import initialState from "../store/initialState";

const authReducer = (state = initialState.auth, action) => {
  switch (action.type) {
    case authTypes.REGISTER:
    case authTypes.LOGIN:
      return {
        ...state,
        loading: true,
      };
    case authTypes.REGISTER_SUCCESS:
    case authTypes.LOGIN_SUCCESS:
      return {
        ...state,
        loading: false,
        error: null,
        success: true,
      };
    case authTypes.REGISTER_ERROR:
    case authTypes.LOGIN_ERROR:
      return {
        ...state,
        loading: false,
        error: action.error,
        success: false,
      };
    default:
      return state;
  }
};

export default authReducer;
