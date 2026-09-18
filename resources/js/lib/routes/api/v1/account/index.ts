import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
export const destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/account/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
destroy.url = (id: string, options?: RouteQueryOptions) => {
  return destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
export const destroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
destroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
export const get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/account/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
get.url = (id: string, options?: RouteQueryOptions) => {
  return get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: get.url(id, options),
  method: 'head',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
export const getForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
getForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
getForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: get.url(id, options),
  method: 'head',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
export const list = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

list.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/account/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
list.url = (options?: RouteQueryOptions) => {
  return list.definition.url + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
list.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
list.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: list.url(options),
  method: 'head',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
export const listForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
listForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
listForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: list.url(options),
  method: 'head',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
export const store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

store.definition = {
  methods: ['post'],
  url: '/api/v1/account/',
} satisfies RouteDefinition<['post']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
store.url = (options?: RouteQueryOptions) => {
  return store.definition.url + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
export const storeForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
storeForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
export const update = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update.url(id, options),
  method: 'put',
});

update.definition = {
  methods: ['put'],
  url: '/api/v1/account/:id',
} satisfies RouteDefinition<['put']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
update.url = (id: string, options?: RouteQueryOptions) => {
  return update.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
update.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update.url(id, options),
  method: 'put',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
export const updateForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update.url(id, options),
  method: 'put',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
updateForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update.url(id, options),
  method: 'put',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
export const update_password = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update_password.url(id, options),
  method: 'put',
});

update_password.definition = {
  methods: ['put'],
  url: '/api/v1/account/:id/password',
} satisfies RouteDefinition<['put']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
update_password.url = (id: string, options?: RouteQueryOptions) => {
  return update_password.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
update_password.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update_password.url(id, options),
  method: 'put',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
export const update_passwordForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update_password.url(id, options),
  method: 'put',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
update_passwordForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update_password.url(id, options),
  method: 'put',
});

export const api_v1_account = {
  destroy: Object.assign(destroy, destroy),
  get: Object.assign(get, get),
  list: Object.assign(list, list),
  store: Object.assign(store, store),
  update: Object.assign(update, update),
  update_password: Object.assign(update_password, update_password),
};

export default api_v1_account;
