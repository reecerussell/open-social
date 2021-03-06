import { post, postForm, get } from "../utils/api";
import { postActions } from "../actions";

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

  return post("posts/like/" + id, null)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      dispatch(postActions.likePostSuccess(id));
    })
    .catch(err => dispatch(postActions.likePostError(err.toString())));
};

export const unlikePost = id => dispatch => {
  dispatch(postActions.unlikePost());

  return post("posts/unlike/" + id, null)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      dispatch(postActions.unlikePostSuccess(id));
    })
    .catch(err => dispatch(postActions.unlikePostError(err.toString())));
};

export const loadPost = id => dispatch => {
  dispatch(postActions.loadPost());

  return get("posts/" + id)
    .then(res => {
      if (!res.ok) {
        throw new Error(res.error);
      }

      dispatch(postActions.loadPostSuccess(res.data));
    })
    .catch(err => dispatch(postActions.loadPostError(err.toString())));
};
