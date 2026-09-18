import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
export const destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/role/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
destroy.url = (id: string, options?: RouteQueryOptions) => {
  return destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
export const destroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
destroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: destroy.url(id, options),
  method: 'delete',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
export const get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/role/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
get.url = (id: string, options?: RouteQueryOptions) => {
  return get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: get.url(id, options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: get.url(id, options),
  method: 'head',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
export const getForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
getForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: get.url(id, options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
getForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: get.url(id, options),
  method: 'head',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
export const list = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

list.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/role/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
list.url = (options?: RouteQueryOptions) => {
  return list.definition.url + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
list.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: list.url(options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
list.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: list.url(options),
  method: 'head',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
export const listForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
listForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: list.url(options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
listForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: list.url(options),
  method: 'head',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
export const store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

store.definition = {
  methods: ['post'],
  url: '/api/v1/role/',
} satisfies RouteDefinition<['post']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
store.url = (options?: RouteQueryOptions) => {
  return store.definition.url + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: store.url(options),
  method: 'post',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
export const storeForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
storeForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: store.url(options),
  method: 'post',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
export const update = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update.url(id, options),
  method: 'put',
});

update.definition = {
  methods: ['put'],
  url: '/api/v1/role/:id',
} satisfies RouteDefinition<['put']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
update.url = (id: string, options?: RouteQueryOptions) => {
  return update.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
update.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update.url(id, options),
  method: 'put',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
export const updateForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update.url(id, options),
  method: 'put',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
updateForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update.url(id, options),
  method: 'put',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
export const update_permissions = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update_permissions.url(id, options),
  method: 'put',
});

update_permissions.definition = {
  methods: ['put'],
  url: '/api/v1/role/:id/permissions',
} satisfies RouteDefinition<['put']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
update_permissions.url = (id: string, options?: RouteQueryOptions) => {
  return update_permissions.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
update_permissions.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: update_permissions.url(id, options),
  method: 'put',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
export const update_permissionsForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update_permissions.url(id, options),
  method: 'put',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
update_permissionsForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: update_permissions.url(id, options),
  method: 'put',
});

export const api_v1_role = {
  destroy: Object.assign(destroy, destroy),
  get: Object.assign(get, get),
  list: Object.assign(list, list),
  store: Object.assign(store, store),
  update: Object.assign(update, update),
  update_permissions: Object.assign(update_permissions, update_permissions),
};

export default api_v1_role;
