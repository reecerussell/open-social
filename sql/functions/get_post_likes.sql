CREATE OR ALTER FUNCTION [dbo].[GetPostLikes] (@PostId INT)
RETURNS INT
BEGIN
	DECLARE @count INT

	SELECT @count=COUNT([PostId]) 
	FROM [PostLikes] 
	WHERE [PostId] = @PostId

	RETURN @count
END