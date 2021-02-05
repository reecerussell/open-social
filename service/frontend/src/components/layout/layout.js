import React from "react";
import { Container, Row, Col, ListGroup, ListGroupItem } from "reactstrap";

const Layout = ({ children }) => (
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
								<a className="text-secondary" href="/">
									<b>test-username</b>
								</a>
								<br />
								<span className="text-muted">6 Followers</span>
							</p>
						</div>
					</div>

					<ListGroup flush className="border-top">
						<ListGroupItem>
							<a href="/" className="d-block text-secondary">
								<i className="fas fa-stream"></i> Feed
							</a>
						</ListGroupItem>
						<ListGroupItem>
							<a href="/" className="d-block text-secondary">
								<i className="fas fa-user-alt"></i> Profile
							</a>
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

export default Layout;
