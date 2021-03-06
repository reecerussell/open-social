import React, { useEffect } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { useParams } from "react-router-dom";
import { userApi } from "../../api";
import environment from "../../environment";
import { Feed } from "../shared";

const ProfileImage = ({ id }) => {
  const url = environment.mediaUrl + id;
  const style = {
    backgroundImage: `url(${url})`,
  };

  return (
    <div className="col-sm-3">
      <div className="mb-3" id="profile-image" style={style} />
    </div>
  );
};

const Profile = ({
  profile,
  error,
  loading,
  fetchProfile,
  followProfile,
  unfollowProfile,
}) => {
  const { username } = useParams();

  useEffect(() => {
    fetchProfile(username);
  }, [fetchProfile, username]);

  const handleFollowChange = e => {
    e.preventDefault();

    if (loading) {
      return;
    }

    if (profile.isFollowing) {
      unfollowProfile(profile.userId);
    } else {
      followProfile(profile.userId);
    }
  };

  return (
    <>
      <div className="section p-4 mb-4" id="profile">
        <div className="row">
          {profile.mediaId && <ProfileImage id={profile.mediaId} />}

          <div className="col-sm-9">
            <h1 className="header-1 text-left">{profile.username}</h1>
            <p>
              <b>
                {profile.followerCount}
                {profile.followerCount === 1 ? " Follower" : " Followers"}
              </b>

              <span className="ml-4">
                {profile.postCount}
                {profile.postCount === 1 ? " Post" : " Posts"}
              </span>

              {profile.isOwner && (
                <span className="ml-4">
                  <a href="/" className="text-secondary underline">
                    Manage Your Account
                  </a>
                </span>
              )}

              {!profile.isOwner && (
                <span
                  className="ml-4 underline follow-link-btn"
                  onClick={handleFollowChange}
                >
                  {!profile.isFollowing ? "Follow" : "Unfollow"}
                </span>
              )}
            </p>

            {profile.bio && <p>{profile.bio}</p>}
          </div>
        </div>
      </div>

      <Feed items={profile.feed} />
    </>
  );
};

Profile.propTypes = {
  profile: PropTypes.shape({
    userId: PropTypes.string.isRequired,
    username: PropTypes.string.isRequired,
    mediaId: PropTypes.string,
    bio: PropTypes.string,
    followerCount: PropTypes.number.isRequired,
    isFollowing: PropTypes.bool.isRequired,
    isOwner: PropTypes.bool.isRequired,
    postCount: PropTypes.number.isRequired,
    feed: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.string.isRequired,
        caption: PropTypes.string.isRequired,
        posted: PropTypes.string.isRequired,
        username: PropTypes.string.isRequired,
        likes: PropTypes.number.isRequired,
        hasLiked: PropTypes.bool.isRequired,
      })
    ),
  }).isRequired,
  error: PropTypes.string,
  loading: PropTypes.bool,
  fetchProfile: PropTypes.func.isRequired,
  followProfile: PropTypes.func.isRequired,
  unfollowProfile: PropTypes.func.isRequired,
};

PropTypes.defaultProps = {
  profile: {
    userId: "",
    username: "",
    mediaId: null,
    bio: null,
    followerCount: 0,
    isFollowing: false,
    isOwner: false,
    postCount: 0,
    feed: [],
  },
  error: null,
  loading: false,
};

const mapStateToProps = state => ({
  profile: state.profile.data,
  error: state.profile.error,
  loading: state.profile.loading,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      fetchProfile: userApi.fetchProfile,
      followProfile: userApi.submitFollow,
      unfollowProfile: userApi.submitUnfollow,
    },
    dispatch
  );

export default connect(mapStateToProps, mapDispatchToProps)(Profile);
