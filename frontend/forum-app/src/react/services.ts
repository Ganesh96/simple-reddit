import type { CommunitySummary, ForumComment, ForumPost, PageInfo } from './types';

type Envelope<T> = {
  status: number;
  message: string;
  code: string;
  data: T;
};

type RawRecord = Record<string, unknown>;

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  if (!headers.has('Content-Type') && init.body) headers.set('Content-Type', 'application/json');
  const response = await fetch(`${BASE_URL}${path}`, { ...init, headers });
  const payload = (await response.json().catch(() => null)) as Envelope<T> | null;
  if (!response.ok) throw new Error(payload?.message || 'Request failed');
  return (payload?.data ?? ({} as T)) as T;
}

function idOf(raw: RawRecord): string {
  const value = raw.id ?? raw.ID ?? raw._id;
  return typeof value === 'string' ? value : '';
}

function textOf(value: unknown): string {
  return typeof value === 'string' ? value : '';
}

function numberOf(value: unknown): number {
  return typeof value === 'number' ? value : 0;
}

function normalizePost(raw: RawRecord): ForumPost {
  return {
    id: idOf(raw),
    title: textOf(raw.title ?? raw.Title),
    text: textOf(raw.text ?? raw.Text ?? raw.body ?? raw.Body),
    community: textOf(raw.community ?? raw.Community),
    username: textOf(raw.username ?? raw.Username),
    upVotes: numberOf(raw.up_votes ?? raw.UpVotes),
    downVotes: numberOf(raw.down_votes ?? raw.DownVotes),
    commentsCount: numberOf(raw.comments_count ?? raw.CommentsCount),
    createdAt: textOf(raw.creation_date ?? raw.CreationDate),
    updatedAt: textOf(raw.updation_date ?? raw.UpdationDate),
  };
}

function normalizeForumComment(raw: RawRecord): ForumComment {
  return {
    id: idOf(raw),
    postId: textOf(raw.post_id ?? raw.PostID),
    text: textOf(raw.text ?? raw.Text ?? raw.body ?? raw.Body),
    username: textOf(raw.username ?? raw.Username),
    upVotes: numberOf(raw.up_votes ?? raw.UpVotes),
    downVotes: numberOf(raw.down_votes ?? raw.DownVotes),
    edited: Boolean(raw.edited ?? raw.Edited),
    createdAt: textOf(raw.creation_date ?? raw.CreationDate),
    updatedAt: textOf(raw.updation_date ?? raw.UpdationDate),
  };
}

function normalizeCommunity(raw: RawRecord): CommunitySummary {
  return {
    id: idOf(raw),
    name: textOf(raw.name ?? raw.Name),
    description: textOf(raw.description ?? raw.Description),
    membersCount: numberOf(raw.members_count ?? raw.MembersCount),
    postsCount: numberOf(raw.posts_count ?? raw.PostsCount),
  };
}

export async function getPosts(params: { after?: string; community?: string } = {}): Promise<{ posts: ForumPost[]; pagination: PageInfo }> {
  const search = new URLSearchParams({ limit: '20' });
  if (params.after) search.set('after', params.after);
  if (params.community) search.set('community', params.community);
  const data = await request<{ posts?: RawRecord[]; pagination?: PageInfo }>(`/posts?${search.toString()}`);
  return { posts: (data.posts ?? []).map(normalizePost), pagination: data.pagination ?? { limit: 20, has_more: false, next_cursor: '' } };
}

export async function getPost(postId: string): Promise<ForumPost> {
  const data = await request<{ post: RawRecord }>(`/posts/${postId}`);
  return normalizePost(data.post);
}

export async function getCommunities(): Promise<CommunitySummary[]> {
  const data = await request<{ communities?: RawRecord[] }>('/communities');
  return (data.communities ?? []).map(normalizeCommunity);
}

export async function getForumComments(postId: string, after = ''): Promise<{ comments: ForumComment[]; pagination: PageInfo }> {
  const search = new URLSearchParams({ limit: '20' });
  if (after) search.set('after', after);
  const data = await request<{ comments?: RawRecord[]; pagination?: PageInfo }>(`/posts/${postId}/comments?${search.toString()}`);
  return { comments: (data.comments ?? []).map(normalizeForumComment), pagination: data.pagination ?? { limit: 20, has_more: false, next_cursor: '' } };
}
