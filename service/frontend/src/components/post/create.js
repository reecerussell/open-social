import React, { useState } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import * as api from "../../api/post";
import { Textarea } from "../form";

const defaultState = {
    caption: "",
};

const Create = ({ error, loading, submitPost }) => {
    const [post, setPost] = useState(defaultState);

    const handleSubmit = e => {
        e.preventDefault();

        submitPost(post);
    };

    const handleUpdate = e => {
        const { name, value } = e.target;
        const newState = { ...post };
        newState[name] = value;
        setPost(newState);
    };

    return (
        <div className="section p-4" id="createPost">
            <h1 className="header-1">Create a Post</h1>
            <p className="text-center small info">
                Make a post to your Social Media account. Upload and image, just
                a caption, or both!
            </p>
            <hr />
            <form onSubmit={handleSubmit}>
                {error && <p className="error">{error}</p>}
                <Textarea
                    name="caption"
                    label="Caption"
                    placeholder="Write a caption!"
                    value={post.caption}
                    onChange={handleUpdate}
                    max="255"
                    required
                />
                <div className="form-group">
                    <div className="button-group float-right">
                        <button
                            type="submit"
                            className="btn btn-primary"
                            disabled={loading}
                        >
                            Submit
                        </button>
                    </div>
                </div>
            </form>
        </div>
    );
};

Create.propTypes = {
    error: PropTypes.string,
    loading: PropTypes.bool.isRequired,
    submitPost: PropTypes.func.isRequired,
};

Create.defaultProps = {
    error: null,
};

const mapStateToProps = state => ({
    error: state.post.error,
    loading: state.post.loading,
});

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            submitPost: api.submitPost,
        },
        dispatch
    );

export default connect(mapStateToProps, mapDispatchToProps)(Create);
