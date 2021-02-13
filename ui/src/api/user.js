import { get } from "../utils/api";
import { profileActions } from "../actions";

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

export { fetchProfile };
