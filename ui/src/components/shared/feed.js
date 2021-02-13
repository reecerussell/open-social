import React from "react";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { postApi } from "../../api";
import { Image, FormattedDate } from "../shared";

const Feed = ({ items, likePost, unlikePost }) => {
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
                    alt={item.caption}
                    onDoubleClick={handleLikePost(item)}
                />
            )}

            <div className="p-4">
                <div className="post-user-info">
                    <p>
                        <Link to={"/p/" + item.id} className="text-secondary">
                            <b>
                                {item.likes === 1
                                    ? item.likes + " Like"
                                    : item.likes + " Likes"}
                            </b>
                        </Link>{" "}
                        <Link to={"/p/" + item.id} className="text-secondary">
                            <b>0 Comments</b>
                        </Link>
                        <br />
                        <b>
                            <Link
                                to={"/u/" + item.username}
                                className="text-secondary"
                            >
                                {item.username}
                            </Link>
                        </b>{" "}
                        <Link to={"/p/" + item.id} className="text-secondary">
                            {item.caption}
                        </Link>
                        <br />
                        <small>
                            <Link to={"/p/" + item.id} className="text-muted">
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
    likePost: PropTypes.func.isRequired,
    unlikePost: PropTypes.func.isRequired,
};

Feed.defaultProps = {
    items: [],
};

const mapStateToProps = state => ({});

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            likePost: postApi.likePost,
            unlikePost: postApi.unlikePost,
        },
        dispatch
    );

export default connect(mapStateToProps, mapDispatchToProps)(Feed);
