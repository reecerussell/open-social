CREATE OR ALTER FUNCTION [dbo].[HasUserLikedPost] (@PostId INT, @UserId INT)
RETURNS BIT
BEGIN
	DECLARE @result BIT = 0

	IF EXISTS(SELECT [PostId] FROM [PostLikes] WHERE [PostId] = @PostId AND [UserId] = @UserId)
	BEGIN
		SET @result = 1
	END

	RETURN @result
END