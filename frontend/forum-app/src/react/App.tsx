import { FormEvent, useEffect, useState } from 'react';
import { Link, NavLink, Route, Routes, useParams } from 'react-router-dom';
import {
  getCommunities,
  getForumComments,
  getPost,
  getPosts,
} from './services';
import type { CommunitySummary, ForumComment, ForumPost, PageInfo } from './types';

const emptyPage: PageInfo = { limit: 20, has_more: false, next_cursor: '' };

export default function App() {
  return (
    <div className="app-shell">
      <Header />
      <main className="main-content">
        <Routes>
          <Route path="/" element={<FeedPage />} />
          <Route path="/home" element={<FeedPage />} />
          <Route path="/communities" element={<CommunitiesPage />} />
          <Route path="/subreddits" element={<CommunitiesPage />} />
          <Route path="/posts/new" element={<ComposePlaceholder />} />
          <Route path="/newpostform" element={<ComposePlaceholder />} />
          <Route path="/post/:postId" element={<PostDetailPage />} />
          <Route path="/login" element={<AuthPlaceholder />} />
          <Route path="/profile" element={<ProfilePlaceholder />} />
          <Route path="/termsandconditions" element={<StaticPage title="Terms" />} />
          <Route path="/privacypolicy" element={<StaticPage title="Privacy" />} />
          <Route path="/contentpolicy" element={<StaticPage title="Content Policy" />} />
          <Route path="/modpolicy" element={<StaticPage title="Moderation Policy" />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
    </div>
  );
}

function Header() {
  return (
    <header className="topbar">
      <Link className="brand" to="/">
        <span className="brand-mark">sr</span>
        <span>Simple Reddit</span>
      </Link>
      <nav className="nav-links" aria-label="Primary navigation">
        <NavLink to="/home">Feed</NavLink>
        <NavLink to="/communities">Communities</NavLink>
        <NavLink to="/posts/new">New post</NavLink>
        <NavLink to="/login">Login</NavLink>
      </nav>
    </header>
  );
}

function FeedPage() {
  const [posts, setPosts] = useState<ForumPost[]>([]);
  const [page, setPage] = useState<PageInfo>(emptyPage);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  async function loadFeed(after = '') {
    setLoading(true);
    setError('');
    try {
      const result = await getPosts(after ? { after } : {});
      setPosts((current) => (after ? [...current, ...result.posts] : result.posts));
      setPage(result.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unable to load posts');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadFeed();
  }, []);

  return (
    <section className="page-grid">
      <div className="stack">
        <PageHeading eyebrow="Feed" title="Latest community posts" description="A faster React feed focused on scanning, reading, and moving to comments quickly." />
        {error && <Notice tone="error">{error}</Notice>}
        {loading && posts.length === 0 ? <SkeletonList /> : null}
        {!loading && posts.length === 0 ? <EmptyState title="No posts yet" body="Create a post after auth write flows are re-enabled." /> : null}
        {posts.map((post) => <PostCard key={post.id} post={post} />)}
        {page.has_more ? (
          <button className="secondary-action" onClick={() => void loadFeed(page.next_cursor)} disabled={loading}>
            {loading ? 'Loading...' : 'Load more'}
          </button>
        ) : null}
      </div>
      <aside className="side-panel">
        <h2>Product direction</h2>
        <p>Minimize choices. Keep the feed readable. Push creation and voting behind clear authenticated actions.</p>
        <Link className="primary-action full" to="/communities">Browse communities</Link>
      </aside>
    </section>
  );
}

function PostCard({ post }: { post: ForumPost }) {
  const score = post.upVotes - post.downVotes;
  return (
    <article className="card post-card">
      <div className="vote-rail" aria-label="Post score">
        <span>{score}</span>
        <small>score</small>
      </div>
      <div className="card-body">
        <p className="meta">Posted by {post.username || 'unknown'} {post.createdAt ? `on ${formatDate(post.createdAt)}` : ''}</p>
        <Link className="post-title" to={`/post/${post.id}`}>{post.title || 'Untitled post'}</Link>
        {post.text ? <p className="post-text">{post.text}</p> : null}
        <div className="card-actions">
          <Link to={`/post/${post.id}`}>{post.commentsCount} comments</Link>
          <span>{post.upVotes} up</span>
          <span>{post.downVotes} down</span>
        </div>
      </div>
    </article>
  );
}

function CommunitiesPage() {
  const [communities, setCommunities] = useState<CommunitySummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    async function run() {
      try {
        setCommunities(await getCommunities());
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unable to load communities');
      } finally {
        setLoading(false);
      }
    }
    void run();
  }, []);

  return (
    <section className="stack narrow">
      <PageHeading eyebrow="Communities" title="Find where posts belong" description="Communities are the strongest information architecture primitive in this app." />
      {error && <Notice tone="error">{error}</Notice>}
      {loading ? <SkeletonList /> : null}
      <div className="community-list">
        {communities.map((community) => (
          <article className="card compact" key={community.id || community.name}>
            <h2>{community.name || 'Unnamed community'}</h2>
            <p>{community.description || 'No description yet.'}</p>
            <div className="card-actions">
              <span>{community.postsCount} posts</span>
              <span>{community.membersCount} members</span>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}

function PostDetailPage() {
  const { postId = '' } = useParams();
  const [post, setPost] = useState<ForumPost | null>(null);
  const [comments, setComments] = useState<ForumComment[]>([]);
  const [page, setPage] = useState<PageInfo>(emptyPage);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  async function loadComments(after = '') {
    const result = await getForumComments(postId, after);
    setComments((current) => (after ? [...current, ...result.comments] : result.comments));
    setPage(result.pagination);
  }

  useEffect(() => {
    async function run() {
      setLoading(true);
      setError('');
      try {
        const [postResult] = await Promise.all([getPost(postId), loadComments()]);
        setPost(postResult);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unable to load post');
      } finally {
        setLoading(false);
      }
    }
    if (postId) void run();
  }, [postId]);

  if (loading) return <SkeletonList />;
  if (error) return <Notice tone="error">{error}</Notice>;
  if (!post) return <NotFound />;

  return (
    <section className="stack narrow">
      <PostCard post={post} />
      <section className="card">
        <h2>Comments</h2>
        {comments.length === 0 ? <EmptyState title="No comments yet" body="The discussion has not started." /> : null}
        <div className="comment-list">
          {comments.map((comment) => <CommentCard key={comment.id} comment={comment} />)}
        </div>
        {page.has_more ? (
          <button className="secondary-action" onClick={() => void loadComments(page.next_cursor)}>Load more comments</button>
        ) : null}
      </section>
    </section>
  );
}

function CommentCard({ comment }: { comment: ForumComment }) {
  return (
    <article className="comment-card">
      <p className="meta">{comment.username || 'unknown'} {comment.edited ? 'edited' : ''}</p>
      <p>{comment.text}</p>
      <div className="card-actions">
        <span>{comment.upVotes - comment.downVotes} score</span>
        <span>{comment.upVotes} up</span>
        <span>{comment.downVotes} down</span>
      </div>
    </article>
  );
}

function ComposePlaceholder() {
  return <EmptyState title="Post creation is next" body="The React shell is live. The blocked auth/write client will be added in the next focused PR so post creation can use the secured backend endpoints." />;
}

function AuthPlaceholder() {
  return <EmptyState title="Auth client is next" body="Login/signup UI was intentionally separated from this read-path migration because the connector blocked the auth-header client commit. The backend auth endpoints remain available." />;
}

function ProfilePlaceholder() {
  return <EmptyState title="Profile depends on auth" body="Profile actions will be re-enabled with the React auth client." />;
}

function StaticPage({ title }: { title: string }) {
  return <EmptyState title={title} body="Policy content should be consolidated into one maintainable content source instead of duplicated page components." />;
}

function NotFound() {
  return <EmptyState title="Page not found" body="Use the primary navigation to return to the feed." />;
}

function PageHeading({ eyebrow, title, description }: { eyebrow: string; title: string; description: string }) {
  return (
    <div className="page-heading">
      <p>{eyebrow}</p>
      <h1>{title}</h1>
      <span>{description}</span>
    </div>
  );
}

function EmptyState({ title, body }: { title: string; body: string }) {
  return (
    <section className="empty-state">
      <h1>{title}</h1>
      <p>{body}</p>
      <Link className="secondary-action" to="/home">Back to feed</Link>
    </section>
  );
}

function Notice({ children, tone }: { children: string; tone: 'error' | 'info' }) {
  return <div className={`notice ${tone}`}>{children}</div>;
}

function SkeletonList() {
  return (
    <div className="stack">
      <div className="skeleton" />
      <div className="skeleton" />
      <div className="skeleton" />
    </div>
  );
}

function formatDate(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return date.toLocaleDateString();
}
