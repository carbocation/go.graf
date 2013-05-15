== Recursive Points ==
Sorting is done based on all points belonging to the post,
as well as all children, children's children, ..., of the 
post.

Tradeoffs include the fact that to list the forum root, 
all posts have to be pulled down (or alternatively, 
the DepthOne... functions could be modified to calculate 
points based on the closure tables, which should work fine).

To convert to a more pure SQL-based approach, GROUP BY ancestor, 
where ancestor is each post.
