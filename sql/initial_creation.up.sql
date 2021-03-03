CREATE TABLE [dbo].[Media] (
	[Id] INT NOT NULL PRIMARY KEY IDENTITY(1,1),
	[ReferenceId] UNIQUEIDENTIFIER NOT NULL UNIQUE,
	[ContentType] VARCHAR(50) NOT NULL
);

CREATE TABLE [dbo].[Users] (
	[Id] INT NOT NULL IDENTITY(1,1) PRIMARY KEY,
	[ReferenceId] UNIQUEIDENTIFIER NOT NULL,
	[MediaId] INT NULL,
	[Username] VARCHAR(20) NOT NULL UNIQUE,
	[PasswordHash] VARCHAR(MAX) NOT NULL,
	[Bio] VARCHAR(255) NOT NULL,
    CONSTRAINT FK_Users_MediaId FOREIGN KEY ([MediaId]) REFERENCES [Users] ([Id])
);

CREATE TABLE [dbo].[UserFollowers] (
	[UserId] INT NOT NULL,
	[FollowerId] INT NOT NULL,
	CONSTRAINT PK_UserFollowers PRIMARY KEY ([UserId], [FollowerId]),
	CONSTRAINT FK_UserFollowers_UserId FOREIGN KEY ([UserId]) REFERENCES [Users] ([Id]),
	CONSTRAINT FK_UserFollowers_FollowerId FOREIGN KEY ([FollowerId]) REFERENCES [Users] ([Id])
);

CREATE TABLE [dbo].[Posts] (
	[Id] INT NOT NULL PRIMARY KEY IDENTITY(1,1),
	[ReferenceId] UNIQUEIDENTIFIER NOT NULL UNIQUE,
	[UserId] INT NOT NULL,
	[MediaId] INT NULL,
	[Posted] DATETIME NOT NULL,
	[Caption] NVARCHAR(255) NOT NULL,
    CONSTRAINT FK_Posts_UserId FOREIGN KEY ([UserId]) REFERENCES [Users] ([Id]),
    CONSTRAINT FK_Posts_MediaId FOREIGN KEY ([MediaId]) REFERENCES [Media] ([Id]),
);

CREATE TABLE [dbo].[PostLikes] (
	[PostId] INT NOT NULL,
	[UserId] INT NOT NULL,
	CONSTRAINT PK_PostLikes PRIMARY KEY ([PostId], [UserId]),
	CONSTRAINT FK_PostLikes_PostId FOREIGN KEY ([PostId]) REFERENCES [Posts] ([Id]),
	CONSTRAINT FK_PostLikes_UserId FOREIGN KEY ([UserId]) REFERENCES [Users] ([Id])
);
