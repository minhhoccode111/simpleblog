# Imagine the flow

**NOTE**: Forms can't send `PUT` and `DELETE` request

- `GET /`:
  - redirect to `/articles`
- `GET /articles`
  - list all published articles
- `GET /articles/:slug`
  - read an article (only if it's published)
- `GET /admin`
  - basic authentication
  - redirect to `/admin/articles`
- `GET /admin/articles`
  - list all articles
- `GET /admin/articles?action=create&slug=some-slug`
  - display the form editor
- `POST /admin/articles`
  - save new article
- `GET /admin/articles/:slug?action=edit`
  - display the form editor fill with existed data
- `POST /admin/articles/:slug?action=edit`
  - save the edited data
- `POST /admin/articles/:slug?action=delete`
  - delete an article
