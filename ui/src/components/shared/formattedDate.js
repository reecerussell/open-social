import moment from "moment";
import PropTypes from "prop-types";

const FormattedDate = ({ value }) => {
    const date = new Date(value);
    return moment(date).format("MMMM Do YYYY, h:mm:ss a");
};

FormattedDate.propTypes = {
    value: PropTypes.string.isRequired,
};

export default FormattedDate;
