import React from "react";
import classNames from "classnames";
import PropTypes from "prop-types";

const Textbox = ({
  label,
  value,
  onChange,
  name,
  required,
  placeholder,
  max,
  text,
}) => {
  const isEmpty = !value;
  const classes = classNames("form-group", "sm-input", {
    empty: isEmpty,
  });

  return (
    <div className={classes}>
      <label htmlFor={name} className="small">
        {label}
      </label>
      <input
        name={name}
        id={name}
        className="form-control"
        type={text}
        placeholder={placeholder}
        max={max}
        required={required}
        value={value}
        onChange={onChange}
      />
    </div>
  );
};

Textbox.propTypes = {
  label: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  name: PropTypes.string.isRequired,
  required: PropTypes.bool,
  placeholder: PropTypes.string,
  max: PropTypes.string,
  type: PropTypes.string,
};

Textbox.defaultProps = {
  required: false,
  placeholder: null,
  max: null,
  type: "text",
};

export default Textbox;
