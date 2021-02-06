import React from "react";
import LazyLoad from "react-lazyload";
import PropTypes from "prop-types";
import environment from "../../environment";

const Image = ({ id, className }) => (
    <LazyLoad once>
        <img src={environment.mediaUrl + id} className={className} />
    </LazyLoad>
);

Image.propTypes = {
    id: PropTypes.string.isRequired,
    className: PropTypes.string,
};

Image.defaultProps = {
    className: null,
};

export default Image;
