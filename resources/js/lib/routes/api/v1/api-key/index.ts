import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
export const destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/api-key/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
destroy.url = (id: string, options?: RouteQueryOptions) => {
  return destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
export const destroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:205
 * @route /api/v1/api-key/:id
 */
destroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});
/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
export const get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/api-key/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
get.url = (id: string, options?: RouteQueryOptions) => {
  return get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: get.url(id, options),
  method: 'head',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
export const getForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
getForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:77
 * @route /api/v1/api-key/:id
 */
getForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: get.url(id, options),
  method: 'head',
});
/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
export const list = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

list.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/api-key/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
list.url = (options?: RouteQueryOptions) => {
  return list.definition.url + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
list.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
list.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: list.url(options),
  method: 'head',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
export const listForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
listForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:46
 * @route /api/v1/api-key/
 */
listForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: list.url(options),
  method: 'head',
});
/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
export const store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

store.definition = {
  methods: ['post'],
  url: '/api/v1/api-key/',
} satisfies RouteDefinition<['post']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
store.url = (options?: RouteQueryOptions) => {
  return store.definition.url + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
export const storeForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:107
 * @route /api/v1/api-key/
 */
storeForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});
/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
export const update = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update.url(id, options),
  method: 'put',
});

update.definition = {
  methods: ['put'],
  url: '/api/v1/api-key/:id',
} satisfies RouteDefinition<['put']>;

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
update.url = (id: string, options?: RouteQueryOptions) => {
  return update.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
update.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update.url(id, options),
  method: 'put',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
export const updateForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update.url(id, options),
  method: 'put',
});

/**
 * @see apikey/internal/app/http/controllers/apikey/apikey.go:160
 * @route /api/v1/api-key/:id
 */
updateForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update.url(id, options),
  method: 'put',
});

export const api_v1_api_key = {
  destroy: Object.assign(destroy, destroy),
  get: Object.assign(get, get),
  list: Object.assign(list, list),
  store: Object.assign(store, store),
  update: Object.assign(update, update),
};

export default api_v1_api_key;
