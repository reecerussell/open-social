import React, { useState, useRef } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import * as api from "../../api/post";
import { Textarea } from "../form";

const defaultState = {
    caption: "",
    file: null,
};
const defaultImageUploadText = "Upload An Image!";

const Create = ({ error, loading, submitPost }) => {
    const [post, setPost] = useState(defaultState);
    const [imageUploadText, setImageUploadText] = useState(
        defaultImageUploadText
    );

    const fileRef = useRef(null);

    const handleSubmit = e => {
        e.preventDefault();

        const formData = new FormData()
        formData.append("caption", post.caption)

        if (post.file) {
            formData.append("file", post.file, post.file.name)
        }

        submitPost(formData);
    };

    const handleUpdate = e => {
        const { name, value } = e.target;
        const newState = { ...post };
        newState[name] = value;
        setPost(newState);
    };

    const handleFileUpdate = e => {
        const { files } = e.target;

        if (files.length > 0) {
            setImageUploadText(files[0].name);
            setPost({...post, file: files[0]})
        } else {
            setImageUploadText(defaultImageUploadText);
            setPost({...post, file: null})
        }
    };

    const selectFileBtnClick = e => {
        e.preventDefault();
        fileRef.current.click();
    };

    return (
        <div className="section p-4" id="createPost">
            <h1 className="header-1">Create a Post</h1>
            <p className="text-center small info">
                Make a post to your Social Media account. Upload an image, just
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
                        <input
                            type="file"
                            className="d-none"
                            ref={fileRef}
                            onChange={handleFileUpdate}
                        />
                        <button
                            type="button"
                            className="btn btn-secondary"
                            onClick={selectFileBtnClick}
                        >
                            {imageUploadText}
                        </button>
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
