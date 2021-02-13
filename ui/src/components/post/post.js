import React, { useEffect } from "react";
import { useParams } from "react-router-dom";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { Image, FormattedDate } from "../shared";
import { postApi } from "../../api";

const Post = ({ post, loading, error, loadPost, likePost, unlikePost }) => {
    const { id } = useParams();

    useEffect(() => {
        loadPost(id);
    }, [loadPost, id]);

    const handleLikeClick = () => {
        if (loading) {
            return;
        }

        if (post.hasLiked) {
            unlikePost(post.id);
        } else {
            likePost(post.id);
        }
    };

    return (
        <div className="section" id="post">
            {post.mediaId && (
                <Image
                    id={post.mediaId}
                    alt={post.caption}
                    className="img-fluid"
                    onDoubleClick={handleLikeClick}
                />
            )}
            <div className="p-4">
                <div className="post-user-info">
                    <a href="/">
                        <div
                            className="post-user-image"
                            style={{
                                backgroundImage:
                                    "url(https://cdn.business2community.com/wp-content/uploads/2017/08/blank-profile-picture-973460_640.png)",
                            }}
                        ></div>
                    </a>
                    <p>
                        {post.likes === 1
                            ? post.likes + " Like"
                            : post.likes + " Likes"}
                        <br />
                        <b>
                            <a href="/" className="text-secondary">
                                {post.username}
                            </a>
                        </b>{" "}
                        {post.caption}
                        <br />
                        <small>
                            <FormattedDate value={post.posted} />{" "}
                            <span
                                className="like-link-btn"
                                onClick={handleLikeClick}
                            >
                                {post.hasLiked ? "Unlike" : "Like"}
                            </span>
                        </small>
                    </p>
                </div>
            </div>
        </div>
    );
};

Post.propTypes = {
    post: PropTypes.shape({
        id: PropTypes.string.isRequired,
        mediaId: PropTypes.string,
        posted: PropTypes.string.isRequired,
        username: PropTypes.string.isRequired,
        caption: PropTypes.string.isRequired,
        likes: PropTypes.number.isRequired,
        hasLiked: PropTypes.bool.isRequired,
        unlikePost: PropTypes.bool.isRequired,
    }).isRequired,
    error: PropTypes.string,
    loading: PropTypes.bool,
    loadPost: PropTypes.func.isRequired,
    likePost: PropTypes.func.isRequired,
};

Post.defaultProps = {
    error: null,
    loading: false,
};

const mapStateToProps = state => ({
    post: state.post,
    error: state.post.error,
    loading: state.post.loading,
});

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadPost: postApi.loadPost,
            likePost: postApi.likePost,
            unlikePost: postApi.unlikePost,
        },
        dispatch
    );

export default connect(mapStateToProps, mapDispatchToProps)(Post);
