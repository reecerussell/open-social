import React from "react";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { Container, Row, Col, ListGroup, ListGroupItem } from "reactstrap";

const Layout = ({ children, userInfo }) => (
  <Container className="sections">
    <Row>
      <Col md="4" className="mb-4 mb-md-0">
        <div className="section" id="sidebar">
          <div className="p-4">
            <div className="post-user-info">
              <a href="/">
                <div
                  className="post-user-image"
                  style={{
                    backgroundImage:
                      "url(https://cdn.business2community.com/wp-content/uploads/2017/08/blank-profile-picture-973460_640.png)",
                  }}
                ></div>
              </a>
              <p>
                <Link to={"/u/" + userInfo.username} className="text-secondary">
                  <b>{userInfo.username}</b>
                </Link>
                <br />
                <span className="text-muted">
                  {userInfo.followerCount}
                  {userInfo.followerCount == 1 ? " Follower" : " Followers"}
                </span>
              </p>
            </div>
          </div>

          <ListGroup flush className="border-top">
            <ListGroupItem>
              <Link to="/" className="d-block text-secondary">
                <i className="fas fa-stream"></i> Feed
              </Link>
            </ListGroupItem>
            <ListGroupItem>
              <Link
                to={"/u/" + userInfo.username}
                className="d-block text-secondary"
              >
                <i className="fas fa-user-alt"></i> Profile
              </Link>
            </ListGroupItem>
            <ListGroupItem>
              <a href="/" className="d-block text-secondary">
                <i className="fas fa-sign-out-alt"></i> Logout
              </a>
            </ListGroupItem>
          </ListGroup>
        </div>
      </Col>
      <Col md="8">{children}</Col>
    </Row>
  </Container>
);

Layout.propTypes = {
  userInfo: PropTypes.shape({
    id: PropTypes.string.isRequired,
    username: PropTypes.string.isRequired,
    mediaId: PropTypes.string,
    followerCount: PropTypes.number.isRequired,
  }),
};

export default Layout;
