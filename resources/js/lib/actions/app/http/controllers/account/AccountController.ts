import { type RouteDefinition, type RouteFormDefinition, type RouteQueryOptions, queryParams } from '@/js/lib/wayfinder';

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
export const Destroy = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

Destroy.definition = {
  methods: ['delete'],
  url: '/api/v1/account/:id',
} satisfies RouteDefinition<['delete']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
Destroy.url = (id: string, options?: RouteQueryOptions) => {
  return Destroy.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
Destroy.delete = (id: string, options?: RouteQueryOptions): RouteDefinition<'delete'> => ({
  url: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
export const DestroyForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:257
 * @route /api/v1/account/:id
 */
DestroyForm.delete = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'delete'> => ({
  action: Destroy.url(id, options),
  method: 'delete',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
export const Get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

Get.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/account/:id',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
Get.url = (id: string, options?: RouteQueryOptions) => {
  return Get.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
Get.get = (id: string, options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: Get.url(id, options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
Get.head = (id: string, options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: Get.url(id, options),
  method: 'head',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
export const GetForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
GetForm.get = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: Get.url(id, options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:78
 * @route /api/v1/account/:id
 */
GetForm.head = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: Get.url(id, options),
  method: 'head',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
export const List = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

List.definition = {
  methods: ['get', 'head'],
  url: '/api/v1/account/',
} satisfies RouteDefinition<['get', 'head']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
List.url = (options?: RouteQueryOptions) => {
  return List.definition.url + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
List.get = (options?: RouteQueryOptions): RouteDefinition<'get'> => ({
  url: List.url(options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
List.head = (options?: RouteQueryOptions): RouteDefinition<'head'> => ({
  url: List.url(options),
  method: 'head',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
export const ListForm = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
ListForm.get = (options?: RouteQueryOptions): RouteFormDefinition<'get'> => ({
  action: List.url(options),
  method: 'get',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:46
 * @route /api/v1/account/
 */
ListForm.head = (options?: RouteQueryOptions): RouteFormDefinition<'head'> => ({
  action: List.url(options),
  method: 'head',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
export const Store = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

Store.definition = {
  methods: ['post'],
  url: '/api/v1/account/',
} satisfies RouteDefinition<['post']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
Store.url = (options?: RouteQueryOptions) => {
  return Store.definition.url + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
Store.post = (options?: RouteQueryOptions): RouteDefinition<'post'> => ({
  url: Store.url(options),
  method: 'post',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
export const StoreForm = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:108
 * @route /api/v1/account/
 */
StoreForm.post = (options?: RouteQueryOptions): RouteFormDefinition<'post'> => ({
  action: Store.url(options),
  method: 'post',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
export const Update = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: Update.url(id, options),
  method: 'put',
});

Update.definition = {
  methods: ['put'],
  url: '/api/v1/account/:id',
} satisfies RouteDefinition<['put']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
Update.url = (id: string, options?: RouteQueryOptions) => {
  return Update.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
Update.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: Update.url(id, options),
  method: 'put',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
export const UpdateForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: Update.url(id, options),
  method: 'put',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:159
 * @route /api/v1/account/:id
 */
UpdateForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: Update.url(id, options),
  method: 'put',
});
/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
export const UpdatePassword = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: UpdatePassword.url(id, options),
  method: 'put',
});

UpdatePassword.definition = {
  methods: ['put'],
  url: '/api/v1/account/:id/password',
} satisfies RouteDefinition<['put']>;

/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
UpdatePassword.url = (id: string, options?: RouteQueryOptions) => {
  return UpdatePassword.definition.url.replace(':id', encodeURIComponent(id)) + queryParams(options);
};

/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
UpdatePassword.put = (id: string, options?: RouteQueryOptions): RouteDefinition<'put'> => ({
  url: UpdatePassword.url(id, options),
  method: 'put',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
export const UpdatePasswordForm = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: UpdatePassword.url(id, options),
  method: 'put',
});

/**
 * @see account/internal/app/http/controllers/account/account.go:208
 * @route /api/v1/account/:id/password
 */
UpdatePasswordForm.put = (id: string, options?: RouteQueryOptions): RouteFormDefinition<'put'> => ({
  action: UpdatePassword.url(id, options),
  method: 'put',
});

export const AccountController = {
  Destroy: Object.assign(Destroy, Destroy),
  Get: Object.assign(Get, Get),
  List: Object.assign(List, List),
  Store: Object.assign(Store, Store),
  Update: Object.assign(Update, Update),
  UpdatePassword: Object.assign(UpdatePassword, UpdatePassword),
};

export default AccountController;
