import { get } from "../utils/api";
import * as actions from "../actions/feedActions";

export const fetchFeed = () => dispatch => {
    dispatch(actions.loadFeed());

    return get("feed")
        .then(res => {
            if (!res.ok) {
                throw new Error(res.error);
            }

            dispatch(actions.loadFeedSuccess(res.data ?? []));
        })
        .catch(err => dispatch(actions.loadFeedError(err.toString())));
};
