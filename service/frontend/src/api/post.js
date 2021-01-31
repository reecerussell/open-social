import { post } from "../utils/api";
import * as actions from "../actions/postActions";

export const submitPost = data => dispatch => {
    dispatch(actions.createPost());

    return post("posts", data)
        .then(res => {
            if (!res.ok) {
                throw new Error(res.error);
            }

            dispatch(actions.createPostSuccess());
        })
        .catch(err => dispatch(actions.createPostError(err.toString())));
};
