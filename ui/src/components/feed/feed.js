import React, { useEffect } from "react";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { feedApi, postApi } from "../../api";
import { Image, FormattedDate } from "../shared";

const Feed = ({ items, error, loading, fetchFeed, likePost, unlikePost }) => {
    useEffect(() => {
        fetchFeed();
    }, [fetchFeed]);

    const handleLikePost = post => e => {
        e.preventDefault();

        if (post.hasUserLiked) {
            unlikePost(post.id);
        } else {
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
                        <Link
                            to={"/post/" + item.id}
                            className="text-secondary"
                        >
                            <b>
                                {item.likes === 1
                                    ? item.likes + " Like"
                                    : item.likes + " Likes"}
                            </b>
                        </Link>{" "}
                        <Link
                            to={"/post/" + item.id}
                            className="text-secondary"
                        >
                            <b>0 Comments</b>
                        </Link>
                        <br />
                        <b>
                            <a href="/" className="text-secondary">
                                {item.username}
                            </a>
                        </b>{" "}
                        <Link
                            to={"/post/" + item.id}
                            className="text-secondary"
                        >
                            {item.caption}
                        </Link>
                        <br />
                        <small>
                            <Link
                                to={"/post/" + item.id}
                                className="text-muted"
                            >
                                <FormattedDate value={item.posted} />
                            </Link>
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
    unlikePost: PropTypes.func.isRequired,
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
            unlikePost: postApi.unlikePost,
        },
        dispatch
    );

export default connect(mapStateToProps, mapDispatchToProps)(Feed);
