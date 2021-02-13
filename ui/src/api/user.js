import { get } from "../utils/api";
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
        .catch(err =>
            dispatch(profileActions.loadProfileError(err.toString()))
        );
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

export { fetchProfile, fetchInfo };
