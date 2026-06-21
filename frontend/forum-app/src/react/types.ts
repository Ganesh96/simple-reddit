export type PageInfo = {
  limit: number;
  has_more: boolean;
  next_cursor: string;
};

export type ForumPost = {
  id: string;
  title: string;
  text: string;
  community: string;
  username: string;
  upVotes: number;
  downVotes: number;
  commentsCount: number;
  createdAt: string;
  updatedAt: string;
};

export type ForumComment = {
  id: string;
  postId: string;
  text: string;
  username: string;
  upVotes: number;
  downVotes: number;
  edited: boolean;
  createdAt: string;
  updatedAt: string;
};

export type CommunitySummary = {
  id: string;
  name: string;
  description: string;
  membersCount: number;
  postsCount: number;
};
