import { post, postForm } from "../utils/api";
import { postActions, feedActions } from "../actions";

export const submitPost = data => dispatch => {
    dispatch(postActions.createPost());

    return postForm("posts", data)
        .then(res => {
            if (!res.ok) {
                throw new Error(res.error);
            }

            dispatch(postActions.createPostSuccess());
        })
        .catch(err => dispatch(postActions.createPostError(err.toString())));
};

export const likePost = id => dispatch => {
    dispatch(postActions.likePost());
    //dispatch(feedActions.likeFeedPost());

    return post("posts/like/" + id, null)
        .then(res => {
            if (!res.ok) {
                throw new Error(res.error);
            }

            dispatch(postActions.likePostSuccess(id));
            //dispatch(feedActions.likeFeedPostSuccess(id));
        })
        .catch(err => {
            dispatch(postActions.likePostError(err.toString()));
            //dispatch(feedActions.likeFeedPostError(err.toString()));
        });
};
