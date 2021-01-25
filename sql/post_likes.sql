CREATE TABLE [dbo].[PostLikes] (
	[PostId] INT NOT NULL,
	[UserId] INT NOT NULL,
	CONSTRAINT PK_PostLikes PRIMARY KEY ([PostId], [UserId]),
	CONSTRAINT FK_PostLikes_PostId FOREIGN KEY ([PostId]) REFERENCES [Posts] ([Id]),
	CONSTRAINT FK_PostLikes_UserId FOREIGN KEY ([UserId]) REFERENCES [Users] ([Id])
)