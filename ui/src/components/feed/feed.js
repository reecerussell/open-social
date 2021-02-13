import React, { useEffect } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { feedApi } from "../../api";
import Shared from "../shared";

const Feed = ({ items, error, loading, fetchFeed }) => {
    useEffect(() => {
        fetchFeed();
    }, [fetchFeed]);

    return <Shared.Feed items={items} />;
};

Feed.propTypes = {
    items: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.string.isRequired,
            caption: PropTypes.string.isRequired,
            posted: PropTypes.string.isRequired,
            username: PropTypes.string.isRequired,
            likes: PropTypes.number.isRequired,
            hasLiked: PropTypes.bool.isRequired,
        }).isRequired
    ),
    error: PropTypes.string,
    loading: PropTypes.bool.isRequired,
    fetchFeed: PropTypes.func.isRequired,
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
        },
        dispatch
    );

export default connect(mapStateToProps, mapDispatchToProps)(Feed);
