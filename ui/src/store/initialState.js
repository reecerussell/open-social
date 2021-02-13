const initialState = {
    feed: {
        items: [],
        loading: false,
        error: null,
    },
    posts: {
        posts: [],
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
    user: {
        id: "",
        username: "",
        followerCount: 0,
        mediaId: null,
        error: null,
        loading: false,
    },
};

export default initialState;
