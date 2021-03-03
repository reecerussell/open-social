import { get, post } from "../utils/api";
import { profileActions, userActions } from "../actions";

const fetchProfile = username => dispatch => {
  dispatch(profileActions.loadProfile());

  return get("profile/" + username)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      dispatch(profileActions.loadProfileSuccess(res.data));
    })
    .catch(err => dispatch(profileActions.loadProfileError(err.toString())));
};

const fetchInfo = () => dispatch => {
  dispatch(userActions.loadInfo());

  return get("me")
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      dispatch(userActions.loadInfoSuccess(res.data));
    })
    .catch(err => dispatch(userActions.loadInfoError(err.toString())));
};

const submitFollow = userId => dispatch => {
  dispatch(profileActions.followProfile(userId));

  return post("users/follow/" + userId)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      dispatch(profileActions.followProfileSuccess(userId));
    })
    .catch(err => dispatch(profileActions.followProfileError(err.toString())));
};

const submitUnfollow = userId => dispatch => {
  dispatch(profileActions.unfollowProfile(userId));

  return post("users/unfollow/" + userId)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      dispatch(profileActions.unfollowProfileSuccess(userId));
    })
    .catch(err =>
      dispatch(profileActions.unfollowProfileError(err.toString()))
    );
};

export { fetchProfile, fetchInfo, submitFollow, submitUnfollow };
