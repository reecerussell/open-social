CREATE TABLE [dbo].[UserFollowers] (
	[UserId] INT NOT NULL,
	[FollowerId] INT NOT NULL,
	CONSTRAINT PK_UserFollowers PRIMARY KEY ([UserId], [FollowerId]),
	CONSTRAINT FK_UserFollowers_UserId FOREIGN KEY ([UserId]) REFERENCES [Users] ([Id]),
	CONSTRAINT FK_UserFollowers_FollowerId FOREIGN KEY ([FollowerId]) REFERENCES [Users] ([Id])
)