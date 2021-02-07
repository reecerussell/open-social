import React from "react";
import classNames from "classnames";
import PropTypes from "prop-types";

const Textarea = ({
    label,
    value,
    onChange,
    name,
    required,
    rows,
    placeholder,
    max,
}) => {
    const isEmpty = !value;
    const classes = classNames("form-control", {
        empty: isEmpty,
    });

    return (
        <div className="form-group">
            <label htmlFor={name} className="small">
                {label}
            </label>
            <textarea
                name={name}
                id={name}
                className={classes}
                rows={rows}
                placeholder={placeholder}
                max={max}
                required={required}
                value={value}
                onChange={onChange}
            ></textarea>
        </div>
    );
};

Textarea.propTypes = {
    label: PropTypes.string.isRequired,
    value: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    name: PropTypes.string.isRequired,
    required: PropTypes.bool,
    rows: PropTypes.string,
    placeholder: PropTypes.string,
    max: PropTypes.string,
};

Textarea.defaultProps = {
    required: false,
    rows: "3",
    placeholder: null,
    max: null,
};

export default Textarea;
