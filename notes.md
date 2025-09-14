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
  - display an input to create new article with a title
- `POST /admin/articles`
  - save new article title
- `GET /admin/articles/:slug`
  - display the form editor fill with existed data
- `POST /admin/articles/:slug?action=edit`
  - save the edited data
- `POST /admin/articles/:slug?action=delete`
  - delete an article
