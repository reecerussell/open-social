import React, { useEffect } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { feedApi, postApi } from "../../api";
import { Image } from "../shared";

const Feed = ({ items, error, loading, fetchFeed, likePost }) => {
    useEffect(() => {
        fetchFeed();
    }, [fetchFeed]);

    const handleLikePost = post => e => {
        e.preventDefault();

        if (!post.hasUserLiked) {
            likePost(post.id);
        }
    };

    return items.map((item, key) => (
        <div className="section mb-4" key={key}>
            {item.mediaId && (
                <Image
                    id={item.mediaId}
                    className="img-fluid"
                    onDoubleClick={handleLikePost(item)}
                />
            )}

            <div className="p-4">
                <div className="post-user-info">
                    <p>
                        <a href="/" className="text-secondary">
                            <b>
                                {item.likes === 1
                                    ? item.likes + " Like"
                                    : item.likes + " Likes"}
                            </b>
                        </a>{" "}
                        <a href="/" className="text-secondary">
                            <b>0 Comments</b>
                        </a>
                        <br />
                        <b>
                            <a href="/" className="text-secondary">
                                {item.username}
                            </a>
                        </b>{" "}
                        <a href="/" className="text-secondary">
                            {item.caption}
                        </a>
                        <br />
                        <small>
                            <a href="/" className="text-muted">
                                {item.posted}
                            </a>
                        </small>
                    </p>
                </div>
            </div>
        </div>
    ));
};

Feed.propTypes = {
    items: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.string.isRequired,
            caption: PropTypes.string.isRequired,
            posted: PropTypes.string.isRequired,
            username: PropTypes.string.isRequired,
            likes: PropTypes.number.isRequired,
            hasUserLiked: PropTypes.bool.isRequired,
        }).isRequired
    ),
    error: PropTypes.string,
    loading: PropTypes.bool.isRequired,
    fetchFeed: PropTypes.func.isRequired,
    likePost: PropTypes.func.isRequired,
};

Feed.defaultProps = {
    items: [],
    error: null,
};

const mapStateToProps = state => ({
    items: state.feed.items,
    error: state.feed.error,
    loading: state.feed.loading,
});

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            fetchFeed: feedApi.fetchFeed,
            likePost: postApi.likePost,
        },
        dispatch
    );

export default connect(mapStateToProps, mapDispatchToProps)(Feed);
