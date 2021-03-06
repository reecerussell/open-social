# SQL

The open-social platform is powered by SQL Server with database migrations provided by [reecerussell/migrations](https://github.com/reecerussell/migrations).

## Migrations

Migrations are configured in [migrations.yaml](migrations.yaml).

| Name                     | Description                                                                                       |
| ------------------------ | ------------------------------------------------------------------------------------------------- |
| InitialCreation          | Initial creation of the database tables, including: Media, Users, UserFollowers, Posts, PostLikes |
| GetPostLikesFunction     | Creates the GetPostLikes SQL function.                                                            |
| HasUserLikedPostFunction | Creates the HasUserLikedPost SQL function.                                                        |