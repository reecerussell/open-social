import React from "react";
import LazyLoad from "react-lazyload";
import PropTypes from "prop-types";
import environment from "../../environment";

const Image = ({ id, className, onDoubleClick }) => (
    <LazyLoad once>
        <img
            src={environment.mediaUrl + id}
            className={className}
            onDoubleClick={onDoubleClick}
        />
    </LazyLoad>
);

Image.propTypes = {
    id: PropTypes.string.isRequired,
    className: PropTypes.string,
    onDoubleClick: PropTypes.func,
};

Image.defaultProps = {
    className: null,
    onDoubleClick: null,
};

export default Image;
