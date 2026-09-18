import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
export const Destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

Destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/role/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
Destroy.url = (id: string, options?: RouteQueryOptions) => {
  return Destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
Destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
export const DestroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:227
 * @route /api/v1/role/:id
 */
DestroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
export const Get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

Get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/role/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
Get.url = (id: string, options?: RouteQueryOptions) => {
  return Get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
Get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
Get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Get.url(id, options),
  method: 'head',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
export const GetForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
GetForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:74
 * @route /api/v1/role/:id
 */
GetForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Get.url(id, options),
  method: 'head',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
export const List = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

List.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/role/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
List.url = (options?: RouteQueryOptions) => {
  return List.definition.url + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
List.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
List.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: List.url(options),
  method: 'head',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
export const ListForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
ListForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:42
 * @route /api/v1/role/
 */
ListForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: List.url(options),
  method: 'head',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
export const Store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

Store.definition = {
  methods: ['post'],
  url: '/api/v1/role/',
} satisfies RouteDefinition<['post']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
Store.url = (options?: RouteQueryOptions) => {
  return Store.definition.url + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
Store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
export const StoreForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:104
 * @route /api/v1/role/
 */
StoreForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
export const Update = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: Update.url(id, options),
  method: 'put',
});

Update.definition = {
  methods: ['put'],
  url: '/api/v1/role/:id',
} satisfies RouteDefinition<['put']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
Update.url = (id: string, options?: RouteQueryOptions) => {
  return Update.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
Update.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: Update.url(id, options),
  method: 'put',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
export const UpdateForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: Update.url(id, options),
  method: 'put',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:144
 * @route /api/v1/role/:id
 */
UpdateForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: Update.url(id, options),
  method: 'put',
});
/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
export const UpdatePermissions = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: UpdatePermissions.url(id, options),
  method: 'put',
});

UpdatePermissions.definition = {
  methods: ['put'],
  url: '/api/v1/role/:id/permissions',
} satisfies RouteDefinition<['put']>;

/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
UpdatePermissions.url = (id: string, options?: RouteQueryOptions) => {
  return UpdatePermissions.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
UpdatePermissions.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: UpdatePermissions.url(id, options),
  method: 'put',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
export const UpdatePermissionsForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: UpdatePermissions.url(id, options),
  method: 'put',
});

/**
 * @see role/internal/app/http/controllers/role/role.go:186
 * @route /api/v1/role/:id/permissions
 */
UpdatePermissionsForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: UpdatePermissions.url(id, options),
  method: 'put',
});

export const RoleController = {
  Destroy: Object.assign(Destroy, Destroy),
  Get: Object.assign(Get, Get),
  List: Object.assign(List, List),
  Store: Object.assign(Store, Store),
  Update: Object.assign(Update, Update),
  UpdatePermissions: Object.assign(UpdatePermissions, UpdatePermissions),
};

export default RoleController;
