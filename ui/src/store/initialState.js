const initialState = {
    feed: {
        items: [],
        loading: false,
        error: null,
    },
    post: {
        loading: false,
        error: null,
    },
    profile: {
        data: {
            username: "",
            mediaId: null,
            bio: null,
            followerCount: 0,
            isFollowing: false,
            isOwner: false,
            postCount: 0,
            feed: [],
        },
        loading: false,
        error: null,
    },
};

export default initialState;
